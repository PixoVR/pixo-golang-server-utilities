package servicetest

import (
	"errors"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)

var (
	lifecycle string
	region    string
	debug     bool
)

type ServerTestSuite struct {
	EnvFilePath *string
	Feature     *ServerTestFeature

	opts Options
}

type Options struct {
	GodogOpts *godog.Options
	Engine    *gin.Engine
}

func NewSuite(serviceClient abstract_client.AbstractClient, opts ...Options) *ServerTestSuite {
	viper.SetDefault("lifecycle", "local")
	viper.SetDefault("region", "na")

	pflag.StringVarP(&region, "region", "r", viper.GetString("region"), "region to run tests against (options: na, saudi)")
	pflag.StringVarP(&lifecycle, "lifecycle", "l", viper.GetString("lifecycle"), "lifecycle to run tests against (options: local, dev, stage, prod)")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind flags")
	}

	suite := &ServerTestSuite{
		Feature: &ServerTestFeature{
			ServiceClient: serviceClient,
		},
		opts: Options{},
	}

	if len(opts) > 0 {
		suite.opts = opts[0]
		suite.Feature.Engine = suite.opts.Engine
	}

	if suite.opts.GodogOpts == nil {
		suite.opts.GodogOpts = &godog.Options{
			Output:    colors.Colored(os.Stdout),
			Randomize: time.Now().UTC().UnixNano(),
			Format:    "pretty",
		}
	}

	suite.setup()

	return suite
}

func (s *ServerTestSuite) Run() {
	gomega.RegisterFailHandler(func(message string, _ ...int) {
		log.Panic().Msg(message)
	})

	status := godog.TestSuite{
		ScenarioInitializer: s.Feature.InitializeScenario,
		Options:             s.opts.GodogOpts,
	}.Run()

	os.Exit(status)
}

func (s *ServerTestSuite) setup() {
	godog.BindCommandLineFlags("godog.", s.opts.GodogOpts)
	pflag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	s.loadEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

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
	pflag.BoolVarP(&debug, "debug", "z", true, "enable debug logging")
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}
