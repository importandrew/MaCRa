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

func TestHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testJWT := []byte("test-jwt-key")

	testData := map[string]struct {
		requestBody  any
		domain       string
		mockSetup    func(*MockPlayerService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
		checkCookie  func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok with domain": {
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			domain: "example.com",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Login(
					mock.Anything,
					"test@example.com",
					"password123",
				).Return(models.Player{ID: 1, Name: "Test User"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, 1, response.ID)
				assert.Equal(t, "Test User", response.Name)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 1)
				assert.Equal(t, "token", cookies[0].Name)
				assert.NotEmpty(t, cookies[0].Value)
				assert.Equal(t, "example.com", cookies[0].Domain)
				assert.True(t, cookies[0].Secure)
				assert.True(t, cookies[0].HttpOnly)
			},
		},
		"ok without domain": {
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			domain: "",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Login(
					mock.Anything,
					"test@example.com",
					"password123",
				).Return(models.Player{ID: 2, Name: "Another User"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, 2, response.ID)
				assert.Equal(t, "Another User", response.Name)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 1)
				assert.Equal(t, "token", cookies[0].Name)
				assert.NotEmpty(t, cookies[0].Value)
				assert.Equal(t, models.DefaultDomain, cookies[0].Domain)
				assert.True(t, cookies[0].Secure)
				assert.True(t, cookies[0].HttpOnly)
			},
		},
		"invalid json": {
			requestBody:  "invalid json string",
			domain:       "",
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
			domain:       "",
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
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			domain: "",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Login(
					mock.Anything,
					"test@example.com",
					"wrongpassword",
				).Return(models.Player{}, errors.New("invalid credentials"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "invalid credentials", response["error"])
			},
		},
		"player not found": {
			requestBody: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			domain: "",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Login(
					mock.Anything,
					"nonexistent@example.com",
					"password123",
				).Return(models.Player{}, models.ErrPlayerNotFound)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, models.ErrPlayerNotFound.Error(), response["error"])
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockPlayerService := NewMockPlayerService(t)
			testCase.mockSetup(mockPlayerService)

			handler := NewHandler(mockPlayerService, testCase.domain, testJWT)

			var bodyBytes []byte
			var err error
			if testCase.requestBody != nil {
				bodyBytes, err = json.Marshal(testCase.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			handler.Login(ctx)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
			if testCase.checkCookie != nil {
				testCase.checkCookie(t, w)
			}
		})
	}
}
