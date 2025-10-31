package routes

import (
	"ipincamp/srikandi-sehat/src/dto"
	"ipincamp/srikandi-sehat/src/handlers"
	menstrualHandler "ipincamp/srikandi-sehat/src/handlers/menstrual"
	"ipincamp/srikandi-sehat/src/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/health", handlers.HealthCheck)
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	registerLimiter := middleware.IPRateLimiter(50, 5*time.Minute)
	loginLimiter := middleware.LoginRateLimiter(5, 1*time.Minute)
	auth.Post("/register",
		registerLimiter,
		middleware.ValidateBody[dto.RegisterRequest],
		handlers.Register,
	)
	auth.Post("/login",
		loginLimiter,
		middleware.ValidateBody[dto.LoginRequest],
		handlers.Login,
	)
	auth.Post("/google",
		middleware.ValidateBody[dto.GoogleLoginRequest],
		handlers.LoginWithGoogle,
	)
	auth.Post("/logout", middleware.AuthMiddleware, handlers.Logout)
	auth.Post(
		"/verify-otp",
		middleware.AuthMiddleware,
		middleware.ValidateBody[dto.VerifyOTPRequest],
		handlers.VerifyOTP,
	)
	auth.Post(
		"/resend-verification",
		middleware.AuthMiddleware,
		handlers.ResendVerification,
	)

	// User routes
	user := api.Group("/me", middleware.AuthMiddleware)
	user.Get("/", handlers.GetMyProfile)
	user.Put("/details", middleware.ValidateBody[dto.UpdateProfileRequest], handlers.UpdateOrCreateProfile)
	user.Patch("/password", middleware.ValidateBody[dto.ChangePasswordRequest], handlers.ChangeMyPassword)
	user.Patch("/fcm-token", handlers.UpdateFcmToken)
	user.Post("/test-notification",
		middleware.ValidateBody[dto.TestNotificationRequest], // Validasi request body
		handlers.SendTestNotification,                        // Panggil handler baru
	)

	// Admin routes
	adminLimiter := middleware.UserRateLimiter(100, 1*time.Minute)
	admin := api.Group("/admin", middleware.AuthMiddleware, middleware.AdminMiddleware, adminLimiter)
	admin.Get("/users/statistics", handlers.GetUserStatistics)
	admin.Post("/reports/generate-csv-link", handlers.GenerateFullReportLink)
	admin.Get("/users", middleware.ValidateQuery[dto.UserQuery], handlers.GetAllUsers)
	admin.Get("/users/:id", middleware.ValidateParams[dto.UserParam], handlers.GetUserByID)

	// Maintenance Management Routes (Admin only)
	maintenance := admin.Group("/maintenance")
	maintenance.Get("/", handlers.GetMaintenanceStatus)
	maintenance.Post("/toggle", middleware.ValidateBody[dto.ToggleMaintenanceRequest], handlers.ToggleMaintenanceMode)
	maintenance.Get("/whitelist", handlers.GetWhitelistedUsers)
	maintenance.Post("/whitelist", middleware.ValidateBody[dto.WhitelistUserRequest], handlers.AddUserToWhitelist)
	maintenance.Delete("/whitelist", middleware.ValidateBody[dto.WhitelistUserRequest], handlers.RemoveUserFromWhitelist)

	// Region routes
	region := api.Group("/regions")
	region.Get("/provinces", handlers.GetAllProvinces)
	region.Get("/regencies", middleware.ValidateQuery[dto.RegencyQuery], handlers.GetRegenciesByProvince)
	region.Get("/districts", middleware.ValidateQuery[dto.DistrictQuery], handlers.GetDistrictsByRegency)
	region.Get("/villages", middleware.ValidateQuery[dto.VillageQuery], handlers.GetVillagesByDistrict)

	// Notification routes
	api.Get("/notifications", middleware.AuthMiddleware, handlers.GetNotificationHistory)
	api.Patch("/notifications/:id/read", middleware.AuthMiddleware, handlers.MarkNotificationAsRead)

	// Rute Unduhan Laporan
	api.Get("/reports/download/:token", handlers.DownloadFullReportByToken)

	// Menstrual health routes
	menstrual := api.Group("/menstrual", middleware.AuthMiddleware, middleware.VerifiedMiddleware)
	menstrual.Get("/cycles/status", menstrualHandler.GetCycleStatus)
	menstrual.Post("/cycles", middleware.ValidateBody[dto.CycleRequest], menstrualHandler.RecordCycle)
	menstrual.Get("/cycles", middleware.ValidateQuery[dto.PaginationQuery], menstrualHandler.GetCycleHistory)
	menstrual.Get("/cycles/:id", middleware.ValidateParams[dto.CycleParam], menstrualHandler.GetCycleByID)
	menstrual.Delete(
		"/cycles/:id",
		middleware.ValidateParams[dto.CycleParam],
		middleware.ValidateBody[dto.DeleteCycleRequest],
		menstrualHandler.DeleteCycleByID,
	)

	// Symptom specific routes
	menstrual.Post("/symptoms/log", middleware.ValidateBody[dto.SymptomLogRequest], menstrualHandler.LogSymptoms)
	menstrual.Get("/symptoms/master", menstrualHandler.GetSymptomsMaster)
	menstrual.Get("/symptoms/history", middleware.ValidateQuery[dto.SymptomHistoryQuery], menstrualHandler.GetSymptomHistory)
	menstrual.Get("/symptoms/log/:id", middleware.ValidateParams[dto.SymptomLogParam], menstrualHandler.GetSymptomLogByID)
	menstrual.Get("/recommendations", menstrualHandler.GetRecommendationsBySymptoms)
}
