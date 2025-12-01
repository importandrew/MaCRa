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

func TestHandler_New(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		requestBody  any
		userID       int
		mockSetup    func(*MockGameService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok": {
			requestBody: NewGameRequest{
				Players: []int{1, 2, 3, 4},
			},
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().NewGame(
					mock.Anything,
					[]int{1, 2, 3, 4},
				).Return(100, nil)
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
			requestBody: NewGameRequest{
				Players: []int{1, 2, 3, 4},
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
		"user not in players list": {
			requestBody: NewGameRequest{
				Players: []int{2, 3, 4, 5},
			},
			userID:       1,
			mockSetup:    func(m *MockGameService) {},
			expectedCode: http.StatusNotAcceptable,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Contains(t, response["error"], "cannot create game without yourself")
			},
		},
		"service error": {
			requestBody: NewGameRequest{
				Players: []int{1, 2, 3, 4},
			},
			userID: 1,
			mockSetup: func(m *MockGameService) {
				m.EXPECT().NewGame(
					mock.Anything,
					[]int{1, 2, 3, 4},
				).Return(0, errors.New("invalid players"))
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

			req := httptest.NewRequest(http.MethodPost, "/api/new", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if testCase.userID > 0 {
				c.Set("userID", testCase.userID)
			}

			handler.New(c)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
