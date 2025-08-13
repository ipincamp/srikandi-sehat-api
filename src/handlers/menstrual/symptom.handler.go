package menstrual

import (
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetSymptomsMaster(c *fiber.Ctx) error {
	var symptoms []menstrual.Symptom
	if err := database.DB.
		Select("id, name, type").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, symptom_id")
		}).
		Find(&symptoms).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch symptoms data")
	}

	var responseData []dto.SymptomMasterResponse

	for _, s := range symptoms {
		var optionDTOs []dto.SymptomOptionResponse
		for _, o := range s.Options {
			optionDTOs = append(optionDTOs, dto.SymptomOptionResponse{
				ID:   o.ID,
				Name: o.Name,
			})
		}

		symptomDTO := dto.SymptomMasterResponse{
			ID:      s.ID,
			Name:    s.Name,
			Type:    s.Type,
			Options: optionDTOs,
		}
		responseData = append(responseData, symptomDTO)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Symptoms master data fetched successfully", responseData)
}

func LogSymptoms(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	input := c.Locals("request_body").(*dto.SymptomLogRequest)

	var symptoms []menstrual.Symptom
	if err := database.DB.
		Select("id, name, type").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, symptom_id")
		}).
		Find(&symptoms).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch symptoms data")
	}

	for _, inputSymptom := range input.Symptoms {
		var dbSymptom *menstrual.Symptom
		for i := range symptoms {
			if inputSymptom.SymptomID == symptoms[i].ID {
				dbSymptom = &symptoms[i]
				break
			}
		}
		if dbSymptom == nil {
			return utils.SendError(c, fiber.StatusBadRequest, "Invalid symptom ID: "+fmt.Sprintf("%d", inputSymptom.SymptomID))
		}
		if inputSymptom.SymptomOptionID != nil {
			foundOption := false
			for _, opt := range dbSymptom.Options {
				if opt.ID == *inputSymptom.SymptomOptionID {
					foundOption = true
					break
				}
			}
			if !foundOption {
				return utils.SendError(c, fiber.StatusBadRequest, "Invalid symptom option ID: "+fmt.Sprintf("%d", *inputSymptom.SymptomOptionID))
			}
		}
	}

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	loggedAt, _ := time.Parse(time.RFC3339, input.LoggedAt)

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	symptomLog := menstrual.SymptomLog{
		UserID:   user.ID,
		LoggedAt: loggedAt,
		Note:     input.Note,
	}
	if err := tx.Create(&symptomLog).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create symptom log")
	}

	var relevantCycle menstrual.MenstrualCycle
	err := tx.Where("user_id = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)", user.ID, loggedAt, loggedAt).
		Order("start_date desc").
		First(&relevantCycle).Error

	if err == nil {
		tx.Model(&symptomLog).Update("menstrual_cycle_id", relevantCycle.ID)
	}

	for _, s := range input.Symptoms {
		detail := menstrual.SymptomLogDetail{
			SymptomLogID: symptomLog.ID,
			SymptomID:    s.SymptomID,
		}
		if s.SymptomOptionID != nil {
			detail.SymptomOptionID.Int64 = int64(*s.SymptomOptionID)
			detail.SymptomOptionID.Valid = true
		}
		if err := tx.Create(&detail).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to log symptom detail")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	responseData := dto.SymptomLogCreateResponse{
		ID: symptomLog.ID,
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Symptoms logged successfully", responseData)
}

func GetSymptomHistory(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	queries := c.Locals("request_queries").(*dto.SymptomHistoryQuery)

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

	baseQuery := database.DB.Model(&menstrual.SymptomLog{}).Where("user_id = ?", user.ID)

	if queries.Date != "" {
		startOfDay, _ := time.Parse("2006-01-02", queries.Date)
		endOfDay := startOfDay.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		baseQuery = baseQuery.Where("logged_at BETWEEN ? AND ?", startOfDay, endOfDay)
	} else if queries.StartDate != "" && queries.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", queries.EndDate)
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		baseQuery = baseQuery.Where("logged_at BETWEEN ? AND ?", queries.StartDate, endDate)
	}

	pagination, paginateScope := utils.GeneratePagination(page, limit, baseQuery, &menstrual.SymptomLog{})

	if pagination.TotalRows == 0 {
		return utils.SendSuccess(c, fiber.StatusOK, "No symptom history found", dto.PaginatedResponse[dto.SymptomHistoryResponse]{
			Data:     []dto.SymptomHistoryResponse{},
			Metadata: pagination,
		})
	}

	var logs []menstrual.SymptomLog
	err := baseQuery.Preload("Details").Scopes(paginateScope).Order("logged_at desc").Find(&logs).Error
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve symptom history")
	}

	var results []dto.SymptomHistoryResponse
	for _, log := range logs {
		results = append(results, dto.SymptomHistoryResponse{
			ID:            log.ID,
			TotalSymptoms: len(log.Details),
			LoggedAt:      log.LoggedAt,
			// Note:          log.Note,
		})
	}

	paginatedResponse := dto.PaginatedResponse[dto.SymptomHistoryResponse]{
		Data:     results,
		Metadata: pagination,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Symptom history fetched successfully", paginatedResponse)
}

