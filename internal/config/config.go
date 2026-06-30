package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const minSecretLength = 32

type Config struct {
	Server    ServerConfig
	DB        DbConfig
	JWT       JWTConfig
	Security  SecurityConfig
	Admin     AdminConfig
	RateLimit RateLimitConfig
}

type AdminConfig struct {
	Email    string
	Password string
}

type SecurityConfig struct {
	SessionSecret string
	CSRFSecret    string
}

type DbConfig struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	Name                   string
	TimeZone               string
	MinIdleConns           int
	MaxOpenConns           int
	DefaultTimeoutDuration int
	ConnMaxLifetime        int
}

type ServerConfig struct {
	Port            string
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type JWTConfig struct {
	SecretKey        string
	ExpiresIn        string
	RefreshExpiresIn string
	Issuer           string
}

type RateLimitConfig struct {
	Limit int
	Reset time.Duration
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return d
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// DSN returns the key-value connection string used by GORM.
// statement_timeout and idle_in_transaction_session_timeout mirror
// DefaultTimeoutDuration so the Postgres side enforces the same ceiling
// as the Go pool's SetConnMaxIdleTime.
func (d DbConfig) DSN() string {
	timeoutMs := d.DefaultTimeoutDuration * 1000
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s"+
			" options='-c statement_timeout=%d -c idle_in_transaction_session_timeout=%d'",
		d.Host, d.User, d.Password, d.Name, d.Port, d.TimeZone, timeoutMs, timeoutMs,
	)
}

// URL returns the postgres:// connection URL used by golang-migrate.
func (d DbConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	if len(getEnv("JWT_SECRET_KEY", "secret")) < minSecretLength || len(getEnv("SESSION_SECRET", "")) < minSecretLength || len(getEnv("CSRF_SECRET", "")) < minSecretLength {
		return nil, fmt.Errorf("JWT_SECRET_KEY, SESSION_SECRET, and CSRF_SECRET must be at least %d characters long", minSecretLength)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8000"),
			Addr:            getEnv("SERVER_ADDR", "0.0.0.0"),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		DB: DbConfig{
			Host:                   getEnv("DB_HOST", "localhost"),
			Port:                   getEnv("DB_PORT", "5432"),
			User:                   getEnv("DB_USER", "user"),
			Password:               getEnv("DB_PASSWORD", "password"),
			Name:                   getEnv("DB_NAME", "dbname"),
			TimeZone:               getEnv("DB_TIMEZONE", "UTC"),
			MinIdleConns:           getEnvAsInt("DB_MIN_IDLE_CONNS", 2),
			MaxOpenConns:           getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
			DefaultTimeoutDuration: getEnvAsInt("DB_DEFAULT_TIMEOUT_DURATION", 60), // in seconds
			ConnMaxLifetime:        getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600),      // in seconds (default 1 hour)
		},
		JWT: JWTConfig{
			SecretKey:        getEnv("JWT_SECRET_KEY", "secret"),
			ExpiresIn:        getEnv("JWT_EXPIRES_IN", "15m"),
			RefreshExpiresIn: getEnv("JWT_REFRESH_EXPIRES_IN", "168h"),
			Issuer:           getEnv("JWT_ISSUER", "your-issuer"),
		},
		Security: SecurityConfig{
			SessionSecret: getEnv("SESSION_SECRET", ""),
			CSRFSecret:    getEnv("CSRF_SECRET", ""),
		},
		Admin: AdminConfig{
			Email:    getEnv("ADMIN_EMAIL", "admin@example.com"),
			Password: getEnv("ADMIN_PASSWORD", "adminpassword"),
		},
		RateLimit: RateLimitConfig{
			Limit: getEnvAsInt("RATE_LIMIT", 10),
			Reset: getEnvAsDuration("RATE_LIMIT_RESET", time.Minute),
		},
	}
	return cfg, nil
}
