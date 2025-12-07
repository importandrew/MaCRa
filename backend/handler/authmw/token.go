package authmw

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/importandrew/MaCRa/models"
)

const expirationTime = 90 * 24 * time.Hour

type UserClaims struct {
	jwt.RegisteredClaims

	UserID int `json:"userID"`
}

func GenerateJWTToken(userID int, jwtKey []byte) (string, error) {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func identifyUserByToken(token string, jwtKey []byte) (int, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrWrongSigningMethod
		}

		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := t.Claims.(*UserClaims)
	if !ok {
		return 0, models.ErrWrongClaims
	}

	return claims.UserID, nil
}
