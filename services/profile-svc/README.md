# 👤 Profile Service (profile-svc)

> **Purpose**: Manages user profiles, display settings, privacy controls, and user preferences. This service handles everything related to user identity and personalization.

## 📋 Table of Contents

- [🎯 Service Overview](#-service-overview)
- [🏗️ Architecture](#️-architecture)
- [✨ Features](#-features)
- [🔌 API Endpoints](#-api-endpoints)
- [📊 Data Models](#-data-models)
- [🚀 Quick Start](#-quick-start)
- [🧪 Testing](#-testing)
- [🔧 Configuration](#-configuration)
- [📈 Monitoring](#-monitoring)

---

## 🎯 Service Overview

The Profile Service is responsible for:

1. **User Profile Management**: Display names, bio, avatar, social links
2. **Privacy Settings**: Control what information is visible to others
3. **Preferences**: Theme, notification settings, language preferences
4. **Anonymous Aliases**: Generate and manage anonymous identities
5. **User Statistics**: Track user activity, posts, engagement metrics

### Why This Service Exists?

- **Separation of Concerns**: Profile data separate from authentication
- **Personalization**: Rich user experience with customizable profiles
- **Privacy Control**: Fine-grained privacy settings
- **Analytics**: User behavior insights without compromising privacy

---

## 🏗️ Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WebApp[React Web App]
        MobileApp[Mobile App]
    end

    subgraph "Profile Service"
        ProfileHandler[Profile Handler<br/>Lambda Function]
        ProfileService[Profile Service<br/>Business Logic]
        AvatarService[Avatar Service<br/>Image Management]
        PrivacyService[Privacy Service<br/>Privacy Controls]
        AliasService[Alias Service<br/>Anonymous Names]
    end

    subgraph "Data Layer"
        DynamoDB[DynamoDB<br/>Profile Data]
        S3[S3 Bucket<br/>Avatar Images]
        CloudFront[CloudFront<br/>Image CDN]
    end

    subgraph "Event Integration"
        EventBridge[EventBridge<br/>Profile Events]
        SNS[SNS<br/>Notifications]
    end

    subgraph "External Services"
        AuthSvc[Auth Service<br/>User Validation]
        ChatSvc[Chat Service<br/>Display Names]
        PostSvc[Post Service<br/>Author Info]
    end

    WebApp -->|HTTPS| ProfileHandler
    MobileApp -->|HTTPS| ProfileHandler
    
    ProfileHandler --> ProfileService
    ProfileHandler --> AvatarService
    ProfileHandler --> PrivacyService
    ProfileHandler --> AliasService
    
    ProfileService --> DynamoDB
    AvatarService --> S3
    AvatarService --> CloudFront
    
    ProfileHandler --> EventBridge
    EventBridge --> SNS
    EventBridge --> ChatSvc
    EventBridge --> PostSvc
    
    ProfileService --> AuthSvc
```

### Data Flow Example

```mermaid
sequenceDiagram
    participant Client
    participant ProfileSvc as Profile Service
    participant DynamoDB
    participant S3
    participant EventBridge

    Client->>ProfileSvc: POST /profile/avatar (multipart)
    ProfileSvc->>S3: Upload image
    S3-->>ProfileSvc: Image URL
    ProfileSvc->>DynamoDB: Update profile with avatar URL
    ProfileSvc->>EventBridge: Publish profile_updated event
    ProfileSvc-->>Client: Updated profile data
    EventBridge-->>ChatSvc: Update display name in chats
    EventBridge-->>PostSvc: Update author info in posts
```

---

## ✨ Features

### 👤 Profile Management

#### Core Profile Data
```json
{
  "profile_features": [
    "Display name and bio",
    "Profile picture/avatar",
    "Social media links",
    "Location and timezone",
    "Interests and skills",
    "Verification badges"
  ]
}
```

#### Avatar System
```json
{
  "avatar_features": [
    "Custom image upload",
    "Auto-generated avatars",
    "Image resizing and optimization",
    "CDN delivery",
    "Fallback to initials"
  ]
}
```

### 🔒 Privacy Controls

#### Visibility Settings
```json
{
  "privacy_levels": {
    "public": "Visible to everyone",
    "authenticated": "Visible to logged-in users only",
    "private": "Hidden from everyone"
  },
  "controllable_fields": [
    "display_name",
    "bio",
    "avatar",
    "location",
    "social_links",
    "activity_status"
  ]
}
```

### 🎭 Anonymous Aliases

#### Dynamic Alias Generation
```json
{
  "alias_features": [
    "Context-aware names (Anonymous_Panda_123)",
    "Session-based aliases",
    "Consistent aliases per room/thread",
    "Conversion tracking",
    "Alias history"
  ]
}
```

### ⚙️ User Preferences

#### Customization Options
```json
{
  "preferences": {
    "theme": ["light", "dark", "auto"],
    "language": ["en", "es", "fr", "de", "ja"],
    "notifications": {
      "email": true,
      "push": true,
      "in_app": true
    },
    "privacy": {
      "show_online_status": false,
      "allow_direct_messages": true
    }
  }
}
```

---

## 🔌 API Endpoints

### Profile Management

#### GET /v1/profile
Get current user's profile.

```http
GET /v1/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "userId": "usr_123456789",
    "displayName": "John Doe",
    "bio": "Software developer passionate about open source",
    "avatar": "https://cdn.multitask.com/avatars/usr_123456789.jpg",
    "location": "San Francisco, CA",
    "timezone": "America/Los_Angeles",
    "socialLinks": {
      "twitter": "https://twitter.com/johndoe",
      "github": "https://github.com/johndoe",
      "linkedin": "https://linkedin.com/in/johndoe"
    },
    "interests": ["golang", "react", "aws", "microservices"],
    "badges": ["verified_developer", "early_adopter"],
    "stats": {
      "postsCount": 42,
      "likesReceived": 156,
      "commentsCount": 89,
      "joinedAt": "2025-01-15T10:00:00Z"
    },
    "privacy": {
      "displayName": "public",
      "bio": "public",
      "avatar": "public",
      "location": "authenticated",
      "socialLinks": "public",
      "activityStatus": "private"
    },
    "preferences": {
      "theme": "dark",
      "language": "en",
      "notifications": {
        "email": true,
        "push": false,
        "inApp": true
      }
    }
  }
}
```

#### PUT /v1/profile
Update user profile.

```http
PUT /v1/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "displayName": "John Smith",
  "bio": "Senior software architect specializing in microservices",
  "location": "New York, NY",
  "timezone": "America/New_York",
  "socialLinks": {
    "twitter": "https://twitter.com/johnsmith",
    "github": "https://github.com/johnsmith"
  },
  "interests": ["golang", "kubernetes", "serverless"]
}
```

#### GET /v1/profile/{userId}
Get another user's public profile.

```http
GET /v1/profile/usr_987654321
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "userId": "usr_987654321",
    "displayName": "Jane Developer",
    "bio": "Full-stack developer and tech blogger",
    "avatar": "https://cdn.multitask.com/avatars/usr_987654321.jpg",
    "badges": ["verified_developer"],
    "stats": {
      "postsCount": 28,
      "likesReceived": 89,
      "joinedAt": "2025-02-20T15:30:00Z"
    }
  }
}
```

### Avatar Management

#### POST /v1/profile/avatar
Upload profile avatar.

```http
POST /v1/profile/avatar
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: multipart/form-data

--boundary
Content-Disposition: form-data; name="avatar"; filename="avatar.jpg"
Content-Type: image/jpeg

[binary image data]
--boundary--
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "avatarUrl": "https://cdn.multitask.com/avatars/usr_123456789.jpg",
    "thumbnailUrl": "https://cdn.multitask.com/avatars/thumbs/usr_123456789.jpg",
    "uploadedAt": "2025-09-29T10:30:00Z"
  }
}
```

#### DELETE /v1/profile/avatar
Remove profile avatar.

```http
DELETE /v1/profile/avatar
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Privacy Settings

#### GET /v1/profile/privacy
Get current privacy settings.

```http
GET /v1/profile/privacy
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### PUT /v1/profile/privacy
Update privacy settings.

```http
PUT /v1/profile/privacy
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "displayName": "public",
  "bio": "authenticated",
  "avatar": "public",
  "location": "private",
  "socialLinks": "authenticated",
  "activityStatus": "private"
}
```

### Anonymous Aliases

#### GET /v1/profile/alias
Get current anonymous alias for a context.

```http
GET /v1/profile/alias?context=chat_room_123&contextType=chat_room
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "aliasId": "alias_abc123def456",
    "displayName": "Anonymous_Panda_789",
    "avatar": "https://cdn.multitask.com/aliases/panda_789.svg",
    "context": "chat_room_123",
    "contextType": "chat_room",
    "createdAt": "2025-09-29T10:00:00Z",
    "isActive": true
  }
}
```

#### POST /v1/profile/alias/regenerate
Generate new anonymous alias for a context.

```http
POST /v1/profile/alias/regenerate
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "context": "chat_room_123",
  "contextType": "chat_room"
}
```

### User Preferences

#### GET /v1/profile/preferences
Get user preferences.

```http
GET /v1/profile/preferences
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### PUT /v1/profile/preferences
Update user preferences.

