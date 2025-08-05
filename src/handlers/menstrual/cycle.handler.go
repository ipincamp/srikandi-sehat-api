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
	defer tx.Rollback()

	activeCycle, err := findActiveCycle(tx, user.ID)
	isStart := true

	if input.StartDate != "" {
		if err == nil {
			return utils.SendError(c, fiber.StatusConflict, "Cannot start a new cycle while another is in progress.")
		}

		startDate, _ := time.Parse("2006-01-02", input.StartDate)

		newCycle := menstrual.MenstrualCycle{UserID: user.ID, StartDate: startDate}
		if err := tx.Create(&newCycle).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to record new cycle")
		}

		updatePreviousCycleLength(tx, user.ID, startDate)
	}

	if input.EndDate != "" {
		if err != nil {
			return utils.SendError(c, fiber.StatusConflict, "No active cycle to end.")
		} else {
			endDate, _ := time.Parse("2006-01-02", input.EndDate)

			updateCurrentCyclePeriod(tx, activeCycle.UserID, endDate)
			isStart = false
		}
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	var message string
	if isStart {
		message = "Cycle started successfully"
	} else {
		message = "Cycle ended successfully"
	}
	return utils.SendSuccess(c, fiber.StatusOK, message, nil)
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
	}
}
