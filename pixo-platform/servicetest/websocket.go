package servicetest

import (
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"strconv"
	"strings"
	"time"
)

func (s *ServerTestFeature) OpenWebsocket(endpoint string) error {
	s.WebsocketConn, s.Response, s.Err = s.ServiceClient.DialWebsocket(endpoint)
	if s.Err != nil {
		return fmt.Errorf("error connecting to websocket: %w", s.Err)
	}

	if s.Response != nil {
		s.StatusCode = s.Response.StatusCode
	}

	return nil
}

func (s *ServerTestFeature) WebsocketIsConnected() error {
	if s.WebsocketConn == nil {
		return errors.New("not connected to websocket")
	}

	return nil
}

func (s *ServerTestFeature) WebsocketIsNotConnected() error {
	if s.WebsocketConn != nil {
		return errors.New("websocket is connected")
	}

	return nil
}

func (s *ServerTestFeature) SetWebsocketReadTimeout(seconds string) error {
	timeout, err := strconv.Atoi(seconds)
	if err != nil {
		return err
	}

	if timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	s.WebsocketReadTimeout = &timeoutDuration

	if err = s.WebsocketConn.SetReadDeadline(time.Now().Add(*s.WebsocketReadTimeout)); err != nil {
		return fmt.Errorf("error setting read deadline: %w", err)
	}

	return nil
}

func (s *ServerTestFeature) SendWebsocketMessage(body *godog.DocString) error {
	if err := s.WebsocketIsConnected(); err != nil {
		return err
	}

	if err := s.ServiceClient.WriteToWebsocket(s.performSubstitutions([]byte(body.Content))); err != nil {
		return fmt.Errorf("error sending message to websocket: %w", err)
	}

	return nil
}

func (s *ServerTestFeature) GetWebsocketMessage() error {
	if err := s.WebsocketIsConnected(); err != nil {
		return err
	}

	s.Message = ""

	n, msg, err := s.ServiceClient.ReadFromWebsocket()
	if err != nil {
		return fmt.Errorf("error reading from websocket: %w", err)
	}

	s.Message = string(msg[:n])
	return nil
}

func (s *ServerTestFeature) CheckMessageNotEmpty() error {
	if s.Message == "" {
		return errors.New("message is empty")
	}
	return nil
}

func (s *ServerTestFeature) TheMessageShouldContainA(key string) error {
	key = string(s.performSubstitutions([]byte(key)))
	actual := TrimString(s.Message)
	if !strings.Contains(actual, key) {
		return fmt.Errorf("expected message to contain %s, but got %s", key, actual)
	}

	return nil
}

func (s *ServerTestFeature) TheMessageShouldNotContainA(key string) error {
	key = string(s.performSubstitutions([]byte(key)))
	actual := TrimString(s.Message)
	if strings.Contains(actual, key) {
		return fmt.Errorf("expected message not to contain %s, but got %s", key, actual)
	}

	return nil
}
