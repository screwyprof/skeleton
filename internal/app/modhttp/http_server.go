package modhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ansel1/merry/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/app/modcfg"
)

var Module = fx.Provide(
	NewHTTPServer,
)

func Register(lifecycle fx.Lifecycle, srv *http.Server, cfg *modcfg.Spec, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Listening to http requests at port " + cfg.HTTPPort)
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					failOnError(fmt.Errorf("failed to start http server: %w", err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := srv.Shutdown(ctx); err != nil {
				return merry.Wrap(err)
			}

			return nil
		},
	})
}

func NewHTTPServer(cfg *modcfg.Spec, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: handler,
	}

	return srv
}

func failOnError(err error) {
	if err != nil {
		fmt.Printf("an error occurred: %v", err)
		os.Exit(1)
	}
}
