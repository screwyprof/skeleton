package modzap

import (
	"context"

	"github.com/screwyprof/golibs/zapbuilder"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
	"github.com/screwyprof/skeleton/internal/pkg/app/version"
)

var Module = fx.Provide(
	New,
)

func Register(lifecycle fx.Lifecycle, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Application init", zap.String("version", version.AppVersion))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// some os.Stdout syncing errors may occur, but we cannot really do anything about them.
			_ = logger.Sync()

			return nil
		},
	})
}

func New(cfg *modcfg.Spec) *zap.Logger {
	return zapbuilder.NewLogger(
		zapbuilder.WithLevel(cfg.LogLevel),
		zapbuilder.WithEncoder(cfg.LogEncType),
	)
}
