package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the single runtime configuration object passed into application
// wiring. Keeping this object explicit avoids reading environment variables
// from random places in the codebase.
type Config struct {
	Env      string         `yaml:"env"`
	App      AppConfig      `yaml:"app"`
	HTTP     HTTPConfig     `yaml:"http"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	MinIO    MinIOConfig    `yaml:"minio"`
	Auth     AuthConfig     `yaml:"auth"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	Docs     DocsConfig     `yaml:"docs"`
}

// AppConfig describes human-facing application metadata.
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// HTTPConfig controls the public API server.
type HTTPConfig struct {
	Addr              string        `yaml:"addr"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	RequestTimeout    time.Duration `yaml:"request_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
	RateLimitRequests int           `yaml:"rate_limit_requests"`
	RateLimitWindow   time.Duration `yaml:"rate_limit_window"`
}

// DatabaseConfig is Postgres-specific for now.
// The DSN should come from a private config file or DATABASE_URL, never from
// committed public documentation with real credentials.
type DatabaseConfig struct {
	Driver       string        `yaml:"driver"`
	DSN          string        `yaml:"dsn"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time"`
}

// RedisConfig is kept as configuration only until the Redis foundation step.
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// MinIOConfig is kept as configuration only until the object storage step.
type MinIOConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	UseSSL    bool   `yaml:"use_ssl"`
}

// AuthConfig controls access and refresh token behavior.
type AuthConfig struct {
	AccessTokenSecret  string        `yaml:"access_token_secret"`
	RefreshTokenSecret string        `yaml:"refresh_token_secret"`
	AccessTokenTTL     time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL    time.Duration `yaml:"refresh_token_ttl"`
}

