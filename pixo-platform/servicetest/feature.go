package servicetest

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/httptest"
)

type ServerTestFeature struct {
	substitutions map[string]string
	Engine        *gin.Engine
	Recorder      *httptest.ResponseRecorder
	Client        *resty.Client

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

func (s *ServerTestFeature) AddStaticSubstitution(key, value string) {
	s.substitutions[key] = value
}

func (s *ServerTestFeature) Reset(interface{}) {
	s.Recorder = httptest.NewRecorder()

	if s.ServiceClient != nil {
		s.ServiceClient.SetToken("")
		s.ServiceClient.SetAPIKey("")
	}

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
