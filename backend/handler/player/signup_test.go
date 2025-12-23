package player

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

func TestHandler_SignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		requestBody  any
		mockSetup    func(*MockPlayerService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok": {
			requestBody: SignupRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Add(
					mock.Anything,
					"test@example.com",
					"Test User",
					"password123",
				).Return(models.Player{ID: 1, Name: "Test User"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		"invalid json": {
			requestBody:  "invalid json string",
			mockSetup:    func(m *MockPlayerService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"empty body": {
			requestBody:  nil,
			mockSetup:    func(m *MockPlayerService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"service error": {
			requestBody: SignupRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Add(
					mock.Anything,
					"test@example.com",
					"Test User",
					"password123",
				).Return(models.Player{}, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "database error", response["error"])
			},
		},
		"player already exists": {
			requestBody: SignupRequest{
				Email:    "existing@example.com",
				Name:     "Existing User",
				Password: "password123",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Add(
					mock.Anything,
					"existing@example.com",
					"Existing User",
					"password123",
				).Return(models.Player{}, models.ErrPlayerExists)
			},
			expectedCode: http.StatusConflict,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, models.ErrPlayerExists.Error(), response["error"])
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockPlayerService := NewMockPlayerService(t)
			testCase.mockSetup(mockPlayerService)

			handler := NewHandler(mockPlayerService, "", nil)

			var bodyBytes []byte
			var err error
			if testCase.requestBody != nil {
				bodyBytes, err = json.Marshal(testCase.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			handler.SignUp(ctx)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
