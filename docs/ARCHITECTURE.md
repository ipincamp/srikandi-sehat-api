# Architecture Documentation

## Hexagonal Architecture Overview

This project implements **Hexagonal Architecture** (also known as **Ports and Adapters**), which organizes code to ensure the business logic (core domain) remains independent of external concerns like databases, web frameworks, or third-party services.

## Core Principles

### 1. Dependency Rule

Dependencies flow **inward only**. The core domain has zero dependencies on outer layers:

```
External → Adapters → Ports → Domain
```

- **Domain**: Pure business logic (no external dependencies)
- **Ports**: Interfaces defining contracts
- **Adapters**: Implementations of ports (database, HTTP, etc.)

### 2. Separation of Concerns

Each layer has a specific responsibility:

- **Core Domain**: Business entities and rules
- **Application Services**: Use cases and workflows
- **Ports**: Abstract interfaces
- **Adapters**: Concrete implementations

## Layer Details

### Domain Layer (`internal/core/domain/`)

Contains pure business entities with no external dependencies.

**Files:**
- `user.go` - User entity
- `otp.go` - OTP (One-Time Password) entity with type constants

**Example:**
```go
type User struct {
    ID        uint
    UUID      string
    Name      string
    Email     string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Key Points:**
- No GORM tags or database annotations
- No JSON tags (those belong in adapters)
- Pure Go structs representing business concepts
- Business logic methods can be added here

### Ports Layer (`internal/core/ports/`)

Defines interfaces that connect the core to the outside world.

**Types of Ports:**

#### Driving Ports (Services)

Interfaces that **drive** the application logic (used by external adapters):

**`services.go`:**
```go
type AuthService interface {
    Register(ctx, name, email, password) (*AuthResponse, error)
    Login(ctx, email, password) (*AuthResponse, error)
    RefreshToken(ctx, refreshToken) (*AuthResponse, error)
    // ... more methods
}

type UserService interface {
    GetUserByID(ctx, uuid) (*User, error)
    CreateUser(ctx, name, email, password) (*User, error)
    UpdateProfile(ctx, userID, newName) (*User, error)
    // ... more methods
}
```

#### Driven Ports (Repositories)

Interfaces that the application **drives** (implemented by infrastructure):

**`repositories.go`:**
```go
type UserRepository interface {
    Save(ctx, user) error
    FindByID(ctx, uuid) (*User, error)
    FindByEmail(ctx, email) (*User, error)
    FindMapByUUIDs(ctx, uuids) (map[string]*User, error)
    Update(ctx, user) error
    Delete(ctx, uuid) error
}

type OTPRepository interface {
    Save(ctx, otp) error
    FindAndConsume(ctx, code, type) (*OTP, error)
}
```

**`unit_of_work.go`:**
```go
type UnitOfWork interface {
    Begin(ctx) (Transaction, error)
}

