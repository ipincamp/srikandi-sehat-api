package menstrual

import (
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	menstrual "ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RecordCycle(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	input := c.Locals("request_body").(*dto.CycleRequest)

	if input == nil || (input.StartDate == "" && input.EndDate == "") {
		return utils.SendError(c, fiber.StatusBadRequest, "StartDate or EndDate must be provided")
	}

	var user models.User
	if err := database.DB.Preload("Profile").First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	log.Print(user.Profile.ID)

	if user.Profile.ID >= 1 {
		// Profile exists, proceed without error
	} else {
		return utils.SendError(c, fiber.StatusForbidden, "Please complete your profile before recording a cycle. Essential data is missing.")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	activeCycle, err := findActiveCycle(tx, user.ID)
	var isStartRequest bool

	if input.StartDate != "" {
		if err == nil {
			return utils.SendError(c, fiber.StatusConflict, "Cannot start a new cycle while another is in progress.")
		}

		startDate, parseErr := time.Parse(time.RFC3339, input.StartDate)
		if parseErr != nil {
			return utils.SendError(c, fiber.StatusBadRequest, "Invalid StartDate format")
		}

		var lastCompletedCycle menstrual.MenstrualCycle
		errLastCompleted := tx.Where("user_id = ? AND end_date IS NOT NULL", user.ID).Order("end_date desc").First(&lastCompletedCycle).Error
		if errLastCompleted == nil {
			if !startDate.After(lastCompletedCycle.EndDate.Time) {
				formattedDate := lastCompletedCycle.EndDate.Time.Format("2 January 2006 15:04:05")
				errorMessage := fmt.Sprintf("Start date cannot be before or equal to the end date of the last completed cycle (%s).", formattedDate)
				return utils.SendError(c, fiber.StatusConflict, errorMessage)
			}
		}

		newCycle := menstrual.MenstrualCycle{UserID: user.ID, StartDate: startDate}
		if err := tx.Create(&newCycle).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to record new cycle")
		}

		updatePreviousCycleLength(tx, user.ID, startDate)
		isStartRequest = true
	}

	if input.EndDate != "" {
		if err != nil {
			return utils.SendError(c, fiber.StatusConflict, "No active cycle to end. Please record a new cycle first.")
		}

		endDate, parseErr := time.Parse(time.RFC3339, input.EndDate)
		if parseErr != nil {
			return utils.SendError(c, fiber.StatusBadRequest, "Invalid EndDate format")
		}

		if endDate.Before(activeCycle.StartDate) {
			formattedDate := activeCycle.StartDate.Format("2 January 2006 15:04:05")
			errorMessage := fmt.Sprintf("Finish date cannot be before the start date of the current cycle (%s).", formattedDate)
			return utils.SendError(c, fiber.StatusBadRequest, errorMessage)
		}

		updateCurrentCyclePeriod(tx, &activeCycle, endDate)
		isStartRequest = false
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	if isStartRequest {
		return utils.SendSuccess(c, fiber.StatusOK, "Cycle started successfully", nil)
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Cycle ended successfully", nil)
}

func GetCycleHistory(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	queries := c.Locals("request_queries").(*dto.PaginationQuery)

	page := queries.Page
	if page <= 0 {
		page = 1
	}
	limit := queries.Limit
	if limit <= 0 {
		limit = 10
	}

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	baseQuery := database.DB.Model(&menstrual.MenstrualCycle{}).Where("user_id = ? AND end_date IS NOT NULL", user.ID)

	pagination, paginateScope := utils.GeneratePagination(page, limit, baseQuery, &menstrual.MenstrualCycle{})

	if pagination.TotalRows == 0 {
		return utils.SendSuccess(c, fiber.StatusOK, "You have no cycle history. Please record a cycle first.", []dto.CycleResponse{})
	}

	var cycles []menstrual.MenstrualCycle
	err := baseQuery.Scopes(paginateScope).Order("start_date desc").Find(&cycles).Error
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch cycle history for user %s: %v", userUUID, err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve cycle history")
	}

	var responseData []dto.CycleResponse
	for _, cycle := range cycles {
		dto := dto.CycleResponse{
			ID:        cycle.ID,
			StartDate: cycle.StartDate,
		}

		if cycle.EndDate.Valid {
			dto.EndDate = &cycle.EndDate.Time
		}
		if cycle.PeriodLength.Valid {
			dto.PeriodLength = &cycle.PeriodLength.Int16
		}
		if cycle.CycleLength.Valid {
			dto.CycleLength = &cycle.CycleLength.Int16
		}
		if cycle.IsPeriodNormal.Valid {
			dto.IsPeriodNormal = &cycle.IsPeriodNormal.Bool
		}
		if cycle.IsCycleNormal.Valid {
			dto.IsCycleNormal = &cycle.IsCycleNormal.Bool
		}
		responseData = append(responseData, dto)
	}

	paginatedResponse := dto.PaginatedResponse[dto.CycleResponse]{
		Data:     responseData,
		Metadata: pagination,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Cycle history fetched successfully", paginatedResponse)
}

func GetCycleByID(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	params := c.Locals("request_params").(*dto.CycleParam)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	var cycle menstrual.MenstrualCycle
	if err := database.DB.Where("id = ? AND user_id = ?", params.ID, user.ID).First(&cycle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "Cycle not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve cycle data")
	}

	var cycleEndDate time.Time
	if cycle.EndDate.Valid {
		cycleEndDate = cycle.EndDate.Time
	} else {
		cycleEndDate = time.Now()
	}

	var symptomLogs []menstrual.SymptomLog
	err := database.DB.
		Preload("Details.Symptom").
		Preload("Details.SymptomOption").
		Where("user_id = ? AND logged_at >= ? AND logged_at <= ?", user.ID, cycle.StartDate, cycleEndDate).
		Order("logged_at desc").
		Find(&symptomLogs).Error

	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve symptom data")
	}

	var symptomGroups []dto.SymptomLogGroupResponse
	for _, log := range symptomLogs {
		var details []dto.SymptomDetail
		for _, detail := range log.Details {
			symptomDetail := dto.SymptomDetail{
				SymptomName:     detail.Symptom.Name,
				SymptomCategory: detail.Symptom.Category,
			}
			if detail.SymptomOptionID.Valid && detail.SymptomOption.Name != "" {
				symptomDetail.SelectedOption = &detail.SymptomOption.Name
			}
			details = append(details, symptomDetail)
		}

		group := dto.SymptomLogGroupResponse{
			ID:       log.ID,
			LoggedAt: log.LoggedAt,
			Details:  details,
		}
		if log.Note != "" {
			group.Note = &log.Note
		}
		symptomGroups = append(symptomGroups, group)
	}

	response := dto.CycleDetailResponse{
		ID:        cycle.ID,
		StartDate: cycle.StartDate,
		Symptoms:  symptomGroups,
	}

	if cycle.EndDate.Valid {
		response.EndDate = &cycle.EndDate.Time
	}
	if cycle.PeriodLength.Valid {
		response.PeriodLength = &cycle.PeriodLength.Int16
	}
	if cycle.CycleLength.Valid {
		response.CycleLength = &cycle.CycleLength.Int16
	}
	if cycle.IsPeriodNormal.Valid {
		response.IsPeriodNormal = &cycle.IsPeriodNormal.Bool
	}
	if cycle.IsCycleNormal.Valid {
		response.IsCycleNormal = &cycle.IsCycleNormal.Bool
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Cycle detail fetched successfully", response)
}

