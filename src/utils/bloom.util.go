package utils

import (
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"log"
	"os"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
)

const bloomFilterFilePath = "database/email_filter.bin"

var emailFilter *bloom.BloomFilter

var mutex = &sync.RWMutex{}

func InitializeBloomFilter() {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(bloomFilterFilePath)
	if err == nil {
		defer file.Close()
		emailFilter = bloom.New(1, 1)
		if _, err := emailFilter.ReadFrom(file); err != nil {
			log.Fatalf("Failed to load Bloom Filter from file: %v", err)
		}
		log.Println("Bloom Filter successfully loaded from file.")
		return
	}

	log.Println("Bloom Filter file not found. Creating new filter from database...")

	emailFilter = bloom.NewWithEstimates(100000, 0.001)

	var emails []string
	database.DB.Model(&models.User{}).Pluck("email", &emails)

	for _, email := range emails {
		emailFilter.AddString(email)
	}

	log.Printf("%d existing emails have been added to the Bloom Filter.", len(emails))
	saveBloomFilter()
}

func saveBloomFilter() {
	file, err := os.Create(bloomFilterFilePath)
	if err != nil {
		log.Printf("ERROR: Failed to create Bloom Filter file: %v", err)
		return
	}
	defer file.Close()

	if _, err := emailFilter.WriteTo(file); err != nil {
		log.Printf("ERROR: Failed to save Bloom Filter to file: %v", err)
	}
}

func AddEmailToFilter(email string) {
	mutex.Lock()
	defer mutex.Unlock()

	emailFilter.AddString(email)
	saveBloomFilter()
}

func CheckEmailExists(email string) bool {
	mutex.RLock()
	defer mutex.RUnlock()

	return emailFilter.TestString(email)
}
