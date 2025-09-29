# 💬 Chat Service (chat-svc)

> **Purpose**: Handles real-time messaging, chat rooms, AI assistant integration, and message management. This service powers the communication backbone of the platform.

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

The Chat Service is responsible for:

1. **Real-time Messaging**: WebSocket connections for instant message delivery
2. **Chat Rooms**: Public rooms, private groups, direct messages
3. **AI Assistant Integration**: Smart responses using Gemini/OpenAI
4. **Message Management**: History, search, reactions, replies
5. **Presence System**: Online status, typing indicators, read receipts

### Why This Service Exists?

- **Real-time Communication**: Instant messaging capabilities
- **Scalability**: Handle thousands of concurrent connections
- **AI Integration**: Smart assistant for enhanced user experience
- **Flexibility**: Support multiple chat contexts and formats

---

## 🏗️ Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WebApp[React Web App<br/>WebSocket Client]
        MobileApp[Mobile App<br/>WebSocket Client]
    end

    subgraph "API Gateway"
        RestAPI[REST API<br/>Message Management]
        WebSocketAPI[WebSocket API<br/>Real-time Communication]
    end

    subgraph "Chat Service"
        ChatHandler[Chat Handler<br/>Lambda Function]
        MessageService[Message Service<br/>Business Logic]
        RoomService[Room Service<br/>Room Management]
        AIService[AI Integration<br/>Assistant Responses]
        PresenceService[Presence Service<br/>Online Status]
    end

    subgraph "Data Layer"
        DynamoDB[DynamoDB<br/>Messages & Rooms]
        ElastiCache[ElastiCache<br/>Active Sessions]
        S3[S3 Bucket<br/>File Attachments]
    end

    subgraph "External Services"
        GeminiAPI[Gemini API<br/>AI Responses]
        OpenAIAPI[OpenAI API<br/>AI Fallback]
        ProfileSvc[Profile Service<br/>User Data]
    end

    subgraph "Event System"
        EventBridge[EventBridge<br/>Message Events]
        SNS[SNS<br/>Push Notifications]
    end

    WebApp -->|HTTPS| RestAPI
    WebApp -->|WSS| WebSocketAPI
    MobileApp -->|HTTPS| RestAPI
    MobileApp -->|WSS| WebSocketAPI
    
    RestAPI --> ChatHandler
    WebSocketAPI --> ChatHandler
    
    ChatHandler --> MessageService
    ChatHandler --> RoomService
    ChatHandler --> AIService
    ChatHandler --> PresenceService
    
    MessageService --> DynamoDB
    RoomService --> DynamoDB
    PresenceService --> ElastiCache
    AIService --> GeminiAPI
    AIService --> OpenAIAPI
    
    ChatHandler --> EventBridge
    EventBridge --> SNS
    ChatHandler --> ProfileSvc
```

### Real-time Message Flow

```mermaid
sequenceDiagram
    participant Client1
    participant WebSocket as WebSocket API
    participant ChatSvc as Chat Service
    participant DynamoDB
    participant Client2
    participant AIAssistant as AI Assistant

    Client1->>WebSocket: Send message via WebSocket
    WebSocket->>ChatSvc: Process message
    ChatSvc->>DynamoDB: Store message
    ChatSvc->>Client2: Broadcast to room members
    ChatSvc->>AIAssistant: Check for AI trigger
    AIAssistant-->>ChatSvc: Generate AI response
    ChatSvc->>DynamoDB: Store AI response
    ChatSvc->>Client1: Send AI response
    ChatSvc->>Client2: Send AI response
