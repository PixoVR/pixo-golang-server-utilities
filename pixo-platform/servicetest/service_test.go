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

	suiteConfig := &SuiteConfig{
		Engine: e.Engine(),
		Reset: func(sc *godog.Scenario) {
			log.Info().Msg("Resetting scenario")
		},
		Steps: []Step{
			{
				Expression: "I can say hello$",
				Handler: func() error {
					log.Info().Msg("I can say hello")
					return nil
				},
			},
		},
	}

	suite := NewSuite(suiteConfig)
	suite.Run()
}
