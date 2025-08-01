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
	auth.Post("/logout", middleware.AuthMiddleware, handlers.Logout)

	// User routes
	user := api.Group("/me", middleware.AuthMiddleware)
	user.Get("/", handlers.Profile)
	user.Patch("/details", handlers.UpdateDetails)
	user.Patch("/password", handlers.ChangePassword)
}
