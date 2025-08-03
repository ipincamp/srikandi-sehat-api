package utils

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"log"
	"sync"
)

var roleCache map[string]models.Role
var cacheMutex = &sync.RWMutex{}

var userRolesCache = make(map[string][]string)
var userCacheMutex = &sync.RWMutex{}

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

	log.Printf("%d roles successfully loaded into cache.", len(roles))
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
	userCacheMutex.RLock()
	roles, found := userRolesCache[userUUID]
	userCacheMutex.RUnlock()

	if found {
		return roles, nil
	}

	userCacheMutex.Lock()
	defer userCacheMutex.Unlock()

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
	userCacheMutex.Lock()
	defer userCacheMutex.Unlock()
	delete(userRolesCache, userUUID)
}
