# ==============================================================================
# Clio Service Makefile
# ==============================================================================

# ------------------------------------------------------------------------------
# Service Metadata Configuration
# Load service metadata from service.json file
# ------------------------------------------------------------------------------
SERVICE_NAME := $(shell cat service.json 2>/dev/null | grep '"name"' | cut -d'"' -f4)
SERVICE_VERSION := $(shell cat service.json 2>/dev/null | grep '"version"' | cut -d'"' -f4)
SERVICE_DESC := $(shell cat service.json 2>/dev/null | grep '"description"' | cut -d'"' -f4)

# Fallback to default values if service.json doesn't exist
ifeq ($(SERVICE_NAME),)
    SERVICE_NAME := unknown-service
endif
ifeq ($(SERVICE_VERSION),)
    SERVICE_VERSION := 0.0.0
endif
ifeq ($(SERVICE_DESC),)
    SERVICE_DESC := No description
endif

# ------------------------------------------------------------------------------
# Build Information
# Generate build-time metadata
# ------------------------------------------------------------------------------
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# ------------------------------------------------------------------------------
# Go Build Configuration
# Metadata package path for ldflags injection
# ------------------------------------------------------------------------------
METADATA_PKG := github.com/muixstudio/clio/services/common/metadata

# ldflags for embedding metadata into binary
LDFLAGS := -ldflags "\
    -X '$(METADATA_PKG).Name=$(SERVICE_NAME)' \
    -X '$(METADATA_PKG).Version=$(SERVICE_VERSION)' \
    -X '$(METADATA_PKG).Description=$(SERVICE_DESC)' \
    -X '$(METADATA_PKG).BuildTime=$(BUILD_TIME)' \
    -X '$(METADATA_PKG).GitCommit=$(GIT_COMMIT)'"

# Output directory for compiled binaries
OUTPUT_DIR := bin

# ------------------------------------------------------------------------------
# Docker Configuration
# Docker registry and image naming conventions
# ------------------------------------------------------------------------------
DOCKER_REGISTRY ?= docker.io
DOCKER_NAMESPACE ?= muixstudio
DOCKER_REPOSITORY ?= clio-$(SERVICE_NAME)
IMAGE_VERSION ?= $(SERVICE_VERSION)
IMAGE_BASE = $(DOCKER_REGISTRY)/$(DOCKER_NAMESPACE)/$(DOCKER_REPOSITORY)
IMAGE_TAG = $(IMAGE_BASE):$(IMAGE_VERSION)

# Docker build arguments
DOCKER_BUILD_ARGS = \
	--build-arg SERVICE_NAME="$(SERVICE_NAME)" \
	--build-arg SERVICE_VERSION=$(SERVICE_VERSION) \
	--build-arg SERVICE_DESC="$(SERVICE_DESC)" \
	--build-arg BUILD_TIME=$(BUILD_TIME) \
	--build-arg GIT_COMMIT=$(GIT_COMMIT)

# ==============================================================================
# Targets
# ==============================================================================

# ------------------------------------------------------------------------------
# Default target
# ------------------------------------------------------------------------------
.DEFAULT_GOAL := help

# ------------------------------------------------------------------------------
# Display build information
# ------------------------------------------------------------------------------
.PHONY: build-info
build-info:
	@echo "==================== Build Information ==================="
	@echo "Service Name:    $(SERVICE_NAME)"
	@echo "Version:         $(SERVICE_VERSION)"
	@echo "Description:     $(SERVICE_DESC)"
	@echo "Build Time:      $(BUILD_TIME)"
	@echo "Git Commit:      $(GIT_COMMIT)"
	@echo "Image Tag:       $(IMAGE_TAG)"
	@echo "========================================================="

# ------------------------------------------------------------------------------
# Clean build artifacts
# ------------------------------------------------------------------------------
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Clean completed"

# ------------------------------------------------------------------------------
# Build multi-platform binaries
# Compile Go binaries for different OS and architectures
# ------------------------------------------------------------------------------
.PHONY: build-binary
build-binary: build-info
	@echo "Building multi-platform binaries..."
	@mkdir -p $(OUTPUT_DIR)

	@echo "Building for Linux amd64..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(SERVICE_NAME)-linux-amd64 ./cmd/main.go

	@echo "Building for Linux arm64..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(SERVICE_NAME)-linux-arm64 ./cmd/main.go

	@echo "Building for MacOS amd64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(SERVICE_NAME)-darwin-amd64 ./cmd/main.go

	@echo "Building for MacOS arm64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(SERVICE_NAME)-darwin-arm64 ./cmd/main.go

	@echo "Multi-platform build completed"
	@ls -lh $(OUTPUT_DIR)

# ------------------------------------------------------------------------------
# Build Docker image for local platform
# Use standard docker build (no multi-arch support)
# ------------------------------------------------------------------------------
.PHONY: build-image
build-image: build-info
	@echo "Building Docker image (local platform)..."
	docker build \
		$(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_TAG) \
		-f ./Dockerfile \
		../../
	@echo "Image build completed: $(IMAGE_TAG)"
	@docker images | grep $(DOCKER_REPOSITORY) || true

# ------------------------------------------------------------------------------
# Build and push multi-architecture Docker image
# Use docker buildx for multi-platform support (amd64 + arm64)
# Note: Requires Docker Buildx and pushes to registry
# ------------------------------------------------------------------------------
.PHONY: build-image-multiarch
build-image-multiarch: build-info
	@echo "Building multi-architecture Docker image..."
	@echo "Warning: Multi-arch images will be pushed directly to the registry"
	docker buildx build \
		$(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_TAG) \
		--platform=linux/amd64,linux/arm64 \
		--push \
		-f ./Dockerfile \
		../../
	@echo "Multi-architecture image pushed to: $(IMAGE_TAG)"

# ------------------------------------------------------------------------------
# Push Docker image to registry
# Push the locally built image to remote registry
# ------------------------------------------------------------------------------
.PHONY: push-image
push-image:
	@echo "Pushing image to registry..."
	docker push $(IMAGE_TAG)
	@echo "Image pushed successfully: $(IMAGE_TAG)"

# ------------------------------------------------------------------------------
# Display help information
# Show all available make targets with descriptions
# ------------------------------------------------------------------------------
.PHONY: help
help:
	@echo "Available Commands:"
	@echo ""
	@echo "  make build-info           - Display build information"
	@echo "  make build-binary         - Build multi-platform binaries (Linux/MacOS, amd64/arm64)"
	@echo "  make build-image          - Build Docker image (local platform)"
	@echo "  make build-image-multiarch- Build and push multi-arch image (amd64 + arm64)"
	@echo "  make push-image           - Push image to remote registry"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make help                 - Display this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build-binary         # Build binaries for all platforms"
	@echo "  make build-image          # Build local Docker image"
	@echo "  make build-image-multiarch# Build and push multi-arch image"
	@echo ""
	@echo "Environment Variables:"
	@echo "  DOCKER_REGISTRY           - Docker registry address (default: docker.io)"
	@echo "  DOCKER_NAMESPACE          - Docker namespace (default: muixstudio)"
	@echo "  IMAGE_VERSION             - Image version (default: read from service.json)"