```

---

## ✨ Features

### 💬 Real-time Messaging

#### Message Types
```json
{
  "message_types": [
    "text",
    "image",
    "file",
    "audio",
    "video",
    "poll",
    "system",
    "ai_response"
  ]
}
```

#### Real-time Features
```json
{
  "realtime_features": [
    "Instant message delivery",
    "Typing indicators",
    "Read receipts",
    "Online presence",
    "Message reactions",
    "Reply threading"
  ]
}
```

### 🏠 Chat Rooms & Channels

#### Room Types
```json
{
  "room_types": {
    "public": "Open to everyone, discoverable",
    "private": "Invite-only, not discoverable",
    "direct": "One-on-one conversations",
    "group": "Small private groups (2-50 people)",
    "announcement": "Read-only for most users"
  }
}
```

#### Room Features
```json
{
  "room_features": [
    "Custom room names and descriptions",
    "Room avatars and banners",
    "Member management (add/remove/ban)",
    "Role-based permissions",
    "Message pinning",
    "File sharing"
  ]
}
```

### 🤖 AI Assistant Integration

#### AI Capabilities
```json
{
  "ai_features": [
    "Context-aware responses",
    "Question answering",
    "Content summarization",
    "Language translation",
    "Sentiment analysis",
    "Content moderation"
  ]
}
```

#### Trigger Patterns
```json
{
  "ai_triggers": [
    "@ai help me with...",
    "@assistant what is...",
    "Can you explain...",
    "Summarize this...",
    "Translate: ..."
  ]
}
```

### 📱 Mobile & Web Support

#### Connection Fallbacks
```json
{
  "connection_methods": [
    "WebSocket (primary)",
    "Server-Sent Events (SSE)",
    "Long polling (fallback)",
    "Push notifications (offline)"
  ]
}
```

---

## 🔌 API Endpoints

### REST API Endpoints

#### GET /v1/chat/rooms
Get list of available chat rooms.

```http
GET /v1/chat/rooms?type=public&limit=20&offset=0
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "rooms": [
      {
        "roomId": "room_general_123",
        "name": "General Discussion",
        "description": "Welcome to our community chat!",
        "type": "public",
        "avatar": "https://cdn.multitask.com/rooms/general.jpg",
        "memberCount": 1247,
        "lastActivity": "2025-09-29T10:30:00Z",
        "isJoined": true,
        "unreadCount": 3
      }
    ],
    "pagination": {
      "total": 50,
      "limit": 20,
      "offset": 0,
      "hasMore": true
    }
  }
}
```

#### POST /v1/chat/rooms
Create a new chat room.

```http
POST /v1/chat/rooms
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Project Alpha Discussion",
  "description": "Private room for Project Alpha team",
  "type": "private",
  "maxMembers": 25,
  "allowInvites": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "roomId": "room_alpha_456",
    "name": "Project Alpha Discussion",
    "description": "Private room for Project Alpha team",
    "type": "private",
    "createdBy": "usr_123456789",
    "createdAt": "2025-09-29T10:45:00Z",
    "memberCount": 1,
    "maxMembers": 25
  }
}
```

#### GET /v1/chat/rooms/{roomId}/messages
Get message history for a room.

```http
GET /v1/chat/rooms/room_general_123/messages?limit=50&before=msg_987654321
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "messageId": "msg_abc123def456",
        "roomId": "room_general_123",
        "authorId": "usr_123456789",
        "authorName": "John Doe",
        "authorAvatar": "https://cdn.multitask.com/avatars/usr_123456789.jpg",
        "content": "Hey everyone! How's the project going?",
        "type": "text",
        "timestamp": "2025-09-29T10:25:00Z",
        "reactions": {
          "👍": ["usr_987654321", "usr_555666777"],
          "❤️": ["usr_444333222"]
        },
        "replyTo": null,
        "isEdited": false,
        "isAnonymous": false
      }
    ],
    "pagination": {
      "hasMore": true,
      "before": "msg_abc123def456",
      "limit": 50
    }
  }
}
```

#### POST /v1/chat/rooms/{roomId}/messages
Send a message to a room.

```http
POST /v1/chat/rooms/room_general_123/messages
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "content": "Just deployed the new feature! 🚀",
  "type": "text",
  "replyTo": "msg_abc123def456",
  "isAnonymous": false
}
```

#### POST /v1/chat/rooms/{roomId}/join
Join a chat room.

```http
POST /v1/chat/rooms/room_general_123/join
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### DELETE /v1/chat/rooms/{roomId}/leave
Leave a chat room.

```http
DELETE /v1/chat/rooms/room_general_123/leave
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### WebSocket API

#### Connection
```javascript
// Connect to WebSocket
const ws = new WebSocket('wss://ws.multitask.com/chat?token=jwt_token');

