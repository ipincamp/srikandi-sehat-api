package handlers

import (
	"fmt"
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// --- Helper functions for CSV Export ---

// calculateAge calculates age based on a birth date.
func calculateAge(birthDate *time.Time) int {
	if birthDate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

// getBMICategory determines the BMI category based on the new value ranges.
func getBMICategory(bmi float32) string {
	if bmi <= 0 {
		return ""
	}
	if bmi < 17.0 {
		return "Sangat Kurus"
	} else if bmi >= 17.0 && bmi < 18.5 {
		return "Kurus"
	} else if bmi >= 18.5 && bmi <= 25.0 {
		return "Normal"
	} else if bmi > 25.0 && bmi <= 27.0 {
		return "Gemuk"
	} else { // > 27.0
		return "Obesitas"
	}
}

// getPeriodCategory determines the period duration category based on the new value ranges.
func getPeriodCategory(length int16) string {
	if length == 0 {
		return "N/A"
	}
	if length < constants.CyclePeriodMinNormalDays {
		return "Pendek (Hipomenorea)"
	} else if length > constants.CyclePeriodMaxNormalDays {
		return "Panjang (Menoragia)"
	}
	return "Normal"
}

// getCycleCategory determines the cycle length category based on the new value ranges.
func getCycleCategory(length int16) string {
	if length == 0 {
		return "N/A"
	}
	if length < constants.CycleLengthMinNormalDays {
		return "Pendek (Polimenorea)"
	} else if length > constants.CycleLengthMaxNormalDays {
		return "Panjang (Oligomenorea)"
	}
	return "Normal"
}

// maskEmail masks the local part of an email for privacy.
// Example: tknhtpX@luHbVdL.edu -> tkn***@luHbVdL.edu
func maskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email // Not a valid email format, return as-is
	}

	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) <= 3 {
		return localPart + "***@" + domainPart
	}

	return localPart[:3] + "***@" + domainPart
}

// --- Handlers ---

