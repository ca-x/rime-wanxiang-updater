package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"rime-wanxiang-updater/internal/base16"
	"rime-wanxiang-updater/internal/theme"

	"github.com/charmbracelet/lipgloss"
)

// Styles 主题化样式集合
type Styles struct {
	// 颜色
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Surface    lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Muted      lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color

	// 样式
	Logo               lipgloss.Style
	Header             lipgloss.Style
	MenuItem           lipgloss.Style
	SelectedMenuItem   lipgloss.Style
	InfoBox            lipgloss.Style
	ErrorText          lipgloss.Style
	SuccessText        lipgloss.Style
	WarningText        lipgloss.Style
	Hint               lipgloss.Style
	ConfigKey          lipgloss.Style
	ConfigValue        lipgloss.Style
	ProgressMsg        lipgloss.Style
	Container          lipgloss.Style
	ScanLine           lipgloss.Style
	Grid               lipgloss.Style
	Blink              lipgloss.Style
	StatusOnline       lipgloss.Style
	StatusProcessing   lipgloss.Style
	StatusError        lipgloss.Style
	Version            lipgloss.Style
	StatusBar          lipgloss.Style
	StatusKey          lipgloss.Style
	StatusValue        lipgloss.Style
	DialogBox          lipgloss.Style
	DialogButton       lipgloss.Style
	DialogActiveButton lipgloss.Style
	DialogCheckbox     lipgloss.Style
	NeonGreen          lipgloss.Style
}

// NewStyles 从主题创建样式
func NewStyles(t *base16.Theme) *Styles {
	s := &Styles{}

	background := softenBackground(t)
	foreground := pickComfortableForeground(t, background)
	muted := pickReadableSecondary(t)
	surface := pickReadableSurface(t)
	if isLightTheme(t) {
		muted = pickComfortableMuted(t, background)
		surface = softenSurface(surface, background)
	}

	// 映射颜色
	s.Primary = t.Cyan
	s.Secondary = t.Magenta
	s.Accent = t.Blue
	s.Surface = surface
	s.Success = t.Green
	s.Warning = t.Yellow
	s.Error = t.Red
	s.Muted = muted
	s.Background = background
	s.Foreground = foreground

	// Logo 样式
	s.Logo = lipgloss.NewStyle().
		Foreground(s.Primary).
		Bold(true)

	// 头部样式
	s.Header = lipgloss.NewStyle().
		Foreground(s.Primary).
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Secondary).
		Align(lipgloss.Center)

	// 菜单项样式
	s.MenuItem = lipgloss.NewStyle().
		Foreground(s.Primary).
		Padding(0, 1)

	// 选中菜单项样式
	s.SelectedMenuItem = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.Surface).
		Padding(0, 1).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.Primary)

	// 信息框样式
	s.InfoBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Secondary).
		Background(s.Surface).
		Padding(1, 2).
		MarginTop(0).
		MarginBottom(0)

	// 错误样式
	s.ErrorText = lipgloss.NewStyle().
		Foreground(s.Error).
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.Error)

	// 成功样式
	s.SuccessText = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.Success)

	// 霓虹绿样式
	s.NeonGreen = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true)

	// 警告样式
	s.WarningText = lipgloss.NewStyle().
		Foreground(s.Warning).
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Orange)

	// 提示样式
	s.Hint = lipgloss.NewStyle().
		Foreground(s.Muted).
		Italic(true).
		Padding(0, 0)

	// 配置项样式
	s.ConfigKey = lipgloss.NewStyle().
		Foreground(s.Secondary).
		Bold(true).
		Width(16)

	s.ConfigValue = lipgloss.NewStyle().
		Foreground(s.Foreground)

	// 进度消息样式
	s.ProgressMsg = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true).
		Padding(1, 2)

	// 容器样式
	s.Container = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.Background).
		Padding(2, 3)

	// 扫描线样式
	s.ScanLine = lipgloss.NewStyle().
		Foreground(s.Primary).
		Faint(true)

	// 网格样式
	s.Grid = lipgloss.NewStyle().
		Foreground(t.Comment)

	// 闪烁效果样式
	s.Blink = lipgloss.NewStyle().
		Foreground(t.Magenta).
		Bold(true).
		Blink(true)

	// 状态指示器样式
	s.StatusOnline = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true)

	s.StatusProcessing = lipgloss.NewStyle().
		Foreground(t.Orange).
		Bold(true)

	s.StatusError = lipgloss.NewStyle().
		Foreground(s.Error).
		Bold(true)

	// 版本号样式
	s.Version = lipgloss.NewStyle().
		Foreground(s.Muted)

	// 状态栏样式
	s.StatusBar = lipgloss.NewStyle().
		Foreground(s.Muted)

	s.StatusKey = lipgloss.NewStyle().
		Foreground(s.Muted).
		Bold(true)

	s.StatusValue = lipgloss.NewStyle().
		Foreground(s.Foreground)

	// 对话框样式
	s.DialogBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Secondary).
		Background(s.Surface).
		Padding(1, 2).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	s.DialogButton = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.Surface).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Comment).
		Padding(0, 2).
		MarginTop(1).
		MarginRight(1)

	s.DialogActiveButton = lipgloss.NewStyle().
		Foreground(s.Background).
		Background(s.Secondary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Secondary).
		Padding(0, 2).
		MarginTop(1).
		MarginRight(1).
		Bold(true)

	s.DialogCheckbox = lipgloss.NewStyle().
		Foreground(s.Primary).
		MarginTop(1)

	return s
}

