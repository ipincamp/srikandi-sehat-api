package utils

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/models/region"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	roleCache      map[string]models.Role
	userRolesCache = make(map[string][]string)

	provincesCache []dto.RegionResponse
	regenciesCache map[string][]dto.RegionResponse
	districtsCache map[string][]dto.RegionResponse
	villagesCache  map[string][]dto.RegionResponse

	// Maintenance Cache
	isMaintenanceActive bool
	maintenanceMessage  string
	whitelistedUserIDs  map[uint]struct{}
	maintenanceMutex    = &sync.RWMutex{}

	// Report Token Cache (BARU)
	reportTokenCache map[string]struct{}
	reportTokenMutex = &sync.RWMutex{}

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

	// Initialize Maintenance Cache
	ReloadMaintenanceStatus()
	ReloadMaintenanceWhitelist()
	log.Println("Maintenance status and whitelist cache initialized.")

	// Initialize Report Token Cache
	reportTokenCache = make(map[string]struct{})
	log.Println("Report token cache initialized.")
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

// ReloadMaintenanceStatus fetches the latest maintenance status from the DB, creating defaults if missing.
func ReloadMaintenanceStatus() {
	maintenanceMutex.Lock()
	defer maintenanceMutex.Unlock()

	// Gunakan FirstOrCreate untuk memastikan record 'maintenance_mode_active' selalu ada
	activeSetting := models.Setting{Key: "maintenance_mode_active"}
	// Mencari record dengan key 'maintenance_mode_active'.
	// Jika tidak ada, buat record baru dengan Key tersebut dan Value default "false".
	err := database.DB.Attrs(models.Setting{Value: "false"}).FirstOrCreate(&activeSetting).Error

	if err != nil {
		// Jika FirstOrCreate gagal (error DB selain record not found), log error serius
		log.Printf("ERROR: Failed to ensure maintenance_mode_active setting exists: %v. Defaulting to false.", err)
		isMaintenanceActive = false
	} else {
		// Jika berhasil (baik menemukan atau membuat), parse nilainya
		parsedValue, parseErr := strconv.ParseBool(activeSetting.Value)
		if parseErr != nil {
			// Jika nilai di DB tidak valid (bukan true/false), log warning dan default ke false
			log.Printf("Warning: Invalid boolean value '%s' for maintenance_mode_active in database. Defaulting to false.", activeSetting.Value)
			isMaintenanceActive = false
		} else {
			isMaintenanceActive = parsedValue
		}
	}

	// Lakukan hal yang sama untuk message (FirstOrCreate dengan default message)
	messageSetting := models.Setting{Key: "maintenance_message"}
	defaultMessage := "Server is currently under maintenance. Please try again later."
	// Mencari record 'maintenance_message'.
	// Jika tidak ada, buat dengan Value defaultMessage.
	err = database.DB.Attrs(models.Setting{Value: defaultMessage}).FirstOrCreate(&messageSetting).Error

	if err != nil {
		// Jika FirstOrCreate gagal, log error dan gunakan default message
		log.Printf("ERROR: Failed to ensure maintenance_message setting exists: %v. Using default message.", err)
		maintenanceMessage = defaultMessage
	} else {
		// Jika berhasil, gunakan nilai dari DB (yang mungkin baru saja dibuat defaultnya)
		maintenanceMessage = messageSetting.Value
	}

	log.Printf("Maintenance status cache reloaded: Active=%v, Message=%s", isMaintenanceActive, maintenanceMessage)
}

// ReloadMaintenanceWhitelist fetches the latest whitelisted user IDs from the DB.
func ReloadMaintenanceWhitelist() {
	maintenanceMutex.Lock()
	defer maintenanceMutex.Unlock()

	var whitelistEntries []models.MaintenanceWhitelist
	err := database.DB.Find(&whitelistEntries).Error
	if err != nil {
		log.Printf("Error reloading maintenance whitelist: %v", err)
		whitelistedUserIDs = make(map[uint]struct{}) // Reset on error
		return
	}

	newWhitelist := make(map[uint]struct{}, len(whitelistEntries))
	for _, entry := range whitelistEntries {
		newWhitelist[entry.UserID] = struct{}{}
	}
	whitelistedUserIDs = newWhitelist
	log.Printf("Maintenance whitelist cache reloaded: %d users whitelisted.", len(whitelistedUserIDs))
}

// GetMaintenanceStatus returns the cached maintenance status and message.
func GetMaintenanceStatus() (bool, string) {
	maintenanceMutex.RLock()
	defer maintenanceMutex.RUnlock()
	return isMaintenanceActive, maintenanceMessage
}

// IsUserWhitelisted checks if a user UUID is in the cached whitelist.
func IsUserWhitelisted(userUUID string) bool {
	maintenanceMutex.RLock()
	defer maintenanceMutex.RUnlock()

	// Need to get user ID from UUID first (could be optimized with another cache if needed)
	var user models.User
	err := database.DB.Select("id").First(&user, "uuid = ?", userUUID).Error
	if err != nil {
		return false // User not found, definitely not whitelisted
	}

	_, exists := whitelistedUserIDs[user.ID]
	return exists
}

// --- Report Token Functions ---

// StoreReportToken menyimpan token unik ke cache dan mengatur masa kedaluwarsa.
func StoreReportToken(token string, expiration time.Duration) {
	reportTokenMutex.Lock()
	reportTokenCache[token] = struct{}{}
	reportTokenMutex.Unlock()

	// Menjadwalkan penghapusan token setelah kedaluwarsa
	time.AfterFunc(expiration, func() {
		reportTokenMutex.Lock()
		delete(reportTokenCache, token)
		reportTokenMutex.Unlock()
		log.Printf("Report token %s expired and was deleted from cache.", token)
	})
}

// UseReportToken mencoba menggunakan token.
// Jika token ada, token akan dihapus (digunakan) dan mengembalikan true.
// Jika token tidak ada, mengembalikan false.
func UseReportToken(token string) bool {
	reportTokenMutex.Lock()
	defer reportTokenMutex.Unlock()

	if _, found := reportTokenCache[token]; found {
		// Token ditemukan, hapus (gunakan) dan kembalikan true
		delete(reportTokenCache, token)
		return true
	}

	// Token tidak ditemukan (sudah digunakan atau kedaluwarsa)
	return false
}
