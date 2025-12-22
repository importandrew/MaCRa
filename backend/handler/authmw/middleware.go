package authmw

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
)

const (
	TokenKey   = "token"
	userCtxKey = "userID"
)

type Auth struct {
	jwt []byte
}

func NewAuth(jwt []byte) Auth {
	return Auth{jwt: jwt}
}

func (a Auth) IdentifyUser(c *gin.Context) {
	cookie, err := c.Cookie(TokenKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		httperr.Error(c, err)
		return
	}

	userID, err := identifyUserByToken(cookie, a.jwt)
	if err != nil {
		httperr.ErrorWithCode(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtxKey, userID)
}

func GetUserIDFromContext(c *gin.Context) (int, error) {
	value, ok := c.Get(userCtxKey)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	id, ok := value.(int)
	if !ok {
		return 0, errors.New("userID should be int")
	}

	return id, nil
}
