package gogenerate

import (
	"bytes"
	"strings"
	"time"

	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

type Generate struct {
	executable Executable
	logs       scribe.Logger
	clock      chronos.Clock
}

func NewGenerate(executable Executable, logs scribe.Logger, clock chronos.Clock) Generate {
	return Generate{
		executable: executable,
		logs:       logs,
		clock:      clock,
	}
}

func (g Generate) Execute(workingDir string, config GenerateConfiguration) error {
	buffer := bytes.NewBuffer(nil)

	args := append([]string{"generate"}, config.Flags...)
	args = append(args, config.Args...)

	g.logs.Process("Executing build process")
	g.logs.Subprocess("Running 'go %s'", strings.Join(args, " "))

	duration, err := g.clock.Measure(func() error {
		return g.executable.Execute(pexec.Execution{
			Args:   args,
			Dir:    workingDir,
			Stdout: buffer,
			Stderr: buffer,
		})
	})
	if err != nil {
		g.logs.Action("Failed after %s", duration.Round(time.Millisecond))
		g.logs.Detail(buffer.String())

		return err
	}

	g.logs.Action("Completed in %s", duration.Round(time.Millisecond))
	g.logs.Break()

	return nil
}
