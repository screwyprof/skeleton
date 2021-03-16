package modsentry

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
	"github.com/screwyprof/skeleton/internal/pkg/app/version"
)

var flushTimeout = 2 * time.Second

func Register(lifecycle fx.Lifecycle, cfg *modcfg.Spec, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if cfg.SentryDSN == "" {
				logger.Info("Sentry disabled")

				return nil
			}

			return sentry.Init(sentry.ClientOptions{
				Dsn:         cfg.SentryDSN,
				Environment: cfg.SentryENV,
				Release:     version.AppVersion,
			})
		},
		OnStop: func(ctx context.Context) error {
			sentry.Flush(flushTimeout)

			return nil
		},
	})
}
