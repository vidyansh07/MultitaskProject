package models

import "time"

// User represents a user in the system
type User struct {
	ID          string    `json:"id" dynamodb:"id"`
	Email       string    `json:"email" dynamodb:"email"`
	Name        string    `json:"name" dynamodb:"name"`
	IsVerified  bool      `json:"is_verified" dynamodb:"is_verified"`
	IsActive    bool      `json:"is_active" dynamodb:"is_active"`
	Roles       []string  `json:"roles" dynamodb:"roles"`
	CreatedAt   time.Time `json:"created_at" dynamodb:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodb:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" dynamodb:"last_login_at,omitempty"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id" dynamodb:"session_id"`
	UserID    string    `json:"user_id" dynamodb:"user_id"`
	Token     string    `json:"-" dynamodb:"token"` // Don't expose token in JSON
	DeviceID  string    `json:"device_id" dynamodb:"device_id"`
	UserAgent string    `json:"user_agent" dynamodb:"user_agent"`
	IPAddress string    `json:"ip_address" dynamodb:"ip_address"`
	CreatedAt time.Time `json:"created_at" dynamodb:"created_at"`
	ExpiresAt time.Time `json:"expires_at" dynamodb:"expires_at"`
	IsActive  bool      `json:"is_active" dynamodb:"is_active"`
}

// AnonymousSession represents an anonymous user session
type AnonymousSession struct {
	ID        string    `json:"id" dynamodb:"anonymous_id"`
	Token     string    `json:"token" dynamodb:"token"`
	CreatedAt time.Time `json:"created_at" dynamodb:"created_at"`
	ExpiresAt time.Time `json:"expires_at" dynamodb:"expires_at"`
}

// Request/Response Models

// LoginRequest represents a login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	DeviceID string `json:"device_id,omitempty"`
}

// RegisterRequest represents a registration request payload
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

// LogoutRequest represents a logout request payload
type LogoutRequest struct {
	SessionID string `json:"session_id,omitempty"` // If empty, logout all sessions
}

// RefreshTokenRequest represents a refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPasswordRequest represents a forgot password request payload
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a password reset request payload
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ChangePasswordRequest represents a password change request payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// VerifyEmailRequest represents an email verification request payload
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents a resend verification request payload
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// AuthResponse represents a successful authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	User         *User  `json:"user"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID    string   `json:"sub"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	SessionID string   `json:"session_id"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	UserID    string    `json:"user_id" dynamodb:"user_id"`
	Token     string    `json:"token" dynamodb:"token"`
	ExpiresAt time.Time `json:"expires_at" dynamodb:"expires_at"`
	Used      bool      `json:"used" dynamodb:"used"`
	CreatedAt time.Time `json:"created_at" dynamodb:"created_at"`
}

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	UserID    string    `json:"user_id" dynamodb:"user_id"`
	Token     string    `json:"token" dynamodb:"token"`
	Email     string    `json:"email" dynamodb:"email"`
	ExpiresAt time.Time `json:"expires_at" dynamodb:"expires_at"`
	Used      bool      `json:"used" dynamodb:"used"`
	CreatedAt time.Time `json:"created_at" dynamodb:"created_at"`
}

// UserCreateRequest represents internal user creation request
type UserCreateRequest struct {
	Email        string   `json:"email"`
	PasswordHash string   `json:"password_hash"`
	Name         string   `json:"name"`
	Roles        []string `json:"roles"`
}

// SessionCreateRequest represents internal session creation request
type SessionCreateRequest struct {
	UserID    string `json:"user_id"`
	DeviceID  string `json:"device_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
}

// Constants for token types
const (
	TokenTypeAccess           = "access"
	TokenTypeRefresh          = "refresh"
	TokenTypePasswordReset    = "password_reset"
	TokenTypeEmailVerification = "email_verification"
)

// Constants for user roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
	RoleMod   = "moderator"
)

// Default values
const (
	DefaultSessionDuration     = 15 * time.Minute    // Access token duration
	DefaultRefreshDuration     = 7 * 24 * time.Hour  // Refresh token duration
	DefaultResetTokenDuration  = 1 * time.Hour       // Password reset token duration
	DefaultVerifyTokenDuration = 24 * time.Hour      // Email verification token duration
	DefaultAnonymousDuration   = 24 * time.Hour      // Anonymous session duration
)

// Validation helper methods

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	validRoles := []string{RoleUser, RoleAdmin, RoleMod}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// SanitizeUser returns a user object safe for public consumption
func (u *User) SanitizeUser() *User {
	return &User{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		IsVerified:  u.IsVerified,
		IsActive:    u.IsActive,
		Roles:       u.Roles,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
	}
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role string) bool {
	for _, userRole := range u.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}

// IsExpired checks if a session is expired
func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// IsExpired checks if an anonymous session is expired
func (a *AnonymousSession) IsExpired() bool {
	return time.Now().UTC().After(a.ExpiresAt)
}