package main

import (
	"context"
	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/database"
	"ipincamp/srikandi-sehat/src/middleware"
	"ipincamp/srikandi-sehat/src/routes"
	"ipincamp/srikandi-sehat/src/utils"
	"ipincamp/srikandi-sehat/src/workers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/robfig/cron/v3"
)

func main() {
	config.LoadConfig()
	config.SetTimeZone()

	utils.InitLogger()

	c := cron.New()
	env := config.Get("APP_ENV")
	if env == "production" {
		log.Println("Running in production mode. Scheduling cron jobs accordingly.")
		c.AddFunc("0 5 * * *", workers.CheckLongMenstrualCycles) // setiap jam 05:00 pagi
		c.AddFunc("0 5 * * *", workers.CheckLateMenstrualCycles) // setiap jam 05:00 pagi
		log.Println("Scheduled cron jobs for production at 05:00 AM daily.")
	} else {
		log.Println("Running in development mode. Scheduling cron jobs for testing.")
		c.AddFunc("@every 1m", workers.CheckLongMenstrualCycles) // setiap 1 menit (testing)
		c.AddFunc("@every 1m", workers.CheckLateMenstrualCycles) // setiap 1 menit (testing)
		log.Println("Scheduled cron jobs for development every 1 minute.")
	}
	c.Start()
	log.Println("Cron job for cycle checking has been scheduled.")
	defer c.Stop()

	utils.InitFCM()
	database.ConnectDB()

	utils.SetupValidator()
	utils.InitializeRegistrationFilter()
	utils.InitializeFrequentLoginFilter()
	utils.InitializeRoleCache()
	utils.InitializeBlocklistCache()

	go utils.CleanupExpiredTokens()

	app := fiber.New(fiber.Config{
		Prefork:        false,
		ServerHeader:   "SrikandiSehat",
		TrustedProxies: []string{config.Get("TRUSTED_PROXIES")},
	})
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Header X-Real-IP: %s", c.Get("X-Real-IP"))
		log.Printf("Header X-Forwarded-For: %s", c.Get("X-Forwarded-For"))
		return c.Next()
	})
	app.Use(middleware.RecoverMiddleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Get("CORS_ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
	}))
	app.Use(logger.New())
	app.Use(middleware.MaintenanceMiddleware())

	routes.SetupRoutes(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := config.Get("API_PORT")
		utils.InfoLogger.Printf("Server is starting on port %s", port)
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			utils.ErrorLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("[APP] Shutting down server...")

	utils.SaveAllBloomFilters()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("[APP] Failed to gracefully shutdown server: %v", err)
	}
}
