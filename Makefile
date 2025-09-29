# 🚀 Multitask Platform - Build and Development Makefile

.PHONY: help build build-all clean test test-unit test-integration deps lint format docker-build docker-push deploy setup dev

# Default target
help: ## Show this help message
	@echo "🚀 Multitask Platform Build System"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# Variables
GO_VERSION := 1.22
SERVICES := auth profile chat post catalog ai
STAGE := dev
REGION := us-east-1

# Build configuration
BUILD_DIR := bin
LDFLAGS := -ldflags="-s -w"
CGO_ENABLED := 0
GOOS := linux
GOARCH := amd64

# Colors for output
COLOR_RESET = \033[0m
COLOR_BLUE = \033[34m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m
COLOR_RED = \033[31m

## 🔨 Build Commands

build: ## Build all services for Lambda deployment
	@echo "$(COLOR_BLUE)🔨 Building all Go services for AWS Lambda...$(COLOR_RESET)"
	@$(MAKE) build-all

build-all: clean $(addprefix build-, $(SERVICES)) ## Build all services
	@echo "$(COLOR_GREEN)✅ All services built successfully!$(COLOR_RESET)"

build-%: ## Build specific service (e.g., make build-auth)
	@echo "$(COLOR_BLUE)🔨 Building $* service...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	@cd services/$*-svc && \
		CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build $(LDFLAGS) -o ../../$(BUILD_DIR)/$* ./cmd/main.go
	@echo "$(COLOR_GREEN)✅ Built $* service -> $(BUILD_DIR)/$*$(COLOR_RESET)"

## 🧪 Testing Commands

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo "$(COLOR_BLUE)🧪 Running unit tests...$(COLOR_RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✅ Unit tests completed! Coverage report: coverage.html$(COLOR_RESET)"

test-integration: ## Run integration tests
	@echo "$(COLOR_BLUE)🧪 Running integration tests...$(COLOR_RESET)"
	@go test -v -tags=integration ./...
	@echo "$(COLOR_GREEN)✅ Integration tests completed!$(COLOR_RESET)"

test-%: ## Run tests for specific service
	@echo "$(COLOR_BLUE)🧪 Testing $* service...$(COLOR_RESET)"
	@cd services/$*-svc && go test -v -race ./...
	@echo "$(COLOR_GREEN)✅ $* service tests completed!$(COLOR_RESET)"

benchmark: ## Run benchmarks
	@echo "$(COLOR_BLUE)📊 Running benchmarks...$(COLOR_RESET)"
	@go test -bench=. -benchmem ./...

## 📦 Dependencies

deps: ## Download and tidy dependencies
	@echo "$(COLOR_BLUE)📦 Managing Go dependencies...$(COLOR_RESET)"
	@go mod download
	@go mod tidy
	@go mod verify
	@echo "$(COLOR_GREEN)✅ Dependencies updated!$(COLOR_RESET)"

deps-update: ## Update all dependencies to latest versions
	@echo "$(COLOR_BLUE)📦 Updating dependencies...$(COLOR_RESET)"
	@go get -u ./...
	@go mod tidy
	@echo "$(COLOR_GREEN)✅ Dependencies updated to latest versions!$(COLOR_RESET)"

deps-vendor: ## Create vendor directory
	@echo "$(COLOR_BLUE)📦 Creating vendor directory...$(COLOR_RESET)"
	@go mod vendor
	@echo "$(COLOR_GREEN)✅ Vendor directory created!$(COLOR_RESET)"

## 🔍 Code Quality

lint: ## Run linting tools
	@echo "$(COLOR_BLUE)🔍 Running linting tools...$(COLOR_RESET)"
	@golangci-lint run ./...
	@echo "$(COLOR_GREEN)✅ Linting completed!$(COLOR_RESET)"

format: ## Format code
	@echo "$(COLOR_BLUE)🎨 Formatting code...$(COLOR_RESET)"
	@gofmt -w .
	@goimports -w .
	@echo "$(COLOR_GREEN)✅ Code formatted!$(COLOR_RESET)"

security: ## Run security scanning
	@echo "$(COLOR_BLUE)🔒 Running security scan...$(COLOR_RESET)"
	@gosec ./...
	@echo "$(COLOR_GREEN)✅ Security scan completed!$(COLOR_RESET)"

## 🧹 Cleanup

clean: ## Clean build artifacts
	@echo "$(COLOR_BLUE)🧹 Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✅ Cleanup completed!$(COLOR_RESET)"

clean-all: clean ## Clean all artifacts including vendor
	@echo "$(COLOR_BLUE)🧹 Deep cleaning...$(COLOR_RESET)"
	@rm -rf vendor/
	@go clean -cache -modcache -testcache
	@echo "$(COLOR_GREEN)✅ Deep cleanup completed!$(COLOR_RESET)"

## 🐳 Docker Commands

docker-build: ## Build Docker images for all services
	@echo "$(COLOR_BLUE)🐳 Building Docker images...$(COLOR_RESET)"
	@for service in $(SERVICES); do \
		echo "Building $$service-svc image..."; \
		docker build -t multitask-$$service:latest -f services/$$service-svc/Dockerfile .; \
	done
	@echo "$(COLOR_GREEN)✅ Docker images built!$(COLOR_RESET)"

docker-push: ## Push Docker images to registry
	@echo "$(COLOR_BLUE)🐳 Pushing Docker images...$(COLOR_RESET)"
	@for service in $(SERVICES); do \
		echo "Pushing $$service-svc image..."; \
		docker push multitask-$$service:latest; \
	done
	@echo "$(COLOR_GREEN)✅ Docker images pushed!$(COLOR_RESET)"

## 🚀 Deployment Commands

deploy: build-all ## Deploy to AWS using Serverless Framework
	@echo "$(COLOR_BLUE)🚀 Deploying to AWS (stage: $(STAGE))...$(COLOR_RESET)"
	@cd infra && npm run deploy:$(STAGE)
	@echo "$(COLOR_GREEN)✅ Deployment completed!$(COLOR_RESET)"

deploy-dev: ## Deploy to development environment
	@$(MAKE) deploy STAGE=dev

deploy-staging: ## Deploy to staging environment  
	@$(MAKE) deploy STAGE=staging

deploy-prod: ## Deploy to production environment
	@$(MAKE) deploy STAGE=prod

undeploy: ## Remove deployment from AWS
	@echo "$(COLOR_YELLOW)⚠️  Removing deployment (stage: $(STAGE))...$(COLOR_RESET)"
	@cd infra && serverless remove --stage $(STAGE)
	@echo "$(COLOR_GREEN)✅ Deployment removed!$(COLOR_RESET)"

## 🛠️ Development Commands

setup: ## Initial project setup
	@echo "$(COLOR_BLUE)🛠️  Setting up development environment...$(COLOR_RESET)"
	@echo "Installing Go dependencies..."
	@$(MAKE) deps
	@echo "Installing infrastructure dependencies..."
	@cd infra && npm install
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "$(COLOR_GREEN)✅ Development environment setup completed!$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_YELLOW)📝 Next steps:$(COLOR_RESET)"
	@echo "1. Configure AWS credentials: aws configure"
	@echo "2. Set up secrets: cd infra && npm run setup:secrets"
	@echo "3. Build services: make build"
	@echo "4. Deploy to dev: make deploy-dev"

dev: ## Start local development environment
	@echo "$(COLOR_BLUE)🔧 Starting local development environment...$(COLOR_RESET)"
	@cd infra && npm run offline &
	@cd frontend && npm run dev &
	@echo "$(COLOR_GREEN)✅ Development servers started!$(COLOR_RESET)"
	@echo "API: http://localhost:3000"
	@echo "Frontend: http://localhost:5173"

stop-dev: ## Stop local development servers
	@echo "$(COLOR_BLUE)⏹️  Stopping development servers...$(COLOR_RESET)"
	@pkill -f "serverless offline" || true
	@pkill -f "vite" || true
	@echo "$(COLOR_GREEN)✅ Development servers stopped!$(COLOR_RESET)"

## 📊 Monitoring Commands

logs: ## Tail AWS Lambda logs
	@echo "$(COLOR_BLUE)📋 Tailing Lambda logs...$(COLOR_RESET)"
	@cd infra && serverless logs -f auth --tail

logs-%: ## Tail logs for specific service
	@echo "$(COLOR_BLUE)📋 Tailing $* service logs...$(COLOR_RESET)"
	@cd infra && serverless logs -f $* --tail

info: ## Show deployment info
	@echo "$(COLOR_BLUE)ℹ️  Deployment information:$(COLOR_RESET)"
	@cd infra && serverless info --stage $(STAGE)

## 🔧 Utility Commands

check-tools: ## Check if required tools are installed
	@echo "$(COLOR_BLUE)🔧 Checking required tools...$(COLOR_RESET)"
	@command -v go >/dev/null 2>&1 || { echo "$(COLOR_RED)❌ Go is not installed$(COLOR_RESET)"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "$(COLOR_RED)❌ Node.js is not installed$(COLOR_RESET)"; exit 1; }
	@command -v aws >/dev/null 2>&1 || { echo "$(COLOR_RED)❌ AWS CLI is not installed$(COLOR_RESET)"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "$(COLOR_RED)❌ Docker is not installed$(COLOR_RESET)"; exit 1; }
	@echo "$(COLOR_GREEN)✅ All required tools are installed!$(COLOR_RESET)"

version: ## Show version information
	@echo "$(COLOR_BLUE)📋 Version Information:$(COLOR_RESET)"
	@echo "Go version: $$(go version)"
	@echo "Node version: $$(node --version)"
	@echo "NPM version: $$(npm --version)"
	@echo "AWS CLI version: $$(aws --version)"
	@echo "Docker version: $$(docker --version)"

size: ## Show binary sizes
	@echo "$(COLOR_BLUE)📏 Binary sizes:$(COLOR_RESET)"
	@ls -lh $(BUILD_DIR)/ 2>/dev/null || echo "No binaries found. Run 'make build' first."

## 📈 Performance

profile-cpu: ## Run CPU profiling
	@echo "$(COLOR_BLUE)📊 Running CPU profiling...$(COLOR_RESET)"
	@go test -cpuprofile cpu.prof -bench . ./...
	@go tool pprof cpu.prof

profile-mem: ## Run memory profiling
	@echo "$(COLOR_BLUE)📊 Running memory profiling...$(COLOR_RESET)"
	@go test -memprofile mem.prof -bench . ./...
	@go tool pprof mem.prof

## 🎯 Quick Actions (Common workflows)

quick-deploy: clean build deploy ## Quick clean, build, and deploy
	@echo "$(COLOR_GREEN)🎯 Quick deployment completed!$(COLOR_RESET)"

quick-test: format lint test ## Quick format, lint, and test
	@echo "$(COLOR_GREEN)🎯 Quick testing completed!$(COLOR_RESET)"

release: clean format lint test build ## Prepare for release
	@echo "$(COLOR_GREEN)🎯 Release preparation completed!$(COLOR_RESET)"

# Default Go commands for convenience
mod-init: ## Initialize Go module
	@go mod init github.com/multitask-platform/backend

mod-tidy: ## Tidy Go modules
	@go mod tidy

# Help command that shows more detailed information
help-detailed: ## Show detailed help with examples
	@echo "🚀 Multitask Platform - Detailed Build System Help"
	@echo ""
	@echo "Common Workflows:"
	@echo "  Development setup:    make setup"
	@echo "  Start development:    make dev"
	@echo "  Quick testing:        make quick-test"
	@echo "  Build and deploy:     make quick-deploy"
	@echo "  Deploy to production: make deploy-prod"
	@echo ""
	@echo "Service-specific commands:"
	@echo "  Build auth service:   make build-auth"
	@echo "  Test chat service:    make test-chat"
	@echo "  Show auth logs:       make logs-auth"
	@echo ""
	@echo "Environment variables:"
	@echo "  STAGE     - Deployment stage (dev/staging/prod)"
	@echo "  REGION    - AWS region (default: us-east-1)"
	@echo ""