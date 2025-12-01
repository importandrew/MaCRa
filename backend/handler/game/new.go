package game

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/handler/httperr"
)

type NewGameRequest struct {
	Players []int `json:"players"`
}

func (h Handler) New(c *gin.Context) {
	var req NewGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	userID, err := authmw.GetUserIDFromContext(c)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if !slices.Contains(req.Players, userID) {
		httperr.ErrorResponse(c, http.StatusNotAcceptable, errors.New("cannot create game without yourself"))
		return
	}

	_, err = h.game.NewGame(c.Request.Context(), req.Players)
	if err != nil {
		httperr.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
