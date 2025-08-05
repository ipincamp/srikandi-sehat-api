package factories

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models/region"
	"log"
	"math/rand"
	"strconv"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func MakeProvince() (region.Province, error) {
	provinceName := fmt.Sprintf("Provinsi %s", faker.Word())
	provinceCode := strconv.Itoa(rand.Intn(90) + 10)

	return region.Province{
		Code: provinceCode,
		Name: provinceName,
	}, nil
}

func MakeRegency(province region.Province) (region.Regency, error) {
	regencyName := fmt.Sprintf("Kabupaten %s", faker.Word())
	regencyCode := fmt.Sprintf("%s%02d", province.Code, rand.Intn(99)+1)

	return region.Regency{
		ProvinceID: province.ID,
		Code:       regencyCode,
		Name:       regencyName,
	}, nil
}

func MakeDistrict(regency region.Regency) (region.District, error) {
	districtName := fmt.Sprintf("Kecamatan %s", faker.Word())
	districtCode := fmt.Sprintf("%s%03d", regency.Code, rand.Intn(999)+1)

	return region.District{
		RegencyID: regency.ID,
		Code:      districtCode,
		Name:      districtName,
	}, nil
}

func MakeVillage(district region.District, classification region.Classification) (region.Village, error) {
	villageName := fmt.Sprintf("Desa %s", faker.Word())
	villageCode := fmt.Sprintf("%s%04d", district.Code, rand.Intn(9999)+1)

	return region.Village{
		DistrictID:       district.ID,
		ClassificationID: classification.ID,
		Code:             villageCode,
		Name:             villageName,
	}, nil
}

func CreateFullAddressSet(tx *gorm.DB) (region.Village, error) {
	province, err := MakeProvince()
	if err != nil {
		return region.Village{}, err
	}
	if err := tx.Create(&province).Error; err != nil {
		return region.Village{}, err
	}

	regency, err := MakeRegency(province)
	if err != nil {
		return region.Village{}, err
	}
	if err := tx.Create(&regency).Error; err != nil {
		return region.Village{}, err
	}

	district, err := MakeDistrict(regency)
	if err != nil {
		return region.Village{}, err
	}
	if err := tx.Create(&district).Error; err != nil {
		return region.Village{}, err
	}

	var classifications []region.Classification
	if err := database.DB.Find(&classifications).Error; err != nil || len(classifications) == 0 {
		log.Fatal("[DB] [FACTORY] [REGION] Classification not found. Please run the classification seeder first.")
	}
	randomClassification := classifications[rand.Intn(len(classifications))]

	village, err := MakeVillage(district, randomClassification)
	if err != nil {
		return region.Village{}, err
	}
	if err := tx.Create(&village).Error; err != nil {
		return region.Village{}, err
	}

	log.Printf("[DB] [FACTORY] [REGION] A complete address set has been created: %s, %s, %s, %s",
		village.Name, district.Name, regency.Name, province.Name)

	return village, nil
}