```http
PUT /v1/profile/preferences
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "theme": "dark",
  "language": "en",
  "notifications": {
    "email": true,
    "push": false,
    "inApp": true
  },
  "privacy": {
    "showOnlineStatus": false,
    "allowDirectMessages": true
  }
}
```

---

## 📊 Data Models

### Profile Model (DynamoDB)
```go
type Profile struct {
    UserID        string            `json:"userId" dynamodb:"user_id"`          // Partition Key
    DisplayName   string            `json:"displayName" dynamodb:"display_name"`
    Bio           string            `json:"bio" dynamodb:"bio"`
    Avatar        string            `json:"avatar" dynamodb:"avatar"`
    Location      string            `json:"location" dynamodb:"location"`
    Timezone      string            `json:"timezone" dynamodb:"timezone"`
    SocialLinks   map[string]string `json:"socialLinks" dynamodb:"social_links"`
    Interests     []string          `json:"interests" dynamodb:"interests"`
    Badges        []string          `json:"badges" dynamodb:"badges"`
    Stats         ProfileStats      `json:"stats" dynamodb:"stats"`
    Privacy       PrivacySettings   `json:"privacy" dynamodb:"privacy"`
    Preferences   UserPreferences   `json:"preferences" dynamodb:"preferences"`
    CreatedAt     time.Time         `json:"createdAt" dynamodb:"created_at"`
    UpdatedAt     time.Time         `json:"updatedAt" dynamodb:"updated_at"`
    IsActive      bool              `json:"isActive" dynamodb:"is_active"`
    Version       int               `json:"version" dynamodb:"version"`        // Optimistic locking
}

type ProfileStats struct {
    PostsCount     int       `json:"postsCount" dynamodb:"posts_count"`
    LikesReceived  int       `json:"likesReceived" dynamodb:"likes_received"`
    CommentsCount  int       `json:"commentsCount" dynamodb:"comments_count"`
    FollowersCount int       `json:"followersCount" dynamodb:"followers_count"`
    FollowingCount int       `json:"followingCount" dynamodb:"following_count"`
    JoinedAt       time.Time `json:"joinedAt" dynamodb:"joined_at"`
    LastActiveAt   time.Time `json:"lastActiveAt" dynamodb:"last_active_at"`
}

type PrivacySettings struct {
    DisplayName    string `json:"displayName" dynamodb:"display_name"`      // public, authenticated, private
    Bio            string `json:"bio" dynamodb:"bio"`
    Avatar         string `json:"avatar" dynamodb:"avatar"`
    Location       string `json:"location" dynamodb:"location"`
    SocialLinks    string `json:"socialLinks" dynamodb:"social_links"`
    ActivityStatus string `json:"activityStatus" dynamodb:"activity_status"`
}

type UserPreferences struct {
    Theme         string                 `json:"theme" dynamodb:"theme"`              // light, dark, auto
    Language      string                 `json:"language" dynamodb:"language"`        // en, es, fr, etc.
    Notifications NotificationSettings   `json:"notifications" dynamodb:"notifications"`
    Privacy       PrivacyPreferences     `json:"privacy" dynamodb:"privacy"`
}

type NotificationSettings struct {
    Email bool `json:"email" dynamodb:"email"`
    Push  bool `json:"push" dynamodb:"push"`
    InApp bool `json:"inApp" dynamodb:"in_app"`
}

type PrivacyPreferences struct {
    ShowOnlineStatus     bool `json:"showOnlineStatus" dynamodb:"show_online_status"`
    AllowDirectMessages  bool `json:"allowDirectMessages" dynamodb:"allow_direct_messages"`
    ShowReadReceipts     bool `json:"showReadReceipts" dynamodb:"show_read_receipts"`
}
```

