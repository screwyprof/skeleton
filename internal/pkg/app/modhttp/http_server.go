package modhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
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
			return srv.Shutdown(ctx)
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
