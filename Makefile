# Variables
APP_NAME = gollama
BUILD_DIR = builds
BIN_DIR = $(BUILD_DIR)/bin

# Go toolchain
GO = go
GOBUILD = $(GO) build
GOTEST = $(GO) test
GOCLEAN = $(GO) clean

# Default target
all: build

# Build the application for specific platforms
build_win: ## Build for Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/win/$(APP_NAME).exe ./gollama.go

build_lin: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/linux/$(APP_NAME) ./gollama.go

build_mac: ## Build for MacOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/mac/$(APP_NAME) ./gollama.go

# Build ggollama for all platforms
build: ## Build both gollama for all platforms and save in the build folder
	$(MAKE) build_lin
	$(MAKE) build_win
	$(MAKE) build_mac

# Install the application globally
install: build ## Install the application
	cp $(BIN_DIR)/linux/$(APP_NAME) ~/go/bin/

# Run the application
run: build ## Run the application
	./$(BIN_DIR)/$(APP_NAME)

# Run tests
test: ## Run tests
	$(GOTEST) ./...

# Lint the code
lint: ## Lint the code
	golangci-lint run

# Clean the build artifacts
clean: ## Clean the build and logs
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) logs/

# Show the directory tree
tr: ## Generate directory tree
	tree > tree.txt

# Help command
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%-12s\033[0m %s\n\n", "Command", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: all build run test lint clean tr help install
