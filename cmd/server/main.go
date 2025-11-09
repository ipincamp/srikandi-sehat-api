package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// --- Internal Dependencies ---
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	// "github.com/ipincamp/srikandi-sehat/internal/adapters/primary/http" (Example)
	// "github.com/ipincamp/srikandi-sehat/internal/core/services"       (Example)
	// "github.com/ipincamp/srikandi-sehat/internal/repositories"        (Example)
)

func main() {
	// --- 1. Load Configuration ---
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded (Env: %s)", cfg.Server.Env)

	// --- 2. Setup Application Context ---
	// Create a context that listens for shutdown signals.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS interrupt signals
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Println("Shutdown signal received, initiating graceful shutdown...")
		cancel()
	}()

	// --- 3. Initialize Driven Adapters (Database) ---
	// Connect to the PostgreSQL database
	dbPool, err := postgres.Connect(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Defer closing the pool until the application exits
	defer dbPool.Close()
	log.Println("Database connection pool established")

	// --- 4. Dependency Injection (Composition Root) ---
	// Here you would initialize your repositories, services, and handlers (primary adapters).
	// This is the core of "Dependency Injection": passing dependencies (like dbPool)
	// into the components that need them.
	//
	// Example:
	// userRepo := repositories.NewGormUserRepository(gormDB) // If using Gorm
	// userRepo := repositories.NewPgxUserRepository(dbPool)   // If using pgx directly
	//
	// authService := services.NewAuthService(userRepo, cfg.Server.JWTSecret)
	//
	// httpHandler := http.NewHandler(authService)
	//
	// server := &http.Server{
	// 	Addr:    ":" + cfg.Server.Port,
	// 	Handler: httpHandler.Router,
	// }

	// --- 5. Start Application ---
	//
	// Example:
	// go func() {
	// 	log.Printf("Starting server on port %s", cfg.Server.Port)
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("Could not start server: %v", err)
	// 	}
	// }()
	//
	// --- TODO: Remove this placeholder and start your server ---
	log.Println("Application dependencies initialized.")
	log.Println("TODO: Start your HTTP server or primary adapter here.")
	// -------------------------------------------------------------

	// --- 6. Wait for Shutdown Signal ---
	// Block here until the context is canceled (e.g., by the OS signal)
	<-ctx.Done()

	// --- 7. Graceful Shutdown ---
	log.Println("Shutting down application...")

	// --- TODO: Add your graceful shutdown logic here ---
	// (e.g., shut down the HTTP server with a timeout)
	//
	// shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer shutdownCancel()
	//
	// if err := server.Shutdown(shutdownCtx); err != nil {
	// 	log.Printf("Server shutdown failed: %v", err)
	// }
	// ---------------------------------------------------

	log.Println("Application shut down gracefully.")
}