### Anonymous Alias Model (DynamoDB)
```go
type AnonymousAlias struct {
    AliasID     string    `json:"aliasId" dynamodb:"alias_id"`         // Partition Key
    UserID      string    `json:"userId" dynamodb:"user_id"`           // Global Secondary Index
    DisplayName string    `json:"displayName" dynamodb:"display_name"`
    Avatar      string    `json:"avatar" dynamodb:"avatar"`            // Generated avatar
    Context     string    `json:"context" dynamodb:"context"`          // chat_room_123, post_456
    ContextType string    `json:"contextType" dynamodb:"context_type"` // chat_room, post, comment
    CreatedAt   time.Time `json:"createdAt" dynamodb:"created_at"`
    ExpiresAt   time.Time `json:"expiresAt" dynamodb:"expires_at"`     // TTL
    IsActive    bool      `json:"isActive" dynamodb:"is_active"`
    UsageCount  int       `json:"usageCount" dynamodb:"usage_count"`   // How often used
}
```

### Avatar Model (S3 Metadata)
```go
type Avatar struct {
    UserID       string    `json:"userId"`
    OriginalURL  string    `json:"originalUrl"`   // Full size image
    ThumbnailURL string    `json:"thumbnailUrl"`  // Small thumbnail
    MediumURL    string    `json:"mediumUrl"`     // Medium size
    ContentType  string    `json:"contentType"`   // image/jpeg, image/png
    FileSize     int64     `json:"fileSize"`      // Bytes
    Width        int       `json:"width"`         // Pixels
    Height       int       `json:"height"`        // Pixels
    UploadedAt   time.Time `json:"uploadedAt"`
    S3Key        string    `json:"s3Key"`         // S3 object key
}
```

