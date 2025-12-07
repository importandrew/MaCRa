package authmw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/models"
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
		httperr.ErrorWithCode(c, http.StatusUnauthorized, models.ErrEmptyToken)
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
		return 0, models.ErrUserIDNotFound
	}

	id, ok := value.(int)
	if !ok {
		return 0, models.ErrUserIDInvalidType
	}

	return id, nil
}
