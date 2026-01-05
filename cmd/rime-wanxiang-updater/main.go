package main

import (
	"fmt"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
	"rime-wanxiang-updater/internal/ui"
)

func main() {
	fmt.Println("======================================")
	fmt.Println("  Rime 万象输入法更新工具 " + types.VERSION)
	fmt.Println("======================================")
	fmt.Printf("\n当前系统: %s\n\n", runtime.GOOS)

	// 加载配置
	cfg, err := config.NewManager()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 Bubble Tea 程序
	p := tea.NewProgram(ui.NewModel(cfg))

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Printf("程序运行失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n感谢使用！")
}
