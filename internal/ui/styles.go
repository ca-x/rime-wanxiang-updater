package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"runtime"
	"strings"
)

var (
	// 赛博朋克色彩定义 - 优化为适配深色和浅色背景
	neonCyan      = lipgloss.AdaptiveColor{Light: "#008B8B", Dark: "#00FFFF"} // 霓虹青色
	neonMagenta   = lipgloss.AdaptiveColor{Light: "#8B008B", Dark: "#FF00FF"} // 霓虹品红
	neonGreen     = lipgloss.AdaptiveColor{Light: "#008000", Dark: "#00FF41"} // 霓虹绿（矩阵代码）
	neonOrange    = lipgloss.AdaptiveColor{Light: "#FF4500", Dark: "#FF6600"} // 霓虹橙
	neonPurple    = lipgloss.AdaptiveColor{Light: "#6A0DAD", Dark: "#B026FF"} // 霓虹紫
	neonYellow    = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFFF00"} // 霓虹黄
	glitchRed     = lipgloss.AdaptiveColor{Light: "#DC143C", Dark: "#FF0040"} // 故障红
	terminalGreen = lipgloss.AdaptiveColor{Light: "#008000", Dark: "#00FF00"} // 终端绿
	gridColor     = lipgloss.AdaptiveColor{Light: "#C0C0C0", Dark: "#2A2F4A"} // 网格线颜色

	scanLine = "▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔"
	gridLine = "░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░"

	// 头部样式 - 赛博朋克风格
	headerStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			Bold(true).
			Padding(1, 3).
			Border(lipgloss.ThickBorder()).
			BorderForeground(neonMagenta).
			Align(lipgloss.Center)

	// Logo 样式
	logoStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			Bold(true)

	// 菜单项样式 - 未选中（赛博朋克）
	menuItemStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			Padding(0, 2).
			MarginLeft(2)

	// 菜单项样式 - 选中（霓虹高亮）
	selectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#0A0E27"}).
				Background(neonMagenta).
				Padding(0, 2).
				Bold(true).
				MarginLeft(2).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(neonCyan).
				BorderLeft(true).
				BorderRight(true)

	// 信息框样式 - 带霓虹边框
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(neonCyan).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// 错误样式 - 故障效果
	errorStyle = lipgloss.NewStyle().
			Foreground(glitchRed).
			Bold(true).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(glitchRed)

	// 成功样式 - 终端绿
	successStyle = lipgloss.NewStyle().
			Foreground(terminalGreen).
			Bold(true).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(neonGreen)

	// 霓虹绿样式 - 用于高亮重要信息
	neonGreenStyle = lipgloss.NewStyle().
			Foreground(neonGreen).
			Bold(true)

	// 警告样式 - 霓虹黄
	warningStyle = lipgloss.NewStyle().
			Foreground(neonYellow).
			Bold(true).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(neonOrange)

	// 提示样式 - 赛博风格
	hintStyle = lipgloss.NewStyle().
			Foreground(neonPurple).
			Italic(true).
			Padding(1, 0).
			Faint(true)

	// 配置项样式
	configKeyStyle = lipgloss.NewStyle().
			Foreground(neonMagenta).
			Bold(true).
			Width(20)

	configValueStyle = lipgloss.NewStyle().
				Foreground(neonCyan)

	// 进度消息样式 - 动态效果
	progressMsgStyle = lipgloss.NewStyle().
				Foreground(neonGreen).
				Bold(true).
				Padding(1, 2)

	// 容器样式
	containerStyle = lipgloss.NewStyle().
			Padding(2, 3)

	// 分隔线样式 - 扫描线效果
	dividerStyle = lipgloss.NewStyle().
			Foreground(gridColor).
			Padding(0, 0)

	// 扫描线样式
	scanLineStyle = lipgloss.NewStyle().
			Foreground(neonCyan).
			Faint(true)

	// 网格样式
	gridStyle = lipgloss.NewStyle().
			Foreground(gridColor).
			Faint(true)

	// 闪烁效果样式
	blinkStyle = lipgloss.NewStyle().
			Foreground(neonMagenta).
			Bold(true).
			Blink(true)

	// 状态指示器样式
	statusOnlineStyle = lipgloss.NewStyle().
				Foreground(neonGreen).
				Bold(true)

	statusProcessingStyle = lipgloss.NewStyle().
				Foreground(neonOrange).
				Bold(true)

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(glitchRed).
				Bold(true)

	// 版本号样式
	versionStyle = lipgloss.NewStyle().
			Foreground(gridColor)

	// 状态栏样式
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			Padding(0, 1)

	statusKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(neonMagenta).
			Padding(0, 1).
			Bold(true)

	statusValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
				Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
				Padding(0, 1)

	// 对话框样式
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(neonMagenta).
			Padding(1, 2).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	dialogButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#FFFDF5"}).
				Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#888B7E"}).
				Padding(0, 3).
				MarginTop(1).
				MarginRight(2)

	dialogActiveButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFF7DB"}).
				Background(neonMagenta).
				Padding(0, 3).
				MarginTop(1).
				MarginRight(2).
				Bold(true)

	dialogCheckboxStyle = lipgloss.NewStyle().
				Foreground(neonCyan).
				MarginTop(1)

	dialogCheckboxCheckedStyle = lipgloss.NewStyle().
					Foreground(neonGreen).
					MarginTop(1).
					Bold(true)
)

