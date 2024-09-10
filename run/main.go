package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)

	var generateEnvironment gogenerate.GenerateEnvironment
	err := env.Parse(&generateEnvironment)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to parse build configuration: %w", err))
		os.Exit(1)
	}

	packit.Run(
		gogenerate.Detect(generateEnvironment),
		gogenerate.Build(
			gogenerate.NewGenerateConfigurationParser(generateEnvironment),
			gogenerate.NewGenerate(
				pexec.NewExecutable("go"),
				logger,
				chronos.DefaultClock,
			),
			logger,
		),
	)
}
