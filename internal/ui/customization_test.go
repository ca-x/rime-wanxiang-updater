package ui

import (
	"os"
	"path/filepath"
	"testing"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

func TestThemePatchTargetForPlatform(t *testing.T) {
	tests := []struct {
		name      string
		platform  string
		engines   []string
		wantFile  string
		wantShown bool
	}{
		{
			name:      "darwin requires squirrel",
			platform:  "darwin",
			engines:   []string{"小企鹅"},
			wantShown: false,
		},
		{
			name:      "darwin shows when squirrel is installed alongside fcitx5",
			platform:  "darwin",
			engines:   []string{"小企鹅", "鼠须管"},
			wantFile:  "squirrel.custom.yaml",
			wantShown: true,
		},
		{
			name:      "windows shows for weasel",
			platform:  "windows",
			engines:   []string{"小狼毫"},
			wantFile:  "weasel.custom.yaml",
			wantShown: true,
		},
		{
			name:      "linux never shows",
			platform:  "linux",
			engines:   []string{"fcitx5"},
			wantShown: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotFile, ok := themePatchTargetForPlatform(tt.platform, tt.engines)
			if ok != tt.wantShown {
				t.Fatalf("themePatchTargetForPlatform() shown = %v, want %v", ok, tt.wantShown)
			}
			if gotFile != tt.wantFile {
				t.Fatalf("themePatchTargetForPlatform() file = %q, want %q", gotFile, tt.wantFile)
			}
		})
	}
}

func TestHandleMenuInputOpensCustomMenu(t *testing.T) {
	m := Model{
		State: ViewMenu,
		Cfg: &config.Manager{
			Config: &types.Config{},
		},
	}

	next, _ := m.handleMenuInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}})
	got := next.(Model)
	if got.State != ViewCustomMenu {
		t.Fatalf("handleMenuInput(6) state = %v, want %v", got.State, ViewCustomMenu)
	}
}

func TestWriteThemePatchPresetsMergesExistingPatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "squirrel.custom.yaml")

	existing := []byte("patch:\n  style/font_point: 15\n")
	if err := os.WriteFile(path, existing, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	definitions := themePatchDefinitions()
	if err := writeThemePatchPresets(path, []themePatchDefinition{definitions[0], definitions[4]}); err != nil {
		t.Fatalf("writeThemePatchPresets() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	patch, ok := normalizeStringMap(doc["patch"])
	if !ok {
		t.Fatalf("patch section missing or invalid: %#v", doc["patch"])
	}

	if got := patch["style/font_point"]; got != 15 {
		t.Fatalf("style/font_point = %#v, want 15", got)
	}
	if _, exists := patch["style/color_scheme"]; exists {
		t.Fatalf("style/color_scheme should not be set in preset step: %#v", patch["style/color_scheme"])
	}
	if _, exists := patch["style/color_scheme_dark"]; exists {
		t.Fatalf("style/color_scheme_dark should not be set in preset step: %#v", patch["style/color_scheme_dark"])
	}

	themeNode, ok := normalizeStringMap(patch["preset_color_schemes/jianchun"])
	if !ok {
		t.Fatalf("theme patch missing: %#v", patch["preset_color_schemes/jianchun"])
	}
	if got := themeNode["name"]; got != "简纯" {
		t.Fatalf("theme name = %#v, want %q", got, "简纯")
	}
	if got := themeNode["border_color"]; got != "0xCE7539" {
		t.Fatalf("border_color = %#v, want %q", got, "0xCE7539")
	}

	if _, ok := normalizeStringMap(patch["preset_color_schemes/wechat"]); !ok {
		t.Fatalf("wechat theme patch missing: %#v", patch["preset_color_schemes/wechat"])
	}
}

func TestWriteThemePatchDefaultUpdatesStyleOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "weasel.custom.yaml")

	definitions := themePatchDefinitions()
	if err := writeThemePatchPresets(path, []themePatchDefinition{definitions[0], definitions[4]}); err != nil {
		t.Fatalf("writeThemePatchPresets() error = %v", err)
	}

	if err := writeThemePatchDefault(path, "wechat"); err != nil {
		t.Fatalf("writeThemePatchDefault() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	patch, ok := normalizeStringMap(doc["patch"])
	if !ok {
		t.Fatalf("patch section missing or invalid: %#v", doc["patch"])
	}

	if got := patch["style/color_scheme"]; got != "wechat" {
		t.Fatalf("style/color_scheme = %#v, want %q", got, "wechat")
	}
	if got := patch["style/color_scheme_dark"]; got != "wechat" {
		t.Fatalf("style/color_scheme_dark = %#v, want %q", got, "wechat")
	}
	if _, ok := normalizeStringMap(patch["preset_color_schemes/jianchun"]); !ok {
		t.Fatalf("existing preset removed: %#v", patch["preset_color_schemes/jianchun"])
	}
	if _, ok := normalizeStringMap(patch["preset_color_schemes/wechat"]); !ok {
		t.Fatalf("selected preset missing: %#v", patch["preset_color_schemes/wechat"])
	}
}

func TestHandleThemePatchDeployPromptInput(t *testing.T) {
	oldDeploy := deployThemePatch
	defer func() { deployThemePatch = oldDeploy }()

	called := false
	deployThemePatch = func(cfg *types.Config) error {
		called = true
		return nil
	}

	m := Model{
		State: ViewThemePatchDeployPrompt,
		Cfg: &config.Manager{
			Config: &types.Config{},
		},
	}

	next, _ := m.handleThemePatchDeployPromptInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	got := next.(Model)
	if got.State != ViewCustomMenu {
		t.Fatalf("handleThemePatchDeployPromptInput(non-enter) state = %v, want %v", got.State, ViewCustomMenu)
	}
	if called {
		t.Fatalf("deployThemePatch should not run on non-enter input")
	}

	next, _ = m.handleThemePatchDeployPromptInput(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)
	if got.State != ViewResult {
		t.Fatalf("handleThemePatchDeployPromptInput(enter) state = %v, want %v", got.State, ViewResult)
	}
	if !called {
		t.Fatalf("deployThemePatch should run on enter")
	}
}
