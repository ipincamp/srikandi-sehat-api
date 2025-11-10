# Developer Guide

Welcome to the Srikandi Sehat GraphQL API development guide! This document will help you understand the codebase and contribute effectively.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Workflow](#development-workflow)
3. [Code Organization](#code-organization)
4. [Adding New Features](#adding-new-features)
5. [Testing](#testing)
6. [Database Migrations](#database-migrations)
7. [Code Style](#code-style)
8. [Common Tasks](#common-tasks)
9. [Troubleshooting](#troubleshooting)

## Getting Started

### Prerequisites

Ensure you have the following installed:
- Go 1.25.4 or higher
- Docker & Docker Compose
- PostgreSQL 16+ (or use Docker)
- Make
- Git

### Initial Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ipincamp/srikandi-sehat-graphql.git
   cd srikandi-sehat-graphql
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Start infrastructure:**
   ```bash
   docker-compose up -d
   ```

5. **Run migrations:**
   ```bash
   make migrate
   ```

6. **Start development server:**
   ```bash
   make dev
   ```

Visit `http://localhost:8000` to access the GraphQL Playground.

## Development Workflow

### Daily Development

1. **Pull latest changes:**
   ```bash
   git pull origin main
   go mod download
   make migrate
   ```

2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make changes and test:**
   ```bash
   make dev  # Start dev server with hot reload
   ```

4. **Run tests:**
   ```bash
   go test ./...
   ```

5. **Commit and push:**
   ```bash
   git add .
   git commit -m "feat: add your feature"
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding/updating tests
- `chore:` Maintenance tasks

**Examples:**
```
feat: add email verification feature
fix: resolve race condition in user registration
docs: update API documentation for login endpoint
refactor: improve error handling in auth service
test: add unit tests for password hashing
chore: update dependencies
```

## Code Organization

### Hexagonal Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    External World                           │
│                  (HTTP, Database, Email)                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
            ┌──────────┴──────────┐
            │                     │
         Adapters             Adapters
         (Driving)            (Driven)
            │                     │
            ↓                     ↓
      ┌─────────┐           ┌─────────┐
      │  Ports  │←─────────→│  Ports  │
      │(Service)│           │ (Repo)  │
      └────┬────┘           └────┬────┘
           │                     │
           └──────────┬──────────┘
                      ↓
                 ┌─────────┐
                 │ Domain  │
                 │(Entities)│
                 └─────────┘
```

### Directory Structure

- **`cmd/`** - Application entry points
  - `server/` - Main API server
  - `migrate/` - Database migration CLI

- **`internal/`** - Private application code
  - `core/` - Business logic (domain, ports, services)
  - `adapters/` - External integrations
    - `driving/` - Inbound (GraphQL, REST)
    - `driven/` - Outbound (Database, Email)

- **`pkg/`** - Reusable packages
  - `config/` - Configuration management
  - `logger/` - Logging utilities
  - `password/` - Password hashing
  - `token/` - Token generation/validation

## Adding New Features

### Step-by-Step Guide

Let's add a "Profile Picture Upload" feature as an example:

#### 1. Define Domain Entity

**File:** `internal/core/domain/user.go`

```go
type User struct {
    ID             uint
    UUID           string
    Name           string
    Email          string
    Password       string
    ProfilePicture string // Add new field
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

#### 2. Update Port Interface

**File:** `internal/core/ports/services.go`

```go
type UserService interface {
    // Existing methods...
    
    // New method
    UpdateProfilePicture(ctx context.Context, userID, pictureURL string) error
}
```

#### 3. Implement Service

**File:** `internal/core/service/user_service.go`

```go
func (s *userService) UpdateProfilePicture(ctx context.Context, userID, pictureURL string) error {
    // 1. Validate input
    if pictureURL == "" {
        return ports.ErrInvalidInput
    }
    
    // 2. Start transaction
    tx, err := s.uow.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)
    
    // 3. Get user
    txUserRepo := tx.GetUserRepository()
    user, err := txUserRepo.FindByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // 4. Update field
    user.ProfilePicture = pictureURL
    
    // 5. Save
    if err := txUserRepo.Update(ctx, user); err != nil {
        return err
    }
    
    // 6. Commit
    return tx.Commit(ctx)
}
```

#### 4. Update Database Model

**File:** `internal/adapters/driven/postgres/models.go`

```go
type UserModel struct {
    ID             uint      `gorm:"primaryKey"`
    UUID           string    `gorm:"unique;not null"`
    Name           string    `gorm:"not null"`
    Email          string    `gorm:"unique;not null"`
    Password       string    `gorm:"not null"`
    ProfilePicture string    `gorm:"type:text"` // Add new field
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

func (m *UserModel) ToDomain() *domain.User {
    return &domain.User{
        // ... existing fields
        ProfilePicture: m.ProfilePicture, // Add mapping
    }
}

func FromDomain(u *domain.User) *UserModel {
    return &UserModel{
        // ... existing fields
        ProfilePicture: u.ProfilePicture, // Add mapping
    }
}
```

#### 5. Create Database Migration

```bash
make create-migration name=add_profile_picture_to_users
```

**File:** `internal/adapters/driven/postgres/db/migrations/YYYYMMDDHHMMSS_add_profile_picture_to_users.go`

```go
package migrations

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/gorm"
)

func AddProfilePictureToUsers() *gormigrate.Migration {
    type User struct {
        ProfilePicture string `gorm:"type:text"`
    }
    
    return &gormigrate.Migration{
        ID: "YYYYMMDDHHMMSS",
        Migrate: func(tx *gorm.DB) error {
            return tx.AutoMigrate(&User{})
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropColumn("users", "profile_picture")
        },
    }
}
```

Register migration in `registry.go`:
```go
func GetMigrations() []*gormigrate.Migration {
    return []*gormigrate.Migration{
        // ... existing migrations
        AddProfilePictureToUsers(),
    }
}
```

#### 6. Update GraphQL Schema

**File:** `internal/adapters/driving/graphql/schema/user.graphqls`

```graphql
type User {
  uuid: String!
  name: String!
  email: String!
  profilePicture: String  # Add new field
  createdAt: String!
  updatedAt: String!
}

input UpdateProfilePictureInput {
  pictureUrl: String!
}

extend type Mutation {
  updateProfilePicture(input: UpdateProfilePictureInput!): User!
}
```

#### 7. Generate GraphQL Code

```bash
make generate
```

#### 8. Implement Resolver

**File:** `internal/adapters/driving/graphql/resolvers/user.resolver.go`

```go
func (r *mutationResolver) UpdateProfilePicture(
    ctx context.Context,
    input model.UpdateProfilePictureInput,
) (*model.User, error) {
    // Get authenticated user
    userID, err := GetAuthUserUUID(ctx)
    if err != nil {
        return nil, err
    }
    
    // Call service
    if err := r.userService.UpdateProfilePicture(ctx, userID, input.PictureURL); err != nil {
        return nil, err
    }
    
    // Return updated user
    user, err := r.userService.GetUserByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    return toGraphQLUser(user), nil
}
```

#### 9. Write Tests

**File:** `internal/core/service/user_service_test.go`

```go
func TestUpdateProfilePicture(t *testing.T) {
    // Setup mocks
    mockUserRepo := &mockUserRepository{}
    mockUOW := &mockUnitOfWork{}
    
    svc := NewUserService(mockUserRepo, nil, logger, mockUOW)
    
    // Test success case
    err := svc.UpdateProfilePicture(ctx, "user-uuid", "https://example.com/pic.jpg")
    assert.NoError(t, err)
    
    // Test validation
    err = svc.UpdateProfilePicture(ctx, "user-uuid", "")
    assert.Error(t, err)
}
```

#### 10. Test Manually

Run migrations:
```bash
make migrate
```

Start server:
```bash
make dev
```

Test in GraphQL Playground:
```graphql
mutation {
  updateProfilePicture(input: {
    pictureUrl: "https://example.com/profile.jpg"
  }) {
    uuid
    name
    profilePicture
  }
}
```

## Testing

### Unit Tests

Test business logic in isolation:

```go
// internal/core/service/auth_service_test.go
func TestRegister_Success(t *testing.T) {
    // Arrange
    mockUserRepo := &mockUserRepository{
        SaveFunc: func(ctx context.Context, user *domain.User) error {
            return nil
        },
    }
    
    // Act
    result, err := authService.Register(ctx, "John", "john@example.com", "password")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.NotEmpty(t, result.AccessToken)
}
```

### Integration Tests

Test with real database:

```go
func TestUserRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := postgres.NewUserRepository(db, logger)
    
    // Test
    user := &domain.User{
        UUID:  uuid.NewString(),
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    err := repo.Save(context.Background(), user)
    assert.NoError(t, err)
    
    found, err := repo.FindByEmail(context.Background(), "test@example.com")
    assert.NoError(t, err)
    assert.Equal(t, user.UUID, found.UUID)
}
```

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/core/service/...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Run specific test
go test -run TestRegister_Success ./internal/core/service/...
```

## Database Migrations

### Creating Migrations

```bash
make create-migration name=create_table_name
```

### Migration Template

```go
package migrations

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/gorm"
)

func CreateTableName() *gormigrate.Migration {
    type YourModel struct {
        ID        uint      `gorm:"primaryKey"`
        Name      string    `gorm:"not null"`
        CreatedAt time.Time
    }
    
    return &gormigrate.Migration{
        ID: "YYYYMMDDHHMMSS",
        Migrate: func(tx *gorm.DB) error {
            return tx.AutoMigrate(&YourModel{})
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable("your_models")
        },
    }
}
```

### Running Migrations

```bash
# Apply pending migrations
make migrate

# Rollback last migration
make migrate-down

# Reset database (development only!)
make migrate-fresh

# Drop all tables (danger!)
make migrate-prune
```

## Code Style

### Go Code

Follow [Effective Go](https://golang.org/doc/effective_go) and:

- Use `gofmt` for formatting
- Use meaningful variable names
- Keep functions small and focused
- Document exported functions
- Handle errors explicitly

**Good:**
```go
// GetUserByID retrieves a user by their UUID.
// Returns ErrUserNotFound if the user doesn't exist.
func (s *userService) GetUserByID(ctx context.Context, uuid string) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, uuid)
    if err != nil {
        if errors.Is(err, ports.ErrUserNotFound) {
            s.logger.Warn().Str("uuid", uuid).Msg("User not found")
            return nil, ports.ErrUserNotFound
        }
        return nil, err
    }
    return user, nil
}
```

**Bad:**
```go
func (s *userService) Get(c context.Context, id string) (*domain.User, error) {
    u, e := s.userRepo.FindByID(c, id)
    if e != nil {
        return nil, e
    }
    return u, nil
}
```

### GraphQL Schema

- Use PascalCase for types
- Use camelCase for fields
- Add descriptions for complex types
- Use input types for mutations

```graphql
"""
User represents a registered user in the system.
"""
type User {
  "Unique identifier for the user"
  uuid: String!
  "Full name of the user"
  name: String!
  "Email address (unique)"
  email: String!
}
```

## Common Tasks

### Adding a New Endpoint

1. Update GraphQL schema (`schema/*.graphqls`)
2. Run `make generate`
3. Implement resolver
4. Add service method if needed
5. Write tests
6. Update documentation

### Adding a New Repository Method

1. Add method to port interface (`ports/repositories.go`)
2. Implement in adapter (`adapters/driven/postgres/*_repository.go`)
3. Write tests
4. Use in service layer

### Changing Database Schema

1. Create migration (`make create-migration`)
2. Update domain entity
3. Update database model
4. Update conversion functions
5. Run migration (`make migrate`)
6. Test thoroughly

### Adding Configuration

1. Add to `pkg/config/config.go`
2. Add to `.env.example`
3. Document in README
4. Use in code

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8000
lsof -i :8000

# Kill process
kill -9 <PID>
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps

# View logs
docker logs srikandi_sehat_postgres_db

# Restart
docker-compose restart db
```

### Migration Errors

```bash
# Check current migration status
go run cmd/migrate/main.go

# Force rollback
make migrate-down

# Reset database (development only!)
make migrate-fresh
```

### Module Issues

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy up
go mod tidy
```

### GraphQL Generation Issues

```bash
# Clean generated files
rm -rf internal/adapters/driving/graphql/generated/*

# Regenerate
make generate
```

## Best Practices

### 1. Always Use Context

```go
func (s *service) DoSomething(ctx context.Context, ...) error {
    // Pass context to all downstream calls
    result, err := s.repo.Find(ctx, ...)
}
```

### 2. Handle Errors Properly

```go
if err != nil {
    // Log with context
    s.logger.Error().
        Err(err).
        Str("user_id", userID).
        Msg("Failed to process request")
    
    // Return appropriate error
    return ports.ErrInternalServer
}
```

### 3. Use Transactions for Multiple Operations

```go
tx, err := s.uow.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Multiple operations...

return tx.Commit(ctx)
```

### 4. Validate Input

```go
if email == "" || !isValidEmail(email) {
    return ports.ErrInvalidInput
}
```

### 5. Keep Layers Separate

- Don't import adapters in core
- Don't put business logic in resolvers
- Convert between types at boundaries

## Resources

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Effective Go](https://golang.org/doc/effective_go)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [GORM Documentation](https://gorm.io/docs/)

## Getting Help

- Check existing documentation in `/docs`
- Search for similar code in the project
- Ask in team chat
- Create an issue on GitHub

Happy coding! 🚀
