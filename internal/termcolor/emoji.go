package termcolor

import (
	"os"
	"runtime"
	"strings"
)

// SupportsEmoji 检测终端是否支持emoji
func SupportsEmoji() bool {
	// 1. 检查是否支持UTF-8
	if !isUTF8Locale() {
		return false
	}

	// 2. 根据操作系统判断
	switch runtime.GOOS {
	case "darwin":
		// macOS终端基本都支持emoji
		return true
	case "windows":
		// Windows 10+ 基本支持
		return isWindowsTerminalOrModern()
	case "linux":
		// Linux需要检查更多条件
		return isLinuxEmojiSupported()
	default:
		return false
	}
}

// isUTF8Locale 检查是否使用UTF-8编码
func isUTF8Locale() bool {
	// 检查LC_ALL, LC_CTYPE, LANG环境变量
	for _, env := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		value := strings.ToLower(os.Getenv(env))
		if value != "" {
			return strings.Contains(value, "utf-8") || strings.Contains(value, "utf8")
		}
	}
	// 默认认为支持UTF-8（现代系统）
	return true
}

// isWindowsTerminalOrModern 检查是否是Windows Terminal或现代终端
func isWindowsTerminalOrModern() bool {
	// Windows Terminal
	if os.Getenv("WT_SESSION") != "" {
		return true
	}

	// Windows 10+ (version 1903+) 的 ConHost 支持emoji
	// 检查 TERM 环境变量
	term := os.Getenv("TERM")
	if strings.Contains(term, "xterm") || strings.Contains(term, "256color") {
		return true
	}

	// PowerShell 7+
	if os.Getenv("PWSH_VERSION") != "" {
		return true
	}

	// 默认假设Windows 10+支持
	return true
}

// isLinuxEmojiSupported 检查Linux终端是否支持emoji
func isLinuxEmojiSupported() bool {
	term := os.Getenv("TERM")
	colorTerm := os.Getenv("COLORTERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	// 已知支持emoji的终端
	supportedTerminals := []string{
		"xterm-256color",
		"screen-256color",
		"tmux-256color",
		"alacritty",
		"kitty",
		"wezterm",
	}

	for _, supported := range supportedTerminals {
		if strings.Contains(term, supported) {
			return true
		}
	}

	// COLORTERM通常表示真彩色支持，现代终端
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return true
	}

	// 常见的现代终端程序
	modernPrograms := []string{
		"vscode",
		"hyper",
		"iTerm.app",
		"WezTerm",
		"Alacritty",
		"kitty",
	}

	for _, modern := range modernPrograms {
		if strings.Contains(termProgram, modern) {
			return true
		}
	}

	// GNOME Terminal
	if os.Getenv("VTE_VERSION") != "" {
		return true
	}

	// Konsole (KDE)
	if os.Getenv("KONSOLE_VERSION") != "" {
		return true
	}

	// 如果在SSH会话中，较保守
	if os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_CLIENT") != "" {
		return false
	}

	// 默认假设现代Linux发行版支持
	return true
}

// GetFallbackIcon 如果不支持emoji，返回替代字符
func GetFallbackIcon(emoji, fallback string) string {
	if SupportsEmoji() {
		return emoji
	}
	return fallback
}
