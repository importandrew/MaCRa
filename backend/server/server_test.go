package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/importandrew/MaCRa/handler/reqctx"
)

func TestClientIP(t *testing.T) {
	router := NewServer()
	router.GET("/ip", func(c *gin.Context) { c.String(http.StatusOK, c.ClientIP()) })

	testData := map[string]struct {
		remoteAddr   string
		forwardedFor string
		expectedIP   string
	}{
		"through nginx on the same machine": {
			remoteAddr:   "127.0.0.1:41000",
			forwardedFor: "203.0.113.7",
			expectedIP:   "203.0.113.7",
		},
		"through nginx over IPv6 loopback": {
			remoteAddr:   "[::1]:41000",
			forwardedFor: "203.0.113.7",
			expectedIP:   "203.0.113.7",
		},
		"direct client faking the header": {
			remoteAddr:   "198.51.100.2:5000",
			forwardedFor: "203.0.113.7",
			expectedIP:   "198.51.100.2",
		},
		"direct client": {
			remoteAddr: "198.51.100.2:5000",
			expectedIP: "198.51.100.2",
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ip", nil)
			req.RemoteAddr = testCase.remoteAddr
			if testCase.forwardedFor != "" {
				req.Header.Set("X-Forwarded-For", testCase.forwardedFor)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedIP, w.Body.String())
		})
	}
}

func TestLogFormatter(t *testing.T) {
	param := gin.LogFormatterParams{
		StatusCode: http.StatusOK,
		ClientIP:   "203.0.113.7",
		Method:     http.MethodGet,
		Path:       "/api/profile",
	}

	param.Keys = map[any]any{reqctx.UserIDKey: 36}
	assert.Contains(t, logFormatter(param), `203.0.113.7 | user 36    | GET     "/api/profile"`)

	param.Keys = map[any]any{reqctx.UserIDKey: reqctx.ZeroUserID}
	assert.Contains(t, logFormatter(param), "| user -     |", "a visitor who isn't signed in")

	param.Keys = nil
	assert.Contains(t, logFormatter(param), "| user -     |", "a route without SetUserID")
}
