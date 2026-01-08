package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/ui"
	"rime-wanxiang-updater/internal/version"
)

var (
	// 赛博朋克色彩 - 自适应深色/浅色背景
	neonCyan    = lipgloss.AdaptiveColor{Light: "#008B8B", Dark: "#00FFFF"}
	neonMagenta = lipgloss.AdaptiveColor{Light: "#8B008B", Dark: "#FF00FF"}
	neonGreen   = lipgloss.AdaptiveColor{Light: "#008000", Dark: "#00FF41"}
	neonPurple  = lipgloss.AdaptiveColor{Light: "#6A0DAD", Dark: "#B026FF"}
	glitchRed   = lipgloss.AdaptiveColor{Light: "#DC143C", Dark: "#FF0040"}
	darkBg      = lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#0A0E27"}

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

func printBootSequence() {
	// 清屏
	fmt.Print("\033[H\033[2J")

	// Logo 样式
	logoStyle := lipgloss.NewStyle().
		Foreground(neonCyan).
		Bold(true)

	fmt.Println(logoStyle.Render(bootLogo))
	fmt.Println() // Logo 后添加空行

	// 版本信息
	versionStyle := lipgloss.NewStyle().
		Foreground(neonMagenta).
		Bold(true)

	versionText := fmt.Sprintf("              >>> UPDATER SYSTEM %s <<<", version.GetVersion())
	fmt.Println(versionStyle.Render(versionText))

	// 添加装饰性分隔线
	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#C0C0C0", Dark: "#2A2F4A"})
	divider := "▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔"
	fmt.Println(dividerStyle.Render(divider))
	fmt.Println()

	// 启动序列
	bootStyle := lipgloss.NewStyle().
		Foreground(neonGreen)

	bootMessages := []string{
		"[✓] INITIALIZING SYSTEM...",
		"[✓] LOADING NEURAL NETWORK...",
		"[✓] CONNECTING TO MATRIX...",
		"[✓] SCANNING HARDWARE: " + runtime.GOOS,
		"[✓] MOUNTING FILE SYSTEMS...",
		"[✓] ESTABLISHING SECURE CHANNELS...",
		"[✓] SYSTEM READY",
	}

	for _, msg := range bootMessages {
		fmt.Println(bootStyle.Render("  " + msg))
		time.Sleep(150 * time.Millisecond)
	}

	// 添加装饰性分隔线
	fmt.Println()
	fmt.Println(dividerStyle.Render(divider))

	fmt.Println()

	// 启动提示
	hintStyle := lipgloss.NewStyle().
		Foreground(neonPurple).
		Italic(true)

	fmt.Println(hintStyle.Render("  ⚡ LAUNCHING MAIN INTERFACE..."))
	fmt.Println()

	time.Sleep(500 * time.Millisecond)
}

func main() {
	// 初始化主题（自动检测终端背景色）
	ui.InitTheme()

	// 显示启动序列
	printBootSequence()

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

	// 创建 Bubble Tea 程序
	p := tea.NewProgram(ui.NewModel(cfg))

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

	// 退出消息
	exitStyle := lipgloss.NewStyle().
		Foreground(neonCyan).
		Bold(true)

	fmt.Println()
	fmt.Println(exitStyle.Render("╔════════════════════════════════════════╗"))
	fmt.Println(exitStyle.Render("║    SYSTEM SHUTDOWN COMPLETE            ║"))
	fmt.Println(exitStyle.Render("║    SEE YOU IN THE NEXT SESSION         ║"))
	fmt.Println(exitStyle.Render("╚════════════════════════════════════════╝"))
	fmt.Println()
}
