package postgres

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// Compile-time interface implementation check.
var _ ports.ActivityLogRepository = (*activityLogRepository)(nil)

// activityLogRepository implements ports.ActivityLogRepository interface using PostgreSQL.
type activityLogRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

// NewActivityLogRepository creates a new activity log repository instance.
func NewActivityLogRepository(db dbExecutor, logger zerolog.Logger) ports.ActivityLogRepository {
	return &activityLogRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists a new activity log entry in the database.
func (r *activityLogRepository) Save(ctx context.Context, log *domain.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (id, actor_user_id, action, target_table, target_id, changes, ip_address, user_agent, request_id, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var targetID sql.NullString
	if log.TargetID != nil {
		targetID = sql.NullString{String: *log.TargetID, Valid: true}
	}

	var changes sql.NullString
	if log.Changes != nil {
		changes = sql.NullString{String: *log.Changes, Valid: true}
	}

	var ipAddress sql.NullString
	if log.IPAddress != nil {
		ipAddress = sql.NullString{String: *log.IPAddress, Valid: true}
	}

	var userAgent sql.NullString
	if log.UserAgent != nil {
		userAgent = sql.NullString{String: *log.UserAgent, Valid: true}
	}

	var requestID sql.NullString
	if log.RequestID != nil {
		requestID = sql.NullString{String: *log.RequestID, Valid: true}
	}

	_, err := r.db.Exec(ctx, query,
		log.ID,
		log.ActorUserID,
		log.Action,
		log.TargetTable,
		targetID,
		changes,
		ipAddress,
		userAgent,
		requestID,
		log.Timestamp,
	)

	if err != nil {
		r.logger.Error().Err(err).Int64("log_id", log.ID).Msg("Failed to save activity log")
		return err
	}

	r.logger.Debug().Int64("log_id", log.ID).Msg("Activity log saved successfully")
	return nil
}

// ListByUserID retrieves all activity logs for a specific user.
func (r *activityLogRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.ActivityLog, error) {
	query := `
		SELECT id, actor_user_id, action, target_table, target_id, changes, ip_address, user_agent, request_id, timestamp
		FROM activity_logs
		WHERE actor_user_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error().Err(err).Str("user_id", userID).Msg("Failed to list activity logs by user ID")
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.ActivityLog
	for rows.Next() {
		log := &domain.ActivityLog{}
		var targetID, changes, ipAddress, userAgent, requestID sql.NullString

		if err := rows.Scan(
			&log.ID,
			&log.ActorUserID,
			&log.Action,
			&log.TargetTable,
			&targetID,
			&changes,
			&ipAddress,
			&userAgent,
			&requestID,
			&log.Timestamp,
		); err != nil {
			r.logger.Error().Err(err).Msg("Error scanning activity log")
			return nil, err
		}

		if targetID.Valid {
			log.TargetID = &targetID.String
		}
		if changes.Valid {
			log.Changes = &changes.String
		}
		if ipAddress.Valid {
			log.IPAddress = &ipAddress.String
		}
		if userAgent.Valid {
			log.UserAgent = &userAgent.String
		}
		if requestID.Valid {
			log.RequestID = &requestID.String
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}
