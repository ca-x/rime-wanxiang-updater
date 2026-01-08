# Rime 安装检测模块

## 功能说明

该模块用于检测系统中是否安装了 Rime 输入法，并在未安装时向用户提供友好的安装建议。

## 特性

- ✅ 跨平台支持（macOS、Linux、Windows）
- ✅ 不阻塞程序运行
- ✅ 提供平台特定的安装建议
- ✅ 在向导和主菜单中自动显示提示

## 检测逻辑

### macOS
- 检查 `/Library/Input Methods/Squirrel.app` 是否存在
- 提供 Homebrew 安装命令和官网下载链接

### Linux
- 检查多个可能的 Rime 配置目录：
  - `~/.local/share/fcitx5/rime` (fcitx5-rime)
  - `/usr/share/rime-data` (系统数据)
  - `~/.config/ibus/rime` (ibus-rime)
  - `~/.config/fcitx/rime` (fcitx-rime)
- 检查 `fcitx5` 和 `ibus-daemon` 可执行文件
- 提供各主流发行版的安装命令

### Windows
- 检查注册表项：
  - `HKCU\Software\Rime\Weasel`
  - `HKLM\SOFTWARE\WOW6432Node\Rime\Weasel`
- 检查默认用户数据目录 `%APPDATA%\Rime`
- 提供官网下载链接和 Scoop 安装命令

## 使用方式

```go
import "rime-wanxiang-updater/internal/detector"

// 检测 Rime 安装状态
status := detector.CheckRimeInstallation()

if !status.Installed {
    fmt.Println(status.Message)
}
```

## 输出示例

### macOS 未安装时
```
⚠️  未检测到 Rime 输入法（鼠须管）

安装方式：
  1. 使用 Homebrew: brew install --cask squirrel
  2. 从官网下载: https://rime.im

提示：程序仍可正常运行，但需要先安装 Rime 才能使用更新的配置。
```

### Linux 未安装时
```
⚠️  未检测到 Rime 输入法

安装方式：
  • Debian/Ubuntu: sudo apt install fcitx5-rime 或 ibus-rime
  • Fedora: sudo dnf install fcitx5-rime 或 ibus-rime
  • Arch Linux: sudo pacman -S fcitx5-rime 或 ibus-rime
  • 从官网下载: https://rime.im

提示：程序仍可正常运行，但需要先安装 Rime 才能使用更新的配置。
```

## 测试

运行测试：
```bash
go test -v ./internal/detector/
```
