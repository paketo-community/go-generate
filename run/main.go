package main

import (
	"os"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
)

func main() {
	logEmitter := gogenerate.NewLogEmitter(os.Stdout)
	packit.Run(
		gogenerate.Detect(),
		gogenerate.Build(
			gogenerate.NewGenerate(pexec.NewExecutable("go"), logEmitter, chronos.DefaultClock),
			logEmitter,
		),
	)
}
