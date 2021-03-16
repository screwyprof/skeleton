package zapbuilder

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Option is a configuration parameter.
type Option func(f *options)

// options logger options.
type options struct {
	logLevel   zapcore.Level
	logEncoder zapcore.Encoder

	fields []zapcore.Field
}

func (o *options) populate(opts ...Option) {
	for _, option := range opts {
		option(o)
	}

	// Set zapcore.InfoLevel by default
	if o.logLevel < zapcore.DebugLevel || o.logLevel > zapcore.FatalLevel {
		o.logLevel = zapcore.InfoLevel
	}

	// Set zapcore.Encoder to json
	if o.logEncoder == nil {
		o.logEncoder = logEncFromString("json")
	}
}

func WithLevel(level string) Option {
	return func(o *options) {
		o.logLevel = logLevelFromString(level)
	}
}

func WithEncoder(encType string) Option {
	return func(o *options) {
		o.logEncoder = logEncFromString(encType)
	}
}

func WithFields(fields ...zap.Field) Option {
	return func(o *options) {
		o.fields = fields
	}
}

func NewLogger(opts ...Option) *zap.Logger {
	o := &options{}
	o.populate(opts...)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(o.logLevel)

	zl := zap.New(zapcore.NewCore(
		o.logEncoder,
		zapcore.Lock(os.Stdout),
		atom,
	), zap.Fields(o.fields...))

	return zl
}

// logLevelFromString Returns zapcore.Level by its name.
// If invalid logLevel is given, zapcore.InfoLevel is returned.
func logLevelFromString(logLevel string) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return zapcore.InfoLevel
	}
	return level
}

// logEncFromString Returns console colored encoder if encType is "simple", json encoder otherwise.
func logEncFromString(encType string) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	switch encType {
	case "pretty":
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	default:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	return encoder
}
