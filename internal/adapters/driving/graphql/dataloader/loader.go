package dataloader

import (
	"context"
	"fmt"
	"time"

	"github.com/graph-gophers/dataloader"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Define a context key for storing the loader
type contextKey string

const loaderKey contextKey = "dataloader"

// userLoader holds the batching function and repository.
// This is intentionally *not* exported.
type userLoader struct {
	repo ports.UserRepository
}

// NewUserLoader creates a new Dataloader instance for fetching users.
// It batches requests over a 2ms window.
func NewUserLoader(repo ports.UserRepository) *dataloader.Loader {
	// Create our internal loader struct
	loader := &userLoader{
		repo: repo,
	}
	// Use the dataloader library to wrap our batch function
	return dataloader.NewBatchedLoader(loader.batchGetUser, dataloader.WithWait((2 * time.Millisecond)))
}

// batchGetUser is the batching function executed by the dataloader.
func (l *userLoader) batchGetUser(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// 1. Get the list of user UUIDs to fetch
	uuids := make([]string, len(keys))
	for i, key := range keys {
		uuids[i] = key.String()
	}

	// 2. Call our efficient, batch-capable repository port
	userMap, err := l.repo.FindMapByUUIDs(ctx, uuids)

	// 3. Create the result slice
	results := make([]*dataloader.Result, len(keys))
	for i, key := range keys {
		// 4. Handle repository-level error
		if err != nil {
			results[i] = &dataloader.Result{Error: err}
			continue
		}

		// 5. Map the result back to the correct key
		uuid := key.String()
		if user, ok := userMap[uuid]; ok {
			// User was found
			results[i] = &dataloader.Result{Data: user}
		} else {
			// User was not found for this key
			results[i] = &dataloader.Result{
				Error: fmt.Errorf("user with uuid '%s' not found", uuid),
			}
		}
	}
	return results
}

// For retrieves the Dataloader from the context.
func For(ctx context.Context) *dataloader.Loader {
	return ctx.Value(loaderKey).(*dataloader.Loader)
}

// WithLoader injects the Dataloader into the context.
func WithLoader(ctx context.Context, loader *dataloader.Loader) context.Context {
	return context.WithValue(ctx, loaderKey, loader)
}
