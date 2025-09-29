package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"github.com/multitask-platform/backend/shared/config"
	"github.com/multitask-platform/backend/shared/logger"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Extract token from Authorization header
		authHeader := request.Headers["Authorization"]
		if authHeader == "" {
			authHeader = request.Headers["authorization"] // case insensitive
		}

		if authHeader == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"missing authorization header"}`,
			}, nil
		}

		// Check for Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"invalid authorization format"}`,
			}, nil
		}

		tokenString := tokenParts[1]

		// Parse and validate JWT token
		cfg := config.Get()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			logger.WarnCtx(ctx, "Invalid JWT token", 
				zap.Error(err),
				zap.String("token_preview", tokenString[:min(len(tokenString), 20)]+"..."),
			)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"invalid or expired token"}`,
			}, nil
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"invalid token claims"}`,
			}, nil
		}

		// Extract user information
		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		roles, _ := claims["roles"].([]interface{})

		if userID == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"missing user ID in token"}`,
			}, nil
		}

		// Add user context
		ctx = logger.WithUserID(ctx, userID)
		ctx = WithUserClaims(ctx, &UserClaims{
			UserID: userID,
			Email:  email,
			Roles:  convertRoles(roles),
		})

		return next(ctx, request)
	}
}

// OptionalAuthMiddleware allows both authenticated and anonymous requests
func OptionalAuthMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		authHeader := request.Headers["Authorization"]
		if authHeader == "" {
			authHeader = request.Headers["authorization"]
		}

		if authHeader != "" {
			// Try to authenticate
			return AuthMiddleware(next)(ctx, request)
		}

		// Continue without authentication
		return next(ctx, request)
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Handle preflight requests
		if request.HTTPMethod == "OPTIONS" {
			cfg := config.Get()
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"Access-Control-Allow-Origin":      cfg.CORSOrigin,
					"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-Correlation-ID",
					"Access-Control-Allow-Credentials": "true",
					"Access-Control-Max-Age":          "3600",
				},
			}, nil
		}

		// Process the request
		response, err := next(ctx, request)
		if err != nil {
			return response, err
		}

		// Add CORS headers to response
		if response.Headers == nil {
			response.Headers = make(map[string]string)
		}

		cfg := config.Get()
		response.Headers["Access-Control-Allow-Origin"] = cfg.CORSOrigin
		response.Headers["Access-Control-Allow-Credentials"] = "true"

		return response, nil
	}
}

// RequestLoggingMiddleware logs request details
func RequestLoggingMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		startTime := time.Now()

		// Add correlation ID to context
		correlationID := request.Headers["X-Correlation-ID"]
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		ctx = logger.WithCorrelationID(ctx, correlationID)

		// Add request ID to context
		requestID := request.RequestContext.RequestID
		ctx = logger.WithRequestID(ctx, requestID)

		logger.InfoCtx(ctx, "HTTP request started",
			zap.String("method", request.HTTPMethod),
			zap.String("path", request.Path),
			zap.String("user_agent", request.Headers["User-Agent"]),
			zap.String("source_ip", request.RequestContext.Identity.SourceIP),
		)

		// Process the request
		response, err := next(ctx, request)

		duration := time.Since(startTime)
		statusCode := http.StatusInternalServerError
		if response.StatusCode != 0 {
			statusCode = response.StatusCode
		}

		// Log response
		logger.LogRequest(ctx, request.HTTPMethod, request.Path, statusCode, duration)

		// Add correlation ID to response headers
		if response.Headers == nil {
			response.Headers = make(map[string]string)
		}
		response.Headers["X-Correlation-ID"] = correlationID

		return response, err
	}
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// TODO: Implement rate limiting using Redis or DynamoDB
		// For now, just pass through
		return next(ctx, request)
	}
}

// ValidationMiddleware validates request payload
func ValidationMiddleware(validator func(request events.APIGatewayProxyRequest) error) func(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			if err := validator(request); err != nil {
				logger.WarnCtx(ctx, "Request validation failed", zap.Error(err))
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: fmt.Sprintf(`{"error":"validation failed","details":"%s"}`, err.Error()),
				}, nil
			}

			return next(ctx, request)
		}
	}
}

// Chain combines multiple middleware functions
func Chain(middlewares ...func(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(next func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// UserClaims represents JWT user claims
type UserClaims struct {
	UserID string   `json:"sub"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

// HasRole checks if user has a specific role
func (u *UserClaims) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Context keys
type contextKeyType string

const userClaimsKey contextKeyType = "user_claims"

// WithUserClaims adds user claims to context
func WithUserClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// GetUserClaims retrieves user claims from context
func GetUserClaims(ctx context.Context) *UserClaims {
	if claims, ok := ctx.Value(userClaimsKey).(*UserClaims); ok {
		return claims
	}
	return nil
}

// Helper functions

func convertRoles(roles []interface{}) []string {
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		if str, ok := role.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}