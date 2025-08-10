package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/region"
	"ipincamp/srikandi-sehat/src/utils"
	"math"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMyProfile(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	var user models.User

	err := database.DB.
		Preload("Roles").
		Preload("Profile.Village.Classification").
		Preload("Profile.Village.District.Regency.Province").
		First(&user, "uuid = ?", userUUID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	if user.Profile.ID == 0 {
		return utils.SendError(c, fiber.StatusNotFound, "Your profile has not been created yet. Please update your profile first.")
	}

	responseData := dto.UserResponseJson(user)
	return utils.SendSuccess(c, fiber.StatusOK, "Profile fetched successfully", responseData)
}

func UpdateOrCreateProfile(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)
	input := c.Locals("request_body").(*dto.UpdateProfileRequest)

	tx := database.DB.Begin()
	if tx.Error != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	var user models.User
	if err := tx.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	var profile models.Profile
	err := tx.Where(models.Profile{UserID: user.ID}).First(&profile).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		profile = models.Profile{
			UserID:              user.ID,
			PhoneNumber:         "",
			HeightCM:            0,
			WeightKG:            0,
			LastEducation:       constants.EduNone,
			ParentLastEducation: constants.EduNone,
			ParentLastJob:       "",
			InternetAccess:      constants.AccessCellular,
			MenarcheAge:         0,
			DateOfBirth:         nil,
			VillageID:           nil,
		}
		if err := tx.Model(&models.Profile{}).Create(&profile).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create profile")
		}
	}

	updateData := make(map[string]interface{})

	if input.Name != nil && *input.Name != user.Name {
		if err := tx.Model(&user).Update("name", *input.Name).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update user name")
		}
	}

	if input.PhoneNumber != nil && *input.PhoneNumber != profile.PhoneNumber {
		updateData["phone_number"] = *input.PhoneNumber
	}
	if input.HeightCM != nil && *input.HeightCM != profile.HeightCM {
		updateData["height_cm"] = *input.HeightCM
	}
	if input.WeightKG != nil && *input.WeightKG != profile.WeightKG {
		updateData["weight_kg"] = *input.WeightKG
	}
	if input.LastEducation != nil && *input.LastEducation != profile.LastEducation {
		updateData["last_education"] = *input.LastEducation
	}
	if input.ParentLastEducation != nil && *input.ParentLastEducation != profile.ParentLastEducation {
		updateData["parent_last_education"] = *input.ParentLastEducation
	}
	if input.ParentLastJob != nil && *input.ParentLastJob != profile.ParentLastJob {
		updateData["parent_last_job"] = *input.ParentLastJob
	}
	if input.InternetAccess != nil && *input.InternetAccess != profile.InternetAccess {
		updateData["internet_access"] = *input.InternetAccess
	}
	if input.MenarcheAge != nil && *input.MenarcheAge != profile.MenarcheAge {
		updateData["menarche_age"] = *input.MenarcheAge
	}

	if input.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *input.DateOfBirth)
		if err == nil {
			loc, locErr := time.LoadLocation(config.Get("TIMEZONE"))
			if locErr == nil {
				dob = dob.In(loc)
			}
			var currentDOBStr string
			if profile.DateOfBirth != nil {
				currentDOBStr = profile.DateOfBirth.Format("2006-01-02")
			}
			newDOBStr := dob.Format("2006-01-02")
			if profile.DateOfBirth == nil || currentDOBStr != newDOBStr {
				updateData["date_of_birth"] = &dob
			}
		}
	}
	if input.VillageCode != nil {
		var village region.Village
		if err := tx.First(&village, "code = ?", *input.VillageCode).Error; err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "Village code not found")
		}
		if profile.VillageID == nil || village.ID != *profile.VillageID {
			updateData["village_id"] = &village.ID
		}
	}

	if len(updateData) > 0 {
		if err := tx.Model(&profile).Updates(updateData).Error; err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update profile")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile updated successfully", nil)
}

