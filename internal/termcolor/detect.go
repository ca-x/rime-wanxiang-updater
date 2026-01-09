package termcolor

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	// 3. 尝试使用 OSC 11 查询
	if bg := detectFromOSC11(); bg != BackgroundUnknown {
		return bg
	}

	// 4. 默认返回深色（大多数开发者使用深色终端）
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

// detectFromOSC11 使用 OSC 11 转义序列查询终端背景色
func detectFromOSC11() Background {
	oldState, err := makeRaw()
	if err != nil {
		return BackgroundUnknown
	}
	defer restore(oldState)

	// 发送 OSC 11 查询
	fmt.Print("\x1b]11;?\x1b\\")

	// 读取响应
	response := make([]byte, 256)
	done := make(chan int)

	go func() {
		n, _ := os.Stdin.Read(response)
		done <- n
	}()

	select {
	case n := <-done:
		if n > 0 {
			return parseOSC11Response(string(response[:n]))
		}
	case <-time.After(100 * time.Millisecond):
		return BackgroundUnknown
	}

	return BackgroundUnknown
}

// parseOSC11Response 解析 OSC 11 响应
func parseOSC11Response(response string) Background {
	// 查找 "rgb:" 或 "rgba:"
	idx := strings.Index(response, "rgb:")
	if idx == -1 {
		return BackgroundUnknown
	}

	colorPart := response[idx+4:]
	endIdx := strings.IndexAny(colorPart, "\x1b\x07\\")
	if endIdx == -1 {
		return BackgroundUnknown
	}

	rgbStr := colorPart[:endIdx]
	parts := strings.Split(rgbStr, "/")
	if len(parts) != 3 {
		return BackgroundUnknown
	}

	// 解析 RGB 值（取前两位十六进制）
	var r, g, b int
	if len(parts[0]) >= 2 {
		fmt.Sscanf(parts[0][:2], "%x", &r)
	}
	if len(parts[1]) >= 2 {
		fmt.Sscanf(parts[1][:2], "%x", &g)
	}
	if len(parts[2]) >= 2 {
		fmt.Sscanf(parts[2][:2], "%x", &b)
	}

	// 计算相对亮度 Y = 0.2126*R + 0.7152*G + 0.0722*B
	luminance := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)

	if luminance < 128 {
		return BackgroundDark
	}
	return BackgroundLight
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