---

## 🚀 Quick Start

### Prerequisites
```bash
# Install Go (1.22+)
go version

# Install AWS CLI and configure
aws configure

# Install dependencies
go mod tidy
```

### Local Development Setup

1. **Environment Variables**
```bash
# Create .env file
cat > .env << EOF
AWS_REGION=us-east-1
DYNAMODB_TABLE_PROFILES=profile-profiles-dev
DYNAMODB_TABLE_ALIASES=profile-aliases-dev
S3_BUCKET_AVATARS=multitask-avatars-dev
CLOUDFRONT_DOMAIN=d123456789.cloudfront.net
AUTH_SERVICE_URL=https://api.multitask.com/v1/auth
LOG_LEVEL=debug
IMAGE_MAX_SIZE=5242880  # 5MB
ALIAS_EXPIRY_HOURS=720  # 30 days
EOF
```

2. **Start Local Development**
```bash
# Run the service locally
go run cmd/lambda/main.go

# Or use AWS SAM for local Lambda environment
sam local start-api --env-vars env.json
```

3. **Test the Service**
```bash
# Get user profile
curl -H "Authorization: Bearer $JWT_TOKEN" \
     http://localhost:8080/v1/profile

# Update profile
curl -X PUT \
     -H "Authorization: Bearer $JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"displayName": "Updated Name"}' \
     http://localhost:8080/v1/profile

# Upload avatar
curl -X POST \
     -H "Authorization: Bearer $JWT_TOKEN" \
     -F "avatar=@./test-avatar.jpg" \
     http://localhost:8080/v1/profile/avatar
```

---

## 🧪 Testing

