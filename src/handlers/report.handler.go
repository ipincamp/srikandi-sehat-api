package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/utils"
	"math"
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
	if length < 2 {
		return "Pendek (Hipomenorea)"
	} else if length > 7 {
		return "Panjang (Menoragia)"
	}
	return "Normal"
}

// getCycleCategory determines the cycle length category based on the new value ranges.
func getCycleCategory(length int16) string {
	if length == 0 {
		return "N/A"
	}
	if length < 21 {
		return "Pendek (Polimenorea)"
	} else if length > 35 {
		return "Panjang (Oligomenorea)"
	}
	return "Normal"
}

// GenerateFullReportLink membuat token sekali pakai dan mengembalikan URL unduhan. (Admin only)
func GenerateFullReportLink(c *fiber.Ctx) error {
	token := uuid.New().String()
	expiration := 5 * time.Minute // Tautan hanya valid selama 5 menit
	expiresAt := time.Now().Add(expiration)

	// Simpan token ke cache
	utils.StoreReportToken(token, expiration)

	// Buat URL lengkap
	// Pastikan aplikasi Anda berjalan di belakang proxy yang mengatur X-Forwarded-Proto atau atur BaseURL secara manual jika perlu
	downloadURL := fmt.Sprintf("%s/api/reports/download/%s", c.BaseURL(), token)

	response := dto.GenerateReportResponse{
		DownloadURL: downloadURL,
		ExpiresAt:   expiresAt,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "One-time download link generated successfully. Link expires in 5 minutes.", response)
}

// DownloadFullReportByToken menghasilkan CSV jika token valid.
func DownloadFullReportByToken(c *fiber.Ctx) error {
	// 1. Validasi token dari URL
	token := c.Params("token")
	if !utils.UseReportToken(token) {
		// Jika token tidak ada (sudah dipakai/kedaluwarsa), kirim error
		return utils.SendError(c, fiber.StatusNotFound, "Link is invalid, has expired, or has already been used.")
	}

	// 2. Jika token valid, lanjutkan dengan logika pembuatan CSV yang ada
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
		utils.ErrorLogger.Println("Failed to fetch cycle data for full export:", err)
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch cycle data")
	}

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

	var records []dto.FullExportRecord
	userCycleCount := make(map[uint]int64)

	for _, cycle := range cycles {
		userCycleCount[cycle.UserID]++
		user := cycle.User
		profile := user.Profile

		record := dto.FullExportRecord{
			// UserUUID:            user.UUID,
			UserName:            user.Name,
			UserEmail:           user.Email,
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

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	header := []string{
		"Nama Pengguna", "Email", "Tanggal Registrasi", "Umur", "No. Telepon",
		"Tinggi (cm)", "Berat (kg)", "IMT", "Kategori IMT", "Usia Menarche", "Pendidikan Terakhir",
		"Pendidikan Ortu", "Pekerjaan Ortu", "Akses Internet", "Desa/Kelurahan", "Kecamatan",
		"Kabupaten/Kota", "Provinsi", "Klasifikasi Alamat", "Siklus Ke-", "Tanggal Mulai", "Tanggal Selesai",
		"Lama Haid (Hari)", "Kategori Lama Haid", "Panjang Siklus (Hari)", "Kategori Panjang Siklus", "Gejala yang Dirasakan",
	}
	w.Write(header)
	for _, rec := range records {
		row := []string{
			rec.UserName, rec.UserEmail, rec.UserRegisteredAt.Format("2006-01-02 15:04:05"), fmt.Sprintf("%d", rec.Age), rec.PhoneNumber,
			fmt.Sprintf("%d", rec.HeightCM), fmt.Sprintf("%.2f", rec.WeightKG), fmt.Sprintf("%.2f", rec.BMI), rec.BMICategory, fmt.Sprintf("%d", rec.MenarcheAge), rec.LastEducation,
			rec.ParentLastEducation, rec.ParentLastJob, rec.InternetAccess, rec.Village, rec.District,
			rec.Regency, rec.Province, rec.Classification, fmt.Sprintf("%d", rec.CycleNumber), rec.StartDate, rec.EndDate,
			fmt.Sprintf("%d", rec.PeriodLength), rec.PeriodCategory, fmt.Sprintf("%d", rec.CycleLength), rec.CycleCategory, rec.Symptoms,
		}
		w.Write(row)
	}
	w.Flush()

	filename := fmt.Sprintf("report_srikandi-sehat_%s.csv", time.Now().Format("2006-01-02_15-04-05")) // Ganti : dengan - agar aman di nama file
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.Send(b.Bytes())
}
