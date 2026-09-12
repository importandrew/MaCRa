package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type CreateEventRequest struct {
	Type         models.EventType       `json:"type" binding:"required"`
	Name         string                 `json:"name" binding:"required"`
	Description  string                 `json:"description"`
	Visibility   models.EventVisibility `json:"visibility"`
	RatingsReach *models.ContentReach   `json:"ratings_reach"`
	GamesReach   *models.ContentReach   `json:"games_reach"`
}

func (h Handler) Create(c *gin.Context) {
	callerID := reqctx.GetUserID(c)

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	params := models.EventCreateParams{
		OwnerID:      callerID,
		Type:         req.Type,
		Name:         req.Name,
		Description:  req.Description,
		Visibility:   req.Visibility,
		RatingsReach: req.RatingsReach,
		GamesReach:   req.GamesReach,
	}

	id, err := h.event.Create(c.Request.Context(), params)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
