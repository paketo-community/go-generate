package gogenerate_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitGoGenerate(t *testing.T) {
	suite := spec.New("go-generate", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("Generate", testGenerate)
	suite("GenerateConfigurationParser", testGenerateConfigurationParser)
	suite.Run(t)
}