func ChangeMyPassword(c *fiber.Ctx) error {
	userUUID := c.Locals("user_id").(string)

	input := c.Locals("request_body").(*dto.ChangePasswordRequest)

	var user models.User
	if err := database.DB.First(&user, "uuid = ?", userUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	match, err := utils.CheckPasswordHash(input.OldPassword, user.Password)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to verify old password")
	}
	if !match {
		return utils.SendError(c, fiber.StatusUnauthorized, "Old password is incorrect")
	}

	newHashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to process new password")
	}

	if err := database.DB.Model(&user).Update("password", newHashedPassword).Error; err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update password")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Password changed successfully", nil)
}

func GetAllUsers(c *fiber.Ctx) error {
	queries := c.Locals("request_queries").(*dto.UserQuery)

	var users []models.User

	subQuery := database.DB.Table("user_roles").
		Select("1").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = users.id AND roles.name = ?", string(constants.AdminRole))

	query := database.DB.Model(&models.User{}).
		Joins("JOIN profiles ON users.id = profiles.user_id").
		Joins("JOIN villages ON profiles.village_id = villages.id").
		Joins("JOIN classifications ON villages.classification_id = classifications.id")

	if queries.Classification != "" {
		query = query.Where("classifications.name = ?", queries.Classification)
	}

	query = query.Where("NOT EXISTS (?)", subQuery)

	page := queries.Page
	if page == 0 {
		page = 1
	}
	limit := queries.Limit
	if limit == 0 {
		limit = 10
	}

	pagination, paginateScope := utils.GeneratePagination(page, limit, query, &models.User{})

	query.Select("users.uuid, users.name, users.created_at").Scopes(paginateScope).Find(&users)

	var responseData []fiber.Map
	for _, user := range users {
		responseData = append(responseData, fiber.Map{
			"id":         user.UUID,
			"name":       user.Name,
			"created_at": user.CreatedAt,
		})
	}

	paginatedResponse := fiber.Map{
		"data": responseData,
		"meta": pagination,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Users fetched successfully", paginatedResponse)
}

func GetUserByID(c *fiber.Ctx) error {
	params := c.Locals("request_params").(*dto.UserParam)
	userUUID := params.ID

	var user models.User
	result := database.DB.
		Preload("Roles").
		Preload("Profile.Village.Classification").
		Preload("Profile.Village.District.Regency.Province").
		First(&user, "uuid = ?", userUUID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return utils.SendError(c, fiber.StatusNotFound, "User not found")
	}

	responseData := dto.UserResponseJson(user)

	return utils.SendSuccess(c, fiber.StatusOK, "User fetched successfully", responseData)
}

func GetUserStatistics(c *fiber.Ctx) error {
	var stats dto.UserStatisticsResponse
	var wg sync.WaitGroup
	var dbErr error

	// Channel to handle potential concurrent errors
	errChan := make(chan error, 4)

	// --- 1. Get Total Rural Users ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		err := database.DB.Model(&models.User{}).
			Joins("JOIN profiles ON users.id = profiles.user_id").
			Joins("JOIN villages ON profiles.village_id = villages.id").
			Joins("JOIN classifications ON villages.classification_id = classifications.id").
			Where("classifications.name = ?", "Perdesaan").
			Count(&count).Error
		if err != nil {
			errChan <- err
			return
		}
		stats.TotalRuralUsers = count
	}()

	// --- 2. Get Total Urban Users ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		err := database.DB.Model(&models.User{}).
			Joins("JOIN profiles ON users.id = profiles.user_id").
			Joins("JOIN villages ON profiles.village_id = villages.id").
			Joins("JOIN classifications ON villages.classification_id = classifications.id").
			Where("classifications.name = ?", "Perkotaan").
			Count(&count).Error
		if err != nil {
			errChan <- err
			return
		}
		stats.TotalUrbanUsers = count
	}()

	// --- 3. Get Total Active Users ---
	// (Users with at least 2 cycles, meaning 1 completed and 1 active/completed)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		err := database.DB.Table("users").
			Joins("JOIN (SELECT user_id, COUNT(id) as cycle_count FROM menstrual_cycles GROUP BY user_id) as mc ON users.id = mc.user_id").
			Where("mc.cycle_count >= 2").
			Count(&count).Error
		if err != nil {
			errChan <- err
			return
		}
		stats.TotalActiveUsers = count
	}()

	// --- 4. Get Total Users (excluding Admins) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		var count int64
		// Subquery to find users who are admins
		subQuery := database.DB.Table("user_roles").
			Select("user_id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.name = ?", string(constants.AdminRole))

		err := database.DB.Model(&models.User{}).
			Where("id NOT IN (?)", subQuery).
			Count(&count).Error
		if err != nil {
			errChan <- err
			return
		}
		stats.TotalUsers = count
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check if any goroutine reported an error
	for err := range errChan {
		if err != nil {
			dbErr = err // Capture the first error
			break
		}
	}

	if dbErr != nil {
		utils.ErrorLogger.Println("Failed to retrieve user statistics:", dbErr)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve user statistics")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User statistics fetched successfully", stats)
}

