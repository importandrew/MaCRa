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

// AuthMailWindow limits the auth emails (verification, password reset) sent to
// one address: at most Max within the last Period.
type AuthMailWindow struct {
	Period time.Duration
	Max    int
}

// authMailsPerHour is "about five per hour" from the agreed abuse protection.
const authMailsPerHour = 5

// AuthMailLimits keep anyone from flooding an address through sign-up,
// resend-verification or forgot-password. Every auth email carries a fresh
// token, so the windows count the address's tokens of any purpose.
var AuthMailLimits = []AuthMailWindow{
	{Period: time.Minute, Max: 1},
	{Period: time.Hour, Max: authMailsPerHour},
}

type AuthToken struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	Purpose   AuthPurpose  `db:"purpose"`
	TokenHash string       `db:"token_hash"`
	ExpiresAt time.Time    `db:"expires_at"`
	UsedAt    sql.NullTime `db:"used_at"`
	CreatedAt time.Time    `db:"created_at"`
}
