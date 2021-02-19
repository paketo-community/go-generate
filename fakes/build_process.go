package fakes

import (
	"sync"

	gogenerate "github.com/paketo-buildpacks/go-generate"
)

type BuildProcess struct {
	ExecuteCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			WorkingDir string
			Config     gogenerate.GenerateConfiguration
		}
		Returns struct {
			Error error
		}
		Stub func(string, gogenerate.GenerateConfiguration) error
	}
}

func (f *BuildProcess) Execute(param1 string, param2 gogenerate.GenerateConfiguration) error {
	f.ExecuteCall.Lock()
	defer f.ExecuteCall.Unlock()
	f.ExecuteCall.CallCount++
	f.ExecuteCall.Receives.WorkingDir = param1
	f.ExecuteCall.Receives.Config = param2
	if f.ExecuteCall.Stub != nil {
		return f.ExecuteCall.Stub(param1, param2)
	}
	return f.ExecuteCall.Returns.Error
}
