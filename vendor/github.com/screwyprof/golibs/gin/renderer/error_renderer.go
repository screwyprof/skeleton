package renderer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError interface {
	Status() int
	Code() int
	Message() string
	Extra() map[string]interface{}
	Cause() error
}

type ErrorRenderer struct {
	Err error
}

func (r *ErrorRenderer) Render(c *gin.Context) {
	apiErr, ok := r.Err.(APIError)
	if ok {
		_ = c.AbortWithError(apiErr.Status(), r.Err)
		return
	}
	_ = c.AbortWithError(http.StatusInternalServerError, r.Err)
}