func GetSymptomLogByID(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	params := c.Locals("request_params").(*dto.SymptomLogParam)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	var symptomLog menstrual.SymptomLog
	err := database.DB.
		Preload("Details.Symptom").
		Preload("Details.SymptomOption").
		Where("id = ? AND user_id = ?", params.ID, user.ID).
		First(&symptomLog).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "Symptom log not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve symptom log")
	}

	var detailsDTO []dto.SymptomLogDetailResponse
	var symptomIDs []uint
	for _, detail := range symptomLog.Details {
		detailsDTO = append(detailsDTO, dto.SymptomLogDetailResponse{
			SymptomID:       detail.SymptomID,
			SymptomName:     detail.Symptom.Name,
			SymptomCategory: detail.Symptom.Category,
			SelectedOption:  detail.SymptomOption.Name,
		})
		symptomIDs = append(symptomIDs, detail.SymptomID)
	}

	var recommendations []menstrual.Recommendation
	if len(symptomIDs) > 0 {
		database.DB.
			Preload("Symptom").
			Where("symptom_id IN ?", symptomIDs).
			Find(&recommendations)
	}
	var recommendationsDTO []dto.RecommendationResponse
	for _, r := range recommendations {
		recommendationsDTO = append(recommendationsDTO, dto.RecommendationResponse{
			ForSymptom:  r.Symptom.Name,
			Title:       r.Title,
			Description: r.Description,
			Source:      r.Source,
		})
	}

	var cycleNumber *int64
	if symptomLog.MenstrualCycleID.Valid {
		var associatedCycle menstrual.MenstrualCycle
		if err := database.DB.First(&associatedCycle, symptomLog.MenstrualCycleID.Int64).Error; err == nil {
			var count int64
			database.DB.Model(&menstrual.MenstrualCycle{}).
				Where("user_id = ? AND start_date <= ?", user.ID, associatedCycle.StartDate).
				Count(&count)
			cycleNumber = &count
		}
	}

	response := dto.SymptomLogDetailViewResponse{
		LoggedAt:        symptomLog.LoggedAt,
		Note:            symptomLog.Note,
		CycleNumber:     cycleNumber,
		Details:         detailsDTO,
		Recommendations: recommendationsDTO,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Symptom log detail fetched successfully", response)
}

func GetRecommendationsBySymptoms(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	recentDate := time.Now().AddDate(0, 0, -30)

	type SymptomFrequency struct {
		SymptomID uint
		Frequency int
	}
	var frequentSymptoms []SymptomFrequency

	err := database.DB.Model(&menstrual.SymptomLogDetail{}).
		Select("symptom_id, COUNT(symptom_id) as frequency").
		Joins("JOIN symptom_logs ON symptom_logs.id = symptom_log_details.symptom_log_id").
		Where("symptom_logs.user_id = ? AND symptom_logs.logged_at >= ?", user.ID, recentDate).
		Group("symptom_id").
		Order("frequency DESC").
		Limit(4).
		Scan(&frequentSymptoms).Error

	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to analyze recent symptoms")
	}

	if len(frequentSymptoms) == 0 {
		return utils.SendSuccess(c, fiber.StatusOK, "No recent symptoms found to generate recommendations", []dto.RecommendationResponse{})
	}

	var topSymptomIDs []uint
	for _, s := range frequentSymptoms {
		topSymptomIDs = append(topSymptomIDs, s.SymptomID)
	}

	var recommendations []menstrual.Recommendation
	if err := database.DB.
		Preload("Symptom").
		Where("symptom_id IN ?", topSymptomIDs).
		Find(&recommendations).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch recommendations")
	}

	var responseData []dto.RecommendationResponse
	for _, r := range recommendations {
		responseData = append(responseData, dto.RecommendationResponse{
			ForSymptom:  r.Symptom.Name,
			Title:       r.Title,
			Description: r.Description,
			Source:      r.Source,
		})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Recommendations fetched successfully based on recent symptoms", responseData)
}
