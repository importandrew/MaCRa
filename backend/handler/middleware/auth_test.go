package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/handler/authcookie"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/handler/token"
)

var testKey = []byte("test-key")

func tokenExpiringIn(t *testing.T, d time.Duration) string {
	t.Helper()

	claims := &token.UserClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
		UserID:    7,
	}
	tkn, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(testKey)
	require.NoError(t, err)

	return tkn
}

func authCookies(w *httptest.ResponseRecorder) []*http.Cookie {
	var cookies []*http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == authcookie.Name && c.Domain == "" {
			cookies = append(cookies, c)
		}
	}
	return cookies
}

func TestAuth_SetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fresh, err := token.GenerateJWT(7, testKey)
	require.NoError(t, err)

	tests := map[string]struct {
		cookie      string
		wantCode    int
		wantUserID  int
		checkCookie func(*testing.T, []*http.Cookie)
	}{
		"no cookie is anonymous": {
			wantCode:   http.StatusOK,
			wantUserID: reqctx.ZeroUserID,
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				assert.Empty(t, cookies)
			},
		},
		"fresh token is kept": {
			cookie:     fresh,
			wantCode:   http.StatusOK,
			wantUserID: 7,
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				assert.Empty(t, cookies)
			},
		},
		"old token is reissued": {
			cookie:     tokenExpiringIn(t, 30*24*time.Hour),
			wantCode:   http.StatusOK,
			wantUserID: 7,
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				require.Len(t, cookies, 1)
				assert.Equal(t, authcookie.MaxAge, cookies[0].MaxAge)

				identity, err := token.IdentifyUser(cookies[0].Value, testKey)
				require.NoError(t, err)
				assert.Equal(t, 7, identity.UserID)
				assert.False(t, identity.NeedsRefresh)
			},
		},
		"expired token is cleared with 401": {
			cookie:   tokenExpiringIn(t, -time.Minute),
			wantCode: http.StatusUnauthorized,
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				require.Len(t, cookies, 1)
				assert.Empty(t, cookies[0].Value)
				assert.Equal(t, -1, cookies[0].MaxAge)
			},
		},
		"invalid token is cleared with 401": {
			cookie:   "garbage",
			wantCode: http.StatusUnauthorized,
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				require.Len(t, cookies, 1)
				assert.Equal(t, -1, cookies[0].MaxAge)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotUserID := -1

			router := gin.New()
			router.GET("/api", NewAuth(testKey).SetUserID, func(c *gin.Context) {
				gotUserID = reqctx.GetUserID(c)
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api", nil)
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: authcookie.Name, Value: tc.cookie})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
			if tc.wantCode == http.StatusOK {
				assert.Equal(t, tc.wantUserID, gotUserID)
			} else {
				assert.Equal(t, -1, gotUserID, "handler must not run")
				assert.JSONEq(t, `{"error":"session_expired"}`, w.Body.String(),
					"the body must not expose token parser details")
			}
			tc.checkCookie(t, authCookies(w))
		})
	}
}
