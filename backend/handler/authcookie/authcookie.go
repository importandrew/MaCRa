// Package authcookie writes the auth cookie that carries the user's JWT.
package authcookie

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/models"
)

const (
	// Name is the auth cookie name.
	Name = "token"
	// MaxAge is the auth cookie lifetime in seconds: 400 days, the longest Max-Age browsers accept.
	MaxAge = 400 * 24 * 60 * 60
)

// cookieAttrs returns the Set-Cookie Domain and Secure flag, derived from config and the request.
// Secure follows the actual connection (TLS or X-Forwarded-Proto), not a hardcoded true.
// Empty config → host-only cookie (not shared with sibling subdomains).
// If configured domain host does not match the request host, Domain is left empty (host-only)
// so LAN misconfigs like cookie_domain=localhost against mj.home still work.
func cookieAttrs(configured, requestHost, forwardedProto string, tls bool) (string, bool) {
	secure := tls || strings.EqualFold(forwardedProto, "https")
	if strings.HasPrefix(strings.ToLower(configured), "http://") {
		secure = false
	} else if strings.HasPrefix(strings.ToLower(configured), "https://") {
		secure = true
	}

	if configured == "" {
		return "", secure
	}

	host := stripCookieHost(configured)
	reqHost := stripCookieHost(requestHost)
	if host != "" && !strings.EqualFold(host, reqHost) {
		return "", secure
	}
	return host, secure
}

func stripCookieHost(raw string) string {
	host := raw
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	// Trailing slash or path from misconfigured URL values.
	if i := strings.IndexAny(host, "/"); i >= 0 {
		host = host[:i]
	}
	return host
}

// Set stores the token in the auth cookie for MaxAge.
func Set(c *gin.Context, configuredDomain, token string) {
	write(c, configuredDomain, token, MaxAge)
}

// Clear removes the auth cookie.
func Clear(c *gin.Context, configuredDomain string) {
	write(c, configuredDomain, "", -1)
}

// write sets or clears the auth cookie for this instance only.
// A legacy parent-domain cookie is also cleared so old shared tokens stop leaking.
func write(c *gin.Context, configuredDomain, value string, maxAge int) {
	cookieDomain, secure := cookieAttrs(
		configuredDomain,
		c.Request.Host,
		c.GetHeader("X-Forwarded-Proto"),
		c.Request.TLS != nil,
	)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(Name, value, maxAge, "/", cookieDomain, secure, true)

	if cookieDomain != models.DefaultDomain {
		// Match how the legacy shared cookie was issued (Secure + SameSite=None).
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(Name, "", -1, "/", models.DefaultDomain, true, true)
	}
}
