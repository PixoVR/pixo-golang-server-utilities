package servicetest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/cucumber/godog"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

func (s *ServerTestFeature) SendRequestWithData(method string, endpoint string, body *godog.DocString) {
	s.MakeRequest(method, endpoint, body)
}

func (s *ServerTestFeature) SendGqlRequestWithVariables(gqlMethodName string, serviceName string, endpoint string, variableBody *godog.DocString) error {
	s.GraphQLOperation = gqlMethodName
	if s.DirectoryFilePath == "" {
		return errors.New("cannot read graphql operation from file")
	}

	fileContent, err := os.ReadFile(s.DirectoryFilePath)
	if err != nil {
		return err
	}

	getVariables := func(variableBody *godog.DocString) (map[string]interface{}, error) {
		if variableBody == nil {
			return nil, nil
		}

		variables := map[string]interface{}{}

		if s.GraphQLResponse == nil {
			s.GraphQLResponse = make(map[string]interface{})
		}

		// loop through graph response and replace all the string that matches $key
		for k, v := range s.GraphQLResponse {
			variableBody.Content = strings.ReplaceAll(variableBody.Content, fmt.Sprintf("$%s", k), fmt.Sprintf("%v", v))
		}

		variableBody.Content = string(s.replaceSubstitutions([]byte(variableBody.Content)))
		log.Debug().Msgf("GraphQL variables body: %s", variableBody.Content)

		if err = json.Unmarshal([]byte(variableBody.Content), &variables); err != nil {
			return nil, err
		}

		return variables, nil
	}

	variables, err := getVariables(variableBody)
	if err != nil {
		return err
	}

	graphqlRequest := struct {
		OperationName string         `json:"operationName"`
		Query         string         `json:"query"`
		Variables     map[string]any `json:"variables,omitempty"`
	}{
		OperationName: gqlMethodName,
		Query:         string(fileContent),
		Variables:     variables,
	}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(graphqlRequest); err != nil {
		return err
	}

	return s.makeGraphQLRequest(endpoint, serviceName, buf.String())
}

func (s *ServerTestFeature) ExtractValueFromResponse(keyName string) error {
	doc, err := jsonquery.Parse(strings.NewReader(s.ResponseString))
	if err != nil {
		return err
	}

	extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("//%s/%s", s.GraphQLOperation, keyName))
	if extractedValue == nil {
		return fmt.Errorf("key %s not found in response", keyName)
	}

	s.GraphQLResponse[keyName] = extractedValue.FirstChild.Data

	return nil
}

func (s *ServerTestFeature) ShouldNotContainForJsonQueryPath(value, jsonQueryPath string) error {
	doc, err := jsonquery.Parse(strings.NewReader(s.ResponseString))
	if err != nil {
		return err
	}

	jsonQueryPath = strings.ReplaceAll(jsonQueryPath, ".", "/")

	extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("//%s", jsonQueryPath))
	if extractedValue == nil {
		return fmt.Errorf("json query '%s' not found in response", jsonQueryPath)
	}

	dataFound := false

	for _, child := range extractedValue.ChildNodes() {
		if child.Value() == value {
			dataFound = true
			break
		}
	}

	if dataFound {
		return fmt.Errorf("value %s found in json query path %s", value, jsonQueryPath)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldContainAThatIsNotNull(jsonQueryPath string) error {
	doc, err := jsonquery.Parse(strings.NewReader(s.ResponseString))
	if err != nil {
		return err
	}

	jsonQueryPath = strings.ReplaceAll(jsonQueryPath, ".", "/")

	extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("//%s", jsonQueryPath))
	if extractedValue == nil {
		return fmt.Errorf("json query '%s' not found in response", jsonQueryPath)
	}

	dataFound := false

	for _, child := range extractedValue.ChildNodes() {
		if child.Value() != nil {
			dataFound = true
			break
		}
	}

	if dataFound {
		return fmt.Errorf("the json query path %s contains s null value", jsonQueryPath)
	}

	return nil
}

func (s *ServerTestFeature) SendRequest(method, endpoint string) {
	s.MakeRequest(method, endpoint, nil)
}

func (s *ServerTestFeature) TheResponseCodeShouldBe(statusCode int) error {
	if s.StatusCode != statusCode {
		return fmt.Errorf("expected response code %d, but got %d", statusCode, s.StatusCode)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldContain(body *godog.DocString) {
	actual := TrimString(s.ResponseString)
	expected := TrimString(body.Content)
	Expect(actual).To(ContainSubstring(expected))
}

func (s *ServerTestFeature) TheResponseHeadersShouldContain(key, value string) {
	headerValue := s.Response.Header.Get(key)
	Expect(headerValue).To(ContainSubstring(value))
}

func (s *ServerTestFeature) TheResponseShouldContainA(key string) {
	key = string(s.replaceSubstitutions([]byte(key)))
	actual := TrimString(s.ResponseString)
	Expect(actual).To(ContainSubstring(key))
}

func (s *ServerTestFeature) TheResponseShouldNotContainA(key string) {
	actual := TrimString(s.ResponseString)
	Expect(actual).NotTo(ContainSubstring(key))
}

func (s *ServerTestFeature) TheResponseShouldMatchJSON(body *godog.DocString) {
	actual := s.ResponseString
	expected := body.Content
	Expect(actual).To(MatchJSON(expected))
}

func (s *ServerTestFeature) WaitForSeconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func (s *ServerTestFeature) TheResponseShouldContainSetTo(property, value string) {
	Expect(s.ResponseString).To(ContainSubstring(fmt.Sprintf("\"%s\":\"%s\"", property, value)))
}

func (s *ServerTestFeature) UseID(id string) {
	s.ID = string(ReplaceRandomID([]byte(id)))
}

func (s *ServerTestFeature) SetFilePath(filename string, directory string) error {
	if directory == "" && filename == "" {
		return fmt.Errorf("directory and filename cannot be empty")
	}
	s.DirectoryFilePath = fmt.Sprintf("./%s/%s", directory, filename)
	return nil
}

func (s *ServerTestFeature) FileToSendInRequest(filename, directory, key string) error {
	if directory == "" && filename == "" {
		return fmt.Errorf("directory and filename cannot be empty")
	}
	s.SendFileKey = key
	s.SendFile = fmt.Sprintf("./%s/%s", directory, filename)

	if _, err := os.Stat(s.SendFile); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", s.SendFile)
	}

	return nil
}
