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
	// 赛博朋克色彩
	neonCyan    = lipgloss.Color("#00FFFF")
	neonMagenta = lipgloss.Color("#FF00FF")
	neonGreen   = lipgloss.Color("#00FF41")
	darkBg      = lipgloss.Color("#0A0E27")

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

	// 版本信息
	versionStyle := lipgloss.NewStyle().
		Foreground(neonMagenta).
		Bold(true)

	versionText := fmt.Sprintf("              >>> UPDATER SYSTEM %s <<<", version.GetVersion())
	fmt.Println(versionStyle.Render(versionText))
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

	fmt.Println()

	// 启动提示
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B026FF")).
		Italic(true)

	fmt.Println(hintStyle.Render("  ⚡ LAUNCHING MAIN INTERFACE..."))
	fmt.Println()

	time.Sleep(500 * time.Millisecond)
}

func main() {
	// 显示启动序列
	printBootSequence()

	// 加载配置
	cfg, err := config.NewManager()
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0040")).
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
			Foreground(lipgloss.Color("#FF0040")).
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
