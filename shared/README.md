# 🔗 Shared Library Documentation

> **Purpose**: Common types, utilities, constants, and event schemas shared across all services and applications. This reduces code duplication and ensures consistency.

## 📋 Table of Contents

- [🎯 Library Overview](#-library-overview)
- [📊 Type Definitions](#-type-definitions)
- [🛠️ Utility Functions](#️-utility-functions)
- [📡 Event Schemas](#-event-schemas)
- [🔧 Constants](#-constants)
- [🚀 Usage Examples](#-usage-examples)
- [📦 Package Structure](#-package-structure)

---

## 🎯 Library Overview

The shared library provides:

1. **Common Types**: User, Message, Post, and other domain models
2. **Utility Functions**: Validation, formatting, encryption helpers
3. **Event Schemas**: EventBridge event definitions
4. **Constants**: API endpoints, error codes, configuration values
5. **API Clients**: Reusable HTTP clients for service communication

### Design Principles

- **DRY (Don't Repeat Yourself)**: Single source of truth for common code
- **Type Safety**: Strong typing for all shared interfaces
- **Versioning**: Semantic versioning for backward compatibility
- **Documentation**: Comprehensive inline documentation

---

## 📊 Type Definitions

### Core Domain Types (types/domain.go)

```go
package types

import (
    "time"
)

// User represents a platform user across all services
type User struct {
    UserID       string            `json:"userId" validate:"required,uuid4"`
    Email        string            `json:"email" validate:"required,email"`
    DisplayName  string            `json:"displayName" validate:"required,min=2,max=50"`
    Avatar       string            `json:"avatar,omitempty" validate:"omitempty,url"`
    Bio          string            `json:"bio,omitempty" validate:"omitempty,max=500"`
    Location     string            `json:"location,omitempty" validate:"omitempty,max=100"`
    Timezone     string            `json:"timezone,omitempty" validate:"omitempty,timezone"`
    SocialLinks  map[string]string `json:"socialLinks,omitempty"`
    Interests    []string          `json:"interests,omitempty" validate:"omitempty,dive,min=1,max=50"`
    Badges       []string          `json:"badges,omitempty"`
    Stats        UserStats         `json:"stats"`
    Privacy      PrivacySettings   `json:"privacy"`
    Preferences  UserPreferences   `json:"preferences"`
    CreatedAt    time.Time         `json:"createdAt"`
    UpdatedAt    time.Time         `json:"updatedAt"`
    IsActive     bool              `json:"isActive"`
    IsVerified   bool              `json:"isVerified"`
}

// UserStats tracks user activity metrics
type UserStats struct {
    PostsCount     int       `json:"postsCount"`
    LikesReceived  int       `json:"likesReceived"`
    CommentsCount  int       `json:"commentsCount"`
    FollowersCount int       `json:"followersCount"`
    FollowingCount int       `json:"followingCount"`
    LastActiveAt   time.Time `json:"lastActiveAt"`
}

// PrivacySettings controls what information is visible
type PrivacySettings struct {
    DisplayName    PrivacyLevel `json:"displayName"`
    Bio            PrivacyLevel `json:"bio"`
    Avatar         PrivacyLevel `json:"avatar"`
    Location       PrivacyLevel `json:"location"`
    SocialLinks    PrivacyLevel `json:"socialLinks"`
    ActivityStatus PrivacyLevel `json:"activityStatus"`
}

// PrivacyLevel defines visibility levels
type PrivacyLevel string

const (
    PrivacyPublic        PrivacyLevel = "public"        // Visible to everyone
    PrivacyAuthenticated PrivacyLevel = "authenticated" // Visible to logged-in users
    PrivacyPrivate       PrivacyLevel = "private"       // Hidden from everyone
)

// UserPreferences stores user customization settings
type UserPreferences struct {
    Theme         ThemeMode             `json:"theme"`
    Language      string                `json:"language" validate:"required,len=2"`
    Notifications NotificationSettings  `json:"notifications"`
    Privacy       PrivacyPreferences    `json:"privacy"`
}

// ThemeMode defines UI theme options
type ThemeMode string

const (
    ThemeLight ThemeMode = "light"
    ThemeDark  ThemeMode = "dark"
    ThemeAuto  ThemeMode = "auto"
)

// NotificationSettings controls notification delivery
type NotificationSettings struct {
    Email    bool `json:"email"`
    Push     bool `json:"push"`
    InApp    bool `json:"inApp"`
    SMS      bool `json:"sms,omitempty"`
}

// PrivacyPreferences controls privacy behavior
type PrivacyPreferences struct {
    ShowOnlineStatus     bool `json:"showOnlineStatus"`
    AllowDirectMessages  bool `json:"allowDirectMessages"`
    ShowReadReceipts     bool `json:"showReadReceipts"`
    AllowActivityTracking bool `json:"allowActivityTracking"`
}
```

### Message and Chat Types (types/chat.go)

```go
package types

// Message represents a chat message
type Message struct {
    MessageID     string            `json:"messageId" validate:"required"`
    RoomID        string            `json:"roomId" validate:"required"`
    AuthorID      string            `json:"authorId" validate:"required"`
    AuthorName    string            `json:"authorName" validate:"required"`
    AuthorAvatar  string            `json:"authorAvatar,omitempty"`
    Content       string            `json:"content" validate:"required,max=5000"`
    Type          MessageType       `json:"type" validate:"required"`
    Timestamp     time.Time         `json:"timestamp"`
    Reactions     map[string][]string `json:"reactions,omitempty"`
    ReplyTo       string            `json:"replyTo,omitempty"`
    IsEdited      bool              `json:"isEdited"`
    EditedAt      *time.Time        `json:"editedAt,omitempty"`
    IsDeleted     bool              `json:"isDeleted"`
    IsAnonymous   bool              `json:"isAnonymous"`
    AnonymousName string            `json:"anonymousName,omitempty"`
    Attachments   []Attachment      `json:"attachments,omitempty"`
}

// MessageType defines different message types
type MessageType string

const (
    MessageTypeText     MessageType = "text"
    MessageTypeImage    MessageType = "image"
    MessageTypeFile     MessageType = "file"
    MessageTypeAudio    MessageType = "audio"
    MessageTypeVideo    MessageType = "video"
    MessageTypePoll     MessageType = "poll"
    MessageTypeSystem   MessageType = "system"
    MessageTypeAI       MessageType = "ai_response"
)

// Attachment represents a file attachment
type Attachment struct {
    ID           string `json:"id" validate:"required"`
    Type         string `json:"type" validate:"required"`
    URL          string `json:"url" validate:"required,url"`
    ThumbnailURL string `json:"thumbnailUrl,omitempty" validate:"omitempty,url"`
    FileName     string `json:"fileName" validate:"required"`
    FileSize     int64  `json:"fileSize" validate:"required,min=1"`
    MimeType     string `json:"mimeType" validate:"required"`
}

// Room represents a chat room or channel
type Room struct {
    RoomID       string       `json:"roomId" validate:"required"`
    Name         string       `json:"name" validate:"required,min=1,max=100"`
    Description  string       `json:"description,omitempty" validate:"omitempty,max=500"`
    Type         RoomType     `json:"type" validate:"required"`
    Avatar       string       `json:"avatar,omitempty" validate:"omitempty,url"`
    Banner       string       `json:"banner,omitempty" validate:"omitempty,url"`
    CreatedBy    string       `json:"createdBy" validate:"required"`
    CreatedAt    time.Time    `json:"createdAt"`
    UpdatedAt    time.Time    `json:"updatedAt"`
    Members      []RoomMember `json:"members,omitempty"`
    MemberCount  int          `json:"memberCount"`
    MaxMembers   int          `json:"maxMembers"`
    LastActivity time.Time    `json:"lastActivity"`
    LastMessage  string       `json:"lastMessage,omitempty"`
    Settings     RoomSettings `json:"settings"`
    IsActive     bool         `json:"isActive"`
    Tags         []string     `json:"tags,omitempty"`
}

// RoomType defines different room types
type RoomType string

const (
    RoomTypePublic       RoomType = "public"
    RoomTypePrivate      RoomType = "private"
    RoomTypeDirect       RoomType = "direct"
    RoomTypeGroup        RoomType = "group"
    RoomTypeAnnouncement RoomType = "announcement"
)

// RoomMember represents a member of a chat room
type RoomMember struct {
    UserID      string        `json:"userId" validate:"required"`
    DisplayName string        `json:"displayName" validate:"required"`
    Avatar      string        `json:"avatar,omitempty"`
    Role        RoomRole      `json:"role" validate:"required"`
    JoinedAt    time.Time     `json:"joinedAt"`
    LastSeen    time.Time     `json:"lastSeen"`
    IsOnline    bool          `json:"isOnline"`
    IsMuted     bool          `json:"isMuted"`
    IsBanned    bool          `json:"isBanned"`
}

// RoomRole defines member roles in a room
type RoomRole string

const (
    RoomRoleOwner     RoomRole = "owner"
    RoomRoleAdmin     RoomRole = "admin"
    RoomRoleModerator RoomRole = "moderator"
    RoomRoleMember    RoomRole = "member"
)

// RoomSettings configures room behavior
type RoomSettings struct {
    AllowInvites     bool     `json:"allowInvites"`
    RequireApproval  bool     `json:"requireApproval"`
    AllowFileUploads bool     `json:"allowFileUploads"`
    MaxFileSize      int64    `json:"maxFileSize"`
    AllowedFileTypes []string `json:"allowedFileTypes"`
    MessageRetention int      `json:"messageRetention"` // days
    EnableAI         bool     `json:"enableAI"`
    ModerationLevel  string   `json:"moderationLevel"`
}
```

### Post and Social Types (types/social.go)

```go
package types

// Post represents a social media post
type Post struct {
    PostID       string        `json:"postId" validate:"required"`
    AuthorID     string        `json:"authorId" validate:"required"`
    AuthorName   string        `json:"authorName" validate:"required"`
    AuthorAvatar string        `json:"authorAvatar,omitempty"`
    Title        string        `json:"title,omitempty" validate:"omitempty,max=200"`
    Content      string        `json:"content" validate:"required,max=10000"`
    Type         PostType      `json:"type" validate:"required"`
    Status       PostStatus    `json:"status" validate:"required"`
    Category     string        `json:"category,omitempty"`
    Tags         []string      `json:"tags,omitempty" validate:"omitempty,dive,min=1,max=50"`
    Images       []string      `json:"images,omitempty" validate:"omitempty,dive,url"`
    Links        []Link        `json:"links,omitempty"`
    Poll         *Poll         `json:"poll,omitempty"`
    Engagement   Engagement    `json:"engagement"`
    CreatedAt    time.Time     `json:"createdAt"`
    UpdatedAt    time.Time     `json:"updatedAt"`
    IsAnonymous  bool          `json:"isAnonymous"`
    AnonymousName string       `json:"anonymousName,omitempty"`
    IsEdited     bool          `json:"isEdited"`
    EditedAt     *time.Time    `json:"editedAt,omitempty"`
}

// PostType defines different post types
type PostType string

const (
    PostTypeText  PostType = "text"
    PostTypeImage PostType = "image"
    PostTypeLink  PostType = "link"
    PostTypePoll  PostType = "poll"
    PostTypeVideo PostType = "video"
)

// PostStatus defines post visibility status
type PostStatus string

const (
    PostStatusDraft     PostStatus = "draft"
    PostStatusPublished PostStatus = "published"
    PostStatusArchived  PostStatus = "archived"
    PostStatusRemoved   PostStatus = "removed"
)

// Link represents an embedded link
type Link struct {
    URL         string `json:"url" validate:"required,url"`
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Image       string `json:"image,omitempty" validate:"omitempty,url"`
    Domain      string `json:"domain,omitempty"`
}

// Poll represents a poll within a post
type Poll struct {
    Question     string       `json:"question" validate:"required,max=500"`
    Options      []PollOption `json:"options" validate:"required,min=2,max=10"`
    AllowMultiple bool        `json:"allowMultiple"`
    ExpiresAt    *time.Time   `json:"expiresAt,omitempty"`
    TotalVotes   int          `json:"totalVotes"`
    IsActive     bool         `json:"isActive"`
}

// PollOption represents a poll option
type PollOption struct {
    ID       string `json:"id" validate:"required"`
    Text     string `json:"text" validate:"required,max=200"`
    Votes    int    `json:"votes"`
    Voters   []string `json:"voters,omitempty"` // User IDs who voted for this option
}

// Engagement tracks post engagement metrics
type Engagement struct {
    LikesCount    int `json:"likesCount"`
    CommentsCount int `json:"commentsCount"`
    SharesCount   int `json:"sharesCount"`
    ViewsCount    int `json:"viewsCount"`
}

// Comment represents a comment on a post
type Comment struct {
    CommentID     string      `json:"commentId" validate:"required"`
    PostID        string      `json:"postId" validate:"required"`
    AuthorID      string      `json:"authorId" validate:"required"`
    AuthorName    string      `json:"authorName" validate:"required"`
    AuthorAvatar  string      `json:"authorAvatar,omitempty"`
    Content       string      `json:"content" validate:"required,max=2000"`
    ParentID      string      `json:"parentId,omitempty"` // For nested comments
    Depth         int         `json:"depth"`              // Comment nesting level
    LikesCount    int         `json:"likesCount"`
    CreatedAt     time.Time   `json:"createdAt"`
    UpdatedAt     time.Time   `json:"updatedAt"`
    IsAnonymous   bool        `json:"isAnonymous"`
    AnonymousName string      `json:"anonymousName,omitempty"`
    IsEdited      bool        `json:"isEdited"`
    EditedAt      *time.Time  `json:"editedAt,omitempty"`
    IsDeleted     bool        `json:"isDeleted"`
}
```

---

## 🛠️ Utility Functions

### Validation Utilities (utils/validation.go)

```go
package utils

import (
    "regexp"
    "strings"
    "unicode"
    "github.com/go-playground/validator/v10"
)

// ValidatorInstance is a shared validator instance
var ValidatorInstance = validator.New()

// ValidateStruct validates a struct using the validator tags
func ValidateStruct(s interface{}) error {
    return ValidatorInstance.Struct(s)
}

// IsValidEmail checks if an email address is valid
func IsValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return emailRegex.MatchString(strings.ToLower(email))
}

// IsValidUsername checks if a username meets requirements
func IsValidUsername(username string) bool {
    if len(username) < 3 || len(username) > 30 {
        return false
    }
    
    // Must start with alphanumeric character
    if !unicode.IsLetter(rune(username[0])) && !unicode.IsDigit(rune(username[0])) {
        return false
    }
    
    // Can contain letters, numbers, underscores, and hyphens
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    return usernameRegex.MatchString(username)
}

// IsValidPassword checks password strength
func IsValidPassword(password string) bool {
    if len(password) < 8 || len(password) > 128 {
        return false
    }
    
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    return hasUpper && hasLower && hasDigit && hasSpecial
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
    // Remove control characters and trim whitespace
    cleaned := strings.TrimSpace(input)
    cleaned = regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(cleaned, "")
    return cleaned
}

// TruncateString truncates a string to maxLength with ellipsis
func TruncateString(s string, maxLength int) string {
    if len(s) <= maxLength {
        return s
    }
    return s[:maxLength-3] + "..."
}

// ValidateContentLength checks if content meets length requirements
func ValidateContentLength(content string, minLength, maxLength int) bool {
    length := len(strings.TrimSpace(content))
    return length >= minLength && length <= maxLength
}
```

### Encryption and Security Utilities (utils/security.go)

```go
package utils

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "golang.org/x/crypto/bcrypt"
    "fmt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// GenerateRandomToken generates a cryptographically secure random token
func GenerateRandomToken(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateSessionID generates a unique session ID
func GenerateSessionID() (string, error) {
    return GenerateRandomToken(32)
}

// HashSHA256 creates a SHA256 hash of the input
func HashSHA256(input string) string {
    hash := sha256.Sum256([]byte(input))
    return hex.EncodeToString(hash[:])
}

// GenerateAnonymousName generates a random anonymous name
func GenerateAnonymousName() string {
    adjectives := []string{
        "Anonymous", "Mystery", "Secret", "Hidden", "Quiet",
        "Swift", "Bright", "Clever", "Gentle", "Bold",
    }
    
    animals := []string{
        "Panda", "Fox", "Wolf", "Eagle", "Dolphin",
        "Tiger", "Bear", "Hawk", "Owl", "Lion",
        "Cat", "Dog", "Rabbit", "Deer", "Swan",
    }
    
    // Generate random indices
    adjIndex := make([]byte, 1)
    animalIndex := make([]byte, 1)
    number := make([]byte, 2)
    
    rand.Read(adjIndex)
    rand.Read(animalIndex)
    rand.Read(number)
    
    adj := adjectives[int(adjIndex[0])%len(adjectives)]
    animal := animals[int(animalIndex[0])%len(animals)]
    num := int(number[0])<<8 + int(number[1])
    
    return fmt.Sprintf("%s_%s_%d", adj, animal, num%1000)
}

// MaskEmail masks an email for privacy (e.g., test@example.com -> t**t@e****e.com)
func MaskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "***@***.***"
    }
    
    local := parts[0]
    domain := parts[1]
    
    if len(local) <= 2 {
        local = strings.Repeat("*", len(local))
    } else {
        local = string(local[0]) + strings.Repeat("*", len(local)-2) + string(local[len(local)-1])
    }
    
    domainParts := strings.Split(domain, ".")
    if len(domainParts) >= 2 {
        mainDomain := domainParts[0]
        if len(mainDomain) <= 2 {
            mainDomain = strings.Repeat("*", len(mainDomain))
        } else {
            mainDomain = string(mainDomain[0]) + strings.Repeat("*", len(mainDomain)-2) + string(mainDomain[len(mainDomain)-1])
        }
        domainParts[0] = mainDomain
        domain = strings.Join(domainParts, ".")
    }
    
    return local + "@" + domain
}
```

### Time and Date Utilities (utils/time.go)

```go
package utils

import (
    "time"
)

// TimeFormats commonly used time formats
var TimeFormats = struct {
    ISO8601      string
    RFC3339      string
    DateOnly     string
    TimeOnly     string
    DateTime     string
    Timestamp    string
}{
    ISO8601:   "2006-01-02T15:04:05Z07:00",
    RFC3339:   time.RFC3339,
    DateOnly:  "2006-01-02",
    TimeOnly:  "15:04:05",
    DateTime:  "2006-01-02 15:04:05",
    Timestamp: "20060102150405",
}

// FormatTime formats time according to the specified format
func FormatTime(t time.Time, format string) string {
    return t.Format(format)
}

// ParseTime parses time string according to the specified format
func ParseTime(timeStr, format string) (time.Time, error) {
    return time.Parse(format, timeStr)
}

// GetTimeAgo returns a human-readable "time ago" string
func GetTimeAgo(t time.Time) string {
    now := time.Now()
    duration := now.Sub(t)
    
    if duration < time.Minute {
        return "just now"
    } else if duration < time.Hour {
        minutes := int(duration.Minutes())
        if minutes == 1 {
            return "1 minute ago"
        }
        return fmt.Sprintf("%d minutes ago", minutes)
    } else if duration < 24*time.Hour {
        hours := int(duration.Hours())
        if hours == 1 {
            return "1 hour ago"
        }
        return fmt.Sprintf("%d hours ago", hours)
    } else if duration < 7*24*time.Hour {
        days := int(duration.Hours() / 24)
        if days == 1 {
            return "1 day ago"
        }
        return fmt.Sprintf("%d days ago", days)
    } else if duration < 30*24*time.Hour {
        weeks := int(duration.Hours() / (7 * 24))
        if weeks == 1 {
            return "1 week ago"
        }
        return fmt.Sprintf("%d weeks ago", weeks)
    } else if duration < 365*24*time.Hour {
        months := int(duration.Hours() / (30 * 24))
        if months == 1 {
            return "1 month ago"
        }
        return fmt.Sprintf("%d months ago", months)
    } else {
        years := int(duration.Hours() / (365 * 24))
        if years == 1 {
            return "1 year ago"
        }
        return fmt.Sprintf("%d years ago", years)
    }
}

// StartOfDay returns the start of the day for the given time
func StartOfDay(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for the given time
func EndOfDay(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsToday checks if the given time is today
func IsToday(t time.Time) bool {
    now := time.Now()
    return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday checks if the given time is yesterday
func IsYesterday(t time.Time) bool {
    yesterday := time.Now().AddDate(0, 0, -1)
    return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}
```

---

## 📡 Event Schemas

### EventBridge Event Definitions (events/schemas.go)

```go
package events

import (
    "time"
)

// EventType defines the type of event
type EventType string

const (
    // User events
    EventUserRegistered    EventType = "user.registered"
    EventUserLoggedIn      EventType = "user.logged_in"
    EventUserLoggedOut     EventType = "user.logged_out"
    EventUserProfileUpdated EventType = "user.profile_updated"
    EventUserDeleted       EventType = "user.deleted"
    
    // Chat events
    EventMessageSent       EventType = "chat.message_sent"
    EventMessageEdited     EventType = "chat.message_edited"
    EventMessageDeleted    EventType = "chat.message_deleted"
    EventRoomCreated       EventType = "chat.room_created"
    EventRoomMemberJoined  EventType = "chat.room_member_joined"
    EventRoomMemberLeft    EventType = "chat.room_member_left"
    
    // Post events
    EventPostCreated       EventType = "post.created"
    EventPostUpdated       EventType = "post.updated"
    EventPostDeleted       EventType = "post.deleted"
    EventPostLiked         EventType = "post.liked"
    EventPostUnliked       EventType = "post.unliked"
    EventCommentCreated    EventType = "post.comment_created"
    
    // Catalog events
    EventProductCreated    EventType = "catalog.product_created"
    EventProductUpdated    EventType = "catalog.product_updated"
    EventProductPurchased  EventType = "catalog.product_purchased"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
    EventID     string    `json:"eventId"`
    EventType   EventType `json:"eventType"`
    Source      string    `json:"source"`      // Service that generated the event
    Timestamp   time.Time `json:"timestamp"`
    Version     string    `json:"version"`     // Event schema version
    CorrelationID string  `json:"correlationId,omitempty"`
}

// UserRegisteredEvent is fired when a new user registers
type UserRegisteredEvent struct {
    BaseEvent
    Data UserRegisteredData `json:"data"`
}

type UserRegisteredData struct {
    UserID      string `json:"userId"`
    Email       string `json:"email"`
    DisplayName string `json:"displayName"`
    Source      string `json:"source"` // email, google, github, etc.
}

// UserProfileUpdatedEvent is fired when user profile is updated
type UserProfileUpdatedEvent struct {
    BaseEvent
    Data UserProfileUpdatedData `json:"data"`
}

type UserProfileUpdatedData struct {
    UserID       string            `json:"userId"`
    ChangedFields []string         `json:"changedFields"`
    OldValues    map[string]interface{} `json:"oldValues"`
    NewValues    map[string]interface{} `json:"newValues"`
}

// MessageSentEvent is fired when a message is sent
type MessageSentEvent struct {
    BaseEvent
    Data MessageSentData `json:"data"`
}

type MessageSentData struct {
    MessageID   string `json:"messageId"`
    RoomID      string `json:"roomId"`
    AuthorID    string `json:"authorId"`
    Content     string `json:"content"`
    MessageType string `json:"messageType"`
    IsAnonymous bool   `json:"isAnonymous"`
}

// PostCreatedEvent is fired when a post is created
type PostCreatedEvent struct {
    BaseEvent
    Data PostCreatedData `json:"data"`
}

type PostCreatedData struct {
    PostID      string   `json:"postId"`
    AuthorID    string   `json:"authorId"`
    Title       string   `json:"title"`
    Content     string   `json:"content"`
    PostType    string   `json:"postType"`
    Category    string   `json:"category"`
    Tags        []string `json:"tags"`
    IsAnonymous bool     `json:"isAnonymous"`
}
```

### Event Publisher Utility (events/publisher.go)

```go
package events

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/eventbridge"
    "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
    "github.com/google/uuid"
)

// Publisher publishes events to EventBridge
type Publisher struct {
    client    *eventbridge.Client
    eventBusName string
    source    string
}

// NewPublisher creates a new event publisher
func NewPublisher(client *eventbridge.Client, eventBusName, source string) *Publisher {
    return &Publisher{
        client:       client,
        eventBusName: eventBusName,
        source:       source,
    }
}

// PublishEvent publishes an event to EventBridge
func (p *Publisher) PublishEvent(ctx context.Context, eventType EventType, data interface{}) error {
    // Create base event
    baseEvent := BaseEvent{
        EventID:   uuid.New().String(),
        EventType: eventType,
        Source:    p.source,
        Timestamp: time.Now().UTC(),
        Version:   "1.0",
    }
    
    // Serialize event data
    eventData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal event data: %w", err)
    }
    
    // Create EventBridge entry
    entry := types.PutEventsRequestEntry{
        Source:       aws.String(p.source),
        DetailType:   aws.String(string(eventType)),
        Detail:       aws.String(string(eventData)),
        EventBusName: aws.String(p.eventBusName),
        Time:         aws.Time(baseEvent.Timestamp),
    }
    
    // Publish event
    input := &eventbridge.PutEventsInput{
        Entries: []types.PutEventsRequestEntry{entry},
    }
    
    result, err := p.client.PutEvents(ctx, input)
    if err != nil {
        return fmt.Errorf("failed to publish event: %w", err)
    }
    
    // Check for failures
    if result.FailedEntryCount > 0 {
        return fmt.Errorf("failed to publish %d events", result.FailedEntryCount)
    }
    
    return nil
}

// PublishUserRegistered publishes a user registered event
func (p *Publisher) PublishUserRegistered(ctx context.Context, userID, email, displayName, source string) error {
    data := UserRegisteredData{
        UserID:      userID,
        Email:       email,
        DisplayName: displayName,
        Source:      source,
    }
    
    return p.PublishEvent(ctx, EventUserRegistered, data)
}

// PublishMessageSent publishes a message sent event
func (p *Publisher) PublishMessageSent(ctx context.Context, messageID, roomID, authorID, content, messageType string, isAnonymous bool) error {
    data := MessageSentData{
        MessageID:   messageID,
        RoomID:      roomID,
        AuthorID:    authorID,
        Content:     content,
        MessageType: messageType,
        IsAnonymous: isAnonymous,
    }
    
    return p.PublishEvent(ctx, EventMessageSent, data)
}

// PublishPostCreated publishes a post created event
func (p *Publisher) PublishPostCreated(ctx context.Context, postID, authorID, title, content, postType, category string, tags []string, isAnonymous bool) error {
    data := PostCreatedData{
        PostID:      postID,
        AuthorID:    authorID,
        Title:       title,
        Content:     content,
        PostType:    postType,
        Category:    category,
        Tags:        tags,
        IsAnonymous: isAnonymous,
    }
    
    return p.PublishEvent(ctx, EventPostCreated, data)
}
```

---

## 🔧 Constants

### API Constants (constants/api.go)

```go
package constants

// API Configuration
const (
    APIVersion = "v1"
    DefaultTimeout = 30 // seconds
    DefaultPageSize = 20
    MaxPageSize = 100
)

// HTTP Headers
const (
    HeaderAuthorization = "Authorization"
    HeaderContentType   = "Content-Type"
    HeaderUserAgent     = "User-Agent"
    HeaderCorrelationID = "X-Correlation-ID"
    HeaderRequestID     = "X-Request-ID"
    HeaderAPIVersion    = "X-API-Version"
)

// Content Types
const (
    ContentTypeJSON = "application/json"
    ContentTypeForm = "application/x-www-form-urlencoded"
    ContentTypeMultipart = "multipart/form-data"
)

// HTTP Status Codes (commonly used)
const (
    StatusOK                    = 200
    StatusCreated               = 201
    StatusAccepted              = 202
    StatusNoContent             = 204
    StatusBadRequest            = 400
    StatusUnauthorized          = 401
    StatusForbidden             = 403
    StatusNotFound              = 404
    StatusConflict              = 409
    StatusUnprocessableEntity   = 422
    StatusTooManyRequests       = 429
    StatusInternalServerError   = 500
    StatusServiceUnavailable    = 503
)

// Rate Limiting
const (
    RateLimitDefault = 1000 // requests per hour
    RateLimitAuth    = 10   // login attempts per minute
    RateLimitUpload  = 5    // file uploads per minute
)
```

### Error Constants (constants/errors.go)

```go
package constants

// Error Codes
const (
    // Authentication errors (1000-1099)
    ErrCodeInvalidCredentials     = "AUTH001"
    ErrCodeTokenExpired           = "AUTH002"
    ErrCodeTokenInvalid           = "AUTH003"
    ErrCodeAccountLocked          = "AUTH004"
    ErrCodeEmailNotVerified       = "AUTH005"
    
    // Validation errors (1100-1199)
    ErrCodeValidationFailed       = "VAL001"
    ErrCodeInvalidEmail           = "VAL002"
    ErrCodeInvalidPassword        = "VAL003"
    ErrCodeInvalidUsername        = "VAL004"
    ErrCodeContentTooLong         = "VAL005"
    
    // Resource errors (1200-1299)
    ErrCodeResourceNotFound       = "RES001"
    ErrCodeResourceExists         = "RES002"
    ErrCodeResourceDeleted        = "RES003"
    ErrCodeInsufficientPermissions = "RES004"
    
    // Rate limiting errors (1300-1399)
    ErrCodeRateLimitExceeded      = "RATE001"
    ErrCodeTooManyRequests        = "RATE002"
    
    // System errors (1400-1499)
    ErrCodeInternalError          = "SYS001"
    ErrCodeServiceUnavailable     = "SYS002"
    ErrCodeDatabaseError          = "SYS003"
    ErrCodeExternalServiceError   = "SYS004"
)

// Error Messages
var ErrorMessages = map[string]string{
    ErrCodeInvalidCredentials:     "Invalid email or password",
    ErrCodeTokenExpired:           "Authentication token has expired",
    ErrCodeTokenInvalid:           "Invalid authentication token",
    ErrCodeAccountLocked:          "Account is temporarily locked due to multiple failed login attempts",
    ErrCodeEmailNotVerified:       "Email address has not been verified",
    ErrCodeValidationFailed:       "Request validation failed",
    ErrCodeInvalidEmail:           "Invalid email address format",
    ErrCodeInvalidPassword:        "Password does not meet security requirements",
    ErrCodeInvalidUsername:        "Username contains invalid characters or is too short/long",
    ErrCodeContentTooLong:         "Content exceeds maximum allowed length",
    ErrCodeResourceNotFound:       "Requested resource was not found",
    ErrCodeResourceExists:         "Resource already exists",
    ErrCodeResourceDeleted:        "Resource has been deleted",
    ErrCodeInsufficientPermissions: "Insufficient permissions to perform this action",
    ErrCodeRateLimitExceeded:      "Rate limit exceeded, please try again later",
    ErrCodeTooManyRequests:        "Too many requests, please slow down",
    ErrCodeInternalError:          "An internal error occurred, please try again",
    ErrCodeServiceUnavailable:     "Service is temporarily unavailable",
    ErrCodeDatabaseError:          "Database operation failed",
    ErrCodeExternalServiceError:   "External service error",
}
```

---

## 🚀 Usage Examples

### Using Shared Types in Services

```go
// In auth-svc/internal/handlers/register.go
package handlers

import (
    "github.com/your-username/MultitaskProject/shared/types"
    "github.com/your-username/MultitaskProject/shared/utils"
    "github.com/your-username/MultitaskProject/shared/constants"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var req types.RegisterRequest
    
    // Decode request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", constants.StatusBadRequest)
        return
    }
    
    // Validate using shared utilities
    if !utils.IsValidEmail(req.Email) {
        http.Error(w, constants.ErrorMessages[constants.ErrCodeInvalidEmail], constants.StatusBadRequest)
        return
    }
    
    if !utils.IsValidPassword(req.Password) {
        http.Error(w, constants.ErrorMessages[constants.ErrCodeInvalidPassword], constants.StatusBadRequest)
        return
    }
    
    // Hash password using shared utility
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        http.Error(w, constants.ErrorMessages[constants.ErrCodeInternalError], constants.StatusInternalServerError)
        return
    }
    
    // Create user
    user := types.User{
        UserID:      uuid.New().String(),
        Email:       req.Email,
        DisplayName: req.DisplayName,
        CreatedAt:   time.Now(),
        IsActive:    true,
    }
    
    // Validate struct using shared validator
    if err := utils.ValidateStruct(&user); err != nil {
        http.Error(w, constants.ErrorMessages[constants.ErrCodeValidationFailed], constants.StatusBadRequest)
        return
    }
    
    // Save to database and respond...
}
```

### Publishing Events

```go
// In chat-svc/internal/handlers/send_message.go
package handlers

import (
    "github.com/your-username/MultitaskProject/shared/events"
    "github.com/your-username/MultitaskProject/shared/types"
)

func SendMessageHandler(eventPublisher *events.Publisher) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.SendMessageRequest
        
        // ... handle request validation ...
        
        // Create message
        message := types.Message{
            MessageID:   uuid.New().String(),
            RoomID:      req.RoomID,
            AuthorID:    userID,
            Content:     req.Content,
            Type:        req.Type,
            Timestamp:   time.Now(),
            IsAnonymous: req.IsAnonymous,
        }
        
        // Save to database
        if err := saveMessage(message); err != nil {
            // handle error
            return
        }
        
        // Publish event using shared event publisher
        err := eventPublisher.PublishMessageSent(
            r.Context(),
            message.MessageID,
            message.RoomID,
            message.AuthorID,
            message.Content,
            string(message.Type),
            message.IsAnonymous,
        )
        if err != nil {
            // Log error but don't fail the request
            log.Printf("Failed to publish message sent event: %v", err)
        }
        
        // Respond with created message
        w.Header().Set("Content-Type", constants.ContentTypeJSON)
        w.WriteHeader(constants.StatusCreated)
        json.NewEncoder(w).Encode(message)
    }
}
```

### Frontend Usage (TypeScript)

```typescript
// In apps/web/src/types/api.ts
export interface User {
  userId: string;
  email: string;
  displayName: string;
  avatar?: string;
  bio?: string;
  location?: string;
  timezone?: string;
  socialLinks?: Record<string, string>;
  interests?: string[];
  badges?: string[];
  stats: UserStats;
  privacy: PrivacySettings;
  preferences: UserPreferences;
  createdAt: string;
  updatedAt: string;
  isActive: boolean;
  isVerified: boolean;
}

export interface Message {
  messageId: string;
  roomId: string;
  authorId: string;
  authorName: string;
  authorAvatar?: string;
  content: string;
  type: MessageType;
  timestamp: string;
  reactions?: Record<string, string[]>;
  replyTo?: string;
  isEdited: boolean;
  editedAt?: string;
  isDeleted: boolean;
  isAnonymous: boolean;
  anonymousName?: string;
  attachments?: Attachment[];
}

// Error handling with shared constants
export const ErrorCodes = {
  INVALID_CREDENTIALS: 'AUTH001',
  TOKEN_EXPIRED: 'AUTH002',
  TOKEN_INVALID: 'AUTH003',
  // ... other error codes
} as const;

export const ErrorMessages: Record<string, string> = {
  [ErrorCodes.INVALID_CREDENTIALS]: 'Invalid email or password',
  [ErrorCodes.TOKEN_EXPIRED]: 'Authentication token has expired',
  [ErrorCodes.TOKEN_INVALID]: 'Invalid authentication token',
  // ... other error messages
};
```

---

## 📦 Package Structure

```
shared/
├── types/                      # Type definitions
│   ├── domain.go              # Core domain types (User, etc.)
│   ├── chat.go                # Chat-related types
│   ├── social.go              # Social/post-related types
│   ├── catalog.go             # Marketplace-related types
│   ├── api.go                 # API request/response types
│   └── common.go              # Common utility types
│
├── utils/                      # Utility functions
│   ├── validation.go          # Input validation utilities
│   ├── security.go            # Security and encryption utilities
│   ├── time.go                # Time and date utilities
│   ├── string.go              # String manipulation utilities
│   ├── http.go                # HTTP utilities
│   └── json.go                # JSON utilities
│
├── events/                     # Event schemas and publishing
│   ├── schemas.go             # EventBridge event definitions
│   ├── publisher.go           # Event publishing utilities
│   └── handlers.go            # Common event handlers
│
├── constants/                  # Application constants
│   ├── api.go                 # API-related constants
│   ├── errors.go              # Error codes and messages
│   ├── config.go              # Configuration constants
│   └── limits.go              # Rate limits and boundaries
│
├── clients/                    # API clients for service communication
│   ├── auth_client.go         # Auth service client
│   ├── profile_client.go      # Profile service client
│   ├── chat_client.go         # Chat service client
│   └── base_client.go         # Base HTTP client
│
├── middleware/                 # Common middleware
│   ├── auth.go                # Authentication middleware
│   ├── cors.go                # CORS middleware
│   ├── logging.go             # Logging middleware
│   ├── rate_limit.go          # Rate limiting middleware
│   └── recovery.go            # Panic recovery middleware
│
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── README.md                  # This documentation
```

---

**Next**: [🌐 Frontend Application Documentation](../apps/web/README.md)

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/MultitaskProject/issues)
- **Documentation**: [Main README](../../README.md)
- **Shared Library Questions**: [Discussions](https://github.com/your-username/MultitaskProject/discussions)