package menstrual

import (
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"math"
	"strconv"
	"strings"
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
	}
	tx.Where("user_id = ? AND DATE(logged_at) = ?", user.ID, loggedAt.Format("2006-01-02")).FirstOrCreate(&symptomLog)
	if input.Note != "" {
		tx.Model(&symptomLog).Update("note", input.Note)
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

	return utils.SendSuccess(c, fiber.StatusOK, "Symptoms logged successfully", nil)
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

	dataQuery := database.DB.Model(&menstrual.SymptomLog{}).
		Select("MIN(symptom_logs.id) as id, COUNT(symptom_log_details.id) as total_symptoms, DATE(symptom_logs.logged_at) as log_date").
		Joins("left join symptom_log_details on symptom_log_details.symptom_log_id = symptom_logs.id").
		Where("symptom_logs.user_id = ?", user.ID).
		Group("DATE(symptom_logs.logged_at)")

	countQuery := database.DB.Model(&menstrual.SymptomLog{}).
		Where("user_id = ?", user.ID)

	if queries.Date != "" {
		dataQuery = dataQuery.Where("DATE(symptom_logs.logged_at) = ?", queries.Date)
		countQuery = countQuery.Where("DATE(logged_at) = ?", queries.Date)
	} else if queries.StartDate != "" && queries.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", queries.EndDate)
		if err != nil {
			return utils.SendError(c, fiber.StatusBadRequest, "Invalid EndDate format")
		}
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		dataQuery = dataQuery.Where("symptom_logs.logged_at BETWEEN ? AND ?", queries.StartDate, endDate)
		countQuery = countQuery.Where("logged_at BETWEEN ? AND ?", queries.StartDate, endDate)
	}

	var totalRows int64
	if err := countQuery.Distinct("DATE(logged_at)").Count(&totalRows).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to count symptom history")
	}

	if totalRows == 0 {
		return utils.SendError(c, fiber.StatusNotFound, "No symptom history found for the given criteria.")
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	offset := (page - 1) * limit

	var results []dto.SymptomHistoryResponse
	err := dataQuery.Offset(offset).Limit(limit).Order("log_date desc").Scan(&results).Error
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve symptom history")
	}

	paginatedResponse := dto.PaginatedResponse[dto.SymptomHistoryResponse]{
		Data: results,
		Metadata: dto.Pagination{
			Limit:       limit,
			TotalRows:   totalRows,
			TotalPages:  totalPages,
			CurrentPage: page,
		},
	}
	if page > 1 {
		prev := page - 1
		paginatedResponse.Metadata.PreviousPage = &prev
	}
	if page < totalPages {
		next := page + 1
		paginatedResponse.Metadata.NextPage = &next
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

func isValidCommaSeparatedInt(input string) bool {
	for _, char := range input {
		if (char < '0' || char > '9') && char != ',' {
			return false
		}
	}
	return true
}

func GetRecommendationsBySymptoms(c *fiber.Ctx) error {
	queries := c.Locals("request_queries").(*dto.RecommendationQuery)

	// Validasi input query parameter
	if queries.SymptomIDs == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "symptom_ids query parameter is required")
	}
	// Jika symptom_ids tidak valid, kembalikan error
	if !isValidCommaSeparatedInt(queries.SymptomIDs) {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid symptom_ids format")
	}

	// 1. Konversi string "1,4,5" menjadi slice of integer []int{1, 4, 5}
	// TODO: remove duplicate. example: "1,1,1,2,3,4" to "1,2,3,4"
	idStrings := strings.Split(queries.SymptomIDs, ",")
	var symptomIDs []int
	for _, idStr := range idStrings {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			symptomIDs = append(symptomIDs, id)
		}
	}

	if len(symptomIDs) == 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid or empty symptom_ids")
	}

	// 2. Cari semua rekomendasi yang cocok di database
	var recommendations []menstrual.Recommendation
	// Gunakan Preload("Symptom") untuk mendapatkan nama gejala terkait
	err := database.DB.
		Preload("Symptom").
		Where("symptom_id IN ?", symptomIDs).
		Find(&recommendations).Error

	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch recommendations")
	}

	// 3. Petakan hasil dari model ke DTO
	var responseData []dto.RecommendationResponse
	for _, r := range recommendations {
		responseData = append(responseData, dto.RecommendationResponse{
			ForSymptom:  r.Symptom.Name,
			Title:       r.Title,
			Description: r.Description,
			Source:      r.Source,
		})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Recommendations fetched successfully", responseData)
}
