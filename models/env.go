package models

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type EnvVars struct {
	Port      int
	Host      string
	JwtSecret string

	DatabaseURL   string
	MigrationsDir string

	LogLevel string
	Timeout  int

	TvDir string

	KitsuBaseURL string
	KitsuToken   string
}

func GetEnv() *EnvVars {
	err := LoadEnvFile(".env")
	if err != nil {
		slog.Error("", "err", err)
		panic(err)
	}
	return &EnvVars{
		Port:          EnvInt("PORT", 8080),
		Host:          EnvString("HOST", "0.0.0.0"),
		JwtSecret:     EnvString("JWT_SECRET", "your-secret-key"),
		DatabaseURL:   EnvString("DATABASE_URL", "./db.sqlite3"),
		MigrationsDir: EnvString("DATABASE_MIGRATION_DIR", "./db/migrations"),
		LogLevel:      EnvString("LOG_LEVEL", "debug"),
		Timeout:       EnvInt("TIMEOUT", 30),
		TvDir:         EnvString("TV_DIR", "/mnt/Thicc32/tv"),
		KitsuBaseURL:  EnvString("KitsuBaseURL", "kitsu.io"),
		// KitsuToken:    EnvString("KitsuToken", ""),
	}
}

func MustEnvString(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic("environment variable not set: " + name)
	}
	return env
}

func MustEnvInt(name string) int {
	env := os.Getenv(name)
	if env == "" {
		panic("environment variable not set: " + name)
	}
	value, err := strconv.Atoi(env)
	if err != nil {
		panic("invalid integer value for environment variable: " + name)
	}
	return value
}

func EnvInt(name string, defaultValue int) int {
	env := os.Getenv(name)
	if env == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(env)
	if err != nil {
		return defaultValue
	}
	return value
}
func EnvString(name string, defaultValue string) string {
	env := os.Getenv(name)
	if env == "" {
		return defaultValue
	}
	return env
}

func LoadEnvFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Split the content into lines
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		// Skip empty lines and comments
		if stripped == "" || stripped[0] == '#' {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(stripped, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in env file: %s", line)
		}

		// check if key and value are not empty
		if strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return fmt.Errorf("invalid line in env file: %s", line)
		}
		// remove any surrounding quotes from the value
		parts[1] = strings.Trim(parts[1], `"'`)

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return nil
}
