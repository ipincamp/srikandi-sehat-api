package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// --- Handler dari gqlgen ---
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	// --- Dependensi Internal ---
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/logger"

	// --- Impor paket GraphQL Anda ---
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/generated"
)

func main() {
	// --- 1. Load Configuration ---
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		// We can't use the configured logger yet, so we use a temporary one.
		// Use a basic zerolog for this fatal error.
		tempLogger := logger.NewLogger("development") // Use "development" for a clean exit message
		tempLogger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// --- 1.5. Initialize Logger ---
	// Now that config is loaded, create the main logger
	// This logger will be used throughout the application
	log := logger.NewLogger(cfg.Server.Env)
	log.Info().Str("Env", cfg.Server.Env).Msg("Configuration loaded")

	// --- 2. Setup Application Context ---
	// Buat context yang me-listen untuk sinyal shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen untuk sinyal interrupt dari OS
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info().Msg("Shutdown signal received, initiating graceful shutdown...")
		cancel()
	}()

	// --- 3. Initialize Driven Adapters (Database) ---
	// Terhubung ke database PostgreSQL
	dbPool, err := postgres.Connect(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	// Tunda penutupan pool koneksi hingga aplikasi keluar
	defer dbPool.Close()
	log.Info().Msg("Database connection pool established")

	// --- 4. Dependency Injection (Composition Root) ---
	// Di sinilah Anda menginisialisasi repositories, services, dan handlers.
	// Ini adalah inti dari "Dependency Injection": meneruskan dependensi (seperti dbPool)
	// ke komponen yang membutuhkannya.

	// 4a. Inisialisasi Repositories (Contoh)
	// userRepo := postgres.NewUserRepository(dbPool)

	// 4b. Inisialisasi Core Services (Contoh)
	// userService := service.NewUserService(userRepo)

	// 4c. Inisialisasi Driving Adapters (GraphQL)
	// "Suntikkan" service ke dalam resolver.
	// Karena kita belum memiliki service untuk tes ini, kita panggil langsung.
	// (Contoh dengan service: gqlResolver := graphql.NewResolver(userService))
	gqlResolver := graphql.NewResolver()

	// 4d. Buat konfigurasi server GraphQL
	gqlConfig := generated.Config{Resolvers: gqlResolver}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))

	// --- 5. Start Application (HTTP Server) ---
	log.Info().Msg("Application dependencies initialized.")

	// Buat HTTP server mux (router)
	httpMux := http.NewServeMux()

	// Atur handler untuk GraphQL Playground di root ("/")
	httpMux.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	// Atur handler untuk endpoint GraphQL utama di "/query"
	httpMux.Handle("/query", gqlServer)

	// Konfigurasi server HTTP
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: httpMux,
		// Anda bisa menambahkan ReadTimeout, WriteTimeout, dll di sini untuk produksi
	}

	// Jalankan server di goroutine terpisah
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Starting HTTP server... (GraphQL Playground at http://localhost:" + cfg.Server.Port + ")")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Could not start server")
		}
	}()

	// --- 6. Wait for Shutdown Signal ---
	// Blok di sini sampai context dibatalkan (misalnya, oleh sinyal OS)
	<-ctx.Done()

	// --- 7. Graceful Shutdown ---
	log.Info().Msg("Shutting down application...")

	// Buat context untuk server shutdown dengan timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Lakukan graceful shutdown server HTTP
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Warn().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Application shut down gracefully.")
}
