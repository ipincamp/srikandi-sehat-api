package postgres

import (
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
)

// User is the database model for a user, optimized for GORM/pgx.
// It is kept separate from the core domain.User model to isolate
// database-specific tags and conventions.
type User struct {
	ID        uint      `db:"id"`
	UUID      string    `db:"uuid"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// toDomain converts the database model (postgres.User)
// into the core business model (domain.User).
func (dbUser *User) toDomain() *domain.User {
	return &domain.User{
		ID:        dbUser.ID,
		UUID:      dbUser.UUID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

// fromDomain converts the core business model (domain.User)
// into the database model (postgres.User) for saving.
func fromDomain(dUser *domain.User) *User {
	return &User{
		ID:        dUser.ID,
		UUID:      dUser.UUID,
		Name:      dUser.Name,
		Email:     dUser.Email,
		Password:  dUser.Password,
		CreatedAt: dUser.CreatedAt,
		UpdatedAt: dUser.UpdatedAt,
	}
}
