package middleware

import (
	"net/http"

	"github.com/ipincamp/srikandi-sehat/internal/adapters/driving/graphql/dataloader"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// DataloaderMiddleware holds dependencies for the dataloader injector.
type DataloaderMiddleware struct {
	userRepo ports.UserRepository
}

// NewDataloaderMiddleware creates a new dataloader middleware.
// It requires the repository port to pass to the loader constructors.
func NewDataloaderMiddleware(repo ports.UserRepository) *DataloaderMiddleware {
	return &DataloaderMiddleware{
		userRepo: repo,
	}
}

// Handler injects the dataloaders into the context for each request.
func (m *DataloaderMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new UserLoader *for this request*
		userLoader := dataloader.NewUserLoader(m.userRepo)

		// Inject it into the request's context
		ctx := dataloader.WithLoader(r.Context(), userLoader)

		// Call the next middleware/handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
