package utils

import (
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"log"
	"sync"
	"time"
)

var blocklistCache = make(map[string]struct{})
var blocklistMutex = &sync.RWMutex{}

func InitializeBlocklistCache() {
	var invalidTokens []models.InvalidToken
	if err := database.DB.Where("expires_at > ?", time.Now()).Find(&invalidTokens).Error; err != nil {
		log.Fatalf("Failed to load blocklisted tokens into cache: %v", err)
	}

	blocklistMutex.Lock()
	defer blocklistMutex.Unlock()
	for _, t := range invalidTokens {
		blocklistCache[t.Token] = struct{}{}
	}

	log.Printf("%d blocked tokens successfully loaded into cache.", len(invalidTokens))
}

func AddToBlocklistCache(token string, duration time.Duration) {
	blocklistMutex.Lock()
	blocklistCache[token] = struct{}{}
	blocklistMutex.Unlock()

	time.AfterFunc(duration, func() {
		blocklistMutex.Lock()
		delete(blocklistCache, token)
		blocklistMutex.Unlock()
	})
}

func IsTokenBlocked(token string) bool {
	blocklistMutex.RLock()
	defer blocklistMutex.RUnlock()
	_, found := blocklistCache[token]
	return found
}
