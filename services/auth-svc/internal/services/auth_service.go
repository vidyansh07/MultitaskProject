package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/multitask-platform/backend/services/auth-svc/internal/models"
	"github.com/multitask-platform/backend/services/auth-svc/internal/repositories"
	"github.com/multitask-platform/backend/shared/config"
	"github.com/multitask-platform/backend/shared/logger"
)

// Common errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotVerified    = errors.New("user not verified")
	ErrUserDisabled       = errors.New("user disabled")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrSessionNotFound    = errors.New("session not found")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	cfg         *config.Config
}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:    repositories.NewDynamoDBUserRepository(),
		sessionRepo: repositories.NewDynamoDBSessionRepository(),
		cfg:         config.Get(),
	}
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	logger.DebugCtx(ctx, "Attempting to login user", zap.String("email", req.Email))

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserDisabled
	}

	// Check if user is verified
	if !user.IsVerified {
		return nil, ErrUserNotVerified
	}

	// Verify password
	err = s.verifyPassword(ctx, user.ID, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Create session
	sessionReq := &models.SessionCreateRequest{
		UserID:    user.ID,
		DeviceID:  req.DeviceID,
		UserAgent: "", // This would be extracted from request headers
		IPAddress: "", // This would be extracted from request context
	}

	session, err := s.createSession(ctx, sessionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login time
	err = s.userRepo.UpdateLastLogin(ctx, user.ID, time.Now().UTC())
	if err != nil {
		logger.WarnCtx(ctx, "Failed to update last login time", zap.Error(err))
		// Don't fail login for this
	}

	// Create response
	response := &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(models.DefaultSessionDuration.Seconds()),
		User:         user.SanitizeUser(),
	}

	logger.InfoCtx(ctx, "User login successful", zap.String("user_id", user.ID))

	return response, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	logger.DebugCtx(ctx, "Attempting to register user", zap.String("email", req.Email))

	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err != repositories.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	userID := uuid.New().String()
	createReq := &models.UserCreateRequest{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Roles:        []string{models.RoleUser},
	}

	user := &models.User{
		ID:         userID,
		Email:      createReq.Email,
		Name:       createReq.Name,
		IsVerified: false, // User needs to verify email
		IsActive:   true,
		Roles:      createReq.Roles,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	err = s.userRepo.CreateUser(ctx, user, createReq.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send verification email
	err = s.sendVerificationEmail(ctx, user)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to send verification email", zap.Error(err))
		// Don't fail registration for this
	}

	logger.InfoCtx(ctx, "User registration successful", zap.String("user_id", user.ID))

	return user, nil
}

// Logout invalidates user sessions
func (s *AuthService) Logout(ctx context.Context, userID, sessionID string) error {
	logger.DebugCtx(ctx, "Attempting to logout user", 
		zap.String("user_id", userID),
		zap.String("session_id", sessionID),
	)

	if sessionID != "" {
		// Logout specific session
		err := s.sessionRepo.DeactivateSession(ctx, sessionID)
		if err != nil {
			return fmt.Errorf("failed to deactivate session: %w", err)
		}
	} else {
		// Logout all sessions
		err := s.sessionRepo.DeactivateUserSessions(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to deactivate user sessions: %w", err)
		}
	}

	logger.InfoCtx(ctx, "User logout successful", zap.String("user_id", userID))

	return nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	logger.DebugCtx(ctx, "Attempting to refresh token")

	// Parse and validate refresh token
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetUser(ctx, claims.UserID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, ErrUserDisabled
	}

	// Verify session exists and is active
	session, err := s.sessionRepo.GetSession(ctx, claims.SessionID)
	if err != nil {
		if err == repositories.ErrSessionNotFound {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if !session.IsActive || session.IsExpired() {
		return nil, ErrTokenExpired
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken(user.ID, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create response
	response := &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(models.DefaultSessionDuration.Seconds()),
		User:         user.SanitizeUser(),
	}

	logger.InfoCtx(ctx, "Token refresh successful", zap.String("user_id", user.ID))

	return response, nil
}

// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	logger.DebugCtx(ctx, "Processing forgot password request", zap.String("email", email))

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			// Don't reveal that user doesn't exist
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Generate reset token
	resetToken, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Store reset token
	err = s.userRepo.CreatePasswordResetToken(ctx, user.ID, resetToken, models.DefaultResetTokenDuration)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send reset email
	err = s.sendPasswordResetEmail(ctx, user, resetToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to send password reset email", zap.Error(err))
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	logger.InfoCtx(ctx, "Password reset initiated", zap.String("user_id", user.ID))

	return nil
}

// ResetPassword resets user password using reset token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	logger.DebugCtx(ctx, "Processing password reset")

	// Verify reset token
	userID, err := s.userRepo.VerifyPasswordResetToken(ctx, token)
	if err != nil {
		if err == repositories.ErrTokenNotFound || err == repositories.ErrTokenExpired {
			return ErrInvalidToken
		}
		return fmt.Errorf("failed to verify reset token: %w", err)
	}

	// Hash new password
	passwordHash, err := s.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	err = s.userRepo.UpdatePassword(ctx, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	err = s.userRepo.MarkPasswordResetTokenUsed(ctx, token)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to mark reset token as used", zap.Error(err))
		// Don't fail reset for this
	}

	// Deactivate all user sessions for security
	err = s.sessionRepo.DeactivateUserSessions(ctx, userID)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to deactivate user sessions", zap.Error(err))
		// Don't fail reset for this
	}

	logger.InfoCtx(ctx, "Password reset successful", zap.String("user_id", userID))

	return nil
}

// ChangePassword changes user password (requires current password)
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	logger.DebugCtx(ctx, "Processing password change", zap.String("user_id", userID))

	// Verify current password
	err := s.verifyPassword(ctx, userID, currentPassword)
	if err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	passwordHash, err := s.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	err = s.userRepo.UpdatePassword(ctx, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	logger.InfoCtx(ctx, "Password change successful", zap.String("user_id", userID))

	return nil
}

// VerifyEmail verifies user email using verification token
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	logger.DebugCtx(ctx, "Processing email verification")

	// Verify email token
	userID, err := s.userRepo.VerifyEmailToken(ctx, token)
	if err != nil {
		if err == repositories.ErrTokenNotFound || err == repositories.ErrTokenExpired {
			return ErrInvalidToken
		}
		return fmt.Errorf("failed to verify email token: %w", err)
	}

	// Mark user as verified
	err = s.userRepo.MarkUserVerified(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to mark user as verified: %w", err)
	}

	// Mark token as used
	err = s.userRepo.MarkEmailTokenUsed(ctx, token)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to mark email token as used", zap.Error(err))
		// Don't fail verification for this
	}

	logger.InfoCtx(ctx, "Email verification successful", zap.String("user_id", userID))

	return nil
}

