package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort      string
	MatcherInterval time.Duration
}

func LoadConfig() *Config {
	serverPort := getEnv("SERVER_PORT", ":8080")
	matcherInterval := getDurationEnv("MATCHER_INTERVAL", 3*time.Second)
	return &Config{
		ServerPort:      serverPort,
		MatcherInterval: matcherInterval,
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return time.Duration(intVal) * time.Second
		}
	}
	return time.Duration(defaultValue) * time.Second
}
