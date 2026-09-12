package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/handler/token"
)

const (
	TokenKey = "token"
)

type Auth struct {
	jwt []byte
}

func NewAuth(jwt []byte) Auth {
	return Auth{jwt: jwt}
}

func (a Auth) SetUserID(c *gin.Context) {
	userID := reqctx.ZeroUserID

	cookie, err := c.Cookie(TokenKey)
	if err == nil {
		id, err := token.IdentifyUser(cookie, a.jwt)
		if err == nil {
			userID = id
		}
	}

	c.Set(reqctx.UserIDKey, userID)
	c.Next()
}
