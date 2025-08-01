package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/models"
	"ipincamp/srikandi-sehat/src/routes"
	"ipincamp/srikandi-sehat/src/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func cleanupExpiredTokens() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running expired token cleanup...")
		now := time.Now()
		result := database.DB.Where("expires_at < ?", now).Delete(&models.InvalidToken{})
		if result.Error != nil {
			log.Printf("Failed to clean up expired tokens: %v", result.Error)
		} else {
			log.Printf("%d expired tokens have been deleted.", result.RowsAffected)
		}
	}
}

func main() {
	config.LoadConfig()
	database.ConnectDB()
	utils.SetupValidator()

	go cleanupExpiredTokens()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Get("CORS_ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
	}))
	app.Use(logger.New())

	routes.SetupRoutes(app)

	port := config.Get("API_PORT")
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
