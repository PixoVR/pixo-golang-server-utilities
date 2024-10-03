package servicetest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/cucumber/godog"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Step struct {
	Expression string
	Handler    interface{}
}

func (s *ServerTestFeature) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.Reset(sc)
		return ctx, nil
	})

	ctx.Step(`^I create a random integer$`, s.CreateRandomInt)
	ctx.Step(`^I create a uuid$`, s.CreateRandomUUID)

	ctx.Step(`^I am signed in as a "([^"]*)"$`, s.SignedInAsA)
	ctx.Step(`^I use the secret key for authentication$`, s.UseSecretKey)

	ctx.Step(`^I send "(GET|POST|DELETE)" request to "([^"]*)"$`, s.SendRequest)
	ctx.Step(`^I send "(GET|POST|DELETE)" request to "([^"]*)" with params$`, s.SendRequestWithParams)
	ctx.Step(`^I send "(PATCH|POST|PUT|DELETE)" request to "([^"]*)" with data$`, s.SendRequestWithData)
	ctx.Step(`^I send "(PATCH|POST|PUT|DELETE)" request to "([^"]*)" with data and "([^"]*)" encoded$`, s.SendRequestWithEncodedData)

	ctx.Step(`^I send "(GET|POST|DELETE)" request to the "([^"]*)" "([^"]*)" service at "([^"]*)"$`, s.SendRequestToService)
	ctx.Step(`^I send "(GET|POST|DELETE)" request to the "([^"]*)" "([^"]*)" service at "([^"]*)" with params$`, s.SendRequestWithParamsToService)
	ctx.Step(`^I send "(PATCH|POST|PUT|DELETE)" request to the "([^"]*)" "([^"]*)" service at "([^"]*)" with data$`, s.SendRequestWithDataToService)
	ctx.Step(`^I send "(PATCH|POST|PUT|DELETE)" request to the "([^"]*)" "([^"]*)" service at "([^"]*)" with data and "([^"]*)" encoded$`, s.SendRequestWithEncodedDataToService)

	ctx.Step(`^I have a file named "([^"]*)" in the "([^"]*)" directory that should be sent in the request with key "([^"]*)"$`, s.FileToSendInRequest)
	ctx.Step(`^I send "([^"]*)" gql request to the "([^"]*)" endpoint "([^"]*)" with the variables$`, s.SendGQLRequestWithVariables)

	ctx.Step(`^the response code should be "([^"]*)"$`, s.TheResponseCodeShouldBe)
	ctx.Step(`^the response code should be (\d+)$`, s.TheResponseCodeShouldBe)
	ctx.Step(`^the response should match json$`, s.TheResponseShouldMatchJSON)
	ctx.Step(`^the response should contain$`, s.TheResponseShouldContain)

	ctx.Step(`^I extract the "([^"]*)" from the response$`, s.ExtractValueFromResponse)
	ctx.Step(`^I extract the "([^"]*)" from the response as "([^"]*)"$`, s.ExtractValueFromResponseAs)

	ctx.Step(`^I download the "([^"]*)" link to the "([^"]*)" directory as "([^"]*)"$`, s.DownloadFileViaLink)
	ctx.Step(`^the file "([^"]*)" in the "([^"]*)" directory should exist$`, s.FileShouldExist)
	ctx.Step(`^the files "([^"]*)" and "([^"]*)" in the "([^"]*)" directory should be the same$`, s.CompareFiles)
	ctx.Step(`^the files "([^"]*)" and "([^"]*)" in the "([^"]*)" directory should not be the same$`, s.CompareFilesForInequality)

	ctx.Step(`^I wait for "([^"]*)" seconds$`, s.WaitForSeconds)

	ctx.Step(`^the websocket is connected$`, s.WebsocketIsConnected)
	ctx.Step(`^the websocket is not connected$`, s.WebsocketIsNotConnected)
	ctx.Step(`^I open a websocket at "([^"]*)"$`, s.OpenWebsocket)
	ctx.Step(`^I set the websocket read timeout to "([^"]*)" seconds$`, s.SetWebsocketReadTimeout)
	ctx.Step(`^I send the following data to the websocket:$`, s.SendWebsocketMessage)
	ctx.Step(`^I read a message from the websocket$`, s.GetWebsocketMessage)

	ctx.Step(`^the response should have a length of (\d+)$`, s.ResponseShouldHaveLength)

	ctx.Step(`^the response should contain a "([^"]*)" header with value "([^"]*)"$`, s.TheResponseHeadersShouldContain)
	ctx.Step(`^the response should contain a "([^"]*)"$`, s.TheResponseShouldContainA)
	ctx.Step(`^the response should not contain a "([^"]*)"$`, s.TheResponseShouldNotContainA)
	ctx.Step(`^the response should not contain "([^"]*)" at path "([^"]*)"$`, s.ShouldNotContainForJsonQueryPath)
	ctx.Step(`^the response should contain a "([^"]*)" that is null$`, s.TheResponseShouldContainAThatIsNull)
	ctx.Step(`^the response should contain a "([^"]*)" that is not null$`, s.TheResponseShouldContainAThatIsNotNull)
	ctx.Step(`^the response should contain a "([^"]*)" that is not empty$`, s.TheResponseShouldContainAThatIsNotEmpty)
	ctx.Step(`^the response should contain a "([^"]*)" set to "([^"]*)"$`, s.TheResponseShouldContainSetTo)
	ctx.Step(`^the response should contain a "([^"]*)" with length (\d+)$`, s.ResponseShouldContainPropertyWithLength)
	ctx.Step(`^the response should contain a "([^"]*)" with "([^"]*)" set to "([^"]*)"$`, s.ResponseContainsObjectWithPropertySetTo)
	ctx.Step(`^the response should contain a "([^"]*)" with item at index (\d+) with "([^"]*)" set to "([^"]*)"$`, s.ResponseContainsObjectWithItemAtIndexWithPropertySetTo)
	ctx.Step(`^the response should contain a "([^"]*)" with item at index (\d+) with "([^"]*)" of length (\d+)$`, s.ResponseContainsObjectWithItemAtIndexWithPropertyOfLength)
	ctx.Step(`^the response should contain a "([^"]*)" with item at index (\d+) with "([^"]*)" that is null$`, s.ResponseContainsObjectWithItemAtIndexWithPropertySetToNull)
	ctx.Step(`^the response should contain a "([^"]*)" with item at index (\d+) with "([^"]*)" that is not null$`, s.ResponseContainsObjectWithItemAtIndexWithPropertySetToNull)

	ctx.Step(`^the response should contain an item at index (\d+) with "([^"]*)" that is not null$`, s.ResponseContainsItemAtIndexWithPropertyNotSetToNull)
	ctx.Step(`^the response should contain an item at index (\d+) with "([^"]*)" that is null$`, s.ResponseContainsItemAtIndexWithPropertySetToNull)
	ctx.Step(`^the response should contain an item at index (\d+) with "([^"]*)" that is not null$`, s.ResponseContainsItemAtIndexWithPropertyNotSetToNull)
	ctx.Step(`^the response should contain an item at index (\d+) with "([^"]*)" set to "([^"]*)"$`, s.ResponseContainsItemAtIndexWithPropertySetTo)
	ctx.Step(`^the response should contain an item at index (\d+) with "([^"]*)" of length (\d+)$`, s.ResponseContainsItemAtIndexWithPropertyOfLength)

	ctx.Step(`^the response should contain an item with "([^"]*)" set to "([^"]*)"$`, s.ResponseContainsItemWithPropertySetTo)

	ctx.Step(`^the message should contain a "([^"]*)"$`, s.TheMessageShouldContainA)
	ctx.Step(`^the message should not contain a "([^"]*)"$`, s.TheMessageShouldNotContainA)
	ctx.Step(`^the message should not be empty$`, s.CheckMessageNotEmpty)
}

