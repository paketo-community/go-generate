package gogenerate_test

import (
	"os"
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
		parser = gogenerate.NewGenerateConfigurationParser()
	})

	context("BP_GO_GENERATE_ARGS is set", func() {
		it.Before(func() {
			os.Setenv("BP_GO_GENERATE_ARGS", "somemodule othermodule")
		})

		it.After(func() {
			os.Unsetenv("BP_GO_GENERATE_ARGS")
		})

		it("uses the values in the env var", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration).To(Equal(gogenerate.GenerateConfiguration{
				Args: []string{"somemodule", "othermodule"},
			}))
		})
	}, spec.Sequential())

	context("BP_GO_GENERATE_ARGS is NOT set", func() {
		it("uses the default value", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration).To(Equal(gogenerate.GenerateConfiguration{
				Args: []string{"./..."},
			}))
		})
	}, spec.Sequential())

	context("BP_GO_GENERATE_FLAGS is set", func() {
		it.Before(func() {
			os.Setenv("BP_GO_GENERATE_FLAGS", `-run="^//go:generate go get" -n`)
		})

		it.After(func() {
			os.Unsetenv("BP_GO_GENERATE_FLAGS")
		})

		it("uses the values in the env var", func() {
			configuration, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())
			Expect(configuration).To(Equal(gogenerate.GenerateConfiguration{
				Args:  []string{"./..."},
				Flags: []string{`-run=^//go:generate go get`, "-n"},
			}))
		})
	}, spec.Sequential())

	context("failure cases", func() {

		context("generate args fail to parse", func() {
			it.Before(func() {
				os.Setenv("BP_GO_GENERATE_ARGS", "\"")
			})

			it("returns an error", func() {
				_, err := parser.Parse()
				Expect(err).To(MatchError(ContainSubstring(`BP_GO_GENERATE_ARGS="\"": invalid command line string`)))
			})

			it.After(func() {
				os.Unsetenv("BP_GO_GENERATE_ARGS")
			})
		}, spec.Sequential())

		context("generate flags fail to parse", func() {
			it.Before(func() {
				os.Setenv("BP_GO_GENERATE_FLAGS", "\"")
			})

			it("returns an error", func() {
				_, err := parser.Parse()
				Expect(err).To(MatchError(ContainSubstring(`BP_GO_GENERATE_FLAGS="\"": invalid command line string`)))
			})

			it.After(func() {
				os.Unsetenv("BP_GO_GENERATE_ARGS")
			})
		}, spec.Sequential())
	})
}
