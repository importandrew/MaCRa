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

func TestHandler_HUD(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		gameID       string
		queryParams  string
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok": {
			gameID:      "100",
			queryParams: "",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().HUD(
					mock.Anything,
					100,
					0,
				).Return(models.GameHUD{
					Live: &models.GameLiveHUD{
						Round: models.Round{
							Num:   1,
							Kanji: "東1",
						},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var hud models.GameHUD
				err := json.Unmarshal(w.Body.Bytes(), &hud)
				assert.NoError(t, err)
				assert.NotNil(t, hud.Live)
				assert.Equal(t, 1, hud.Live.Round.Num)
			},
		},
		"ok with linger": {
			gameID:      "100",
			queryParams: "?linger=30",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().HUD(
					mock.Anything,
					100,
					30,
				).Return(models.GameHUD{
					Final: []models.PlayerResult{
						{ID: 1, Name: "Alice", GP: 10, TP: 4},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var hud models.GameHUD
				err := json.Unmarshal(w.Body.Bytes(), &hud)
				assert.NoError(t, err)
				assert.Nil(t, hud.Live)
				assert.Len(t, hud.Final, 1)
			},
		},
		"invalid id": {
			gameID:       "not-a-number",
			queryParams:  "",
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
			gameID:      "200",
			queryParams: "",
			mockSetup: func(m *MockGameService) {
				m.EXPECT().HUD(
					mock.Anything,
					200,
					0,
				).Return(models.GameHUD{}, errors.New("game not found"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/game/"+testCase.gameID+"/hud"+testCase.queryParams, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: testCase.gameID}}

			handler.HUD(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
