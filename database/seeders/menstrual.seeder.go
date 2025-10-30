package seeders

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"log"

	"gorm.io/gorm"
)

func SeedMenstrualData(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [SYMPTOMS] Seeding menstrual symptoms, options, and recommendations...")

	// --- TAHAP 1: Seed Symptoms (Gejala) ---
	symptomMap := make(map[string]uint)

	symptoms := []menstrual.Symptom{
		{Name: "Dismenore", Category: "Nyeri", Type: constants.SymptomTypeBasic},
		{Name: "Kram Perut", Category: "Nyeri", Type: constants.SymptomTypeBasic},
		{Name: "5L", Category: "Fisik", Type: constants.SymptomTypeBasic},
		{Name: "Mood Swing", Category: "Emosional", Type: constants.SymptomTypeOptions},
	}

	for i := range symptoms {
		// Pointer ke elemen slice (&symptoms[i])
		s := &symptoms[i]
		if err := tx.FirstOrCreate(s, menstrual.Symptom{Name: s.Name}).Error; err != nil {
			log.Printf("[DB] [SEED] [SYMPTOMS] Error seeding symptom: %s\n", s.Name)
			return err
		}
		// Simpan ID ke "cache" map
		symptomMap[s.Name] = s.ID
	}
	log.Println("[DB] [SEED] [SYMPTOMS] Menstrual symptoms seeded successfully.")

	// --- TAHAP 2: Seed Symptom Options (Pilihan Gejala) ---
	// Ambil ID dari map (tanpa kueri DB)
	moodSwingID, ok := symptomMap["Mood Swing"]
	if !ok {
		log.Println("[DB] [SEED] [SYMPTOM_OPTIONS] Error finding Mood Swing symptom ID in cache map.")
		// Fail fast if map is wrong
		return gorm.ErrRecordNotFound
	}

	symptomOptions := []menstrual.SymptomOption{
		{SymptomID: moodSwingID, Name: "Senang", Value: string(constants.MoodTypeHappy)},
		{SymptomID: moodSwingID, Name: "Biasa", Value: string(constants.MoodTypeNeutral)},
		{SymptomID: moodSwingID, Name: "Galau", Value: string(constants.MoodTypeAnxious)},
		{SymptomID: moodSwingID, Name: "Sedih", Value: string(constants.MoodTypeSad)},
		{SymptomID: moodSwingID, Name: "Marah", Value: string(constants.MoodTypeAngry)},
	}

	log.Println("[DB] [SEED] [SYMPTOM_OPTIONS] Seeding symptom options...")
	for _, so := range symptomOptions {
		// Buat salinan lokal agar GORM tidak bingung dengan pointer loop
		option := so
		if err := tx.FirstOrCreate(&option, menstrual.SymptomOption{SymptomID: option.SymptomID, Name: option.Name}).Error; err != nil {
			log.Printf("[DB] [SEED] [SYMPTOM_OPTIONS] Error seeding symptom option: %s\n", option.Name)
			return err
		}
	}
	log.Println("[DB] [SEED] [SYMPTOM_OPTIONS] Symptom options seeded successfully.")

	// --- TAHAP 3: Seed Recommendations (Rekomendasi) ---
	log.Println("[DB] [SEED] [RECOMMENDATIONS] Seeding recommendations...")

	// Ambil semua ID yang dibutuhkan dari map (tanpa kueri DB)
	dismenoreID := symptomMap["Dismenore"]
	crampSymptomID := symptomMap["Kram Perut"]
	fiveLSymptomID := symptomMap["5L"]
	moodSwingSymptomID := symptomMap["Mood Swing"]

	recommendations := []menstrual.Recommendation{
		// dismenoreID
		{
			SymptomID:   dismenoreID,
			Title:       "Tips Alami untuk Meredakan Nyeri Haid Ringan",
			Description: "Untuk menekan rasa sakit, cukup dilakukan kompres hangat, olahraga teratur, istirahat yang cukup, minum air kelapa hijau, minuman jahe, sereh, serta jamu kunir-asem, dan akupresur/ Sanyinjiao Hegu atau senam dismenorea. Apabila nyeri haid yang dirasakan sampai mengganggu aktivitas sehari-hari, bisa diberikan obat anti peradangan yang bersifat non steroid atau berkonsultasi langsung dengan tenaga kesehatan.",
			Source:      "",
		},
		{
			SymptomID:   dismenoreID,
			Title:       "Akupresur Titik Sanyinjiao dan Hegu Untuk Mengatasi Dismenore",
			Description: "",
			Source:      "https://youtu.be/l7Z91rEHD6w",
		},
		{
			SymptomID:   dismenoreID,
			Title:       "Senam Dismenorea untuk Remaja",
			Description: "",
			Source:      "https://youtu.be/z_wcXr-gIiU",
		},
		// crampSymptomID
		{
			SymptomID:   crampSymptomID,
			Title:       "Asupan Kalium untuk Redakan Kram Perut Saat Haid",
			Description: "Konsumsi makanan yang tinggi kalium seperti ubi jalar, pisang, salmon, kismis, kacang-kacangan, dan yoghurt. Mengolah makanan dengan cara dikukus atau dipanggang juga dapat membantu meningkatkan asupan kalium dalam tubuh.",
			Source:      "",
		},
		// fiveLSymptomID
		{
			SymptomID:   fiveLSymptomID,
			Title:       "Lemah, Letih, Lesu, Lemas, Lunglai",
			Description: "Untuk mencegah anemia, saat menstruasi, minumlah 1 tablet penambah darah (tablet Fe) selama menstruasi setiap hari dan sekali seminggu ketika tidak menstruasi.",
			Source:      "",
		},
		// moodSwingSymptomID
		{
			SymptomID:   moodSwingSymptomID,
			Title:       "Cara Mengatasi Mood Swing Saat Menstruasi",
			Description: "Lakukan aroma terapi, meditasi, atau aktivitas relaksasi lainnya untuk membantu menstabilkan suasana hati.",
			Source:      "",
		},
	}

	for _, r := range recommendations {
		rec := r
		if err := tx.FirstOrCreate(&rec, menstrual.Recommendation{SymptomID: rec.SymptomID, Title: rec.Title}).Error; err != nil {
			log.Printf("[DB] [SEED] [RECOMMENDATIONS] Error seeding recommendation: %s\n", rec.Title)
			return err
		}
	}

	log.Println("[DB] [SEED] [RECOMMENDATIONS] Recommendations seeded successfully.")
	return nil
}
