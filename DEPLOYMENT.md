# 📄 Multitask Platform - Deployment Guide

## 🚀 Quick Start Deployment

### Prerequisites
- Node.js 18+ and npm
- Go 1.22+
- AWS CLI configured
- Git

### 1. Clone and Setup
```bash
git clone https://github.com/your-org/MultitaskProject.git
cd MultitaskProject

# Setup environment
cp .env.example .env
# Edit .env with your actual values

# Install dependencies
npm run setup  # or make setup
```

### 2. Configure AWS Secrets
```bash
cd infra
npm run setup:secrets
```

### 3. Deploy Infrastructure
```bash
# Development
npm run deploy:dev

# Staging
npm run deploy:staging

# Production
npm run deploy:prod
```

## 🔧 Manual Deployment Steps

### Backend Services Deployment

1. **Build All Services**
```bash
# Using Make (Linux/Mac)
make build-all

# Using PowerShell (Windows)
cd services/auth-svc
$env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o ../../bin/auth ./cmd/main.go

# Repeat for other services: profile, chat, post, catalog, ai
```

2. **Deploy with Serverless Framework**
```bash
cd infra
npm install
npx serverless deploy --stage dev
```

### Frontend Deployment

1. **Build Frontend**
```bash
cd frontend
npm install
npm run build
```

2. **Deploy to S3**
```bash
aws s3 sync dist/ s3://multitask-frontend-dev --delete
```

3. **Invalidate CloudFront**
```bash
aws cloudfront create-invalidation --distribution-id YOUR_DISTRIBUTION_ID --paths "/*"
```

## 🌍 Environment-Specific Configurations

### Development Environment
- **Stage**: `dev`
- **Domain**: `api-dev.multitask.com`
- **Features**: Debug logging, development tools enabled
- **Scaling**: Minimal (cost-optimized)

### Staging Environment
- **Stage**: `staging`
- **Domain**: `api-staging.multitask.com`
- **Features**: Production-like, testing enabled
- **Scaling**: Moderate

### Production Environment
- **Stage**: `prod`
- **Domain**: `api.multitask.com`
- **Features**: Full monitoring, auto-scaling, backups
- **Scaling**: High availability

## 📊 Post-Deployment Verification

### Health Checks
```bash
# API Health Check
curl https://api-dev.multitask.com/v1/auth/health

# WebSocket Connection Test
wscat -c wss://api-dev.multitask.com/ws

# Frontend Test
curl https://app-dev.multitask.com
```

### Service Status
```bash
# Check Lambda functions
aws lambda list-functions --query 'Functions[?starts_with(FunctionName, `multitask`)].{Name:FunctionName,Status:State}'

# Check DynamoDB tables
aws dynamodb list-tables --query 'TableNames[?starts_with(@, `multitask`)]'

# Check S3 buckets
aws s3 ls | grep multitask
```

## 🔄 CI/CD Pipeline Deployment

### GitHub Actions
The repository includes automated CI/CD pipeline that:

1. **Triggers on**:
   - Push to `main` → Deploy to Production
   - Push to `staging` → Deploy to Staging
   - Push to `develop` → Deploy to Development
   - Manual workflow dispatch

2. **Pipeline Stages**:
   - Code quality checks
   - Security scanning
   - Unit/Integration tests
   - Build services
   - Deploy infrastructure
   - Deploy frontend
   - Health checks
   - Notifications

### Required Secrets
Add these to GitHub repository secrets:

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_ACCESS_KEY_ID_PROD  # Separate for production
AWS_SECRET_ACCESS_KEY_PROD

# Application Secrets
JWT_SECRET
GEMINI_API_KEY
OPENAI_API_KEY

# Notifications
SLACK_WEBHOOK_URL
```

## 🛠️ Troubleshooting

### Common Issues

1. **Lambda Function Timeout**
```bash
# Increase timeout in serverless.yml
timeout: 60  # seconds
```

2. **DynamoDB Throttling**
```bash
# Switch to on-demand billing
billingMode: PAY_PER_REQUEST
```

3. **CORS Issues**
```bash
# Check CORS settings in serverless.yml
cors:
  origin: https://yourdomain.com
  headers:
    - Content-Type
    - Authorization
```

4. **Cold Start Performance**
```bash
# Enable provisioned concurrency for critical functions
provisionedConcurrency: 2
```

### Logs and Debugging

```bash
# View Lambda logs
serverless logs -f auth --tail

# View all logs
aws logs describe-log-groups --log-group-name-prefix "/aws/lambda/multitask"

# Debug mode
STAGE=dev DEBUG=true serverless deploy
```

## 📈 Scaling and Optimization

### Performance Tuning

1. **Lambda Optimization**
   - Use appropriate memory allocation
   - Enable X-Ray tracing
   - Implement connection pooling

2. **DynamoDB Optimization**
   - Design efficient partition keys
   - Use Global Secondary Indexes wisely
   - Implement caching with ElastiCache

3. **Frontend Optimization**
   - Enable CloudFront compression
   - Implement service workers
   - Use lazy loading

### Monitoring Setup

```bash
# CloudWatch Alarms
aws cloudwatch put-metric-alarm \
  --alarm-name "MultitaskPlatform-HighErrorRate" \
  --alarm-description "High error rate detected" \
  --metric-name Errors \
  --namespace AWS/Lambda \
  --statistic Sum \
  --period 300 \
  --threshold 10 \
  --comparison-operator GreaterThanThreshold
```

## 🔐 Security Considerations

### Production Security Checklist

- [ ] Use AWS WAF for API Gateway
- [ ] Enable CloudTrail logging
- [ ] Rotate secrets regularly
- [ ] Implement least privilege IAM roles
- [ ] Enable VPC endpoints for private resources
- [ ] Use AWS KMS for encryption
- [ ] Set up Security Hub monitoring

### Backup Strategy

```bash
# DynamoDB Point-in-time Recovery
aws dynamodb put-backup-policy \
  --table-name multitask-profiles-prod \
  --backup-policy BackupEnabled=true

# S3 Versioning
aws s3api put-bucket-versioning \
  --bucket multitask-frontend-prod \
  --versioning-configuration Status=Enabled
```

## 📞 Support and Resources

- **Documentation**: [Link to full docs]
- **API Reference**: [Link to API docs]
- **Monitoring Dashboard**: [Link to monitoring]
- **Status Page**: [Link to status page]

## 🔄 Rollback Procedures

### Emergency Rollback

```bash
# Rollback infrastructure
cd infra
serverless rollback --timestamp TIMESTAMP

# Rollback frontend
aws s3 sync s3://multitask-frontend-prod-backup/ s3://multitask-frontend-prod/
aws cloudfront create-invalidation --distribution-id ID --paths "/*"
```

### Database Migration Rollback

```bash
# Restore from backup
aws dynamodb restore-table-from-backup \
  --target-table-name multitask-profiles-prod \
  --backup-arn arn:aws:dynamodb:region:account:table/multitask-profiles-prod/backup/backup-id
```