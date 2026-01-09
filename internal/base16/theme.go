package base16

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme 基于 base16 配色方案的主题，包含预定义样式
type Theme struct {
	Scheme *Scheme

	// 语义化颜色
	Background    lipgloss.Color
	BackgroundAlt lipgloss.Color
	Selection     lipgloss.Color
	Comment       lipgloss.Color
	Foreground    lipgloss.Color
	ForegroundAlt lipgloss.Color

	Red     lipgloss.Color
	Orange  lipgloss.Color
	Yellow  lipgloss.Color
	Green   lipgloss.Color
	Cyan    lipgloss.Color
	Blue    lipgloss.Color
	Magenta lipgloss.Color
	Brown   lipgloss.Color

	// 预定义样式
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Text      lipgloss.Style
	Muted     lipgloss.Style
	Keyword   lipgloss.Style
	String    lipgloss.Style
	Number    lipgloss.Style
	Function  lipgloss.Style
	Variable  lipgloss.Style
	Comment_  lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Info      lipgloss.Style
	Highlight lipgloss.Style
	Selected  lipgloss.Style
	Border    lipgloss.Style
	StatusBar lipgloss.Style
	StatusKey lipgloss.Style
}

// NewTheme 从配色方案创建主题
func NewTheme(scheme *Scheme) *Theme {
	t := &Theme{Scheme: scheme}

	// 映射语义化颜色
	t.Background = scheme.Color(0x00)
	t.BackgroundAlt = scheme.Color(0x01)
	t.Selection = scheme.Color(0x02)
	t.Comment = scheme.Color(0x03)
	t.Foreground = scheme.Color(0x05)
	t.ForegroundAlt = scheme.Color(0x04)

	t.Red = scheme.Color(0x08)
	t.Orange = scheme.Color(0x09)
	t.Yellow = scheme.Color(0x0A)
	t.Green = scheme.Color(0x0B)
	t.Cyan = scheme.Color(0x0C)
	t.Blue = scheme.Color(0x0D)
	t.Magenta = scheme.Color(0x0E)
	t.Brown = scheme.Color(0x0F)

	// 创建预定义样式
	t.Title = lipgloss.NewStyle().
		Foreground(t.Blue).
		Bold(true)

	t.Subtitle = lipgloss.NewStyle().
		Foreground(t.Cyan).
		Italic(true)

	t.Text = lipgloss.NewStyle().
		Foreground(t.Foreground)

	t.Muted = lipgloss.NewStyle().
		Foreground(t.Comment).
		Faint(true)

	t.Keyword = lipgloss.NewStyle().
		Foreground(t.Magenta).
		Bold(true)

	t.String = lipgloss.NewStyle().
		Foreground(t.Green)

	t.Number = lipgloss.NewStyle().
		Foreground(t.Orange)

	t.Function = lipgloss.NewStyle().
		Foreground(t.Blue)

	t.Variable = lipgloss.NewStyle().
		Foreground(t.Red)

	t.Comment_ = lipgloss.NewStyle().
		Foreground(t.Comment).
		Italic(true)

	t.Error = lipgloss.NewStyle().
		Foreground(t.Red).
		Bold(true)

	t.Success = lipgloss.NewStyle().
		Foreground(t.Green).
		Bold(true)

	t.Warning = lipgloss.NewStyle().
		Foreground(t.Yellow).
		Bold(true)

	t.Info = lipgloss.NewStyle().
		Foreground(t.Cyan).
		Bold(true)

	t.Highlight = lipgloss.NewStyle().
		Foreground(t.Yellow).
		Background(t.Selection)

	t.Selected = lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Magenta).
		Bold(true)

	t.Border = lipgloss.NewStyle().
		BorderForeground(t.Cyan)

	t.StatusBar = lipgloss.NewStyle().
		Foreground(t.ForegroundAlt).
		Background(t.BackgroundAlt)

	t.StatusKey = lipgloss.NewStyle().
		Foreground(scheme.Color(0x07)).
		Background(t.Magenta).
		Bold(true)

	return t
}

