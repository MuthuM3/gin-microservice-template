package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	// Load .env file first if it exists
	if err := loadDotConfig(".env"); err != nil {
		// Don't fail if .env doesn't exist, just log it
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// set default first
	setDefaults(cfg)

	// Load from file if path is provided
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from env: %w", err)
	}

	// Validate configuration
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// LoadForEnvironment loads configuration based on the environment
func LoadForEnvironment(env string) (*Config, error) {
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = "configs"
	}

	configFile := fmt.Sprintf("%s.yaml", env)
	configPath := filepath.Join(configDir, configFile)

	return Load(configPath)
}


// set Defaults sets default values from the configration
func setDefaults(cfg *Config) {
	// Server defaults
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 8000
	cfg.Server.ReadTimeout = 15 * time.Second
	cfg.Server.WriteTimeout = 15 * time.Second
	cfg.Server.IdleTimeout = 60 * time.Second
	cfg.Server.Environment = "development"

	// Database defaults
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.User = "postgres"
	cfg.Database.Password = "root"
	cfg.Database.Database = "todo"
	cfg.Database.SSLMode = "disable"
	cfg.Database.MaxOpenConns = 25
	cfg.Database.MaxIdleConns = 5
	cfg.Database.ConnMaxLifetime = 5 * time.Minute

	// JWT defaults
	cfg.JWT.Expiration = 24 * time.Hour
	cfg.JWT.Issuer = "todo-api"

	// Logger defaults
	cfg.Logger.Level = "info"
	cfg.Logger.Format = "json"
	cfg.Logger.OutputPath = "stdout"

	// Rate limit defaults
	cfg.RateLimit.Enabled = true
	cfg.RateLimit.RequestsPerWindow = 100
	cfg.RateLimit.Window = time.Minute
	cfg.RateLimit.UserBased = false

	// CORS defaults
	cfg.CORS.Enabled = true
	cfg.CORS.AllowedOrigins = []string{"*"}
	cfg.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	cfg.CORS.AllowedHeaders = []string{"Content-Type", "Authorization"}
	cfg.CORS.MaxAge = 86400

	// Redis defaults
	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = 6379
	cfg.Redis.Password = ""
	cfg.Redis.Database = 0
}

func loadDotConfig(fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skipe empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Spilt on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Set environment variable only if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

func loadFromFile(path string, cfg *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(cfg)
}

func loadFromEnv(cfg *Config) error {
	v := reflect.ValueOf(cfg).Elem()
	return loadStructFromEnv(v, "")
}

func loadStructFromEnv(v reflect.Value, prefix string) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		envTag := fieldType.Tag.Get("env")

		if envTag == "" && field.Kind() == reflect.Struct {
			// Recursively load nested structs
			newPrefix := prefix
			if newPrefix != "" {
				newPrefix += "_"
			}
			newPrefix += strings.ToUpper(fieldType.Name)
			if err := loadStructFromEnv(field, newPrefix); err != nil {
				return err
			}
			continue
		}

		if envTag == "" {
			envValue := os.Getenv(envTag)
			if envValue != "" {
				if err := setFieldValue(field, envValue); err != nil {
					return fmt.Errorf("failed to set %s: %w", envTag, err)
				}
			}
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	// handle time.duration first before checking other types
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(duration))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		field.SetBool(boolValue)
	case reflect.Slice:
		// Handle string slices (comma-sperated)
		if field.Type().Elem().Kind() == reflect.String {
			values := strings.Split(value, ",")
			field.Set(reflect.ValueOf(values))
		}
	}
	return nil
}

// validate validates the configuration
func validate(cfg *Config) error {
	// Validate JWT secret
	if cfg.JWT.Secret == "" {
		// try to generate a default for developement
		if cfg.Server.Environment == "development" {
			cfg.JWT.Secret = "development-seecret-do-not-use-in-production"
		} else {
			return fmt.Errorf("JWT secret is required")
		}
	}

	if cfg.Database.Host == "" || cfg.Database.Database == "" {
		return fmt.Errorf("database host and name are required")
	}

	// Validate server configuration
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	return nil
}
