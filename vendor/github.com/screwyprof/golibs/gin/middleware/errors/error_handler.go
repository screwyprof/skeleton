package errors

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/screwyprof/golibs/gin/middleware/ctxzap"
)

type APIError interface {
	Status() int
	Code() int
	Message() string
	Extra() map[string]interface{}
	Cause() error
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			reportErrors(c)
			displayLastErrorErr(c)
		}
	}
}

func reportErrors(c *gin.Context) {
	logger := ctxzap.FromContext(c)
	hub := sentrygin.GetHubFromContext(c)
	for _, err := range c.Errors {
		apiErr := extractAPIErr(err.Err)
		logger.Error("An error occurred", zap.Error(apiErr))
		// send to sentry
		if hub != nil {
			hub.CaptureException(apiErr)
		}
	}
}

func extractAPIErr(err error) error {
	apiErr, ok := err.(APIError)
	if ok {
		return apiErr.Cause()
	}
	return err
}

func displayLastErrorErr(c *gin.Context) {
	apiErr, ok := c.Errors.Last().Err.(APIError)
	if ok {
		res := gin.H{
			"code":    apiErr.Code(),
			"message": apiErr.Message(),
		}
		if e := apiErr.Extra(); e != nil {
			res["extra"] = e
		}
		c.JSON(apiErr.Status(), res)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": "Internal Server Error",
	})
}
