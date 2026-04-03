package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"rime-wanxiang-updater/internal/i18n"
)

func TestRenderBootScreenShowsAllSteps(t *testing.T) {
	styles := bootScreenStyles{
		logo:    lipgloss.NewStyle(),
		version: lipgloss.NewStyle(),
		muted:   lipgloss.NewStyle(),
		pending: lipgloss.NewStyle(),
		done:    lipgloss.NewStyle(),
		active:  lipgloss.NewStyle(),
		hint:    lipgloss.NewStyle(),
		panel:   lipgloss.NewStyle(),
	}

	steps := []string{"初始化系统", "加载更新模块", "连接发布源"}
	rendered := renderBootScreen(i18n.LocaleZhCN, 80, 24, steps, 1, false, "", styles)

	for _, step := range steps {
		if !strings.Contains(rendered, step) {
			t.Fatalf("renderBootScreen() missing step %q", step)
		}
	}
}

func TestRenderExitScreenIncludesBothLines(t *testing.T) {
	styles := exitScreenStyles{
		primary:   lipgloss.NewStyle(),
		secondary: lipgloss.NewStyle(),
		muted:     lipgloss.NewStyle(),
	}

	rendered := renderExitScreen(i18n.LocaleZhCN, 80, 24, styles)

	for _, want := range []string{"本次会话已结束", "下次更新再见"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("renderExitScreen() missing %q", want)
		}
	}
}