ws.onopen = () => {
  console.log('Connected to chat');
  
  // Join a room
  ws.send(JSON.stringify({
    type: 'join_room',
    roomId: 'room_general_123'
  }));
};
```

#### Send Message
```javascript
// Send a text message
ws.send(JSON.stringify({
  type: 'send_message',
  roomId: 'room_general_123',
  content: 'Hello everyone!',
  messageType: 'text'
}));
```

#### Receive Messages
```javascript
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  switch(data.type) {
    case 'new_message':
      displayMessage(data.message);
      break;
    case 'user_joined':
      showUserJoined(data.user);
      break;
    case 'typing_indicator':
      showTypingIndicator(data.user);
      break;
  }
};
```

#### WebSocket Message Types
```json
{
  "client_to_server": [
    "join_room",
    "leave_room", 
    "send_message",
    "typing_start",
    "typing_stop",
    "message_reaction",
    "mark_read"
  ],
  "server_to_client": [
    "new_message",
    "message_updated",
    "message_deleted",
    "user_joined",
    "user_left",
    "typing_indicator",
    "presence_update",
    "ai_response"
  ]
}
```

### AI Assistant Endpoints

#### POST /v1/chat/ai/ask
Ask the AI assistant a question.

```http
POST /v1/chat/ai/ask
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "question": "Can you explain how microservices work?",
  "context": "software_development",
  "roomId": "room_general_123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "response": "Microservices are a software architecture pattern where applications are built as a collection of small, independent services that communicate over well-defined APIs...",
    "sources": [
      "https://martinfowler.com/articles/microservices.html"
    ],
    "confidence": 0.92,
    "responseTime": 1200
  }
}
```

---

## 📊 Data Models

### Message Model (DynamoDB)
```go
type Message struct {
    MessageID    string            `json:"messageId" dynamodb:"message_id"`     // Partition Key
    RoomID       string            `json:"roomId" dynamodb:"room_id"`           // Global Secondary Index
    AuthorID     string            `json:"authorId" dynamodb:"author_id"`
    AuthorName   string            `json:"authorName" dynamodb:"author_name"`
    AuthorAvatar string            `json:"authorAvatar" dynamodb:"author_avatar"`
    Content      string            `json:"content" dynamodb:"content"`
    Type         string            `json:"type" dynamodb:"type"`                // text, image, file, etc.
    Timestamp    time.Time         `json:"timestamp" dynamodb:"timestamp"`
    Reactions    map[string][]string `json:"reactions" dynamodb:"reactions"`    // emoji -> user_ids
    ReplyTo      string            `json:"replyTo" dynamodb:"reply_to"`         // Parent message ID
    IsEdited     bool              `json:"isEdited" dynamodb:"is_edited"`
    EditedAt     time.Time         `json:"editedAt" dynamodb:"edited_at"`
    IsDeleted    bool              `json:"isDeleted" dynamodb:"is_deleted"`
    IsAnonymous  bool              `json:"isAnonymous" dynamodb:"is_anonymous"`
    AnonymousName string           `json:"anonymousName" dynamodb:"anonymous_name"`
    Attachments  []Attachment      `json:"attachments" dynamodb:"attachments"`
    Metadata     MessageMetadata   `json:"metadata" dynamodb:"metadata"`
}

type Attachment struct {
    ID          string `json:"id" dynamodb:"id"`
    Type        string `json:"type" dynamodb:"type"`        // image, file, audio, video
    URL         string `json:"url" dynamodb:"url"`
    ThumbnailURL string `json:"thumbnailUrl" dynamodb:"thumbnail_url"`
    FileName    string `json:"fileName" dynamodb:"file_name"`
    FileSize    int64  `json:"fileSize" dynamodb:"file_size"`
    MimeType    string `json:"mimeType" dynamodb:"mime_type"`
}

type MessageMetadata struct {
    IPAddress   string `json:"ipAddress" dynamodb:"ip_address"`
    UserAgent   string `json:"userAgent" dynamodb:"user_agent"`
    DeviceType  string `json:"deviceType" dynamodb:"device_type"`
    ReadBy      []ReadReceipt `json:"readBy" dynamodb:"read_by"`
}

