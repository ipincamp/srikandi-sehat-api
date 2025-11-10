package domain

import "time"

// Tipe OTP untuk membedakan tujuan
const (
	OTPTypeVerification    = "email_verification"
	OTPTypePasswordReset   = "password_reset"
	OTPTypeEmailChange     = "email_change"
	OTPTypeAccountDeletion = "account_deletion"
)

// OTP adalah entitas domain untuk One-Time Password.
type OTP struct {
	ID        uint
	UserID    *uint  // Nullable, untuk OTP verifikasi email sebelum user dibuat
	Email     string // Untuk verifikasi
	Code      string // Kode 6 digit
	Type      string // Tipe OTP (lihat konstanta di atas)
	ExpiresAt time.Time
	CreatedAt time.Time
}
