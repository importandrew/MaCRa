package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/handler/event"
	"github.com/importandrew/MaCRa/handler/player"
)

// TestRegisterRoutes_AnonymousNeedsLogin checks that an anonymous request to a route that
// changes data, or to a route about the caller's own game, gets 401 before any handler runs.
func TestRegisterRoutes_AnonymousNeedsLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := New(event.Handler{}, player.Handler{}, []byte("test-key"))
	h.RegisterRoutes(router)

	params := strings.NewReplacer(":eventId", "1", ":gameId", "2", ":playerId", "3")

	var checked int
	for _, route := range router.Routes() {
		if !strings.HasPrefix(route.Path, "/api/") || route.Method == http.MethodGet {
			continue
		}
		checked++
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(route.Method, params.Replace(route.Path), nil))
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
	require.NotZero(t, checked, "no mutation routes found")

	for _, path := range []string{"/api/players", "/api/profile", "/api/event/1/now", "/api/event/1/last"} {
		t.Run("GET "+path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}