type ReadReceipt struct {
    UserID    string    `json:"userId" dynamodb:"user_id"`
    ReadAt    time.Time `json:"readAt" dynamodb:"read_at"`
}
```

### Room Model (DynamoDB)
```go
type Room struct {
    RoomID       string            `json:"roomId" dynamodb:"room_id"`          // Partition Key
    Name         string            `json:"name" dynamodb:"name"`
    Description  string            `json:"description" dynamodb:"description"`
    Type         string            `json:"type" dynamodb:"type"`               // public, private, direct, group
    Avatar       string            `json:"avatar" dynamodb:"avatar"`
    Banner       string            `json:"banner" dynamodb:"banner"`
    CreatedBy    string            `json:"createdBy" dynamodb:"created_by"`
    CreatedAt    time.Time         `json:"createdAt" dynamodb:"created_at"`
    UpdatedAt    time.Time         `json:"updatedAt" dynamodb:"updated_at"`
    Members      []RoomMember      `json:"members" dynamodb:"members"`
    MemberCount  int               `json:"memberCount" dynamodb:"member_count"`
    MaxMembers   int               `json:"maxMembers" dynamodb:"max_members"`
    LastActivity time.Time         `json:"lastActivity" dynamodb:"last_activity"`
    LastMessage  string            `json:"lastMessage" dynamodb:"last_message"`
    Settings     RoomSettings      `json:"settings" dynamodb:"settings"`
    IsActive     bool              `json:"isActive" dynamodb:"is_active"`
    Tags         []string          `json:"tags" dynamodb:"tags"`
}

type RoomMember struct {
    UserID      string    `json:"userId" dynamodb:"user_id"`
    DisplayName string    `json:"displayName" dynamodb:"display_name"`
    Avatar      string    `json:"avatar" dynamodb:"avatar"`
    Role        string    `json:"role" dynamodb:"role"`           // owner, admin, moderator, member
    JoinedAt    time.Time `json:"joinedAt" dynamodb:"joined_at"`
    LastSeen    time.Time `json:"lastSeen" dynamodb:"last_seen"`
    IsOnline    bool      `json:"isOnline" dynamodb:"is_online"`
    IsMuted     bool      `json:"isMuted" dynamodb:"is_muted"`
    IsBanned    bool      `json:"isBanned" dynamodb:"is_banned"`
}

type RoomSettings struct {
    AllowInvites     bool     `json:"allowInvites" dynamodb:"allow_invites"`
    RequireApproval  bool     `json:"requireApproval" dynamodb:"require_approval"`
    AllowFileUploads bool     `json:"allowFileUploads" dynamodb:"allow_file_uploads"`
    MaxFileSize      int64    `json:"maxFileSize" dynamodb:"max_file_size"`
    AllowedFileTypes []string `json:"allowedFileTypes" dynamodb:"allowed_file_types"`
    MessageRetention int      `json:"messageRetention" dynamodb:"message_retention"` // days
    EnableAI         bool     `json:"enableAI" dynamodb:"enable_ai"`
    ModerationLevel  string   `json:"moderationLevel" dynamodb:"moderation_level"`   // none, basic, strict
}
```

### Presence Model (ElastiCache)
```go
type UserPresence struct {
    UserID       string            `json:"userId" redis:"user_id"`
    Status       string            `json:"status" redis:"status"`           // online, away, busy, offline
    LastSeen     time.Time         `json:"lastSeen" redis:"last_seen"`
    CurrentRooms []string          `json:"currentRooms" redis:"current_rooms"`
    DeviceInfo   PresenceDevice    `json:"deviceInfo" redis:"device_info"`
    CustomStatus string            `json:"customStatus" redis:"custom_status"`
    IsTyping     map[string]bool   `json:"isTyping" redis:"is_typing"`      // roomId -> isTyping
}

type PresenceDevice struct {
    Type      string `json:"type" redis:"type"`           // web, mobile, desktop
    Platform  string `json:"platform" redis:"platform"`   // browser, ios, android
    Version   string `json:"version" redis:"version"`     // app version
}
```

### AI Context Model
```go
type AIContext struct {
    ContextID     string            `json:"contextId" dynamodb:"context_id"`
    RoomID        string            `json:"roomId" dynamodb:"room_id"`
    UserID        string            `json:"userId" dynamodb:"user_id"`
    Messages      []ContextMessage  `json:"messages" dynamodb:"messages"`
    CreatedAt     time.Time         `json:"createdAt" dynamodb:"created_at"`
    UpdatedAt     time.Time         `json:"updatedAt" dynamodb:"updated_at"`
    ExpiresAt     time.Time         `json:"expiresAt" dynamodb:"expires_at"`     // TTL
    TokensUsed    int               `json:"tokensUsed" dynamodb:"tokens_used"`
}

