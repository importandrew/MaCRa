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

var (
	ErrNotFound    = errors.New("user id not found in context")
	ErrInvalidType = errors.New("user id in context is of invalid type")
	ErrEmptyToken  = errors.New("token is empty")
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
		httperr.ErrorResponse(c, http.StatusUnauthorized, ErrEmptyToken)
		return
	}

	userID, err := identifyUserByToken(cookie, a.jwt)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtxKey, userID)
}

func GetUserIDFromContext(c *gin.Context) (int, error) {
	value, ok := c.Get(userCtxKey)
	if !ok {
		return 0, ErrNotFound
	}

	id, ok := value.(int)
	if !ok {
		return 0, ErrInvalidType
	}

	return id, nil
}
