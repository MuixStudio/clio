# Makefile for building with embedded metadata

# 从 service.json 读取元数据
SERVICE_NAME := $(shell cat service.json 2>/dev/null | grep '"name"' | cut -d'"' -f4)
SERVICE_VERSION := $(shell cat service.json 2>/dev/null | grep '"version"' | cut -d'"' -f4)
SERVICE_DESC := $(shell cat service.json 2>/dev/null | grep '"description"' | cut -d'"' -f4)

# 如果 service.json 不存在，使用默认值
ifeq ($(SERVICE_NAME),)
    SERVICE_NAME := unknown-service
endif
ifeq ($(SERVICE_VERSION),)
    SERVICE_VERSION := 0.0.0
endif
ifeq ($(SERVICE_DESC),)
    SERVICE_DESC := No description
endif

# 构建信息
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | cut -d' ' -f3)

# 包路径（根据你的项目调整）
METADATA_PKG := github.com/muixstudio/clio/services/web/metadata

# ldflags 构建参数
LDFLAGS := -ldflags "\
    -X '$(METADATA_PKG).Name=$(SERVICE_NAME)' \
    -X '$(METADATA_PKG).Version=$(SERVICE_VERSION)' \
    -X '$(METADATA_PKG).Description=$(SERVICE_DESC)' \
    -X '$(METADATA_PKG).BuildTime=$(BUILD_TIME)' \
    -X '$(METADATA_PKG).GitCommit=$(GIT_COMMIT)' \
    -X '$(METADATA_PKG).GitBranch=$(GIT_BRANCH)' \
    -X '$(METADATA_PKG).GoVersion=$(GO_VERSION)'"

# 输出目录
OUTPUT_DIR := bin
BINARY_NAME := clio

.PHONY: all build build-info clean test help install

# 默认目标
all: build

# 显示构建信息
build-info:
	@echo "=== Build Information ==="
	@echo "Service Name:    $(SERVICE_NAME)"
	@echo "Version:         $(SERVICE_VERSION)"
	@echo "Description:     $(SERVICE_DESC)"
	@echo "Build Time:      $(BUILD_TIME)"
	@echo "Git Commit:      $(GIT_COMMIT)"
	@echo "Git Branch:      $(GIT_BRANCH)"
	@echo "Go Version:      $(GO_VERSION)"
	@echo "========================="

# 构建二进制文件
build: build-info
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./main.go
	@echo "✅ Build completed: $(OUTPUT_DIR)/$(BINARY_NAME)"

# 构建并运行
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	$(OUTPUT_DIR)/$(BINARY_NAME) run

# 显示版本信息
version: build
	@$(OUTPUT_DIR)/$(BINARY_NAME) version

# 清理构建产物
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(OUTPUT_DIR)
	@echo "✅ Cleaned"

# 运行测试
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# 安装到 GOPATH
install: build-info
	@echo "📦 Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./main.go
	@echo "✅ Installed"

# 构建多平台版本
build-all: build-info
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(OUTPUT_DIR)
	
	@echo "Building for Linux amd64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 ./main.go
	
	@echo "Building for Linux arm64..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 ./main.go
	
	@echo "Building for macOS amd64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 ./main.go
	
	@echo "Building for macOS arm64 (M1/M2)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 ./main.go
	
	@echo "Building for Windows amd64..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe ./main.go
	
	@echo "✅ Multi-platform build completed"
	@ls -lh $(OUTPUT_DIR)

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  make build       - Build the binary with embedded metadata"
	@echo "  make run         - Build and run the application"
	@echo "  make version     - Build and show version info"
	@echo "  make build-info  - Show build information"
	@echo "  make build-all   - Build for multiple platforms"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make install     - Install to GOPATH"
	@echo "  make help        - Show this help"