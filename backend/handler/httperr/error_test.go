package httperr

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/models"
)

// sliceErrors is an error of an incomparable type, like
// validator.ValidationErrors.
type sliceErrors []string

func (e sliceErrors) Error() string { return fmt.Sprint([]string(e)) }

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		err      error
		wantCode int
		wantBody string
	}{
		"sentinel": {
			err:      models.ErrGameNotFound,
			wantCode: http.StatusNotFound,
			wantBody: `{"error":"game_not_found"}`,
		},
		"wrapped sentinel, IDs stay in the log": {
			err:      fmt.Errorf("user 5: %w", models.ErrEventForbidden),
			wantCode: http.StatusForbidden,
			wantBody: `{"error":"event_forbidden"}`,
		},
		"data": {
			err: fmt.Errorf("player 3 in game 9: %w",
				models.WithData(models.ErrPlayerInAnotherGame, models.ErrorData{"name": "Ivan"})),
			wantCode: http.StatusConflict,
			wantBody: `{"error":"player_in_another_game","name":"Ivan"}`,
		},
		"data can't replace the code": {
			err:      models.WithData(models.ErrNoRounds, models.ErrorData{"error": "x"}),
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"no_rounds"}`,
		},
		"unknown error is internal, text not sent": {
			err:      errors.New("connection refused to 10.0.0.1"),
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"internal"}`,
		},
		"incomparable error in the tree": {
			err:      fmt.Errorf("%w: %w", models.ErrInvalidRequest, sliceErrors{"a", "b"}),
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"invalid_request"}`,
		},
		"outermost sentinel wins": {
			err:      fmt.Errorf("%w: %w", models.ErrValidatePlayers, models.ErrGameNotFound),
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"invalid_players"}`,
		},
		"joined errors": {
			err:      errors.Join(errors.New("x"), models.ErrTooLate),
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"too_late"}`,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, testCase.err)

			assert.Equal(t, testCase.wantCode, w.Code)
			assert.JSONEq(t, testCase.wantBody, w.Body.String())
			assert.True(t, c.IsAborted())
		})
	}
}

func TestInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InvalidRequest(c, sliceErrors{"Players: len"})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid_request"}`, w.Body.String())
}
