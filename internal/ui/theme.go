package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// DetectTerminalBackground 尝试检测终端背景色是深色还是浅色
// 返回 true 表示深色背景，false 表示浅色背景
func DetectTerminalBackground() bool {
	// 1. 检查环境变量 COLORFGBG (常见于许多终端)
	// 格式通常为 "foreground;background"，其中 0-7 是暗色，8-15 是亮色
	if colorFgBg := os.Getenv("COLORFGBG"); colorFgBg != "" {
		parts := strings.Split(colorFgBg, ";")
		if len(parts) >= 2 {
			bg := parts[len(parts)-1]
			// 0-7 或 "0" 被认为是深色背景
			// 15 或 "15" 被认为是浅色背景
			if bg == "0" || bg == "8" || bg == "7" {
				return true // 深色背景
			}
			if bg == "15" || bg == "7" {
				return false // 浅色背景
			}
		}
	}

	// 2. 检查 TERM_PROGRAM (某些终端设置)
	termProgram := os.Getenv("TERM_PROGRAM")
	if termProgram == "Apple_Terminal" || termProgram == "iTerm.app" {
		// macOS 终端默认是浅色背景，但这不可靠
		// 我们继续尝试其他方法
	}

	// 3. 尝试使用 OSC 11 查询终端背景色
	// 这是最准确的方法，但不是所有终端都支持
	bgColor := queryTerminalBackgroundColor()
	if bgColor != "" {
		// 解析 RGB 值并判断亮度
		isDark := isColorDark(bgColor)
		return isDark
	}

	// 4. 如果都失败，默认假设是深色背景（因为大多数开发者使用深色终端）
	return true
}

// queryTerminalBackgroundColor 使用 OSC 11 转义序列查询终端背景色
func queryTerminalBackgroundColor() string {
	// 保存当前终端设置
	oldState, err := makeRaw()
	if err != nil {
		return ""
	}
	defer restore(oldState)

	// 发送 OSC 11 查询序列
	fmt.Print("\x1b]11;?\x1b\\")

	// 读取响应（设置超时）
	response := make([]byte, 256)
	done := make(chan int)

	go func() {
		n, _ := os.Stdin.Read(response)
		done <- n
	}()

	select {
	case n := <-done:
		if n > 0 {
			// 响应格式通常是: \x1b]11;rgb:RRRR/GGGG/BBBB\x1b\\
			responseStr := string(response[:n])
			return parseOSC11Response(responseStr)
		}
	case <-time.After(100 * time.Millisecond):
		// 超时，终端可能不支持
		return ""
	}

	return ""
}

// parseOSC11Response 解析 OSC 11 响应
func parseOSC11Response(response string) string {
	// 查找 "rgb:" 或 "rgba:"
	if idx := strings.Index(response, "rgb:"); idx != -1 {
		// 提取颜色部分
		colorPart := response[idx+4:]
		if endIdx := strings.IndexAny(colorPart, "\x1b\x07\\"); endIdx != -1 {
			return colorPart[:endIdx]
		}
	}
	return ""
}

// isColorDark 判断 RGB 颜色是否为深色
// 使用相对亮度公式: Y = 0.2126*R + 0.7152*G + 0.0722*B
func isColorDark(rgbStr string) bool {
	// rgbStr 格式: "RRRR/GGGG/BBBB" (16位十六进制)
	parts := strings.Split(rgbStr, "/")
	if len(parts) != 3 {
		return true // 默认深色
	}

	// 解析 RGB 值 (取前两位十六进制即可，范围 0-255)
	var r, g, b int
	fmt.Sscanf(parts[0][:2], "%x", &r)
	fmt.Sscanf(parts[1][:2], "%x", &g)
	fmt.Sscanf(parts[2][:2], "%x", &b)

	// 计算相对亮度
	luminance := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)

	// 如果亮度 < 128，认为是深色背景
	return luminance < 128
}

// InitTheme 初始化主题，根据终端背景色设置环境变量
// 这样 lipgloss.AdaptiveColor 就能正确工作
func InitTheme() {
	isDark := DetectTerminalBackground()

	if isDark {
		// 设置为深色主题
		os.Setenv("TERM_BACKGROUND", "dark")
		// lipgloss 会检查这些环境变量
		lipgloss.SetHasDarkBackground(true)
	} else {
		// 设置为浅色主题
		os.Setenv("TERM_BACKGROUND", "light")
		lipgloss.SetHasDarkBackground(false)
	}
}
