package ctxtags

import (
	"context"

	"github.com/gin-gonic/gin"
)

var ctxKey = "ctxtags"

func FromContext(c context.Context) Tags {
	t, ok := c.Value(ctxKey).(Tags)
	if !ok {
		return NoopTags
	}

	return t
}

func ToContext(c *gin.Context, tags Tags) *gin.Context {
	c.Set(ctxKey, tags)
	return c
}
