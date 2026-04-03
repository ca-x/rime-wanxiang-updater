package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/controller"
	"rime-wanxiang-updater/internal/i18n"
	"rime-wanxiang-updater/internal/termcolor"
	"rime-wanxiang-updater/internal/theme"
	"rime-wanxiang-updater/internal/ui"
	"rime-wanxiang-updater/internal/version"
)

var (
	// 赛博朋克色彩 - 自适应深色/浅色背景
	neonCyan    = lipgloss.AdaptiveColor{Light: "#008B8B", Dark: "#00FFFF"}
	neonGreen   = lipgloss.AdaptiveColor{Light: "#008000", Dark: "#00FF41"}
	glitchRed   = lipgloss.AdaptiveColor{Light: "#DC143C", Dark: "#FF0040"}
	darkBg      = lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#0A0E27"}
	mutedText   = lipgloss.AdaptiveColor{Light: "#64748B", Dark: "#6B7280"}
	panelBorder = lipgloss.AdaptiveColor{Light: "#94A3B8", Dark: "#334155"}

	bootLogo = `
╔═══════════════════════════════════════════════════════════════╗
║  ██████╗ ██╗███╗   ███╗███████╗    ██╗    ██╗ █████╗ ███╗   ██╗║
║  ██╔══██╗██║████╗ ████║██╔════╝    ██║    ██║██╔══██╗████╗  ██║║
║  ██████╔╝██║██╔████╔██║█████╗      ██║ █╗ ██║███████║██╔██╗ ██║║
║  ██╔══██╗██║██║╚██╔╝██║██╔══╝      ██║███╗██║██╔══██║██║╚██╗██║║
║  ██║  ██║██║██║ ╚═╝ ██║███████╗    ╚███╔███╔╝██║  ██║██║ ╚████║║
║  ╚═╝  ╚═╝╚═╝╚═╝     ╚═╝╚══════╝     ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═══╝║
╚═══════════════════════════════════════════════════════════════╝`
)

type bootScreenStyles struct {
	logo    lipgloss.Style
	version lipgloss.Style
	muted   lipgloss.Style
	pending lipgloss.Style
	done    lipgloss.Style
	active  lipgloss.Style
	hint    lipgloss.Style
	panel   lipgloss.Style
}

type exitScreenStyles struct {
	primary   lipgloss.Style
	secondary lipgloss.Style
	muted     lipgloss.Style
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func terminalSize() (int, int) {
	for _, fd := range []int{int(os.Stdout.Fd()), int(os.Stderr.Fd()), int(os.Stdin.Fd())} {
		width, height, err := term.GetSize(fd)
		if err == nil && width > 0 && height > 0 {
			return width, height
		}
	}

	return 80, 24
}

func placeBlock(block string, width, height int) string {
	if width <= 0 || height <= 0 {
		return block
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, block)
}

func placeHorizontally(block string, width int) string {
	if width <= 0 {
		return block
	}

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, block)
}

func renderBootScreen(
	locale i18n.Locale,
	width int,
	height int,
	steps []string,
	completed int,
	showLaunch bool,
	launchSuffix string,
	styles bootScreenStyles,
) string {
	stepPanel := renderBootStepPanel(steps, completed, showLaunch, width, styles)
	lines := []string{
		styles.logo.Render(bootLogo),
		"",
		styles.version.Render(i18n.Text(locale, "boot.version", version.GetVersion())),
		"",
		stepPanel,
	}

	if showLaunch {
		lines = append(lines, "")
		lines = append(lines, styles.hint.Render(i18n.Text(locale, "boot.launch")+launchSuffix))
	}

	return placeBlock(strings.Join(lines, "\n"), width, height)
}

func renderBootStepPanel(
	steps []string,
	completed int,
	showLaunch bool,
	screenWidth int,
	styles bootScreenStyles,
) string {
	lines := make([]string, 0, len(steps))
	for i, step := range steps {
		markerStyle := styles.pending
		marker := "·"
		textStyle := styles.muted

		switch {
		case showLaunch || i < completed-1:
			markerStyle = styles.done
			marker = "●"
			textStyle = styles.done
		case completed > 0 && i == completed-1:
			markerStyle = styles.active
			marker = "◉"
			textStyle = styles.active
		}

		lines = append(lines, markerStyle.Render(marker)+" "+textStyle.Render(step))
	}

	return styles.panel.
		Width(bootPanelWidth(steps, screenWidth)).
		Render(strings.Join(lines, "\n"))
}

func bootPanelWidth(steps []string, screenWidth int) int {
	width := 34
	for _, step := range steps {
		stepWidth := lipgloss.Width(step) + 7
		if stepWidth > width {
			width = stepWidth
		}
	}

	if width > 54 {
		width = 54
	}

	maxWidth := screenWidth - 18
	if maxWidth < 34 {
		maxWidth = 34
	}
	if width > maxWidth {
		width = maxWidth
	}

	return width
}

