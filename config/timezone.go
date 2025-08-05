package config

import (
	"log"
	"time"
)

func SetTimeZone() {
	jakartaTime, err := time.LoadLocation(Get("TIMEZONE"))
	if err != nil {
		log.Fatalf("[APP] Failed to load %s timezone: %v", Get("TIMEZONE"), err)
	}
	time.Local = jakartaTime
}