// RenderGradientTitle 渲染渐变色标题（无边框）
func RenderGradientTitle(text string) string {
	title := lipgloss.NewStyle().
		Foreground(neonCyan).
		Bold(true).
		Render(text)

	rule := lipgloss.NewStyle().
		Foreground(gridColor).
		Render("────")

	line := lipgloss.JoinHorizontal(lipgloss.Center, rule, " ", title, " ", rule)
	return "\n" + lipgloss.PlaceHorizontal(65, lipgloss.Center, line) + "\n"
}

// RenderStatusBar 渲染底部状态栏（保留用于向后兼容）
func RenderStatusBar(version, engine, source string) string {
	const width = 65

	versionKey := statusKeyStyle.Render("版本")
	versionVal := statusValueStyle.Render(version)

	engineKey := statusKeyStyle.Render("引擎")
	engineVal := statusValueStyle.Render(engine)

	sourceKey := statusKeyStyle.Render("下载源")
	sourceVal := statusValueStyle.Render(source)

	// 拼接状态栏
	bar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		versionKey, versionVal, " ",
		engineKey, engineVal, " ",
		sourceKey, sourceVal,
	)

	// 填充到固定宽度
	barWidth := lipgloss.Width(bar)
	if barWidth < width {
		bar += strings.Repeat(" ", width-barWidth)
	}

	return statusBarStyle.Width(width).Render(bar)
}

// RenderStatusBarThemed 渲染底部状态栏（主题化版本）。
func RenderStatusBarThemed(
	s *Styles,
	width int,
	versionLabel string,
	engineLabel string,
	sourceLabel string,
	schemeLabel string,
	version string,
	engine string,
	source string,
	scheme string,
) string {
	bar := buildStatusBarLine(s, width, versionLabel, engineLabel, sourceLabel, schemeLabel, version, engine, source, scheme)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(bar)
}

func buildStatusBarLine(
	s *Styles,
	width int,
	versionLabel string,
	engineLabel string,
	sourceLabel string,
	schemeLabel string,
	version string,
	engine string,
	source string,
	scheme string,
) string {
	separator := s.StatusBar.Render("  ·  ")
	fullParts := []string{
		s.StatusKey.Render(versionLabel) + " " + s.StatusValue.Render(version),
		s.StatusKey.Render(engineLabel) + " " + s.StatusValue.Render(engine),
		s.StatusKey.Render(sourceLabel) + " " + s.StatusValue.Render(source),
		s.StatusKey.Render(schemeLabel) + " " + s.StatusValue.Render(scheme),
	}
	fullBar := strings.Join(fullParts, separator)
	if lipgloss.Width(fullBar) <= width {
		return fullBar
	}

	values := []string{version, engine, source, scheme}
	compactBar := renderStatusBarValues(s, values)
	if lipgloss.Width(compactBar) <= width {
		return compactBar
	}

	return renderStatusBarValues(s, compressStatusBarValues(values, width, lipgloss.Width(separator)))
}

func renderStatusBarValues(s *Styles, values []string) string {
	renderedValues := make([]string, 0, len(values))
	for _, value := range values {
		renderedValues = append(renderedValues, s.StatusValue.Render(value))
	}

	return strings.Join(renderedValues, s.StatusBar.Render("  ·  "))
}

func compressStatusBarValues(values []string, width int, separatorWidth int) []string {
	compressed := append([]string(nil), values...)
	if len(compressed) == 0 {
		return compressed
	}

	minWidths := make([]int, len(compressed))
	for i := range compressed {
		minWidths[i] = 4
		if lipgloss.Width(compressed[i]) < minWidths[i] {
			minWidths[i] = lipgloss.Width(compressed[i])
		}
	}

	for statusBarValuesWidth(compressed, separatorWidth) > width {
		index := widestStatusBarValue(compressed)
		if lipgloss.Width(compressed[index]) <= minWidths[index] {
			break
		}

		nextWidth := lipgloss.Width(compressed[index]) - 1
		if nextWidth < minWidths[index] {
			nextWidth = minWidths[index]
		}
		compressed[index] = truncateStatusBarValue(compressed[index], nextWidth)
	}

	return compressed
}