func renderExitScreen(locale i18n.Locale, width int, height int, styles exitScreenStyles) string {
	lines := []string{
		styles.muted.Render(i18n.Text(locale, "boot.version", version.GetVersion())),
		"",
		styles.primary.Render(i18n.Text(locale, "boot.exit.line1")),
		styles.secondary.Render(i18n.Text(locale, "boot.exit.line2")),
	}

	return placeBlock(strings.Join(lines, "\n"), width, height)
}

func printBootSequence(locale i18n.Locale) {
	width, height := terminalSize()

	styles := bootScreenStyles{
		logo: lipgloss.NewStyle().
			Foreground(neonCyan).
			Bold(true),
		version: lipgloss.NewStyle().
			Foreground(mutedText).
			Bold(true),
		muted: lipgloss.NewStyle().
			Foreground(mutedText),
		pending: lipgloss.NewStyle().
			Foreground(mutedText),
		done: lipgloss.NewStyle().
			Foreground(neonGreen),
		active: lipgloss.NewStyle().
			Foreground(neonCyan).
			Bold(true),
		hint: lipgloss.NewStyle().
			Foreground(mutedText).
			Italic(true),
		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(panelBorder).
			Padding(0, 2),
	}

	bootMessages := []string{
		i18n.Text(locale, "boot.step.init"),
		i18n.Text(locale, "boot.step.model"),
		i18n.Text(locale, "boot.step.connect"),
		i18n.Text(locale, "boot.step.hardware", runtime.GOOS),
		i18n.Text(locale, "boot.step.files"),
		i18n.Text(locale, "boot.step.channel"),
		i18n.Text(locale, "boot.step.ready"),
	}

	for i := range bootMessages {
		clearScreen()
		fmt.Print(renderBootScreen(
			locale,
			width,
			height,
			bootMessages,
			i+1,
			false,
			"",
			styles,
		))
		time.Sleep(85 * time.Millisecond)
	}

	for _, suffix := range []string{".", "..", "..."} {
		clearScreen()
		fmt.Print(renderBootScreen(
			locale,
			width,
			height,
			bootMessages,
			len(bootMessages),
			true,
			suffix,
			styles,
		))
		time.Sleep(70 * time.Millisecond)
	}

	time.Sleep(180 * time.Millisecond)
}

func main() {
	// 初始化终端颜色检测（自动检测终端背景色）
	termcolor.InitLipgloss()

	// 加载配置
	cfg, err := config.NewManager()
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(glitchRed).
			Background(darkBg).
			Bold(true).
			Padding(0, 1)

		fmt.Println(errorStyle.Render("⚠ FATAL ERROR: " + err.Error()))
		os.Exit(1)
	}

	bootLocale := i18n.Normalize(cfg.Config.Language)

	// 显示启动序列
	printBootSequence(bootLocale)

	// 初始化主题管理器
	themeMgr := theme.NewManager()

	// 从配置加载主题设置
	if cfg.Config.ThemeAdaptive {
		light := cfg.Config.ThemeLight
		dark := cfg.Config.ThemeDark
		if light == "" {
			light = "cyberpunk-light"
		}
		if dark == "" {
			dark = "cyberpunk"
		}
		themeMgr.SetAdaptiveTheme(light, dark)
	} else if cfg.Config.ThemeFixed != "" {
		themeMgr.SetTheme(cfg.Config.ThemeFixed)
	}

	// 创建通信通道
	commandChan := make(chan controller.Command, 10)
	eventChan := make(chan controller.Event, 100)

	// 创建控制器
	ctrl := controller.NewController(cfg, commandChan, eventChan)

	// 在后台启动控制器
	go ctrl.Run()

	// 创建 UI 模型
	model := ui.NewModel(cfg, themeMgr, commandChan, eventChan)

	// 创建 Bubble Tea 程序
	// 使用 WithAltScreen 避免启动序列输出影响终端焦点
	p := tea.NewProgram(model, tea.WithAltScreen())

	// 运行程序
	if _, err := p.Run(); err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(glitchRed).
			Background(darkBg).
			Bold(true).
			Padding(0, 1)

		fmt.Println(errorStyle.Render("⚠ RUNTIME ERROR: " + err.Error()))
		os.Exit(1)
	}

	// 停止控制器
	ctrl.Stop()

	width, _ := terminalSize()
	exitBlock := renderExitScreen(bootLocale, width, 8, exitScreenStyles{
		primary: lipgloss.NewStyle().
			Foreground(neonCyan).
			Bold(true),
		secondary: lipgloss.NewStyle().
			Foreground(mutedText),
		muted: lipgloss.NewStyle().
			Foreground(mutedText),
	})

	clearScreen()
	fmt.Print(exitBlock)
}
