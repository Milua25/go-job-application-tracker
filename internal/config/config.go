package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	DB       DbConfig
	JWT      JWTConfig
	Security SecurityConfig
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
	Port string
	Addr string
}

type JWTConfig struct {
	SecretKey        string
	ExpiresIn        string
	RefreshExpiresIn string
	Issuer           string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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

// // URL returns the postgres:// connection URL required by golang-migrate.
// func (d DbConfig) URL() string {
// 	return fmt.Sprintf(
// 		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		d.User, d.Password, d.Host, d.Port, d.Name,
// 	)
// }

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Addr: getEnv("SERVER_ADDR", "0.0.0.0"),
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
			ConnMaxLifetime:        getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600),       // in seconds (default 1 hour)
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
	}
	return cfg, nil
}
