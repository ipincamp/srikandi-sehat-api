package middleware

import (
	"context"
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

// Handler is the actual HTTP middleware function.
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token, proceed without authentication.
			// The 'Me' resolver will block this.
			next.ServeHTTP(w, r)
			return
		}

		// 2. Validate header format "Bearer <token>"
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			// Invalid format, proceed without auth
			next.ServeHTTP(w, r)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			// Unsupported auth type
			next.ServeHTTP(w, r)
			return
		}

		// 3. Validate the token
		accessToken := fields[1]
		payload, err := am.tokenMaker.ValidateToken(accessToken)
		if err != nil {
			// Invalid or expired token
			next.ServeHTTP(w, r)
			return
		}

		// 4. Check if it's an access token
		if payload.UseFor != token.UseForAccessToken {
			// Wrong token type (e.g., a refresh token)
			next.ServeHTTP(w, r)
			return
		}

		// 5. SUCCESS! Inject the user UUID into the context
		// This is the key that the 'Me' resolver is looking for.
		ctx := context.WithValue(r.Context(), resolvers.AuthUserUUIDKey, payload.UserID)

		// 6. Call the next handler in the chain with the *new* context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
