package gogenerate

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		run := os.Getenv("BP_GO_GENERATE")
		if run != "true" {
			return packit.DetectResult{}, packit.Fail.WithMessage("BP_GO_GENERATE is empty")
		}

		return packit.DetectResult{}, nil
	}
}
