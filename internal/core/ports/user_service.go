package ports

import (
	"context"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// UserService mendefinisikan logika bisnis untuk User.
type UserService interface {
	// GetByID mengambil satu user berdasarkan ID.
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
