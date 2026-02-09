# Rime 万象输入法更新工具

这是一个用 Go 语言编写的 Rime 万象输入法自动更新工具，支持 Windows、macOS 和 Linux 平台。

## ✨ 特性

- 🎨 **精美的 TUI 界面**: 使用 Bubble Tea 和 Lipgloss 构建的现代化终端界面
- 🔄 **自动更新**: 支持词库、方案、模型的自动检测和更新
- 🌍 **跨平台支持**: 原生支持 Windows、macOS 和 Linux
- 📦 **多种安装方式**: 支持 AUR (Arch Linux)、Chocolatey (Windows)、Homebrew (macOS)
- 🚀 **自动化发布**: GitHub Actions 自动构建多平台二进制文件并发布到包管理器
- 🔌 **代理支持**: 支持 SOCKS5 和 HTTP 代理
- 🪞 **镜像加速**: 支持 CNB 镜像，国内访问更快
- 💾 **断点续传**: 下载支持断点续传，节省流量
- 🔐 **SHA256 校验**: 确保文件完整性和安全性

## 项目结构

```
rime-wanxiang-updater/
├── cmd/
│   └── rime-wanxiang-updater/
│       └── main.go                 # 程序入口
├── internal/
│   ├── types/                      # 类型定义
│   │   └── types.go                # 核心数据结构
│   ├── fileutil/                   # 文件工具
│   │   ├── hash.go                 # SHA256 哈希计算
│   │   ├── download.go             # 文件下载（支持断点续传）
│   │   └── extract.go              # ZIP 压缩包解压
│   ├── api/                        # API 客户端
│   │   ├── client.go               # HTTP 客户端（支持代理）
│   │   ├── github.go               # GitHub API（带重试机制）
│   │   └── cnb.go                  # CNB 镜像 API
│   ├── deployer/                   # 部署管理（平台特定）
│   │   ├── deployer.go             # 部署接口
│   │   ├── windows.go              # Windows 部署 (//go:build windows)
│   │   ├── darwin.go               # macOS 部署 (//go:build darwin)
│   │   └── linux.go                # Linux 部署 (//go:build linux)
│   ├── config/                     # 配置管理
│   │   ├── config.go               # 配置读写和管理
│   │   ├── paths_windows.go        # Windows 路径检测
│   │   ├── paths_darwin.go         # macOS 路径检测
│   │   └── paths_linux.go          # Linux 路径检测
│   ├── updater/                    # 更新器
│   │   ├── base.go                 # 基础更新器
│   │   ├── scheme.go               # 方案更新器
│   │   ├── dict.go                 # 词库更新器
│   │   ├── model.go                # 模型更新器
│   │   └── combined.go             # 组合更新器
│   └── ui/                         # 用户界面
│       ├── model.go                # Bubble Tea 模型
│       └── styles.go               # Lipgloss 样式定义
└── .github/
    └── workflows/
        └── release.yml             # 自动发布工作流
```

## 📦 安装

### Arch Linux (AUR)

```bash
# 使用 yay
yay -S rime-wanxiang-updater

# 使用 paru
paru -S rime-wanxiang-updater

# 手动安装
git clone https://aur.archlinux.org/rime-wanxiang-updater.git
cd rime-wanxiang-updater
makepkg -si
```

AUR 包页面：https://aur.archlinux.org/packages/rime-wanxiang-updater

### Windows (Chocolatey)

```powershell
# 安装
choco install rime-wanxiang-updater

# 升级到最新版本
choco upgrade rime-wanxiang-updater

# 卸载
choco uninstall rime-wanxiang-updater
```

