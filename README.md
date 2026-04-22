# Rime 万象输入法更新工具

这是一个用 Go 语言编写的 Rime 万象输入法自动更新工具，支持 Windows、macOS 和 Linux 平台。
> 如果你不想使用万象的方案，倾向于更自由地选输入法方案，推荐使用 https://github.com/ca-x/snout ,在支持万象输入法方案的基础上支持了其他的输入法方案以及自动patch使用万象语言模型

## ✨ 特性

- 🔄 **自动更新**: 支持词库、方案、模型的一键检查、下载和部署
- 🎨 **图形化终端界面**: 使用键盘即可完成更新、配置和主题设置
- 🌐 **中英双语界面**: 主界面与主要功能页支持中文 / English 切换
- 🌍 **跨平台支持**: 原生支持 Windows、macOS 和 Linux
- 🧩 **主题定制**: 支持程序界面主题切换，并按平台提供 Rime / Fcitx5 主题设置能力
- ⚡ **自动部署/重载**: 更新完成后自动部署；主题设置后也可直接在程序内执行部署或重载
- 📦 **多种安装方式**: 支持 AUR (Arch Linux)、Chocolatey (Windows)、Homebrew (macOS)
- 🔌 **代理支持**: 支持 SOCKS5 和 HTTP 代理
- 🪞 **镜像加速**: 支持 CNB 镜像，国内访问更快
- 💾 **断点续传**: 下载支持断点续传，节省流量
- 🔐 **SHA256 校验**: 确保文件完整性和安全性

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

### 1. 一键自动更新

- 自动检查方案、词库、模型是否有新版本
- 按当前配置批量下载并部署
- 适合日常维护，直接作为默认入口使用

### 2. 分项更新

- **词库更新**: 只更新词库，适合追新词条
- **方案更新**: 更新完整方案包，适合升级版本
- **模型更新**: 单独更新模型文件，不影响其他资源

### 3. 配置管理

- 修改下载源、代理、自动更新倒计时
- 管理更新前/更新后 Hook
- 管理排除文件规则
- 管理界面语言与程序主题

### 4. 自定义主题

- 切换程序自身的 TUI 主题
- 在支持的平台为 Rime 写入主题 Patch
- 在 Linux 上安装并设置 Fcitx5 主题

## 🎨 TUI 界面

程序是键盘优先的终端界面，常用操作如下：

- **导航**: 使用数字键 (`1-8`) 或方向键 (↑↓) / vim 键 (`j/k`) 选择主菜单项
- **确认**: 按 Enter 或数字键执行操作
- **退出**: 按 `Q` 或 `Ctrl+C` 退出程序
- **返回**: 在子页面按 `Q` 或 `Esc` 返回上一层/主菜单

## 🧩 主题与自定义

主菜单中的 `6. 自定义` 提供和主题相关的功能。程序会根据当前系统和已检测到的输入法引擎显示可用选项。

### 你可以做什么

- **切换程序界面主题**: 只影响本工具自己的终端界面
- **为 Rime 添加主题 Patch**: 将预设主题写入 `weasel.custom.yaml` 或 `squirrel.custom.yaml`
- **设置当前活动主题**: 从已写入的主题里选择一个作为默认主题
- **安装 Fcitx5 主题**: Linux 下可把内置主题安装到本地并直接设置

### 平台支持矩阵

| 功能 | Windows | macOS | Linux |
|------|---------|-------|-------|
| 程序 TUI 主题切换 | ✅ | ✅ | ✅ |
| Rime 主题 Patch | ✅ 小狼毫 | ✅ 鼠须管 | ❌ |
| Fcitx5 主题安装/设置 | ❌ | ❌ | ✅ 仅 `fcitx5` |

### Rime 主题 Patch

- **Windows**: 检测到 **小狼毫** 时可用，写入 `weasel.custom.yaml`
- **macOS**: 检测到 **鼠须管** 时可用，写入 `squirrel.custom.yaml`
- **Linux**: 不提供这个功能

使用流程分为两步：

1. 先用 `Space` 多选想要写入的主题预设
2. 再从这些已选主题里选择一个默认主题

完成后程序可直接执行重新部署，无需再手动处理。

> 主题 Patch 会直接写入 Rime 用户目录根目录下的 `weasel.custom.yaml` 或 `squirrel.custom.yaml`，不需要手动复制模板文件。

### Linux Fcitx5 主题

仅在 **Linux 且检测到 `fcitx5` 已安装** 时显示。

你只需要：

1. 先多选需要保留的内置主题
2. 分别选择浅色模式和深色模式下的 Fcitx5 主题
3. 让程序自动写入 `Theme` / `DarkTheme` 并启用跟随系统深色模式
4. 按提示决定是否立即重载 `fcitx5`

主题会安装到 `~/.local/share/fcitx5/themes/`。

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
