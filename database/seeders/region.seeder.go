package seeders

import (
	"ipincamp/srikandi-sehat/src/constants"
	"ipincamp/srikandi-sehat/src/models/region"
	"log"

	"gorm.io/gorm"
)

func SeedRegions(tx *gorm.DB) error {
	log.Println("[DB] [SEED] [REGION] Seeding regions...")

	perkotaan, perdesaan, err := seedClassifications(tx)
	if err != nil {
		return err
	}

	if err := seedBanyumasRegion(tx, perkotaan, perdesaan); err != nil {
		return err
	}

	log.Println("[DB] [SEED] [REGION] Seeding regions completed successfully.")
	return nil
}

func seedClassifications(tx *gorm.DB) (region.Classification, region.Classification, error) {
	log.Println("[DB] [SEED] [CLASSIFICATION] Seeding classifications...")
	var perkotaan = region.Classification{Name: string(constants.UrbanClassification)}
	var perdesaan = region.Classification{Name: string(constants.RuralClassification)}

	if err := tx.FirstOrCreate(&perkotaan, perkotaan).Error; err != nil {
		return perkotaan, perdesaan, err
	}
	if err := tx.FirstOrCreate(&perdesaan, perdesaan).Error; err != nil {
		return perkotaan, perdesaan, err
	}

	log.Println("[DB] [SEED] [CLASSIFICATION] Classifications seeded successfully.")
	return perkotaan, perdesaan, nil
}

