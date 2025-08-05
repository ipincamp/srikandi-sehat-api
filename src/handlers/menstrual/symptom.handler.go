package menstrual

import (
	"ipincamp/srikandi-sehat/database"
	dto "ipincamp/srikandi-sehat/src/dto/menstrual"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LogSymptoms(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	input := c.Locals("request_body").(*dto.SymptomLogRequest)

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

	tx.Where("symptom_log_id = ?", symptomLog.ID).Delete(&menstrual.SymptomLogDetail{})

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
