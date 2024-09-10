package gogenerate

type GenerateEnvironment struct {
	RunGoGenerate   bool   `env:BP_GO_GENERATE`
	GoGenerateArgs  string `env:BP_GO_GENERATE_ARGS`
	GoGenerateFlags string `env:BP_GO_GENERATE_FLAGS`
}
