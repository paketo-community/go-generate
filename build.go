package gogenerate

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface ConfigurationParser --output fakes/configuration_parser.go
type ConfigurationParser interface {
	Parse() (GenerateConfiguration, error)
}

//go:generate faux --interface BuildProcess --output fakes/build_process.go
type BuildProcess interface {
	Execute(workingDir string, config GenerateConfiguration) error
}

func Build(parser ConfigurationParser, buildProcess BuildProcess, logs scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		config, err := parser.Parse()
		if err != nil {
			return packit.BuildResult{}, err
		}

		err = buildProcess.Execute(context.WorkingDir, config)
		if err != nil {
			return packit.BuildResult{}, err
		}
		return packit.BuildResult{}, nil
	}
}