### Unit Tests
```go
// Example unit tests
func TestGetProfile(t *testing.T) {
    // Given: Valid user ID
    userID := "usr_123456789"
    
    // When: Getting profile
    profile, err := GetProfile(userID)
    
    // Then: Should return profile
    assert.NoError(t, err)
    assert.NotNil(t, profile)
    assert.Equal(t, userID, profile.UserID)
}

func TestGenerateAlias(t *testing.T) {
    // Given: User and context
    userID := "usr_123456789"
    context := "chat_room_123"
    
    // When: Generating alias
    alias, err := GenerateAlias(userID, context, "chat_room")
    
    // Then: Should create valid alias
    assert.NoError(t, err)
    assert.NotEmpty(t, alias.AliasID)
    assert.Contains(t, alias.DisplayName, "Anonymous_")
    assert.Equal(t, context, alias.Context)
}

func TestUpdatePrivacySettings(t *testing.T) {
    // Given: Valid privacy settings
    userID := "usr_123456789"
    settings := PrivacySettings{
        DisplayName: "public",
        Bio:         "authenticated",
        Avatar:      "public",
    }
    
    // When: Updating privacy settings
    err := UpdatePrivacySettings(userID, settings)
    
    // Then: Should update successfully
    assert.NoError(t, err)
}
```

### Integration Tests
```bash
# Test with local DynamoDB
docker run -p 8000:8000 amazon/dynamodb-local
go test -tags=integration ./...

# Test avatar upload to S3
AWS_PROFILE=test go test -run TestAvatarUpload ./...
```

### Load Testing
```yaml
# artillery.yml for profile service load testing
config:
  target: 'https://api.multitask.com'
  phases:
    - duration: 60
      arrivalRate: 10
  variables:
    jwt_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'

scenarios:
  - name: "Profile Operations"
    weight: 70
    flow:
      - get:
          url: "/v1/profile"
          headers:
            Authorization: "Bearer {{ jwt_token }}"
      - put:
          url: "/v1/profile"
          headers:
            Authorization: "Bearer {{ jwt_token }}"
            Content-Type: "application/json"
          json:
            displayName: "Load Test User {{ $randomString() }}"
            
  - name: "Alias Generation"
    weight: 30
    flow:
      - get:
          url: "/v1/profile/alias"
          headers:
            Authorization: "Bearer {{ jwt_token }}"
          qs:
            context: "load_test_room"
            contextType: "chat_room"
```

---

## 🔧 Configuration

### AWS Resources Required

#### DynamoDB Tables
```yaml
# Profiles table
ProfilesTable:
  Type: AWS::DynamoDB::Table
  Properties:
    TableName: profile-profiles-${env}
    BillingMode: PAY_PER_REQUEST
    AttributeDefinitions:
      - AttributeName: user_id
        AttributeType: S
    KeySchema:
      - AttributeName: user_id
        KeyType: HASH
    StreamSpecification:
      StreamViewType: NEW_AND_OLD_IMAGES

# Aliases table
AliasesTable:
  Type: AWS::DynamoDB::Table
  Properties:
    TableName: profile-aliases-${env}
    BillingMode: PAY_PER_REQUEST
    AttributeDefinitions:
      - AttributeName: alias_id
        AttributeType: S
      - AttributeName: user_id
        AttributeType: S
    KeySchema:
      - AttributeName: alias_id
        KeyType: HASH
    GlobalSecondaryIndexes:
      - IndexName: user-id-index
        KeySchema:
          - AttributeName: user_id
            KeyType: HASH
        Projection:
          ProjectionType: ALL
    TimeToLiveSpecification:
      AttributeName: expires_at
      Enabled: true
```

#### S3 Bucket for Avatars
```yaml
AvatarsBucket:
  Type: AWS::S3::Bucket
  Properties:
    BucketName: multitask-avatars-${env}
    PublicAccessBlockConfiguration:
      BlockPublicAcls: true
      BlockPublicPolicy: true
      IgnorePublicAcls: true
      RestrictPublicBuckets: true
    NotificationConfiguration:
      LambdaConfigurations:
        - Event: s3:ObjectCreated:*
          Function: !GetAtt ImageProcessorFunction.Arn
    CorsConfiguration:
      CorsRules:
        - AllowedHeaders: ["*"]
          AllowedMethods: [PUT, POST]
          AllowedOrigins: ["https://multitask.com"]
          MaxAge: 3600
```

### Environment Variables
```bash
# Required
AWS_REGION=us-east-1
DYNAMODB_TABLE_PROFILES=profile-profiles-prod
DYNAMODB_TABLE_ALIASES=profile-aliases-prod
S3_BUCKET_AVATARS=multitask-avatars-prod
CLOUDFRONT_DOMAIN=cdn.multitask.com

# Optional
LOG_LEVEL=info
IMAGE_MAX_SIZE=5242880         # 5MB max avatar size
IMAGE_ALLOWED_TYPES=jpeg,png,gif
ALIAS_EXPIRY_HOURS=720         # 30 days
PROFILE_CACHE_TTL=300          # 5 minutes
GENERATE_THUMBNAILS=true       # Generate thumbnail images
THUMBNAIL_SIZE=150             # Thumbnail size in pixels
```

