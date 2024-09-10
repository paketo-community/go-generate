package gogenerate

import (
	"fmt"

	"github.com/mattn/go-shellwords"
)

type GenerateConfiguration struct {
	Args  []string
	Flags []string
}

type GenerateConfigurationParser struct {
	generateEnvironment GenerateEnvironment
}

func NewGenerateConfigurationParser(generateEnvironment GenerateEnvironment) GenerateConfigurationParser {
	return GenerateConfigurationParser{
		generateEnvironment: generateEnvironment,
	}
}

func (p GenerateConfigurationParser) WithEnvironment(generateEnvironment GenerateEnvironment) GenerateConfigurationParser {
	p.generateEnvironment = generateEnvironment
	return p
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
	if p.generateEnvironment.GoGenerateArgs != "" {
		generateConfiguration.Args, err = shellwordsParser.Parse(p.generateEnvironment.GoGenerateArgs)

		if err != nil {
			return GenerateConfiguration{}, fmt.Errorf("BP_GO_GENERATE_ARGS=%q: %w", p.generateEnvironment.GoGenerateArgs, err)
		}
	}

	if p.generateEnvironment.GoGenerateFlags != "" {
		generateConfiguration.Flags, err = shellwordsParser.Parse(p.generateEnvironment.GoGenerateFlags)

		if err != nil {
			return GenerateConfiguration{}, fmt.Errorf("BP_GO_GENERATE_FLAGS=%q: %w", p.generateEnvironment.GoGenerateFlags, err)
		}
	}
	return generateConfiguration, nil
}
