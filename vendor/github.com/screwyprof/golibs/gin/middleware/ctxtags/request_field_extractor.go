package ctxtags

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

const XRequestID = "X-Request-Id"

// RequestFieldExtractorFunc is a user-provided function that extracts field information from a http request.
// It is called from tags middleware on arrival of request.
// Keys and values will be added to the context tags. If there are no fields, you should return a nil.
type RequestFieldExtractorFunc func(ctx *gin.Context) map[string]interface{}

func RequestID(c *gin.Context) map[string]interface{} {
	requestID := c.Request.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = ksuid.New().String()
	}

	return map[string]interface{}{
		XRequestID: requestID,
	}
}
