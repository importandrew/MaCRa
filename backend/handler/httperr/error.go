package httperr

import (
	"errors"
	"fmt"
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/importandrew/MaCRa/logger"
	"github.com/importandrew/MaCRa/models"
)

// Error aborts the request with err's status and {"error": "<code>"} plus the
// data of a models.DataError in err (see models.Errors). An error not found
// there is a 500 "internal". The full error goes only to the log: at WARN for
// a client error (4xx), at ERROR otherwise.
func Error(c *gin.Context, err error) {
	apiErr := lookup(err)
	if apiErr.Status < http.StatusInternalServerError {
		logger.Warn(err)
	} else {
		logger.Error(err)
	}

	body := gin.H{}
	if dataErr, ok := errors.AsType[*models.DataError](err); ok {
		maps.Copy(body, dataErr.Data)
	}
	body["error"] = apiErr.Code

	c.AbortWithStatusJSON(apiErr.Status, body)
}

// InvalidRequest is Error for a request that can't be read, e.g. a binding
// error: a 400 "invalid_request".
func InvalidRequest(c *gin.Context, err error) {
	Error(c, fmt.Errorf("%w: %w", models.ErrInvalidRequest, err))
}

// lookup returns the APIError of the outermost error in err's tree that is in
// models.Errors, so a wrapping error wins over the one it wraps.
func lookup(err error) models.APIError {
	queue := []error{err}
	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		// Not models.Errors[e]: indexing a map with an error of an
		// incomparable type (e.g. validator.ValidationErrors, a slice)
		// panics. == between different types is just false.
		for target, apiErr := range models.Errors {
			if e == target { //nolint:errorlint // lookup walks the tree itself, one level at a time
				return apiErr
			}
		}

		switch u := e.(type) { //nolint:errorlint // the next level of the tree, not a match
		case interface{ Unwrap() error }:
			if next := u.Unwrap(); next != nil {
				queue = append(queue, next)
			}
		case interface{ Unwrap() []error }:
			queue = append(queue, u.Unwrap()...)
		}
	}

	return models.ErrInternal
}
