package player

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/config"
)

func TestHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		viaNginx       bool // over https, through nginx on the same machine
		expectedSecure bool
	}{
		"over https through nginx": {viaNginx: true, expectedSecure: true},
		"over plain http":          {viaNginx: false, expectedSecure: false},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			handler := NewHandler(NewMockPlayerService(t), config.Config{})

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			req.Host = "macra.pro"
			if testCase.viaNginx {
				req.RemoteAddr = "127.0.0.1:40000"
				req.Header.Set("X-Forwarded-Proto", "https")
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			handler.Logout(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			cookies := w.Result().Cookies()
			assert.Len(t, cookies, 2)

			auth := cookies[0]
			assert.Equal(t, "token", auth.Name)
			assert.Empty(t, auth.Value)
			assert.Equal(t, -1, auth.MaxAge)
			assert.Empty(t, auth.Domain, "host-only")
			assert.Equal(t, testCase.expectedSecure, auth.Secure)
			assert.True(t, auth.HttpOnly)
			assert.Equal(t, http.SameSiteLaxMode, auth.SameSite)

			legacy := cookies[1]
			assert.Equal(t, "token", legacy.Name)
			assert.Equal(t, -1, legacy.MaxAge)
			assert.Equal(t, "macra.pro", legacy.Domain)
		})
	}
}
