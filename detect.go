package gogenerate

import (
	"os"

	"github.com/paketo-buildpacks/packit"
)

func Detect(parser ConfigurationParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		run := os.Getenv("BP_GO_GENERATE")
		if run != "true" {
			return packit.DetectResult{}, packit.Fail.WithMessage("BP_GO_GENERATE is empty")
		}

		if _, err := parser.Parse(); err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("failed to parse go generate configuration: %w", err)
		}

		return packit.DetectResult{}, nil
	}
}
