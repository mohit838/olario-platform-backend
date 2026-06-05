package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppEnv   string
	AppDebug bool
	AppName  string
	AppPort  string
	AppDB    string
}

func LoadConfig(path string) (*EnvConfig, error) {

	err := godotenv.Load(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	// Get from env files
	appEnv := os.Getenv("APP_ENV")
	appDebugValue := os.Getenv("APP_DEBUG")
	appName := os.Getenv("APP_NAME")
	appPort := os.Getenv("APP_PORT")
	appDB := os.Getenv("DATABASE_URL")

	if strings.TrimSpace(appPort) == "" {
		appPort = "8080"
	}

	if strings.TrimSpace(appDB) == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	appDebug := false
	if parsed, parseErr := strconv.ParseBool(appDebugValue); parseErr == nil {
		appDebug = parsed
	}
	return &EnvConfig{
		AppEnv:   appEnv,
		AppDebug: appDebug,
		AppName:  appName,
		AppPort:  appPort,
		AppDB:    appDB,
	}, nil

}
