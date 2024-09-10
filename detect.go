package gogenerate

import (
	"github.com/paketo-buildpacks/packit/v2"
)

func Detect(generateEnvironement GenerateEnvironment) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		if !generateEnvironement.RunGoGenerate {
			return packit.DetectResult{}, packit.Fail.WithMessage("BP_GO_GENERATE is not truthy")
		}

		return packit.DetectResult{}, nil
	}
}