func GetCycleStatus(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	var activeCycle menstrual.MenstrualCycle
	err := database.DB.Where("user_id = ? AND end_date IS NULL", user.ID).Order("start_date desc").First(&activeCycle).Error

	// Case 1: User is currently in a cycle.
	if err == nil {
		today := time.Now()
		currentPeriodDay := int(today.Sub(activeCycle.StartDate).Hours()/24) + 1
		isPeriodNormal := currentPeriodDay >= 2 && currentPeriodDay <= 7

		response := dto.CycleStatusResponse{
			IsOnCycle:        true,
			CurrentPeriodDay: &currentPeriodDay,
			IsPeriodNormal:   &isPeriodNormal,
			Message:          fmt.Sprintf("Anda sedang berada di hari ke-%d siklus menstruasi.", currentPeriodDay),
		}

		var previousCycle menstrual.MenstrualCycle
		errPrev := database.DB.Where("user_id = ? AND end_date IS NOT NULL", user.ID).Order("start_date desc").First(&previousCycle).Error
		if errPrev == nil {
			currentCycleLength := int(activeCycle.StartDate.Sub(previousCycle.StartDate).Hours() / 24)
			isCycleNormal := currentCycleLength >= 21 && currentCycleLength <= 35
			response.CurrentCycleLength = &currentCycleLength
			response.IsCycleNormal = &isCycleNormal
		}

		return utils.SendSuccess(c, fiber.StatusOK, "Cycle status fetched.", response)
	}

	// Case 2: User is not currently in a cycle.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var completedCycles []menstrual.MenstrualCycle
		database.DB.Where("user_id = ? AND end_date IS NOT NULL", user.ID).Order("start_date desc").Limit(6).Find(&completedCycles)

		if len(completedCycles) == 0 {
			return utils.SendSuccess(c, fiber.StatusOK, "No cycle data available.", dto.CycleStatusResponse{
				IsOnCycle: false,
				Message:   "Belum ada data siklus ditemukan. Silakan catat siklus menstruasi Anda untuk melihat status dan prediksi.",
			})
		}

		response := dto.CycleStatusResponse{
			IsOnCycle: false,
			Message:   "Anda tidak sedang dalam siklus menstruasi saat ini.",
		}

		// Add last period length
		lastCompletedCycle := completedCycles[0]
		if lastCompletedCycle.PeriodLength.Valid {
			lastPeriodLen := int(lastCompletedCycle.PeriodLength.Int16)
			response.LastPeriodLength = &lastPeriodLen
		}

		// Calculate last cycle length if possible
		if len(completedCycles) >= 2 {
			lastCycleLength := int(completedCycles[0].StartDate.Sub(completedCycles[1].StartDate).Hours() / 24)
			isCycleNormal := lastCycleLength >= 21 && lastCycleLength <= 35
			response.LastCycleLength = &lastCycleLength
			response.IsCycleNormal = &isCycleNormal
		}

		// Predict next period
		var totalCycleLength int
		var validCyclesForAvg int
		for _, cycle := range completedCycles {
			if cycle.CycleLength.Valid {
				totalCycleLength += int(cycle.CycleLength.Int16)
				validCyclesForAvg++
			}
		}

		if validCyclesForAvg > 0 {
			averageCycleLength := totalCycleLength / validCyclesForAvg
			lastStartDate := completedCycles[0].StartDate
			predictedDate := lastStartDate.AddDate(0, 0, averageCycleLength)
			daysUntil := int(time.Until(predictedDate).Hours() / 24)

			if daysUntil >= 0 {
				predictedDateStr := predictedDate.Format("2006-01-02")
				response.DaysUntilNextPeriod = &daysUntil
				response.PredictedPeriodDate = &predictedDateStr
				response.Message = fmt.Sprintf("Periode menstruasi Anda berikutnya diprediksi dalam %d hari.", daysUntil)
			} else {
				response.Message = "Tanggal prediksi menstruasi Anda telah lewat. Silakan catat siklus baru jika sudah dimulai."
			}
		} else if response.LastCycleLength == nil {
			response.Message = "Data belum cukup untuk memprediksi periode berikutnya. Silakan catat minimal satu siklus menstruasi lengkap."
		}

		return utils.SendSuccess(c, fiber.StatusOK, "Cycle status fetched.", response)
	}

	return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve cycle data")
}

