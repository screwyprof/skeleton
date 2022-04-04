//go:build integration

package app

import (
	"log"

	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/app/fxlogger"
	"github.com/screwyprof/skeleton/internal/app/modcfg"
	"github.com/screwyprof/skeleton/internal/app/modrel"
	"github.com/screwyprof/skeleton/internal/app/modzap"
)

var TestModule fx.Option

func init() {
	cfg, err := modcfg.New()
	if err != nil {
		log.Fatalf("cannot init config: %v\n", err)
	}

	TestModule = fx.Options(
		fx.Logger(fxlogger.New(modzap.New(cfg))),
		modcfg.Module,
		modzap.Module,
		// modtracer.Module,

		modrel.Module,
		// modstorerep.Module,
		// modqueryer.Module,
		// modcmdhdlr.Module,

		// modgin.Module,
		// modhttp.Module,

		// fx.Invoke(modtracer.RegisterTracer),
		fx.Invoke(modrel.Register),
		// fx.Invoke(modhttp.Register),
		// fx.Invoke(modzap.Register),
		// fx.Invoke(modsentry.Register),
	)
}