// DefaultStyles 返回默认样式（使用当前主题）
func DefaultStyles(mgr *theme.Manager) *Styles {
	t := mgr.Current()
	if t == nil {
		// 回退到 cyberpunk
		t = base16.NewTheme(base16.Cyberpunk())
	}
	return NewStyles(t)
}

// RefreshStyles 刷新样式（主题切换后调用）
func (m *Model) RefreshStyles() {
	m.Styles = DefaultStyles(m.ThemeManager)
}

func pickReadableSecondary(t *base16.Theme) lipgloss.Color {
	if contrastRatio(t.ForegroundAlt, t.Background) >= 3.0 {
		return t.ForegroundAlt
	}

	return t.Foreground
}

func pickReadableSurface(t *base16.Theme) lipgloss.Color {
	if contrastRatio(t.BackgroundAlt, t.Background) >= 1.08 {
		return t.BackgroundAlt
	}

	return t.Selection
}

func isLightTheme(t *base16.Theme) bool {
	return relativeLuminance(t.Background) >= 0.55
}

func softenBackground(t *base16.Theme) lipgloss.Color {
	if !isLightTheme(t) {
		return t.Background
	}

	return blendColors(t.Background, t.BackgroundAlt, 0.12)
}

func softenSurface(surface, background lipgloss.Color) lipgloss.Color {
	return blendColors(surface, background, 0.38)
}

func pickComfortableForeground(t *base16.Theme, background lipgloss.Color) lipgloss.Color {
	if !isLightTheme(t) {
		return t.Foreground
	}

	candidate := blendColors(t.Foreground, t.Comment, 0.22)
	if contrastRatio(candidate, background) >= 6.5 {
		return candidate
	}

	return t.Foreground
}

func pickComfortableMuted(t *base16.Theme, background lipgloss.Color) lipgloss.Color {
	candidate := blendColors(t.Foreground, t.Comment, 0.55)
	if contrastRatio(candidate, background) >= 3.2 {
		return candidate
	}

	return pickReadableSecondary(t)
}

func contrastRatio(a, b lipgloss.Color) float64 {
	aLum := relativeLuminance(a)
	bLum := relativeLuminance(b)
	if aLum < bLum {
		aLum, bLum = bLum, aLum
	}

	return (aLum + 0.05) / (bLum + 0.05)
}

func relativeLuminance(c lipgloss.Color) float64 {
	r, g, b := parseHexColor(string(c))
	return 0.2126*channelLuminance(r) + 0.7152*channelLuminance(g) + 0.0722*channelLuminance(b)
}

func channelLuminance(v uint8) float64 {
	srgb := float64(v) / 255.0
	if srgb <= 0.03928 {
		return srgb / 12.92
	}

	return math.Pow((srgb+0.055)/1.055, 2.4)
}

func parseHexColor(raw string) (uint8, uint8, uint8) {
	hex := strings.TrimPrefix(strings.TrimSpace(raw), "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		panic(fmt.Errorf("parse red channel %q: %w", raw, err))
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		panic(fmt.Errorf("parse green channel %q: %w", raw, err))
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		panic(fmt.Errorf("parse blue channel %q: %w", raw, err))
	}

	return uint8(r), uint8(g), uint8(b)
}

func blendColors(a, b lipgloss.Color, ratio float64) lipgloss.Color {
	if ratio <= 0 {
		return a
	}
	if ratio >= 1 {
		return b
	}

	ar, ag, ab := parseHexColor(string(a))
	br, bg, bb := parseHexColor(string(b))

	mix := func(x, y uint8) uint8 {
		return uint8(math.Round(float64(x)*(1-ratio) + float64(y)*ratio))
	}

	return lipgloss.Color(
		fmt.Sprintf("#%02x%02x%02x", mix(ar, br), mix(ag, bg), mix(ab, bb)),
	)
}