type Transaction interface {
    GetUserRepository() UserRepository
    GetOTPRepository() OTPRepository
    Commit(ctx) error
    Rollback(ctx) error
}
```

#### Other Ports

**`errors.go`:**
Defines domain-specific errors:
```go
var (
    ErrUserNotFound        = errors.New("user not found")
    ErrEmailExists         = errors.New("email already exists")
    ErrInvalidCredentials  = errors.New("invalid credentials")
    ErrTokenExpired        = errors.New("token expired")
    // ... more errors
)
```

### Service Layer (`internal/core/service/`)

Implements the driving ports (business logic).

**Files:**
- `auth_service.go` - Authentication business logic
- `user_service.go` - User management business logic
- `auth_service_test.go` - Unit tests for auth service

**Key Characteristics:**
- Depends only on domain entities and port interfaces
- Orchestrates use cases
- Transaction management via Unit of Work pattern
- No knowledge of database or web framework

**Example Flow (Register):**
```go
func (s *authService) Register(ctx, name, email, password) (*AuthResponse, error) {
    // 1. Hash password (using injected hasher)
    hashedPassword := s.hasher.Hash(password)
    
    // 2. Create domain entity
    user := &domain.User{
        UUID:     uuid.New(),
        Name:     name,
        Email:    email,
        Password: hashedPassword,
    }
    
    // 3. Begin transaction
    tx := s.uow.Begin(ctx)
    defer handleRollback(tx, &err)
    
    // 4. Save using transactional repository
    txUserRepo := tx.GetUserRepository()
    if err := txUserRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Commit transaction
    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }
    
    // 6. Generate tokens
    return s.createTokenSet(user)
}
```

### Driven Adapters (`internal/adapters/driven/`)

Implement the driven ports (infrastructure).

#### PostgreSQL Adapter (`driven/postgres/`)

**Files:**
- `postgres.go` - Connection pool management
- `gorm.go` - GORM instance creation
- `models.go` - Database models (with GORM tags)
- `user_repository.go` - UserRepository implementation
- `otp_repository.go` - OTPRepository implementation
- `unit_of_work.go` - Transaction management
- `db/migrations/` - Database migrations

**Key Pattern: Separation of Domain and DB Models**

```go
// models.go - Database model
type UserModel struct {
    ID        uint      `gorm:"primaryKey"`
    UUID      string    `gorm:"unique;not null"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (UserModel) TableName() string {
    return "users"
}

// Conversion functions
func (m *UserModel) ToDomain() *domain.User {
    return &domain.User{
        ID:        m.ID,
        UUID:      m.UUID,
        Name:      m.Name,
        Email:     m.Email,
        Password:  m.Password,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

func FromDomain(u *domain.User) *UserModel {
    return &UserModel{
        ID:        u.ID,
        UUID:      u.UUID,
        Name:      u.Name,
        Email:     u.Email,
        Password:  u.Password,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
```

**Unit of Work Implementation:**

The Unit of Work pattern ensures that multiple repository operations can be executed within a single database transaction:

```go
type unitOfWork struct {
    pool   *pgxpool.Pool
    logger zerolog.Logger
}

func (uow *unitOfWork) Begin(ctx) (*transaction, error) {
    tx, err := uow.pool.Begin(ctx)
    return &transaction{tx: tx, logger: uow.logger}, err
}

type transaction struct {
    tx     pgx.Tx
    logger zerolog.Logger
}

func (t *transaction) GetUserRepository() ports.UserRepository {
    db := gorm.New(t.tx, t.logger)
    return NewUserRepository(db, t.logger)
}

func (t *transaction) Commit(ctx) error {
    return t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx) error {
    return t.tx.Rollback(ctx)
}
```

#### Mailgun Adapter (`driven/mailgun/`)

**Files:**
- `mail_service.go` - Email service implementation

Supports two drivers:
- **SMTP** - For local development (Mailpit)
- **Mailgun** - For production email delivery

### Driving Adapters (`internal/adapters/driving/`)

Implement the driving ports (API layer).

#### GraphQL Adapter (`driving/graphql/`)

**Structure:**
```
graphql/
├── schema/              # GraphQL schema definitions
│   ├── auth.graphqls
│   ├── user.graphqls
│   └── health.graphqls
├── resolvers/           # Resolver implementations
│   ├── resolver.go
│   ├── auth.resolver.go
│   ├── user.resolver.go
│   └── health.resolver.go
├── models/              # GraphQL models (generated)
├── generated/           # Generated GraphQL code
├── middleware/          # HTTP middleware
│   ├── auth.go
│   └── dataloader.go
└── dataloader/          # DataLoader for N+1 prevention
    └── loader.go
```

**Resolver Example:**

```go
func (r *mutationResolver) Register(ctx, input) (*model.AuthResponse, error) {
    // Call service layer
    result, err := r.authService.Register(
        ctx,
        input.Name,
        input.Email,
        input.Password,
    )
    if err != nil {
        // Map domain errors to GraphQL errors
        return nil, mapError(err)
    }
    
    // Convert to GraphQL model
    return &model.AuthResponse{
        AccessToken:  result.AccessToken,
        RefreshToken: result.RefreshToken,
    }, nil
}
```

**Middleware: Authentication**

Implements "fail-open" for public requests, "fail-close" for authenticated:

```go
func (am *AuthMiddleware) Handler(next) http.Handler {
    return http.HandlerFunc(func(w, r) {
        authHeader := r.Header.Get("Authorization")
        
        // No token? Allow public access
        if authHeader == "" {
            next.ServeHTTP(w, r)
            return
        }
        
        // Token present? Must be valid
        token := extractToken(authHeader)
        payload, err := am.tokenMaker.ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", 401)
            return
        }
        
        // Inject user ID into context
        ctx := context.WithValue(r.Context(), AuthUserUUIDKey, payload.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**DataLoader: N+1 Query Prevention**

Batches multiple database queries into single queries:

```go
type Loaders struct {
    UserByUUID *dataloader.Loader
}

func NewLoaders(userRepo) *Loaders {
    return &Loaders{
        UserByUUID: dataloader.NewBatchedLoader(
            func(ctx, keys) []*dataloader.Result {
                // Batch load all users at once
                users := userRepo.FindMapByUUIDs(ctx, keys)
                
                // Return in same order as keys
                results := make([]*dataloader.Result, len(keys))
                for i, key := range keys {
                    results[i] = &dataloader.Result{
                        Data: users[key],
                    }
                }
                return results
            },
        ),
    }
}
```

## Package Layer (`pkg/`)

Reusable utilities that can be used across projects.

### Config (`pkg/config/`)

Environment-based configuration loading:
- Reads from `.env` file (development)
- Falls back to environment variables (production)
- Type-safe configuration structs

### Logger (`pkg/logger/`)

Structured logging with zerolog:
- JSON output for production
- Pretty console output for development
- Context-aware logging

### Password (`pkg/password/`)

Secure password hashing:
- **Argon2id** algorithm
- Memory-hard and resistant to GPU attacks
- Configurable parameters

### Token (`pkg/token/`)

PASETO token implementation:
- **PASETO v2** (Platform-Agnostic SEcurity TOkens)
- Authenticated encryption (no signature verification needed)
- Separate access and refresh tokens

## Dependency Injection

Dependencies are wired together in `cmd/server/main.go`:

```go
func main() {
    // 1. Load configuration
    cfg := config.Load()
    log := logger.NewLogger(cfg.Server.Env)
    
    // 2. Connect to database
    dbPool := postgres.Connect(ctx, cfg.Database.DSN())
    
    // 3. Initialize utilities
    hasher := password.NewArgon2idHasher()
    tokenMaker := token.NewPasetoMaker(cfg.Token.SymmetricKey)
    
    // 4. Initialize repositories
    userRepo := postgres.NewUserRepository(dbPool, log)
    otpRepo := postgres.NewOTPRepository(dbPool, log)
    mailSvc := mailgun.NewMailService(cfg.Mail, log)
    
    // 5. Initialize Unit of Work
    uow := postgres.NewUnitOfWork(dbPool, log)
    
    // 6. Initialize services (inject dependencies)
    authService := service.NewAuthService(
        userRepo, tokenMaker, hasher, cfg.Token, log, uow, otpRepo, mailSvc,
    )
    userService := service.NewUserService(
        userRepo, hasher, log, uow,
    )
    
    // 7. Initialize resolvers
    resolver := resolvers.NewResolver(authService, userService)
    
    // 8. Create GraphQL handler
    srv := handler.NewDefaultServer(
        generated.NewExecutableSchema(
            generated.Config{Resolvers: resolver},
        ),
    )
    
    // 9. Setup HTTP server with middleware
    mux := http.NewServeMux()
    mux.Handle("/query", authMiddleware.Handler(srv))
    
    // 10. Start server
    http.ListenAndServe(":8000", mux)
}
```

## Benefits of This Architecture

### 1. **Testability**

Core business logic can be tested without external dependencies:

```go
func TestRegister(t *testing.T) {
    // Create mock repositories
    mockUserRepo := &mockUserRepository{}
    mockUOW := &mockUnitOfWork{}
    
    // Create service with mocks
    svc := NewAuthService(mockUserRepo, ...)
    
    // Test business logic
    result, err := svc.Register(ctx, "John", "john@example.com", "password")
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 2. **Flexibility**

Easy to swap implementations:
- Switch from PostgreSQL to MySQL
- Replace PASETO with JWT
- Change from GraphQL to REST

### 3. **Maintainability**

Clear boundaries between layers:
- Database changes don't affect business logic
- API changes don't affect domain entities
- Easy to locate and fix bugs

### 4. **Scalability**

Each layer can evolve independently:
- Add new features without touching existing code
- Optimize specific layers without affecting others
- Multiple teams can work on different layers

### 5. **Domain-Driven Design**

Focus on business logic:
- Domain entities reflect business concepts
- Services implement use cases
- Technical details are isolated

## Best Practices

### 1. **Always Use Interfaces**

```go
// Good: Depend on interface
type AuthService struct {
    userRepo ports.UserRepository
}

// Bad: Depend on concrete type
type AuthService struct {
    userRepo *postgres.UserRepository
}
```

### 2. **Keep Domain Pure**

```go
// Good: Pure domain entity
type User struct {
    ID    uint
    Name  string
    Email string
}

// Bad: Domain with framework dependencies
type User struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `gorm:"not null" json:"name"`
    Email string `gorm:"unique" json:"email"`
}
```

### 3. **Use Unit of Work for Transactions**

```go
// Good: Transactional operation
tx := uow.Begin(ctx)
defer handleRollback(tx, &err)
txRepo := tx.GetUserRepository()
txRepo.Save(ctx, user)
tx.Commit(ctx)

// Bad: Non-transactional when multiple operations needed
userRepo.Save(ctx, user)
otpRepo.Save(ctx, otp) // What if this fails?
```

### 4. **Map Errors at Boundaries**

```go
// In resolver (adapter layer)
result, err := r.authService.Login(ctx, email, password)
if err != nil {
    if errors.Is(err, ports.ErrInvalidCredentials) {
        return nil, gqlerror.Errorf("Invalid email or password")
    }
    return nil, gqlerror.Errorf("Internal server error")
}
```

## Conclusion

This hexagonal architecture provides a robust, maintainable, and testable foundation for the Srikandi Sehat GraphQL API. By keeping the core domain independent and using clear interfaces, the codebase remains flexible and easy to evolve as requirements change.
