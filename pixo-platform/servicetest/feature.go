package servicetest

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

type ServerTestFeature struct {
	staticSubstitutions  map[string]string
	dynamicSubstitutions map[string]SubstitutionFunc

	BeforeRequest func(body []byte)
	Engine        *gin.Engine
	Recorder      *httptest.ResponseRecorder
	Client        *resty.Client

	UDPConn              *net.UDPConn
	WebsocketConn        *websocket.Conn
	WebsocketReadTimeout *time.Duration

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

type SubstitutionFunc func(data []byte) []byte

func NewServerTestFeature() *ServerTestFeature {
	return &ServerTestFeature{
		Client:               resty.New(),
		staticSubstitutions:  make(map[string]string),
		dynamicSubstitutions: make(map[string]SubstitutionFunc),
	}
}

func (s *ServerTestFeature) AddStaticSubstitution(key, value string) {
	s.staticSubstitutions[key] = value
}

func (s *ServerTestFeature) AddDynamicSubstitution(key string, value SubstitutionFunc) {
	s.dynamicSubstitutions[key] = value
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
