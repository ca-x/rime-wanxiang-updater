package ui

import "github.com/charmbracelet/lipgloss"

var (
	// 颜色定义
	primaryColor   = lipgloss.Color("#00BFFF") // 蓝色
	secondaryColor = lipgloss.Color("#FF6B6B") // 红色
	successColor   = lipgloss.Color("#51CF66") // 绿色
	warningColor   = lipgloss.Color("#FFD93D") // 黄色
	mutedColor     = lipgloss.Color("#6C7A89") // 灰色
	bgColor        = lipgloss.Color("#1A1B26") // 深色背景
	fgColor        = lipgloss.Color("#C0CAF5") // 前景色

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// 头部样式
	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Align(lipgloss.Center)

	// 菜单项样式 - 未选中
	menuItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 2)

	// 菜单项样式 - 选中
	selectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(primaryColor).
				Padding(0, 2).
				Bold(true)

	// 信息框样式
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// 错误样式
	errorStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Padding(0, 1)

	// 成功样式
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)

	// 警告样式
	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Padding(0, 1)

	// 提示样式
	hintStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(1, 0)

	// 配置项样式
	configKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Width(15)

	configValueStyle = lipgloss.NewStyle().
				Foreground(fgColor)

	// 进度消息样式
	progressMsgStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Padding(1, 2)

	// 容器样式
	containerStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// 分隔线样式
	dividerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(0, 2)
)
