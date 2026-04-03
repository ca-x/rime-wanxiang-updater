package ui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFcitxThemeSupportedForPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		engines  []string
		want     bool
	}{
		{
			name:     "linux with fcitx5",
			platform: "linux",
			engines:  []string{"fcitx5", "ibus"},
			want:     true,
		},
		{
			name:     "linux without fcitx5",
			platform: "linux",
			engines:  []string{"ibus"},
			want:     false,
		},
		{
			name:     "darwin with fcitx5 is not supported",
			platform: "darwin",
			engines:  []string{"fcitx5"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fcitxThemeSupportedForPlatform(tt.platform, tt.engines); got != tt.want {
				t.Fatalf("fcitxThemeSupportedForPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFcitxClassicUIConfigPreservesOtherLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "classicui.conf")

	initial := "[Groups/0]\nName=ClassicUI\nTheme=old-theme\nVertical Candidate List=False\n"
	if err := os.WriteFile(path, []byte(initial), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := writeFcitxClassicUIConfig(path, "new-theme"); err != nil {
		t.Fatalf("writeFcitxClassicUIConfig() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	content := string(data)
	for _, want := range []string{
		"Name=ClassicUI",
		"Theme=new-theme",
		"Vertical Candidate List=False",
	} {
		if !containsLine(content, want) {
			t.Fatalf("classicui.conf missing line %q in %q", want, content)
		}
	}
}

func TestInstallFcitxThemeCopiesThemeDirectory(t *testing.T) {
	source := fstest.MapFS{
		"demo/theme.conf":  &fstest.MapFile{Data: []byte("Theme=demo\n")},
		"demo/panel.svg":   &fstest.MapFile{Data: []byte("<svg></svg>")},
		"demo/arrow.png":   &fstest.MapFile{Data: []byte("png")},
		"other/theme.conf": &fstest.MapFile{Data: []byte("Theme=other\n")},
	}

	destRoot := t.TempDir()
	if err := installFcitxTheme(source, "demo", destRoot); err != nil {
		t.Fatalf("installFcitxTheme() error = %v", err)
	}

	for _, rel := range []string{
		"demo/theme.conf",
		"demo/panel.svg",
		"demo/arrow.png",
	} {
		path := filepath.Join(destRoot, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat(%q) error = %v", path, err)
		}
		if info.IsDir() {
			t.Fatalf("%q should be a file", path)
		}
	}

	if _, err := os.Stat(filepath.Join(destRoot, "other")); !os.IsNotExist(err) {
		t.Fatalf("unexpected unrelated theme copied")
	}
}

func TestSetFcitxThemePrefersDBusAndFallsBackToConfig(t *testing.T) {
	t.Run("dbus success skips file fallback", func(t *testing.T) {
		configPath := filepath.Join(t.TempDir(), "classicui.conf")
		called := false

		err := setFcitxThemeWithFallback("demo", configPath, func(themeName string) error {
			called = true
			if themeName != "demo" {
				t.Fatalf("themeName = %q, want %q", themeName, "demo")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("setFcitxThemeWithFallback() error = %v", err)
		}
		if !called {
			t.Fatalf("dbus setter was not called")
		}
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			t.Fatalf("config fallback should not run when dbus succeeds")
		}
	})

	t.Run("dbus failure falls back to config", func(t *testing.T) {
		configPath := filepath.Join(t.TempDir(), "classicui.conf")

		err := setFcitxThemeWithFallback("fallback-theme", configPath, func(string) error {
			return os.ErrPermission
		})
		if err != nil {
			t.Fatalf("setFcitxThemeWithFallback() error = %v", err)
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if !containsLine(string(data), "Theme=fallback-theme") {
			t.Fatalf("config fallback did not write theme: %q", string(data))
		}
	})
}

func TestInstalledFcitxThemeSelectionsLoadsExistingBuiltinThemes(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"demo", "other", "third-party"} {
		if err := os.MkdirAll(filepath.Join(root, name), 0755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", name, err)
		}
	}

	selections, err := installedFcitxThemeSelections(root, []string{"demo", "other", "missing"})
	if err != nil {
		t.Fatalf("installedFcitxThemeSelections() error = %v", err)
	}

	if !selections["demo"] || !selections["other"] {
		t.Fatalf("existing builtin themes should be selected: %#v", selections)
	}
	if selections["missing"] {
		t.Fatalf("missing builtin theme should not be selected: %#v", selections)
	}
	if selections["third-party"] {
		t.Fatalf("non-builtin theme should not be selected: %#v", selections)
	}
}

func TestSyncInstalledFcitxThemesCopiesSelectedRemovesUnselectedBuiltinAndKeepsCustom(t *testing.T) {
	source := fstest.MapFS{
		"demo/theme.conf":  &fstest.MapFile{Data: []byte("Theme=demo\n")},
		"other/theme.conf": &fstest.MapFile{Data: []byte("Theme=other\n")},
	}
	destRoot := t.TempDir()

	if err := os.MkdirAll(filepath.Join(destRoot, "demo"), 0755); err != nil {
		t.Fatalf("MkdirAll(demo) error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(destRoot, "third-party"), 0755); err != nil {
		t.Fatalf("MkdirAll(third-party) error = %v", err)
	}

	if err := syncInstalledFcitxThemes(source, destRoot, []string{"demo", "other"}, map[string]bool{
		"other": true,
	}); err != nil {
		t.Fatalf("syncInstalledFcitxThemes() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(destRoot, "demo")); !os.IsNotExist(err) {
		t.Fatalf("unselected builtin theme directory should be removed")
	}
	if _, err := os.Stat(filepath.Join(destRoot, "other", "theme.conf")); err != nil {
		t.Fatalf("selected builtin theme should be copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destRoot, "third-party")); err != nil {
		t.Fatalf("custom theme directory should be preserved: %v", err)
	}
}

func TestOpenFcitxThemeListPreselectsInstalledThemes(t *testing.T) {
	oldList := listAvailableFcitxThemes
	oldSelections := loadInstalledFcitxThemeSelections
	defer func() {
		listAvailableFcitxThemes = oldList
		loadInstalledFcitxThemeSelections = oldSelections
	}()

	listAvailableFcitxThemes = func() ([]string, error) {
		return []string{"demo", "other"}, nil
	}
	loadInstalledFcitxThemeSelections = func([]string) (map[string]bool, error) {
		return map[string]bool{"other": true}, nil
	}

	m := Model{
		Cfg: &config.Manager{
			Config: &types.Config{
				Language:         "zh-CN",
				InstalledEngines: []string{"fcitx5"},
			},
		},
	}

	next, _ := m.openFcitxThemeList()
	got := next.(Model)

	if got.State != ViewFcitxThemeList {
		t.Fatalf("openFcitxThemeList() state = %v, want %v", got.State, ViewFcitxThemeList)
	}
	if !got.FcitxThemeSelections["other"] {
		t.Fatalf("existing theme should be preselected: %#v", got.FcitxThemeSelections)
	}
	if got.FcitxThemeSelections["demo"] {
		t.Fatalf("non-existing theme should not be preselected: %#v", got.FcitxThemeSelections)
	}
}

func TestApplyFcitxThemeChoiceSyncsSelectionAndMovesToDefaultList(t *testing.T) {
	oldSync := syncInstalledFcitxThemeSelections
	defer func() { syncInstalledFcitxThemeSelections = oldSync }()

	called := false
	syncInstalledFcitxThemeSelections = func(themeNames []string, selections map[string]bool) error {
		called = true
		if !selections["other"] {
			t.Fatalf("selection should include existing checked theme: %#v", selections)
		}
		if !selections["demo"] {
			t.Fatalf("selection should include newly checked theme: %#v", selections)
		}
		return nil
	}

	m := Model{
		State:                ViewFcitxThemeList,
		FcitxThemeList:       []string{"demo", "other"},
		FcitxThemeChoice:     0,
		FcitxThemeSelections: map[string]bool{"other": true},
	}

	next, _ := m.handleFcitxThemeListInput(tea.KeyMsg{Type: tea.KeySpace})
	got := next.(Model)
	next, _ = got.handleFcitxThemeListInput(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)

	if !called {
		t.Fatalf("syncInstalledFcitxThemeSelections should be called")
	}
	if got.State != ViewFcitxThemeDefaultList {
		t.Fatalf("state after syncing fcitx themes = %v, want %v", got.State, ViewFcitxThemeDefaultList)
	}
}

func TestApplyFcitxThemeDefaultChoiceSetsSelectedTheme(t *testing.T) {
	oldSet := setFcitxThemeDefault
	defer func() { setFcitxThemeDefault = oldSet }()

	calledWith := ""
	setFcitxThemeDefault = func(themeName string) error {
		calledWith = themeName
		return nil
	}

	m := Model{
		State:                   ViewFcitxThemeDefaultList,
		FcitxThemeChoice:        1,
		FcitxThemeDefaultChoice: 1,
		FcitxThemeList:          []string{"demo", "other"},
		FcitxThemeSelections:    map[string]bool{"demo": true, "other": true},
	}

	next, _ := m.applyFcitxThemeDefaultChoice()
	got := next.(Model)

	if calledWith != "other" {
		t.Fatalf("setFcitxThemeDefault() called with %q, want %q", calledWith, "other")
	}
	if got.FcitxThemeSelected != "other" {
		t.Fatalf("FcitxThemeSelected = %q, want %q", got.FcitxThemeSelected, "other")
	}
	if got.State != ViewFcitxThemeDeployPrompt {
		t.Fatalf("state after applying default = %v, want %v", got.State, ViewFcitxThemeDeployPrompt)
	}
}

func TestApplyFcitxThemeChoiceReportsSyncError(t *testing.T) {
	oldSync := syncInstalledFcitxThemeSelections
	defer func() { syncInstalledFcitxThemeSelections = oldSync }()

	syncInstalledFcitxThemeSelections = func(themeNames []string, selections map[string]bool) error {
		return errors.New("boom")
	}

	m := Model{
		State:                ViewFcitxThemeList,
		FcitxThemeList:       []string{"demo"},
		FcitxThemeSelections: map[string]bool{"demo": true},
		Cfg: &config.Manager{
			Config: &types.Config{
				Language: "zh-CN",
			},
		},
	}

	next, _ := m.applyFcitxThemeChoice()
	got := next.(Model)

	if got.State != ViewResult {
		t.Fatalf("state after sync error = %v, want %v", got.State, ViewResult)
	}
}

func containsLine(content, line string) bool {
	for _, item := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if item == line {
			return true
		}
	}
	return false
}
