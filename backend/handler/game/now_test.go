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

func TestHandler_Now(t *testing.T) {
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
				m.EXPECT().Now(
					mock.Anything,
					1,
				).Return(models.GameNow{
					ID: 100,
					Round: models.Round{
						Num:   5,
						Kanji: "南1",
					},
					Me: models.PlayerWindScore{
						ID:   1,
						Name: "Player 1",
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var now models.GameNow
				err := json.Unmarshal(w.Body.Bytes(), &now)
				assert.NoError(t, err)
				assert.Equal(t, 100, now.ID)
				assert.Equal(t, 5, now.Round.Num)
				assert.Equal(t, 1, now.Me.ID)
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
				m.EXPECT().Now(
					mock.Anything,
					1,
				).Return(models.GameNow{}, errors.New("no active game"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/now", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if testCase.userID > 0 {
				c.Set("userID", testCase.userID)
			}

			handler.Now(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
