package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"time"
)

var (
	l *zerolog.Logger
)

func init() {
	lifecycle := os.Getenv("LIFECYCLE")

	if strings.ToLower(lifecycle) == "prod" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if strings.ToLower(lifecycle) == "test" {
		zerolog.SetGlobalLevel(zerolog.Level(0)) // trace level
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	l = &log.Logger
}

func LoggingMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = l
	}

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		logger.Info().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("duration", time.Since(startTime).String()).
			Msg("Request processed")

		logger.Debug().Msgf("Request Headers: %v", c.Request.Header)

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			body, _ := io.ReadAll(c.Request.Body)
			println(string(body))
			logger.Debug().Msgf("Request Body: %s", body)
		}

		c.Next()
	}
}
