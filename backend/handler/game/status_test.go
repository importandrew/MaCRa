package game

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Status(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		gameID       string
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok": {
			gameID: "100",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().Status(
					mock.Anything,
					100,
				).Return(models.GameLive, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, float64(100), response["id"])
				assert.Equal(t, string(models.GameLive), response["status"])
			},
		},
		"ok ended": {
			gameID: "200",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().Status(
					mock.Anything,
					200,
				).Return(models.GameEnded, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, float64(200), response["id"])
				assert.Equal(t, string(models.GameEnded), response["status"])
			},
		},
		"invalid id": {
			gameID:       "not-a-number",
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"service error": {
			gameID: "300",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().Status(
					mock.Anything,
					300,
				).Return(models.GameStatus(""), errors.New("game not found"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameService := NewMockGameService(t)
			testCase.mockSetup(mockGameService)

			handler := NewHandler(mockGameService, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/game/"+testCase.gameID+"/status", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: testCase.gameID}}

			handler.Status(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
