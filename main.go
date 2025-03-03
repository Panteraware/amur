package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	e := echo.New()

	e.RouteNotFound("/*", ServeFile)
	e.POST("/upload", UploadFile)

	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(Config.Cors))

	log.Info().Str("version", version).Str("commit", commit).Str("date", date)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", Config.Port)))
}
