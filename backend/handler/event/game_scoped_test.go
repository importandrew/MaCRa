package event

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func newGameScopedContext(userID, eventID, gameID int) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{
		{Key: "eventId", Value: strconv.Itoa(eventID)},
		{Key: "gameId", Value: strconv.Itoa(gameID)},
	}
	c.Set(reqctx.UserIDKey, userID)
	return c, w
}

func TestHandler_GamePlayersAndLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const userID, eventID, gameID = 5, 1, 7

	routes := map[string]struct {
		handle    func(Handler, *gin.Context)
		expectGet func(*MockGameService)
	}{
		"players": {
			handle: Handler.GamePlayers,
			expectGet: func(m *MockGameService) {
				m.EXPECT().Players(mock.Anything, gameID).Return(models.GamePlayers{}, nil)
			},
		},
		"logs": {
			handle: Handler.GameLogs,
			expectGet: func(m *MockGameService) {
				m.EXPECT().Logs(mock.Anything, gameID).Return(models.GameLog{}, nil)
			},
		},
	}

	statuses := map[models.GameStatus]models.Capability{
		models.GameCreated: models.CapWatchLive,
		models.GameLive:    models.CapWatchLive,
		models.GameEnded:   models.CapWatchResults,
		models.GameAborted: models.CapWatchResults,
	}

	for name, route := range routes {
		for status, capability := range statuses {
			t.Run(name+" "+string(status)+" allowed", func(t *testing.T) {
				eventSvc := NewMockEventService(t)
				gameSvc := NewMockGameService(t)
				gameSvc.EXPECT().Status(mock.Anything, gameID).Return(status, nil)
				eventSvc.EXPECT().HasCapability(mock.Anything, eventID, userID, gameID, capability).Return(nil)
				route.expectGet(gameSvc)

				c, w := newGameScopedContext(userID, eventID, gameID)
				route.handle(NewHandler(eventSvc, gameSvc), c)

				assert.Equal(t, http.StatusOK, w.Code)
			})

			t.Run(name+" "+string(status)+" forbidden", func(t *testing.T) {
				eventSvc := NewMockEventService(t)
				gameSvc := NewMockGameService(t)
				gameSvc.EXPECT().Status(mock.Anything, gameID).Return(status, nil)
				eventSvc.EXPECT().HasCapability(mock.Anything, eventID, userID, gameID, capability).Return(models.ErrEventForbidden)

				c, w := newGameScopedContext(userID, eventID, gameID)
				route.handle(NewHandler(eventSvc, gameSvc), c)

				assert.Equal(t, http.StatusForbidden, w.Code)
			})
		}

		t.Run(name+" status error", func(t *testing.T) {
			gameSvc := NewMockGameService(t)
			gameSvc.EXPECT().Status(mock.Anything, gameID).Return(models.GameStatus(""), errors.New("database error"))

			c, w := newGameScopedContext(userID, eventID, gameID)
			route.handle(NewHandler(NewMockEventService(t), gameSvc), c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	}
}
