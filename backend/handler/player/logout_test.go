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
				assert.Len(t, cookies, 1)
				assert.Equal(t, "token", cookies[0].Name)
				assert.Empty(t, cookies[0].Value)
				assert.Equal(t, -1, cookies[0].MaxAge)
				assert.Equal(t, "example.com", cookies[0].Domain)
				assert.True(t, cookies[0].Secure)
				assert.True(t, cookies[0].HttpOnly)
			},
		},
		"ok without domain": {
			domain:       "",
			expectedCode: http.StatusOK,
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				assert.Len(t, cookies, 1)
				assert.Equal(t, "token", cookies[0].Name)
				assert.Empty(t, cookies[0].Value)
				assert.Equal(t, -1, cookies[0].MaxAge)
				assert.Equal(t, models.DefaultDomain, cookies[0].Domain)
				assert.True(t, cookies[0].Secure)
				assert.True(t, cookies[0].HttpOnly)
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
