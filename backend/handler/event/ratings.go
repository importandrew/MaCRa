package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func (h Handler) Ratings(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapViewRatings); err != nil {
		httperr.Error(c, err)
		return
	}

	ratings, err := h.event.Ratings(c.Request.Context(), eventID, userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, ratings)
}

// TourRating returns a tournament's rating by tour
// (GET /event/:id/tours/rating). It needs models.CapViewRatings.
func (h Handler) TourRating(c *gin.Context) {
	userID := reqctx.GetUserID(c)
	eventID := reqctx.GetEventID(c)

	if err := h.event.HasCapability(c.Request.Context(), eventID, userID, 0, models.CapViewRatings); err != nil {
		httperr.Error(c, err)
		return
	}

	rating, err := h.event.TourRating(c.Request.Context(), eventID, userID)
	if err != nil {
		httperr.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, rating)
}