func findActiveCycle(tx *gorm.DB, userID uint) (menstrual.MenstrualCycle, error) {
	var activeCycle menstrual.MenstrualCycle
	err := tx.Where("user_id = ? AND end_date IS NULL", userID).
		Order("start_date desc").
		First(&activeCycle).Error
	return activeCycle, err
}

func updatePreviousCycleLength(tx *gorm.DB, userID uint, newStartDate time.Time) {
	var previousCycle menstrual.MenstrualCycle
	err := tx.Where("user_id = ? AND start_date < ?", userID, newStartDate).
		Order("start_date desc").
		First(&previousCycle).Error

	if err == nil {
		cycleLength := int16(newStartDate.Sub(previousCycle.StartDate).Hours() / 24)
		isNormal := cycleLength >= 21 && cycleLength <= 35

		tx.Model(&previousCycle).Updates(map[string]interface{}{
			"cycle_length":    cycleLength,
			"is_cycle_normal": isNormal,
		})
	}
}

func updateCurrentCyclePeriod(tx *gorm.DB, currentCycle *menstrual.MenstrualCycle, endDate time.Time) {
	loc := time.Local

	startDay := time.Date(currentCycle.StartDate.Year(), currentCycle.StartDate.Month(), currentCycle.StartDate.Day(), 0, 0, 0, 0, loc)
	endDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, loc)

	periodLength := int16(endDay.Sub(startDay).Hours()/24) + 1
	isNormal := periodLength >= 2 && periodLength <= 7

	tx.Model(currentCycle).Updates(map[string]interface{}{
		"end_date":         endDate,
		"period_length":    periodLength,
		"is_period_normal": isNormal,
	})
}
