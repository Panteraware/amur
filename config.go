package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type ConfigStruct struct {
	TimeZone      string
	Port          int
	UseRedis      bool
	CanConvertHLS bool
	CanScaleVideo bool
	VideoScale    string
	RedisHost     string
	RedisPass     string
	RedisDb       int
	Domain        string
	UploadKey     string
	PublicFolder  string
	Cors          middleware.CORSConfig
	Prometheus    prom
}

type prom struct {
	Username string
	Password string
}

func ConfigNew() {
	Config = &ConfigStruct{
		TimeZone:      getEnv("TZ", "UTC"),
		UseRedis:      getEnvAsBool("USE_REDIS", false),
		CanConvertHLS: getEnvAsBool("CAN_CONVERT_HLS", false),
		CanScaleVideo: getEnvAsBool("CAN_SCALE_VIDEO", false),
		VideoScale:    getEnv("VIDEO_SCALE", "720"),
		Port:          getEnvAsInt("PORT", 3000),
		RedisHost:     getEnv("REDIS_HOST", "localhost:6379"),
		RedisPass:     getEnv("REDIS_PASS", ""),
		RedisDb:       getEnvAsInt("REDIS_DB", 1),
		Domain:        getEnv("DOMAIN", "localhost"),
		UploadKey:     getEnv("UPLOAD_KEY", ""),
		PublicFolder:  getEnv("PUBLIC_FOLDER", "/public/"),
		Cors: middleware.CORSConfig{
			AllowOrigins:     strings.Split(getEnv("CORS_ALLOW_ORIGINS", "localhost"), ","),
			AllowHeaders:     strings.Split(getEnv("CORS_ALLOW_HEADERS", "GET,HEAD,POST,PUT,PATCH,DELETE"), ","),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
		},
		Prometheus: prom{
			Username: getEnv("PROM_USERNAME", "admin"),
			Password: getEnv("PROM_PASSWORD", ""),
		},
	}

	if strings.Contains(Config.Domain, "localhost") {
		Config.Domain = fmt.Sprintf("http://localhost:%d", Config.Port)
	} else if !strings.Contains(Config.Domain, "https://") {
		Config.Domain = fmt.Sprintf("https://%s", Config.Domain)
	}

	if getEnv("PROM_PASSWORD", "") == "" {
		// Essentially create a really long unguessable password to protect metrics
		b := make([]byte, 48)
		_, err := rand.Read(b)
		if err != nil {
			return
		}

		Config.Prometheus.Password = base64.StdEncoding.EncodeToString(b)

		log.Info().Msg("no prometheus password provided so prometheus is disabled")
	}
}

var Config *ConfigStruct

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultVal
}
