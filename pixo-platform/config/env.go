package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func LoadEnvVars() {
	envPath := filepath.Join(GetProjectRoot(), ".env")
	err := godotenv.Load(envPath)

	if err != nil {
		log.Warn().Msgf("No .env file loaded: %s", err)
	}
}

func GetLifecycle() string {
	lifecycle, ok := os.LookupEnv("LIFECYCLE")
	if !ok {
		lifecycle = "dev"
	}

	return lifecycle
}
