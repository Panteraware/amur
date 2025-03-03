package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	ConfigNew()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	r := engine()

	r.Use(gin.Recovery())

	log.Info().Str("version", version).Str("commit", commit).Str("date", date)

	if err := engine().Run(fmt.Sprintf(":%d", Config.Port)); err != nil {
		log.Err(err).Msg("Error starting server")
	}
}

func engine() *gin.Engine {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if !strings.Contains(Config.Domain, c.Request.Host) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info().
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("user_agent", c.Request.UserAgent()).
			Msg("")
	})

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Debug().
			Str("method", httpMethod).
			Str("path", absolutePath).
			Str("handler", handlerName).
			Int("nuHandlers", nuHandlers).
			Msg("")
	}

	r.NoRoute(ServeFile)

	r.POST("/upload", UploadFile)

	return r
}
