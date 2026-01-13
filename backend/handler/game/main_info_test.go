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

func TestHandler_MainInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		userID       int
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok both": {
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().LastGameResult(
					mock.Anything,
					1,
				).Return(models.GameResult{
					ID: 100,
					PlayerResults: []models.PlayerResult{
						{ID: 1, Name: "Alice", GP: 10, TP: 4},
					},
				}, nil)
				m.EXPECT().Now(
					mock.Anything,
					1,
				).Return(models.GameNow{
					ID:     200,
					Status: models.GameLive,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var info MainInfo
				err := json.Unmarshal(w.Body.Bytes(), &info)
				assert.NoError(t, err)
				assert.True(t, info.HasLast)
				assert.True(t, info.HasCurrent)
			},
		},
		"ok only last": {
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().LastGameResult(
					mock.Anything,
					1,
				).Return(models.GameResult{
					ID: 100,
					PlayerResults: []models.PlayerResult{
						{ID: 1, Name: "Alice", GP: 10, TP: 4},
					},
				}, nil)
				m.EXPECT().Now(
					mock.Anything,
					1,
				).Return(models.GameNow{}, errors.New("no active game"))
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var info MainInfo
				err := json.Unmarshal(w.Body.Bytes(), &info)
				assert.NoError(t, err)
				assert.True(t, info.HasLast)
				assert.False(t, info.HasCurrent)
			},
		},
		"ok only current": {
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().LastGameResult(
					mock.Anything,
					1,
				).Return(models.GameResult{}, errors.New("game not found"))
				m.EXPECT().Now(
					mock.Anything,
					1,
				).Return(models.GameNow{
					ID:     200,
					Status: models.GameLive,
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var info MainInfo
				err := json.Unmarshal(w.Body.Bytes(), &info)
				assert.NoError(t, err)
				assert.False(t, info.HasLast)
				assert.True(t, info.HasCurrent)
			},
		},
		"ok neither": {
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().LastGameResult(
					mock.Anything,
					1,
				).Return(models.GameResult{}, errors.New("game not found"))
				m.EXPECT().Now(
					mock.Anything,
					1,
				).Return(models.GameNow{}, errors.New("no active game"))
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var info MainInfo
				err := json.Unmarshal(w.Body.Bytes(), &info)
				assert.NoError(t, err)
				assert.False(t, info.HasLast)
				assert.False(t, info.HasCurrent)
			},
		},
		"user id not in context": {
			userID:       0,
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var info MainInfo
				err := json.Unmarshal(w.Body.Bytes(), &info)
				assert.NoError(t, err)
				assert.False(t, info.HasLast)
				assert.False(t, info.HasCurrent)
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockGameService := NewMockGameService(t)
			testCase.mockSetup(mockGameService)

			handler := NewHandler(mockGameService, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/main-info", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if testCase.userID > 0 {
				c.Set("userID", testCase.userID)
			}

			handler.MainInfo(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
