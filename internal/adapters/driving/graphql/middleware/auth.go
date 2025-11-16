package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/resolvers"
	"github.com/ipincamp/srikandi-sehat/pkg/token"
)

// AuthMiddleware is a struct that holds dependencies for the middleware.
type AuthMiddleware struct {
	tokenMaker token.Maker
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(maker token.Maker) *AuthMiddleware {
	return &AuthMiddleware{
		tokenMaker: maker,
	}
}

// writeError is a helper to return a standard HTTP error.
// We are in middleware, so we must respond with HTTP, not GraphQL JSON.
func writeError(w http.ResponseWriter, message string, statusCode int) {
	// Set the content type to JSON, as GraphQL clients expect this.
	w.Header().Set("Content-Type", "application/json")
	// Set the HTTP status code (e.g., 401 Unauthorized).
	w.WriteHeader(statusCode)

	// Construct a standard GraphQL error object.
	// This structure {"errors": [{"message": "..."}]} is what the client (Playground) expects.
	errResponse := map[string]any{
		"errors": []map[string]any{
			{
				"message": message,
			},
		},
		"data": nil,
	}

	// Encode the error map to JSON and write it to the response.
	json.NewEncoder(w).Encode(errResponse)
}

// Handler is the actual HTTP middleware function.
// This middleware implements a "fail-open" for public requests
// (no token) and a "fail-close" for authenticated requests (token present).
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the Authorization header
		authHeader := r.Header.Get("Authorization")

		// 2. No Authorization header found.
		// This is a public request (e.g., login, register, public queries).
		// We proceed without injecting a user.
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// 3. Validate header format "Bearer <token>"
		// We *must* validate it ("fail-close").
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			writeError(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return // BLOCK
		}

		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			writeError(w, "Unsupported Authorization type, must be Bearer", http.StatusUnauthorized)
			return // BLOCK
		}

		// 4. Validate the token itself
		accessToken := fields[1]
		payload, err := am.tokenMaker.ValidateToken(accessToken)
		if err != nil {
			// Token is invalid (expired, bad signature, etc.)
			writeError(w, "Invalid or expired token", http.StatusUnauthorized)
			return // BLOCK
		}

		// 5. Check if it's an access token
		if payload.UseFor != token.UseForAccessToken {
			writeError(w, "Invalid token type provided", http.StatusUnauthorized)
			return // BLOCK
		}

		// 6. SUCCESS! Inject the user UUID into the context
		// This is the key that the 'Me' resolver is looking for.
		ctx := context.WithValue(r.Context(), resolvers.AuthUserUUIDKey, payload.UserID)

		// 7. Call the next handler in the chain with the *new* context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