func seedBanyumasRegion(tx *gorm.DB, perkotaan, perdesaan region.Classification) error {
	log.Println("[DB] [SEED] [REGION] [BANYUMAS] Seeding Banyumas region...")

	provinceData := region.Province{
		Code: "33",
		Name: "JAWA TENGAH",
		Regencies: []region.Regency{
			{
				Code: "3302",
				Name: "BANYUMAS",
				Districts: []region.District{
					// Kecamatan Lumbir
					{
						Code: "3302010",
						Name: "LUMBIR",
						Villages: []region.Village{
							{Code: "3302010001", Name: "CINGEBUL", ClassificationID: perdesaan.ID},
							{Code: "3302010002", Name: "KEDUNGGEDE", ClassificationID: perdesaan.ID},
							{Code: "3302010003", Name: "CIDORA", ClassificationID: perdesaan.ID},
							{Code: "3302010004", Name: "BESUKI", ClassificationID: perdesaan.ID},
							{Code: "3302010005", Name: "PARUNGKAMAL", ClassificationID: perdesaan.ID},
							{Code: "3302010006", Name: "CIRAHAB", ClassificationID: perdesaan.ID},
							{Code: "3302010007", Name: "CANDUK", ClassificationID: perdesaan.ID},
							{Code: "3302010008", Name: "KARANGGAYAM", ClassificationID: perdesaan.ID},
							{Code: "3302010009", Name: "LUMBIR", ClassificationID: perdesaan.ID},
							{Code: "3302010010", Name: "DERMAJI", ClassificationID: perdesaan.ID},
						},
					},
					// Kecamatan Wangon
					{
						Code: "3302020",
						Name: "WANGON",
						Villages: []region.Village{
							{Code: "3302020001", Name: "RANDEGAN", ClassificationID: perkotaan.ID},
							{Code: "3302020002", Name: "RAWAHENG", ClassificationID: perkotaan.ID},
							{Code: "3302020003", Name: "PENGADEGAN", ClassificationID: perdesaan.ID},
							{Code: "3302020004", Name: "KLAPAGADING", ClassificationID: perkotaan.ID},
							{Code: "3302020005", Name: "KLAPAGADING KULON", ClassificationID: perkotaan.ID},
							{Code: "3302020006", Name: "WANGON", ClassificationID: perkotaan.ID},
							{Code: "3302020007", Name: "BANTERAN", ClassificationID: perkotaan.ID},
							{Code: "3302020008", Name: "JAMBU", ClassificationID: perkotaan.ID},
							{Code: "3302020009", Name: "JURANGBAHAS", ClassificationID: perdesaan.ID},
							{Code: "3302020010", Name: "CIKAKAK", ClassificationID: perdesaan.ID},
							{Code: "3302020011", Name: "WLAHAR", ClassificationID: perkotaan.ID},
							{Code: "3302020012", Name: "WINDUNEGARA", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Jatilawang
					{
						Code: "3302030",
						Name: "JATILAWANG",
						Villages: []region.Village{
							{Code: "3302030001", Name: "GUNUNG WETAN", ClassificationID: perdesaan.ID},
							{Code: "3302030002", Name: "PEKUNCEN", ClassificationID: perdesaan.ID},
							{Code: "3302030003", Name: "KARANGLEWAS", ClassificationID: perdesaan.ID},
							{Code: "3302030004", Name: "KARANGANYAR", ClassificationID: perdesaan.ID},
							{Code: "3302030005", Name: "MARGASANA", ClassificationID: perkotaan.ID},
							{Code: "3302030006", Name: "ADISARA", ClassificationID: perkotaan.ID},
							{Code: "3302030007", Name: "KEDUNGWRINGIN", ClassificationID: perkotaan.ID},
							{Code: "3302030008", Name: "BANTAR", ClassificationID: perdesaan.ID},
							{Code: "3302030009", Name: "TINGGARJAYA", ClassificationID: perkotaan.ID},
							{Code: "3302030010", Name: "TUNJUNG", ClassificationID: perkotaan.ID},
							{Code: "3302030011", Name: "GENTAWANGI", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Rawalo
					{
						Code: "3302040",
						Name: "RAWALO",
						Villages: []region.Village{
							{Code: "3302040001", Name: "LOSARI", ClassificationID: perdesaan.ID},
							{Code: "3302040002", Name: "MENGANTI", ClassificationID: perkotaan.ID},
							{Code: "3302040003", Name: "BANJARPARAKAN", ClassificationID: perkotaan.ID},
							{Code: "3302040004", Name: "RAWALO", ClassificationID: perkotaan.ID},
							{Code: "3302040005", Name: "TAMBAKNEGARA", ClassificationID: perkotaan.ID},
							{Code: "3302040006", Name: "SIDAMULIH", ClassificationID: perkotaan.ID},
							{Code: "3302040007", Name: "PESAWAHAN", ClassificationID: perkotaan.ID},
							{Code: "3302040008", Name: "TIPAR", ClassificationID: perkotaan.ID},
							{Code: "3302040009", Name: "SANGGREMAN", ClassificationID: perdesaan.ID},
						},
					},
					// Kecamatan Kebasen
					{
						Code: "3302050",
						Name: "KEBASEN",
						Villages: []region.Village{
							{Code: "3302050001", Name: "ADISANA", ClassificationID: perkotaan.ID},
							{Code: "3302050002", Name: "BANGSA", ClassificationID: perkotaan.ID},
							{Code: "3302050003", Name: "KARANGSARI", ClassificationID: perkotaan.ID},
							{Code: "3302050004", Name: "RANDEGAN", ClassificationID: perkotaan.ID},
							{Code: "3302050005", Name: "KALIWEDI", ClassificationID: perkotaan.ID},
							{Code: "3302050006", Name: "SAWANGAN", ClassificationID: perkotaan.ID},
							{Code: "3302050007", Name: "KALISALAK", ClassificationID: perkotaan.ID},
							{Code: "3302050008", Name: "CINDAGA", ClassificationID: perkotaan.ID},
							{Code: "3302050009", Name: "KEBASEN", ClassificationID: perkotaan.ID},
							{Code: "3302050010", Name: "GAMBARSARI", ClassificationID: perkotaan.ID},
							{Code: "3302050011", Name: "TUMIYANG", ClassificationID: perkotaan.ID},
							{Code: "3302050012", Name: "MANDIRANCAN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Kemranjen
					{
						Code: "3302060",
						Name: "KEMRANJEN",
						Villages: []region.Village{
							{Code: "3302060001", Name: "GRUJUGAN", ClassificationID: perkotaan.ID},
							{Code: "3302060002", Name: "SIRAU", ClassificationID: perkotaan.ID},
							{Code: "3302060003", Name: "SIBALUNG", ClassificationID: perkotaan.ID},
							{Code: "3302060004", Name: "SIBRAMA", ClassificationID: perkotaan.ID},
							{Code: "3302060005", Name: "KEDUNGPRING", ClassificationID: perkotaan.ID},
							{Code: "3302060006", Name: "KECILA", ClassificationID: perkotaan.ID},
							{Code: "3302060007", Name: "NUSAMANGIR", ClassificationID: perdesaan.ID},
							{Code: "3302060008", Name: "KARANGJATI", ClassificationID: perkotaan.ID},
							{Code: "3302060009", Name: "KEBARONGAN", ClassificationID: perkotaan.ID},
							{Code: "3302060010", Name: "SIDAMULYA", ClassificationID: perkotaan.ID},
							{Code: "3302060011", Name: "PAGERALANG", ClassificationID: perkotaan.ID},
							{Code: "3302060012", Name: "ALASMALANG", ClassificationID: perkotaan.ID},
							{Code: "3302060013", Name: "PETARANGAN", ClassificationID: perkotaan.ID},
							{Code: "3302060014", Name: "KARANGGINTUNG", ClassificationID: perdesaan.ID},
							{Code: "3302060015", Name: "KARANGSALAM", ClassificationID: perdesaan.ID},
						},
					},
					// Kecamatan Sumpiuh
					{
						Code: "3302070",
						Name: "SUMPIUH",
						Villages: []region.Village{
							{Code: "3302070001", Name: "PANDAK", ClassificationID: perkotaan.ID},
							{Code: "3302070002", Name: "KUNTILI", ClassificationID: perkotaan.ID},
							{Code: "3302070003", Name: "KEMIRI", ClassificationID: perkotaan.ID},
							{Code: "3302070004", Name: "KARANGGEDANG", ClassificationID: perdesaan.ID},
							{Code: "3302070005", Name: "NUSADADI", ClassificationID: perdesaan.ID},
							{Code: "3302070006", Name: "SELANDAKA", ClassificationID: perkotaan.ID},
							{Code: "3302070007", Name: "SUMPIUH", ClassificationID: perkotaan.ID},
							{Code: "3302070008", Name: "KRADENAN", ClassificationID: perkotaan.ID},
							{Code: "3302070009", Name: "SELANEGARA", ClassificationID: perkotaan.ID},
							{Code: "3302070010", Name: "KEBOKURA", ClassificationID: perkotaan.ID},
							{Code: "3302070011", Name: "LEBENG", ClassificationID: perdesaan.ID},
							{Code: "3302070012", Name: "KETANDA", ClassificationID: perkotaan.ID},
							{Code: "3302070013", Name: "BANJARPANEPEN", ClassificationID: perdesaan.ID},
							{Code: "3302070014", Name: "BOGANGIN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Tambak
					{
						Code: "3302080",
						Name: "TAMBAK",
						Villages: []region.Village{
							{Code: "3302080001", Name: "PLANGKAPAN", ClassificationID: perdesaan.ID},
							{Code: "3302080002", Name: "GUMELAR LOR", ClassificationID: perkotaan.ID},
							{Code: "3302080003", Name: "GUMELAR KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302080004", Name: "KARANGPETIR", ClassificationID: perkotaan.ID},
							{Code: "3302080005", Name: "GEBANGSARI", ClassificationID: perdesaan.ID},
							{Code: "3302080006", Name: "KARANGPUCUNG", ClassificationID: perkotaan.ID},
							{Code: "3302080007", Name: "PREMBUN", ClassificationID: perkotaan.ID},
							{Code: "3302080008", Name: "PESANTREN", ClassificationID: perkotaan.ID},
							{Code: "3302080009", Name: "BUNIAYU", ClassificationID: perkotaan.ID},
							{Code: "3302080010", Name: "PURWODADI", ClassificationID: perkotaan.ID},
							{Code: "3302080011", Name: "KAMULYAN", ClassificationID: perkotaan.ID},
							{Code: "3302080012", Name: "WATUAGUNG", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Somagede
					{
						Code: "3302090",
						Name: "SOMAGEDE",
						Villages: []region.Village{
							{Code: "3302090001", Name: "TANGGERAN", ClassificationID: perdesaan.ID},
							{Code: "3302090002", Name: "SOKAWERA", ClassificationID: perkotaan.ID},
							{Code: "3302090003", Name: "SOMAGEDE", ClassificationID: perkotaan.ID},
							{Code: "3302090004", Name: "KLINTING", ClassificationID: perkotaan.ID},
							{Code: "3302090005", Name: "KEMAWI", ClassificationID: perdesaan.ID},
							{Code: "3302090006", Name: "PIASA KULON", ClassificationID: perkotaan.ID},
							{Code: "3302090007", Name: "KANDING", ClassificationID: perkotaan.ID},
							{Code: "3302090008", Name: "SOMAKATON", ClassificationID: perdesaan.ID},
							{Code: "3302090009", Name: "PLANA", ClassificationID: perdesaan.ID},
						},
					},
					// Kecamatan Kalibagor
					{
						Code: "3302100",
						Name: "KALIBAGOR",
						Villages: []region.Village{
							{Code: "3302100001", Name: "SROWOT", ClassificationID: perkotaan.ID},
							{Code: "3302100002", Name: "SURO", ClassificationID: perdesaan.ID},
							{Code: "3302100003", Name: "KALIORI", ClassificationID: perkotaan.ID},
							{Code: "3302100004", Name: "WLAHAR WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302100005", Name: "PEKAJA", ClassificationID: perkotaan.ID},
							{Code: "3302100006", Name: "KARANGDADAP", ClassificationID: perkotaan.ID},
							{Code: "3302100007", Name: "KALIBAGOR", ClassificationID: perkotaan.ID},
							{Code: "3302100008", Name: "PAJERUKAN", ClassificationID: perkotaan.ID},
							{Code: "3302100009", Name: "PETIR", ClassificationID: perkotaan.ID},
							{Code: "3302100010", Name: "KALICUPAK KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302100011", Name: "KALICUPAK LOR", ClassificationID: perkotaan.ID},
							{Code: "3302100012", Name: "KALISOGRA WETAN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Banyumas
					{
						Code: "3302110",
						Name: "BANYUMAS",
						Villages: []region.Village{
							{Code: "3302110001", Name: "BINANGUN", ClassificationID: perdesaan.ID},
							{Code: "3302110002", Name: "PASINGGANGAN", ClassificationID: perkotaan.ID},
							{Code: "3302110003", Name: "KEDUNGGEDE", ClassificationID: perkotaan.ID},
							{Code: "3302110004", Name: "KARANGRAU", ClassificationID: perkotaan.ID},
							{Code: "3302110005", Name: "KEJAWAR", ClassificationID: perkotaan.ID},
							{Code: "3302110006", Name: "DANARAJA", ClassificationID: perkotaan.ID},
							{Code: "3302110007", Name: "KEDUNGUTER", ClassificationID: perkotaan.ID},
							{Code: "3302110008", Name: "SUDAGARAN", ClassificationID: perkotaan.ID},
							{Code: "3302110009", Name: "PEKUNDEN", ClassificationID: perkotaan.ID},
							{Code: "3302110010", Name: "KALISUBE", ClassificationID: perkotaan.ID},
							{Code: "3302110011", Name: "DAWUHAN", ClassificationID: perkotaan.ID},
							{Code: "3302110012", Name: "PAPRINGAN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Patikraja
					{
						Code: "3302120",
						Name: "PATIKRAJA",
						Villages: []region.Village{
							{Code: "3302120001", Name: "SAWANGAN WETAN", ClassificationID: perdesaan.ID},
							{Code: "3302120002", Name: "KARANGENDEP", ClassificationID: perdesaan.ID},
							{Code: "3302120003", Name: "NOTOG", ClassificationID: perkotaan.ID},
							{Code: "3302120004", Name: "PATIKRAJA", ClassificationID: perkotaan.ID},
							{Code: "3302120005", Name: "PEGALONGAN", ClassificationID: perkotaan.ID},
							{Code: "3302120006", Name: "SOKAWERA", ClassificationID: perkotaan.ID},
							{Code: "3302120007", Name: "WLAHAR KULON", ClassificationID: perkotaan.ID},
							{Code: "3302120008", Name: "KEDUNGRANDU", ClassificationID: perkotaan.ID},
							{Code: "3302120009", Name: "KEDUNGWULUH KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302120010", Name: "KEDUNGWULUH LOR", ClassificationID: perkotaan.ID},
							{Code: "3302120011", Name: "KARANGANYAR", ClassificationID: perkotaan.ID},
							{Code: "3302120012", Name: "SIDABOWA", ClassificationID: perkotaan.ID},
							{Code: "3302120013", Name: "KEDUNGWRINGIN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Purwojati
					{
						Code: "3302130",
						Name: "PURWOJATI",
						Villages: []region.Village{
							{Code: "3302130001", Name: "GERDUREN", ClassificationID: perdesaan.ID},
							{Code: "3302130002", Name: "KARANGTALUN KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302130003", Name: "KALIURIP", ClassificationID: perdesaan.ID},
							{Code: "3302130004", Name: "KARANGTALUN LOR", ClassificationID: perkotaan.ID},
							{Code: "3302130005", Name: "PURWOJATI", ClassificationID: perkotaan.ID},
							{Code: "3302130006", Name: "KLAPASAWIT", ClassificationID: perkotaan.ID},
							{Code: "3302130007", Name: "KARANGMANGU", ClassificationID: perdesaan.ID},
							{Code: "3302130008", Name: "KALIPUTIH", ClassificationID: perdesaan.ID},
							{Code: "3302130009", Name: "KALIWANGI", ClassificationID: perdesaan.ID},
							{Code: "3302130010", Name: "KALITAPEN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Ajibarang
					{
						Code: "3302140",
						Name: "AJIBARANG",
						Villages: []region.Village{
							{Code: "3302140001", Name: "DARMAKRADENAN", ClassificationID: perdesaan.ID},
							{Code: "3302140002", Name: "TIPAR KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302140003", Name: "SAWANGAN", ClassificationID: perkotaan.ID},
							{Code: "3302140004", Name: "JINGKANG", ClassificationID: perdesaan.ID},
							{Code: "3302140005", Name: "BANJARSARI", ClassificationID: perkotaan.ID},
							{Code: "3302140006", Name: "KALIBENDA", ClassificationID: perkotaan.ID},
							{Code: "3302140007", Name: "PANCURENDANG", ClassificationID: perkotaan.ID},
							{Code: "3302140008", Name: "PANCASAN", ClassificationID: perkotaan.ID},
							{Code: "3302140009", Name: "KARANGBAWANG", ClassificationID: perkotaan.ID},
							{Code: "3302140010", Name: "KRACAK", ClassificationID: perkotaan.ID},
							{Code: "3302140011", Name: "AJIBARANG KULON", ClassificationID: perkotaan.ID},
							{Code: "3302140012", Name: "AJIBARANG WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302140013", Name: "LESMANA", ClassificationID: perkotaan.ID},
							{Code: "3302140014", Name: "PANDANSARI", ClassificationID: perkotaan.ID},
							{Code: "3302140015", Name: "CIBERUNG", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Gumelar
					{
						Code: "3302150",
						Name: "GUMELAR",
						Villages: []region.Village{
							{Code: "3302150001", Name: "CILANGKAP", ClassificationID: perdesaan.ID},
							{Code: "3302150002", Name: "CIHONJE", ClassificationID: perdesaan.ID},
							{Code: "3302150003", Name: "PANINGKABAN", ClassificationID: perdesaan.ID},
							{Code: "3302150004", Name: "KARANGKEMOJING", ClassificationID: perdesaan.ID},
							{Code: "3302150005", Name: "GANCANG", ClassificationID: perkotaan.ID},
							{Code: "3302150006", Name: "KEDUNGURANG", ClassificationID: perkotaan.ID},
							{Code: "3302150007", Name: "GUMELAR", ClassificationID: perkotaan.ID},
							{Code: "3302150008", Name: "TLAGA", ClassificationID: perdesaan.ID},
							{Code: "3302150009", Name: "SAMUDRA", ClassificationID: perdesaan.ID},
							{Code: "3302150010", Name: "SAMUDRA KULON", ClassificationID: perdesaan.ID},
						},
					},
					// Kecamatan Pekuncen
					{
						Code: "3302160",
						Name: "PEKUNCEN",
						Villages: []region.Village{
							{Code: "3302160001", Name: "CIBANGKONG", ClassificationID: perdesaan.ID},
							{Code: "3302160002", Name: "PETAHUNAN", ClassificationID: perdesaan.ID},
							{Code: "3302160003", Name: "SEMEDO", ClassificationID: perkotaan.ID},
							{Code: "3302160004", Name: "CIKAWUNG", ClassificationID: perkotaan.ID},
							{Code: "3302160005", Name: "KARANGKLESEM", ClassificationID: perkotaan.ID},
							{Code: "3302160006", Name: "CANDINEGARA", ClassificationID: perkotaan.ID},
							{Code: "3302160007", Name: "CIKEMBULAN", ClassificationID: perkotaan.ID},
							{Code: "3302160008", Name: "TUMIYANG", ClassificationID: perkotaan.ID},
							{Code: "3302160009", Name: "GLEMPANG", ClassificationID: perdesaan.ID},
							{Code: "3302160010", Name: "PEKUNCEN", ClassificationID: perkotaan.ID},
							{Code: "3302160011", Name: "PASIRAMAN LOR", ClassificationID: perkotaan.ID},
							{Code: "3302160012", Name: "PASIRAMAN KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302160013", Name: "BANJARANYAR", ClassificationID: perkotaan.ID},
							{Code: "3302160014", Name: "KARANGKEMIRI", ClassificationID: perkotaan.ID},
							{Code: "3302160015", Name: "KRANGGAN", ClassificationID: perkotaan.ID},
							{Code: "3302160016", Name: "KRAJAN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Cilongok
					{
						Code: "3302170",
						Name: "CILONGOK",
						Villages: []region.Village{
							{Code: "3302170001", Name: "BATUANTEN", ClassificationID: perkotaan.ID},
							{Code: "3302170002", Name: "KASEGERAN", ClassificationID: perkotaan.ID},
							{Code: "3302170003", Name: "JATISABA", ClassificationID: perdesaan.ID},
							{Code: "3302170004", Name: "PANUSUPAN", ClassificationID: perkotaan.ID},
							{Code: "3302170005", Name: "PEJOGOL", ClassificationID: perkotaan.ID},
							{Code: "3302170006", Name: "PAGERAJI", ClassificationID: perkotaan.ID},
							{Code: "3302170007", Name: "SUDIMARA", ClassificationID: perkotaan.ID},
							{Code: "3302170008", Name: "CILONGOK", ClassificationID: perkotaan.ID},
							{Code: "3302170009", Name: "CIPETE", ClassificationID: perkotaan.ID},
							{Code: "3302170010", Name: "CIKIDANG", ClassificationID: perkotaan.ID},
							{Code: "3302170011", Name: "PERNASIDI", ClassificationID: perkotaan.ID},
							{Code: "3302170012", Name: "LANGGONGSARI", ClassificationID: perkotaan.ID},
							{Code: "3302170013", Name: "RANCAMAYA", ClassificationID: perkotaan.ID},
							{Code: "3302170014", Name: "PANEMBANGAN", ClassificationID: perkotaan.ID},
							{Code: "3302170015", Name: "KARANGLO", ClassificationID: perkotaan.ID},
							{Code: "3302170016", Name: "KALISARI", ClassificationID: perkotaan.ID},
							{Code: "3302170017", Name: "KARANGTENGAH", ClassificationID: perdesaan.ID},
							{Code: "3302170018", Name: "SAMBIRATA", ClassificationID: perdesaan.ID},
							{Code: "3302170019", Name: "GUNUNGLURAH", ClassificationID: perkotaan.ID},
							{Code: "3302170020", Name: "SOKAWERA", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Karanglewas
					{
						Code: "3302180",
						Name: "KARANGLEWAS",
						Villages: []region.Village{
							{Code: "3302180001", Name: "KEDIRI", ClassificationID: perkotaan.ID},
							{Code: "3302180002", Name: "PANGEBATAN", ClassificationID: perkotaan.ID},
							{Code: "3302180003", Name: "KARANGLEWAS KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302180004", Name: "TAMANSARI", ClassificationID: perkotaan.ID},
							{Code: "3302180005", Name: "KARANGKEMIRI", ClassificationID: perkotaan.ID},
							{Code: "3302180006", Name: "KARANGGUDE KULON", ClassificationID: perkotaan.ID},
							{Code: "3302180007", Name: "PASIR KULON", ClassificationID: perkotaan.ID},
							{Code: "3302180008", Name: "PASIR WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302180009", Name: "PASIR LOR", ClassificationID: perkotaan.ID},
							{Code: "3302180010", Name: "JIPANG", ClassificationID: perkotaan.ID},
							{Code: "3302180011", Name: "SINGASARI", ClassificationID: perkotaan.ID},
							{Code: "3302180012", Name: "BABAKAN", ClassificationID: perkotaan.ID},
							{Code: "3302180013", Name: "SUNYALANGU", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Kedungbanteng
					{
						Code: "3302230",
						Name: "KEDUNGBANTENG",
						Villages: []region.Village{
							{Code: "3302230001", Name: "KEDUNGBANTENG", ClassificationID: perkotaan.ID},
							{Code: "3302230002", Name: "KEBOCORAN", ClassificationID: perkotaan.ID},
							{Code: "3302230003", Name: "KARANGSALAM KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302230004", Name: "BEJI", ClassificationID: perkotaan.ID},
							{Code: "3302230005", Name: "KARANGNANGKA", ClassificationID: perkotaan.ID},
							{Code: "3302230006", Name: "KENITEN", ClassificationID: perkotaan.ID},
							{Code: "3302230007", Name: "DAWUHAN WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302230008", Name: "DAWUHAN KULON", ClassificationID: perkotaan.ID},
							{Code: "3302230009", Name: "BASEH", ClassificationID: perdesaan.ID},
							{Code: "3302230010", Name: "KALISALAK", ClassificationID: perdesaan.ID},
							{Code: "3302230011", Name: "WINDUJAYA", ClassificationID: perdesaan.ID},
							{Code: "3302230012", Name: "KALIKESUR", ClassificationID: perdesaan.ID},
							{Code: "3302230013", Name: "KUTALIMAN", ClassificationID: perkotaan.ID},
							{Code: "3302230014", Name: "MELUNG", ClassificationID: perdesaan.ID},
						},
					},
					// Kecamatan Baturraden
					{
						Code: "3302220",
						Name: "BATURRADEN",
						Villages: []region.Village{
							{Code: "3302220001", Name: "PURWOSARI", ClassificationID: perkotaan.ID},
							{Code: "3302220002", Name: "KUTASARI", ClassificationID: perkotaan.ID},
							{Code: "3302220003", Name: "PANDAK", ClassificationID: perkotaan.ID},
							{Code: "3302220004", Name: "PAMIJEN", ClassificationID: perkotaan.ID},
							{Code: "3302220005", Name: "REMPOAH", ClassificationID: perkotaan.ID},
							{Code: "3302220006", Name: "KEBUMEN", ClassificationID: perkotaan.ID},
							{Code: "3302220007", Name: "KARANG TENGAH", ClassificationID: perkotaan.ID},
							{Code: "3302220008", Name: "KEMUTUG KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302220009", Name: "KARANGSALAM", ClassificationID: perkotaan.ID},
							{Code: "3302220010", Name: "KEMUTUG LOR", ClassificationID: perkotaan.ID},
							{Code: "3302220011", Name: "KARANGMANGU", ClassificationID: perkotaan.ID},
							{Code: "3302220012", Name: "KETENGER", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Sumbang
					{
						Code: "3302210",
						Name: "SUMBANG",
						Villages: []region.Village{
							{Code: "3302210001", Name: "KARANGGINTUNG", ClassificationID: perkotaan.ID},
							{Code: "3302210002", Name: "TAMBAKSOGRA", ClassificationID: perkotaan.ID},
							{Code: "3302210003", Name: "KARANGCEGAK", ClassificationID: perkotaan.ID},
							{Code: "3302210004", Name: "KARANGTURI", ClassificationID: perdesaan.ID},
							{Code: "3302210005", Name: "SILADO", ClassificationID: perkotaan.ID},
							{Code: "3302210006", Name: "SUSUKAN", ClassificationID: perkotaan.ID},
							{Code: "3302210007", Name: "SUMBANG", ClassificationID: perkotaan.ID},
							{Code: "3302210008", Name: "KEBANGGAN", ClassificationID: perkotaan.ID},
							{Code: "3302210009", Name: "KAWUNGCARANG", ClassificationID: perkotaan.ID},
							{Code: "3302210010", Name: "DATAR", ClassificationID: perkotaan.ID},
							{Code: "3302210011", Name: "BANJARSARI KULON", ClassificationID: perkotaan.ID},
							{Code: "3302210012", Name: "BANJARSARI WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302210013", Name: "BANTERAN", ClassificationID: perkotaan.ID},
							{Code: "3302210014", Name: "CIBEREM", ClassificationID: perkotaan.ID},
							{Code: "3302210015", Name: "SIKAPAT", ClassificationID: perkotaan.ID},
							{Code: "3302210016", Name: "GANDATAPA", ClassificationID: perkotaan.ID},
							{Code: "3302210017", Name: "KOTAYASA", ClassificationID: perkotaan.ID},
							{Code: "3302210018", Name: "LIMPAKUWUS", ClassificationID: perdesaan.ID},
							{Code: "3302210019", Name: "KEDUNGMALANG", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Kembaran
					{
						Code: "3302200",
						Name: "KEMBARAN",
						Villages: []region.Village{
							{Code: "3302200001", Name: "LEDUG", ClassificationID: perkotaan.ID},
							{Code: "3302200002", Name: "PLIKEN", ClassificationID: perkotaan.ID},
							{Code: "3302200003", Name: "PURWODADI", ClassificationID: perkotaan.ID},
							{Code: "3302200004", Name: "KARANG TENGAH", ClassificationID: perkotaan.ID},
							{Code: "3302200005", Name: "KRAMAT", ClassificationID: perkotaan.ID},
							{Code: "3302200006", Name: "SAMBENG WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302200007", Name: "SAMBENG KULON", ClassificationID: perkotaan.ID},
							{Code: "3302200008", Name: "PURBADANA", ClassificationID: perkotaan.ID},
							{Code: "3302200009", Name: "KEMBARAN", ClassificationID: perkotaan.ID},
							{Code: "3302200010", Name: "BOJONGSARI", ClassificationID: perkotaan.ID},
							{Code: "3302200011", Name: "KARANGSOKA", ClassificationID: perkotaan.ID},
							{Code: "3302200012", Name: "DUKUHWALUH", ClassificationID: perkotaan.ID},
							{Code: "3302200013", Name: "TAMBAKSARI KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302200014", Name: "BANTARWUNI", ClassificationID: perkotaan.ID},
							{Code: "3302200015", Name: "KARANGSARI", ClassificationID: perkotaan.ID},
							{Code: "3302200016", Name: "LINGGASARI", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Sokaraja
					{
						Code: "3302190",
						Name: "SOKARAJA",
						Villages: []region.Village{
							{Code: "3302190001", Name: "KALIKIDANG", ClassificationID: perkotaan.ID},
							{Code: "3302190002", Name: "SOKARAJA TENGAH", ClassificationID: perkotaan.ID},
							{Code: "3302190003", Name: "SOKARAJA KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302190004", Name: "SOKARAJA WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302190005", Name: "KLAHANG", ClassificationID: perkotaan.ID},
							{Code: "3302190006", Name: "BANJARSARI KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302190007", Name: "JOMPO KULON", ClassificationID: perkotaan.ID},
							{Code: "3302190008", Name: "BANJARANYAR", ClassificationID: perkotaan.ID},
							{Code: "3302190009", Name: "LEMBERANG", ClassificationID: perkotaan.ID},
							{Code: "3302190010", Name: "KARANGDUREN", ClassificationID: perkotaan.ID},
							{Code: "3302190011", Name: "SOKARAJA LOR", ClassificationID: perkotaan.ID},
							{Code: "3302190012", Name: "KEDONDONG", ClassificationID: perkotaan.ID},
							{Code: "3302190013", Name: "PAMIJEN", ClassificationID: perkotaan.ID},
							{Code: "3302190014", Name: "SOKARAJA KULON", ClassificationID: perkotaan.ID},
							{Code: "3302190015", Name: "KARANGKEDAWUNG", ClassificationID: perkotaan.ID},
							{Code: "3302190016", Name: "WIRADADI", ClassificationID: perkotaan.ID},
							{Code: "3302190017", Name: "KARANGNANAS", ClassificationID: perkotaan.ID},
							{Code: "3302190018", Name: "KARANGRAU", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Purwokerto Selatan
					{
						Code: "3302710",
						Name: "PURWOKERTO SELATAN",
						Villages: []region.Village{
							{Code: "3302710001", Name: "KARANGKLESEM", ClassificationID: perkotaan.ID},
							{Code: "3302710002", Name: "TELUK", ClassificationID: perkotaan.ID},
							{Code: "3302710003", Name: "BERKOH", ClassificationID: perkotaan.ID},
							{Code: "3302710004", Name: "PURWOKERTO KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302710005", Name: "PURWOKERTO KULON", ClassificationID: perkotaan.ID},
							{Code: "3302710006", Name: "KARANGPUCUNG", ClassificationID: perkotaan.ID},
							{Code: "3302710007", Name: "TANJUNG", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Purwokerto Barat
					{
						Code: "3302720",
						Name: "PURWOKERTO BARAT",
						Villages: []region.Village{
							{Code: "3302720001", Name: "KARANGLEWAS LOR", ClassificationID: perkotaan.ID},
							{Code: "3302720002", Name: "PASIR KIDUL", ClassificationID: perkotaan.ID},
							{Code: "3302720003", Name: "REJASARI", ClassificationID: perkotaan.ID},
							{Code: "3302720004", Name: "PASIRMUNCANG", ClassificationID: perkotaan.ID},
							{Code: "3302720005", Name: "BANTARSOKA", ClassificationID: perkotaan.ID},
							{Code: "3302720006", Name: "KOBER", ClassificationID: perkotaan.ID},
							{Code: "3302720007", Name: "KEDUNGWULUH", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Purwokerto Timur
					{
						Code: "3302730",
						Name: "PURWOKERTO TIMUR",
						Villages: []region.Village{
							{Code: "3302730001", Name: "SOKANEGARA", ClassificationID: perkotaan.ID},
							{Code: "3302730002", Name: "KRANJI", ClassificationID: perkotaan.ID},
							{Code: "3302730003", Name: "PURWOKERTO LOR", ClassificationID: perkotaan.ID},
							{Code: "3302730004", Name: "PURWOKERTO WETAN", ClassificationID: perkotaan.ID},
							{Code: "3302730005", Name: "MERSI", ClassificationID: perkotaan.ID},
							{Code: "3302730006", Name: "ARCAWINANGUN", ClassificationID: perkotaan.ID},
						},
					},
					// Kecamatan Purwokerto Utara
					{
						Code: "3302740",
						Name: "PURWOKERTO UTARA",
						Villages: []region.Village{
							{Code: "3302740001", Name: "BOBOSAN", ClassificationID: perkotaan.ID},
							{Code: "3302740002", Name: "PURWANEGARA", ClassificationID: perkotaan.ID},
							{Code: "3302740003", Name: "BANCARKEMBAR", ClassificationID: perkotaan.ID},
							{Code: "3302740004", Name: "SUMAMPIR", ClassificationID: perkotaan.ID},
							{Code: "3302740005", Name: "PABUWARAN", ClassificationID: perkotaan.ID},
							{Code: "3302740006", Name: "GRENDENG", ClassificationID: perkotaan.ID},
							{Code: "3302740007", Name: "KARANGWANGKAL", ClassificationID: perkotaan.ID},
						},
					},
				},
			},
		},
	}

	// 1. Buat atau Cari Provinsi
	provinceDB := region.Province{Code: provinceData.Code, Name: provinceData.Name}
	if err := tx.FirstOrCreate(&provinceDB, region.Province{Code: provinceDB.Code}).Error; err != nil {
		log.Printf("[DB] [SEED] [REGION] Error on Province %s: %v", provinceDB.Name, err)
		return err
	}

	// 2. Iterasi data Regencies (dari struct provinceData)
	for _, regData := range provinceData.Regencies {
		regDB := region.Regency{
			Code:       regData.Code,
			Name:       regData.Name,
			ProvinceID: provinceDB.ID,
		}
		// Cari atau buat Regency
		if err := tx.FirstOrCreate(&regDB, region.Regency{Code: regDB.Code}).Error; err != nil {
			log.Printf("[DB] [SEED] [REGION] Error on Regency %s: %v", regDB.Name, err)
			continue // Lanjut ke regency berikutnya jika error
		}

		// 3. Iterasi data Districts (dari struct regData)
		for _, distData := range regData.Districts {
			distDB := region.District{
				Code:      distData.Code,
				Name:      distData.Name,
				RegencyID: regDB.ID,
			}
			// Cari atau buat District
			if err := tx.FirstOrCreate(&distDB, region.District{Code: distDB.Code}).Error; err != nil {
				log.Printf("[DB] [SEED] [REGION] Error on District %s: %v", distDB.Name, err)
				continue // Lanjut ke district berikutnya jika error
			}

			// 4. Cek apakah villages untuk district ini sudah ada
			var count int64
			tx.Model(&region.Village{}).Where("district_id = ?", distDB.ID).Count(&count)

			if count == 0 {
				// Belum ada, lakukan "1 query" (bulk insert)
				villagesToCreate := distData.Villages
				// Set Foreign Key untuk semua village di batch ini
				for i := range villagesToCreate {
					villagesToCreate[i].DistrictID = distDB.ID
					// ClassificationID sudah di-set saat struct dibuat
				}

				// Lakukan CreateInBatches
				if err := tx.CreateInBatches(villagesToCreate, 50).Error; err != nil {
					// Jika gagal di sini, ini adalah error serius
					log.Printf("[DB] [SEED] [REGION] FAILED to bulk insert villages for District %s: %v", distDB.Name, err)
				} else {
					log.Printf("[DB] [SEED] [REGION] [BANYUMAS] Bulk inserted %d villages for District %s", len(villagesToCreate), distDB.Name)
				}
			}
			// Jika count > 0, desa sudah ada, jadi lewati (idempotent)
		}
	}

	log.Println("[DB] [SEED] [REGION] [BANYUMAS] Banyumas region seeded successfully.")
	return nil
}
