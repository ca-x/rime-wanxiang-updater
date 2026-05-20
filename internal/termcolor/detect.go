package termcolor

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Background 表示终端背景类型
type Background int

const (
	BackgroundUnknown Background = iota
	BackgroundDark
	BackgroundLight
)

// String 返回背景类型的字符串表示
func (b Background) String() string {
	switch b {
	case BackgroundDark:
		return "dark"
	case BackgroundLight:
		return "light"
	default:
		return "unknown"
	}
}

// IsDark 返回是否为深色背景
func (b Background) IsDark() bool {
	return b == BackgroundDark || b == BackgroundUnknown
}

// DetectBackground 检测终端背景色
func DetectBackground() Background {
	// 1. 检查环境变量 COLORFGBG
	if bg := detectFromColorFgBg(); bg != BackgroundUnknown {
		return bg
	}

	// 2. 检查 TERMINAL_BACKGROUND 环境变量
	if bg := detectFromEnv(); bg != BackgroundUnknown {
		return bg
	}

	// Do not probe OSC 11 here. This is a Bubble Tea app, and reading stdin
	// before Bubble Tea starts can leave a competing reader that consumes user
	// input escape sequences on terminals that do not answer promptly.
	// Users can still force the mode with TERM_BACKGROUND/TERMINAL_BACKGROUND.
	return BackgroundDark
}

// detectFromColorFgBg 从 COLORFGBG 环境变量检测
func detectFromColorFgBg() Background {
	colorFgBg := os.Getenv("COLORFGBG")
	if colorFgBg == "" {
		return BackgroundUnknown
	}

	parts := strings.Split(colorFgBg, ";")
	if len(parts) < 2 {
		return BackgroundUnknown
	}

	bg := parts[len(parts)-1]
	switch bg {
	case "0", "8":
		return BackgroundDark
	case "7", "15":
		return BackgroundLight
	}

	return BackgroundUnknown
}

// detectFromEnv 从环境变量检测
func detectFromEnv() Background {
	// 检查常见的终端背景环境变量
	for _, envVar := range []string{"TERMINAL_BACKGROUND", "TERM_BACKGROUND", "COLORTHEME"} {
		value := strings.ToLower(os.Getenv(envVar))
		if value == "dark" {
			return BackgroundDark
		}
		if value == "light" {
			return BackgroundLight
		}
	}

	return BackgroundUnknown
}

// InitLipgloss 根据终端背景初始化 lipgloss
func InitLipgloss() Background {
	bg := DetectBackground()

	if bg.IsDark() {
		os.Setenv("TERM_BACKGROUND", "dark")
		lipgloss.SetHasDarkBackground(true)
	} else {
		os.Setenv("TERM_BACKGROUND", "light")
		lipgloss.SetHasDarkBackground(false)
	}

	return bg
}

// SetBackground 手动设置背景模式
func SetBackground(dark bool) {
	if dark {
		os.Setenv("TERM_BACKGROUND", "dark")
		lipgloss.SetHasDarkBackground(true)
	} else {
		os.Setenv("TERM_BACKGROUND", "light")
		lipgloss.SetHasDarkBackground(false)
	}
}
