package main

import (
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/routes"
	"ipincamp/srikandi-sehat/src/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()
	utils.SetupValidator()

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
