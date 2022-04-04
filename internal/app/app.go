package app

import (
	"log"

	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/app/fxlogger"
	"github.com/screwyprof/skeleton/internal/app/modcfg"
	"github.com/screwyprof/skeleton/internal/app/modcmdhdlr"
	"github.com/screwyprof/skeleton/internal/app/modgin"
	"github.com/screwyprof/skeleton/internal/app/modhttp"
	"github.com/screwyprof/skeleton/internal/app/modqueryer"
	"github.com/screwyprof/skeleton/internal/app/modrel"
	"github.com/screwyprof/skeleton/internal/app/modsentry"
	"github.com/screwyprof/skeleton/internal/app/modstorerep"
	"github.com/screwyprof/skeleton/internal/app/modzap"
)

var Module = newAppModule()

func newAppModule() fx.Option {
	cfg, err := modcfg.New()
	if err != nil {
		log.Fatalf("cannot init config: %v\n", err)
	}

	app := fx.Options(
		fx.Logger(fxlogger.New(modzap.New(cfg))),
		modcfg.Module,
		modzap.Module,
		// modtracer.Module,

		modrel.Module,
		modstorerep.Module,
		modqueryer.Module,
		modcmdhdlr.Module,

		modgin.Module,
		modhttp.Module,

		// fx.Invoke(modtracer.RegisterTracer),
		fx.Invoke(modrel.Register),
		fx.Invoke(modhttp.Register),
		fx.Invoke(modzap.Register),
		fx.Invoke(modsentry.Register),
	)

	return app
}