// AdaptiveTheme 支持明暗模式的自适应主题
type AdaptiveTheme struct {
	LightScheme *Scheme
	DarkScheme  *Scheme

	// 自适应颜色
	Background    lipgloss.AdaptiveColor
	BackgroundAlt lipgloss.AdaptiveColor
	Selection     lipgloss.AdaptiveColor
	Comment       lipgloss.AdaptiveColor
	Foreground    lipgloss.AdaptiveColor
	ForegroundAlt lipgloss.AdaptiveColor

	Red     lipgloss.AdaptiveColor
	Orange  lipgloss.AdaptiveColor
	Yellow  lipgloss.AdaptiveColor
	Green   lipgloss.AdaptiveColor
	Cyan    lipgloss.AdaptiveColor
	Blue    lipgloss.AdaptiveColor
	Magenta lipgloss.AdaptiveColor
	Brown   lipgloss.AdaptiveColor

	// 预定义样式
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Text      lipgloss.Style
	Muted     lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Info      lipgloss.Style
	Highlight lipgloss.Style
	Selected  lipgloss.Style
	Border    lipgloss.Style
	StatusBar lipgloss.Style
	StatusKey lipgloss.Style
}

// NewAdaptiveTheme 创建自适应主题
func NewAdaptiveTheme(light, dark *Scheme) *AdaptiveTheme {
	t := &AdaptiveTheme{
		LightScheme: light,
		DarkScheme:  dark,
	}

	// 创建自适应颜色
	t.Background = AdaptiveColor(light, dark, 0x00)
	t.BackgroundAlt = AdaptiveColor(light, dark, 0x01)
	t.Selection = AdaptiveColor(light, dark, 0x02)
	t.Comment = AdaptiveColor(light, dark, 0x03)
	t.Foreground = AdaptiveColor(light, dark, 0x05)
	t.ForegroundAlt = AdaptiveColor(light, dark, 0x04)

	t.Red = AdaptiveColor(light, dark, 0x08)
	t.Orange = AdaptiveColor(light, dark, 0x09)
	t.Yellow = AdaptiveColor(light, dark, 0x0A)
	t.Green = AdaptiveColor(light, dark, 0x0B)
	t.Cyan = AdaptiveColor(light, dark, 0x0C)
	t.Blue = AdaptiveColor(light, dark, 0x0D)
	t.Magenta = AdaptiveColor(light, dark, 0x0E)
	t.Brown = AdaptiveColor(light, dark, 0x0F)

	// 创建样式
	t.Title = lipgloss.NewStyle().
		Foreground(t.Blue).
		Bold(true)

	t.Subtitle = lipgloss.NewStyle().
		Foreground(t.Cyan).
		Italic(true)

	t.Text = lipgloss.NewStyle().
		Foreground(t.Foreground)

	t.Muted = lipgloss.NewStyle().
		Foreground(t.Comment).
		Faint(true)

	t.Error = lipgloss.NewStyle().
		Foreground(t.Red).
		Bold(true)

	t.Success = lipgloss.NewStyle().
		Foreground(t.Green).
		Bold(true)

	t.Warning = lipgloss.NewStyle().
		Foreground(t.Yellow).
		Bold(true)

	t.Info = lipgloss.NewStyle().
		Foreground(t.Cyan).
		Bold(true)

	t.Highlight = lipgloss.NewStyle().
		Foreground(t.Yellow).
		Background(t.Selection)

	t.Selected = lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Magenta).
		Bold(true)

	t.Border = lipgloss.NewStyle().
		BorderForeground(t.Cyan)

	t.StatusBar = lipgloss.NewStyle().
		Foreground(t.ForegroundAlt).
		Background(t.BackgroundAlt)

	t.StatusKey = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#" + light.Base07, Dark: "#" + dark.Base07}).
		Background(t.Magenta).
		Bold(true)

	return t
}
