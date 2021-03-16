package fxlogger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type Adaptor struct {
	logger Logger
}

func New(logger *zap.Logger) *Adaptor {
	return &Adaptor{logger: logger}
}

func (la *Adaptor) Printf(format string, args ...interface{}) {
	la.logger.Debug(fmt.Sprintf(format, args...))
}
