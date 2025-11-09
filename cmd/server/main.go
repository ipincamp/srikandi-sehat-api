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
	"github.com/ipincamp/srikandi-sehat/internal/core/service"

	"github.com/ipincamp/srikandi-sehat/pkg/config"
	"github.com/ipincamp/srikandi-sehat/pkg/logger"
	"github.com/ipincamp/srikandi-sehat/pkg/password"
	"github.com/ipincamp/srikandi-sehat/pkg/token"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/generated"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/middleware"
	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/resolvers"
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

	// 4a. Initialize 'pkg' helpers (implementations)
	hasher := password.NewArgon2idHasher()
	tokenMaker, err := token.NewPasetoMaker(cfg.Token.SymmetricKey, cfg.Token.Issuer)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Paseto token maker")
	}

	// 4b. Initialize Driven Adapters (Repositories)
	// Inject the logger with component context
	userRepoLogger := log.With().Str("component", "UserRepository(Non-TX)").Logger()
	// Create the *non-transactional* user repository.
	// We pass dbPool, which satisfies the dbExecutor interface.
	// This repo is used for read-only operations like Login.
	userRepo := postgres.NewUserRepository(dbPool, userRepoLogger)

	// 4c. Initialize Unit of Work
	uowLogger := log.With().Str("component", "UnitOfWork").Logger()
	// Create the Unit of Work factory, passing the pool
	uow := postgres.NewUnitOfWork(dbPool, uowLogger)

	// 4d. Initialize Core Services
	authServiceLogger := log.With().Str("component", "AuthService").Logger()
	userService := service.NewUserService(
		userRepo,          // ports.UserRepository
		hasher,            // password.Hasher
		authServiceLogger, // zerolog.Logger
		uow,               // ports.UnitOfWork
	)
	authService := service.NewAuthService(
		userRepo, // Pass the non-tx repo for reads
		tokenMaker,
		hasher,
		cfg.Token,
		authServiceLogger,
		uow,
	)

	// 4e. Initialize Driving Adapters (GraphQL)
	// Inject the service and a logger
	resolverLogger := log.With().Str("component", "GraphQLResolver").Logger()
	gqlResolver := resolvers.NewResolver(
		authService,
		userService,
		resolverLogger,
	)
	gqlConfig := generated.Config{Resolvers: gqlResolver}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))

	// --- 5. Start Application (HTTP Server) ---
	log.Info().Msg("Application dependencies initialized.")

	// Create Auth Middleware
	authMw := middleware.NewAuthMiddleware(tokenMaker)

	// Create Dataloader Middleware
	// We pass the non-transactional repo, which is perfect for read-only batching.
	dataloaderMw := middleware.NewDataloaderMiddleware(userRepo)

	httpMux := http.NewServeMux()
	httpMux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// --- IMPORTANT: Chain the middleware ---
	// The request will flow: dataloaderMw -> authMw -> gqlServer
	// This ensures the loader is in the context *before* the auth (or any other) middleware runs.
	httpMux.Handle("/query", dataloaderMw.Handler(authMw.Handler(gqlServer)))

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
