package main

import (
	"flag"

	"go.uber.org/fx"

	"github.com/screwyprof/skeleton/internal/app"
	"github.com/screwyprof/skeleton/internal/app/version"
)

// flags .
var (
	showBuildInfo = flag.Bool("version", false, "Display build info and exit")
)

func main() {
	flag.Parse()

	if *showBuildInfo {
		version.PrintBuildInfo()

		return
	}

	application := fx.New(app.Module)
	application.Run()
}
