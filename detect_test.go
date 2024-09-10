package gogenerate_test

import (
	"os"
	"testing"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		detect     packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		detect = gogenerate.Detect(gogenerate.GenerateEnvironment{
			RunGoGenerate: true,
		})
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("detects", func() {
		_, err := detect(packit.DetectContext{
			WorkingDir: workingDir,
		})

		Expect(err).NotTo(HaveOccurred())
	})

	context("when BP_GO_GENERATE is empty", func() {
		it.Before(func() {
			detect = gogenerate.Detect(gogenerate.GenerateEnvironment{
				RunGoGenerate: false,
			})
		})

		it("fails detection", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail.WithMessage("BP_GO_GENERATE is not truthy")))
		})
	})
}
