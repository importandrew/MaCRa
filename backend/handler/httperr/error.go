package httperr

import (
	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/logger"
)

var ErrReadingBody = "error reading body: %w"

type msg struct {
	Error string `json:"error"`
}

func ErrorResponse(c *gin.Context, code int, err error) {
	logger.Error(err)
	c.AbortWithStatusJSON(code, msg{Error: err.Error()})
}
