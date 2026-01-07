package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"runtime"
	"strings"
)

var (
	// 赛博朋克色彩定义 - 优化为适配深色和浅色背景
	neonCyan      = lipgloss.AdaptiveColor{Light: "#008B8B", Dark: "#00FFFF"} // 霓虹青色
	neonMagenta   = lipgloss.AdaptiveColor{Light: "#8B008B", Dark: "#FF00FF"} // 霓虹品红
	neonGreen     = lipgloss.AdaptiveColor{Light: "#008000", Dark: "#00FF41"} // 霓虹绿（矩阵代码）
	neonPink      = lipgloss.AdaptiveColor{Light: "#C71585", Dark: "#FF10F0"} // 霓虹粉
	neonOrange    = lipgloss.AdaptiveColor{Light: "#FF4500", Dark: "#FF6600"} // 霓虹橙
	neonBlue      = lipgloss.AdaptiveColor{Light: "#0000CD", Dark: "#0080FF"} // 霓虹蓝
	neonPurple    = lipgloss.AdaptiveColor{Light: "#6A0DAD", Dark: "#B026FF"} // 霓虹紫
	neonYellow    = lipgloss.AdaptiveColor{Light: "#DAA520", Dark: "#FFFF00"} // 霓虹黄
	darkBg        = lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#0A0E27"} // 深色背景
	darkBg2       = lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#1A1F3A"} // 次级深色背景
	glitchRed     = lipgloss.AdaptiveColor{Light: "#DC143C", Dark: "#FF0040"} // 故障红
	terminalGreen = lipgloss.AdaptiveColor{Light: "#008000", Dark: "#00FF00"} // 终端绿
	shadowGray    = lipgloss.AdaptiveColor{Light: "#A9A9A9", Dark: "#1C1C28"} // 阴影灰
	gridColor     = lipgloss.AdaptiveColor{Light: "#C0C0C0", Dark: "#2A2F4A"} // 网格线颜色

	// ASCII 艺术标题
	asciiLogo = `
╔═══════════════════════════════════════════════════════════════╗
║  ██████╗ ██╗███╗   ███╗███████╗    ██╗    ██╗ █████╗ ███╗   ██╗║
║  ██╔══██╗██║████╗ ████║██╔════╝    ██║    ██║██╔══██╗████╗  ██║║
║  ██████╔╝██║██╔████╔██║█████╗      ██║ █╗ ██║███████║██╔██╗ ██║║
║  ██╔══██╗██║██║╚██╔╝██║██╔══╝      ██║███╗██║██╔══██║██║╚██╗██║║
║  ██║  ██║██║██║ ╚═╝ ██║███████╗    ╚███╔███╔╝██║  ██║██║ ╚████║║
║  ╚═╝  ╚═╝╚═╝╚═╝     ╚═╝╚══════╝     ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═══╝║
╚═══════════════════════════════════════════════════════════════╝`

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
			Foreground(neonPink).
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
			Foreground(neonPurple).
			Italic(true)

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

	// 渐变色调色板（从粉色到黄色）
	titleGradient = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 30)

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
	var result strings.Builder
	chars := []rune(text)

	for i, ch := range chars {
		colorIdx := (i * len(titleGradient)) / len(chars)
		if colorIdx >= len(titleGradient) {
			colorIdx = len(titleGradient) - 1
		}

		color, _ := colorful.MakeColor(titleGradient[colorIdx])
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color.Hex())).
			Bold(true)

		result.WriteString(style.Render(string(ch)))
	}

	// 居中显示，宽度65
	centeredResult := lipgloss.PlaceHorizontal(65, lipgloss.Center, result.String())

	return "\n" + centeredResult + "\n"
}

// RenderStatusBar 渲染底部状态栏
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

// RenderBootSequence 渲染启动序列状态
func RenderBootSequence(appVersion string) string {
	var b strings.Builder

	// 标题行
	title := versionStyle.Render(">>> UPDATER SYSTEM " + appVersion + " <<<")
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
	launchMsg := "  " + neonGreenStyle.Render("⚡") + " " + configValueStyle.Render("LAUNCHING MAIN INTERFACE...")
	b.WriteString(launchMsg + "\n")

	return b.String()
}
