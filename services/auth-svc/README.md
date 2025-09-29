# 🔐 Authentication Service (auth-svc)

> **Purpose**: Handles user authentication, authorization, and anonymous user management. This service is the security gateway for the entire platform.

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

The Authentication Service is responsible for:

1. **User Registration & Login**: Email/password, social logins (Google, GitHub)
2. **Anonymous Users**: Generate temporary identities for anonymous browsing
3. **JWT Token Management**: Issue, validate, and refresh tokens
4. **Session Management**: Handle user sessions across devices
5. **Security**: Rate limiting, brute force protection, password policies

### Why This Service Exists?

- **Centralized Security**: Single point for all authentication logic
- **Scalability**: Handle millions of concurrent users
- **Flexibility**: Support multiple authentication methods
- **Compliance**: Meet security standards (GDPR, SOC2)

---

## 🏗️ Architecture

```mermaid
graph TB
    subgraph "Client Applications"
        WebApp[React Web App]
        MobileApp[Mobile App]
        API[API Client]
    end

    subgraph "Authentication Service"
        AuthHandler[Auth Handler<br/>Lambda Function]
        TokenService[Token Service<br/>JWT Management]
        UserService[User Service<br/>User Operations]
        AnonymousService[Anonymous Service<br/>Temp Users]
    end

    subgraph "Data Layer"
        Cognito[AWS Cognito<br/>User Pool]
        DynamoDB[DynamoDB<br/>Sessions & Anonymous Users]
        SSM[SSM<br/>JWT Secrets]
    end

    subgraph "Security Layer"
        WAF[AWS WAF<br/>Rate Limiting]
        Shield[AWS Shield<br/>DDoS Protection]
    end

    WebApp -->|HTTPS| WAF
    MobileApp -->|HTTPS| WAF
    API -->|HTTPS| WAF
    
    WAF --> AuthHandler
    Shield --> AuthHandler
    
    AuthHandler --> TokenService
    AuthHandler --> UserService
    AuthHandler --> AnonymousService
    
    TokenService --> SSM
    UserService --> Cognito
    UserService --> DynamoDB
    AnonymousService --> DynamoDB
```

### Service Interactions

```mermaid
sequenceDiagram
    participant Client
    participant AuthSvc as Auth Service
    participant Cognito
    participant DynamoDB
    participant EventBridge

    Client->>AuthSvc: POST /login
    AuthSvc->>Cognito: Verify credentials
    Cognito-->>AuthSvc: User details
    AuthSvc->>DynamoDB: Create session
    AuthSvc->>EventBridge: Publish user_logged_in event
    AuthSvc-->>Client: JWT tokens (access + refresh)
```

---

## ✨ Features

### 🔑 Authentication Methods

#### 1. Email/Password Authentication
```json
{
  "method": "email_password",
  "features": [
    "Password strength validation",
    "Account verification via email",
    "Password reset flow",
    "Brute force protection"
  ]
}
```

#### 2. Social Authentication
```json
{
  "providers": ["google", "github", "facebook"],
  "features": [
    "OAuth 2.0 flow",
    "Profile data import",
    "Account linking",
    "Privacy controls"
  ]
}
```

#### 3. Anonymous Authentication
```json
{
  "anonymous_features": [
    "Temporary user IDs",
    "Generated display names",
    "Limited session duration",
    "Conversion to registered user"
  ]
}
```

### 🛡️ Security Features

- **JWT Tokens**: Short-lived access tokens (15 min) + long-lived refresh tokens (7 days)
- **Rate Limiting**: 5 login attempts per minute per IP
- **Session Management**: Multiple device support with session invalidation
- **Password Policy**: Minimum 8 characters, mixed case, numbers, symbols
- **Account Lockout**: Temporary lockout after 5 failed attempts

---

## 🔌 API Endpoints

### Authentication Endpoints

#### POST /v1/auth/register
Register a new user account.

```http
POST /v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "displayName": "John Doe",
  "acceptTerms": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Account created successfully. Please verify your email.",
  "data": {
    "userId": "usr_123456789",
    "email": "user@example.com",
    "displayName": "John Doe",
    "emailVerified": false,
    "createdAt": "2025-09-29T10:00:00Z"
  }
}
```

#### POST /v1/auth/login
Authenticate user and return JWT tokens.

