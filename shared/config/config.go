package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config holds all configuration values
type Config struct {
	// Environment
	Stage  string
	Region string

	// AWS Resources
	DynamoDB struct {
		AuthSessions    string
		AuthAnonymous   string
		Profiles        string
		ProfileAliases  string
		ChatMessages    string
		ChatRooms       string
		Posts           string
		Comments        string
	}

	S3 struct {
		AvatarsBucket     string
		AttachmentsBucket string
		FrontendBucket    string
	}

	Cognito struct {
		UserPoolID string
		ClientID   string
	}

	EventBridge struct {
		BusName string
	}

	// External APIs
	JWTSecret    string
	GeminiAPIKey string
	OpenAIAPIKey string

	// CORS
	CORSOrigin string

	// Service Settings
	ServiceName        string
	LogLevel           string
	CloudFrontDomain   string
	WebSocketEndpoint  string

	// Rate Limiting
	RateLimit struct {
		RequestsPerMinute int
		BurstSize         int
	}

	// Timeouts
	Timeouts struct {
		DatabaseTimeout time.Duration
		HTTPTimeout     time.Duration
		WebSocketWrite  time.Duration
	}
}

var globalConfig *Config

// Load initializes and returns the configuration
func Load() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	config := &Config{
		Stage:  getEnv("STAGE", "dev"),
		Region: getEnv("REGION", "us-east-1"),

		ServiceName:       getEnv("SERVICE_NAME", "unknown"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		CloudFrontDomain:  getEnv("CLOUDFRONT_DOMAIN", ""),
		WebSocketEndpoint: getEnv("WEBSOCKET_API_ENDPOINT", ""),

		JWTSecret:    getEnv("JWT_SECRET", ""),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),

		CORSOrigin: getEnv("CORS_ORIGIN", "*"),
	}

	// DynamoDB table names
	config.DynamoDB.AuthSessions = getEnv("DYNAMODB_TABLE_AUTH_SESSIONS", "")
	config.DynamoDB.AuthAnonymous = getEnv("DYNAMODB_TABLE_AUTH_ANONYMOUS", "")
	config.DynamoDB.Profiles = getEnv("DYNAMODB_TABLE_PROFILES", "")
	config.DynamoDB.ProfileAliases = getEnv("DYNAMODB_TABLE_PROFILE_ALIASES", "")
	config.DynamoDB.ChatMessages = getEnv("DYNAMODB_TABLE_CHAT_MESSAGES", "")
	config.DynamoDB.ChatRooms = getEnv("DYNAMODB_TABLE_CHAT_ROOMS", "")
	config.DynamoDB.Posts = getEnv("DYNAMODB_TABLE_POSTS", "")
	config.DynamoDB.Comments = getEnv("DYNAMODB_TABLE_COMMENTS", "")

	// S3 bucket names
	config.S3.AvatarsBucket = getEnv("S3_BUCKET_AVATARS", "")
	config.S3.AttachmentsBucket = getEnv("S3_BUCKET_ATTACHMENTS", "")
	config.S3.FrontendBucket = getEnv("S3_BUCKET_FRONTEND", "")

	// Cognito
	config.Cognito.UserPoolID = getEnv("COGNITO_USER_POOL_ID", "")
	config.Cognito.ClientID = getEnv("COGNITO_CLIENT_ID", "")

	// EventBridge
	config.EventBridge.BusName = getEnv("EVENTBRIDGE_BUS_NAME", "")

	// Rate limiting
	config.RateLimit.RequestsPerMinute = getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60)
	config.RateLimit.BurstSize = getEnvInt("RATE_LIMIT_BURST_SIZE", 10)

	// Timeouts
	config.Timeouts.DatabaseTimeout = getEnvDuration("DATABASE_TIMEOUT", 5*time.Second)
	config.Timeouts.HTTPTimeout = getEnvDuration("HTTP_TIMEOUT", 30*time.Second)
	config.Timeouts.WebSocketWrite = getEnvDuration("WEBSOCKET_WRITE_TIMEOUT", 10*time.Second)

	globalConfig = config
	return config, nil
}

// Get returns the global configuration (must call Load first)
func Get() *Config {
	if globalConfig == nil {
		panic("config not loaded - call config.Load() first")
	}
	return globalConfig
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.Stage == "prod"
}

// IsDevelopment returns true if running in development
func (c *Config) IsDevelopment() bool {
	return c.Stage == "dev"
}

// GetLogLevel returns the zap log level
func (c *Config) GetLogLevel() zap.AtomicLevel {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn", "warning":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	required := map[string]string{
		"JWT_SECRET":                    c.JWTSecret,
		"DYNAMODB_TABLE_AUTH_SESSIONS": c.DynamoDB.AuthSessions,
		"COGNITO_USER_POOL_ID":         c.Cognito.UserPoolID,
		"COGNITO_CLIENT_ID":            c.Cognito.ClientID,
	}

	for key, value := range required {
		if value == "" {
			return &ValidationError{Field: key, Message: "is required"}
		}
	}

	return nil
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "config validation failed: " + e.Field + " " + e.Message
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}