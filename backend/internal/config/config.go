package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	App      AppConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret        string
	ExpiryMinutes int
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string
	LogLevel    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Server config
	serverPort := getEnv("SERVER_PORT", "8080")
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "60"))

	// Database config
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "papo_reto")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	// Redis config
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	// JWT config
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtExpiryMinutes, _ := strconv.Atoi(getEnv("JWT_EXPIRY_MINUTES", "60"))

	// App config
	environment := getEnv("APP_ENV", "development")
	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		Server: ServerConfig{
			Port:         serverPort,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
		Redis: RedisConfig{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPassword,
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:        jwtSecret,
			ExpiryMinutes: jwtExpiryMinutes,
		},
		App: AppConfig{
			Environment: environment,
			LogLevel:    logLevel,
		},
	}, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr returns the Redis connection string
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
