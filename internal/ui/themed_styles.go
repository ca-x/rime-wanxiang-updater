package ui

import (
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

	// 映射颜色
	s.Primary = t.Cyan
	s.Secondary = t.Magenta
	s.Accent = t.Blue
	s.Success = t.Green
	s.Warning = t.Yellow
	s.Error = t.Red
	s.Muted = t.Comment
	s.Background = t.Background
	s.Foreground = t.Foreground

	// Logo 样式
	s.Logo = lipgloss.NewStyle().
		Foreground(s.Primary).
		Bold(true)

	// 头部样式
	s.Header = lipgloss.NewStyle().
		Foreground(s.Primary).
		Bold(true).
		Padding(1, 3).
		Border(lipgloss.ThickBorder()).
		BorderForeground(s.Secondary).
		Align(lipgloss.Center)

	// 菜单项样式
	s.MenuItem = lipgloss.NewStyle().
		Foreground(s.Primary).
		Padding(0, 2).
		MarginLeft(2)

	// 选中菜单项样式
	s.SelectedMenuItem = lipgloss.NewStyle().
		Foreground(s.Background).
		Background(s.Secondary).
		Padding(0, 2).
		Bold(true).
		MarginLeft(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(s.Primary).
		BorderLeft(true).
		BorderRight(true)

	// 信息框样式
	s.InfoBox = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(s.Primary).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	// 错误样式
	s.ErrorText = lipgloss.NewStyle().
		Foreground(s.Error).
		Bold(true).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(s.Error)

	// 成功样式
	s.SuccessText = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(s.Success)

	// 霓虹绿样式
	s.NeonGreen = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true)

	// 警告样式
	s.WarningText = lipgloss.NewStyle().
		Foreground(s.Warning).
		Bold(true).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(t.Orange)

	// 提示样式
	s.Hint = lipgloss.NewStyle().
		Foreground(t.Magenta).
		Italic(true).
		Padding(1, 0).
		Faint(true)

	// 配置项样式
	s.ConfigKey = lipgloss.NewStyle().
		Foreground(s.Secondary).
		Bold(true).
		Width(20)

	s.ConfigValue = lipgloss.NewStyle().
		Foreground(s.Primary)

	// 进度消息样式
	s.ProgressMsg = lipgloss.NewStyle().
		Foreground(s.Success).
		Bold(true).
		Padding(1, 2)

	// 容器样式
	s.Container = lipgloss.NewStyle().
		Padding(2, 3)

	// 扫描线样式
	s.ScanLine = lipgloss.NewStyle().
		Foreground(s.Primary).
		Faint(true)

	// 网格样式
	s.Grid = lipgloss.NewStyle().
		Foreground(s.Muted).
		Faint(true)

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
		Foreground(t.Magenta).
		Italic(true)

	// 状态栏样式
	s.StatusBar = lipgloss.NewStyle().
		Foreground(t.ForegroundAlt).
		Background(t.BackgroundAlt).
		Padding(0, 1)

	s.StatusKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(s.Secondary).
		Padding(0, 1).
		Bold(true)

	s.StatusValue = lipgloss.NewStyle().
		Foreground(t.ForegroundAlt).
		Background(t.BackgroundAlt).
		Padding(0, 1)

	// 对话框样式
	s.DialogBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Secondary).
		Padding(1, 2).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	s.DialogButton = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(t.BackgroundAlt).
		Padding(0, 3).
		MarginTop(1).
		MarginRight(2)

	s.DialogActiveButton = lipgloss.NewStyle().
		Foreground(s.Background).
		Background(s.Secondary).
		Padding(0, 3).
		MarginTop(1).
		MarginRight(2).
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
