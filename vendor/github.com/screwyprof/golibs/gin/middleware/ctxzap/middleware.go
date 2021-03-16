package ctxzap

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func CtxZap(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		newCtx := ToContext(c, logger)
		fields, statusCode := performRequest(c, utc, timeFormat)

		// may add new fields extracted from ctxtags which could have been populated by a handler.
		log := FromContext(newCtx).With(fields...)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			log.Warn("Handled request")
		case statusCode >= 500:
			log.Error("Bad request")
		default:
			log.Info("Handled request")
		}
	}
}

func performRequest(c *gin.Context, utc bool, timeFormat string) ([]zapcore.Field, int) {
	start := time.Now()

	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	if utc {
		end = end.UTC()
	}

	fields := []zap.Field{
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("start_time", start.Format(timeFormat)),
		zap.String("end_time", end.Format(timeFormat)),
		zap.Duration("latency", latency),
	}

	if len(c.Errors) > 0 {
		fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
	}

	return fields, c.Writer.Status()
}
