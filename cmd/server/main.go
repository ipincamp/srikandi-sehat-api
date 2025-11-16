package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/generated"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/middleware"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/resolvers" // <-- IMPORT PORTS
	"github.com/ipincamp/srikandi-sehat/internal/core/service"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/logger"
	"github.com/ipincamp/srikandi-sehat/pkg/password"
	"github.com/ipincamp/srikandi-sehat/pkg/token"
)

func main() {
	// --- 1. Load Configuration ---
	cfg, err := config.Load()
	if err != nil {
		tempLogger := logger.NewLogger("development")
		tempLogger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// --- 1.5. Initialize Logger ---
	log := logger.NewLogger(cfg.Server.Env)
	log.Info().Str("Env", cfg.Server.Env).Msg("Configuration loaded")

	// --- 2. Setup Graceful Shutdown ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info().Msg("Shutdown signal received, initiating graceful shutdown...")
		cancel()
	}()

	// --- 3. Initialize Driven Adapters (Database) ---
	db, err := postgres.ConnectGORM(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database using GORM")
	}
	log.Info().Msg("Database connection established (GORM)")

	// --- 4. Dependency Injection (Merakit Arsitektur) ---

	// 4a. Initialize 'pkg' helpers (implementations)
	hasher := password.NewArgon2idHasher()
	tokenMaker, err := token.NewPasetoMaker(cfg.Token.SymmetricKey, cfg.Token.Issuer)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Paseto token maker")
	}

	// 4b. Inisialisasi Repositories (Driven Adapters)
	// This repo is for non-transactional reads (e.g., UserService.GetByID)
	userRepoLogger := log.With().Str("component", "UserRepository").Logger()
	userRepo := postgres.NewUserRepository(db, userRepoLogger)

	// 4c. Inisialisasi Unit of Work (for transactional mutations)
	uowLogger := log.With().Str("component", "UnitOfWork").Logger()
	uow := postgres.NewUnitOfWork(db, uowLogger)

	// 4d. Inisialisasi Core Services
	userServiceLogger := log.With().Str("component", "UserService").Logger()
	// UserService still uses the non-transactional repo for reads
	userService := service.NewUserService(userRepo, userServiceLogger)

	authServiceLogger := log.With().Str("component", "AuthService").Logger()
	// AuthService now uses the Unit of Work for its mutations
	authService := service.NewAuthService(
		uow,
		hasher,
		tokenMaker,
		cfg.Token,
		authServiceLogger,
	)

	// 4e. Inisialisasi GraphQL (Driving Adapter)
	resolverLogger := log.With().Str("component", "GraphQLResolver").Logger()
	gqlResolver := resolvers.NewResolver(userService, authService, resolverLogger)

	gqlConfig := generated.Config{Resolvers: gqlResolver}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))

	// 4f. Inisialisasi Middleware
	authMw := middleware.NewAuthMiddleware(tokenMaker)

	// --- 5. Setup HTTP Server & Routing ---
	httpMux := http.NewServeMux()

	// Rute untuk GraphQL Playground
	httpMux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// Rute untuk API query, dilindungi oleh middleware
	httpMux.Handle("/query", authMw.Handler(gqlServer))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: httpMux,
	}

	// --- 6. Start Server (Goroutine) ---
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msgf("Starting HTTP server at http://localhost:%s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Could not start server")
		}
	}()

	// --- 7. Wait for Shutdown Signal ---
	<-ctx.Done() // Blokir hingga sinyal shutdown diterima

	// --- 8. Graceful Shutdown ---
	log.Info().Msg("Shutting down application...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Warn().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Application shut down gracefully.")
}
