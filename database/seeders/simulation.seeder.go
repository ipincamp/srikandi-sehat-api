package seeders

import (
	"database/sql"
	"fmt"
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"ipincamp/srikandi-sehat/src/models/region"
	"ipincamp/srikandi-sehat/src/utils"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type UserRole struct {
	UserID uint
	RoleID uint
}

func SeedSimulationData(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [SIMULATION] Starting user and cycle simulation...")

	const totalUsers = 50
	const userBatchSize = 50
	const dataBatchSize = 100

	// --- TAHAP 1: Ambil Master Data (Dependensi) ---
	// Ambil 1x di awal
	log.Println("[DB] [SEED] [SIMULATION] Fetching master data (roles, villages, symptoms)...")
	var userRole models.Role
	if err := tx.First(&userRole, "name = ?", constants.UserRole).Error; err != nil {
		return fmt.Errorf("failed to find user role: %w", err)
	}

	var villages []region.Village
	if err := tx.Find(&villages).Error; err != nil {
		return fmt.Errorf("failed to find villages: %w", err)
	}
	if len(villages) == 0 {
		return fmt.Errorf("no villages found in database, please seed regions first")
	}

	var symptoms []menstrual.Symptom
	if err := tx.Preload("Options").Find(&symptoms).Error; err != nil {
		return fmt.Errorf("failed to find symptoms: %w", err)
	}

	hashedPassword, _ := utils.HashPassword("password")

	// --- TAHAP 2: Generate Users di Memori ---
	log.Printf("[DB] [SEED] [SIMULATION] Generating %d users in memory...", totalUsers)
	usersToCreate := make([]models.User, 0, totalUsers)
	for i := 0; i < totalUsers; i++ {
		usersToCreate = append(usersToCreate, models.User{
			Name:  faker.Name(),
			Email: faker.Email(),
			// Password: hashedPassword,
		})
	}

	// --- TAHAP 3: Batch Insert Users ---
	if err := tx.CreateInBatches(&usersToCreate, userBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert users: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d users.", len(usersToCreate))

	// --- TAHAP 4: Generate Data Turunan (Profiles, Roles, Cycles, AuthProviders) di Memori ---
	log.Println("[DB] [SEED] [SIMULATION] Generating profiles, roles, auth providers, and cycles in memory...")
	profilesToCreate := make([]models.Profile, 0, totalUsers)
	userRolesToCreate := make([]UserRole, 0, totalUsers)
	cyclesToCreate := make([]menstrual.MenstrualCycle, 0, totalUsers*7) // Estimasi 7 siklus per user
	authProvidersToCreate := make([]models.UserAuthProvider, 0, totalUsers)

	for _, user := range usersToCreate {
		// Generate Profile
		dob := time.Now().AddDate(-(rand.Intn(7) + 13), rand.Intn(12), rand.Intn(28)) // Usia 13-20
		profilesToCreate = append(profilesToCreate, models.Profile{
			UserID:              user.ID,
			PhoneNumber:         faker.Phonenumber(),
			DateOfBirth:         &dob,
			HeightCM:            uint(rand.Intn(25) + 145),   // 145-170 cm
			WeightKG:            float32(rand.Intn(30) + 40), // 40-70 kg
			MenarcheAge:         uint(rand.Intn(10) + 11),    // 11-20 tahun
			VillageID:           &villages[rand.Intn(len(villages))].ID,
			LastEducation:       constants.EducationLevel("SMP"),
			ParentLastEducation: constants.EducationLevel("SMA"),
			ParentLastJob:       "Wiraswasta",
			InternetAccess:      constants.AccessCellular,
		})

		// Generate UserRole
		userRolesToCreate = append(userRolesToCreate, UserRole{
			UserID: user.ID,
			RoleID: userRole.ID,
		})

		authProvidersToCreate = append(authProvidersToCreate, models.UserAuthProvider{
			UserID:   user.ID,
			Provider: "local",
			Password: sql.NullString{String: hashedPassword, Valid: true},
		})

		// Generate Cycles
		totalCycles := rand.Intn(10) + 3 // 3-12 siklus per user
		lastStartDate := time.Now().AddDate(0, -totalCycles, 0)

		for j := 0; j < totalCycles; j++ {
			periodLength := int16(rand.Intn(5) + 3)  // 3-7 hari
			cycleLength := int16(rand.Intn(10) + 23) // 23-32 hari

			startDate := lastStartDate.AddDate(0, 0, int(cycleLength))
			endDate := startDate.AddDate(0, 0, int(periodLength-1))

			cyclesToCreate = append(cyclesToCreate, menstrual.MenstrualCycle{
				UserID:         user.ID,
				StartDate:      startDate,
				EndDate:        sql.NullTime{Time: endDate, Valid: true},
				PeriodLength:   sql.NullInt16{Int16: periodLength, Valid: true},
				CycleLength:    sql.NullInt16{Int16: cycleLength, Valid: true},
				IsPeriodNormal: sql.NullBool{Bool: periodLength >= 2 && periodLength <= 7, Valid: true},
				IsCycleNormal:  sql.NullBool{Bool: cycleLength >= 21 && cycleLength <= 35, Valid: true},
			})
			lastStartDate = startDate
		}
	}

	// --- TAHAP 5: Batch Insert Profiles, Roles, dan Cycles ---
	if err := tx.CreateInBatches(&profilesToCreate, dataBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert profiles: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d profiles.", len(profilesToCreate))

	if err := tx.Table("user_roles").CreateInBatches(&userRolesToCreate, dataBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert user_roles: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d user role links.", len(userRolesToCreate))

	if err := tx.CreateInBatches(&cyclesToCreate, dataBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert cycles: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d total cycles.", len(cyclesToCreate))

	if err := tx.CreateInBatches(&authProvidersToCreate, dataBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert auth providers: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d auth providers.", len(authProvidersToCreate))

	// --- TAHAP 6: Generate Symptom Logs di Memori ---
	log.Println("[DB] [SEED] [SIMULATION] Generating symptom logs in memory...")
	logsToCreate := make([]menstrual.SymptomLog, 0, len(cyclesToCreate))

	for _, cycle := range cyclesToCreate {
		// 75% kemungkinan mencatat gejala
		if rand.Intn(100) < 75 {
			periodLength := int(cycle.PeriodLength.Int16)
			if periodLength <= 0 {
				periodLength = 1
			}
			logDate := cycle.StartDate.AddDate(0, 0, rand.Intn(periodLength))

			logsToCreate = append(logsToCreate, menstrual.SymptomLog{
				UserID:           cycle.UserID,
				LoggedAt:         logDate,
				Note:             "Ini adalah catatan dummy dari seeder.",
				MenstrualCycleID: sql.NullInt64{Int64: int64(cycle.ID), Valid: true},
			})
		}
	}

	// --- TAHAP 7: Batch Insert Symptom Logs ---
	if err := tx.CreateInBatches(&logsToCreate, dataBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert symptom logs: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d symptom logs.", len(logsToCreate))

	// --- TAHAP 8: Generate Symptom Log Details di Memori ---
	log.Println("[DB] [SEED] [SIMULATION] Generating symptom log details in memory...")
	detailsToCreate := make([]menstrual.SymptomLogDetail, 0, len(logsToCreate)*2) // Estimasi 2 detail per log

	for _, logEntry := range logsToCreate {
		numSymptomsToLog := rand.Intn(3) + 1 // 1-3 gejala per log
		for k := 0; k < numSymptomsToLog; k++ {
			symptom := symptoms[rand.Intn(len(symptoms))]
			detail := menstrual.SymptomLogDetail{
				SymptomLogID: logEntry.ID,
				SymptomID:    symptom.ID,
			}
			// Jika gejala punya opsi, pilih salah satu secara acak
			if symptom.Type == constants.SymptomTypeOptions && len(symptom.Options) > 0 {
				option := symptom.Options[rand.Intn(len(symptom.Options))]
				detail.SymptomOptionID = sql.NullInt64{Int64: int64(option.ID), Valid: true}
			}
			detailsToCreate = append(detailsToCreate, detail)
		}
	}

	// --- TAHAP 9: Batch Insert Symptom Log Details ---
	if err := tx.CreateInBatches(&detailsToCreate, dataBatchSize).Error; err != nil {
		return fmt.Errorf("failed to batch insert symptom log details: %w", err)
	}
	log.Printf("[DB] [SEED] [SIMULATION] Successfully inserted %d symptom log details.", len(detailsToCreate))

	log.Println("[DB] [SEED] [SIMULATION] 50 dummy users and their cycle data have been seeded successfully.")
	return nil
}
