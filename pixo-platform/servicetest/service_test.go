package servicetest_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/engine"
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/servicetest"
	"github.com/cucumber/godog"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestMain(m *testing.M) {
	config := engine.Config{BasePath: "/api"}
	e := engine.NewEngine(config)

	helloStep := Step{
		Expression: "I can say hello$",
		Handler: func() error {
			log.Info().Msg("I can say hello")
			return nil
		},
	}

	resetFunc := func(sc *godog.Scenario) {
		log.Info().Msg("Resetting scenario")
	}

	suiteConfig := &SuiteConfig{
		Engine: e.Engine(),
		Reset:  resetFunc,
		Steps:  []Step{helloStep},
	}

	suite := NewSuite(suiteConfig)

	if suite.Lifecycle == "" {
		log.Fatal().Msg("Failed to get lifecycle")
	}

	if suite.Region == "" {
		log.Fatal().Msg("Failed to get region")
	}

	suite.Run()
}