func statusBarValuesWidth(values []string, separatorWidth int) int {
	width := 0
	for i, value := range values {
		if i > 0 {
			width += separatorWidth
		}
		width += lipgloss.Width(value)
	}

	return width
}

func widestStatusBarValue(values []string) int {
	index := 0
	maxWidth := lipgloss.Width(values[0])
	for i := 1; i < len(values); i++ {
		valueWidth := lipgloss.Width(values[i])
		if valueWidth > maxWidth {
			index = i
			maxWidth = valueWidth
		}
	}

	return index
}

func truncateStatusBarValue(value string, width int) string {
	if width <= 0 || lipgloss.Width(value) <= width {
		return value
	}
	if width == 1 {
		return "…"
	}

	ellipsisWidth := lipgloss.Width("…")
	runes := []rune(value)
	currentWidth := 0
	var b strings.Builder
	for _, r := range runes {
		runeWidth := lipgloss.Width(string(r))
		if currentWidth+runeWidth+ellipsisWidth > width {
			break
		}
		b.WriteRune(r)
		currentWidth += runeWidth
	}

	return b.String() + "…"
}

// RenderBootSequence 渲染启动序列状态
func RenderBootSequence(appVersion string) string {
	var b strings.Builder

	// 标题行
	title := versionStyle.Render("Updater System " + appVersion)
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Render(title) + "\n\n")

	// 使用 list 渲染启动项
	checkmark := neonGreenStyle.Render("[✓]")
	bootList := list.New(
		"INITIALIZING SYSTEM...",
		"LOADING NEURAL NETWORK...",
		"CONNECTING TO MATRIX...",
		"SCANNING HARDWARE: "+runtime.GOOS,
		"MOUNTING FILE SYSTEMS...",
		"ESTABLISHING SECURE CHANNELS...",
		"SYSTEM READY",
	).
		Enumerator(func(items list.Items, index int) string {
			return checkmark
		}).
		EnumeratorStyleFunc(func(items list.Items, index int) lipgloss.Style {
			return lipgloss.NewStyle().PaddingRight(1).MarginLeft(2)
		}).
		ItemStyleFunc(func(items list.Items, index int) lipgloss.Style {
			return configValueStyle
		})

	// 直接使用 list 的输出
	b.WriteString(bootList.String() + "\n")

	// 添加空行和启动界面提示
	b.WriteString("\n")
	launchMsg := "  " + neonGreenStyle.Render("•") + " " + configValueStyle.Render("Launching main interface...")
	b.WriteString(launchMsg + "\n")

	return b.String()
}

// RenderHeader 渲染简洁的页面标题（用于主界面，不显示启动序列）
func RenderHeader(appVersion string) string {
	title := versionStyle.Render("Rime Wanxiang Updater " + appVersion)
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Render(title) + "\n"
}

// RenderCheckList 渲染带 checkmark 的列表
// title: 列表标题(英文)
// items: 列表项(中文组件名)
// isUpdated: true=已更新(绿色), false=已是最新(灰色)
// versions: 组件版本信息
func RenderCheckList(title string, items []string, isUpdated bool, versions map[string]string) string {
	var b strings.Builder

	// 标题样式
	var titleStyle lipgloss.Style
	var checkmark string
	var itemStyle lipgloss.Style

	if isUpdated {
		// 已更新：绿色
		titleStyle = lipgloss.NewStyle().Foreground(neonGreen).Bold(true)
		checkmark = "✓"
		itemStyle = configValueStyle
	} else {
		// 已是最新：灰色
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#808080", Dark: "#808080"}).Bold(true)
		checkmark = "✓"
		itemStyle = hintStyle
	}

	// 渲染标题
	b.WriteString(titleStyle.Render(title+":") + "\n")

	// 将中文组件名映射到英文
	componentMap := map[string]string{
		"方案": "Scheme",
		"词库": "Dictionary",
		"模型": "Model",
	}

	// 使用 list 渲染列表项
	listItems := make([]any, len(items))
	for i, item := range items {
		englishName, ok := componentMap[item]
		var itemText string
		if ok {
			// 中英文都显示
			itemText = englishName + " | " + item
		} else {
			// 只有中文时直接使用
			itemText = item
		}

		// 如果有版本信息，追加版本号
		if versions != nil {
			if version, ok := versions[item]; ok && version != "" {
				itemText += " (" + version + ")"
			}
		}

		listItems[i] = itemText
	}

	checkList := list.New(listItems...).
		Enumerator(func(items list.Items, index int) string {
			return titleStyle.Render(checkmark)
		}).
		EnumeratorStyleFunc(func(items list.Items, index int) lipgloss.Style {
			return lipgloss.NewStyle().PaddingRight(1).MarginLeft(1)
		}).
		ItemStyleFunc(func(items list.Items, index int) lipgloss.Style {
			return itemStyle
		})

	b.WriteString(checkList.String())

	return b.String()
}
