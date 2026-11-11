// Package authcookie writes the auth cookie that carries the user's JWT.
package authcookie

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/token"
)

const (
	// Name is the auth cookie name.
	Name = "token"
	// MaxAge is the auth cookie lifetime in seconds. It is derived from token.Lifetime, so the
	// cookie and the token it carries always expire together.
	MaxAge = int(token.Lifetime / time.Second)

	// legacyDomain is the parent domain an old version issued a shared cookie for. Every write
	// clears that cookie too, so an old token can't shadow the current one on any subdomain.
	legacyDomain = "macra.pro"
)

// Set stores the token in the auth cookie for MaxAge.
func Set(c *gin.Context, token string) {
	write(c, token, MaxAge)
}

// Clear removes the auth cookie.
func Clear(c *gin.Context) {
	write(c, "", -1)
}

// write sets or clears the auth cookie. The cookie is host-only (no Domain attribute): it
// belongs to the site's own address, so macra.pro and dev.macra.pro never see each other's
// sessions. It is Secure when the client's connection is https.
func write(c *gin.Context, value string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(Name, value, maxAge, "/", "", isHTTPS(c), true)

	// Match how the legacy shared cookie was issued (Secure + SameSite=None).
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(Name, "", -1, "/", legacyDomain, true, true)
}

// isHTTPS reports whether the client reached the site over https: directly over TLS, or
// through nginx on the same machine, which says so in X-Forwarded-Proto. The header is
// ignored from anyone else.
func isHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}

	ip := net.ParseIP(c.RemoteIP())
	return ip != nil && ip.IsLoopback() && strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
