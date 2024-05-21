package servicetest

import (
	"errors"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/engine"
	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)

type ServerTestSuite struct {
	EnvFilePath *string
	Feature     *ServerTestFeature

	opts Options
}

type Options struct {
	GodogOpts    godog.Options
	CustomEngine *engine.CustomEngine
}

func NewSuite(opts Options, serviceClient abstract_client.AbstractClient) *ServerTestSuite {
	suite := &ServerTestSuite{
		opts: opts,
		Feature: &ServerTestFeature{
			ServiceClient: serviceClient,
		},
	}

	if opts.CustomEngine != nil {
		suite.Feature.Engine = opts.CustomEngine.Engine()
	}

	suite.setup()

	return suite
}

func (s *ServerTestSuite) Run() {
	status := godog.TestSuite{
		ScenarioInitializer: s.Feature.InitializeScenario,
		Options:             &s.opts.GodogOpts,
	}.Run()

	os.Exit(status)
}

func (s *ServerTestSuite) setup() {
	godog.BindCommandLineFlags("godog.", &s.opts.GodogOpts)
	pflag.Parse()

	gomega.RegisterFailHandler(func(message string, _ ...int) {
		log.Panic().Msg(message)
	})

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	s.loadEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("lifecycle", "local")
	viper.SetDefault("region", "na")

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Warn().Msg("Config file not found; ignore error if desired")
		}
	}

	initLogger()
}

func (s *ServerTestSuite) loadEnv() {
	if s.EnvFilePath != nil {
		if _, err := os.Stat(*s.EnvFilePath); err != nil {
			_ = godotenv.Load(*s.EnvFilePath)
		}
	} else {
		if _, err := os.Stat(".env"); os.IsNotExist(err) {
			if _, err = os.Stat("../.env"); os.IsNotExist(err) {
				log.Warn().Msg("No env file found")
			} else {
				_ = godotenv.Load("../.env")
			}
		} else {
			_ = godotenv.Load(".env")
		}
	}
}

func initLogger() {
	var enableDebug bool
	pflag.BoolVarP(&enableDebug, "debug", "z", true, "enable debug logging")

	pflag.StringP("region", "r", viper.GetString("region"), "region to run tests against (options: na, saudi)")
	pflag.StringP("lifecycle", "l", viper.GetString("lifecycle"), "lifecycle to run tests against (options: local, dev, stage, prod)")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind flags")
	}

	if enableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}
