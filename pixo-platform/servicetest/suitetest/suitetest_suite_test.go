package suitetest_test

import (
	"io"
	"testing"

	"github.com/cucumber/godog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/servicetest"
)

var (
	suite          *servicetest.ServerTestSuite
	opts           *godog.Options
	overrideCalled bool
)

func TestSuitetest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Test Suite")
}

var _ = BeforeSuite(func() {
	opts = &godog.Options{Format: "progress", Output: io.Discard}
	suite = servicetest.NewSuite(&servicetest.SuiteConfig{
		Lifecycle: "internal",
		Opts:      opts,
		Strict:    true,
		Steps: []servicetest.Step{
			{
				Expression: `^the response is '(.*)'$`,
				Handler: func(body string) {
					suite.Feature.ResponseString = body
				},
			},
			{
				Expression: `^I create a random integer$`,
				Handler: func() {
					overrideCalled = true
				},
			},
		},
	})
})

func run(feature string) int {
	runOpts := *opts
	runOpts.FeatureContents = []godog.Feature{{Name: "test.feature", Contents: []byte(feature)}}

	return godog.TestSuite{
		ScenarioInitializer: suite.InitializeScenario,
		Options:             &runOpts,
	}.Run()
}
