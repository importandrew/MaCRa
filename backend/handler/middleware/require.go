package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
)

func RequireUser(c *gin.Context) {
	if reqctx.GetUserID(c) <= reqctx.ZeroUserID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func RequireEvent(c *gin.Context) {
	if reqctx.GetEventID(c) <= 0 {
		httperr.ErrorWithCode(c, http.StatusBadRequest, errors.New("invalid event id"))
		return
	}
	c.Next()
}

func RequireGame(c *gin.Context) {
	if reqctx.GetGameID(c) <= 0 {
		httperr.ErrorWithCode(c, http.StatusBadRequest, errors.New("invalid game id"))
		return
	}
	c.Next()
}

func RequirePlayer(c *gin.Context) {
	if reqctx.GetPlayerID(c) <= 0 {
		httperr.ErrorWithCode(c, http.StatusBadRequest, errors.New("invalid player id"))
		return
	}
	c.Next()
}