func (s *ServerTestFeature) SendRequest(method, endpoint string) error {
	return s.SendRequestToService(method, "", "", endpoint)
}

func (s *ServerTestFeature) SendRequestToService(method, tenant, service, endpoint string) error {
	return s.MakeRequest(method, tenant, service, endpoint, nil, nil)
}

func (s *ServerTestFeature) SendRequestWithParams(method, endpoint string, params *godog.DocString) error {
	return s.SendRequestWithParamsToService(method, "", "", endpoint, params)
}

func (s *ServerTestFeature) SendRequestWithParamsToService(method, tenant, service, endpoint string, params *godog.DocString) error {
	var paramsMap map[string]string
	if err := json.Unmarshal(s.PerformSubstitutions([]byte(params.Content)), &paramsMap); err != nil {
		return fmt.Errorf("error unmarshalling params: %w", err)
	}

	return s.MakeRequest(method, tenant, service, endpoint, nil, paramsMap)
}

func (s *ServerTestFeature) SendRequestWithData(method, endpoint string, body *godog.DocString) error {
	return s.SendRequestWithDataToService(method, "", "", endpoint, body)
}

func (s *ServerTestFeature) SendRequestWithDataToService(method, tenant, service, endpoint string, body *godog.DocString) error {
	return s.MakeRequest(method, tenant, service, endpoint, body, nil)
}

