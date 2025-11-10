# Quick Reference Guide

A quick reference for common tasks and commands in the Srikandi Sehat GraphQL API project.

## Quick Start

```bash
# Clone and setup
git clone https://github.com/ipincamp/srikandi-sehat-graphql.git
cd srikandi-sehat-graphql
cp .env.example .env
# Edit .env with your settings

# Start infrastructure
docker-compose up -d

# Run migrations
make migrate

# Start development server
make dev
```

Access GraphQL Playground at `http://localhost:8000/`

---

## Common Commands

### Development

```bash
make dev              # Start with hot reload (Air)
make build            # Build binary
make run              # Build and run
make clean            # Clean build artifacts
make debug            # Start with debugger (Delve)
```

### Database

```bash
make create-migration name=your_migration  # Create new migration
make migrate          # Run pending migrations
make migrate-down     # Rollback last migration
make migrate-fresh    # Reset database (dev only!)
make migrate-prune    # Drop all tables (danger!)
```

### GraphQL

```bash
make generate         # Generate GraphQL code from schema
```

### Testing

```bash
go test ./...                              # Run all tests
go test -v ./...                           # Verbose output
go test -cover ./...                       # With coverage
go test ./internal/core/service/...       # Specific package
go test -run TestRegister_Success ./...   # Specific test
```

---

## GraphQL Queries & Mutations

### Authentication

#### Register
```graphql
mutation {
  register(input: {
    name: "John Doe"
    email: "john@example.com"
    password: "SecurePass123!"
  }) {
    accessToken
    refreshToken
  }
}
```

#### Login
```graphql
mutation {
  login(input: {
    email: "john@example.com"
    password: "SecurePass123!"
  }) {
    accessToken
    refreshToken
  }
}
```

#### Refresh Token
```graphql
mutation {
  refreshToken(refreshToken: "v2.local.eyJ...") {
    accessToken
    refreshToken
  }
}
```

#### Logout
```graphql
mutation {
  logout(refreshToken: "v2.local.eyJ...")
}
```

### User Management

#### Get My Profile
```graphql
query {
  me {
    uuid
    name
    email
    createdAt
    updatedAt
  }
}
```
*Requires: Authorization header*

#### Update Profile
```graphql
mutation {
  updateProfile(input: {
    name: "Jane Doe"
  }) {
    uuid
    name
    email
  }
}
```
*Requires: Authorization header*

### Password Management

#### Change Password
```graphql
mutation {
  changePassword(input: {
    oldPassword: "OldPassword123"
    newPassword: "NewPassword456"
  })
}
```
*Requires: Authorization header*

#### Forgot Password
```graphql
mutation {
  forgotPassword(email: "john@example.com")
}
```

#### Verify Email
```graphql
mutation {
  verifyEmail(otp: "123456")
}
```

### Health Check

```graphql
query {
  health
}
```

---

## HTTP Headers

### Public Endpoints
No authentication required:
- `register`
- `login`
- `forgotPassword`
- `verifyEmail`
- `health`

### Protected Endpoints
Requires authentication header:
```
Authorization: Bearer <access_token>
```

Endpoints:
- `me`
- `updateProfile`
- `changePassword`
- `requestEmailChange`
- `deleteMyAccount`
- `refreshToken` (uses refresh token in body)

---

## cURL Examples

### Register
```bash
curl -X POST http://localhost:8000/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation($input: RegisterInput!) { register(input: $input) { accessToken refreshToken } }",
    "variables": {
      "input": {
        "name": "John Doe",
        "email": "john@example.com",
        "password": "SecurePass123!"
      }
    }
  }'
```

### Login
```bash
curl -X POST http://localhost:8000/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation($input: LoginInput!) { login(input: $input) { accessToken refreshToken } }",
    "variables": {
      "input": {
        "email": "john@example.com",
        "password": "SecurePass123!"
      }
    }
  }'
```

### Get Profile (Authenticated)
```bash
curl -X POST http://localhost:8000/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "query": "query { me { uuid name email createdAt updatedAt } }"
  }'
```

---

## Environment Variables

### Required

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=srikandi_sehat
DB_USER=postgres
DB_PASS=your_password

# Token (must be 32 bytes)
TOKEN_SYMMETRIC_KEY=your_32_byte_secret_key_here_now
```

### Optional

```bash
# Application
APP_ENV=development           # or production
APP_PORT=8000
APP_HOST=0.0.0.0
APP_TIMEZONE=Asia/Jakarta

# Token
TOKEN_ISSUER=srikandi-sehat-api
TOKEN_ACCESS_DURATION=15m
TOKEN_REFRESH_DURATION=720h

# Database
DB_SSL_MODE=disable           # use 'require' in production
DB_TIMEZONE=Asia/Jakarta

# Email
MAIL_DRIVER=smtp              # or 'mailgun' for production
MAIL_HOST=localhost
MAIL_PORT=1025
MAIL_FROM_ADDRESS=no-reply@srikandi-sehat.com
MAIL_FROM_NAME=Srikandi Sehat

# Mailgun (production only)
MAILGUN_DOMAIN=mg.yourdomain.com
MAILGUN_API_KEY=your_api_key
```

---

## Docker Commands

### Development

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Restart service
docker-compose restart db

# Execute command in container
docker-compose exec db psql -U postgres
```

