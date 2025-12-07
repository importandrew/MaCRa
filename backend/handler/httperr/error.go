package httperr

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

type msg struct {
	Error string `json:"error"`
}

func ErrorWithCode(c *gin.Context, code int, err error) {
	logger.Error(err)
	c.AbortWithStatusJSON(code, msg{Error: err.Error()})
}

func Error(c *gin.Context, err error) {
	logger.Error(err)
	c.AbortWithStatusJSON(codeFromErr(err), msg{Error: err.Error()})
}

func codeFromErr(err error) int {
	for targetErr, code := range models.Errors {
		if errors.Is(err, targetErr) {
			return code
		}
	}
	return http.StatusInternalServerError
}
