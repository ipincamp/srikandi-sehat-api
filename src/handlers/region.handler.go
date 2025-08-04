package handlers

import (
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllProvinces(c *fiber.Ctx) error {
	provinces := utils.GetProvincesFromCache()
	if len(provinces) == 0 {
		return utils.SendError(c, fiber.StatusNotFound, "No provinces found")
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Provinces fetched successfully", provinces)
}

func GetRegenciesByProvince(c *fiber.Ctx) error {
	queries := c.Locals("request_queries").(*dto.RegencyQuery)

	regencies, found := utils.GetRegenciesByProvinceCodeFromCache(queries.ProvinceCode)
	if !found {
		return utils.SendError(c, fiber.StatusNotFound, "No regencies found for the given province code")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Regencies fetched successfully", regencies)
}

func GetDistrictsByRegency(c *fiber.Ctx) error {
	queries := c.Locals("request_queries").(*dto.DistrictQuery)

	districts, found := utils.GetDistrictsByRegencyCodeFromCache(queries.RegencyCode)
	if !found {
		return utils.SendError(c, fiber.StatusNotFound, "No districts found for the given regency code")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Districts fetched successfully", districts)
}

func GetVillagesByDistrict(c *fiber.Ctx) error {
	queries := c.Locals("request_queries").(*dto.VillageQuery)

	villages, found := utils.GetVillagesByDistrictCodeFromCache(queries.DistrictCode)
	if !found {
		return utils.SendError(c, fiber.StatusNotFound, "No villages found for the given district code")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Villages fetched successfully", villages)
}