// SMTPConfig holds email provider settings for invitations and future
// notifications. The app can run without SMTP until those flows are enabled.
type SMTPConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Secure      bool   `yaml:"secure"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	FromAddress string `yaml:"from_address"`
}

// DocsConfig controls development API documentation.
// Swagger stays disabled in production and protected in local development so
// the API shape is not exposed accidentally.
type DocsConfig struct {
	SwaggerEnabled  bool   `yaml:"swagger_enabled"`
	SwaggerUsername string `yaml:"swagger_username"`
	SwaggerPassword string `yaml:"swagger_password"`
}

// Load builds Config from defaults, optional YAML, environment variables, and
// finally command-line flags. That order lets local/private values override
// committed examples while keeping CLI flags useful for safe runtime options.
func Load(args []string) (Config, error) {
	cfg := Default()

	configPath := findConfigPath(args)
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath != "" {
		if err := loadYAML(configPath, &cfg); err != nil {
			return Config{}, err
		}
	}

	applyEnv(&cfg)

	fs := flag.NewFlagSet("olario-api", flag.ContinueOnError)
	fs.StringVar(&configPath, "config", configPath, "path to YAML config file")
	fs.StringVar(&cfg.Env, "env", cfg.Env, "application environment")
	fs.StringVar(&cfg.HTTP.Addr, "http-addr", cfg.HTTP.Addr, "HTTP server address")
	fs.DurationVar(&cfg.HTTP.RequestTimeout, "http-request-timeout", cfg.HTTP.RequestTimeout, "HTTP request timeout")
	fs.IntVar(&cfg.HTTP.RateLimitRequests, "http-rate-limit-requests", cfg.HTTP.RateLimitRequests, "allowed requests per rate-limit window")
	fs.DurationVar(&cfg.HTTP.RateLimitWindow, "http-rate-limit-window", cfg.HTTP.RateLimitWindow, "rate-limit window")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// LoadFile loads a specific YAML file and applies environment overrides.
// It is useful for command entrypoints such as migrations that have their own
// CLI flags and only need shared config parsing.
func LoadFile(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		if err := loadYAML(path, &cfg); err != nil {
			return Config{}, err
		}
	}
	applyEnv(&cfg)
	return cfg, cfg.Validate()
}

// Default returns safe local defaults with no real secrets.
func Default() Config {
	return Config{
		Env: "local",
		App: AppConfig{
			Name:    "olario-platform-backend",
			Version: "0.1.0",
		},
		HTTP: HTTPConfig{
			Addr:              ":8080",
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			RequestTimeout:    10 * time.Second,
			ShutdownTimeout:   10 * time.Second,
			RateLimitRequests: 60,
			RateLimitWindow:   time.Minute,
		},
		Database: DatabaseConfig{
			Driver:       "postgres",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  15 * time.Minute,
		},
		Redis: RedisConfig{
			Addr: "localhost:6379",
			DB:   0,
		},
		MinIO: MinIOConfig{
			Endpoint: "localhost:9000",
			Bucket:   "olario-local",
			UseSSL:   false,
		},
		Auth: AuthConfig{
			AccessTokenSecret:  "local-dev-access-secret-change-me",
			RefreshTokenSecret: "local-dev-refresh-secret-change-me",
			AccessTokenTTL:     15 * time.Minute,
			RefreshTokenTTL:    7 * 24 * time.Hour,
		},
		SMTP: SMTPConfig{
			Port: 587,
		},
		Docs: DocsConfig{
			SwaggerEnabled:  false,
			SwaggerUsername: "superadmin",
		},
	}
}

// Validate catches invalid config early before the server starts.
func (cfg Config) Validate() error {
	if cfg.Env == "" {
		return fmt.Errorf("env is required")
	}
	if cfg.HTTP.Addr == "" {
		return fmt.Errorf("http addr is required")
	}
	if cfg.HTTP.RateLimitRequests <= 0 {
		return fmt.Errorf("http rate limit requests must be greater than zero")
	}
	if cfg.HTTP.RateLimitWindow <= 0 {
		return fmt.Errorf("http rate limit window must be greater than zero")
	}
	if cfg.Database.Driver != "" && cfg.Database.Driver != "postgres" {
		return fmt.Errorf("database driver must be postgres")
	}
	return nil
}

func loadYAML(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config file %q: %w", path, err)
	}
	return nil
}

func applyEnv(cfg *Config) {
	cfg.Env = envString("APP_ENV", cfg.Env)
	cfg.HTTP.Addr = envString("HTTP_ADDR", cfg.HTTP.Addr)
	cfg.Database.DSN = envString("DATABASE_URL", cfg.Database.DSN)
	cfg.Redis.Addr = envString("REDIS_ADDR", cfg.Redis.Addr)
	cfg.Redis.Username = envString("REDIS_USERNAME", cfg.Redis.Username)
	cfg.Redis.Password = envString("REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Redis.DB = envInt("REDIS_DB", cfg.Redis.DB)
	cfg.MinIO.Endpoint = envString("MINIO_ENDPOINT", cfg.MinIO.Endpoint)
	cfg.MinIO.AccessKey = envString("MINIO_ACCESS_KEY", cfg.MinIO.AccessKey)
	cfg.MinIO.SecretKey = envString("MINIO_SECRET_KEY", cfg.MinIO.SecretKey)
	cfg.MinIO.Bucket = envString("MINIO_BUCKET", cfg.MinIO.Bucket)
	cfg.MinIO.UseSSL = envBool("MINIO_USE_SSL", cfg.MinIO.UseSSL)
	cfg.Auth.AccessTokenSecret = envString("AUTH_ACCESS_TOKEN_SECRET", cfg.Auth.AccessTokenSecret)
	cfg.Auth.RefreshTokenSecret = envString("AUTH_REFRESH_TOKEN_SECRET", cfg.Auth.RefreshTokenSecret)
	cfg.SMTP.Host = envString("PROD_SMTP_HOST", cfg.SMTP.Host)
	cfg.SMTP.Port = envInt("PROD_SMTP_PORT", cfg.SMTP.Port)
	cfg.SMTP.Secure = envBool("PROD_SMTP_SECURE", cfg.SMTP.Secure)
	cfg.SMTP.User = envString("PROD_SMTP_USER", cfg.SMTP.User)
	cfg.SMTP.Password = envString("PROD_SMTP_PASSWORD", cfg.SMTP.Password)
	cfg.SMTP.FromAddress = envString("PROD_SMTP_FROM_ADDRESS", cfg.SMTP.FromAddress)
	cfg.Docs.SwaggerEnabled = envBool("SWAGGER_ENABLED", cfg.Docs.SwaggerEnabled)
	cfg.Docs.SwaggerUsername = envString("SWAGGER_USERNAME", cfg.Docs.SwaggerUsername)
	cfg.Docs.SwaggerPassword = envString("SWAGGER_PASSWORD", cfg.Docs.SwaggerPassword)
}

func findConfigPath(args []string) string {
	for i, arg := range args {
		if arg == "--config" || arg == "-config" {
			if i+1 < len(args) {
				return args[i+1]
			}
			return ""
		}
		if strings.HasPrefix(arg, "--config=") {
			return strings.TrimPrefix(arg, "--config=")
		}
		if strings.HasPrefix(arg, "-config=") {
			return strings.TrimPrefix(arg, "-config=")
		}
	}
	return ""
}

func envString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
