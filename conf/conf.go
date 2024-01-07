package conf

import (
	"os"
	"strconv"
	"time"
)

// Config represents the server configuration
type Config struct {
	Port          int
	CertFile      string
	KeyFile       string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	MaxFileSizeMB int
}

// Load loads the configuration from environment variables
func Load() Config {
	return Config{
		Port:          getEnvInt("PORT", 8000),
		CertFile:      getEnv("CERT_FILE", "cert.pem"),
		KeyFile:       getEnv("KEY_FILE", "key.pem"),
		ReadTimeout:   getEnvDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:  getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:   getEnvDuration("IDLE_TIMEOUT", 15*time.Second),
		MaxFileSizeMB: getEnvInt("MAX_FILE_SIZE_MB", 10),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return result
}
