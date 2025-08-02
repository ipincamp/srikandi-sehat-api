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
	user.Get("/", handlers.GetMyProfile)
	user.Put("/details", handlers.UpdateOrCreateProfile)
	user.Patch("/password", handlers.ChangeMyPassword)

	// Admin routes
	admin := api.Group("/admin", middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.Get("/users", handlers.GetAllUsers)
	admin.Get("/users/:id", handlers.GetUserByID)

	// Region routes
	region := api.Group("/regions")
	region.Get("/provinces", handlers.GetAllProvinces)
	region.Get("/regencies", handlers.GetRegenciesByProvince)
	region.Get("/districts", handlers.GetDistrictsByRegency)
	region.Get("/villages", handlers.GetVillagesByDistrict)
}
