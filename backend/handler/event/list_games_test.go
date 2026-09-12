package event

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newListGamesContext(userID, eventID int, query string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/api/event/"+strconv.Itoa(eventID)+"/games?"+query,
		nil,
	)
	c.Params = gin.Params{{Key: "eventId", Value: strconv.Itoa(eventID)}}
	c.Set(reqctx.UserIDKey, userID)
	return c, w
}

func TestHandler_ListGames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("live ok", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().
			HasCapability(mock.Anything, 1, 5, 0, models.CapWatchLive).
			Return(nil)
		eventSvc.EXPECT().
			ListGamesByStatus(mock.Anything, 1, models.GameLive, 50, 0).
			Return(models.EventGames{
				Games: []models.EventGameItem{
					{ID: 9, Status: models.GameLive},
				},
			}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newListGamesContext(5, 1, "status=live")
		handler.ListGames(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var body models.EventGames
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Len(t, body.Games, 1)
		assert.Equal(t, 9, body.Games[0].ID)
	})

	t.Run("ended uses watch_results", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().
			HasCapability(mock.Anything, 1, 5, 0, models.CapWatchResults).
			Return(nil)
		eventSvc.EXPECT().
			ListGamesByStatus(mock.Anything, 1, models.GameEnded, 50, 0).
			Return(models.EventGames{}, nil)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newListGamesContext(5, 1, "status=ended")
		handler.ListGames(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid status", func(t *testing.T) {
		handler := NewHandler(NewMockEventService(t), NewMockGameService(t))
		c, w := newListGamesContext(5, 1, "status=aborted")
		handler.ListGames(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().
			HasCapability(mock.Anything, 1, 5, 0, models.CapWatchLive).
			Return(models.ErrEventForbidden)

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newListGamesContext(5, 1, "status=live")
		handler.ListGames(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		eventSvc := NewMockEventService(t)
		eventSvc.EXPECT().
			HasCapability(mock.Anything, 1, 5, 0, models.CapWatchLive).
			Return(nil)
		eventSvc.EXPECT().
			ListGamesByStatus(mock.Anything, 1, models.GameLive, 50, 0).
			Return(models.EventGames{}, errors.New("db down"))

		handler := NewHandler(eventSvc, NewMockGameService(t))
		c, w := newListGamesContext(5, 1, "status=live")
		handler.ListGames(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
