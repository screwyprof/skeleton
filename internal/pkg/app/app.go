package app

import (
	"log"

	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/pkg/app/fxlogger"
	"github.com/screwyprof/skeleton/internal/pkg/app/modcfg"
	"github.com/screwyprof/skeleton/internal/pkg/app/modcmdhdlr"
	"github.com/screwyprof/skeleton/internal/pkg/app/modgin"
	"github.com/screwyprof/skeleton/internal/pkg/app/modhttp"
	"github.com/screwyprof/skeleton/internal/pkg/app/modpostgres"
	"github.com/screwyprof/skeleton/internal/pkg/app/modqueryer"
	"github.com/screwyprof/skeleton/internal/pkg/app/modsentry"
	"github.com/screwyprof/skeleton/internal/pkg/app/modstorerep"
	"github.com/screwyprof/skeleton/internal/pkg/app/modzap"
)

var Module fx.Option

func init() {
	cfg, err := modcfg.New()
	if err != nil {
		log.Fatalf("cannot init config: %v\n", err)
	}

	Module = fx.Options(
		fx.Logger(fxlogger.New(modzap.New(cfg))),
		modcfg.Module,
		modzap.Module,
		// modtracer.Module,

		modpostgres.Module,
		modstorerep.Module,
		modqueryer.Module,
		modcmdhdlr.Module,

		modgin.Module,
		modhttp.Module,

		// fx.Invoke(modtracer.RegisterTracer),
		fx.Invoke(modpostgres.Register),
		fx.Invoke(modhttp.Register),
		fx.Invoke(modzap.Register),
		fx.Invoke(modsentry.Register),
	)
}
