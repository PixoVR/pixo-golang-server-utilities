package servicetest

import (
	"context"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/httptest"
)

type ServerTestFeature struct {
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder
	Client   *resty.Client

	ServiceClient abstract_client.AbstractClient

	Response       *http.Response
	ResponseString string
	StatusCode     int

	Token   string
	ID      string
	UserID  int
	APIKey  string
	Message string
	Err     error

	DirectoryFilePath string
	GraphQLOperation  string
	SendFileKey       string
	SendFile          string
	GraphQLResponse   map[string]interface{}
}

func (s *ServerTestFeature) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.Reset(sc)
		return ctx, nil
	})

	ctx.Step(`^I use the id "([^"]*)" for the following requests$`, s.UseID)
	ctx.Step(`^I send "(GET|POST|DELETE)" request to "([^"]*)"$`, s.SendRequest)
	ctx.Step(`^I send "(PATCH|POST|PUT|DELETE)" request to "([^"]*)" with data$`, s.SendRequestWithData)
	ctx.Step(`^I send "([^"]*)" gql request to the "([^"]*)" endpoint "([^"]*)" with the variables$`, s.SendGqlRequestWithVariables)

	ctx.Step(`^the response code should be "([^"]*)"$`, s.TheResponseCodeShouldBe)
	ctx.Step(`^the response code should be (\d+)$`, s.TheResponseCodeShouldBe)
	ctx.Step(`^the response should match json$`, s.TheResponseShouldMatchJSON)
	ctx.Step(`^the response should contain$`, s.TheResponseShouldContain)
	ctx.Step(`^the response should contain a "([^"]*)" header with value "([^"]*)"$`, s.TheResponseHeadersShouldContain)
	ctx.Step(`^the response should contain a "([^"]*)"$`, s.TheResponseShouldContainA)
	ctx.Step(`^the response should not contain a "([^"]*)"$`, s.TheResponseShouldNotContainA)
	ctx.Step(`^I extract the "([^"]*)" from the response$`, s.ExtractValueFromResponse)

	ctx.Step(`^I have a file named "([^"]*)" in the "([^"]*)" directory$`, s.SetFilePath)
	ctx.Step(`^I have a file named "([^"]*)" in the "([^"]*)" directory that should be sent in the request with key "([^"]*)"$`, s.FileToSendInRequest)

	ctx.Step(`^I wait for "([^"]*)" seconds$`, s.WaitForSeconds)
	ctx.Step(`^the response should contain "([^"]*)" set to "([^"]*)"$`, s.TheResponseShouldContainSetTo)
	ctx.Step(`^it should not contain "([^"]*)" for the path "([^"]*)"$`, s.ShouldNotContainForJsonQueryPath)
	ctx.Step(`^the response should contain a "([^"]*)" that is not null$`, s.TheResponseShouldContainAThatIsNotNull)
}

func (s *ServerTestFeature) Reset(interface{}) {
	s.Recorder = httptest.NewRecorder()

	s.ServiceClient.SetToken("")
	s.ServiceClient.SetAPIKey("")

	s.Response = nil
	s.ResponseString = ""
	s.StatusCode = 0

	s.Token = ""
	s.ID = ""
	s.UserID = 0
	s.APIKey = ""
	s.Err = nil

	s.GraphQLResponse = make(map[string]interface{})
	s.DirectoryFilePath = ""
	s.SendFileKey = ""
	s.SendFileKey = ""
}
