package authcookie

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		remoteAddr     string
		forwardedProto string
		tls            bool
		expectedSecure bool
	}{
		"through nginx over https": {
			remoteAddr: "127.0.0.1:40000", forwardedProto: "https", expectedSecure: true,
		},
		"through nginx over IPv6 loopback": {
			remoteAddr: "[::1]:40000", forwardedProto: "https", expectedSecure: true,
		},
		"through nginx over http": {
			remoteAddr: "127.0.0.1:40000", forwardedProto: "http", expectedSecure: false,
		},
		"direct TLS": {
			remoteAddr: "198.51.100.2:5000", tls: true, expectedSecure: true,
		},
		"direct http faking the header": {
			remoteAddr: "198.51.100.2:5000", forwardedProto: "https", expectedSecure: false,
		},
		"local development over http": {
			remoteAddr: "127.0.0.1:40000", expectedSecure: false,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", nil)
			c.Request.Host = "macra.pro"
			c.Request.RemoteAddr = testCase.remoteAddr
			if testCase.forwardedProto != "" {
				c.Request.Header.Set("X-Forwarded-Proto", testCase.forwardedProto)
			}
			if testCase.tls {
				c.Request.TLS = &tls.ConnectionState{}
			}

			Set(c, "jwt")

			cookies := w.Result().Cookies()
			require.Len(t, cookies, 2)

			auth := cookies[0]
			assert.Equal(t, Name, auth.Name)
			assert.Equal(t, "jwt", auth.Value)
			assert.Equal(t, MaxAge, auth.MaxAge)
			assert.Empty(t, auth.Domain, "host-only: subdomains don't share it")
			assert.Equal(t, testCase.expectedSecure, auth.Secure)
			assert.True(t, auth.HttpOnly)
			assert.Equal(t, http.SameSiteLaxMode, auth.SameSite)

			legacy := cookies[1]
			assert.Equal(t, Name, legacy.Name)
			assert.Equal(t, -1, legacy.MaxAge)
			assert.Equal(t, legacyDomain, legacy.Domain)
		})
	}
}

func TestClear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/logout", nil)

	Clear(c)

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 2)
	assert.Empty(t, cookies[0].Value)
	assert.Equal(t, -1, cookies[0].MaxAge)
	assert.Empty(t, cookies[0].Domain)
}
