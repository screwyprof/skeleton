package ctxzap

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/screwyprof/golibs/gin/middleware/ctxtags"
)

type ctxLogger struct {
	logger *zap.Logger
	fields []zapcore.Field
}

var (
	ctxKey    = "ctxzap"
	nopLogger = zap.NewNop()
)

// FromContext creates the call-scoped Logger with ctxtags values set as fields.
func FromContext(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(ctxKey).(*ctxLogger)
	if !ok || l.logger == nil {
		return nopLogger
	}
	fields := tagsToFields(ctx)
	fields = append(fields, l.fields...)
	return l.logger.With(fields...)
}

// tagsToFields transforms the Tags on the supplied context into zap fields.
func tagsToFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field
	tags := ctxtags.FromContext(ctx)
	for k, v := range tags.Values() {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// ToContext adds zap.Logger to the context for later extraction.
func ToContext(c *gin.Context, logger *zap.Logger) *gin.Context {
	l := &ctxLogger{
		logger: logger,
	}
	c.Set(ctxKey, l)
	return c
}
