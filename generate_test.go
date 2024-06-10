package gogenerate_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/go-generate/fakes"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGenerate(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		executable *fakes.Executable
		logs       *bytes.Buffer

		generate gogenerate.Generate
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-directory")
		Expect(err).NotTo(HaveOccurred())

		executable = &fakes.Executable{}

		logs = bytes.NewBuffer(nil)

		now := time.Now()
		times := []time.Time{now, now.Add(1 * time.Second)}

		clock := chronos.NewClock(func() time.Time {
			if len(times) == 0 {
				return time.Now()
			}

			t := times[0]
			times = times[1:]
			return t
		})

		generate = gogenerate.NewGenerate(executable, scribe.NewLogger(logs), clock)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("Execute", func() {
		it("executes the go generate process", func() {
			err := generate.Execute(workingDir, gogenerate.GenerateConfiguration{
				Args:  []string{"main.go", "somemodule"},
				Flags: []string{"-n"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(executable.ExecuteCall.Receives.Execution.Args).To(Equal([]string{
				"generate",
				"-n",
				"main.go",
				"somemodule",
			}))
			Expect(executable.ExecuteCall.Receives.Execution.Dir).To(Equal(workingDir))

			Expect(logs.String()).To(ContainSubstring("  Executing build process"))
			Expect(logs.String()).To(ContainSubstring(`    Running 'go generate -n main.go somemodule'`))
			Expect(logs.String()).To(ContainSubstring("      Completed in 1s"))
		})

		context("failure cases", func() {
			context("the executable fails", func() {
				it.Before(func() {
					executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
						fmt.Fprintln(execution.Stdout, "build error stdout")
						fmt.Fprintln(execution.Stderr, "build error stderr")

						return errors.New("executable failed")
					}
				})

				it("returns an error", func() {
					err := generate.Execute(workingDir, gogenerate.GenerateConfiguration{Args: []string{"./..."}})
					Expect(err).To(MatchError(ContainSubstring("executable failed")))

					Expect(logs.String()).To(ContainSubstring("      Failed after 1s"))
					Expect(logs.String()).To(ContainSubstring("        build error stdout"))
					Expect(logs.String()).To(ContainSubstring("        build error stderr"))
				})
			})
		})
	})
}
