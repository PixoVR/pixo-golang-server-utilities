package servicetest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/httptest"
)

type ServerTestFeature struct {
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder

	Client   *resty.Client
	Response *http.Response

	ID             string
	UserID         int
	APIKey         string
	Err            error
	Token          string
	ResponseString string
	StatusCode     int

	DirectoryFilePath string
	GraphQLOperation  string
	SendFileKey       string
	SendFile          string
	GraphQLResponse   map[string]interface{}
}

func (s *ServerTestFeature) resetResponse(interface{}) {
	s.Recorder = httptest.NewRecorder()

	s.Client = resty.New()
	s.Response = nil

	s.Token = ""
	s.APIKey = ""
	s.Err = nil
	s.GraphQLResponse = make(map[string]interface{})
	s.DirectoryFilePath = ""
	s.SendFileKey = ""
	s.SendFileKey = ""
}