```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "rt_987654321abcdef",
    "expiresIn": 900,
    "tokenType": "Bearer",
    "user": {
      "userId": "usr_123456789",
      "email": "user@example.com",
      "displayName": "John Doe",
      "roles": ["user"]
    }
  }
}
```

#### POST /v1/auth/anonymous
Create an anonymous user session.

```http
POST /v1/auth/anonymous
Content-Type: application/json

{
  "deviceInfo": {
    "userAgent": "Mozilla/5.0...",
    "platform": "web"
  }
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "anonymousId": "anon_abc123def456",
    "displayName": "Anonymous_Panda_789",
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600,
    "capabilities": ["read", "post_anonymous", "chat_anonymous"]
  }
}
```

#### POST /v1/auth/refresh
Refresh access token using refresh token.

```http
POST /v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "rt_987654321abcdef"
}
```

#### POST /v1/auth/logout
Logout user and invalidate tokens.

```http
POST /v1/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "refreshToken": "rt_987654321abcdef"
}
```

### User Management Endpoints

#### GET /v1/auth/profile
Get current user profile.

```http
GET /v1/auth/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### POST /v1/auth/change-password
Change user password.

```http
POST /v1/auth/change-password
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "currentPassword": "OldPass123!",
  "newPassword": "NewSecurePass456!"
}
```

#### POST /v1/auth/forgot-password
Request password reset.

```http
POST /v1/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

### Social Authentication Endpoints

#### GET /v1/auth/social/{provider}
Initiate social login flow.

```http
GET /v1/auth/social/google?redirect_uri=https://yourapp.com/callback
```

#### POST /v1/auth/social/{provider}/callback
Handle social login callback.

```http
POST /v1/auth/social/google/callback
Content-Type: application/json

{
  "code": "authorization_code_from_provider",
  "state": "random_state_value"
}
```

---

## 📊 Data Models

### User Model (Cognito User Pool)
```go
type User struct {
    UserID       string    `json:"userId" cognito:"sub"`
    Email        string    `json:"email" cognito:"email"`
    EmailVerified bool     `json:"emailVerified" cognito:"email_verified"`
    DisplayName  string    `json:"displayName" cognito:"name"`
    GivenName    string    `json:"givenName" cognito:"given_name"`
    FamilyName   string    `json:"familyName" cognito:"family_name"`
    PhoneNumber  string    `json:"phoneNumber" cognito:"phone_number"`
    Picture      string    `json:"picture" cognito:"picture"`
    Roles        []string  `json:"roles" cognito:"custom:roles"`
    CreatedAt    time.Time `json:"createdAt" cognito:"created_at"`
    UpdatedAt    time.Time `json:"updatedAt" cognito:"updated_at"`
    LastLogin    time.Time `json:"lastLogin" cognito:"last_login"`
    Status       string    `json:"status" cognito:"status"` // ACTIVE, INACTIVE, SUSPENDED
}
```

### Session Model (DynamoDB)
```go
type Session struct {
    SessionID    string    `json:"sessionId" dynamodb:"session_id"`     // Partition Key
    UserID       string    `json:"userId" dynamodb:"user_id"`           // Global Secondary Index
    RefreshToken string    `json:"refreshToken" dynamodb:"refresh_token"`
    DeviceInfo   DeviceInfo `json:"deviceInfo" dynamodb:"device_info"`
    CreatedAt    time.Time `json:"createdAt" dynamodb:"created_at"`
    ExpiresAt    time.Time `json:"expiresAt" dynamodb:"expires_at"`     // TTL
    IsActive     bool      `json:"isActive" dynamodb:"is_active"`
    IPAddress    string    `json:"ipAddress" dynamodb:"ip_address"`
    UserAgent    string    `json:"userAgent" dynamodb:"user_agent"`
}

type DeviceInfo struct {
    Platform    string `json:"platform"`    // web, ios, android
    DeviceType  string `json:"deviceType"`  // desktop, mobile, tablet
    OS          string `json:"os"`          // windows, macos, ios, android
    Browser     string `json:"browser"`     // chrome, safari, firefox
    AppVersion  string `json:"appVersion"`  // 1.0.0
}
```

