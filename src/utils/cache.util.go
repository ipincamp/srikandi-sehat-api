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
