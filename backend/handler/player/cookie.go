package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/handler/authmw"
	"github.com/importandrew/MaCRa/models"
)

// writeAuthCookie sets or clears the auth cookie for this instance only.
// An empty h.domain produces a host-only cookie (not shared with sibling subdomains).
// A legacy parent-domain cookie is also cleared so old shared tokens stop leaking.
func (h Handler) writeAuthCookie(c *gin.Context, value string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(authmw.TokenKey, value, maxAge, "/", h.domain, true, true)

	if h.domain != models.DefaultDomain {
		c.SetCookie(authmw.TokenKey, "", -1, "/", models.DefaultDomain, true, true)
	}
}
