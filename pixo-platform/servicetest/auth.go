package servicetest

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func (s *ServerTestFeature) UseSecretKey() error {
	key := strings.ToUpper(fmt.Sprintf("%s_%s", viper.GetString("lifecycle"), "SECRET_KEY"))

	s.SecretKey = os.Getenv(key)
	if s.SecretKey == "" {
		return errors.New("secret key not found")
	}

	return nil
}

func (s *ServerTestFeature) SignedInAsA(role string) error {
	username := getUsername(role)
	if username == "" {
		return errors.New("username not found")
	}

	password := getPassword(role)
	if password == "" {
		return errors.New("password not found")
	}

	if s.PlatformClient == nil {
		return errors.New("platform client not initialized")
	}

	if err := s.PlatformClient.Login(username, password); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	s.Token = s.ServiceClient.GetToken()
	if s.Token == "" {
		return errors.New("token not found")
	}

	s.UserID = s.PlatformClient.ActiveUserID()
	if s.UserID == 0 {
		return errors.New("user id not found")
	}

	return nil
}

func getUsername(role string) string {
	return os.Getenv(strings.ToUpper(role) + "_USERNAME")
}

func getPassword(role string) string {
	return os.Getenv(strings.ToUpper(role) + "_PASSWORD")
}
