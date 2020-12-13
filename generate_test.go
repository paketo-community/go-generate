package gogenerate_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/go-generate/fakes"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGenerate(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir  string
		environment []string
		executable  *fakes.Executable
		logs        *bytes.Buffer

		generate gogenerate.Generate
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-directory")
		Expect(err).NotTo(HaveOccurred())

		environment = os.Environ()
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

		generate = gogenerate.NewGenerate(executable, gogenerate.NewLogEmitter(logs), clock)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("Execute", func() {
		it("runs go generate ./...", func() {
			err := generate.Execute(workingDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(executable.ExecuteCall.Receives.Execution.Args).To(Equal([]string{"generate", "./..."}))
			Expect(executable.ExecuteCall.Receives.Execution.Env).To(Equal(append(environment, fmt.Sprintf("GOPATH=%s", modCachePath))))
			Expect(executable.ExecuteCall.Receives.Execution.Dir).To(Equal(workingDir))

			Expect(logs.String()).To(ContainSubstring("  Executing build process"))
			Expect(logs.String()).To(ContainSubstring("    Running 'go generate ./...'"))
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
					err := generate.Execute(workingDir)
					Expect(err).To(MatchError(ContainSubstring("executable failed")))

					Expect(logs.String()).To(ContainSubstring("      Failed after 1s"))
					Expect(logs.String()).To(ContainSubstring("        build error stdout"))
					Expect(logs.String()).To(ContainSubstring("        build error stderr"))
				})
			})
		})
	})
}
