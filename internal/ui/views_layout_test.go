package ui

import (
	"strings"
	"testing"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"
	"rime-wanxiang-updater/internal/types"
	"rime-wanxiang-updater/internal/version"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderHeaderBlockUsesCompactHeaderLines(t *testing.T) {
	m := Model{
		Width: 80,
		Styles: &Styles{
			Primary:    lipgloss.Color(""),
			Accent:     lipgloss.Color(""),
			Border:     lipgloss.Color(""),
			Foreground: lipgloss.Color(""),
			Muted:      lipgloss.Color(""),
		},
	}

	rendered := strings.TrimSpace(m.renderHeaderBlock())
	lines := nonEmptyLines(rendered)
	if len(lines) != 2 {
		t.Fatalf("renderHeaderBlock() non-empty line count = %d, want 2: %q", len(lines), rendered)
	}

	for _, line := range lines {
		if got := lipgloss.Width(line); got > m.pageWidth() {
			t.Fatalf("renderHeaderBlock() line width = %d, want <= %d: %q", got, m.pageWidth(), line)
		}
	}
}

func TestRenderMenuEntrySelectedFitsPageWidth(t *testing.T) {
	m := Model{
		Width: 80,
		Styles: &Styles{
			Primary:    lipgloss.Color(""),
			Accent:     lipgloss.Color(""),
			Foreground: lipgloss.Color(""),
			Muted:      lipgloss.Color(""),
		},
	}

	rendered := m.renderMenuEntry(
		5,
		"⚙",
		"查看配置",
		"检查下载源、自动更新、代理和 Hook 等设置。",
		true,
	)

	if got := lipgloss.Width(rendered); got > m.pageWidth() {
		t.Fatalf("renderMenuEntry() width = %d, want <= %d: %q", got, m.pageWidth(), rendered)
	}
}

func TestRenderMenuDoesNotRepeatStatusBarSummary(t *testing.T) {
	themeMgr := theme.NewManager()
	if err := themeMgr.SetTheme("one-dark"); err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}

	cfg := &config.Manager{
		Config: &types.Config{
			SchemeType:       "base",
			UseMirror:        true,
			Language:         "zh-CN",
			InstalledEngines: []string{"fcitx5"},
		},
	}

	m := Model{
		Width:        80,
		Cfg:          cfg,
		ThemeManager: themeMgr,
		Styles:       &Styles{},
		RimeInstallStatus: detector.InstallationStatus{
			Installed: true,
		},
	}

	rendered := m.renderMenu()
	statusBar := RenderStatusBarThemed(
		m.Styles,
		m.pageWidth(),
		m.t("menu.summary.version"),
		m.t("menu.summary.engine"),
		m.t("menu.summary.source"),
		m.t("menu.summary.scheme"),
		version.GetVersion(),
		m.Cfg.GetEngineDisplayName(),
		m.configuredSourceLabel(),
		m.schemeLabel(m.Cfg.GetSchemeDisplayName()),
	)

	if strings.Contains(rendered, statusBar) {
		t.Fatalf("renderMenu() contains duplicated status bar summary: %q", statusBar)
	}
}

func TestRenderSummaryCardUsesCompactWidth(t *testing.T) {
	m := Model{
		Width: 80,
		Styles: &Styles{
			StatusKey:   lipgloss.NewStyle(),
			Foreground:  lipgloss.Color(""),
			Border:      lipgloss.Color(""),
			StatusValue: lipgloss.NewStyle(),
		},
	}

	rendered := strings.TrimSpace(m.renderSummaryCard([][2]string{
		{"当前方案:", "基础版"},
		{"引擎:", "fcitx5"},
		{"下载源:", "CNB 镜像"},
		{"自动更新:", "关闭"},
	}))

	if got := maxTrimmedLineWidth(rendered); got >= m.pageWidth() {
		t.Fatalf("renderSummaryCard() width = %d, want < %d", got, m.pageWidth())
	}
}

func TestRenderMenuUsesHintStripInsteadOfLegacyHintLine(t *testing.T) {
	themeMgr := theme.NewManager()
	if err := themeMgr.SetTheme("one-dark"); err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}

	cfg := &config.Manager{
		Config: &types.Config{
			SchemeType:       "base",
			UseMirror:        true,
			Language:         "zh-CN",
			InstalledEngines: []string{"fcitx5"},
		},
	}

	m := Model{
		Width:        80,
		Cfg:          cfg,
		ThemeManager: themeMgr,
		Styles:       DefaultStyles(themeMgr),
		RimeInstallStatus: detector.InstallationStatus{
			Installed: true,
		},
	}

	rendered := m.renderMenu()

	if strings.Contains(rendered, m.t("menu.hint")) {
		t.Fatalf("renderMenu() contains legacy hint line: %q", m.t("menu.hint"))
	}

	for _, want := range []string{
		m.t("ui.hint.shortcuts"),
		m.t("ui.hint.nav"),
		m.t("ui.hint.select"),
		m.t("ui.hint.quit"),
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("renderMenu() missing hint chip %q", want)
		}
	}
}

func TestRenderHintStripUsesSinglePanel(t *testing.T) {
	m := Model{
		Width: 80,
		Styles: &Styles{
			Muted:  lipgloss.Color(""),
			Border: lipgloss.Color(""),
		},
	}

	rendered := m.renderHintStrip("1-8 快捷操作", "↑↓ / J K", "Enter 选择", "A 关于", "Q 退出")

	if got := strings.Count(rendered, "╭"); got != 1 {
		t.Fatalf("renderHintStrip() top border count = %d, want 1: %q", got, rendered)
	}
	if got := strings.Count(rendered, "╰"); got != 1 {
		t.Fatalf("renderHintStrip() bottom border count = %d, want 1: %q", got, rendered)
	}
}

func nonEmptyLines(text string) []string {
	rawLines := strings.Split(text, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}

	return lines
}

func maxTrimmedLineWidth(text string) int {
	maxWidth := 0
	for line := range strings.SplitSeq(text, "\n") {
		width := lipgloss.Width(strings.TrimSpace(line))
		if width > maxWidth {
			maxWidth = width
		}
	}

	return maxWidth
}
