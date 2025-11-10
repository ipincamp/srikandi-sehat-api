# API Documentation

## GraphQL Endpoint

- **URL**: `http://localhost:8000/query`
- **Playground**: `http://localhost:8000/`

## Authentication

Most mutations and queries require authentication using a Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Token Types

1. **Access Token** - Short-lived (15 minutes), used for API requests
2. **Refresh Token** - Long-lived (30 days), used to obtain new access tokens

## Queries

### 1. Health Check

Check API health status.

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

**Authentication:** Not required

---

### 2. Get My Profile

Get the authenticated user's profile information.

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

**Response:**
```json
{
  "data": {
    "me": {
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "john@example.com",
      "createdAt": "2025-11-10T10:00:00Z",
      "updatedAt": "2025-11-10T10:00:00Z"
    }
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Errors:**
- `401 Unauthorized` - Invalid or expired access token
- `404 Not Found` - User not found

---

## Mutations

### 1. Register

Create a new user account.

**Mutation:**
```graphql
mutation Register($input: RegisterInput!) {
  register(input: $input) {
    accessToken
    refreshToken
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }
}
```

**Response:**
```json
{
  "data": {
    "register": {
      "accessToken": "v2.local.eyJ...",
      "refreshToken": "v2.local.eyJ..."
    }
  }
}
```

**Validation Rules:**
- `name`: Required, non-empty string
- `email`: Required, valid email format
- `password`: Required, minimum 8 characters recommended

**Errors:**
- `409 Conflict` - Email already exists
- `400 Bad Request` - Invalid input data

---

### 2. Login

Authenticate with email and password.

**Mutation:**
```graphql
mutation Login($input: LoginInput!) {
  login(input: $input) {
    accessToken
    refreshToken
  }
}
```

**Variables:**
```json
{
  "input": {
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }
}
```

**Response:**
```json
{
  "data": {
    "login": {
      "accessToken": "v2.local.eyJ...",
      "refreshToken": "v2.local.eyJ..."
    }
  }
}
```

**Errors:**
- `401 Unauthorized` - Invalid email or password
- `400 Bad Request` - Invalid input data

---

### 3. Refresh Token

Obtain a new pair of tokens using a refresh token.

**Mutation:**
```graphql
mutation RefreshToken($refreshToken: String!) {
  refreshToken(refreshToken: $refreshToken) {
    accessToken
    refreshToken
  }
}
```

**Variables:**
```json
{
  "refreshToken": "v2.local.eyJ..."
}
```

**Response:**
```json
{
  "data": {
    "refreshToken": {
      "accessToken": "v2.local.eyJ...",
      "refreshToken": "v2.local.eyJ..."
    }
  }
}
```

**Errors:**
- `401 Unauthorized` - Invalid or expired refresh token
- `401 Unauthorized` - Token type mismatch (access token provided)

---

### 4. Logout

Logout the user (client should delete tokens).

**Mutation:**
```graphql
mutation Logout($refreshToken: String!) {
  logout(refreshToken: $refreshToken)
}
```

**Variables:**
```json
{
  "refreshToken": "v2.local.eyJ..."
}
```

**Response:**
```json
{
  "data": {
    "logout": true
  }
}
```

**Note:** In the current stateless implementation, this is a no-op. The client is responsible for deleting tokens. Future versions may implement token blacklisting.

---

### 5. Change Password

Change the authenticated user's password.

**Mutation:**
```graphql
mutation ChangePassword($input: ChangePasswordInput!) {
  changePassword(input: $input)
}
```

**Variables:**
```json
{
  "input": {
    "oldPassword": "OldPassword123!",
    "newPassword": "NewSecurePassword456!"
  }
}
```

**Response:**
```json
{
  "data": {
    "changePassword": true
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Errors:**
- `401 Unauthorized` - Invalid or expired access token
- `401 Unauthorized` - Old password is incorrect
- `400 Bad Request` - New password doesn't meet requirements

---

### 6. Update Profile

Update the authenticated user's profile information.

**Mutation:**
```graphql
mutation UpdateProfile($input: UpdateProfileInput!) {
  updateProfile(input: $input) {
    uuid
    name
    email
    createdAt
    updatedAt
  }
}
```

**Variables:**
```json
{
  "input": {
    "name": "Jane Doe"
  }
}
```

**Response:**
```json
{
  "data": {
    "updateProfile": {
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Jane Doe",
      "email": "john@example.com",
      "createdAt": "2025-11-10T10:00:00Z",
      "updatedAt": "2025-11-10T12:30:00Z"
    }
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Errors:**
- `401 Unauthorized` - Invalid or expired access token
- `404 Not Found` - User not found
- `400 Bad Request` - Invalid input data

---

### 7. Forgot Password

Request a password reset OTP via email.

**Mutation:**
```graphql
mutation ForgotPassword($email: String!) {
  forgotPassword(email: $email)
}
```

**Variables:**
```json
{
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "data": {
    "forgotPassword": true
  }
}
```

**Note:** 
- Always returns `true` to prevent email enumeration attacks
- OTP is sent to email if the address exists
- OTP is valid for 15 minutes

**Email Example:**
```
Subject: Password Reset Request

Hello John,

You requested a password reset. Your OTP code is: 123456

This code will expire in 15 minutes.

If you didn't request this, please ignore this email.
```

---

### 8. Verify Email (Blueprint)

Verify email address using OTP code.

**Mutation:**
```graphql
mutation VerifyEmail($otp: String!) {
  verifyEmail(otp: $otp)
}
```

**Variables:**
```json
{
  "otp": "123456"
}
```

**Response:**
```json
{
  "data": {
    "verifyEmail": true
  }
}
```

**Errors:**
- `400 Bad Request` - Invalid or expired OTP
- `404 Not Found` - OTP not found

**Note:** This feature sends an OTP during registration or email change verification.

---

### 9. Request Email Change (Blueprint)

Request to change email address (sends OTP to new email).

**Mutation:**
```graphql
mutation RequestEmailChange($newEmail: String!) {
  requestEmailChange(newEmail: $newEmail)
}
```

**Variables:**
```json
{
  "newEmail": "newemail@example.com"
}
```

**Response:**
```json
{
  "data": {
    "requestEmailChange": true
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Errors:**
- `401 Unauthorized` - Invalid or expired access token
- `409 Conflict` - New email already in use
- `400 Bad Request` - Invalid email format

**Note:** OTP will be sent to the new email address for verification.

---

### 10. Delete My Account (Blueprint)

Delete the authenticated user's account.

**Mutation:**
```graphql
mutation DeleteMyAccount {
  deleteMyAccount
}
```

**Response:**
```json
{
  "data": {
    "deleteMyAccount": true
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Errors:**
- `401 Unauthorized` - Invalid or expired access token
- `404 Not Found` - User not found

**Note:** This performs a soft delete. Data is marked as deleted but retained in the database.

---

## Error Handling

All errors follow GraphQL error format:

```json
{
  "errors": [
    {
      "message": "Error message here",
      "path": ["mutation", "register"],
      "extensions": {
        "code": "ERROR_CODE"
      }
    }
  ],
  "data": null
}
```

### Common Error Codes

- `UNAUTHENTICATED` - Missing or invalid authentication
- `FORBIDDEN` - Authenticated but not authorized
- `BAD_USER_INPUT` - Invalid input data
- `INTERNAL_SERVER_ERROR` - Unexpected server error
- `NOT_FOUND` - Resource not found

### HTTP Status Codes

- `200 OK` - Successful GraphQL request (even if contains errors)
- `400 Bad Request` - Malformed GraphQL query
- `401 Unauthorized` - Authentication failed (middleware level)
- `500 Internal Server Error` - Server error

---

## Data Types

### User

```graphql
type User {
  uuid: String!
  name: String!
  email: String!
  createdAt: String!  # ISO 8601 format
  updatedAt: String!  # ISO 8601 format
}
```

### AuthResponse

```graphql
type AuthResponse {
  accessToken: String!
  refreshToken: String!
}
```

---

## Rate Limiting

**Recommendation:** Implement rate limiting at the reverse proxy level (e.g., Nginx, Cloudflare).

Suggested limits:
- **Register/Login**: 5 requests per minute per IP
- **Password Reset**: 3 requests per hour per email
- **General API**: 100 requests per minute per user

---

## Best Practices

### 1. Token Management

```javascript
// Store tokens securely
localStorage.setItem('accessToken', response.accessToken);
localStorage.setItem('refreshToken', response.refreshToken);

// Include in requests
const headers = {
  'Authorization': `Bearer ${localStorage.getItem('accessToken')}`
};

// Refresh when expired
if (error.extensions.code === 'UNAUTHENTICATED') {
  const newTokens = await refreshToken();
  localStorage.setItem('accessToken', newTokens.accessToken);
  // Retry original request
}

// Clear on logout
localStorage.removeItem('accessToken');
localStorage.removeItem('refreshToken');
```

### 2. Error Handling

```javascript
try {
  const result = await graphqlClient.query(query);
  return result.data;
} catch (error) {
  if (error.graphQLErrors) {
    error.graphQLErrors.forEach(({ message, extensions }) => {
      console.error(`GraphQL Error: ${message} (${extensions.code})`);
    });
  }
  if (error.networkError) {
    console.error('Network Error:', error.networkError);
  }
}
```

### 3. Input Validation

Always validate on both client and server:

```javascript
// Client-side validation
function validateEmail(email) {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return regex.test(email);
}

function validatePassword(password) {
  return password.length >= 8;
}
```

### 4. Secure Password Storage

**Never:**
- Store passwords in plain text
- Log passwords
- Send passwords in URLs
- Store passwords in cookies

**Always:**
- Use HTTPS in production
- Hash passwords server-side (Argon2id)
- Use secure password policies

---

## Testing with cURL

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

## Versioning

This API follows semantic versioning. Breaking changes will be announced in advance and deprecated endpoints will be maintained for at least 6 months.

Current version: **v1.0.0**

---

## Support

For issues or questions:
- Create an issue on GitHub
- Contact the development team
- Check the documentation at `/docs`

---

## Changelog

### v1.0.0 (2025-11-10)
- Initial release
- User registration and authentication
- Profile management
- Password reset functionality
- Email verification (blueprint)