// ResendVerification sends a new verification email
func (s *AuthService) ResendVerification(ctx context.Context, email string) error {
	logger.DebugCtx(ctx, "Processing resend verification", zap.String("email", email))

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			// Don't reveal that user doesn't exist
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is already verified
	if user.IsVerified {
		// Don't reveal that user is already verified
		return nil
	}

	// Send verification email
	err = s.sendVerificationEmail(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	logger.InfoCtx(ctx, "Verification email resent", zap.String("user_id", user.ID))

	return nil
}

// GetUser returns user information
func (s *AuthService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user.SanitizeUser(), nil
}

// CreateAnonymousSession creates an anonymous session
func (s *AuthService) CreateAnonymousSession(ctx context.Context) (*models.AnonymousSession, error) {
	logger.DebugCtx(ctx, "Creating anonymous session")

	sessionID := uuid.New().String()
	token, err := s.generateAnonymousToken(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate anonymous token: %w", err)
	}

	session := &models.AnonymousSession{
		ID:        sessionID,
		Token:     token,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(models.DefaultAnonymousDuration),
	}

	err = s.sessionRepo.CreateAnonymousSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create anonymous session: %w", err)
	}

	logger.InfoCtx(ctx, "Anonymous session created", zap.String("session_id", sessionID))

	return session, nil
}

// GetUserSessions returns user's active sessions
func (s *AuthService) GetUserSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	sessions, err := s.sessionRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	return sessions, nil
}

// RevokeSession revokes a specific session
func (s *AuthService) RevokeSession(ctx context.Context, userID, sessionID string) error {
	// Verify session belongs to user
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		if err == repositories.ErrSessionNotFound {
			return ErrSessionNotFound
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session.UserID != userID {
		return ErrSessionNotFound // Don't reveal that session exists for another user
	}

	err = s.sessionRepo.DeactivateSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}

	logger.InfoCtx(ctx, "Session revoked", 
		zap.String("user_id", userID),
		zap.String("session_id", sessionID),
	)

	return nil
}

// Private helper methods

func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) verifyPassword(ctx context.Context, userID, password string) error {
	passwordHash, err := s.userRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get password hash: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("password verification failed: %w", err)
	}

	return nil
}

func (s *AuthService) generateAccessToken(user *models.User, sessionID string) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(models.DefaultSessionDuration)

	claims := &models.TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		SessionID: sessionID,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        claims.UserID,
		"email":      claims.Email,
		"name":       claims.Name,
		"roles":      claims.Roles,
		"session_id": claims.SessionID,
		"iat":        claims.IssuedAt,
		"exp":        claims.ExpiresAt,
		"type":       models.TokenTypeAccess,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) generateRefreshToken(userID, sessionID string) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(models.DefaultRefreshDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        userID,
		"session_id": sessionID,
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"type":       models.TokenTypeRefresh,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) generateAnonymousToken(sessionID string) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(models.DefaultAnonymousDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_id": sessionID,
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"type":       "anonymous",
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign anonymous token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) parseRefreshToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Verify token type
	tokenType, _ := claims["type"].(string)
	if tokenType != models.TokenTypeRefresh {
		return nil, ErrInvalidToken
	}

	userID, _ := claims["sub"].(string)
	sessionID, _ := claims["session_id"].(string)

	if userID == "" || sessionID == "" {
		return nil, ErrInvalidToken
	}

	return &models.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
	}, nil
}

func (s *AuthService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) createSession(ctx context.Context, req *models.SessionCreateRequest) (*models.Session, error) {
	sessionID := uuid.New().String()
	now := time.Now().UTC()

	session := &models.Session{
		ID:        sessionID,
		UserID:    req.UserID,
		DeviceID:  req.DeviceID,
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
		CreatedAt: now,
		ExpiresAt: now.Add(models.DefaultRefreshDuration),
		IsActive:  true,
	}

	err := s.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *AuthService) sendVerificationEmail(ctx context.Context, user *models.User) error {
	// Generate verification token
	verificationToken, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store verification token
	err = s.userRepo.CreateEmailVerificationToken(ctx, user.ID, user.Email, verificationToken, models.DefaultVerifyTokenDuration)
	if err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// TODO: Send actual email using SES
	logger.InfoCtx(ctx, "Verification email would be sent", 
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("token", verificationToken),
	)

	return nil
}

func (s *AuthService) sendPasswordResetEmail(ctx context.Context, user *models.User, resetToken string) error {
	// TODO: Send actual email using SES
	logger.InfoCtx(ctx, "Password reset email would be sent", 
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("token", resetToken),
	)

	return nil
}