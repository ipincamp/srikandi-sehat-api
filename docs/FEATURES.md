# Feature Documentation

This document provides detailed information about all implemented features in the Srikandi Sehat GraphQL API.

## Table of Contents

1. [Authentication System](#authentication-system)
2. [User Management](#user-management)
3. [Password Management](#password-management)
4. [Email Services](#email-services)
5. [Security Features](#security-features)
6. [Database Management](#database-management)
7. [Logging and Monitoring](#logging-and-monitoring)

---

## Authentication System

### Overview

The authentication system uses PASETO (Platform-Agnostic SEcurity TOkens) for secure, stateless authentication with separate access and refresh tokens.

### Token Strategy

#### Access Tokens
- **Purpose**: Authorize API requests
- **Lifetime**: 15 minutes (configurable)
- **Usage**: Included in Authorization header
- **Security**: Short-lived to minimize exposure risk

#### Refresh Tokens
- **Purpose**: Obtain new access tokens
- **Lifetime**: 30 days (configurable)
- **Usage**: Submitted to refresh endpoint
- **Security**: Long-lived but only used for token renewal

### Registration Flow

```
┌─────────┐     Register      ┌─────────────┐
│ Client  │ ─────────────────→│   Server    │
└─────────┘                    └─────────────┘
     ↓                                ↓
     │                         1. Validate input
     │                         2. Hash password (Argon2id)
     │                         3. Begin transaction
     │                         4. Check email uniqueness (DB constraint)
     │                         5. Save user
     │                         6. Commit transaction
     │                         7. Generate tokens
     ↓                                ↓
┌─────────┐  Access + Refresh ┌─────────────┐
│ Client  │ ←─────────────────│   Server    │
└─────────┘      Tokens        └─────────────┘
```

**Implementation Details:**

**File:** `internal/core/service/auth_service.go`

```go
func (s *authService) Register(ctx, name, email, password) (*AuthResponse, error) {
    // 1. Hash password with Argon2id
    hashedPassword, err := s.hasher.Hash(password)
    
    // 2. Create domain user
    user := &domain.User{
        UUID:     uuid.NewString(),
        Name:     name,
        Email:    email,
        Password: hashedPassword,
    }
    
    // 3. Transactional save
    tx, err := s.uow.Begin(ctx)
    defer handleRollback(tx, &err)
    
    txUserRepo := tx.GetUserRepository()
    if err := txUserRepo.Save(ctx, user); err != nil {
        if errors.Is(err, ports.ErrDuplicateEmail) {
            return nil, ports.ErrEmailExists
        }
        return nil, err
    }
    
    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }
    
    // 4. Generate token pair
    return s.createTokenSet(user)
}
```

**Key Features:**
- ✅ Email uniqueness enforced by database constraint
- ✅ Transactional safety (all-or-nothing)
- ✅ Secure password hashing (Argon2id)
- ✅ Automatic token generation
- ✅ Race condition prevention via DB constraints

### Login Flow

```
┌─────────┐       Login        ┌─────────────┐
│ Client  │ ─────────────────→│   Server    │
└─────────┘                    └─────────────┘
     ↓                                ↓
     │                         1. Find user by email
     │                         2. Verify password
     │                         3. Generate tokens
     ↓                                ↓
┌─────────┐  Access + Refresh ┌─────────────┐
│ Client  │ ←─────────────────│   Server    │
└─────────┘      Tokens        └─────────────┘
```

**Implementation:**

```go
func (s *authService) Login(ctx, email, password) (*AuthResponse, error) {
    // 1. Find user
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, ports.ErrUserNotFound) {
            return nil, ports.ErrInvalidCredentials
        }
        return nil, err
    }
    
    // 2. Verify password
    if !s.hasher.Compare(user.Password, password) {
        return nil, ports.ErrInvalidCredentials
    }
    
    // 3. Generate tokens
    return s.createTokenSet(user)
}
```

**Security Features:**
- ✅ Generic error messages (prevents email enumeration)
- ✅ Constant-time password comparison
- ✅ Automatic token rotation

### Token Refresh Flow

```
┌─────────┐  Refresh Token    ┌─────────────┐
│ Client  │ ─────────────────→│   Server    │
└─────────┘                    └─────────────┘
     ↓                                ↓
     │                         1. Validate refresh token
     │                         2. Check token type
     │                         3. Find user
     │                         4. Generate new tokens
     ↓                                ↓
┌─────────┐  New Access +     ┌─────────────┐
│ Client  │  Refresh Tokens   │   Server    │
└─────────┘ ←─────────────────└─────────────┘
```

**Implementation:**

```go
func (s *authService) RefreshToken(ctx, refreshToken) (*AuthResponse, error) {
    // 1. Validate token
    payload, err := s.maker.ValidateToken(refreshToken)
    if err != nil {
        if errors.Is(err, token.ErrTokenExpired) {
            return nil, ports.ErrTokenExpired
        }
        return nil, ports.ErrInvalidToken
    }
    
    // 2. Check token type
    if payload.UseFor != token.UseForRefreshToken {
        return nil, ports.ErrTokenUseMismatch
    }
    
    // 3. Find user
    user, err := s.userRepo.FindByID(ctx, payload.UserID)
    if err != nil {
        return nil, err
    }
    
    // 4. Generate new token pair
    return s.createTokenSet(user)
}
```

**Features:**
- ✅ Automatic token rotation
- ✅ Token type validation
- ✅ Expired token detection
- ✅ User existence verification

### Authentication Middleware

**File:** `internal/adapters/driving/graphql/middleware/auth.go`

The middleware implements a "fail-open for public, fail-close for authenticated" strategy:

```go
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
            return // BLOCK
        }
        
        // Check token type
        if payload.UseFor != token.UseForAccessToken {
            http.Error(w, "Invalid token type", 401)
            return // BLOCK
        }
        
        // Inject user ID into context
        ctx := context.WithValue(r.Context(), AuthUserUUIDKey, payload.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Behavior:**
- ✅ Public endpoints (login, register) work without token
- ✅ Protected endpoints require valid token
- ✅ Invalid tokens are rejected immediately
- ✅ User ID available in context for resolvers

---

## User Management

### Get User Profile

Retrieve the authenticated user's profile information.

**Query:**
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

**Implementation:**

**Resolver:** `internal/adapters/driving/graphql/resolvers/user.resolver.go`

```go
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
    // 1. Get user ID from context (set by auth middleware)
    userID, err := GetAuthUserUUID(ctx)
    if err != nil {
        return nil, err
    }
    
    // 2. Fetch user from service
    user, err := r.userService.GetUserByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Convert to GraphQL model
    return toGraphQLUser(user), nil
}
```

**Service:** `internal/core/service/user_service.go`

```go
func (s *userService) GetUserByID(ctx, uuid) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, uuid)
    if err != nil {
        if errors.Is(err, ports.ErrUserNotFound) {
            return nil, ports.ErrUserNotFound
        }
        return nil, err
    }
    return user, nil
}
```

### Update Profile

Update user's profile information (currently name only).

**Mutation:**
```graphql
mutation {
  updateProfile(input: {
    name: "New Name"
  }) {
    uuid
    name
    email
  }
}
```

**Implementation:**

```go
func (s *userService) UpdateProfile(ctx, userID, newName) (*domain.User, error) {
    // 1. Begin transaction
    tx, err := s.uow.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)
    
    // 2. Get user
    txUserRepo := tx.GetUserRepository()
    user, err := txUserRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Update name
    user.Name = newName
    
    // 4. Save
    if err := txUserRepo.Update(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Commit
    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

**Features:**
- ✅ Transactional update
- ✅ Authentication required
- ✅ Returns updated user object

---

## Password Management

### Change Password

Allows authenticated users to change their password.

**Mutation:**
```graphql
mutation {
  changePassword(input: {
    oldPassword: "OldPassword123"
    newPassword: "NewPassword456"
  })
}
```

**Flow:**

```
┌─────────┐  Change Password  ┌─────────────┐
│ Client  │ ─────────────────→│   Server    │
└─────────┘                    └─────────────┘
     ↓                                ↓
     │                         1. Get user from context
     │                         2. Verify old password
     │                         3. Hash new password
     │                         4. Begin transaction
     │                         5. Update password
     │                         6. Commit transaction
     ↓                                ↓
┌─────────┐      Success      ┌─────────────┐
│ Client  │ ←─────────────────│   Server    │
└─────────┘                    └─────────────┘
```

**Implementation:**

```go
func (s *authService) ChangePassword(ctx, userID, oldPassword, newPassword) error {
    // 1. Get user
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. Verify old password
    if !s.hasher.Compare(user.Password, oldPassword) {
        return ports.ErrInvalidCredentials
    }
    
    // 3. Hash new password
    newHashedPassword, err := s.hasher.Hash(newPassword)
    if err != nil {
        return err
    }
    
    // 4. Transactional update
    tx, err := s.uow.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)
    
    txUserRepo := tx.GetUserRepository()
    user.Password = newHashedPassword
    
    if err := txUserRepo.Update(ctx, user); err != nil {
        return err
    }
    
    return tx.Commit(ctx)
}
```

**Security:**
- ✅ Old password verification required
- ✅ New password hashed with Argon2id
- ✅ Transactional update
- ✅ Authentication required

### Forgot Password (OTP-based)

Allows users to reset their password via email OTP.

**Mutation:**
```graphql
mutation {
  forgotPassword(email: "user@example.com")
}
```

**Flow:**

```
┌─────────┐  Forgot Password  ┌─────────────┐
│ Client  │ ─────────────────→│   Server    │
└─────────┘                    └─────────────┘
     ↓                                ↓
     │                         1. Find user by email
     │                         2. Generate 6-digit OTP
     │                         3. Save OTP to database
     │                         4. Send email with OTP
     ↓                                ↓
┌─────────┐      Success      ┌─────────────┐
│ Client  │ ←─────────────────│   Server    │
└─────────┘   (Always true)    └─────────────┘
```

**Implementation:**

```go
func (s *authService) ForgotPassword(ctx, email) error {
    // 1. Find user (but don't reveal if email exists)
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        // Always return success to prevent email enumeration
        return nil
    }
    
    // 2. Generate OTP
    otpCode, err := generateOTP(6)
    if err != nil {
        return err
    }
    
    // 3. Save OTP
    otp := &domain.OTP{
        UserID:    &user.ID,
        Email:     user.Email,
        Code:      otpCode,
        Type:      domain.OTPTypePasswordReset,
        ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
    }
    
    if err := s.otpRepo.Save(ctx, otp); err != nil {
        return err
    }
    
    // 4. Send email
    err = s.emailSvc.SendPasswordResetEmail(ctx, user.Email, user.Name, otpCode)
    if err != nil {
        s.logger.Error().Err(err).Msg("Failed to send reset email")
        // Don't reveal error to user
    }
    
    return nil
}
```

**Security Features:**
- ✅ Always returns success (prevents email enumeration)
- ✅ OTP expires in 15 minutes
- ✅ OTP is single-use (consumed on verification)
- ✅ Secure random OTP generation

**OTP Generation:**

```go
func generateOTP(length int) (string, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return "", err
    }
    
    otpChars := "0123456789"
    for i := range buffer {
        buffer[i] = otpChars[int(buffer[i])%len(otpChars)]
    }
    
    return string(buffer), nil
}
```

---

## Email Services

### Email Service Architecture

**File:** `internal/adapters/driven/mailgun/mail_service.go`

Supports two drivers:

#### 1. SMTP (Development)
- Uses Mailpit for local testing
- No external API required
- Emails viewable at `http://localhost:8025`

#### 2. Mailgun (Production)
- Production-ready email delivery
- High deliverability
- Tracking and analytics

**Configuration:**

```go
type Mail struct {
    Driver        string // "smtp" or "mailgun"
    Host          string
    Port          int
    FromAddress   string
    FromName      string
    MailgunDomain string
    MailgunAPIKey string
}
```

### Email Templates

#### Password Reset Email

```go
func (s *mailService) SendPasswordResetEmail(ctx, email, name, otp) error {
    subject := "Password Reset Request"
    body := fmt.Sprintf(`
Hello %s,

You requested a password reset for your account.

Your OTP code is: %s

This code will expire in 15 minutes.

If you didn't request this, please ignore this email.

Best regards,
Srikandi Sehat Team
    `, name, otp)
    
    return s.send(ctx, email, subject, body)
}
```

#### Email Verification

```go
func (s *mailService) SendEmailVerificationEmail(ctx, email, name, otp) error {
    subject := "Verify Your Email Address"
    body := fmt.Sprintf(`
Hello %s,

Thank you for registering with Srikandi Sehat!

Please verify your email address using this OTP: %s

This code will expire in 15 minutes.

Best regards,
Srikandi Sehat Team
    `, name, otp)
    
    return s.send(ctx, email, subject, body)
}
```

**Features:**
- ✅ Template-based emails
- ✅ Configurable sender information
- ✅ Error handling and logging
- ✅ Context-aware (can be cancelled)

---

## Security Features

### Password Hashing (Argon2id)

**File:** `pkg/password/argon.go`

Argon2id is the recommended algorithm for password hashing:
- Memory-hard (resistant to GPU attacks)
- Time-hard (resistant to brute force)
- Winner of Password Hashing Competition (2015)

**Configuration:**

```go
type Argon2idHasher struct {
    time    uint32 // Number of iterations
    memory  uint32 // Memory in KiB
    threads uint8  // Parallelism
    keyLen  uint32 // Output length
    saltLen uint32 // Salt length
}

func NewArgon2idHasher() *Argon2idHasher {
    return &Argon2idHasher{
        time:    1,      // 1 iteration
        memory:  64*1024, // 64 MB
        threads: 4,      // 4 threads
        keyLen:  32,     // 32 bytes
        saltLen: 16,     // 16 bytes
    }
}
```

**Hash Format:**

```
$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
```

**Implementation:**

```go
func (h *Argon2idHasher) Hash(password string) (string, error) {
    // 1. Generate random salt
    salt := make([]byte, h.saltLen)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }
    
    // 2. Generate hash
    hash := argon2.IDKey(
        []byte(password),
        salt,
        h.time,
        h.memory,
        h.threads,
        h.keyLen,
    )
    
    // 3. Encode in standard format
    return fmt.Sprintf(
        "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version,
        h.memory,
        h.time,
        h.threads,
        base64.RawStdEncoding.EncodeToString(salt),
        base64.RawStdEncoding.EncodeToString(hash),
    ), nil
}

func (h *Argon2idHasher) Compare(hashedPassword, password string) bool {
    // 1. Parse stored hash
    parts := strings.Split(hashedPassword, "$")
    if len(parts) != 6 {
        return false
    }
    
    // 2. Extract parameters
    var memory, time uint32
    var threads uint8
    fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
    
    // 3. Decode salt and hash
    salt, _ := base64.RawStdEncoding.DecodeString(parts[4])
    storedHash, _ := base64.RawStdEncoding.DecodeString(parts[5])
    
    // 4. Hash input password with same parameters
    inputHash := argon2.IDKey(
        []byte(password),
        salt,
        time,
        memory,
        threads,
        uint32(len(storedHash)),
    )
    
    // 5. Constant-time comparison
    return subtle.ConstantTimeCompare(storedHash, inputHash) == 1
}
```

**Security Benefits:**
- ✅ Cryptographically secure random salt
- ✅ Memory-hard (expensive to parallelize)
- ✅ Configurable work factor
- ✅ Constant-time comparison (prevents timing attacks)

### PASETO Tokens

**File:** `pkg/token/paseto.go`

PASETO (Platform-Agnostic SEcurity TOkens) v2 provides:
- Authenticated encryption (no need for separate signature)
- Eliminates common JWT vulnerabilities
- Simpler to use correctly

**Token Structure:**

```
v2.local.<encrypted-payload>
```

**Payload:**

```go
type Payload struct {
    UserID    string
    IssuedAt  time.Time
    ExpiredAt time.Time
    UseFor    string // "access" or "refresh"
}
```

**Implementation:**

```go
func (maker *PasetoMaker) CreateToken(userID string, duration time.Duration, useFor string) (string, error) {
    payload := &Payload{
        UserID:    userID,
        IssuedAt:  time.Now(),
        ExpiredAt: time.Now().Add(duration),
        UseFor:    useFor,
    }
    
    return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) ValidateToken(token string) (*Payload, error) {
    payload := &Payload{}
    
    err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    if time.Now().After(payload.ExpiredAt) {
        return nil, ErrTokenExpired
    }
    
    return payload, nil
}
```

**Benefits:**
- ✅ No algorithm confusion attacks
- ✅ Authenticated encryption
- ✅ Simple API
- ✅ Version-specific implementation

---

## Database Management

### Unit of Work Pattern

Ensures multiple operations execute within a single transaction:

**File:** `internal/adapters/driven/postgres/unit_of_work.go`

```go
type unitOfWork struct {
    pool   *pgxpool.Pool
    logger zerolog.Logger
}

func (uow *unitOfWork) Begin(ctx context.Context) (ports.Transaction, error) {
    tx, err := uow.pool.Begin(ctx)
    if err != nil {
        return nil, err
    }
    
    return &transaction{
        tx:     tx,
        logger: uow.logger,
    }, nil
}

type transaction struct {
    tx     pgx.Tx
    logger zerolog.Logger
}

func (t *transaction) GetUserRepository() ports.UserRepository {
    db := NewGormDB(t.tx, t.logger)
    return NewUserRepository(db, t.logger)
}

func (t *transaction) Commit(ctx context.Context) error {
    return t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) error {
    return t.tx.Rollback(ctx)
}
```

**Usage Pattern:**

```go
// Begin transaction
tx, err := uow.Begin(ctx)
if err != nil {
    return err
}

// Automatic rollback on error
defer func() {
    if err != nil {
        tx.Rollback(ctx)
    }
}()

// Get transactional repositories
userRepo := tx.GetUserRepository()
otpRepo := tx.GetOTPRepository()

// Perform operations
userRepo.Save(ctx, user)
otpRepo.Save(ctx, otp)

// Commit
return tx.Commit(ctx)
```

### DataLoader (N+1 Prevention)

**File:** `internal/adapters/driving/graphql/dataloader/loader.go`

Batches multiple database queries into single queries:

```go
type Loaders struct {
    UserByUUID *dataloader.Loader
}

func NewLoaders(userRepo ports.UserRepository) *Loaders {
    return &Loaders{
        UserByUUID: dataloader.NewBatchedLoader(
            func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
                // Convert keys to string slice
                uuids := keys.Keys()
                
                // Single batch query
                usersMap, err := userRepo.FindMapByUUIDs(ctx, uuids)
                if err != nil {
                    // Return error for all
                    return createErrorResults(len(keys), err)
                }
                
                // Return results in same order as keys
                results := make([]*dataloader.Result, len(keys))
                for i, key := range keys.Keys() {
                    user, ok := usersMap[key]
                    if !ok {
                        results[i] = &dataloader.Result{
                            Error: ports.ErrUserNotFound,
                        }
                    } else {
                        results[i] = &dataloader.Result{
                            Data: user,
                        }
                    }
                }
                return results
            },
        ),
    }
}
```

**Usage:**

```go
// Without DataLoader (N+1 problem):
// Query 1: SELECT * FROM posts
// Query 2: SELECT * FROM users WHERE uuid = 'user1'
// Query 3: SELECT * FROM users WHERE uuid = 'user2'
// Query 4: SELECT * FROM users WHERE uuid = 'user3'
// ... (N queries)

// With DataLoader (batched):
// Query 1: SELECT * FROM posts
// Query 2: SELECT * FROM users WHERE uuid IN ('user1', 'user2', 'user3', ...)
// (2 queries total)
```

---

## Logging and Monitoring

### Structured Logging

**File:** `pkg/logger/logger.go`

Uses zerolog for high-performance structured logging:

```go
func NewLogger(env string) zerolog.Logger {
    var logger zerolog.Logger
    
    if env == "production" {
        // JSON output for production
        logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
    } else {
        // Pretty console output for development
        output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
        logger = zerolog.New(output).With().Timestamp().Logger()
    }
    
    return logger
}
```

**Usage:**

```go
// Info level
logger.Info().
    Str("user_id", userID).
    Str("email", email).
    Msg("User registered successfully")

// Error level
logger.Error().
    Err(err).
    Str("user_id", userID).
    Msg("Failed to update user")

// With context
log := logger.With().
    Str("component", "AuthService").
    Logger()
```

**Output (Development):**
```
2025-11-10T10:00:00Z INF User registered successfully user_id=uuid-123 email=user@example.com
```

**Output (Production):**
```json
{
  "level": "info",
  "time": "2025-11-10T10:00:00Z",
  "message": "User registered successfully",
  "user_id": "uuid-123",
  "email": "user@example.com"
}
```

### Health Check

**Query:**
```graphql
query {
  health
}
```

**Response:**
```json
{
  "data": {
    "health": "OK"
  }
}
```

Simple endpoint for monitoring service availability.

---

## Summary

This documentation covers all implemented features in the Srikandi Sehat GraphQL API. Each feature is designed with security, scalability, and maintainability in mind, following hexagonal architecture principles.

For implementation details, see:
- [Architecture Guide](./ARCHITECTURE.md)
- [API Reference](./API.md)
- [Developer Guide](./DEVELOPER.md)
