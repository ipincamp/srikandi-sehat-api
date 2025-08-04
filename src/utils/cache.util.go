package utils

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/region"
	"log"
	"sync"
)

var (
	roleCache      map[string]models.Role
	userRolesCache = make(map[string][]string)

	provincesCache []dto.RegionResponse
	regenciesCache map[string][]dto.RegionResponse
	districtsCache map[string][]dto.RegionResponse
	villagesCache  map[string][]dto.RegionResponse

	cacheMutex = &sync.RWMutex{}
)

func InitializeRoleCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	roleCache = make(map[string]models.Role)

	var roles []models.Role
	if err := database.DB.Find(&roles).Error; err != nil {
		log.Fatalf("Failed to load roles into cache: %v", err)
	}

	for _, role := range roles {
		roleCache[role.Name] = role
	}

	log.Printf("Role cache initialized with %d roles", len(roleCache))

	var provinces []region.Province
	database.DB.Find(&provinces)
	provincesCache = make([]dto.RegionResponse, len(provinces))
	for i, p := range provinces {
		provincesCache[i] = dto.RegionResponse{Code: p.Code, Name: p.Name}
	}

	var regencies []region.Regency
	database.DB.Find(&regencies)
	regenciesCache = make(map[string][]dto.RegionResponse)
	for _, r := range regencies {
		provinceCode := getParentCode(r.Code)
		regenciesCache[provinceCode] = append(regenciesCache[provinceCode], dto.RegionResponse{Code: r.Code, Name: r.Name})
	}

	var districts []region.District
	database.DB.Find(&districts)
	districtsCache = make(map[string][]dto.RegionResponse)
	for _, d := range districts {
		regencyCode := getParentCode(d.Code)
		districtsCache[regencyCode] = append(districtsCache[regencyCode], dto.RegionResponse{Code: d.Code, Name: d.Name})
	}

	var villages []region.Village
	database.DB.Preload("Classification").Find(&villages)
	villagesCache = make(map[string][]dto.RegionResponse)
	for _, v := range villages {
		districtCode := getParentCode(v.Code)
		villagesCache[districtCode] = append(villagesCache[districtCode], dto.RegionResponse{
			Code:           v.Code,
			Name:           v.Name,
			Classification: v.Classification.Name,
		})
	}

	log.Printf("Region data cached: %d provinces, %d regencies, %d districts, %d villages",
		len(provincesCache), len(regenciesCache), len(districtsCache), len(villagesCache))
}

func GetRoleByName(name string) (models.Role, error) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	role, found := roleCache[name]
	if !found {
		return models.Role{}, fmt.Errorf("role '%s' not found in cache", name)
	}
	return role, nil
}

func GetUserRoles(userUUID string) ([]string, error) {
	cacheMutex.RLock()
	roles, found := userRolesCache[userUUID]
	cacheMutex.RUnlock()

	if found {
		return roles, nil
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	roles, found = userRolesCache[userUUID]
	if found {
		return roles, nil
	}

	var user models.User
	if err := database.DB.Preload("Roles").First(&user, "uuid = ?", userUUID).Error; err != nil {
		return nil, err
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	userRolesCache[userUUID] = roleNames

	return roleNames, nil
}

// InvalidateUserRolesCache menghapus data role user dari cache.
/*
func AssignRoleToUser(c *fiber.Ctx) error {
	// ... (logika untuk mencari user dan role)

	// Setelah berhasil menetapkan role:
	database.DB.Model(&user).Association("Roles").Append(&role)

	// ===== INVALIDASI CACHE DI SINI =====
	utils.InvalidateUserRolesCache(user.UUID)
	// =====================================

	return utils.SendSuccess(c, fiber.StatusOK, "Role assigned successfully", nil)
}
*/
func InvalidateUserRolesCache(userUUID string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	delete(userRolesCache, userUUID)
}

func getParentCode(code string) string {
	switch len(code) {
	case 4: // Regency code (e.g., "3302"), return Province code
		return code[:2]
	case 7: // District code (e.g., "3302010"), return Regency code
		return code[:4]
	case 10: // Village code (e.g., "3302010001"), return District code
		return code[:7]
	default:
		return ""
	}
}

func GetProvincesFromCache() []dto.RegionResponse {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return provincesCache
}

func GetRegenciesByProvinceCodeFromCache(provinceCode string) ([]dto.RegionResponse, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	data, found := regenciesCache[provinceCode]
	return data, found
}

func GetDistrictsByRegencyCodeFromCache(regencyCode string) ([]dto.RegionResponse, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	data, found := districtsCache[regencyCode]
	return data, found
}

func GetVillagesByDistrictCodeFromCache(districtCode string) ([]dto.RegionResponse, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	data, found := villagesCache[districtCode]
	return data, found
}
