package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/multitask-platform/backend/services/auth-svc/internal/models"
)

// Common repository errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenExpired    = errors.New("token expired")
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User CRUD operations
	CreateUser(ctx context.Context, user *models.User, passwordHash string) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error

	// Password operations
	GetPasswordHash(ctx context.Context, userID string) (string, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID string, loginTime time.Time) error

	// Email verification
	CreateEmailVerificationToken(ctx context.Context, userID, email, token string, duration time.Duration) error
	VerifyEmailToken(ctx context.Context, token string) (string, error) // returns userID
	MarkEmailTokenUsed(ctx context.Context, token string) error
	MarkUserVerified(ctx context.Context, userID string) error

	// Password reset
	CreatePasswordResetToken(ctx context.Context, userID, token string, duration time.Duration) error
	VerifyPasswordResetToken(ctx context.Context, token string) (string, error) // returns userID
	MarkPasswordResetTokenUsed(ctx context.Context, token string) error
}

// SessionRepository defines the interface for session data operations
type SessionRepository interface {
	// Session CRUD operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, sessionID string) (*models.Session, error)
	GetUserSessions(ctx context.Context, userID string) ([]*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, sessionID string) error

	// Session management
	DeactivateSession(ctx context.Context, sessionID string) error
	DeactivateUserSessions(ctx context.Context, userID string) error
	CleanupExpiredSessions(ctx context.Context) error

	// Anonymous sessions
	CreateAnonymousSession(ctx context.Context, session *models.AnonymousSession) error
	GetAnonymousSession(ctx context.Context, sessionID string) (*models.AnonymousSession, error)
	DeleteAnonymousSession(ctx context.Context, sessionID string) error
	CleanupExpiredAnonymousSessions(ctx context.Context) error
}

// Mock implementations for now (will be replaced with DynamoDB implementations)

type MockUserRepository struct{}

func NewDynamoDBUserRepository() UserRepository {
	return &MockUserRepository{}
}

func (r *MockUserRepository) CreateUser(ctx context.Context, user *models.User, passwordHash string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	// TODO: Implement DynamoDB operations
	return nil, ErrUserNotFound
}

func (r *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// TODO: Implement DynamoDB operations
	return nil, ErrUserNotFound
}

func (r *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) DeleteUser(ctx context.Context, userID string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) GetPasswordHash(ctx context.Context, userID string) (string, error) {
	// TODO: Implement DynamoDB operations
	return "", ErrUserNotFound
}

func (r *MockUserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string, loginTime time.Time) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) CreateEmailVerificationToken(ctx context.Context, userID, email, token string, duration time.Duration) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) VerifyEmailToken(ctx context.Context, token string) (string, error) {
	// TODO: Implement DynamoDB operations
	return "", ErrTokenNotFound
}

func (r *MockUserRepository) MarkEmailTokenUsed(ctx context.Context, token string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) MarkUserVerified(ctx context.Context, userID string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) CreatePasswordResetToken(ctx context.Context, userID, token string, duration time.Duration) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockUserRepository) VerifyPasswordResetToken(ctx context.Context, token string) (string, error) {
	// TODO: Implement DynamoDB operations
	return "", ErrTokenNotFound
}

func (r *MockUserRepository) MarkPasswordResetTokenUsed(ctx context.Context, token string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

type MockSessionRepository struct{}

func NewDynamoDBSessionRepository() SessionRepository {
	return &MockSessionRepository{}
}

func (r *MockSessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	// TODO: Implement DynamoDB operations
	return nil, ErrSessionNotFound
}

func (r *MockSessionRepository) GetUserSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	// TODO: Implement DynamoDB operations
	return nil, nil
}

func (r *MockSessionRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) DeactivateSession(ctx context.Context, sessionID string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) DeactivateUserSessions(ctx context.Context, userID string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) CreateAnonymousSession(ctx context.Context, session *models.AnonymousSession) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) GetAnonymousSession(ctx context.Context, sessionID string) (*models.AnonymousSession, error) {
	// TODO: Implement DynamoDB operations
	return nil, ErrSessionNotFound
}

func (r *MockSessionRepository) DeleteAnonymousSession(ctx context.Context, sessionID string) error {
	// TODO: Implement DynamoDB operations
	return nil
}

func (r *MockSessionRepository) CleanupExpiredAnonymousSessions(ctx context.Context) error {
	// TODO: Implement DynamoDB operations
	return nil
}