package gogenerate_test

import (
	"testing"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGenerateConfigurationParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		parser gogenerate.GenerateConfigurationParser
	)

	it.Before(func() {
		parser = gogenerate.NewGenerateConfigurationParser(gogenerate.GenerateEnvironment{})
	})

	context("BP_GO_GENERATE_ARGS is set", func() {
		it.Before(func() {
			parser = parser.WithEnvironment(gogenerate.GenerateEnvironment{
				GoGenerateArgs: "somemodule othermodule",
			})
		})

		it("uses the values in the env var", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration).To(Equal(gogenerate.GenerateConfiguration{
				Args: []string{"somemodule", "othermodule"},
			}))
		})
	})

	context("BP_GO_GENERATE_ARGS is NOT set", func() {
		it("uses the default value", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration).To(Equal(gogenerate.GenerateConfiguration{
				Args: []string{"./..."},
			}))
		})
	})

	context("BP_GO_GENERATE_FLAGS is set", func() {
		it.Before(func() {
			parser = parser.WithEnvironment(gogenerate.GenerateEnvironment{
				GoGenerateFlags: `-run="^//go:generate go get" -n`,
			})
		})

		it("uses the values in the env var", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration).To(Equal(gogenerate.GenerateConfiguration{
				Args:  []string{"./..."},
				Flags: []string{`-run=^//go:generate go get`, "-n"},
			}))
		})
	})

	context("failure cases", func() {

		context("generate args fail to parse", func() {
			it.Before(func() {
				parser = parser.WithEnvironment(gogenerate.GenerateEnvironment{
					GoGenerateArgs: "\"",
				})
			})

			it("returns an error", func() {
				_, err := parser.Parse()
				Expect(err).To(MatchError(ContainSubstring(`BP_GO_GENERATE_ARGS="\"": invalid command line string`)))
			})
		})

		context("generate flags fail to parse", func() {
			it.Before(func() {
				parser = parser.WithEnvironment(gogenerate.GenerateEnvironment{
					GoGenerateFlags: "\"",
				})
			})

			it("returns an error", func() {
				_, err := parser.Parse()
				Expect(err).To(MatchError(ContainSubstring(`BP_GO_GENERATE_FLAGS="\"": invalid command line string`)))
			})
		})
	})
}
