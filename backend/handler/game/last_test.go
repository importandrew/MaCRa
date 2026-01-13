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

func TestHandler_Last(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		userID       int
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok": {
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().LastGameResult(
					mock.Anything,
					1,
				).Return(models.GameResult{
					ID: 100,
					PlayerResults: []models.PlayerResult{
						{ID: 1, Name: "Alice", GP: 10, TP: 4},
						{ID: 2, Name: "Bob", GP: -5, TP: 2},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result models.GameResult
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, 100, result.ID)
				assert.Len(t, result.PlayerResults, 2)
			},
		},
		"user id not in context": {
			userID:       0,
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusUnauthorized,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				// When AbortWithStatus is called, no JSON body is written
				assert.Empty(t, w.Body.String())
			},
		},
		"service error": {
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().LastGameResult(
					mock.Anything,
					1,
				).Return(models.GameResult{}, errors.New("game not found"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/last", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if testCase.userID > 0 {
				c.Set("userID", testCase.userID)
			}

			handler.Last(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
