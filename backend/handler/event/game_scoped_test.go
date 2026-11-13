package event

import (
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

	// Which of watch_live / watch_results applies is the service's business
	// now (view_game); the handler asks for view_game.
	for name, route := range routes {
		t.Run(name+" allowed", func(t *testing.T) {
			eventSvc := NewMockEventService(t)
			gameSvc := NewMockGameService(t)
			eventSvc.EXPECT().HasCapability(mock.Anything, eventID, userID, gameID, models.CapViewGame).Return(nil)
			route.expectGet(gameSvc)

			c, w := newGameScopedContext(userID, eventID, gameID)
			route.handle(NewHandler(eventSvc, gameSvc), c)

			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run(name+" forbidden", func(t *testing.T) {
			eventSvc := NewMockEventService(t)
			eventSvc.EXPECT().HasCapability(mock.Anything, eventID, userID, gameID, models.CapViewGame).Return(models.ErrEventForbidden)

			c, w := newGameScopedContext(userID, eventID, gameID)
			route.handle(NewHandler(eventSvc, NewMockGameService(t)), c)

			assert.Equal(t, http.StatusForbidden, w.Code)
		})
	}
}
