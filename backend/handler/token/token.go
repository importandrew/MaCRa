package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const expirationTime = 90 * 24 * time.Hour

var (
	ErrWrongClaims        = errors.New("token claims are of wrong type")
	ErrWrongSigningMethod = errors.New("invalid signing method")
)

type UserClaims struct {
	jwt.RegisteredClaims

	UserID int `json:"userID"`
}

func GenerateJWT(userID int, jwtKey []byte) (string, error) {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func IdentifyUser(token string, jwtKey []byte) (int, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrWrongSigningMethod
		}

		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := t.Claims.(*UserClaims)
	if !ok {
		return 0, ErrWrongClaims
	}

	return claims.UserID, nil
}
