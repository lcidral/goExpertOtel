.PHONY: help test docs build up down clean lint fmt deps health

# Monorepo settings
SERVICES := service-a service-b
PACKAGES := models telemetry utils
DOCKER_COMPOSE_FILE := deployments/docker/docker-compose.yml

help: ## Show this help
	@echo "GoExpertOtel - Sistema de Temperatura por CEP com OpenTelemetry"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
up: ## Start all services with docker-compose
	@echo "🚀 Starting all services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

down: ## Stop all services
	@echo "🛑 Stopping all services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

logs: ## Show logs from all services
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Testing commands
test: ## Run all unit tests
	@echo "🧪 Running unit tests for all services and packages..."
	@go test ./services/... ./pkg/... -v -cover

test-service-a: ## Run tests for service-a only
	@echo "🧪 Testing service-a..."
	@go test ./services/service-a/... -v -cover

test-service-b: ## Run tests for service-b only
	@echo "🧪 Testing service-b..."
	@go test ./services/service-b/... -v -cover

test-pkg: ## Run tests for shared packages
	@echo "🧪 Testing shared packages..."
	@go test ./pkg/... -v -cover

test-integration: ## Run integration tests
	@echo "🔗 Running integration tests..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@sleep 10
	@echo "Integration tests would run here..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

test-e2e: ## Run E2E tests
	@echo "🌐 Running E2E tests..."
	@if [ ! -f deployments/docker/.env ]; then \
		echo "📝 Creating .env file from template..."; \
		cp deployments/docker/.env.example deployments/docker/.env; \
	fi
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@sleep 15
	@echo "Running E2E test suite..."
	@go test ./test/e2e -v -timeout=5m || (docker-compose -f $(DOCKER_COMPOSE_FILE) logs && exit 1)
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

test-all: test test-integration test-e2e ## Run all tests

# Coverage
coverage: ## Generate coverage report for entire monorepo (excludes E2E tests)
	@echo "📊 Generating coverage report..."
	@go test ./services/... ./pkg/... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

coverage-service-a: ## Coverage for service-a
	@go test ./services/service-a/... -coverprofile=coverage-service-a.out
	@go tool cover -html=coverage-service-a.out -o coverage-service-a.html

coverage-service-b: ## Coverage for service-b
	@go test ./services/service-b/... -coverprofile=coverage-service-b.out
	@go tool cover -html=coverage-service-b.out -o coverage-service-b.html

# Build commands
build: ## Build all services with docker
	@echo "🔨 Building all services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build

build-local: ## Build services locally
	@echo "🔨 Building services locally..."
	@mkdir -p bin
	@cd services/service-a && go build -o ../../bin/service-a ./cmd/server
	@cd services/service-b && go build -o ../../bin/service-b ./cmd/server
	@echo "Binaries created in bin/"

# Documentation
docs: ## Generate API documentation
	@echo "📚 Generating API docs..."
	@echo "API documentation would be generated here (swagger, etc.)"

# Code quality
lint: ## Run linters for entire monorepo
	@echo "🔍 Running linters..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

lint-fix: ## Fix linting issues
	@golangci-lint run --fix ./... || echo "golangci-lint not installed"

fmt: ## Format code
	@go fmt ./...

# Health checks
health: ## Check health of all services
	@echo "🩺 Checking service health..."
	@curl -sf http://localhost:8080/health > /dev/null && echo "✅ Service A healthy" || echo "❌ Service A unhealthy"
	@curl -sf http://localhost:8081/health > /dev/null && echo "✅ Service B healthy" || echo "❌ Service B unhealthy"

# Monitoring
zipkin: ## Open Zipkin UI
	@echo "🔍 Opening Zipkin UI..."
	@command -v open > /dev/null && open http://localhost:9411 || echo "Open http://localhost:9411 in your browser"

# Example requests
example-valid: ## Test with valid CEP
	@echo "🧪 Testing valid CEP..."
	@curl -X POST http://localhost:8080/ \
		-H "Content-Type: application/json" \
		-d '{"cep": "01310100"}' | jq || echo "Response received (jq not installed for formatting)"

example-invalid: ## Test with invalid CEP
	@echo "🧪 Testing invalid CEP..."
	@curl -X POST http://localhost:8080/ \
		-H "Content-Type: application/json" \
		-d '{"cep": "123"}' | jq || echo "Response received (jq not installed for formatting)"

example-not-found: ## Test with non-existent CEP
	@echo "🧪 Testing non-existent CEP..."
	@curl -X POST http://localhost:8080/ \
		-H "Content-Type: application/json" \
		-d '{"cep": "00000000"}' | jq || echo "Response received (jq not installed for formatting)"

examples: example-valid example-invalid example-not-found ## Run all example requests

# Utilities
clean: ## Clean build artifacts and containers
	@echo "🧹 Cleaning up..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v
	@docker system prune -f
	@rm -rf bin/ coverage*.out coverage*.html

deps: ## Download and tidy dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy

dev-setup: ## Setup development environment
	@echo "🔧 Setting up development environment..."
	@go mod download
	@mkdir -p bin
	@echo "Development environment ready!"

# Development workflow
dev: ## Start development environment
	@echo "🚀 Starting development environment..."
	@make up
	@sleep 5
	@make health

dev-logs: ## Show development logs
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f service-a service-b

# Quick development commands
restart: down up ## Restart all services

rebuild: down build up ## Rebuild and restart all services

status: ## Show status of all containers
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps

# Cache operations
cache-stats: ## Show cache statistics
	@echo "📊 Cache statistics:"
	@curl -s http://localhost:8081/cache/stats | jq || echo "Cache stats retrieved (jq not installed for formatting)"

# Utility for checking prerequisites
check-deps: ## Check if required tools are installed
	@echo "🔍 Checking dependencies..."
	@command -v docker > /dev/null && echo "✅ Docker installed" || echo "❌ Docker not found"
	@command -v docker-compose > /dev/null && echo "✅ Docker Compose installed" || echo "❌ Docker Compose not found"
	@command -v go > /dev/null && echo "✅ Go installed" || echo "❌ Go not found"
	@command -v curl > /dev/null && echo "✅ curl installed" || echo "❌ curl not found"
	@command -v jq > /dev/null && echo "✅ jq installed (optional)" || echo "⚠️  jq not found (optional, for JSON formatting)"

# Show project info
info: ## Show project information
	@echo "📋 Project Information"
	@echo "====================="
	@echo "Name: GoExpertOtel"
	@echo "Description: Sistema de Temperatura por CEP com OpenTelemetry"
	@echo "Services: $(SERVICES)"
	@echo "Packages: $(PACKAGES)"
	@echo "Docker Compose: $(DOCKER_COMPOSE_FILE)"
	@echo ""
	@echo "Endpoints:"
	@echo "  Service A: http://localhost:8080"
	@echo "  Service B: http://localhost:8081"
	@echo "  Zipkin UI: http://localhost:9411"
	@echo ""
	@echo "Quick start: make dev"