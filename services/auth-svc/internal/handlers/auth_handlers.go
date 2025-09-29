package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/multitask-platform/backend/services/auth-svc/internal/models"
	"github.com/multitask-platform/backend/services/auth-svc/internal/services"
	"github.com/multitask-platform/backend/shared/logger"
	"github.com/multitask-platform/backend/shared/middleware"
)

// AuthHandlers contains all auth-related HTTP handlers
type AuthHandlers struct {
	authService *services.AuthService
	validator   *validator.Validate
}

// NewAuthHandlers creates a new instance of AuthHandlers
func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{
		authService: services.NewAuthService(),
		validator:   validator.New(),
	}
}

// Login handles user login requests
func (h *AuthHandlers) Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing login request")

	// Parse request body
	var loginReq models.LoginRequest
	if err := json.Unmarshal([]byte(request.Body), &loginReq); err != nil {
		logger.WarnCtx(ctx, "Invalid login request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&loginReq); err != nil {
		logger.WarnCtx(ctx, "Login request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Authenticate user
	authResponse, err := h.authService.Login(ctx, &loginReq)
	if err != nil {
		logger.WarnCtx(ctx, "Login failed", zap.Error(err))
		
		switch err {
		case services.ErrInvalidCredentials:
			return h.errorResponse(http.StatusUnauthorized, "invalid credentials"), nil
		case services.ErrUserNotVerified:
			return h.errorResponse(http.StatusForbidden, "email not verified"), nil
		case services.ErrUserDisabled:
			return h.errorResponse(http.StatusForbidden, "account disabled"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "authentication failed"), nil
		}
	}

	logger.InfoCtx(ctx, "Login successful", zap.String("user_id", authResponse.User.ID))

	return h.successResponse(http.StatusOK, authResponse), nil
}

// Register handles user registration requests
func (h *AuthHandlers) Register(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing registration request")

	// Parse request body
	var registerReq models.RegisterRequest
	if err := json.Unmarshal([]byte(request.Body), &registerReq); err != nil {
		logger.WarnCtx(ctx, "Invalid registration request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&registerReq); err != nil {
		logger.WarnCtx(ctx, "Registration request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Register user
	user, err := h.authService.Register(ctx, &registerReq)
	if err != nil {
		logger.WarnCtx(ctx, "Registration failed", zap.Error(err))
		
		switch err {
		case services.ErrUserAlreadyExists:
			return h.errorResponse(http.StatusConflict, "user already exists"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "registration failed"), nil
		}
	}

	logger.InfoCtx(ctx, "Registration successful", zap.String("user_id", user.ID))

	response := map[string]interface{}{
		"message": "registration successful, please verify your email",
		"user_id": user.ID,
	}

	return h.successResponse(http.StatusCreated, response), nil
}

// Logout handles user logout requests
func (h *AuthHandlers) Logout(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing logout request")

	// Get user from context
	userClaims := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		return h.errorResponse(http.StatusUnauthorized, "unauthorized"), nil
	}

	// Parse request to get session ID (optional - can logout specific session or all)
	var logoutReq models.LogoutRequest
	if request.Body != "" {
		if err := json.Unmarshal([]byte(request.Body), &logoutReq); err != nil {
			logger.WarnCtx(ctx, "Invalid logout request body", zap.Error(err))
			// Continue with logout all sessions if body is invalid
		}
	}

	// Logout
	err := h.authService.Logout(ctx, userClaims.UserID, logoutReq.SessionID)
	if err != nil {
		logger.ErrorCtx(ctx, "Logout failed", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, "logout failed"), nil
	}

	logger.InfoCtx(ctx, "Logout successful", zap.String("user_id", userClaims.UserID))

	return h.successResponse(http.StatusOK, map[string]string{"message": "logout successful"}), nil
}

// RefreshToken handles token refresh requests
func (h *AuthHandlers) RefreshToken(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing token refresh request")

	// Parse request body
	var refreshReq models.RefreshTokenRequest
	if err := json.Unmarshal([]byte(request.Body), &refreshReq); err != nil {
		logger.WarnCtx(ctx, "Invalid refresh token request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&refreshReq); err != nil {
		logger.WarnCtx(ctx, "Refresh token request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Refresh token
	authResponse, err := h.authService.RefreshToken(ctx, refreshReq.RefreshToken)
	if err != nil {
		logger.WarnCtx(ctx, "Token refresh failed", zap.Error(err))
		
		switch err {
		case services.ErrInvalidToken:
			return h.errorResponse(http.StatusUnauthorized, "invalid refresh token"), nil
		case services.ErrTokenExpired:
			return h.errorResponse(http.StatusUnauthorized, "refresh token expired"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "token refresh failed"), nil
		}
	}

	logger.InfoCtx(ctx, "Token refresh successful", zap.String("user_id", authResponse.User.ID))

	return h.successResponse(http.StatusOK, authResponse), nil
}

// ForgotPassword handles forgot password requests
func (h *AuthHandlers) ForgotPassword(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing forgot password request")

	// Parse request body
	var forgotReq models.ForgotPasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &forgotReq); err != nil {
		logger.WarnCtx(ctx, "Invalid forgot password request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&forgotReq); err != nil {
		logger.WarnCtx(ctx, "Forgot password request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Process forgot password
	err := h.authService.ForgotPassword(ctx, forgotReq.Email)
	if err != nil {
		logger.ErrorCtx(ctx, "Forgot password processing failed", zap.Error(err))
		// Don't reveal if user exists or not
	}

	// Always return success to prevent email enumeration
	response := map[string]string{
		"message": "if the email exists, a password reset link has been sent",
	}

	return h.successResponse(http.StatusOK, response), nil
}

// ResetPassword handles password reset requests
func (h *AuthHandlers) ResetPassword(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing password reset request")

	// Parse request body
	var resetReq models.ResetPasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &resetReq); err != nil {
		logger.WarnCtx(ctx, "Invalid reset password request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&resetReq); err != nil {
		logger.WarnCtx(ctx, "Reset password request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Reset password
	err := h.authService.ResetPassword(ctx, resetReq.Token, resetReq.NewPassword)
	if err != nil {
		logger.WarnCtx(ctx, "Password reset failed", zap.Error(err))
		
		switch err {
		case services.ErrInvalidToken:
			return h.errorResponse(http.StatusBadRequest, "invalid or expired reset token"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "password reset failed"), nil
		}
	}

	logger.InfoCtx(ctx, "Password reset successful")

	response := map[string]string{
		"message": "password reset successful",
	}

	return h.successResponse(http.StatusOK, response), nil
}

// ChangePassword handles password change requests
func (h *AuthHandlers) ChangePassword(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing password change request")

	// Get user from context
	userClaims := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		return h.errorResponse(http.StatusUnauthorized, "unauthorized"), nil
	}

	// Parse request body
	var changeReq models.ChangePasswordRequest
	if err := json.Unmarshal([]byte(request.Body), &changeReq); err != nil {
		logger.WarnCtx(ctx, "Invalid change password request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&changeReq); err != nil {
		logger.WarnCtx(ctx, "Change password request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Change password
	err := h.authService.ChangePassword(ctx, userClaims.UserID, changeReq.CurrentPassword, changeReq.NewPassword)
	if err != nil {
		logger.WarnCtx(ctx, "Password change failed", zap.Error(err))
		
		switch err {
		case services.ErrInvalidCredentials:
			return h.errorResponse(http.StatusBadRequest, "current password is incorrect"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "password change failed"), nil
		}
	}

	logger.InfoCtx(ctx, "Password change successful", zap.String("user_id", userClaims.UserID))

	response := map[string]string{
		"message": "password changed successfully",
	}

	return h.successResponse(http.StatusOK, response), nil
}

// VerifyEmail handles email verification requests
func (h *AuthHandlers) VerifyEmail(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing email verification request")

	// Parse request body
	var verifyReq models.VerifyEmailRequest
	if err := json.Unmarshal([]byte(request.Body), &verifyReq); err != nil {
		logger.WarnCtx(ctx, "Invalid verify email request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&verifyReq); err != nil {
		logger.WarnCtx(ctx, "Verify email request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Verify email
	err := h.authService.VerifyEmail(ctx, verifyReq.Token)
	if err != nil {
		logger.WarnCtx(ctx, "Email verification failed", zap.Error(err))
		
		switch err {
		case services.ErrInvalidToken:
			return h.errorResponse(http.StatusBadRequest, "invalid or expired verification token"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "email verification failed"), nil
		}
	}

	logger.InfoCtx(ctx, "Email verification successful")

	response := map[string]string{
		"message": "email verified successfully",
	}

	return h.successResponse(http.StatusOK, response), nil
}

// ResendVerification handles resend verification email requests
func (h *AuthHandlers) ResendVerification(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing resend verification request")

	// Parse request body
	var resendReq models.ResendVerificationRequest
	if err := json.Unmarshal([]byte(request.Body), &resendReq); err != nil {
		logger.WarnCtx(ctx, "Invalid resend verification request body", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "invalid request body"), nil
	}

	// Validate request
	if err := h.validator.Struct(&resendReq); err != nil {
		logger.WarnCtx(ctx, "Resend verification request validation failed", zap.Error(err))
		return h.errorResponse(http.StatusBadRequest, "validation failed: "+err.Error()), nil
	}

	// Resend verification
	err := h.authService.ResendVerification(ctx, resendReq.Email)
	if err != nil {
		logger.ErrorCtx(ctx, "Resend verification failed", zap.Error(err))
		// Don't reveal if user exists or not
	}

	// Always return success to prevent email enumeration
	response := map[string]string{
		"message": "if the email exists and is not verified, a verification email has been sent",
	}

	return h.successResponse(http.StatusOK, response), nil
}

// GetCurrentUser returns current user information
func (h *AuthHandlers) GetCurrentUser(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing get current user request")

	// Get user from context
	userClaims := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		return h.errorResponse(http.StatusUnauthorized, "unauthorized"), nil
	}

	// Get user details
	user, err := h.authService.GetUser(ctx, userClaims.UserID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, "failed to get user information"), nil
	}

	return h.successResponse(http.StatusOK, user), nil
}

// CreateAnonymousSession creates an anonymous session
func (h *AuthHandlers) CreateAnonymousSession(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing create anonymous session request")

	// Create anonymous session
	session, err := h.authService.CreateAnonymousSession(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to create anonymous session", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, "failed to create anonymous session"), nil
	}

	logger.InfoCtx(ctx, "Anonymous session created", zap.String("session_id", session.ID))

	return h.successResponse(http.StatusCreated, session), nil
}

// GetUserSessions returns user's active sessions
func (h *AuthHandlers) GetUserSessions(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing get user sessions request")

	// Get user from context
	userClaims := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		return h.errorResponse(http.StatusUnauthorized, "unauthorized"), nil
	}

	// Get user sessions
	sessions, err := h.authService.GetUserSessions(ctx, userClaims.UserID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user sessions", zap.Error(err))
		return h.errorResponse(http.StatusInternalServerError, "failed to get sessions"), nil
	}

	return h.successResponse(http.StatusOK, map[string]interface{}{"sessions": sessions}), nil
}

// RevokeSession revokes a specific session
func (h *AuthHandlers) RevokeSession(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.InfoCtx(ctx, "Processing revoke session request")

	// Get user from context
	userClaims := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		return h.errorResponse(http.StatusUnauthorized, "unauthorized"), nil
	}

	// Extract session ID from path
	pathParts := strings.Split(strings.TrimPrefix(request.Path, "/v1/auth/sessions/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		return h.errorResponse(http.StatusBadRequest, "session ID required"), nil
	}
	sessionID := pathParts[0]

	// Revoke session
	err := h.authService.RevokeSession(ctx, userClaims.UserID, sessionID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to revoke session", zap.Error(err))
		
		switch err {
		case services.ErrSessionNotFound:
			return h.errorResponse(http.StatusNotFound, "session not found"), nil
		default:
			return h.errorResponse(http.StatusInternalServerError, "failed to revoke session"), nil
		}
	}

	logger.InfoCtx(ctx, "Session revoked successfully", 
		zap.String("user_id", userClaims.UserID),
		zap.String("session_id", sessionID),
	)

	return h.successResponse(http.StatusOK, map[string]string{"message": "session revoked successfully"}), nil
}

// Helper methods

func (h *AuthHandlers) successResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func (h *AuthHandlers) errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}