type ContextMessage struct {
    Role      string `json:"role" dynamodb:"role"`           // user, assistant, system
    Content   string `json:"content" dynamodb:"content"`
    Timestamp time.Time `json:"timestamp" dynamodb:"timestamp"`
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

# Install Redis (for local development)
redis-server --version

# Install dependencies
go mod tidy
```

### Local Development Setup

1. **Environment Variables**
```bash
# Create .env file
cat > .env << EOF
AWS_REGION=us-east-1
DYNAMODB_TABLE_MESSAGES=chat-messages-dev
DYNAMODB_TABLE_ROOMS=chat-rooms-dev
ELASTICACHE_CLUSTER=chat-presence-dev
S3_BUCKET_ATTACHMENTS=chat-attachments-dev
WEBSOCKET_API_ENDPOINT=wss://ws-dev.multitask.com
GEMINI_API_KEY=your-gemini-api-key
OPENAI_API_KEY=your-openai-api-key
LOG_LEVEL=debug
MAX_MESSAGE_LENGTH=5000
MAX_FILE_SIZE=10485760  # 10MB
TYPING_TIMEOUT=3000     # 3 seconds
EOF
```

2. **Start Local Services**
```bash
# Start Redis for presence (if testing locally)
redis-server

# Run the chat service
go run cmd/lambda/main.go

# Or use AWS SAM for WebSocket testing
sam local start-api --env-vars env.json
```

3. **Test the Service**
```bash
# Create a room
curl -X POST http://localhost:8080/v1/chat/rooms \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Room", "type": "public"}'

# Get rooms
curl -H "Authorization: Bearer $JWT_TOKEN" \
     http://localhost:8080/v1/chat/rooms

# Send a message (via REST)
curl -X POST http://localhost:8080/v1/chat/rooms/room_123/messages \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Hello world!", "type": "text"}'
```

4. **Test WebSocket Connection**
```javascript
// Simple WebSocket test
const ws = new WebSocket('ws://localhost:8080/chat?token=' + JWT_TOKEN);

ws.onopen = () => {
  console.log('Connected');
  ws.send(JSON.stringify({
    type: 'join_room',
    roomId: 'room_123'
  }));
};

ws.onmessage = (event) => {
  console.log('Received:', JSON.parse(event.data));
};
```

---

## 🧪 Testing

### Unit Tests
```go
// Example unit tests
func TestCreateRoom(t *testing.T) {
    // Given: Valid room data
    roomData := CreateRoomRequest{
        Name:        "Test Room",
        Description: "A test room",
        Type:        "public",
    }
    
    // When: Creating room
    room, err := CreateRoom("usr_123", roomData)
    
    // Then: Should create room successfully
    assert.NoError(t, err)
    assert.NotEmpty(t, room.RoomID)
    assert.Equal(t, roomData.Name, room.Name)
    assert.Equal(t, "usr_123", room.CreatedBy)
}

func TestSendMessage(t *testing.T) {
    // Given: Valid message data
    msg := SendMessageRequest{
        RoomID:  "room_123",
        Content: "Test message",
        Type:    "text",
    }
    
    // When: Sending message
    message, err := SendMessage("usr_123", msg)
    
    // Then: Should send message successfully
    assert.NoError(t, err)
    assert.NotEmpty(t, message.MessageID)
    assert.Equal(t, msg.Content, message.Content)
    assert.Equal(t, "usr_123", message.AuthorID)
}

func TestAIResponse(t *testing.T) {
    // Given: AI question
    question := "What is Go programming language?"
    
    // When: Getting AI response
    response, err := GetAIResponse(question, "general")
    
    // Then: Should return valid response
    assert.NoError(t, err)
    assert.NotEmpty(t, response.Response)
    assert.Greater(t, response.Confidence, 0.5)
}
```

### Integration Tests
```bash
# Test with local DynamoDB and Redis
docker-compose up -d dynamodb redis
go test -tags=integration ./...

# Test WebSocket connections
go test -run TestWebSocketConnection ./...

# Test AI integration
GEMINI_API_KEY=test_key go test -run TestAI ./...
```

### Load Testing
```yaml
# artillery.yml for chat service load testing
config:
  target: 'wss://ws.multitask.com'
  phases:
    - duration: 60
      arrivalRate: 10
  processor: "./test-functions.js"

scenarios:
  - name: "WebSocket Chat"
    weight: 80
    engine: ws
    flow:
      - connect:
          url: "/chat?token={{ jwt_token }}"
      - send:
          data: '{"type": "join_room", "roomId": "load_test_room"}'
      - loop:
        - send:
            data: '{"type": "send_message", "roomId": "load_test_room", "content": "Load test message {{ $randomString() }}", "messageType": "text"}'
        - think: 2
        count: 10
        
  - name: "REST API"
    weight: 20
    flow:
      - get:
          url: "/v1/chat/rooms"
          headers:
            Authorization: "Bearer {{ jwt_token }}"
      - post:
          url: "/v1/chat/rooms/load_test_room/messages"
          headers:
            Authorization: "Bearer {{ jwt_token }}"
            Content-Type: "application/json"
          json:
            content: "REST API test message"
            type: "text"
```

---

## 🔧 Configuration

### AWS Resources Required

#### DynamoDB Tables
```yaml
# Messages table
MessagesTable:
  Type: AWS::DynamoDB::Table
  Properties:
    TableName: chat-messages-${env}
    BillingMode: PAY_PER_REQUEST
    AttributeDefinitions:
      - AttributeName: message_id
        AttributeType: S
      - AttributeName: room_id
        AttributeType: S
      - AttributeName: timestamp
        AttributeType: S
    KeySchema:
      - AttributeName: message_id
        KeyType: HASH
    GlobalSecondaryIndexes:
      - IndexName: room-messages-index
        KeySchema:
          - AttributeName: room_id
            KeyType: HASH
          - AttributeName: timestamp
            KeyType: RANGE
        Projection:
          ProjectionType: ALL
    StreamSpecification:
      StreamViewType: NEW_AND_OLD_IMAGES

# Rooms table
RoomsTable:
  Type: AWS::DynamoDB::Table
  Properties:
    TableName: chat-rooms-${env}
    BillingMode: PAY_PER_REQUEST
    AttributeDefinitions:
      - AttributeName: room_id
        AttributeType: S
      - AttributeName: type
        AttributeType: S
    KeySchema:
      - AttributeName: room_id
        KeyType: HASH
    GlobalSecondaryIndexes:
      - IndexName: room-type-index
        KeySchema:
          - AttributeName: type
            KeyType: HASH
        Projection:
          ProjectionType: ALL
```

#### ElastiCache Redis Cluster
```yaml
PresenceCache:
  Type: AWS::ElastiCache::ReplicationGroup
  Properties:
    Description: Chat presence and real-time data
    ReplicationGroupId: chat-presence-${env}
    NumCacheClusters: 2
    Engine: redis
    CacheNodeType: cache.t3.micro
    Port: 6379
    SecurityGroupIds:
      - !Ref PresenceCacheSecurityGroup
    SubnetGroupName: !Ref PresenceCacheSubnetGroup
```

#### API Gateway WebSocket
```yaml
WebSocketAPI:
  Type: AWS::ApiGatewayV2::Api
  Properties:
    Name: chat-websocket-${env}
    ProtocolType: WEBSOCKET
    RouteSelectionExpression: "$request.body.type"
    
ConnectRoute:
  Type: AWS::ApiGatewayV2::Route
  Properties:
    ApiId: !Ref WebSocketAPI
    RouteKey: $connect
    Target: !Sub "integrations/${ConnectIntegration}"
    
DisconnectRoute:
  Type: AWS::ApiGatewayV2::Route
  Properties:
    ApiId: !Ref WebSocketAPI
    RouteKey: $disconnect
    Target: !Sub "integrations/${DisconnectIntegration}"
```

### Environment Variables
```bash
# Required
AWS_REGION=us-east-1
DYNAMODB_TABLE_MESSAGES=chat-messages-prod
DYNAMODB_TABLE_ROOMS=chat-rooms-prod
ELASTICACHE_CLUSTER=chat-presence-prod
WEBSOCKET_API_ENDPOINT=wss://ws.multitask.com
GEMINI_API_KEY=your-gemini-api-key

# Optional
LOG_LEVEL=info
MAX_MESSAGE_LENGTH=5000
MAX_FILE_SIZE=10485760        # 10MB
TYPING_TIMEOUT=3000           # 3 seconds
MESSAGE_RETENTION_DAYS=365    # 1 year
AI_RESPONSE_TIMEOUT=10000     # 10 seconds
PRESENCE_TTL=300              # 5 minutes
WEBSOCKET_IDLE_TIMEOUT=600    # 10 minutes
```

---

## 📈 Monitoring

### Key Metrics to Track

#### Application Metrics
- **Message Throughput**: Messages sent per second
- **Active Connections**: Concurrent WebSocket connections
- **Room Activity**: Messages per room, active rooms
- **AI Usage**: AI requests per hour, response times

#### Performance Metrics
- **Response Time**: API endpoint response times
- **WebSocket Latency**: Message delivery latency
- **Error Rate**: Failed message deliveries
- **Database Performance**: DynamoDB read/write latency

#### Business Metrics
- **User Engagement**: Daily active users in chat
- **Message Volume**: Total messages per day
- **AI Adoption**: Percentage of users using AI features
- **Room Popularity**: Most active chat rooms

### CloudWatch Dashboard
```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "metrics": [
          ["AWS/Lambda", "Invocations", "FunctionName", "chat-svc-${env}"],
          ["AWS/Lambda", "Errors", "FunctionName", "chat-svc-${env}"],
          ["AWS/Lambda", "Duration", "FunctionName", "chat-svc-${env}"],
          ["AWS/ApiGateway", "Count", "ApiName", "chat-websocket-${env}"],
          ["AWS/DynamoDB", "ConsumedReadCapacityUnits", "TableName", "chat-messages-${env}"],
          ["AWS/ElastiCache", "CurrConnections", "CacheClusterId", "chat-presence-${env}"]
        ],
        "period": 300,
        "stat": "Sum",
        "region": "us-east-1",
        "title": "Chat Service Metrics"
      }
    }
  ]
}
```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. WebSocket Connection Issues
```bash
# Check API Gateway WebSocket configuration
aws apigatewayv2 get-api --api-id your-websocket-api-id

