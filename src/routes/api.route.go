package routes

import (
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/handlers"
	"ipincamp/srikandi-sehat/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateBody[dto.RegisterRequest], handlers.Register)
	auth.Post("/login", middleware.ValidateBody[dto.LoginRequest], handlers.Login)
	auth.Post("/logout", middleware.AuthMiddleware, handlers.Logout)

	// User routes
	user := api.Group("/me", middleware.AuthMiddleware)
	user.Get("/", handlers.GetMyProfile)
	user.Put("/details", middleware.ValidateBody[dto.UpdateProfileRequest], handlers.UpdateOrCreateProfile)
	user.Patch("/password", middleware.ValidateBody[dto.ChangePasswordRequest], handlers.ChangeMyPassword)

	// Admin routes
	admin := api.Group("/admin", middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.Get("/users", middleware.ValidateQuery[dto.UserQuery], handlers.GetAllUsers)
	admin.Get("/users/:id", handlers.GetUserByID)

	// Region routes
	region := api.Group("/regions")
	region.Get("/provinces", handlers.GetAllProvinces)
	region.Get("/regencies", middleware.ValidateQuery[dto.RegencyQuery], handlers.GetRegenciesByProvince)
	region.Get("/districts", middleware.ValidateQuery[dto.DistrictQuery], handlers.GetDistrictsByRegency)
	region.Get("/villages", middleware.ValidateQuery[dto.VillageQuery], handlers.GetVillagesByDistrict)
}
