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

func TestHandler_Players(t *testing.T) {
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
				m.EXPECT().Players(
					mock.Anything,
					100,
				).Return(models.GamePlayers{
					ID: 100,
					Players: []models.Player{
						{ID: 1, Name: "Player 1"},
						{ID: 2, Name: "Player 2"},
						{ID: 3, Name: "Player 3"},
						{ID: 4, Name: "Player 4"},
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var players models.GamePlayers
				err := json.Unmarshal(w.Body.Bytes(), &players)
				assert.NoError(t, err)
				assert.Equal(t, 100, players.ID)
				assert.Len(t, players.Players, 4)
				assert.Equal(t, 1, players.Players[0].ID)
				assert.Equal(t, "Player 1", players.Players[0].Name)
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
				m.EXPECT().Players(
					mock.Anything,
					200,
				).Return(models.GamePlayers{}, errors.New("game not found"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/game/"+testCase.gameID+"/players", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: testCase.gameID}}

			handler.Players(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
