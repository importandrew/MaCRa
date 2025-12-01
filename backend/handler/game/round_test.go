package game

import (
	"bytes"
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

func TestHandler_EnterRoundResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		requestBody  any
		userID       int
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok deal-in": {
			requestBody: EnterRoundResultRequest{
				Outcome:  models.Outcome("deal-in"),
				Value:    8,
				WinnerID: 1,
				LoserID:  2,
			},
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().EnterRound(
					mock.Anything,
					1,
					models.Outcome("deal-in"),
					8,
					1,
					2,
				).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		"ok self": {
			requestBody: EnterRoundResultRequest{
				Outcome:  models.Outcome("self"),
				Value:    10,
				WinnerID: 1,
				LoserID:  0,
			},
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().EnterRound(
					mock.Anything,
					1,
					models.Outcome("self"),
					10,
					1,
					0,
				).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		"ok draw": {
			requestBody: EnterRoundResultRequest{
				Outcome:  models.Outcome("draw"),
				Value:    0,
				WinnerID: 0,
				LoserID:  0,
			},
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().EnterRound(
					mock.Anything,
					1,
					models.Outcome("draw"),
					0,
					0,
					0,
				).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		"invalid json": {
			requestBody:  "invalid json",
			userID:       1,
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"user id not in context": {
			requestBody: EnterRoundResultRequest{
				Outcome:  models.Outcome("deal-in"),
				Value:    8,
				WinnerID: 1,
				LoserID:  2,
			},
			userID:       0,
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"service error": {
			requestBody: EnterRoundResultRequest{
				Outcome:  models.Outcome("deal-in"),
				Value:    8,
				WinnerID: 1,
				LoserID:  2,
			},
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().EnterRound(
					mock.Anything,
					1,
					models.Outcome("deal-in"),
					8,
					1,
					2,
				).Return(errors.New("invalid hand"))
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

			req := httptest.NewRequest(http.MethodPost, "/api/round", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if testCase.userID > 0 {
				c.Set("userID", testCase.userID)
			}

			handler.EnterRoundResult(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
