package config

import (
	"errors"
	"os"
	"strconv"

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
