.PHONY: build install clean version help

# 版本信息 / Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志 / Build flags
LDFLAGS := -X 'cc-provider/cmd.Version=$(VERSION)' \
           -X 'cc-provider/cmd.BuildTime=$(BUILD_TIME)' \
           -X 'cc-provider/cmd.GitCommit=$(GIT_COMMIT)'

# 默认目标 / Default target
all: build

# 构建二进制文件 / Build binary
build:
	@echo "Building cc-provider $(VERSION)..."
	@go build -ldflags "$(LDFLAGS)" -o bin/cc-provider .
	@echo "Build complete: bin/cc-provider"

# 安装到 GOPATH/bin / Install to GOPATH/bin
install:
	@echo "Installing cc-provider $(VERSION)..."
	@go install -ldflags "$(LDFLAGS)" .
	@echo "Installation complete"

# 清理构建产物 / Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete"

# 显示版本信息 / Display version information
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

# 显示帮助信息 / Display help information
help:
	@echo "Available targets:"
	@echo "  make build   - Build the binary with version information"
	@echo "  make install - Install the binary to GOPATH/bin"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make version - Display version information"
	@echo "  make help    - Display this help message"
