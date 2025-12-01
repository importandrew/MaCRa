package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_FalseMahjong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		gameID       string
		requestBody  any
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok": {
			gameID: "100",
			requestBody: FalseMahjongRequest{
				LoserID:        2,
				Value:          30,
				ForceNextRound: false,
			},
			mockSetup: func(m *MockGameService) {
				m.EXPECT().FalseMahjong(
					mock.Anything,
					100,
					2,
					30,
					false,
				).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		"ok with force next round": {
			gameID: "100",
			requestBody: FalseMahjongRequest{
				LoserID:        3,
				Value:          60,
				ForceNextRound: true,
			},
			mockSetup: func(m *MockGameService) {
				m.EXPECT().FalseMahjong(
					mock.Anything,
					100,
					3,
					60,
					true,
				).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		"invalid game id": {
			gameID:       "not-a-number",
			requestBody:  FalseMahjongRequest{LoserID: 2, Value: 30},
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"invalid json": {
			gameID:       "100",
			requestBody:  "invalid json",
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
			gameID: "100",
			requestBody: FalseMahjongRequest{
				LoserID:        2,
				Value:          30,
				ForceNextRound: false,
			},
			mockSetup: func(m *MockGameService) {
				m.EXPECT().FalseMahjong(
					mock.Anything,
					100,
					2,
					30,
					false,
				).Return(errors.New("game not found"))
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

			var bodyBytes []byte
			var err error
			if testCase.requestBody != nil {
				bodyBytes, err = json.Marshal(testCase.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/game/"+testCase.gameID+"/false-mahjong",
				bytes.NewBuffer(bodyBytes),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: testCase.gameID}}

			handler.FalseMahjong(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
