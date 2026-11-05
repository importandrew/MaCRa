package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/authcookie"
	"github.com/importandrew/MaCRa/handler/httperr"
	"github.com/importandrew/MaCRa/handler/reqctx"
	"github.com/importandrew/MaCRa/handler/token"
	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

type Auth struct {
	jwt          []byte
	cookieDomain string
}

func NewAuth(jwt []byte, cookieDomain string) Auth {
	return Auth{jwt: jwt, cookieDomain: cookieDomain}
}

// SetUserID identifies the caller by the auth cookie. Without a cookie the caller is anonymous.
// A cookie with an expired or otherwise invalid token is cleared and the request fails with 401,
// so the client knows its session is gone instead of silently getting anonymous responses.
// A valid token older than token.RefreshAfter is reissued, so active sessions never expire.
func (a Auth) SetUserID(c *gin.Context) {
	userID := reqctx.ZeroUserID

	cookie, err := c.Cookie(authcookie.Name)
	if err == nil && cookie != "" {
		identity, err := token.IdentifyUser(cookie, a.jwt)
		if err != nil {
			authcookie.Clear(c, a.cookieDomain)
			httperr.Error(c, fmt.Errorf("identify user: %w: %w", err, models.ErrSessionExpired))
			return
		}

		userID = identity.UserID

		if identity.NeedsRefresh {
			a.refresh(c, userID)
		}
	}

	c.Set(reqctx.UserIDKey, userID)
	c.Next()
}

// refresh reissues the token. A failure keeps the current, still valid cookie.
func (a Auth) refresh(c *gin.Context, userID int) {
	tkn, err := token.GenerateJWT(userID, a.jwt)
	if err != nil {
		logger.Error(fmt.Errorf("refresh token for user %d: %w", userID, err))
		return
	}

	authcookie.Set(c, a.cookieDomain, tkn)
}
