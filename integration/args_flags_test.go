package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testArgsAndFlags(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack   occam.Pack
		docker occam.Docker

		image     occam.Image
		container occam.Container

		name   string
		source string
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose().WithNoColor()
		docker = occam.NewDocker()
	})

	context("building a simple app using BP_GO_GENERATE_ARGS and BP_GO_GENERATE_FLAGS", func() {
		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("successfully passes args to go generate", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "args_and_flags"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.Build.
				WithPullPolicy("never").
				WithBuildpacks(
					settings.Buildpacks.GoDist.Online,
					settings.Buildpacks.GoGenerate.Online,
					settings.Buildpacks.BuildPlan.Online,
				).
				WithEnv(map[string]string{
					"BP_GO_GENERATE":       "true",
					"BP_GO_GENERATE_ARGS":  "main.go",
					"BP_GO_GENERATE_FLAGS": `-run "^//go:generate sleep"`,
				}).
				Execute(name, source)
			Expect(err).ToNot(HaveOccurred(), logs.String)

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, settings.Buildpack.Name)),
				"  Executing build process",
				"    Running 'go generate -run ^//go:generate sleep main.go'",
				MatchRegexp(`      Completed in ([0-9]*(\.[0-9]*)?[a-z]+)+`),
			))

			container, err = docker.Container.Run.
				WithCommand("ls -alR /workspace").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			logs, err = docker.Container.Logs.Execute(container.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(logs.String()).NotTo(ContainSubstring("main_moved.txt"))
			Expect(logs.String()).NotTo(ContainSubstring("internal_moved.txt"))
		})
	})
}
