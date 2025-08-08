package seeders

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models/menstrual"
	"log"

	"gorm.io/gorm"
)

func SeedMenstrualData(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [SYMPTOMS] Seeding menstrual symptoms and recommendations...")
	// Gejala (Symptoms)
	symptoms := []menstrual.Symptom{
		{Name: "Dismenore", Category: "Nyeri", Type: constants.SymptomTypeBasic},
		{Name: "Kram Perut", Category: "Nyeri", Type: constants.SymptomTypeBasic},
		{Name: "5L", Category: "Fisik", Type: constants.SymptomTypeBasic},
		{Name: "Mood Swing", Category: "Emosional", Type: constants.SymptomTypeOptions},
	}

	for _, s := range symptoms {
		if err := tx.FirstOrCreate(&s, menstrual.Symptom{Name: s.Name}).Error; err != nil {
			log.Printf("[DB] [SEED] [SYMPTOMS] Error seeding symptom: %s\n", s.Name)
			return err
		}
	}
	log.Println("[DB] [SEED] [SYMPTOMS] Menstrual symptoms seeded successfully.")

	// Pilihan Gejala (Symptom Options)
	var moodSymptom menstrual.Symptom
	if err := tx.First(&moodSymptom, "name = ?", "Mood Swing").Error; err != nil {
		log.Println("[DB] [SEED] [SYMPTOM_OPTIONS] Error finding Mood Swing symptom:", err)
		return err
	}

	symptomOptions := []menstrual.SymptomOption{
		{SymptomID: moodSymptom.ID, Name: "Senang", Value: string(constants.MoodTypeHappy)},
		{SymptomID: moodSymptom.ID, Name: "Biasa", Value: string(constants.MoodTypeNeutral)},
		{SymptomID: moodSymptom.ID, Name: "Galau", Value: string(constants.MoodTypeAnxious)},
		{SymptomID: moodSymptom.ID, Name: "Sedih", Value: string(constants.MoodTypeSad)},
		{SymptomID: moodSymptom.ID, Name: "Marah", Value: string(constants.MoodTypeAngry)},
	}

	log.Println("[DB] [SEED] [SYMPTOM_OPTIONS] Seeding symptom options...")
	for _, so := range symptomOptions {
		if err := tx.FirstOrCreate(&so, menstrual.SymptomOption{SymptomID: so.SymptomID, Name: so.Name}).Error; err != nil {
			log.Printf("[DB] [SEED] [SYMPTOM_OPTIONS] Error seeding symptom option: %s\n", so.Name)
			return err
		}
	}
	log.Println("[DB] [SEED] [SYMPTOM_OPTIONS] Symptom options seeded successfully.")

	// Rekomendasi (Recommendations)
	log.Println("[DB] [SEED] [RECOMMENDATIONS] Seeding recommendations for menstrual symptoms...")
	// "Dismenore"
	var dysmenorrheaSymptom menstrual.Symptom
	if err := tx.First(&dysmenorrheaSymptom, "name = ?", "Dismenore").Error; err != nil {
		log.Println("[DB] [SEED] [RECOMMENDATIONS] Error finding Dismenore symptom:", err)
		return err
	}
	// "Kram Perut"
	var crampSymptom menstrual.Symptom
	if err := tx.First(&crampSymptom, "name = ?", "Kram Perut").Error; err != nil {
		log.Println("[DB] [SEED] [RECOMMENDATIONS] Error finding Kram Perut symptom:", err)
		return err
	}

	// "5L"
	var fiveLSymptom menstrual.Symptom
	if err := tx.First(&fiveLSymptom, "name = ?", "5L").Error; err != nil {
		log.Println("[DB] [SEED] [RECOMMENDATIONS] Error finding 5L symptom:", err)
		return err
	}

	// "Mood Swing"
	var moodSwingSymptom menstrual.Symptom
	if err := tx.First(&moodSwingSymptom, "name = ?", "Mood Swing").Error; err != nil {
		log.Println("[DB] [SEED] [RECOMMENDATIONS] Error finding Mood Swing symptom:", err)
		return err
	}

	recommendations := []menstrual.Recommendation{
		// dismenorrheaSymptom
		{
			SymptomID:   dysmenorrheaSymptom.ID,
			Title:       "Tips Alami untuk Meredakan Nyeri Haid Ringan",
			Description: "Untuk menekan rasa sakit, cukup dilakukan kompres hangat, olahraga teratur, istirahat yang cukup, minum air kelapa hijau, minuman jahe, sereh, serta jamu kunir-asem, dan akupresur/ Sanyinjiao Hegu atau senam dismenorea. Apabila nyeri haid yang dirasakan sampai mengganggu aktivitas sehari-hari, bisa diberikan obat anti peradangan yang bersifat non steroid atau berkonsultasi langsung dengan tenaga kesehatan.",
			Source:      "",
		},
		{
			SymptomID:   dysmenorrheaSymptom.ID,
			Title:       "Akupresur Titik Sanyinjiao dan Hegu Untuk Mengatasi Dismenore",
			Description: "Kamu juga bisa mencoba teknik pijat akupresur di titik Sanyinjiao dan Hegu. Titik ini dipercaya bisa membantu mengurangi nyeri haid dengan merangsang sistem saraf tertentu.",
			Source:      "https://youtu.be/l7Z91rEHD6w",
		},
		{
			SymptomID:   dysmenorrheaSymptom.ID,
			Title:       "Senam Dismenorea untuk Remaja",
			Description: "Selain bantu meringankan nyeri, senam ini juga bisa bikin tubuh terasa lebih segar dan mood jadi lebih baik. Gerakannya sederhana dan bisa dilakukan di rumah.",
			Source:      "https://youtu.be/z_wcXr-gIiU",
		},
		// crampSymptom
		{
			SymptomID:   crampSymptom.ID,
			Title:       "Asupan Kalium untuk Redakan Kram Perut Saat Haid",
			Description: "Konsumsi makanan yang tinggi kalium seperti ubi jalar, pisang, salmon, kismis, kacang-kacangan, dan yoghurt. Mengolah makanan dengan cara dikukus atau dipanggang juga dapat membantu meningkatkan asupan kalium dalam tubuh.",
			Source:      "",
		},
		// fiveLSymptom
		{
			SymptomID:   fiveLSymptom.ID,
			Title:       "Lemah, Letih, Lesu, Lemas, Lunglai",
			Description: "Untuk mencegah anemia, saat menstruasi, minumlah 1 tablet penambah darah (tablet Fe) selama menstruasi setiap hari dan sekali seminggu ketika tidak menstruasi.",
			Source:      "",
		},
		// moodSymptom
		{
			SymptomID:   moodSwingSymptom.ID,
			Title:       "Cara Mengatasi Mood Swing Saat Menstruasi",
			Description: "Lakukan aroma terapi, meditasi, atau aktivitas relaksasi lainnya untuk membantu menstabilkan suasana hati.",
			Source:      "",
		},
	}

	for _, r := range recommendations {
		if err := tx.FirstOrCreate(&r, menstrual.Recommendation{SymptomID: r.SymptomID, Title: r.Title}).Error; err != nil {
			log.Printf("[DB] [SEED] [RECOMMENDATIONS] Error seeding recommendation: %s\n", r.Title)
			return err
		}
	}

	log.Println("[DB] [SEED] [RECOMMENDATIONS] Recommendations seeded successfully.")
	return nil
}