> **注意**: 首次使用 Chocolatey 需要先[安装 Chocolatey](https://chocolatey.org/install)

Chocolatey 包页面：https://community.chocolatey.org/packages/rime-wanxiang-updater

### macOS (Homebrew)

```bash
# 添加 tap
brew tap tinypkg/tap

# 安装
brew install rime-wanxiang-updater
```
> 该tap更多的软件列表，请参考 https://github.com/tinypkg/homebrew-tap
### 其他平台 - 下载预编译版本

从 [Releases](https://github.com/ca-x/rime-wanxiang-updater/releases) 页面下载适合您系统的版本：

- **Windows**: `rime-wanxiang-updater-windows-amd64.exe`
- **macOS (Apple Silicon)**: `rime-wanxiang-updater-darwin-arm64`
- **macOS (Intel)**: `rime-wanxiang-updater-darwin-amd64`
- **Linux**: `rime-wanxiang-updater-linux-amd64`

## 🚀 快速开始

### 运行程序

```bash
# Windows
.\rime-wanxiang-updater.exe

# macOS/Linux (需要先添加执行权限)
chmod +x rime-wanxiang-updater
./rime-wanxiang-updater
```

### 首次运行

首次运行会启动配置向导，引导您完成初始设置：

1. 选择方案版本（基础版 / 增强版）
2. 如选择增强版，选择辅助码方案
3. 自动获取并保存配置

## 📦 核心功能

### 1. 词库更新
- 自动检测词库版本
- 增量更新，只下载变化的文件
- SHA256 校验确保完整性
- 自动部署到输入法目录

### 2. 方案更新
- 支持多种辅助码方案
- 自动清理旧的 build 文件
- 智能文件替换

### 3. 模型更新
- 下载最新语言模型
- 自动部署到指定目录

### 4. 自动更新
- 一键检测所有组件更新
- 批量下载和部署
- 完成后自动重启输入法

### 5. 配置管理
- JSON 格式配置文件
- 支持代理设置
- 支持镜像源切换
- 支持文件排除规则

## 🎨 TUI 界面

程序使用 Bubble Tea 和 Lipgloss 构建精美的终端界面：

- **导航**: 使用数字键 (1-6) 或方向键 (↑↓) / vim 键 (j/k) 选择菜单项
- **确认**: 按 Enter 或数字键执行操作
- **退出**: 按 q 或 Ctrl+C 退出程序
- **返回**: 在子页面按 q 或 ESC 返回主菜单

## 🔧 配置文件

配置文件位置：

- **Windows**: `%APPDATA%\rime-updater\config.json`
- **macOS**: `~/Library/Application Support/rime-updater/config.json`
- **Linux**: `~/.config/rime-updater/config.json`

配置示例：

```json
{
  "engine": "weasel",
  "scheme_type": "pro",
  "scheme_file": "wanxiang-xhup-fuzhu.zip",
  "dict_file": "wanxiang-xhup-dicts.zip",
  "use_mirror": false,
  "github_token": "",
  "exclude_files": [".DS_Store", ".git"],
  "auto_update": false,
  "proxy_enabled": false,
  "proxy_type": "socks5",
  "proxy_address": "127.0.0.1:1080"
}
```

## 🛠️ 开发指南

### 环境要求

- Go 1.21 或更高版本
- Git

### 依赖库

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles/progress@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/cloudflare/backoff@latest
go get golang.org/x/net/proxy
```

### 本地构建

```bash
# 克隆仓库
git clone https://github.com/your-username/rime-wanxiang-updater.git
cd rime-wanxiang-updater

# 安装依赖
go mod download

# 构建当前平台
go build -o rime-wanxiang-updater ./cmd/rime-wanxiang-updater

# 交叉编译
GOOS=windows GOARCH=amd64 go build -o rime-wanxiang-updater.exe ./cmd/rime-wanxiang-updater
GOOS=darwin GOARCH=arm64 go build -o rime-wanxiang-updater-mac ./cmd/rime-wanxiang-updater
GOOS=linux GOARCH=amd64 go build -o rime-wanxiang-updater-linux ./cmd/rime-wanxiang-updater
```

### 发布新版本

```bash
# 创建并推送 tag
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 会自动：
1. 构建所有平台的二进制文件
2. 创建 GitHub Release
3. 发布到 AUR (Arch Linux)
4. 发布到 Chocolatey (Windows)
5. 上传编译好的文件

## 🏗️ 架构设计

### 模块化设计

- **types**: 中央类型定义，避免循环依赖
- **fileutil**: 文件操作工具集，可独立测试
- **api**: API 客户端，支持重试和代理
- **deployer**: 平台特定部署逻辑，使用构建约束隔离
- **config**: 配置管理，支持平台特定路径检测
- **updater**: 更新器模块，实现单一职责原则
- **ui**: 界面层，与业务逻辑解耦

### 平台构建约束

使用 Go 的构建标签 (`//go:build`) 实现平台特定代码：

```go
//go:build windows
// Windows 特定代码

//go:build darwin
// macOS 特定代码

//go:build linux
// Linux 特定代码
```

### 重试机制

使用 Cloudflare Backoff 库实现指数退避重试：

- 初始延迟: 1 秒
- 最大延迟: 10 秒
- 最大重试次数: 3 次

## 🌟 技术亮点

1. **优雅的错误处理**: 所有错误都带有上下文信息
2. **代码复用**: 基础更新器模式减少重复代码
3. **类型安全**: 充分利用 Go 的类型系统
4. **可测试性**: 接口驱动设计，便于单元测试
5. **可维护性**: 清晰的模块划分和文档注释

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

待定

## 🙏 致谢

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 优秀的 TUI 框架
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - 精美的样式库
- [Cloudflare Backoff](https://github.com/cloudflare/backoff) - 可靠的重试机制
- [Rime万象更新工具](https://github.com/rimeinn/rime-wanxiang-update-tools)
- [Rime 万象输入法](https://github.com/amzxyz/rime_wanxiang)