### Anonymous User Model (DynamoDB)
```go
type AnonymousUser struct {
    AnonymousID  string    `json:"anonymousId" dynamodb:"anonymous_id"`  // Partition Key
    DisplayName  string    `json:"displayName" dynamodb:"display_name"`  // Generated name
    CreatedAt    time.Time `json:"createdAt" dynamodb:"created_at"`
    ExpiresAt    time.Time `json:"expiresAt" dynamodb:"expires_at"`      // TTL (24 hours)
    DeviceInfo   DeviceInfo `json:"deviceInfo" dynamodb:"device_info"`
    Capabilities []string  `json:"capabilities" dynamodb:"capabilities"` // What they can do
    ConvertedTo  string    `json:"convertedTo" dynamodb:"converted_to"`  // UserID if converted
}
```

### Rate Limiting Model (DynamoDB)
```go
type RateLimit struct {
    Key          string    `json:"key" dynamodb:"key"`                    // IP:endpoint or user:endpoint
    RequestCount int       `json:"requestCount" dynamodb:"request_count"`
    WindowStart  time.Time `json:"windowStart" dynamodb:"window_start"`
    ExpiresAt    time.Time `json:"expiresAt" dynamodb:"expires_at"`       // TTL
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
COGNITO_USER_POOL_ID=us-east-1_abcd1234
COGNITO_CLIENT_ID=1234567890abcdefghij
COGNITO_CLIENT_SECRET=secret123
JWT_SECRET=your-super-secret-jwt-key
DYNAMODB_TABLE_SESSIONS=auth-sessions-dev
DYNAMODB_TABLE_ANONYMOUS=auth-anonymous-dev
DYNAMODB_TABLE_RATE_LIMITS=auth-rate-limits-dev
LOG_LEVEL=debug
EOF
```

2. **Start Local Development**
```bash
# Run the service locally
go run cmd/lambda/main.go

# Or use AWS SAM for local Lambda environment
sam local start-api
```

3. **Test the Service**
```bash
# Register a new user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!",
    "displayName": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!"
  }'
```

---

## 🧪 Testing

### Unit Tests
```go
// Example unit test
func TestAuthenticateUser(t *testing.T) {
    // Given: Valid user credentials
    email := "test@example.com"
    password := "ValidPass123!"
    
    // When: Authenticating user
    user, token, err := AuthenticateUser(email, password)
    
    // Then: Should return user and valid token
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotEmpty(t, token)
    assert.Equal(t, email, user.Email)
}

func TestGenerateAnonymousUser(t *testing.T) {
    // When: Generating anonymous user
    anonUser, err := GenerateAnonymousUser()
    
    // Then: Should create valid anonymous user
    assert.NoError(t, err)
    assert.NotEmpty(t, anonUser.AnonymousID)
    assert.Contains(t, anonUser.DisplayName, "Anonymous_")
    assert.NotEmpty(t, anonUser.Capabilities)
}
```

### Integration Tests
```bash
# Run integration tests against local DynamoDB
go test -tags=integration ./...

# Run tests against real AWS services (staging)
AWS_PROFILE=staging go test -tags=integration ./...
```

### Load Testing
```bash
# Install artillery for load testing
npm install -g artillery

# Run load test
artillery run tests/load/auth-load-test.yml
```

---

## 🔧 Configuration

### AWS Resources Required

#### Cognito User Pool
```yaml
CognitoUserPool:
  Type: AWS::Cognito::UserPool
  Properties:
    UserPoolName: multitask-users-${env}
    UsernameAttributes:
      - email
    PasswordPolicy:
      MinimumLength: 8
      RequireUppercase: true
      RequireLowercase: true
      RequireNumbers: true
      RequireSymbols: true
    Schema:
      - Name: email
        Required: true
        Mutable: false
      - Name: name
        Required: true
        Mutable: true
      - Name: roles
        AttributeDataType: String
        Mutable: true
```

#### DynamoDB Tables
```yaml
SessionsTable:
  Type: AWS::DynamoDB::Table
  Properties:
    TableName: auth-sessions-${env}
    BillingMode: PAY_PER_REQUEST
    AttributeDefinitions:
      - AttributeName: session_id
        AttributeType: S
      - AttributeName: user_id
        AttributeType: S
    KeySchema:
      - AttributeName: session_id
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

### Environment Variables
```bash
# Required
AWS_REGION=us-east-1
COGNITO_USER_POOL_ID=us-east-1_abcd1234
COGNITO_CLIENT_ID=1234567890abcdefghij
JWT_SECRET=your-super-secret-jwt-key

