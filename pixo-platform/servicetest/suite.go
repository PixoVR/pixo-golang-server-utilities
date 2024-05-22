package servicetest

import (
	"context"
	"errors"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "github.com/onsi/gomega"
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
	Lifecycle   string
	Region      string

	config *SuiteConfig
}

type SuiteConfig struct {
	Opts          *godog.Options
	ServiceClient abstract_client.AbstractClient
	Engine        *gin.Engine
	BeforeRequest func(body []byte)
	Reset         func(sc *godog.Scenario)
	Steps         []Step
}

func NewSuite(config *SuiteConfig) *ServerTestSuite {
	viper.SetDefault("lifecycle", "local")
	viper.SetDefault("region", "na")

	pflag.StringVarP(&region, "region", "r", viper.GetString("region"), "region to run tests against (options: na, saudi)")
	pflag.StringVarP(&lifecycle, "lifecycle", "l", viper.GetString("lifecycle"), "lifecycle to run tests against (options: local, dev, stage, prod)")
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind flags")
	}

	viper.Set("region", region)
	viper.Set("lifecycle", lifecycle)

	suite := &ServerTestSuite{
		Feature:   NewServerTestFeature(),
		Lifecycle: lifecycle,
		Region:    region,

		config: config,
	}

	if lifecycle == "" || lifecycle == "local" || lifecycle == "internal" {
		suite.Feature.PlatformClient = &graphql_api.MockGraphQLClient{}

		if suite.Lifecycle == "local" {
			suite.Feature.PlatformClient = graphql_api.NewClient(urlfinder.ClientConfig{
				Lifecycle: suite.Lifecycle,
				APIKey:    os.Getenv("PIXO_API_KEY"),
			})
		}
	} else {
		suite.Feature.PlatformClient = graphql_api.NewClient(urlfinder.ClientConfig{
			Lifecycle: suite.Lifecycle,
			Region:    suite.Region,
		})
	}

	suite.Feature.Engine = config.Engine

	if suite.config.Opts == nil {
		suite.config.Opts = &godog.Options{
			Output:    colors.Colored(os.Stdout),
			Randomize: time.Now().UTC().UnixNano(),
			Format:    "pretty",
		}
	}

	suite.setup()

	pflag.Parse()
	return suite
}

func (s *ServerTestSuite) AddSteps(steps ...Step) {
	s.config.Steps = append(s.config.Steps, steps...)
}

func (s *ServerTestSuite) Run() {
	RegisterFailHandler(func(message string, _ ...int) {
		log.Panic().Msg(message)
	})

	status := godog.TestSuite{
		ScenarioInitializer: s.InitializeScenario,
		Options:             s.config.Opts,
	}.Run()

	os.Exit(status)
}

func (s *ServerTestSuite) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		if s.config.Reset != nil {
			s.config.Reset(sc)
		}
		return ctx, nil
	})

	for _, step := range s.config.Steps {
		ctx.Step(step.Expression, step.Handler)
	}

	s.Feature.InitializeScenario(ctx)
}

func (s *ServerTestSuite) setup() {
	godog.BindCommandLineFlags("godog.", s.config.Opts)

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
