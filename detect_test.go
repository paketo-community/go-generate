package gogenerate_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/go-generate/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		parser     *fakes.ConfigurationParser
		detect     packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		parser = &fakes.ConfigurationParser{}
		parser.ParseCall.Returns.GenerateConfiguration.Args = []string{"./..."}

		os.Setenv("BP_GO_GENERATE", "true")

		detect = gogenerate.Detect(parser)
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
			os.Unsetenv("BP_GO_GENERATE")
		})

		it("fails detection", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail.WithMessage("BP_GO_GENERATE is empty")))
		})
	})

	context("when the configuration parser fails", func() {
		it.Before(func() {
			parser.ParseCall.Returns.Error = errors.New("failed to parse configuration")
		})

		it("returns an error", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(ContainSubstring("failed to parse configuration")))
		})
	})
}
