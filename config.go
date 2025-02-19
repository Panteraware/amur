package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ConfigStruct struct {
	TimeZone     string
	Port         int
	UseRedis     bool
	RedisHost    string
	RedisPass    string
	RedisDb      int
	Domain       string
	UploadKey    string
	PublicFolder string
}

func ConfigNew() {
	Config = &ConfigStruct{
		TimeZone:     getEnv("TZ", "UTC"),
		UseRedis:     getEnvAsBool("USE_REDIS", false),
		Port:         getEnvAsInt("PORT", 3000),
		RedisHost:    getEnv("REDIS_HOST", "localhost:6379"),
		RedisPass:    getEnv("REDIS_PASS", ""),
		RedisDb:      getEnvAsInt("REDIS_DB", 1),
		Domain:       getEnv("DOMAIN", "localhost"),
		UploadKey:    getEnv("UPLOAD_KEY", ""),
		PublicFolder: getEnv("PUBLIC_FOLDER", "/public/"),
	}

	if strings.Contains(Config.Domain, "localhost") {
		Config.Domain = fmt.Sprintf("http://localhost:%d", Config.Port)
	} else if !strings.Contains(Config.Domain, "https://") {
		Config.Domain = fmt.Sprintf("https://%s", Config.Domain)
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
