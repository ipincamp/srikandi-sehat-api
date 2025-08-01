package handlers

import (
	"errors"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models/region"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllProvinces(c *fiber.Ctx) error {
	var provinces []region.Province
	database.DB.Find(&provinces)

	var responseData []dto.RegionResponse
	for _, p := range provinces {
		responseData = append(responseData, dto.RegionResponse{Code: p.Code, Name: p.Name})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Provinces fetched successfully", responseData)
}

func GetRegenciesByProvince(c *fiber.Ctx) error {
	provinceCode := c.Query("province_code")
	if provinceCode == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Query 'province_code' is required")
	}

	var province region.Province
	if err := database.DB.First(&province, "code = ?", provinceCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "Province not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	var regencies []region.Regency
	database.DB.Where("province_id = ?", province.ID).Find(&regencies)

	var responseData []dto.RegionResponse
	for _, r := range regencies {
		responseData = append(responseData, dto.RegionResponse{Code: r.Code, Name: r.Name})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Regencies fetched successfully", responseData)
}

func GetDistrictsByRegency(c *fiber.Ctx) error {
	regencyCode := c.Query("regency_code")
	if regencyCode == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Query 'regency_code' is required")
	}

	var regency region.Regency
	if err := database.DB.First(&regency, "code = ?", regencyCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "Regency not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	var districts []region.District
	database.DB.Where("regency_id = ?", regency.ID).Find(&districts)

	var responseData []dto.RegionResponse
	for _, d := range districts {
		responseData = append(responseData, dto.RegionResponse{Code: d.Code, Name: d.Name})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Districts fetched successfully", responseData)
}

func GetVillagesByDistrict(c *fiber.Ctx) error {
	districtCode := c.Query("district_code")
	if districtCode == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Query 'district_code' is required")
	}

	var district region.District
	if err := database.DB.First(&district, "code = ?", districtCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, fiber.StatusNotFound, "District not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Database error")
	}

	var villages []region.Village
	database.DB.Where("district_id = ?", district.ID).Find(&villages)

	var responseData []dto.RegionResponse
	for _, v := range villages {
		responseData = append(responseData, dto.RegionResponse{Code: v.Code, Name: v.Name})
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Villages fetched successfully", responseData)
}
