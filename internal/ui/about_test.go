package ui

import (
	"strings"
	"testing"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"
	"rime-wanxiang-updater/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleMenuInputOpensAboutViewWithA(t *testing.T) {
	m := Model{
		State: ViewMenu,
		Cfg: &config.Manager{
			Config: &types.Config{Language: "zh-CN"},
		},
	}

	next, _ := m.handleMenuInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	got := next.(Model)
	if got.State != ViewAbout {
		t.Fatalf("handleMenuInput(a) state = %v, want %v", got.State, ViewAbout)
	}
}

func TestRenderMenuContainsAboutHint(t *testing.T) {
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
	if !strings.Contains(rendered, m.t("ui.hint.about")) {
		t.Fatalf("renderMenu() missing about hint %q", m.t("ui.hint.about"))
	}
}

func TestRenderAboutUsesLocalizedCopy(t *testing.T) {
	themeMgr := theme.NewManager()
	if err := themeMgr.SetTheme("one-dark"); err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}

	m := Model{
		Width: 80,
		Cfg: &config.Manager{
			Config: &types.Config{Language: "en"},
		},
		ThemeManager: themeMgr,
		Styles:       DefaultStyles(themeMgr),
		State:        ViewAbout,
	}

	rendered := m.renderAbout()
	for _, want := range []string{
		m.t("about.title"),
		m.t("about.name.en"),
		m.t("about.name.zh"),
		m.t("about.homepage"),
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("renderAbout() missing %q", want)
		}
	}
}
