package config

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	ServerPort	  string
	WorkerCount   int
	PollInterval  int
}

func LoadConfig() *Config {
	return &Config{
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		WorkerCount:   getEnvAsInt("WORKER_COUNT", 5),
		PollInterval:  getEnvAsInt("POLL_INTERVAL", 5),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}