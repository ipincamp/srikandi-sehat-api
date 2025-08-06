package menstrual

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
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

	logDate, _ := time.Parse("2006-01-02", input.LogDate)

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	symptomLog := menstrual.SymptomLog{
		UserID:  user.ID,
		LogDate: logDate,
	}
	tx.Where("user_id = ? AND log_date = ?", user.ID, logDate).FirstOrCreate(&symptomLog)
	if input.Note != "" {
		tx.Model(&symptomLog).Update("note", input.Note)
	}

	var activeCycle menstrual.MenstrualCycle
	err := tx.Where("user_id = ? AND start_date <= ? AND end_date >= ?", user.ID, logDate, logDate).
		First(&activeCycle).Error
	if err == nil {
		tx.Model(&symptomLog).Update("menstrual_cycle_id", activeCycle.ID)
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

func GetSymptomLogsByDateRange(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	queries := c.Locals("request_queries").(*dto.SymptomLogQuery)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	startDate, _ := time.Parse("2006-01-02", queries.StartDate)
	endDate, _ := time.Parse("2006-01-02", queries.EndDate)

	var logs []menstrual.SymptomLog
	database.DB.
		Preload("Details.Symptom").
		Preload("Details.SymptomOption").
		Where("user_id = ? AND log_date BETWEEN ? AND ?", user.ID, startDate, endDate).
		Order("log_date asc").
		Find(&logs)

	var responseData []dto.SymptomLogResponse
	for _, log := range logs {
		var detailsDTO []dto.SymptomLogDetailResponse
		for _, detail := range log.Details {
			detailDTO := dto.SymptomLogDetailResponse{
				SymptomName:     detail.Symptom.Name,
				SymptomCategory: detail.Symptom.Category,
			}
			if detail.SymptomOptionID.Valid {
				detailDTO.SelectedOption = detail.SymptomOption.Name
			}
			detailsDTO = append(detailsDTO, detailDTO)
		}

		responseData = append(responseData, dto.SymptomLogResponse{
			LogDate: log.LogDate,
			Note:    log.Note,
			Details: detailsDTO,
		})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Symptom logs fetched successfully", responseData)
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
