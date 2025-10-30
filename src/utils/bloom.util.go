package utils

import (
	"fmt"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"log"
	"os"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
)

const registrationFilterFilePath = "database/email_registration_filter.bin"
const frequentLoginFilterFilePath = "database/frequent_login_filter.bin"

var registrationEmailFilter *bloom.BloomFilter
var registrationMutex = &sync.RWMutex{}

var frequentLoginFilter *bloom.BloomFilter
var frequentLoginMutex = &sync.RWMutex{}

func InitializeRegistrationFilter() {
	registrationMutex.Lock()
	defer registrationMutex.Unlock()

	file, err := os.Open(registrationFilterFilePath)
	if err == nil {
		defer file.Close()
		registrationEmailFilter = bloom.New(1, 1)
		if _, err := registrationEmailFilter.ReadFrom(file); err != nil {
			log.Printf("ERROR: Failed to read Registration Bloom Filter from file: %v", err)
			return
		}
		log.Println("Registration Bloom Filter loaded from file.")
		return
	}

	log.Println("Creating new Registration Bloom Filter...")
	registrationEmailFilter = bloom.NewWithEstimates(100000, 0.001)

	var emails []string
	database.DB.Model(&models.User{}).Pluck("email", &emails)

	for _, email := range emails {
		registrationEmailFilter.AddString(email)
	}

	log.Printf("%d existing emails have been added to the registration filter.", len(emails))
	saveRegistrationFilter()
}

func saveRegistrationFilter() {
	file, err := os.Create(registrationFilterFilePath)
	if err != nil {
		log.Printf("ERROR: Failed to create Registration Filter file: %v", err)
		return
	}
	defer file.Close()
	if _, err := registrationEmailFilter.WriteTo(file); err != nil {
		log.Printf("ERROR: Failed to save Registration Filter to file: %v", err)
	}
}

func AddEmailToRegistrationFilter(email string) {
	registrationMutex.Lock()
	defer registrationMutex.Unlock()
	registrationEmailFilter.AddString(email)
	// saveRegistrationFilter()
}

func CheckEmailExistsInRegistrationFilter(email string) bool {
	registrationMutex.RLock()
	defer registrationMutex.RUnlock()
	return registrationEmailFilter.TestString(email)
}

func InitializeFrequentLoginFilter() {
	frequentLoginMutex.Lock()
	defer frequentLoginMutex.Unlock()

	file, err := os.Open(frequentLoginFilterFilePath)
	if err == nil {
		defer file.Close()
		frequentLoginFilter = bloom.New(1, 1)
		if _, err := frequentLoginFilter.ReadFrom(file); err != nil {
			log.Fatalf("FATAL: Failed to read Frequent Login Bloom Filter from file: %v", err)
		}
		log.Println("Frequent Login Bloom Filter loaded from file.")
		return
	}

	log.Println("Creating new Frequent Login Bloom Filter...")
	frequentLoginFilter = bloom.NewWithEstimates(10000, 0.01)
	// saveFrequentLoginFilter()
}

func AddUserToFrequentLoginFilter(user models.User) {
	frequentLoginMutex.Lock()
	defer frequentLoginMutex.Unlock()

	entry := fmt.Sprintf("%s:%s", user.UUID, user.Email)
	frequentLoginFilter.AddString(entry)
	// saveFrequentLoginFilter()
}

func SaveAllBloomFilters() {
	log.Println("Saving all Bloom Filters...")

	registrationMutex.Lock()
	defer registrationMutex.Unlock()
	regFile, err := os.Create(registrationFilterFilePath)
	if err != nil {
		log.Printf("ERROR: Failed to create file Filter Registrasi: %v", err)
	} else {
		defer regFile.Close()
		if _, err := registrationEmailFilter.WriteTo(regFile); err != nil {
			log.Printf("ERROR: Failed to save Filter Registrasi to file: %v", err)
		}
	}

	frequentLoginMutex.Lock()
	defer frequentLoginMutex.Unlock()
	loginFile, err := os.Create(frequentLoginFilterFilePath)
	if err != nil {
		log.Printf("ERROR: Failed to create file Filter Login Sering: %v", err)
	} else {
		defer loginFile.Close()
		if _, err := frequentLoginFilter.WriteTo(loginFile); err != nil {
			log.Printf("ERROR: Failed to save Filter Login Sering to file: %v", err)
		}
	}

	log.Println("Saving Bloom Filters completed.")
}