---

## 📈 Monitoring

### Key Metrics to Track

#### Application Metrics
- **Profile Updates**: Frequency of profile modifications
- **Avatar Uploads**: Image upload success rate and size distribution
- **Alias Generation**: Anonymous alias creation rate
- **Privacy Changes**: How often users modify privacy settings

#### Performance Metrics
- **Response Time**: API endpoint response times
- **Error Rate**: Failed requests percentage
- **DynamoDB Performance**: Read/write latency and throttling
- **S3 Upload Performance**: Avatar upload success rate and speed

#### Business Metrics
- **Profile Completion**: Percentage of users with complete profiles
- **Anonymous Usage**: How often users use anonymous features
- **Social Links**: Most popular social media integrations
- **Avatar Usage**: Percentage of users with custom avatars

### CloudWatch Dashboard
```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "metrics": [
          ["AWS/Lambda", "Invocations", "FunctionName", "profile-svc-${env}"],
          ["AWS/Lambda", "Errors", "FunctionName", "profile-svc-${env}"],
          ["AWS/Lambda", "Duration", "FunctionName", "profile-svc-${env}"],
          ["AWS/DynamoDB", "ConsumedReadCapacityUnits", "TableName", "profile-profiles-${env}"],
          ["AWS/DynamoDB", "ConsumedWriteCapacityUnits", "TableName", "profile-profiles-${env}"],
          ["AWS/S3", "NumberOfObjects", "BucketName", "multitask-avatars-${env}"]
        ],
        "period": 300,
        "stat": "Average",
        "region": "us-east-1",
        "title": "Profile Service Metrics"
      }
    }
  ]
}
```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. Avatar Upload Failures
```bash
# Check S3 bucket permissions
aws s3api get-bucket-policy --bucket multitask-avatars-dev

# Test direct upload
aws s3 cp test-image.jpg s3://multitask-avatars-dev/test/

# Check CloudFront cache
curl -I https://cdn.multitask.com/avatars/test-image.jpg
```

#### 2. Profile Data Inconsistency
```bash
# Check DynamoDB item
aws dynamodb get-item \
  --table-name profile-profiles-dev \
  --key '{"user_id":{"S":"usr_123456789"}}'

# Check for concurrent updates (optimistic locking)
# Look for version conflicts in logs
```

#### 3. Alias Generation Issues
```bash
# Check alias table
aws dynamodb scan \
  --table-name profile-aliases-dev \
  --filter-expression "user_id = :uid" \
  --expression-attribute-values '{":uid":{"S":"usr_123456789"}}'
```

---

## 🚀 Deployment

### Using Serverless Framework
```yaml
# serverless.yml for profile-svc
service: profile-svc

provider:
  name: aws
  runtime: go1.x
  region: ${env:AWS_REGION, 'us-east-1'}
  stage: ${env:STAGE, 'dev'}
  environment:
    DYNAMODB_TABLE_PROFILES: !Ref ProfilesTable
    DYNAMODB_TABLE_ALIASES: !Ref AliasesTable
    S3_BUCKET_AVATARS: !Ref AvatarsBucket

functions:
  profile:
    handler: bin/profile
    events:
      - http:
          path: /{proxy+}
          method: ANY
          cors: true
    
resources:
  Resources:
    ProfilesTable:
      Type: AWS::DynamoDB::Table
      Properties:
        TableName: profile-profiles-${self:provider.stage}
        # ... table definition
```

### Build and Deploy
```bash
# Build Go binary
GOOS=linux go build -o bin/profile cmd/lambda/main.go

# Deploy to AWS
serverless deploy --stage prod

# Deploy only function
serverless deploy function --function profile --stage prod
```

---

**Next**: [💬 Chat Service Documentation](../chat-svc/README.md)

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/MultitaskProject/issues)
- **Documentation**: [Main README](../../README.md)
- **API Docs**: [OpenAPI Specification](./docs/openapi.yml)