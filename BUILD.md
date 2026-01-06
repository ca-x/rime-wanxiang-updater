# 构建说明

本项目使用 ldflags 在编译时注入版本号。

## 版本管理

版本号存储在 `internal/version` 包中，默认值为 `dev`。在编译时可以通过 `-ldflags` 参数注入实际的版本号。

## 本地构建

### 使用 Makefile（推荐）

```bash
# 构建开发版本 (version=dev)
make build

# 构建指定版本的发布版本
make build-release VERSION=v1.0.0

# 构建所有平台版本
make build-all VERSION=v1.0.0

# 运行程序
make run

# 运行测试
make test

# 清理构建产物
make clean

# 显示帮助信息
make help
```

### 使用 go build

```bash
# 开发版本（默认 version=dev）
go build -o rime-wanxiang-updater ./cmd/rime-wanxiang-updater

# 指定版本号
go build -ldflags="-X rime-wanxiang-updater/internal/version.Version=v1.0.0" \
  -o rime-wanxiang-updater ./cmd/rime-wanxiang-updater

# 发布版本（压缩优化）
go build -trimpath \
  -ldflags="-s -w -X rime-wanxiang-updater/internal/version.Version=v1.0.0" \
  -o rime-wanxiang-updater ./cmd/rime-wanxiang-updater
```

### 跨平台编译

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags="-X rime-wanxiang-updater/internal/version.Version=v1.0.0" \
  -o rime-wanxiang-updater-linux-amd64 ./cmd/rime-wanxiang-updater

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -ldflags="-X rime-wanxiang-updater/internal/version.Version=v1.0.0" \
  -o rime-wanxiang-updater-darwin-arm64 ./cmd/rime-wanxiang-updater

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags="-X rime-wanxiang-updater/internal/version.Version=v1.0.0" \
  -o rime-wanxiang-updater-windows-amd64.exe ./cmd/rime-wanxiang-updater
```

## GitHub Actions 自动发布

项目配置了 GitHub Actions 自动构建和发布流程。

### 创建发布版本

```bash
# 创建并推送 tag
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 会自动：
1. 构建所有平台版本（Windows/macOS/Linux，AMD64/ARM64）
2. 使用 tag 名称作为版本号注入到二进制文件中
3. 创建 GitHub Release
4. 上传所有构建产物到 Release

### 支持的平台

- Windows: AMD64, ARM64
- macOS: AMD64 (Intel), ARM64 (Apple Silicon)
- Linux: AMD64, ARM64

## ldflags 说明

项目使用以下 ldflags 参数：

- `-X rime-wanxiang-updater/internal/version.Version=<版本号>`: 注入版本号
- `-s`: 去除符号表
- `-w`: 去除 DWARF 调试信息
- `-trimpath`: 去除文件系统路径信息

发布版本会应用所有优化参数以减小二进制文件大小。
