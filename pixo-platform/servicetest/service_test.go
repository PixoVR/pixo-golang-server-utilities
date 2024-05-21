package servicetest_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/heartbeat"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/engine"
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/servicetest"
	"github.com/cucumber/godog"
	"github.com/rs/zerolog/log"
	"testing"
)

var (
	helloStep = Step{
		Expression: "I can say hello$",
		Handler: func() error {
			log.Info().Msg("I can say hello")
			return nil
		},
	}

	goodbyeStep = Step{
		Expression: "I can say goodbye$",
		Handler: func() error {
			log.Info().Msg("I can say goodbye")
			return nil
		},
	}

	resetFunc = func(sc *godog.Scenario) {
		log.Info().Msg("Resetting scenario")
	}
)

func TestMain(m *testing.M) {
	config := engine.Config{BasePath: "/heartbeat"}
	e := engine.NewEngine(config)

	suiteConfig := &SuiteConfig{
		Engine: e.Engine(),
		Reset:  resetFunc,
		Steps:  []Step{helloStep},
	}

	suite := NewSuite(suiteConfig)

	suite.AddSteps(goodbyeStep)

	suite.Feature.AddStaticSubstitution("$STATIC_ID", "1234")
	suite.Feature.AddDynamicSubstitution("$DYNAMIC_ID", func() string { return "4321" })

	suite.Feature.ServiceClient = heartbeat.NewClient(urlfinder.ClientConfig{
		Region:    "na",
		Lifecycle: "dev",
	})

	if suite.Lifecycle == "" {
		log.Fatal().Msg("Failed to get lifecycle")
	}

	if suite.Region == "" {
		log.Fatal().Msg("Failed to get region")
	}

	suite.Run()
}