func (s *ServerTestFeature) SendRequestWithEncodedData(method, endpoint, encodedPath string, body *godog.DocString) error {
	return s.SendRequestWithEncodedDataToService(method, "", "", endpoint, encodedPath, body)
}

func (s *ServerTestFeature) SendRequestWithEncodedDataToService(method, tenant, service, endpoint, encodedPath string, body *godog.DocString) error {
	var data map[string]interface{}
	if err := json.Unmarshal(s.PerformSubstitutions([]byte(body.Content)), &data); err != nil {
		return fmt.Errorf("error unmarshalling body: %w", err)
	}

	encodedData, err := EncodeData(data, encodedPath)
	if err != nil {
		return fmt.Errorf("error encoding data: %w", err)
	}

	encodedBytes, err := json.Marshal(encodedData)
	if err != nil {
		return fmt.Errorf("error marshalling encoded data: %w", err)
	}

	body.Content = string(encodedBytes)

	return s.MakeRequest(method, tenant, service, endpoint, body, nil)
}

func (s *ServerTestFeature) SendGQLRequestWithVariables(gqlMethodName string, serviceName string, endpoint string, variableBody *godog.DocString) error {
	s.GraphQLOperation = gqlMethodName
	s.DirectoryFilePath = fmt.Sprintf("./gql/%s.gql", gqlMethodName)

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

		variableBody.Content = string(s.PerformSubstitutions([]byte(variableBody.Content)))
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

	return s.MakeGraphQLRequest(endpoint, serviceName, buf.String())
}

func (s *ServerTestFeature) ExtractValueFromResponse(keyName string) error {
	if err := s.ExtractValueFromResponseAs(keyName, keyName); err != nil {
		return err
	}

	if strings.ToLower(keyName) == "id" {
		if s.GraphQLResponse["id"] != nil {
			s.ID = s.GraphQLResponse["id"].(string)
		}
	}

	return nil
}

func (s *ServerTestFeature) ExtractValueFromResponseAs(keyName, identifier string) error {
	doc, err := jsonquery.Parse(strings.NewReader(s.ResponseString))
	if err != nil {
		return err
	}

	extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("/%s", keyName))
	if extractedValue == nil {
		return fmt.Errorf("key %s not found in response", keyName)
	}

	s.GraphQLResponse[identifier] = extractedValue.FirstChild.Data

	if identifier != "id" {
		s.dynamicSubstitutions[identifier] = func(data []byte) string {
			if s.GraphQLResponse[identifier] == nil {
				return ""
			}
			return s.GraphQLResponse[identifier].(string)
		}
	}

	return nil
}

func (s *ServerTestFeature) FileShouldExist(filename, directory string) error {
	if directory == "" && filename == "" {
		return fmt.Errorf("directory and filename cannot be empty")
	}
	filePath := fmt.Sprintf("./%s/%s", directory, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filePath)
	}
	return nil
}

func (s *ServerTestFeature) CompareFiles(file1, file2, dir string) error {
	file1Contents, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, file1))
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", file1, err)
	}

	file2Contents, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, file2))
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", file2, err)
	}

	if !bytes.Equal(file1Contents, file2Contents) {
		return fmt.Errorf("file contents do not match")
	}

	return nil
}

