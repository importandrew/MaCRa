package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func RequireUser(c *gin.Context) {
	if reqctx.GetUserID(c) <= reqctx.ZeroUserID {
		httperr.Error(c, models.ErrLoginRequired)
		return
	}
	c.Next()
}

func RequireEvent(c *gin.Context) {
	if reqctx.GetEventID(c) <= 0 {
		httperr.InvalidRequest(c, errors.New("invalid event id"))
		return
	}
	c.Next()
}

func RequireGame(c *gin.Context) {
	if reqctx.GetGameID(c) <= 0 {
		httperr.InvalidRequest(c, errors.New("invalid game id"))
		return
	}
	c.Next()
}

func RequirePlayer(c *gin.Context) {
	if reqctx.GetPlayerID(c) <= 0 {
		httperr.InvalidRequest(c, errors.New("invalid player id"))
		return
	}
	c.Next()
}
