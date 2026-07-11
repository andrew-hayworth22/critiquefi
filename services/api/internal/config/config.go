// Package config builds the configuration for the application
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"time"
)

// Config represents all the configuration values for the application
type Config struct {
	Env                      string        `env:"ENV" envDefault:"local"`
	Port                     string        `env:"PORT" envDefault:"8080"`
	ReadTimeout              time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout             time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
	IdleTimeout              time.Duration `env:"IDLE_TIMEOUT" envDefault:"5s"`
	LogLevel                 string        `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL              string        `env:"DB_URL,required"`
	MaxOpenDBConns           int           `env:"MAX_OPEN_DB_CONNS" envDefault:"25"`
	MaxIdleDBConns           int           `env:"MAX_IDLE_DB_CONNS" envDefault:"25"`
	DBConnMaxLifetime        time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
	DBConnMaxIdleTime        time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"1m"`
	HealthCheckPeriod        time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"1m"`
	JWTSecret                string        `env:"JWT_SECRET,required"`
	AccessTokenTTL           time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL          time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`
	PasswordResetTokenTTL    time.Duration `env:"PASSWORD_RESET_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenCookieName   string        `env:"REFRESH_TOKEN_COOKIE_NAME" envDefault:"rt"`
	RefreshTokenCookieDomain string        `env:"REFRESH_TOKEN_COOKIE_DOMAIN" envDefault:"localhost"`
	CORSOrigins              []string      `env:"CORS_ORIGINS" envDefault:"http://localhost:8080"`
	EmailProvider            string        `env:"EMAIL_PROVIDER" envDefault:"logmail"`
	AWSAccessKeyID           string        `env:"AWS_ACCESS_KEY_ID" envDefault:""`
	AWSSecretAccessKey       string        `env:"AWS_SECRET_ACCESS_KEY" envDefault:""`
	AWSRegion                string        `env:"AWS_REGION" envDefault:"us-east-1"`
	FromAddress              string        `env:"FROM_ADDRESS" envDefault:"noreply@localhost"`
	FromName                 string        `env:"FROM_NAME" envDefault:"critiquefi"`
	BaseURL                  string        `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

// Load gets the configuration values from the application environment
func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}

// SlogLevel converts the logmail level config string to a slog.Level
func (c *Config) SlogLevel() slog.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewLogger creates a new logger for the application
func (c *Config) NewLogger() *slog.Logger {
	level := c.SlogLevel()

	var handler slog.Handler
	if c.Env == "local" {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource: true,
			Level:     level,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
	}

	return slog.New(handler)
}
