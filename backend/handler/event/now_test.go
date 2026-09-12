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
)

func newNowTestContext(userID, eventID int) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/event/"+strconv.Itoa(eventID)+"/now", nil)
	c.Params = gin.Params{{Key: "eventId", Value: strconv.Itoa(eventID)}}
	c.Set(reqctx.UserIDKey, userID)
	return c, w
}

func TestHandler_Now(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameSvc := NewMockGameService(t)
	gameSvc.EXPECT().Now(mock.Anything, 1, 1).Return(models.GameNow{ID: 100}, nil)

	handler := NewHandler(NewMockEventService(t), gameSvc)
	c, w := newNowTestContext(1, 1)

	handler.Now(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var now models.GameNow
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &now))
	assert.Equal(t, 100, now.ID)
}

func TestHandler_Now_serviceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gameSvc := NewMockGameService(t)
	gameSvc.EXPECT().Now(mock.Anything, 1, 1).Return(models.GameNow{}, errors.New("no active game"))

	handler := NewHandler(NewMockEventService(t), gameSvc)
	c, w := newNowTestContext(1, 1)

	handler.Now(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
