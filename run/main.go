package main

import (
	"os"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)

	packit.Run(
		gogenerate.Detect(),
		gogenerate.Build(
			gogenerate.NewGenerateConfigurationParser(),
			gogenerate.NewGenerate(
				pexec.NewExecutable("go"),
				logger,
				chronos.DefaultClock,
			),
			logger,
		),
	)
}
