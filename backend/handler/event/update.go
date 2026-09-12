package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

type UpdateEventRequest struct {
	Name         *string                 `json:"name"`
	Description  *string                 `json:"description"`
	Visibility   *models.EventVisibility `json:"visibility"`
	RatingsReach *models.ContentReach    `json:"ratings_reach"`
	GamesReach   *models.ContentReach    `json:"games_reach"`
}

func (h Handler) Update(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.ErrorWithCode(c, http.StatusBadRequest, err)
		return
	}

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapManageEvent); err != nil {
		httperr.Error(c, err)
		return
	}

	params := models.UpdateEventParams{
		ID:           eventID,
		Name:         req.Name,
		Description:  req.Description,
		Visibility:   req.Visibility,
		RatingsReach: req.RatingsReach,
		GamesReach:   req.GamesReach,
	}

	if err := h.event.Update(c.Request.Context(), params); err != nil {
		httperr.Error(c, err)
		return
	}

	c.Status(http.StatusOK)
}
