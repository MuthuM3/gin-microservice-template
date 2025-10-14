package config

import (
	"fmt"
	"time"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig    `yaml:"database"`
	JWT         JWTConfig         `yaml:"jwt"`
	Logger      LoggerConfig      `yaml:"logger"`
	RateLimit   RateLimitConfig   `yaml:"rate_Limit"`
	CORS        CORSConfig        `yaml:"cors"`
	Redis       RedisConfig       `yaml:"radis"`
	Cache       CacheConfig       `yaml:"cache"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	Security    SecurityConfig    `yaml:"security"`
	Performance PerformanceConfig `yaml:"performance"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string        `yaml:"host" env:"SERVER_HOST" default:"localhost"`
	Port         int           `yaml:"port" env:"PORT" default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" default:"15s"`
	WriteTimeout time.Duration `yaml:"write_timeout" default:"15s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" default:"60s"`
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
	Window           time.Duration `yaml:"window" env:"RATE_LIMIT_WINDOW" default:"1m"`
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
	Host     string `yaml:"host" env:"REDIS_HOST" default:"localhost"`
	Port     int    `yaml:"port" env:"REDIS_PORT" default:"6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" default:""`
	Database string `yaml:"database" env:"REDIS_DATABASE" default:"0"`
}

// GetConnectionString return the database connection string
func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
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

// CacheConfig hold cache related configuration
type CacheConfig struct {
	Enabled    bool          `yaml:"enabled" env:"CACHE_ENABLED" default:"true"`
	DefaultTTL time.Duration `yaml:"default_ttl" default:"1h"`
	LongTTL    time.Duration `yaml:"long_ttl" default:"24h"`
	ShortTTL   time.Duration `yaml:"short_ttl" default:"5m"`
	SessionTTL time.Duration `yaml:"session_ttl" default:"30m"`
	StateTTL   time.Duration `yaml:"state_ttl" default:"15m"`
	MaxMemory  string        `yaml:"max_memory" default:"256md"`
	KeyPrefix  string        `yaml:"key_prefix" default:"go-microservice-api"`
}

// MetricsConfig holds metrics-related configuration
type MetricsConfig struct {
	Enabled            bool          `yaml:"enabled" env:"METRICS_ENABLED" default:"true"`
	CollectionInterval time.Duration `yaml:"collection_interval" default:"30s"`
	RetentionPeriod    time.Duration `yaml:"retention_period" default:"24h"`
	ExportPrometheus   bool          `yaml:"export_prometheus" default:"false"`
	PrometheusPath     string        `yaml:"prometheus_path" default:"/matrics"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	PasswordMinLength       int           `yaml:"password_min_length" default:"8"`
	PasswordRequiredUpper   string        `yaml:"password_required_upper" default:"true"`
	PasswordRequiredLower   string        `yaml:"password_required_lower" default:"true"`
	PasswordRequiredDigital string        `yaml:"password_required_digital" default:"true"`
	PasswordRequiredSymbol  string        `yaml:"password_required_symbol" default:"false"`
	MaxLoginAttempts        int           `yaml:"max_login_attempts" default:"5"`
	LoginLogoutDuration     time.Duration `yaml:"login_logout_duration" default:"15m"`
	SessionTimeout          time.Duration `yaml:"session_timeout" default:"24h"`
	CSRFEnabled             bool          `yaml:"csrf_enabled" default:"true"`
	CSRFTokenLength         int           `yaml:"csrf_token_length" default:"32"`
	SecureHeaders           bool          `yaml:"secure_header" default:"true"`
	ContentTypeValidation   bool          `yaml:"content_type_validation" default:"true"`
	MaxRequestSize          int64         `yaml:"max_request_size" default:"10485760"`
}

// PerformanceConfig holds performance-related configuration
type PerformanceConfig struct {
	EnableCompression     bool          `yaml:"enable_compression" default:"true"`
	CompressionLevel      int           `yaml:"compression_level" default:"6"`
	CompressionMinLength  int           `yaml:"compression_min_length" default:"1024"`
	EnableCaching         bool          `yaml:"enable_caching" default:"true"`
	CacheControlMaxAge    int           `yaml:"cache_control_max_age" default:"3600"`
	EnableETag            bool          `yaml:"enable_etag" default:"true"`
	MaxConcurrentRequests int           `yaml:"max_concurrent_requests" default:"1000"`
	RequestTimeout        time.Duration `yaml:"request_timeout" default:"30s"`
	KeepAliveTimeout      time.Duration `yaml:"keep_alive_timeout" default:"60s"`
	EnableProfiling       bool          `yaml:"enable_profiling" default:"false"`
	ProfilingPath         string        `yaml:"profiling_path" default:"/debug/pprof"`
}

// IsSecure returns true if security features are enabled
func (c *SecurityConfig) IsSecure() bool {
	return c.SecureHeaders && c.ContentTypeValidation
}

// GetMaxRequestSizeMB returns max request size in MB
func (c *SecurityConfig) GetMaxRequestSizeMB() int64 {
	return c.MaxRequestSize / (1024 * 1024)
}

// IsCompressionEnabled returns true if compression should be enabled
func (c *PerformanceConfig) IsCompressionEnabled() bool {
	return c.EnableCompression
}

// IsCachingEnabled returns true if caching should be enabled
func (c *PerformanceConfig) IsCachingEnabled() bool {
	return c.EnableCaching
}

// IsProfilingEnabled returns true if profiling should be enabled
func (c *PerformanceConfig) IsProfilingEnabled() bool {
	return c.EnableProfiling
}
