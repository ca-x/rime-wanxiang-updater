package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderStatusBarThemedCompactsToSingleLine(t *testing.T) {
	styles := &Styles{
		StatusBar:   lipgloss.NewStyle(),
		StatusKey:   lipgloss.NewStyle(),
		StatusValue: lipgloss.NewStyle(),
	}

	rendered := RenderStatusBarThemed(
		styles,
		64,
		"版本:",
		"引擎:",
		"下载源:",
		"当前方案:",
		"dev",
		"鼠须管",
		"CNB 镜像",
		"基础版",
	)

	if strings.Contains(rendered, "\n") {
		t.Fatalf("RenderStatusBarThemed() contains newline, want single line: %q", rendered)
	}
	if got := lipgloss.Width(rendered); got > 64 {
		t.Fatalf("RenderStatusBarThemed() width = %d, want <= 64", got)
	}
	for _, unwanted := range []string{"版本:", "引擎:", "下载源:", "当前方案:"} {
		if strings.Contains(rendered, unwanted) {
			t.Fatalf("RenderStatusBarThemed() contains %q in compact mode: %q", unwanted, rendered)
		}
	}

	for _, want := range []string{"dev", "鼠须管", "CNB 镜像", "基础版"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("RenderStatusBarThemed() missing %q: %q", want, rendered)
		}
	}
}