func (s *ServerTestFeature) CompareFilesForInequality(file1, file2, dir string) error {
	if s.CompareFiles(file1, file2, dir) == nil {
		return fmt.Errorf("file contents match")
	}

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
	val, err := s.getJSONNodeFromResponse(jsonQueryPath)
	if err != nil {
		return err
	}

	dataFound := false

	for _, child := range val.ChildNodes() {
		if child.Value() != nil && child.Value() != "<nil>" {
			dataFound = true
			break
		}
	}

	if !dataFound {
		return fmt.Errorf("the json query path %s contains a null value", jsonQueryPath)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldContainAThatIsNull(jsonQueryPath string) error {
	val, err := s.getJSONNodeFromResponse(jsonQueryPath)
	if err != nil {
		return err
	}

	dataFound := false

	for _, child := range val.ChildNodes() {
		if child.Value() != nil && child.Value() != "<nil>" {
			dataFound = true
			break
		}
	}

	if dataFound {
		return fmt.Errorf("the json query path %s does not contain a null value: %s", jsonQueryPath, s.ResponseString)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldContainAThatIsNotEmpty(jsonQueryPath string) error {
	val, err := s.getJSONNodeFromResponse(jsonQueryPath)
	if err != nil {
		return err
	}

	dataFound := false

	for _, child := range val.ChildNodes() {
		if child.Value() != nil && child.Value() != "<nil>" && child.Value() != "" {
			dataFound = true
			break
		}
	}

	if !dataFound {
		return fmt.Errorf("the json query path %s contains a null or empty value", jsonQueryPath)
	}

	return nil
}

func (s *ServerTestFeature) getJSONNodeFromResponse(queryPath string) (*jsonquery.Node, error) {
	doc, err := jsonquery.Parse(strings.NewReader(s.ResponseString))
	if err != nil {
		return nil, err
	}

	queryPath = strings.ReplaceAll(queryPath, ".", "/")

	extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("//%s", queryPath))
	if extractedValue == nil {
		return nil, fmt.Errorf("json query '%s' not found in response", queryPath)
	}

	return extractedValue, nil
}

func (s *ServerTestFeature) TheResponseCodeShouldBe(statusCode int) error {
	if s.StatusCode != statusCode {
		return fmt.Errorf("expected response code %d, but got %d: %s", statusCode, s.StatusCode, s.ResponseString)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldContain(body *godog.DocString) error {
	actual := TrimString(s.ResponseString)
	expected := TrimString(body.Content)
	if !strings.Contains(actual, expected) {
		return fmt.Errorf("expected response to contain %s, but got %s", expected, actual)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseHeadersShouldContain(key, value string) error {
	headerValue := s.HTTPResponse.Header.Get(key)
	if !strings.Contains(headerValue, value) {
		return fmt.Errorf("expected response header %s to contain %s, but got %s", key, value, headerValue)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldContainA(key string) error {
	key = string(s.PerformSubstitutions([]byte(key)))
	actual := TrimString(s.ResponseString)

	if !strings.Contains(actual, key) {
		return fmt.Errorf("expected response to contain %s, but got %s", key, actual)
	}

	return nil
}

func (s *ServerTestFeature) TheResponseShouldNotContainA(key string) error {
	actual := TrimString(s.ResponseString)
	if strings.Contains(actual, key) {
		return fmt.Errorf("expected response to not contain %s, but got %s", key, actual)

	}
	return nil
}

func (s *ServerTestFeature) TheResponseShouldMatchJSON(body *godog.DocString) {
	actual := s.ResponseString
	expected := body.Content
	Expect(actual).To(MatchJSON(expected))
}

func (s *ServerTestFeature) WaitForSeconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func (s *ServerTestFeature) TheResponseShouldContainSetTo(property, value string) error {
	value = string(s.PerformSubstitutions([]byte(value)))
	expectedString := fmt.Sprintf("\"%s\":\"%s\"", property, value)
	savedByString := strings.Contains(s.ResponseString, expectedString)

	isNullOrBool := value == "null" || value == "true" || value == "false"
	expectedBool := fmt.Sprintf("\"%s\":%s", property, value)
	savedByBool := isNullOrBool && strings.Contains(s.ResponseString, expectedBool)

	intValue, err := strconv.Atoi(value)
	isNum := err == nil
	expectedInt := fmt.Sprintf("\"%s\":%d", property, intValue)
	savedByInt := isNum && strings.Contains(s.ResponseString, expectedInt)

	_, err = strconv.ParseFloat(value, 64)
	isFloat := err == nil
	expectedFloat := fmt.Sprintf("\"%s\":%s", property, value)
	savedByFloat := isFloat && strings.Contains(s.ResponseString, expectedFloat)

	if !(savedByString || savedByBool || savedByInt || savedByFloat) {
		return fmt.Errorf("expected response to contain %s set to %s", property, value)
	}

	return nil
}

func (s *ServerTestFeature) FileToSendInRequest(filename, directory, key string) error {
	if directory == "" && filename == "" {
		return fmt.Errorf("directory and filename cannot be empty")
	}

	filepath := fmt.Sprintf("./%s/%s", directory, filename)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filepath)
	}

	s.FilesToSend = append(s.FilesToSend, FileToSend{
		Key:  key,
		Path: filepath,
	})
	return nil
}

func (s *ServerTestFeature) DownloadFileViaLink(keyName, downloadDirectory, filename string) error {
	if keyName == "" {
		return fmt.Errorf("key name cannot be empty")
	}

	dataFromMap := s.GraphQLResponse[keyName]
	if dataFromMap == nil {
		return fmt.Errorf("key %s not found in response", keyName)
	}

	url := dataFromMap.(string)

	if !strings.HasPrefix(url, "http") {
		return fmt.Errorf("url '%s' is not valid", url)
	}

	filepath := fmt.Sprintf("./%s/%s", downloadDirectory, filename)

	return s.DownloadFile(filepath, url)
}

func (s *ServerTestFeature) DownloadFile(filepath, url string) error {
	log.Debug().Msgf("Downloading file from %s to %s", url, filepath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)

	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200, got %d, Body: %v", response.StatusCode, response)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	reader := bytes.NewReader(bodyBytes)
	_, err = io.Copy(out, reader)
	return err
}

func (s *ServerTestFeature) ResponseShouldHaveLength(length int) error {
	items := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(s.ResponseString), &items); err != nil {
		return fmt.Errorf("failed to unmarshal response into list: %v", err)
	}

	if len(items) != length {
		return fmt.Errorf("the response contains %d items, expected %d: %s", len(items), length, s.ResponseString)
	}

	return nil
}

func (a *ServerTestFeature) getFirstNodeFromResponse(queryPath string) (*jsonquery.Node, error) {
	doc, err := jsonquery.Parse(strings.NewReader(a.ResponseString))
	if err != nil {
		return nil, err
	}

	queryPath = strings.ReplaceAll(queryPath, ".", "/")

	extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("//%s", queryPath))
	if extractedValue == nil {
		return nil, fmt.Errorf("json query '%s' not found in response", queryPath)
	}

	return extractedValue, nil
}

func (a *ServerTestFeature) ResponseShouldContainPropertyWithLength(jsonQueryPath string, length int) error {
	val, err := a.getFirstNodeFromResponse(jsonQueryPath)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Val: %v", val.Value())
	if val.Value() == nil || len(val.ChildNodes()) != length {
		return fmt.Errorf("the json query path %s does not contain %d items: %s", jsonQueryPath, length, a.ResponseString)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsItemWithPropertySetTo(property, value string) error {
	value = string(s.PerformSubstitutions([]byte(value)))

	items := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(s.ResponseString), &items); err != nil {
		return fmt.Errorf("failed to unmarshal response into list: %v", err)
	}

	for _, item := range items {
		itemMap := item.(map[string]interface{})
		itemValue := fmt.Sprint(itemMap[property])
		if itemValue == value {
			return nil
		}
	}

	return fmt.Errorf("no item found with %s set to %s", property, value)
}

func (s *ServerTestFeature) ResponseContainsObjectWithPropertySetTo(objectName, property, value string) error {
	value = string(s.PerformSubstitutions([]byte(value)))

	response := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s.ResponseString), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response into map: %v", err)
	}

	objectValue, ok := response[objectName]
	if !ok {
		return fmt.Errorf("response does not have %s in response", objectName)
	}

	objectMap, ok := objectValue.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%s is not an object: got %v", objectName, objectValue)
	}

	itemValue := fmt.Sprint(objectMap[property])
	if itemValue != value {
		return fmt.Errorf("%s does not have %s set to %s, found %s", objectName, property, value, itemValue)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsObjectWithItemAtIndexWithPropertySetTo(objectName string, index int, property, value string) error {
	value = string(s.PerformSubstitutions([]byte(value)))

	response := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s.ResponseString), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response into map: %v", err)
	}

	objectValue, ok := response[objectName]
	if !ok {
		return fmt.Errorf("response does not have %s in response", objectName)
	}

	itemsList, ok := objectValue.([]interface{})
	if !ok {
		return fmt.Errorf("%s is not a list: got %v", objectName, objectValue)
	}

	if index >= len(itemsList) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(itemsList))
	}

	itemMap := itemsList[index].(map[string]interface{})
	itemValue := fmt.Sprint(itemMap[property])
	if itemValue != value {
		return fmt.Errorf("item at index %d does not have %s set to %s, found %s", index, property, value, itemValue)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsObjectWithItemAtIndexWithPropertyOfLength(objectName string, index int, property string, length int) error {
	response := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s.ResponseString), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response into map: %v", err)
	}

	objectValue, ok := response[objectName]
	if !ok {
		return fmt.Errorf("response does not have %s in response", objectName)
	}

	itemsList, ok := objectValue.([]interface{})
	if !ok {
		return fmt.Errorf("%s is not a list: got %v", objectName, objectValue)
	}

	if index >= len(itemsList) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(itemsList))
	}

	itemMap := itemsList[index].(map[string]interface{})

	itemValue := itemMap[property]
	if itemValue == nil {
		if length == 0 {
			return nil
		}
		return fmt.Errorf("item at index %d does not have %s set to a slice. got %v", index, property, itemValue)
	}

	if len(itemValue.([]interface{})) != length {
		return fmt.Errorf("item at index %d does not have %s set to a slice of length %d, found %d", index, property, length, len(itemValue.([]interface{})))
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsObjectWithItemAtIndexWithPropertySetToNull(objectName string, index int, property string) error {
	response := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s.ResponseString), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response into map: %v", err)
	}

	objectValue, ok := response[objectName]
	if !ok {
		return fmt.Errorf("response does not have %s in response", objectName)
	}

	itemsList, ok := objectValue.([]interface{})
	if !ok {
		return fmt.Errorf("%s is not a list: got %v", objectName, objectValue)
	}

	if index >= len(itemsList) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(itemsList))
	}

	itemMap := itemsList[index].(map[string]interface{})
	itemValue := itemMap[property]
	if itemValue != nil {
		return fmt.Errorf("item at index %d does not have %s set to null, found %v", index, property, itemValue)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsObjectWithItemAtIndexWithPropertyIsNotNull(objectName string, index int, property string) error {
	response := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s.ResponseString), &response); err != nil {
		return fmt.Errorf("failed to unmarshal response into map: %v", err)
	}

	objectValue, ok := response[objectName]
	if !ok {
		return fmt.Errorf("response does not have %s in response", objectName)
	}

	itemsList, ok := objectValue.([]interface{})
	if !ok {
		return fmt.Errorf("%s is not a list: got %v", objectName, objectValue)
	}

	if index >= len(itemsList) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(itemsList))
	}

	itemMap := itemsList[index].(map[string]interface{})
	itemValue := itemMap[property]
	if itemValue == nil {
		return fmt.Errorf("item at index %d does not have %s set to a non-null value, found nil", index, property)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsItemAtIndexWithPropertySetToNull(index int, property string) error {
	items := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(s.ResponseString), &items); err != nil {
		return fmt.Errorf("failed to unmarshal response into list: %v", err)
	}

	if index >= len(items) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(items))
	}

	itemMap := items[index].(map[string]interface{})
	itemValue := itemMap[property]
	if itemValue != nil {
		return fmt.Errorf("item at index %d does not have %s set to null, found %v", index, property, itemValue)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsItemAtIndexWithPropertyNotSetToNull(index int, property string) error {
	items := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(s.ResponseString), &items); err != nil {
		return fmt.Errorf("failed to unmarshal response into list: %v", err)
	}

	if index >= len(items) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(items))
	}

	itemMap := items[index].(map[string]interface{})
	itemValue := itemMap[property]
	if itemValue == nil {
		return fmt.Errorf("item at index %d does not have %s set to a non-null value, found nil", index, property)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsItemAtIndexWithPropertySetTo(index int, property, value string) error {
	value = string(s.PerformSubstitutions([]byte(value)))

	items := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(s.ResponseString), &items); err != nil {
		return fmt.Errorf("failed to unmarshal response into list: %v", err)
	}

	if index >= len(items) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(items))
	}

	itemMap := items[index].(map[string]interface{})
	itemValue := fmt.Sprint(itemMap[property])
	if itemValue != value {
		return fmt.Errorf("item at index %d does not have %s set to %s, found %s", index, property, value, itemValue)
	}

	return nil
}

func (s *ServerTestFeature) ResponseContainsItemAtIndexWithPropertyOfLength(index int, property string, length int) error {
	items := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(s.ResponseString), &items); err != nil {
		return fmt.Errorf("failed to unmarshal response into list: %v", err)
	}

	if index >= len(items) {
		return fmt.Errorf("not enough items in response to get item at index %d, found %d", index, len(items))
	}

	itemMap := items[index].(map[string]interface{})
	itemValue := itemMap[property]

	if itemValue == nil {
		return fmt.Errorf("item at index %d does not have %s set to a slice", index, property)
	} else if len(itemValue.([]interface{})) != length {
		return fmt.Errorf("item at index %d does not have %s set to a slice of length %d, found %d", index, property, length, len(itemValue.([]interface{})))
	}

	return nil
}