# Test WebSocket connection
wscat -c wss://ws.multitask.com/chat?token=jwt_token

# Check Lambda logs for connection/disconnection events
aws logs filter-log-events \
  --log-group-name /aws/lambda/chat-svc-prod \
  --filter-pattern "WebSocket"
```

#### 2. Message Delivery Problems
```bash
# Check DynamoDB for stored messages
aws dynamodb query \
  --table-name chat-messages-prod \
  --index-name room-messages-index \
  --key-condition-expression "room_id = :roomId" \
  --expression-attribute-values '{":roomId":{"S":"room_123"}}'

# Check EventBridge for message events
aws events describe-rule --name chat-message-events-prod
```

#### 3. AI Assistant Not Responding
```bash
# Check AI service configuration
echo $GEMINI_API_KEY | base64  # Verify API key is set

# Test AI API directly
curl -X POST "https://generativeai.googleapis.com/v1/models/gemini-pro:generateContent" \
  -H "Authorization: Bearer $GEMINI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"Hello"}]}]}'
```

---

## 🚀 Deployment

### Using Serverless Framework
```yaml
# serverless.yml for chat-svc
service: chat-svc

provider:
  name: aws
  runtime: go1.x
  region: ${env:AWS_REGION, 'us-east-1'}
  stage: ${env:STAGE, 'dev'}
  websocketsApiName: chat-websocket-${self:provider.stage}
  websocketsApiRouteSelectionExpression: $request.body.type
  
functions:
  chat:
    handler: bin/chat
    events:
      - http:
          path: /{proxy+}
          method: ANY
          cors: true
      - websocket:
          route: $connect
      - websocket:
          route: $disconnect
      - websocket:
          route: join_room
      - websocket:
          route: send_message
    environment:
      DYNAMODB_TABLE_MESSAGES: !Ref MessagesTable
      DYNAMODB_TABLE_ROOMS: !Ref RoomsTable
      WEBSOCKET_API_ENDPOINT: !Sub "${WebsocketsApi}.execute-api.${AWS::Region}.amazonaws.com/${self:provider.stage}"
```

### Build and Deploy
```bash
# Build Go binary
GOOS=linux go build -o bin/chat cmd/lambda/main.go

# Deploy to AWS
serverless deploy --stage prod

# Deploy WebSocket API separately if needed
serverless deploy --stage prod --function chat
```

---

**Next**: [📝 Post Service Documentation](../post-svc/README.md)

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/MultitaskProject/issues)
- **Documentation**: [Main README](../../README.md)
- **API Docs**: [OpenAPI Specification](./docs/openapi.yml)