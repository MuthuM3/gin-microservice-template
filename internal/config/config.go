package config

import (
	"fmt"
	"time"
)

type Config struct {
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string        `yaml:"host" env:"SERVER_HOST" default:"localhost"`
	Port         int           `yaml:"port" env:"PORT" default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" default:"15s"`
	WriteTimeout time.Duration `yaml:"write_timeout" default:"15s"`
	idleTimeout  time.Duration `yaml:"idle_timeout" default:"60s"`
	Environment  string        `yaml:"environment" env:"APP_ENV" default:"development"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `yaml:"host" env:"DB_HOST" default:"localhost"`
	Port            int           `yaml:"host" env:"DB_PORT" default:"5432"`
	User            string        `yaml:"user" env:"DB_USER" default:"postgres"`
	Password        string        `yaml:"password" env:"DB_PASSWORD" default:"root"`
	Database        string        `yaml:"database" env:"DB_NAME" default:"go-microservice"`
	SSLMode         string        `yaml:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int           `yaml:"max_open_conns" default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" default:"5"`
	ConnMaxLifeTime time.Duration `yaml:"conn_max_lifetime" default:"5m"`
}

// JWTConfig holds the jwt-related configuration
type JWTConfig struct {
	Secret     string        `yaml:"secret" env:"JWT_SECRET"`
	Expiration time.Duration `yaml:"expiration" default:"24h"`
	Issuer     string        `yaml:"issuer" default:"microservice-api"`
}

// Logger config holds logger related configuration
type LoggerConfig struct {
	Level      string `yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format     string `yaml:"format" default:"json"`
	OutputPath string `yaml:"output_path" default:"stdout"`
}

// RateLimitConfig holds rate limit configuration
type RateLimitConfig struct {
	Enabled          bool          `yaml:"enabled" env:"RATE_LIMIT_ENABLED" default:"true"`
	RequestPerWindow int           `yaml:"request_per_window" env:"RATE_LIMIT_REQUEST_PER_WINDOW" default:"100"`
	window           time.Duration `yaml:"window" env:"RATE_LIMIT_WINDOW" default:"1m"`
	UserBased        bool          `yaml:"user_based" default:"false"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Enabled        bool     `yaml:"enabled" default:"true"`
	AllowedOrigins []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	MaxAge         int      `yaml:"max_age" default:"86400"`
}

// RedisConfig holds redis related configuration
type RedisConfig struct {
	Host      string `yaml:"host" env:"REDIS_HOST" default:"localhost"`
	Port      int    `yaml:"port" env:"REDIS_PORT" default:"6379"`
	Password  string `yaml:"password" env:"REDIS_PASSWORD" default:""`
	Databasee string `yaml:"database" env:"REDIS_DATABASE" default:"0"`
}

// GetConnectionString return the database connection string
func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// IsDevelopment returns true if the server is running in development mode
func (c *ServerConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the server is running in production mode
func (c *ServerConfig) IsProduction() bool {
	return c.Environment == "production"
}

// GetAddress returns the server address
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetConnectionString returns the Redis connection string
func (c *RedisConfig) GetConnectionString() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetAddress returns the Redis address
func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
