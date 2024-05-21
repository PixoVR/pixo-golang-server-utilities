package servicetest_test

import (
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/engine"
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/servicetest"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var godogOpts = godog.Options{
	Output:    colors.Colored(os.Stdout),
	Randomize: time.Now().UTC().UnixNano(),
	Format:    "pretty",
}

func TestMain(m *testing.M) {
	config := engine.Config{
		BasePath:       "/v2",
		InternalRoutes: true,
		ExternalRoutes: true,
	}
	e := engine.NewEngine(config)

	serviceClient := graphql_api.NewClient(urlfinder.ClientConfig{Lifecycle: "dev"})

	opts := Options{
		CustomEngine: e,
		GodogOpts:    godogOpts,
	}

	suite := NewSuite(opts, serviceClient)
	suite.Run()
}
