package routes

import (
	"ipincamp/srikandi-sehat/src/handlers"
	"ipincamp/srikandi-sehat/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Get("/profile", middleware.AuthMiddleware, handlers.Profile)
	auth.Post("/logout", middleware.AuthMiddleware, handlers.Logout)
}
