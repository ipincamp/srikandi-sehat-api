# Srikandi Sehat GraphQL API

A production-ready GraphQL API built with Go, following **Hexagonal Architecture** (Ports and Adapters) principles. This project provides a secure, scalable authentication and user management system.

## 🏗️ Architecture

This project follows **Hexagonal Architecture** (also known as Ports and Adapters), which promotes:

- **Clean separation of concerns** between business logic and infrastructure
- **Technology independence** - easily swap databases, frameworks, or external services
- **Testability** - core business logic can be tested without external dependencies
- **Maintainability** - clear boundaries between different layers

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                     Driving Adapters                        │
│              (GraphQL Resolvers, Middleware)                │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    Driving Ports                            │
│              (Service Interfaces)                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    Core Domain                              │
│        (Business Logic & Entities)                          │
│    User, OTP, AuthService, UserService                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                    Driven Ports                             │
│         (Repository Interfaces)                             │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                   Driven Adapters                           │
│         (PostgreSQL, Mailgun, etc.)                         │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Features

### ✅ Implemented Features

- **User Registration** - Create new user accounts with email and password
- **User Login** - Authenticate users with email and password
- **Token-based Authentication** - Secure PASETO tokens (access + refresh)
- **Token Refresh** - Renew expired access tokens without re-login
- **Change Password** - Authenticated users can update their password
- **Update Profile** - Users can update their profile information
- **Get User Profile** - Retrieve authenticated user information
- **Password Reset** - Forgot password flow with OTP via email
- **Email Verification** - Verify email addresses with OTP codes
- **Request Email Change** - Change email address with verification
- **Account Deletion** - Soft delete user accounts
- **Database Migrations** - Version-controlled database schema
- **Health Check** - API health status endpoint
- **DataLoader** - Efficient batch loading to prevent N+1 queries
- **Structured Logging** - JSON-structured logs with zerolog
- **Environment Configuration** - Flexible configuration via environment variables

### 🔮 Planned Features

Additional features can be easily added following the established patterns.

## 📋 Prerequisites

- **Go 1.25.4+**
- **PostgreSQL 16+**
- **Docker & Docker Compose** (for local development)
- **Make** (for running Makefile commands)

## 🛠️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/ipincamp/srikandi-sehat-graphql.git
cd srikandi-sehat-graphql
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Application
APP_ENV=development
APP_PORT=8000
APP_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=srikandi_sehat
DB_USER=postgres
DB_PASS=your_password
DB_SSL_MODE=disable

# PASETO Token (must be 32 bytes)
TOKEN_SYMMETRIC_KEY=your_32_byte_secret_key_here_now
TOKEN_ISSUER=srikandi-sehat-api
TOKEN_ACCESS_DURATION=15m
TOKEN_REFRESH_DURATION=720h

# Email Service
MAIL_DRIVER=smtp
MAIL_HOST=localhost
MAIL_PORT=1025
MAIL_FROM_ADDRESS=no-reply@srikandi-sehat.com
MAIL_FROM_NAME=Srikandi Sehat
```

### 4. Start Infrastructure Services

Start PostgreSQL and Mailpit (local email testing):

```bash
docker-compose up -d
```

### 5. Run Database Migrations

```bash
make migrate
```

### 6. Start the Development Server

```bash
make dev
```

The GraphQL API will be available at `http://localhost:8000/query` with the GraphQL Playground at `http://localhost:8000/`.

## 📚 Documentation

- **[Architecture Guide](./docs/ARCHITECTURE.md)** - Detailed explanation of hexagonal architecture implementation
- **[API Reference](./docs/API.md)** - Complete GraphQL API documentation
- **[Developer Guide](./docs/DEVELOPER.md)** - Guide for developers contributing to this project
- **[Feature Documentation](./docs/FEATURES.md)** - Detailed documentation of all features

## 🔧 Available Commands

```bash
# Development
make dev              # Run with hot reload
make build            # Build binary
make run              # Build and run production binary
make clean            # Clean build artifacts

# Database Migrations
make create-migration name=create_something_table  # Create new migration
make migrate          # Run all pending migrations
make migrate-down     # Rollback last migration
make migrate-fresh    # Drop all tables and re-run migrations
make migrate-prune    # Drop all tables (dangerous!)

# GraphQL
make generate         # Generate GraphQL code from schema

# Help
make help             # Show all available commands
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/core/service/...
```

## 🏛️ Project Structure

```
.
├── cmd/
│   ├── migrate/         # Database migration CLI
│   └── server/          # Main application entry point
├── internal/
│   ├── adapters/
│   │   ├── driven/      # Infrastructure adapters (DB, Email, etc.)
│   │   │   ├── mailgun/
│   │   │   └── postgres/
│   │   └── driving/     # API adapters (GraphQL, REST, etc.)
│   │       └── graphql/
│   └── core/
│       ├── domain/      # Domain entities
│       ├── ports/       # Port interfaces
│       └── service/     # Business logic services
├── pkg/
│   ├── config/          # Configuration management
│   ├── logger/          # Logging utilities
│   ├── password/        # Password hashing
│   └── token/           # Token generation/validation
├── tools/               # Development tools
├── docker-compose.yml   # Local infrastructure
├── Makefile            # Build automation
└── README.md           # This file
```

## 🔐 Security

- **Password Hashing**: Argon2id algorithm (memory-hard, secure)
- **Token Security**: PASETO v2 (authenticated encryption)
- **SQL Injection**: Protected via GORM parameterization
- **Rate Limiting**: Recommended to add at reverse proxy level
- **HTTPS**: Recommended for production (use reverse proxy)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Contributing

Please read [DEVELOPER.md](./docs/DEVELOPER.md) for details on our code of conduct and the process for submitting pull requests.

## 📧 Contact

For questions or support, please contact the development team.

---

Built with ❤️ using Go and GraphQL
