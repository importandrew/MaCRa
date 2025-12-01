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

func TestHandler_Logs(t *testing.T) {
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
				m.EXPECT().Logs(
					mock.Anything,
					100,
				).Return(models.GameLog{
					ID: 100,
					RoundLogs: []models.RoundLog{
						{Round: 1, Outcome: models.Outcome("draw")},
						{Round: 2, Outcome: models.Outcome("deal-in")},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var log models.GameLog
				err := json.Unmarshal(w.Body.Bytes(), &log)
				assert.NoError(t, err)
				assert.Equal(t, 100, log.ID)
				assert.Len(t, log.RoundLogs, 2)
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
			gameID: "200",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().Logs(
					mock.Anything,
					200,
				).Return(models.GameLog{}, errors.New("game not found"))
			},
			expectedCode: http.StatusBadRequest,
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

			req := httptest.NewRequest(http.MethodGet, "/api/game/"+testCase.gameID+"/logs", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: testCase.gameID}}

			handler.Logs(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
