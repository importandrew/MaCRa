package models

import (
	"database/sql"
	"time"
)

type AuthPurpose string

const (
	AuthPurposeEmailVerify   AuthPurpose = "email_verify"
	AuthPurposePasswordReset AuthPurpose = "password_reset"

	EmailVerifyTTL   = 24 * time.Hour
	PasswordResetTTL = time.Hour
)

type AuthToken struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	Purpose   AuthPurpose  `db:"purpose"`
	TokenHash string       `db:"token_hash"`
	ExpiresAt time.Time    `db:"expires_at"`
	UsedAt    sql.NullTime `db:"used_at"`
	CreatedAt time.Time    `db:"created_at"`
}
