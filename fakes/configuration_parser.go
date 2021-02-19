package fakes

import (
	"sync"

	gogenerate "github.com/paketo-buildpacks/go-generate"
)

type ConfigurationParser struct {
	ParseCall struct {
		sync.Mutex
		CallCount int
		Returns   struct {
			GenerateConfiguration gogenerate.GenerateConfiguration
			Error                 error
		}
		Stub func() (gogenerate.GenerateConfiguration, error)
	}
}

func (f *ConfigurationParser) Parse() (gogenerate.GenerateConfiguration, error) {
	f.ParseCall.Lock()
	defer f.ParseCall.Unlock()
	f.ParseCall.CallCount++
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub()
	}
	return f.ParseCall.Returns.GenerateConfiguration, f.ParseCall.Returns.Error
}
