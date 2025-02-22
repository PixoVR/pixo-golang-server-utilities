package servicetest

import (
	"context"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
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
	ServiceClient abstract.AbstractClient
	Engine        *gin.Engine
	BeforeRequest func(body []byte)
	Reset         func(sc *godog.Scenario)
	Steps         []Step
}

func NewSuite(config *SuiteConfig) *ServerTestSuite {
	if config == nil {
		config = &SuiteConfig{}
	}

	if config.Opts == nil {
		config.Opts = &godog.Options{
			Output:    colors.Colored(os.Stdout),
			Randomize: time.Now().UTC().UnixNano(),
			Format:    "pretty",
		}
	}

	pflag.BoolVarP(&debug, "debug", "v", false, "enable debug logging")
	pflag.StringVarP(&region, "region", "r", "na", "region to run tests against (options: na, saudi)")
	pflag.StringVarP(&lifecycle, "lifecycle", "l", "local", "lifecycle to run tests against (options: local, dev, stage, prod)")

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind flags")
	}

	godog.BindCommandLineFlags("godog.", config.Opts)
	pflag.Parse()

	viper.Set("region", region)
	viper.Set("lifecycle", lifecycle)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	suite := &ServerTestSuite{
		Feature:   NewServerTestFeature(),
		Lifecycle: lifecycle,
		Region:    region,

		config: config,
	}

	if lifecycle == "" || lifecycle == "internal" {
		suite.Feature.PlatformClient = &platform.MockClient{}
		suite.Feature.Engine = config.Engine
	} else {
		suite.Feature.PlatformClient = platform.NewClient(urlfinder.ClientConfig{
			Lifecycle: suite.Lifecycle,
			Region:    suite.Region,
			APIKey:    os.Getenv("PIXO_API_KEY"),
		})
	}

	suite.setup()

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
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	s.loadEnv()

	viper.SetConfigName("test-config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Warn().Msg("Config file not found; ignore error if desired")
		}
	}
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
