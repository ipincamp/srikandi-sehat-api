package routes

import (
	"ipincamp/srikandi-sehat/handlers"
	"ipincamp/srikandi-sehat/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes
	publicAuth := api.Group("/auth")
	publicAuth.Post("/register", handlers.Register)
	publicAuth.Post("/login", handlers.Login)

	protectedAuth := api.Group("/auth", middleware.AuthMiddleware)
	protectedAuth.Post("/logout", handlers.Logout)
	protectedAuth.Get("/me", handlers.Profile)
}
