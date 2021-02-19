package gogenerate

import (
	"os"

	"github.com/mattn/go-shellwords"
)

type GenerateConfiguration struct {
	Args  []string
	Flags []string
}

type GenerateConfigurationParser struct {
}

func NewGenerateConfigurationParser() GenerateConfigurationParser {
	return GenerateConfigurationParser{}
}

func (p GenerateConfigurationParser) Parse() (GenerateConfiguration, error) {
	var (
		generateConfiguration GenerateConfiguration
		err                   error
		shellwordsParser      *shellwords.Parser
	)

	shellwordsParser = shellwords.NewParser()
	shellwordsParser.ParseEnv = true

	generateConfiguration.Args = []string{"./..."}
	if val, ok := os.LookupEnv("BP_GO_GENERATE_ARGS"); ok {
		generateConfiguration.Args, err = shellwordsParser.Parse(val)

		if err != nil {
			return GenerateConfiguration{}, err
		}
	}

	if val, ok := os.LookupEnv("BP_GO_GENERATE_FLAGS"); ok {
		generateConfiguration.Flags, err = shellwordsParser.Parse(val)

		if err != nil {
			return GenerateConfiguration{}, err
		}
	}
	return generateConfiguration, nil
}
