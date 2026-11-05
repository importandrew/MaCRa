package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Lifetime matches the auth cookie's MaxAge (400 days, the browser maximum), so a token
	// never expires while the cookie that carries it is still sent.
	Lifetime = 400 * 24 * time.Hour
	// RefreshAfter is how old a token gets before a request reissues it. Each visit then pushes
	// the expiry forward, so active users never have to log in again.
	RefreshAfter = 24 * time.Hour
)

var (
	ErrWrongClaims        = errors.New("token claims are of wrong type")
	ErrWrongSigningMethod = errors.New("invalid signing method")
	ErrNoExpiration       = errors.New("token has no expiration")
)

type UserClaims struct {
	jwt.RegisteredClaims

	UserID int `json:"userID"`
}

func GenerateJWT(userID int, jwtKey []byte) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(Lifetime)),
		UserID:    userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Identity is the result of a valid token.
type Identity struct {
	UserID int
	// NeedsRefresh reports that the token is older than RefreshAfter (judged by its expiry,
	// so tokens issued with a shorter lifetime count as old too) and should be reissued.
	NeedsRefresh bool
}

func IdentifyUser(token string, jwtKey []byte) (Identity, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrWrongSigningMethod
		}

		return jwtKey, nil
	})

	if err != nil {
		return Identity{}, err
	}

	claims, ok := t.Claims.(*UserClaims)
	if !ok {
		return Identity{}, ErrWrongClaims
	}

	if claims.ExpiresAt == nil {
		return Identity{}, ErrNoExpiration
	}

	remaining := time.Until(claims.ExpiresAt.Time)

	return Identity{
		UserID:       claims.UserID,
		NeedsRefresh: remaining < Lifetime-RefreshAfter,
	}, nil
}
