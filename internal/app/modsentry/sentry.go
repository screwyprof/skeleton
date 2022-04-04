package modsentry

import (
	"context"
	"time"

	"github.com/ansel1/merry/v2"
	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/app/modcfg"
	"github.com/screwyprof/skeleton/internal/app/version"
)

var flushTimeout = 2 * time.Second

func Register(lifecycle fx.Lifecycle, cfg *modcfg.Spec, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if cfg.SentryDSN == "" {
				logger.Info("Sentry disabled")

				return nil
			}

			if err := sentry.Init(sentry.ClientOptions{
				Dsn:         cfg.SentryDSN,
				Environment: cfg.SentryENV,
				Release:     version.AppVersion,
			}); err != nil {
				return merry.Wrap(err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			sentry.Flush(flushTimeout)

			return nil
		},
	})
}
