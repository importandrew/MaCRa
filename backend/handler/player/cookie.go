package player

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authmw "github.com/importandrew/MaCRa/handler/middleware"
	"github.com/importandrew/MaCRa/models"
)

// cookieAttrs derives Set-Cookie Domain / Secure from config and the request.
// Secure follows the actual connection (TLS or X-Forwarded-Proto), not a hardcoded true.
// Empty config → host-only cookie (not shared with sibling subdomains).
// If configured domain host does not match the request host, Domain is left empty (host-only)
// so LAN misconfigs like cookie_domain=localhost against mj.home still work.
func cookieAttrs(configured, requestHost, forwardedProto string, tls bool) (domain string, secure bool) {
	secure = tls || strings.EqualFold(forwardedProto, "https")
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

// writeAuthCookie sets or clears the auth cookie for this instance only.
// A legacy parent-domain cookie is also cleared so old shared tokens stop leaking.
func (h Handler) writeAuthCookie(c *gin.Context, value string, maxAge int) {
	cookieDomain, secure := cookieAttrs(
		h.domain,
		c.Request.Host,
		c.GetHeader("X-Forwarded-Proto"),
		c.Request.TLS != nil,
	)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(authmw.TokenKey, value, maxAge, "/", cookieDomain, secure, true)

	if cookieDomain != models.DefaultDomain {
		// Match how the legacy shared cookie was issued (Secure + SameSite=None).
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(authmw.TokenKey, "", -1, "/", models.DefaultDomain, true, true)
	}
}
