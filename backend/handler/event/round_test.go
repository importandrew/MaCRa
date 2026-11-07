package event

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/models"
)

func newEntryTestContext(path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/event/1"+path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "eventId", Value: "1"}}
	c.Set(reqctx.UserIDKey, 1)
	return c, w
}

// TestHandler_RoundEntries_requireLogSize checks that an entry without the
// log size the table showed is refused before it reaches the game.
func TestHandler_RoundEntries_requireLogSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		path  string
		body  string
		enter func(Handler, *gin.Context)
	}{
		"round":         {"/round", `{"outcome":"draw"}`, Handler.EnterRoundResult},
		"penalty":       {"/penalty", `{"loser":2,"value":10}`, Handler.Penalty},
		"false mahjong": {"/false-mahjong", `{"loser":2,"value":30}`, Handler.FalseMahjong},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			handler := NewHandler(NewMockEventService(t), NewMockGameService(t))
			c, w := newEntryTestContext(testCase.path, testCase.body)

			testCase.enter(handler, c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestHandler_EnterRoundResult_passesLogSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameSvc := NewMockGameService(t)
	gameSvc.EXPECT().EnterRound(mock.Anything, 1, models.Draw, 0, 0, 0, 1, 0).
		Return(models.GameNow{}, models.ErrRoundAlreadyEntered)

	handler := NewHandler(NewMockEventService(t), gameSvc)
	c, w := newEntryTestContext("/round", `{"outcome":"draw","log_size":0}`)

	handler.EnterRoundResult(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}
