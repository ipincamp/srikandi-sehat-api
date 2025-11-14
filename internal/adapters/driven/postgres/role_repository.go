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
var _ ports.RoleRepository = (*roleRepository)(nil)

// roleRepository implements ports.RoleRepository interface using PostgreSQL.
type roleRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewRoleRepository creates a new role repository instance.
func NewRoleRepository(db dbExecutor, logger zerolog.Logger) ports.RoleRepository {
	return &roleRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a new role or updates an existing role in the database.
func (r *roleRepository) Save(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (id, name, description)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET name = $2, description = $3
	`

	_, err := r.db.Exec(ctx, query, role.ID, role.Name, role.Description)
	if err != nil {
		r.logger.Error().Err(err).Str("role_id", role.ID).Msg("Failed to save role")
		return err
	}

	r.logger.Debug().Str("role_id", role.ID).Msg("Role saved successfully")
	return nil
}

// FindByID retrieves a role by its unique identifier.
func (r *roleRepository) FindByID(ctx context.Context, roleID string) (*domain.Role, error) {
	query := `SELECT id, name, description FROM roles WHERE id = $1`

	role := &domain.Role{}
	err := r.db.QueryRow(ctx, query, roleID).Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("role_id", roleID).Msg("Role not found")
			return nil, ports.ErrRoleNotFound
		}
		r.logger.Error().Err(err).Str("role_id", roleID).Msg("Error finding role by ID")
		return nil, err
	}

	return role, nil
}

// FindByName retrieves a role by its name.
func (r *roleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = $1`

	role := &domain.Role{}
	err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug().Str("name", name).Msg("Role not found")
			return nil, ports.ErrRoleNotFound
		}
		r.logger.Error().Err(err).Str("name", name).Msg("Error finding role by name")
		return nil, err
	}

	return role, nil
}

// Delete removes a role from the database by its ID.
func (r *roleRepository) Delete(ctx context.Context, roleID string) error {
	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.db.Exec(ctx, query, roleID)
	if err != nil {
		r.logger.Error().Err(err).Str("role_id", roleID).Msg("Failed to delete role")
		return err
	}

	if result.RowsAffected() == 0 {
		return ports.ErrRoleNotFound
	}

	r.logger.Debug().Str("role_id", roleID).Msg("Role deleted successfully")
	return nil
}

// List retrieves all roles in the system.
func (r *roleRepository) List(ctx context.Context) ([]*domain.Role, error) {
	query := `SELECT id, name, description FROM roles ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list roles")
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			r.logger.Error().Err(err).Msg("Error scanning role")
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

// AddPermissionToRole associates a permission with a role.
func (r *roleRepository) AddPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		r.logger.Error().Err(err).Str("role_id", roleID).Str("permission_id", permissionID).Msg("Failed to add permission to role")
		return err
	}

	r.logger.Debug().Str("role_id", roleID).Str("permission_id", permissionID).Msg("Permission added to role")
	return nil
}

// RemovePermissionFromRole removes a permission from a role.
func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	result, err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		r.logger.Error().Err(err).Str("role_id", roleID).Str("permission_id", permissionID).Msg("Failed to remove permission from role")
		return err
	}

	if result.RowsAffected() == 0 {
		return ports.ErrPermissionNotFound
	}

	r.logger.Debug().Str("role_id", roleID).Str("permission_id", permissionID).Msg("Permission removed from role")
	return nil
}
