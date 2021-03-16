package modpostgres

import (
	"context"

	"github.com/go-pg/pg/v9"
	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
)

var Module = fx.Provide(
	NewDB,
)

func Register(lifecycle fx.Lifecycle, db *pg.DB) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})
}

func NewDB(cfg *modcfg.Spec) (*pg.DB, error) {
	opt, err := pg.ParseURL(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	return pg.Connect(opt), nil
}
