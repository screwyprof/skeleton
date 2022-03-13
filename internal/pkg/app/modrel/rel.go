package modrel

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/ansel1/merry/v2"
	"github.com/go-rel/postgres"
	"github.com/go-rel/rel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
)

var Module = fx.Provide(
	NewRel,
	NewRelRepo,
)

func Register(lifecycle fx.Lifecycle, db rel.Adapter) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := db.Close(); err != nil {
				return merry.Wrap(err)
			}

			return nil
		},
	})
}

func NewRel(cfg *modcfg.Spec) (rel.Adapter, error) {
	config, err := pgx.ParseConfig(cfg.PostgresDSN)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	adapter, err := sql.Open("pgx", stdlib.RegisterConnConfig(config))
	if err != nil {
		return nil, merry.Wrap(err)
	}

	return postgres.New(adapter), nil
}

func NewRelRepo(adapter rel.Adapter, logger *zap.Logger) rel.Repository {
	r := rel.New(adapter)
	r.Instrumentation(func(ctx context.Context, op string, message string) func(err error) {
		// no op for rel functions.
		if strings.HasPrefix(op, "rel-") {
			return func(error) {}
		}

		t := time.Now()

		return func(err error) {
			duration := time.Since(t)

			if err != nil {
				logger.Error(message,
					zap.Error(err),
					zap.Duration("duration", duration),
					zap.String("operation", op))

				return
			}

			logger.Info(message,
				zap.Duration("duration", duration),
				zap.String("operation", op))
		}
	})

	return r
}