// GenerateFullReportLink generates XLSX report, saves it as a static file, and returns an encrypted token (Admin only)
func GenerateFullReportLink(c *fiber.Ctx) error {
	// 1. Generate unique token and password
	token := uuid.New().String()
	password := utils.GenerateRandomPassword()
	expiration := 30 * time.Minute // Link valid for 30 minutes
	expiresAt := time.Now().Add(expiration)

	// 2. Query report data
	subQuery := database.DB.Table("user_roles").
		Select("user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", string(constants.AdminRole))

	var cycles []menstrual.MenstrualCycle
	err := database.DB.
		Preload("User.Profile.Village.Classification").
		Preload("User.Profile.Village.District.Regency.Province").
		Where("user_id NOT IN (?)", subQuery).
		Order("user_id, start_date ASC").
		Find(&cycles).Error

	if err != nil {
		utils.ErrorLogger.Println("Failed to fetch cycle data for report:", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch cycle data")
	}

	// 3. Fetch symptoms
	var allSymptomLogs []menstrual.SymptomLog
	database.DB.Preload("Details.Symptom").
		Where("menstrual_cycle_id IS NOT NULL").
		Find(&allSymptomLogs)

	symptomsByCycleID := make(map[int64][]string)
	for _, log := range allSymptomLogs {
		if log.MenstrualCycleID.Valid {
			for _, detail := range log.Details {
				symptomsByCycleID[log.MenstrualCycleID.Int64] = append(symptomsByCycleID[log.MenstrualCycleID.Int64], detail.Symptom.Name)
			}
		}
	}

	// 4. Build report records
	var records []dto.FullExportRecord
	userCycleCount := make(map[uint]int64)

	for _, cycle := range cycles {
		userCycleCount[cycle.UserID]++
		user := cycle.User
		profile := user.Profile

		record := dto.FullExportRecord{
			UserName:            user.Name,
			UserEmail:           maskEmail(user.Email),
			UserRegisteredAt:    user.CreatedAt,
			Age:                 calculateAge(profile.DateOfBirth),
			PhoneNumber:         profile.PhoneNumber,
			HeightCM:            profile.HeightCM,
			WeightKG:            profile.WeightKG,
			MenarcheAge:         profile.MenarcheAge,
			LastEducation:       string(profile.LastEducation),
			ParentLastEducation: string(profile.ParentLastEducation),
			ParentLastJob:       profile.ParentLastJob,
			InternetAccess:      string(profile.InternetAccess),
		}

		if profile.Village.ID > 0 {
			record.Village = profile.Village.Name
			record.District = profile.Village.District.Name
			record.Regency = profile.Village.District.Regency.Name
			record.Province = profile.Village.District.Regency.Province.Name
			record.Classification = profile.Village.Classification.Name
		}

		if record.HeightCM > 0 && record.WeightKG > 0 {
			heightInMeters := float32(record.HeightCM) / 100
			bmi := record.WeightKG / (heightInMeters * heightInMeters)
			record.BMI = float32(math.Round(float64(bmi)*100) / 100)
			record.BMICategory = getBMICategory(record.BMI)
		}

		endDate := ""
		if cycle.EndDate.Valid {
			endDate = cycle.EndDate.Time.Format("2006-01-02")
		}

		symptoms := "Tidak ada gejala tercatat"
		if symptomNames, found := symptomsByCycleID[int64(cycle.ID)]; found {
			uniqueSymptoms := make(map[string]bool)
			for _, name := range symptomNames {
				uniqueSymptoms[name] = true
			}
			var uniqueNames []string
			for name := range uniqueSymptoms {
				uniqueNames = append(uniqueNames, name)
			}
			symptoms = strings.Join(uniqueNames, "; ")
		}

		record.CycleNumber = userCycleCount[cycle.UserID]
		record.StartDate = cycle.StartDate.Format("2006-01-02")
		record.EndDate = endDate
		record.PeriodLength = cycle.PeriodLength.Int16
		record.PeriodCategory = getPeriodCategory(cycle.PeriodLength.Int16)
		record.CycleLength = cycle.CycleLength.Int16
		record.CycleCategory = getCycleCategory(cycle.CycleLength.Int16)
		record.Symptoms = symptoms

		records = append(records, record)
	}

	// 5. Create reports directory if not exists
	reportsDir := "./reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		utils.ErrorLogger.Println("Failed to create reports directory:", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to create reports directory")
	}

	// 6. Generate XLSX file
	filename := fmt.Sprintf("report_srikandi-sehat_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))
	filePath := filepath.Join(reportsDir, filename)

	if err := utils.GenerateReportXLSX(records, filePath); err != nil {
		utils.ErrorLogger.Println("Failed to generate XLSX:", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to generate report file")
	}

	// 7. Encrypt token with password
	encryptedToken, err := utils.EncryptToken(token, password)
	if err != nil {
		utils.ErrorLogger.Println("Failed to encrypt token:", err)
		os.Remove(filePath) // Clean up the file
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to generate secure token")
	}

	// 8. Store metadata in cache
	metadata := utils.ReportMetadata{
		Filename:  filename,
		Password:  password,
		Token:     token,
		ExpiresAt: expiresAt,
		FilePath:  filePath,
		Used:      false,
	}
	utils.StoreReportMetadata(encryptedToken, metadata)

	// 9. Build download URL
	downloadURL := fmt.Sprintf("%s/api/reports/download", config.Get("APP_BASE_URL"))

	response := dto.GenerateReportResponse{
		DownloadURL:    downloadURL,
		EncryptedToken: encryptedToken,
		ExpiresAt:      expiresAt,
	}

	return utils.SendSuccess(c, fiber.StatusOK,
		fmt.Sprintf("Report generated successfully. Use the encrypted token and password to download. Password: %s", password),
		response)
}

// DownloadFullReport validates password and encrypted token, then downloads the XLSX file
func DownloadFullReport(c *fiber.Ctx) error {
	// 1. Parse and validate request body
	input := c.Locals("request_body").(*dto.ValidateReportRequest)
	token := input.Token

	// 2. Get metadata from cache
	metadata, found := utils.GetReportMetadata(token)
	if !found {
		return utils.SendError(c, fiber.StatusNotFound, "Link is invalid, has expired, or has already been used")
	}

	// 3. Check if already used
	if metadata.Used {
		return utils.SendError(c, fiber.StatusGone, "This download link has already been used")
	}

	// 4. Check expiration
	if time.Now().After(metadata.ExpiresAt) {
		utils.DeleteReportMetadata(token, true)
		os.Remove(metadata.FilePath) // Clean up expired file
		return utils.SendError(c, fiber.StatusGone, "Download link has expired")
	}

	// 5. Decrypt and validate token using password
	decryptedToken, err := utils.DecryptToken(token, input.Password)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid password")
	}

	if decryptedToken != metadata.Token {
		return utils.SendError(c, fiber.StatusUnauthorized, "Token validation failed")
	}

	// 6. Check if file exists
	if _, err := os.Stat(metadata.FilePath); os.IsNotExist(err) {
		utils.DeleteReportMetadata(token, false)
		return utils.SendError(c, fiber.StatusNotFound, "Report file not found")
	}

	// 7. Mark as used
	if !utils.MarkReportAsUsed(token) {
		return utils.SendError(c, fiber.StatusConflict, "Failed to mark report as used")
	}

	// 8. Send the file
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", metadata.Filename))

	// Schedule file deletion after successful download
	go func() {
		time.Sleep(5 * time.Second)
		if err := os.Remove(metadata.FilePath); err != nil {
			utils.ErrorLogger.Printf("Failed to delete report file %s: %v", metadata.FilePath, err)
		} else {
			utils.InfoLogger.Printf("Report file deleted: %s", metadata.FilePath)
		}
		utils.DeleteReportMetadata(token, false)
	}()

	return c.SendFile(metadata.FilePath)
}
