package player

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/handler/authcookie"
)

func (h Handler) Logout(c *gin.Context) {
	authcookie.Clear(c)
	c.Status(http.StatusOK)
}
