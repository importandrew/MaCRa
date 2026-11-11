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
	jwt []byte
}

func NewAuth(jwt []byte) Auth {
	return Auth{jwt: jwt}
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
			// The client only needs to know the session is gone; parser details stay in the log.
			logger.Warn(fmt.Errorf("identify user: %w", err))
			authcookie.Clear(c)
			httperr.Error(c, models.ErrSessionExpired)
			return
		}

		userID = identity.UserID

		// Parallel requests carrying the same old token each reissue it; the browser keeps the
		// last cookie. That happens once a day per user and costs a few HMAC signatures, so it
		// is accepted rather than coordinated.
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

	authcookie.Set(c, tkn)
}
