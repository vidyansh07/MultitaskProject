package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/multitask-platform/backend/services/auth-svc/internal/handlers"
	"github.com/multitask-platform/backend/shared/config"
	"github.com/multitask-platform/backend/shared/logger"
	"github.com/multitask-platform/backend/shared/middleware"
	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	if err := logger.Initialize(cfg.GetLogLevel(), cfg.IsDevelopment()); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Configuration validation failed", zap.Error(err))
	}

	logger.Info("Starting Auth Service",
		zap.String("service", cfg.ServiceName),
		zap.String("stage", cfg.Stage),
		zap.String("region", cfg.Region),
	)

	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers()

	// Create router with middleware
	router := createRouter(authHandlers)

	// Start Lambda function
	lambda.Start(router)
}

// createRouter sets up the HTTP routing with middleware
func createRouter(authHandlers *handlers.AuthHandlers) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return middleware.Chain(
		middleware.CORSMiddleware,
		middleware.RequestLoggingMiddleware,
		middleware.RateLimitMiddleware,
	)(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Parse the path to determine the route
		path := strings.TrimPrefix(request.Path, "/v1/auth")
		method := request.HTTPMethod

		logger.DebugCtx(ctx, "Routing request",
			zap.String("method", method),
			zap.String("path", path),
		)

		// Route to appropriate handler
		switch {
		// Authentication endpoints
		case path == "/login" && method == "POST":
			return authHandlers.Login(ctx, request)
		case path == "/logout" && method == "POST":
			return middleware.AuthMiddleware(authHandlers.Logout)(ctx, request)
		case path == "/register" && method == "POST":
			return authHandlers.Register(ctx, request)
		case path == "/refresh" && method == "POST":
			return authHandlers.RefreshToken(ctx, request)

		// Password management
		case path == "/forgot-password" && method == "POST":
			return authHandlers.ForgotPassword(ctx, request)
		case path == "/reset-password" && method == "POST":
			return authHandlers.ResetPassword(ctx, request)
		case path == "/change-password" && method == "POST":
			return middleware.AuthMiddleware(authHandlers.ChangePassword)(ctx, request)

		// Email verification
		case path == "/verify-email" && method == "POST":
			return authHandlers.VerifyEmail(ctx, request)
		case path == "/resend-verification" && method == "POST":
			return authHandlers.ResendVerification(ctx, request)

		// User info
		case path == "/me" && method == "GET":
			return middleware.AuthMiddleware(authHandlers.GetCurrentUser)(ctx, request)

		// Anonymous session management
		case path == "/anonymous" && method == "POST":
			return authHandlers.CreateAnonymousSession(ctx, request)

		// Health check
		case path == "/health" && method == "GET":
			return healthCheck(ctx, request)

		// Session management
		case strings.HasPrefix(path, "/sessions") && method == "GET":
			return middleware.AuthMiddleware(authHandlers.GetUserSessions)(ctx, request)
		case strings.HasPrefix(path, "/sessions/") && method == "DELETE":
			return middleware.AuthMiddleware(authHandlers.RevokeSession)(ctx, request)

		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"endpoint not found"}`,
			}, nil
		}
	})
}

// healthCheck returns service health status
func healthCheck(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg := config.Get()

	health := map[string]interface{}{
		"service":   "auth-svc",
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"stage":     cfg.Stage,
		"region":    cfg.Region,
		"uptime":    time.Since(time.Now()).String(), // This would be tracked in a real app
	}

	body, err := json.Marshal(health)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to marshal health response", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"internal server error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}