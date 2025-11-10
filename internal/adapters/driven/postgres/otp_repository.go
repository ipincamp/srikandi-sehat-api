package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// Compile-time check
var _ ports.OTPRepository = (*otpRepository)(nil)

type otpRepository struct {
	db     dbExecutor
	logger zerolog.Logger
}

func NewOTPRepository(db dbExecutor, logger zerolog.Logger) ports.OTPRepository {
	return &otpRepository{
		db:     db,
		logger: logger,
	}
}

// Save menyimpan OTP baru ke database.
func (r *otpRepository) Save(ctx context.Context, otp *domain.OTP) error {
	dbOTP := otpFromDomain(otp)
	dbOTP.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO otps (user_id, email, code, type, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query,
		dbOTP.UserID,
		dbOTP.Email,
		dbOTP.Code,
		dbOTP.Type,
		dbOTP.ExpiresAt,
		dbOTP.CreatedAt,
	).Scan(&dbOTP.ID)

	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to save OTP")
		return err
	}

	otp.ID = dbOTP.ID // Set ID yang di-generate DB
	return nil
}

// FindAndConsume mencari OTP, memvalidasi, dan menghapusnya.
// Ini adalah operasi atomik menggunakan `DELETE ... RETURNING`.
func (r *otpRepository) FindAndConsume(ctx context.Context, code string, otpType string) (*domain.OTP, error) {
	query := `
		DELETE FROM otps
		WHERE code = $1 AND type = $2 AND expires_at > $3
		RETURNING id, user_id, email, code, type, expires_at, created_at
	`
	now := time.Now().UTC()
	dbOTP := &OTP{}

	err := r.db.QueryRow(ctx, query, code, otpType, now).Scan(
		&dbOTP.ID,
		&dbOTP.UserID,
		&dbOTP.Email,
		&dbOTP.Code,
		&dbOTP.Type,
		&dbOTP.ExpiresAt,
		&dbOTP.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn().Str("type", otpType).Msg("OTP not found, expired, or already used")
			// Gunakan error port yang generik
			return nil, ports.ErrInvalidToken
		}
		r.logger.Error().Err(err).Msg("Error finding/consuming OTP")
		return nil, ports.ErrUnexpectedFind
	}

	r.logger.Debug().Str("type", otpType).Msg("OTP consumed successfully")
	return dbOTP.toDomain(), nil
}