### Production

```bash
# Build and start
docker-compose -f docker-compose.prod.yml up -d --build

# View logs
docker-compose -f docker-compose.prod.yml logs -f app

# Stop
docker-compose -f docker-compose.prod.yml down

# Update and restart
git pull
docker-compose -f docker-compose.prod.yml up -d --build
```

---

## Database Commands

### PostgreSQL CLI

```bash
# Connect to database
psql -h localhost -U postgres -d srikandi_sehat

# List databases
\l

# List tables
\dt

# Describe table
\d users

# Run query
SELECT * FROM users;

# Exit
\q
```

### Backup & Restore

```bash
# Backup
pg_dump -U postgres -d srikandi_sehat > backup.sql
pg_dump -U postgres -d srikandi_sehat | gzip > backup.sql.gz

# Restore
psql -U postgres -d srikandi_sehat < backup.sql
gunzip -c backup.sql.gz | psql -U postgres -d srikandi_sehat
```

---

## Common Errors

### Port Already in Use
```bash
# Find process
lsof -i :8000

# Kill process
kill -9 <PID>
```

### Database Connection Failed
```bash
# Check PostgreSQL is running
docker ps

# Check logs
docker logs srikandi_sehat_postgres_db

# Restart database
docker-compose restart db
```

### Migration Failed
```bash
# Rollback
make migrate-down

# Reset (dev only!)
make migrate-fresh
```

### Module Issues
```bash
go clean -modcache
go mod download
go mod tidy
```

---

## Project Structure

```
.
├── cmd/
│   ├── migrate/          # Migration CLI
│   └── server/           # Main server
├── internal/
│   ├── core/
│   │   ├── domain/       # Business entities
│   │   ├── ports/        # Interfaces
│   │   └── service/      # Business logic
│   └── adapters/
│       ├── driven/       # Infrastructure
│       │   ├── postgres/
│       │   └── mailgun/
│       └── driving/      # API
│           └── graphql/
├── pkg/
│   ├── config/           # Configuration
│   ├── logger/           # Logging
│   ├── password/         # Password hashing
│   └── token/            # Token handling
├── docs/                 # Documentation
├── .env.example          # Environment template
├── docker-compose.yml    # Dev infrastructure
├── Makefile              # Build automation
└── README.md
```

---

## Important Files

| File | Purpose |
|------|---------|
| `cmd/server/main.go` | Application entry point |
| `internal/core/domain/*.go` | Domain entities |
| `internal/core/ports/*.go` | Port interfaces |
| `internal/core/service/*.go` | Business logic |
| `internal/adapters/driving/graphql/schema/*.graphqls` | GraphQL schema |
| `internal/adapters/driving/graphql/resolvers/*.go` | GraphQL resolvers |
| `internal/adapters/driven/postgres/*_repository.go` | Database repositories |
| `.env` | Environment variables |
| `Makefile` | Build commands |

---

## Development Workflow

1. **Create feature branch**
   ```bash
   git checkout -b feature/your-feature
   ```

2. **Make changes**
   - Update domain entities if needed
   - Add/update port interfaces
   - Implement service logic
   - Update GraphQL schema
   - Generate GraphQL code: `make generate`
   - Implement resolvers

3. **Create migration if needed**
   ```bash
   make create-migration name=your_migration
   # Edit generated file
   make migrate
   ```

4. **Test changes**
   ```bash
   go test ./...
   make dev  # Manual testing
   ```

5. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: your feature description"
   git push origin feature/your-feature
   ```

6. **Create Pull Request**

---

## Useful Links

- **GraphQL Playground**: http://localhost:8000/
- **Mailpit UI**: http://localhost:8025/
- **PostgreSQL**: localhost:5432

---

## Getting Help

- Check [Documentation](./docs/)
- Review [Architecture Guide](./docs/ARCHITECTURE.md)
- See [API Reference](./docs/API.md)
- Read [Developer Guide](./docs/DEVELOPER.md)
- Check [Feature Docs](./docs/FEATURES.md)
- Review [Deployment Guide](./docs/DEPLOYMENT.md)

---

## Tips & Tricks

### Generate Secure Keys
```bash
# 32-byte key for PASETO
openssl rand -base64 32

# Strong password
openssl rand -base64 24
```

### Pretty Print JSON
```bash
# In GraphQL responses
curl ... | jq '.'
```

### Watch Logs
```bash
# Application logs
sudo journalctl -u srikandi-sehat -f

# Docker logs
docker-compose logs -f app

# Multiple files
tail -f /var/log/nginx/*.log
```

### Database Performance
```sql
-- Check slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Check table sizes
SELECT
  schemaname,
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## Security Checklist

- [ ] Use HTTPS in production
- [ ] Secure database password (20+ characters)
- [ ] Secure token symmetric key (32 bytes)
- [ ] Enable database SSL mode in production
- [ ] Disable GraphQL playground in production
- [ ] Set up rate limiting
- [ ] Configure CORS properly
- [ ] Keep dependencies updated
- [ ] Use environment variables for secrets
- [ ] Enable security headers
- [ ] Set up monitoring and alerting

---

**Version**: 1.0.0  
**Last Updated**: 2025-11-10  
**License**: MIT
