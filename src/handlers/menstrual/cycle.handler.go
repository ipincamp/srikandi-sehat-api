package menstrual

import (
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	menstrual "ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
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
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
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

		startDate, _ := time.Parse("2006-01-02", input.StartDate)

		var lastCompletedCycle menstrual.MenstrualCycle
		errLastCompleted := tx.Where("user_id = ? AND end_date IS NOT NULL", user.ID).Order("end_date desc").First(&lastCompletedCycle).Error
		if errLastCompleted == nil {
			if !startDate.After(lastCompletedCycle.EndDate.Time) {
				return utils.SendError(c, fiber.StatusConflict, "The new cycle's start date cannot overlap with the previous cycle.")
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

		endDate, _ := time.Parse("2006-01-02", input.EndDate)

		if endDate.Before(activeCycle.StartDate) {
			return utils.SendError(c, fiber.StatusBadRequest, "The end date cannot be before the start date of the current cycle.")
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

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	var cycles []menstrual.MenstrualCycle
	database.DB.Where("user_id = ?", user.ID).Order("start_date desc").Find(&cycles)

	var responseData []dto.CycleResponse
	for _, cycle := range cycles {
		responseData = append(responseData, dto.CycleResponse{
			ID:             cycle.ID,
			StartDate:      cycle.StartDate,
			EndDate:        cycle.EndDate,
			PeriodLength:   cycle.PeriodLength,
			CycleLength:    cycle.CycleLength,
			IsPeriodNormal: cycle.IsPeriodNormal,
			IsCycleNormal:  cycle.IsCycleNormal,
		})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Cycle history fetched successfully", responseData)
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
	periodLength := int16(endDate.Sub(currentCycle.StartDate).Hours()/24) + 1
	isNormal := periodLength >= 2 && periodLength <= 7

	tx.Model(currentCycle).Updates(map[string]interface{}{
		"end_date":         endDate,
		"period_length":    periodLength,
		"is_period_normal": isNormal,
	})
}
