package main

import (
	"os"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
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