# Optional
LOG_LEVEL=info                    # debug, info, warn, error
JWT_EXPIRY=15m                    # Access token expiry
REFRESH_TOKEN_EXPIRY=168h         # Refresh token expiry (7 days)
ANONYMOUS_SESSION_EXPIRY=24h      # Anonymous session expiry
RATE_LIMIT_REQUESTS=60            # Requests per minute
RATE_LIMIT_WINDOW=1m              # Rate limit window
PASSWORD_MIN_LENGTH=8             # Minimum password length
```

---

## 📈 Monitoring

### Key Metrics to Track

#### Application Metrics
- **Authentication Rate**: Successful logins per minute
- **Registration Rate**: New user signups per hour
- **Error Rate**: Failed authentication attempts percentage
- **Response Time**: Average API response time
- **Token Refresh Rate**: Refresh token usage frequency

#### Security Metrics
- **Failed Login Attempts**: Brute force detection
- **Account Lockouts**: Suspicious activity indicators
- **Anonymous User Creation**: Usage patterns
- **Session Duration**: Average user session length

#### Business Metrics
- **User Growth**: Daily/Monthly active users
- **Conversion Rate**: Anonymous to registered users
- **Social Login Usage**: Preferred authentication methods

### CloudWatch Dashboard
```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "metrics": [
          ["AWS/Lambda", "Invocations", "FunctionName", "auth-svc-${env}"],
          ["AWS/Lambda", "Errors", "FunctionName", "auth-svc-${env}"],
          ["AWS/Lambda", "Duration", "FunctionName", "auth-svc-${env}"]
        ],
        "period": 300,
        "stat": "Sum",
        "region": "us-east-1",
        "title": "Auth Service Metrics"
      }
    }
  ]
}
```

### Alarms
```yaml
HighErrorRateAlarm:
  Type: AWS::CloudWatch::Alarm
  Properties:
    AlarmName: auth-svc-high-error-rate-${env}
    AlarmDescription: High error rate in auth service
    MetricName: Errors
    Namespace: AWS/Lambda
    Statistic: Sum
    Period: 300
    EvaluationPeriods: 2
    Threshold: 10
    ComparisonOperator: GreaterThanThreshold
    AlarmActions:
      - !Ref SNSTopicAlerts
```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. "Invalid JWT Token"
```bash
# Check token expiry
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." | base64 -d

# Verify JWT secret in SSM
aws ssm get-parameter --name "/multitask/${env}/jwt-secret" --with-decryption
```

#### 2. Cognito User Pool Issues
```bash
# Check user pool configuration
aws cognito-idp describe-user-pool --user-pool-id us-east-1_abcd1234

# List users in pool
aws cognito-idp list-users --user-pool-id us-east-1_abcd1234
```

#### 3. DynamoDB Connection Issues
```bash
# Test DynamoDB connectivity
aws dynamodb describe-table --table-name auth-sessions-dev

# Check table item count
aws dynamodb scan --table-name auth-sessions-dev --select COUNT
```

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug

# Run with verbose output
go run -v cmd/lambda/main.go
```

---

## 🚀 Deployment

### Using Serverless Framework
```yaml
# serverless.yml for auth-svc
service: auth-svc

provider:
  name: aws
  runtime: go1.x
  region: ${env:AWS_REGION, 'us-east-1'}
  stage: ${env:STAGE, 'dev'}
  
functions:
  auth:
    handler: bin/auth
    events:
      - http:
          path: /{proxy+}
          method: ANY
          cors: true
    environment:
      COGNITO_USER_POOL_ID: !Ref CognitoUserPool
      JWT_SECRET: ${ssm:/multitask/${self:provider.stage}/jwt-secret}
```

### Deployment Commands
```bash
# Build and deploy
make build
serverless deploy --stage prod

# Deploy specific function
serverless deploy function --function auth --stage prod

# View logs
serverless logs --function auth --tail
```

---

**Next**: [👤 Profile Service Documentation](../profile-svc/README.md)

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/MultitaskProject/issues)
- **Documentation**: [Main README](../../README.md)
- **API Docs**: [OpenAPI Specification](./docs/openapi.yml)