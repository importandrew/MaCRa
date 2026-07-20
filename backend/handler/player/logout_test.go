package player

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		domain       string
		expectedCode int
		checkCookie  func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok with domain": {
			domain:       "example.com",
			expectedCode: http.StatusOK,
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 2)

				auth := cookies[0]
				assert.Equal(t, "token", auth.Name)
				assert.Empty(t, auth.Value)
				assert.Equal(t, -1, auth.MaxAge)
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
		"ok without domain": {
			domain:       "",
			expectedCode: http.StatusOK,
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 2)

				auth := cookies[0]
				assert.Equal(t, "token", auth.Name)
				assert.Empty(t, auth.Value)
				assert.Equal(t, -1, auth.MaxAge)
				assert.Empty(t, auth.Domain)
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
		"ok with default domain": {
			domain:       models.DefaultDomain,
			expectedCode: http.StatusOK,
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 1)
				assert.Equal(t, "token", cookies[0].Name)
				assert.Empty(t, cookies[0].Value)
				assert.Equal(t, -1, cookies[0].MaxAge)
				assert.Equal(t, models.DefaultDomain, cookies[0].Domain)
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockPlayerService := NewMockPlayerService(t)

			handler := NewHandler(mockPlayerService, testCase.domain, nil)

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			handler.Logout(ctx)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.checkCookie != nil {
				testCase.checkCookie(t, w)
			}
		})
	}
}
