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
				assert.Len(t, cookies, 2)

				auth := cookies[0]
				assert.Equal(t, "token", auth.Name)
				assert.NotEmpty(t, auth.Value)
				assert.Equal(t, "example.com", auth.Domain)
				assert.True(t, auth.Secure)
				assert.True(t, auth.HttpOnly)
				assert.Equal(t, http.SameSiteLaxMode, auth.SameSite)

				legacy := cookies[1]
				assert.Equal(t, "token", legacy.Name)
				assert.Empty(t, legacy.Value)
				assert.Equal(t, -1, legacy.MaxAge)
				assert.Equal(t, models.DefaultDomain, legacy.Domain)
			},
		},
		"ok without domain http": {
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
				assert.Len(t, cookies, 2)

				auth := cookies[0]
				assert.Equal(t, "token", auth.Name)
				assert.NotEmpty(t, auth.Value)
				assert.Equal(t, "", auth.Domain)
				assert.False(t, auth.Secure)
				assert.True(t, auth.HttpOnly)
				assert.Equal(t, http.SameSiteLaxMode, auth.SameSite)

				legacy := cookies[1]
				assert.Equal(t, "token", legacy.Name)
				assert.Empty(t, legacy.Value)
				assert.Equal(t, models.DefaultDomain, legacy.Domain)
			},
		},
		"ok http lan mismatched cookie domain": {
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			domain: "localhost",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Login(
					mock.Anything,
					"test@example.com",
					"password123",
				).Return(models.Player{ID: 3, Name: "LAN User"}, nil)
			},
			expectedCode: http.StatusOK,
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 2)
				assert.Equal(t, "", cookies[0].Domain)
				assert.False(t, cookies[0].Secure)
				assert.Equal(t, models.DefaultDomain, cookies[1].Domain)
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
			expectedCode: http.StatusNotFound,
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
			switch testCase.domain {
			case "example.com":
				req.Host = "example.com"
				req.Header.Set("X-Forwarded-Proto", "https")
			case "localhost":
				req.Host = "mj.home:3001"
			default:
				req.Host = "mj.home:3001"
			}

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
