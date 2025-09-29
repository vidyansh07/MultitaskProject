package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	sugar        *zap.SugaredLogger
)

// Initialize sets up the global logger
func Initialize(level zap.AtomicLevel, isDevelopment bool) error {
	var config zap.Config

	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.Development = true
		config.Level = level
	} else {
		config = zap.NewProductionConfig()
		config.Level = level
		// Add timestamp format for production
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339))
		}
	}

	// Add service name and correlation ID support
	config.InitialFields = map[string]interface{}{
		"service": os.Getenv("SERVICE_NAME"),
		"stage":   os.Getenv("STAGE"),
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1), // Skip wrapper functions in stack trace
	)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	globalLogger = logger
	sugar = logger.Sugar()

	return nil
}

// WithContext adds context information to the logger
func WithContext(ctx context.Context) *zap.Logger {
	if globalLogger == nil {
		panic("logger not initialized - call logger.Initialize() first")
	}

	logger := globalLogger

	// Add correlation ID if present in context
	if correlationID := getCorrelationID(ctx); correlationID != "" {
		logger = logger.With(zap.String("correlation_id", correlationID))
	}

	// Add user ID if present in context
	if userID := getUserID(ctx); userID != "" {
		logger = logger.With(zap.String("user_id", userID))
	}

	// Add request ID if present in context
	if requestID := getRequestID(ctx); requestID != "" {
		logger = logger.With(zap.String("request_id", requestID))
	}

	return logger
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Debug(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Fatal(msg, fields...)
}

// InfoCtx logs an info message with context
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Info(msg, fields...)
}

// DebugCtx logs a debug message with context
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Debug(msg, fields...)
}

// WarnCtx logs a warning message with context
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Warn(msg, fields...)
}

// ErrorCtx logs an error message with context
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Error(msg, fields...)
}

// LogError logs an error with full details
func LogError(ctx context.Context, err error, msg string, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	WithContext(ctx).Error(msg, allFields...)
}

// LogRequest logs HTTP request details
func LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, fields ...zap.Field) {
	allFields := append(fields,
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	)
	WithContext(ctx).Info("HTTP request completed", allFields...)
}

// LogDatabaseOperation logs database operations
func LogDatabaseOperation(ctx context.Context, operation, table string, duration time.Duration, err error, fields ...zap.Field) {
	allFields := append(fields,
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Duration("duration", duration),
	)

	if err != nil {
		allFields = append(allFields, zap.Error(err))
		WithContext(ctx).Error("Database operation failed", allFields...)
	} else {
		WithContext(ctx).Debug("Database operation completed", allFields...)
	}
}

// LogWebSocketEvent logs WebSocket events
func LogWebSocketEvent(ctx context.Context, event string, roomID string, fields ...zap.Field) {
	allFields := append(fields,
		zap.String("event", event),
		zap.String("room_id", roomID),
	)
	WithContext(ctx).Info("WebSocket event", allFields...)
}

// LogEventBridge logs EventBridge events
func LogEventBridge(ctx context.Context, eventType string, source string, err error, fields ...zap.Field) {
	allFields := append(fields,
		zap.String("event_type", eventType),
		zap.String("source", source),
	)

	if err != nil {
		allFields = append(allFields, zap.Error(err))
		WithContext(ctx).Error("EventBridge event failed", allFields...)
	} else {
		WithContext(ctx).Info("EventBridge event published", allFields...)
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// Context key types for type safety
type contextKeyType string

const (
	correlationIDKey contextKeyType = "correlation_id"
	userIDKey        contextKeyType = "user_id"
	requestIDKey     contextKeyType = "request_id"
)

// Context helper functions
func getCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}

func getUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

func getRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithCorrelationID adds correlation ID to context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}