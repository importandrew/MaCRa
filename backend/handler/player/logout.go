package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) Logout(c *gin.Context) {
	h.writeAuthCookie(c, "", -1)
	c.Status(http.StatusOK)
}
