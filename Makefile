# ====================================================================================
#  Makefile for Go Projects
# ====================================================================================

# ------------------------------------------------------------------------------------
#  Configuration
# ------------------------------------------------------------------------------------

# Binary name for the compiled application
BINARY_NAME=main

# The path to the main package to build/run
CMD_PATH=./cmd/main.go

# Go command
GO=go

# Linker flags for building a smaller binary in production
LDFLAGS=-ldflags="-s -w"

# ====================================================================================
#  Commands
# ====================================================================================

.DEFAULT_GOAL := help

.PHONY: help build dev prod clean

help: ## ✨ Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 📦 Build the application binary
	@echo "Building binary..."
	@$(GO) build $(LDFLAGS) -o ./cmd/$(BINARY_NAME) $(CMD_PATH)
	@echo "Binary created at cmd/$(BINARY_NAME)"

dev: ## 🚀 Run the application in development mode (with hot-reload)
	@echo "Starting dev server with hot-reload (requires 'air')..."
	@echo "Install with: go install github.com/air-verse/air@latest"
	@air -c .air.dev.toml

local: ## 🚀 Run the application in local mode (with hot-reload)
	@echo "Starting dev server with hot-reload (requires 'air')..."
	@echo "Install with: go install github.com/air-verse/air@latest"
	@air -c .air.local.toml

prod: build ## ⚙️  Run the application in production mode
	@echo "Starting application in production mode..."
	@./cmd/$(BINARY_NAME) --mode http --stage prod

clean: ## 🧹 Remove build artifacts
	@echo "Cleaning up..."
	@rm -rf ./cmd/main
	@echo "Done."