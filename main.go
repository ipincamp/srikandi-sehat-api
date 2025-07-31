package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	app := fiber.New()
	app.Use(logger.New())

	routes.SetupRoutes(app)

	port := config.Get("API_PORT")
	log.Fatal(app.Listen(":" + port))
}
