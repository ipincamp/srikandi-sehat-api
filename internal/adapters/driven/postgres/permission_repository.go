package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Compile-time interface implementation check.
var _ ports.PermissionRepository = (*permissionRepository)(nil)

// permissionRepository implements ports.PermissionRepository interface using PostgreSQL.
type permissionRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewPermissionRepository creates a new permission repository instance.
func NewPermissionRepository(db dbExecutor, logger zerolog.Logger) ports.PermissionRepository {
	return &permissionRepository{
		db:     db,
		logger: logger,
	}
}

// FindByID retrieves a permission by its unique identifier.
func (r *permissionRepository) FindByID(ctx context.Context, permissionID string) (*domain.Permission, error) {
	query := `SELECT id, name, description FROM permissions WHERE id = $1`

	permission := &domain.Permission{}
	err := r.db.QueryRow(ctx, query, permissionID).Scan(&permission.ID, &permission.Name, &permission.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("permission_id", permissionID).Msg("Permission not found")
			return nil, ports.ErrPermissionNotFound
		}
		r.logger.Error().Err(err).Str("permission_id", permissionID).Msg("Error finding permission by ID")
		return nil, err
	}

	return permission, nil
}

// FindByName retrieves a permission by its name.
func (r *permissionRepository) FindByName(ctx context.Context, name string) (*domain.Permission, error) {
	query := `SELECT id, name, description FROM permissions WHERE name = $1`

	permission := &domain.Permission{}
	err := r.db.QueryRow(ctx, query, name).Scan(&permission.ID, &permission.Name, &permission.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("name", name).Msg("Permission not found")
			return nil, ports.ErrPermissionNotFound
		}
		r.logger.Error().Err(err).Str("name", name).Msg("Error finding permission by name")
		return nil, err
	}

	return permission, nil
}

// List retrieves all permissions in the system.
func (r *permissionRepository) List(ctx context.Context) ([]*domain.Permission, error) {
	query := `SELECT id, name, description FROM permissions ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list permissions")
		return nil, err
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		permission := &domain.Permission{}
		if err := rows.Scan(&permission.ID, &permission.Name, &permission.Description); err != nil {
			r.logger.Error().Err(err).Msg("Error scanning permission")
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// Save persists a new permission or updates an existing permission in the database.
func (r *permissionRepository) Save(ctx context.Context, permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (id, name, description)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET name = $2, description = $3
	`

	_, err := r.db.Exec(ctx, query, permission.ID, permission.Name, permission.Description)
	if err != nil {
		r.logger.Error().Err(err).Str("permission_id", permission.ID).Msg("Failed to save permission")
		return err
	}

	r.logger.Debug().Str("permission_id", permission.ID).Msg("Permission saved successfully")
	return nil
}
