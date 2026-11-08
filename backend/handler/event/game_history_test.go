package event

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func newGameHistoryContext(userID, eventID, playerID int) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{
		{Key: "eventId", Value: strconv.Itoa(eventID)},
		{Key: "playerId", Value: strconv.Itoa(playerID)},
	}
	c.Set(reqctx.UserIDKey, userID)
	return c, w
}

func TestHandler_GameHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ok", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().HasCapability(mock.Anything, 1, 5, 0, models.CapWatchResults).Return(nil)
		eventSvc.EXPECT().GameHistory(mock.Anything, 1).Return([]int{7, 8}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 0)
		handler.GameHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var body struct {
			GameIDs []int `json:"game_ids"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Equal(t, []int{7, 8}, body.GameIDs)
	})

	t.Run("no watch_results", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().HasCapability(mock.Anything, 1, 0, 0, models.CapWatchResults).Return(models.ErrEventForbidden)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(0, 1, 0)
		handler.GameHistory(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestHandler_PlayerGameHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ok", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().HasCapability(mock.Anything, 1, 5, 0, models.CapWatchResults).Return(nil)
		eventSvc.EXPECT().PlayerGameHistory(mock.Anything, 1, 3).Return([]int{9}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 3)
		handler.PlayerGameHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("no watch_results", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().HasCapability(mock.Anything, 1, 5, 0, models.CapWatchResults).Return(models.ErrEventForbidden)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 3)
		handler.PlayerGameHistory(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("own history needs no watch_results", func(t *testing.T) {
		// No HasCapability expectation: the mock fails the test if it is called.
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().PlayerGameHistory(mock.Anything, 1, 5).Return([]int{4}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newGameHistoryContext(5, 1, 5)
		handler.PlayerGameHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