func DownloadUsersCSV(c *fiber.Ctx) error {
	subQuery := database.DB.Table("user_roles").
		Select("user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", string(constants.AdminRole))

	var records []dto.UserCSVRecord
	err := database.DB.Model(&models.User{}).
		Select(`
			users.uuid, users.name, users.email,
			DATE(profiles.date_of_birth) as date_of_birth,
			profiles.phone_number, profiles.height_cm, profiles.weight_kg,
			profiles.menarche_age, profiles.last_education, profiles.parent_last_education,
			profiles.parent_last_job, profiles.internet_access,
			villages.name as village, districts.name as district,
			regencies.name as regency, provinces.name as province,
			classifications.name as classification,
			users.created_at as registered_at
		`).
		Joins("LEFT JOIN profiles ON users.id = profiles.user_id").
		Joins("LEFT JOIN villages ON profiles.village_id = villages.id").
		Joins("LEFT JOIN districts ON villages.district_id = districts.id").
		Joins("LEFT JOIN regencies ON districts.regency_id = regencies.id").
		Joins("LEFT JOIN provinces ON regencies.province_id = provinces.id").
		Joins("LEFT JOIN classifications ON villages.classification_id = classifications.id").
		Where("users.id NOT IN (?)", subQuery).
		Order("users.created_at DESC").
		Scan(&records).Error

	if err != nil {
		utils.ErrorLogger.Println("Failed to fetch user data for CSV export:", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch user data")
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	header := []string{
		"UUID", "Nama", "Email", "Tanggal Lahir", "No. Telepon", "Tinggi (cm)", "Berat (kg)", "IMT", "Kategori IMT",
		"Usia Menarche", "Pendidikan Terakhir", "Pendidikan Ortu", "Pekerjaan Ortu",
		"Akses Internet", "Desa/Kelurahan", "Kecamatan", "Kabupaten/Kota", "Provinsi",
		"Klasifikasi Alamat", "Tanggal Registrasi",
	}
	if err := w.Write(header); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to write CSV header")
	}

	for _, record := range records {
		if record.HeightCM > 0 && record.WeightKG > 0 {
			heightInMeters := float32(record.HeightCM) / 100
			bmi := record.WeightKG / (heightInMeters * heightInMeters)
			record.BMI = float32(math.Round(float64(bmi)*100) / 100)
			record.BMICategory = getBMICategory(record.BMI)
		}

		dob := ""
		if record.DateOfBirth != nil {
			dob = *record.DateOfBirth
		}

		csvRow := []string{
			record.UUID, record.Name, record.Email, dob, record.PhoneNumber,
			fmt.Sprintf("%d", record.HeightCM), fmt.Sprintf("%.2f", record.WeightKG), fmt.Sprintf("%.2f", record.BMI), record.BMICategory,
			fmt.Sprintf("%d", record.MenarcheAge), record.LastEducation, record.ParentLastEducation, record.ParentLastJob,
			record.InternetAccess, record.Village, record.District, record.Regency, record.Province,
			record.Classification, record.RegisteredAt.Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(csvRow); err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to write CSV row")
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Error writing CSV data")
	}

	filename := fmt.Sprintf("users_export_%s.csv", time.Now().Format("2006-01-02"))
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.Send(b.Bytes())
}

func getBMICategory(bmi float32) string {
	if bmi <= 0 {
		return ""
	}
	if bmi < 18.5 {
		return "Kurus"
	} else if bmi >= 18.5 && bmi <= 24.9 {
		return "Normal"
	} else if bmi >= 25.0 && bmi <= 29.9 {
		return "Gemuk"
	} else {
		return "Obesitas"
	}
}
