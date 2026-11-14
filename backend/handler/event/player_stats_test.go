package event

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestHandler_PlayerStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("own statistics need no rating access", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().PlayerStats(mock.Anything, 1, 5).Return(models.PlayerStats{PlayerID: 5, Games: 3}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 5)
		handler.PlayerStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var body models.PlayerStats
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Equal(t, 3, body.Games)
	})

	t.Run("another player's statistics follow the rating", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().HasCapability(mock.Anything, 1, 5, 0, models.CapViewRatings).Return(nil)
		eventSvc.EXPECT().PlayerStats(mock.Anything, 1, 3).Return(models.PlayerStats{PlayerID: 3}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 3)
		handler.PlayerStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("without rating access", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().HasCapability(mock.Anything, 1, 5, 0, models.CapViewRatings).Return(models.ErrEventForbidden)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 3)
		handler.PlayerStats(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
