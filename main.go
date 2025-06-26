package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
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

	NewAsynqClient()
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

	log.Info().Str("version", version).Str("commit", commit).Str("date", date).Msg("")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go WatchFolder()

	if Config.UseRedis {
		go func() {
			NewAsynqServer()
		}()
	}

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", Config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msg("error starting server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("error shutting down server")
	} else {
		log.Info().Msg("shutting down server")
	}
}
