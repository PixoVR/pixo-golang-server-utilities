package servicetest_test

import (
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/engine"
	. "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/servicetest"
	"testing"
)

func TestMain(m *testing.M) {
	config := engine.Config{
		BasePath:       "/api",
		InternalRoutes: true,
		ExternalRoutes: true,
	}
	e := engine.NewEngine(config)

	serviceClient := graphql_api.NewClient(urlfinder.ClientConfig{Lifecycle: "dev"})

	opts := Options{CustomEngine: e}

	suite := NewSuite(serviceClient, opts)
	suite.Run()
}
