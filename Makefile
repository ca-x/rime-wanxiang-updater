.PHONY: build build-release clean test help

# 默认版本号
VERSION ?= dev
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 编译参数
LDFLAGS := -X rime-wanxiang-updater/internal/version.Version=$(VERSION)
OUTPUT := rime-wanxiang-updater

help: ## 显示帮助信息
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用的 targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建开发版本 (VERSION=dev)
	@echo "构建开发版本..."
	go build -ldflags="$(LDFLAGS)" -o $(OUTPUT) ./cmd/rime-wanxiang-updater
	@echo "构建完成: $(OUTPUT) (version: $(VERSION))"

build-release: ## 构建发布版本 (需要设置 VERSION 变量)
	@echo "构建发布版本..."
	go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o $(OUTPUT) ./cmd/rime-wanxiang-updater
	@echo "构建完成: $(OUTPUT) (version: $(VERSION))"

build-all: ## 构建所有平台版本 (需要设置 VERSION 变量)
	@echo "构建所有平台版本 (version: $(VERSION))..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o dist/$(OUTPUT)-linux-amd64 ./cmd/rime-wanxiang-updater
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o dist/$(OUTPUT)-linux-arm64 ./cmd/rime-wanxiang-updater
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o dist/$(OUTPUT)-darwin-amd64 ./cmd/rime-wanxiang-updater
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o dist/$(OUTPUT)-darwin-arm64 ./cmd/rime-wanxiang-updater
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o dist/$(OUTPUT)-windows-amd64.exe ./cmd/rime-wanxiang-updater
	GOOS=windows GOARCH=arm64 go build -trimpath -ldflags="-s -w $(LDFLAGS)" -o dist/$(OUTPUT)-windows-arm64.exe ./cmd/rime-wanxiang-updater
	@echo "所有平台构建完成！"

test: ## 运行测试
	go test -v ./...

clean: ## 清理构建产物
	rm -f $(OUTPUT)
	rm -rf dist/
	@echo "清理完成"

run: ## 运行程序
	go run ./cmd/rime-wanxiang-updater

version: ## 显示版本信息
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
