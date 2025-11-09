package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// --- gqlgen Handlers ---
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	// --- Internal Dependencies ---
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driven/postgres"
	"github.com/ipincamp/srikandi-sehat/internal/core/service"
	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/logger"
	"github.com/ipincamp/srikandi-sehat/pkg/password"
	"github.com/ipincamp/srikandi-sehat/pkg/token"

	// --- GraphQL Package Imports ---
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/generated"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/resolvers"
)

func main() {
	// --- 1. Load Configuration ---
	cfg, err := config.Load()
	if err != nil {
		tempLogger := logger.NewLogger("development") // Use "development" for a clean exit message
		tempLogger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// --- 1.1. Initialize Logger ---
	log := logger.NewLogger(cfg.Server.Env)
	log.Info().Str("Env", cfg.Server.Env).Msg("Configuration loaded")

	// --- 2. Setup Application Context ---
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
	dbPool, err := postgres.Connect(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbPool.Close()
	log.Info().Msg("Database connection pool established")

	// --- 4. Dependency Injection (Composition Root) ---
	// This is the only place in the app that knows about concrete implementations.
	// We assemble the application components here, following Hexagonal Architecture.

	// 4a. Initialize 'pkg' helpers (implementations)
	hasher := password.NewArgon2idHasher()
	tokenMaker, err := token.NewPasetoMaker(cfg.Token.SymmetricKey, cfg.Token.Issuer)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Paseto token maker")
	}

	// 4b. Initialize Driven Adapters (Repositories)
	// We create the concrete 'postgres.userRepository' implementation.
	userRepo := postgres.NewUserRepository(dbPool)

	// 4c. Initialize Core Services
	// We inject the repository (an interface) into the service.
	authService := service.NewAuthService(userRepo, tokenMaker, hasher, cfg.Token)

	// 4d. Initialize Driving Adapters (GraphQL)
	// We inject the service (an interface) into the resolver.
	gqlResolver := resolvers.NewResolver(authService)

	// 4e. Create GraphQL server configuration
	gqlConfig := generated.Config{Resolvers: gqlResolver}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))

	// --- 5. Start Application (HTTP Server) ---
	log.Info().Msg("Application dependencies initialized.")

	httpMux := http.NewServeMux()
	httpMux.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	httpMux.Handle("/query", gqlServer)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: httpMux,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Starting HTTP server... (GraphQL Playground at http://localhost:" + cfg.Server.Port + ")")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Could not start server")
		}
	}()

	// --- 6. Wait for Shutdown Signal ---
	<-ctx.Done()

	// --- 7. Graceful Shutdown ---
	log.Info().Msg("Shutting down application...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Warn().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Application shut down gracefully.")
}
