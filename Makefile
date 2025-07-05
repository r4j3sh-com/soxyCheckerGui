.PHONY: help dev build build-all build-release test test-coverage test-frontend lint format vet security-scan clean copyright tag release-notes docker-build docker-run setup-hooks update-deps docs bench profile install-tools check-all setup

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := soxyCheckerGui
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

help: ## Show this help message
    @echo "$(BLUE)SoxyChecker GUI - Development Commands$(RESET)"
    @echo ""
    @awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Development
dev: ## Start development server
    @echo "$(YELLOW)Starting development server...$(RESET)"
    wails dev

install-deps: ## Install all dependencies
    @echo "$(YELLOW)Installing Go dependencies...$(RESET)"
    go mod download
    @echo "$(YELLOW)Installing frontend dependencies...$(RESET)"
    cd frontend && npm ci

# Building
build: ## Build for current platform
    @echo "$(YELLOW)Building for current platform...$(RESET)"
    wails build

build-all: ## Build for all platforms
    @echo "$(YELLOW)Building for all platforms...$(RESET)"
    wails build -platform linux/amd64
    wails build -platform windows/amd64
    wails build -platform darwin/amd64
    wails build -platform darwin/arm64

build-release: ## Build optimized release version
    @echo "$(YELLOW)Building optimized release...$(RESET)"
    wails build -clean -s -trimpath

# Testing
test: ## Run all tests
    @echo "$(YELLOW)Running Go tests...$(RESET)"
    go test -v -race ./...

test-coverage: ## Run tests with coverage
    @echo "$(YELLOW)Running tests with coverage...$(RESET)"
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

test-frontend: ## Run frontend tests
    @echo "$(YELLOW)Running frontend tests...$(RESET)"
    cd frontend && npm test

# Code quality
lint: ## Run linters
    @echo "$(YELLOW)Running Go linters...$(RESET)"
    golangci-lint run
    @echo "$(YELLOW)Running frontend linters...$(RESET)"
    cd frontend && npm run lint

format: ## Format code
    @echo "$(YELLOW)Formatting Go code...$(RESET)"
    gofmt -s -w .
    @echo "$(YELLOW)Formatting frontend code...$(RESET)"
    cd frontend && npm run format

vet: ## Run go vet
    @echo "$(YELLOW)Running go vet...$(RESET)"
    go vet ./...

# Security
security-scan: ## Run security scans
    @echo "$(YELLOW)Running security scans...$(RESET)"
    gosec ./...
    @echo "$(YELLOW)Running npm audit...$(RESET)"
    cd frontend && npm audit

# Utilities
clean: ## Clean build artifacts
    @echo "$(YELLOW)Cleaning build artifacts...$(RESET)"
    rm -rf build/
    rm -rf frontend/dist/
    rm -rf coverage.out coverage.html
    @echo "$(GREEN)Clean complete!$(RESET)"

copyright: ## Add copyright headers to source files
    @echo "$(YELLOW)Adding copyright headers...$(RESET)"
    chmod +x scripts/add-copyright.sh
    ./scripts/add-copyright.sh

# Release management
tag: ## Create a new tag (usage: make tag VERSION=v1.0.0)
    @if [ -z "$(VERSION)" ]; then \
        echo "$(RED)Error: VERSION is required. Usage: make tag VERSION=v1.0.0$(RESET)"; \
        exit 1; \
    fi
    @echo "$(YELLOW)Creating tag $(VERSION)...$(RESET)"
    git tag -a $(VERSION) -m "Release $(VERSION)"
    git push origin $(VERSION)
    @echo "$(GREEN)Tag $(VERSION) created and pushed!$(RESET)"

release-notes: ## Generate release notes
    @echo "$(YELLOW)Generating release notes...$(RESET)"
    @git log --pretty=format:"- %s (%h)" $(shell git describe --tags --abbrev=0 HEAD^)..HEAD > RELEASE_NOTES.md
    @echo "$(GREEN)Release notes generated: RELEASE_NOTES.md$(RESET)"

# Docker (optional)
docker-build: ## Build Docker image
    @echo "$(YELLOW)Building Docker image...$(RESET)"
    docker build -t $(APP_NAME) .

docker-run: ## Run Docker container
    @echo "$(YELLOW)Running Docker container...$(RESET)"
    docker run -p 8080:8080 $(APP_NAME)

# Pre-commit setup
setup-hooks: ## Setup pre-commit hooks
    @echo "$(YELLOW)Setting up pre-commit hooks...$(RESET)"
    pip install pre-commit
    pre-commit install
    @echo "$(GREEN)Pre-commit hooks installed!$(RESET)"

# Dependencies update
update-deps: ## Update all dependencies
    @echo "$(YELLOW)Updating Go dependencies...$(RESET)"
    go get -u ./...
    go mod tidy
    @echo "$(YELLOW)Updating frontend dependencies...$(RESET)"
    cd frontend && npm update
    @echo "$(GREEN)Dependencies updated!$(RESET)"

# Documentation
docs: ## Generate documentation
    @echo "$(YELLOW)Starting documentation server...$(RESET)"
    godoc -http=:6060

# Benchmarks
bench: ## Run benchmarks
    @echo "$(YELLOW)Running benchmarks...$(RESET)"
    go test -bench=. -benchmem ./...

# Profile
profile: ## Run with profiling
    @echo "$(YELLOW)Building with profiling...$(RESET)"
    go build -o bin/$(APP_NAME)-profile
    ./bin/$(APP_NAME)-profile -cpuprofile=cpu.prof -memprofile=mem.prof

# Install tools
install-tools: ## Install development tools
    @echo "$(YELLOW)Installing development tools...$(RESET)"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    @echo "$(GREEN)Development tools installed!$(RESET)"

# Check everything before commit
check-all: lint test security-scan ## Run all checks before committing
    @echo "$(GREEN)✅ All checks passed!$(RESET)"

# Quick development setup
setup: install-deps install-tools setup-hooks ## Complete development setup
    @echo "$(GREEN)🚀 Development environment setup complete!$(RESET)"
    @echo "$(BLUE)Run 'make dev' to start development server$(RESET)"

# Build info
info: ## Show build information
    @echo "$(BLUE)Build Information:$(RESET)"
    @echo "  App Name: $(APP_NAME)"
    @echo "  Version: $(VERSION)"
    @echo "  Build Time: $(BUILD_TIME)"
    @echo "  Git Commit: $(GIT_COMMIT)"
