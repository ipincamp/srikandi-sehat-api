package menstrual

import (
	"ipincamp/srikandi-sehat/database"
	dto "ipincamp/srikandi-sehat/src/dto/menstrual"
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

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	if input.StartDate != "" {
		startDate, _ := time.Parse("2006-01-02", input.StartDate)

		newCycle := menstrual.MenstrualCycle{UserID: user.ID, StartDate: startDate}
		if err := tx.Create(&newCycle).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to record new cycle")
		}

		updatePreviousCycleLength(tx, user.ID, startDate)
	}

	if input.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", input.EndDate)
		updateCurrentCyclePeriod(tx, user.ID, endDate)
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Cycle data recorded successfully", nil)
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
		log.Printf("Updated previous cycle for user %d. Length: %d days.", userID, cycleLength)
	}
}

func updateCurrentCyclePeriod(tx *gorm.DB, userID uint, endDate time.Time) {
	var currentCycle menstrual.MenstrualCycle
	err := tx.Where("user_id = ? AND end_date IS NULL", userID).
		Order("start_date desc").
		First(&currentCycle).Error

	if err == nil {
		periodLength := int16(endDate.Sub(currentCycle.StartDate).Hours()/24) + 1
		isNormal := periodLength >= 2 && periodLength <= 7

		tx.Model(&currentCycle).Updates(map[string]interface{}{
			"end_date":         endDate,
			"period_length":    periodLength,
			"is_period_normal": isNormal,
		})
		log.Printf("Updated current period for user %d. Duration: %d days.", userID, periodLength)
	}
}
