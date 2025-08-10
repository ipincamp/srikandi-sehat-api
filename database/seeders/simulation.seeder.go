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

// SeedSimulationData creates 50 dummy users with realistic cycle and symptom data.
func SeedSimulationData(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [SIMULATION] Starting user and cycle simulation...")

	// 1. Get necessary master data
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

	// 2. Loop to create 50 users
	for i := 0; i < 50; i++ {
		// --- Create User and Profile ---
		email := faker.Email()
		hashedPassword, _ := utils.HashPassword("password") // Default password for all dummy users

		user := models.User{
			Name:     faker.Name(),
			Email:    email,
			Password: hashedPassword,
		}
		if err := tx.Create(&user).Error; err != nil {
			log.Printf("Failed to create dummy user %s: %v", email, err)
			continue
		}
		tx.Model(&user).Association("Roles").Append(&userRole)

		// Create profile for the user
		dob := time.Now().AddDate(-(rand.Intn(7) + 13), rand.Intn(12), rand.Intn(28)) // Age 13-20
		profile := models.Profile{
			UserID:              user.ID,
			PhoneNumber:         faker.Phonenumber(),
			DateOfBirth:         &dob,
			HeightCM:            uint(rand.Intn(25) + 145),   // 145-170 cm
			WeightKG:            float32(rand.Intn(30) + 40), // 40-70 kg
			MenarcheAge:         uint(rand.Intn(4) + 11),     // 11-15 years old
			VillageID:           &villages[rand.Intn(len(villages))].ID,
			LastEducation:       constants.EducationLevel("SMP"),
			ParentLastEducation: constants.EducationLevel("SMA"),
			ParentLastJob:       "Wiraswasta",
			InternetAccess:      constants.InternetAccess("Seluler"),
		}
		tx.Create(&profile)

		// --- Simulate Menstrual Cycles ---
		totalCycles := rand.Intn(10) + 3 // Generate 3 to 12 cycles per user
		lastStartDate := time.Now().AddDate(0, -totalCycles, 0)

		for j := 0; j < totalCycles; j++ {
			periodLength := int16(rand.Intn(5) + 3)  // 3-7 days
			cycleLength := int16(rand.Intn(10) + 23) // 23-32 days

			startDate := lastStartDate.AddDate(0, 0, int(cycleLength))
			endDate := startDate.AddDate(0, 0, int(periodLength-1))

			cycle := menstrual.MenstrualCycle{
				UserID:         user.ID,
				StartDate:      startDate,
				EndDate:        sql.NullTime{Time: endDate, Valid: true},
				PeriodLength:   sql.NullInt16{Int16: periodLength, Valid: true},
				CycleLength:    sql.NullInt16{Int16: cycleLength, Valid: true},
				IsPeriodNormal: sql.NullBool{Bool: periodLength >= 2 && periodLength <= 7, Valid: true},
				IsCycleNormal:  sql.NullBool{Bool: cycleLength >= 21 && cycleLength <= 35, Valid: true},
			}
			tx.Create(&cycle)

			// --- Simulate Symptoms for this cycle (75% chance) ---
			if rand.Intn(100) < 75 {
				logDate := startDate.AddDate(0, 0, rand.Intn(int(periodLength)))
				symptomLog := menstrual.SymptomLog{
					UserID:           user.ID,
					LoggedAt:         logDate,
					Note:             "Ini adalah catatan dummy dari seeder.",
					MenstrualCycleID: sql.NullInt64{Int64: int64(cycle.ID), Valid: true},
				}
				tx.Create(&symptomLog)

				// Add 1 to 3 random symptoms
				numSymptomsToLog := rand.Intn(3) + 1
				for k := 0; k < numSymptomsToLog; k++ {
					symptom := symptoms[rand.Intn(len(symptoms))]
					detail := menstrual.SymptomLogDetail{
						SymptomLogID: symptomLog.ID,
						SymptomID:    symptom.ID,
					}
					// If symptom has options, pick one randomly
					if symptom.Type == constants.SymptomTypeOptions && len(symptom.Options) > 0 {
						option := symptom.Options[rand.Intn(len(symptom.Options))]
						detail.SymptomOptionID = sql.NullInt64{Int64: int64(option.ID), Valid: true}
					}
					tx.Create(&detail)
				}
			}
			lastStartDate = startDate
		}
		log.Printf("[DB] [SEED] [SIMULATION] Created user %d (%s) with %d cycles.", i+1, email, totalCycles)
	}

	log.Println("[DB] [SEED] [SIMULATION] 50 dummy users and their cycle data have been seeded successfully.")
	return nil
}
