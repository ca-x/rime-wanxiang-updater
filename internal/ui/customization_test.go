package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/theme"
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

func TestReadThemePatchSelectionsLoadsExistingBuiltinThemes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "squirrel.custom.yaml")
	content := []byte(`patch:
  preset_color_schemes/jianchun:
    name: 简纯
  preset_color_schemes/custom_theme:
    name: 自定义
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	selections, err := readThemePatchSelections(path)
	if err != nil {
		t.Fatalf("readThemePatchSelections() error = %v", err)
	}

	if !selections["jianchun"] {
		t.Fatalf("expected builtin theme to be selected: %#v", selections)
	}
	if selections["custom_theme"] {
		t.Fatalf("non-builtin theme should not be selected: %#v", selections)
	}
}

func TestSyncThemePatchPresetsRemovesUnselectedBuiltinThemesAndKeepsCustom(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "weasel.custom.yaml")
	content := []byte(`patch:
  preset_color_schemes/jianchun:
    name: 简纯
  preset_color_schemes/wechat:
    name: 微信
  preset_color_schemes/custom_theme:
    name: 自定义
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := syncThemePatchPresets(path, map[string]bool{
		"wechat": true,
	}); err != nil {
		t.Fatalf("syncThemePatchPresets() error = %v", err)
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

	if _, ok := patch["preset_color_schemes/jianchun"]; ok {
		t.Fatalf("unselected builtin theme should be removed: %#v", patch["preset_color_schemes/jianchun"])
	}
	if _, ok := normalizeStringMap(patch["preset_color_schemes/wechat"]); !ok {
		t.Fatalf("selected builtin theme should remain: %#v", patch["preset_color_schemes/wechat"])
	}
	if _, ok := normalizeStringMap(patch["preset_color_schemes/custom_theme"]); !ok {
		t.Fatalf("custom theme should be preserved: %#v", patch["preset_color_schemes/custom_theme"])
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

func TestThemePatchFlowSpaceRuneSelectionCanBeConfirmedAndDeployed(t *testing.T) {
	oldTargetResolver := themePatchTargetResolver
	oldDataDirResolver := themePatchDataDir
	oldDeploy := deployThemePatch
	defer func() {
		themePatchTargetResolver = oldTargetResolver
		themePatchDataDir = oldDataDirResolver
		deployThemePatch = oldDeploy
	}()

	dir := t.TempDir()
	themePatchTargetResolver = func(platform string, installedEngines []string) (string, string, bool) {
		return "小狼毫", "weasel.custom.yaml", true
	}
	themePatchDataDir = func(engineName string) string {
		if engineName != "小狼毫" {
			t.Fatalf("engineName = %q, want %q", engineName, "小狼毫")
		}
		return dir
	}

	deployCalled := false
	deployThemePatch = func(cfg *types.Config) error {
		deployCalled = true
		return nil
	}

	m := newThemePatchTestModel(t)
	m.Cfg.Config.InstalledEngines = []string{"小狼毫"}
	m.InitThemePatchListView()
	m.ThemePatchSearchQuery = "lumk"
	m.syncThemePatchFilterState()

	next, _ := m.handleThemePatchListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	got := next.(Model)
	if !got.ThemePatchSelections["Lumk_light"] {
		t.Fatalf("ThemePatchSelections should contain filtered key %q: %#v", "Lumk_light", got.ThemePatchSelections)
	}

	next, _ = got.handleThemePatchListInput(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)
	if got.State != ViewThemePatchDefaultList {
		t.Fatalf("state after preset confirm = %v, want %v", got.State, ViewThemePatchDefaultList)
	}

	path := filepath.Join(dir, "weasel.custom.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after preset confirm error = %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("yaml.Unmarshal() after preset confirm error = %v", err)
	}

	patch, ok := normalizeStringMap(doc["patch"])
	if !ok {
		t.Fatalf("patch section missing or invalid after preset confirm: %#v", doc["patch"])
	}
	if _, ok := normalizeStringMap(patch["preset_color_schemes/Lumk_light"]); !ok {
		t.Fatalf("selected preset missing after preset confirm: %#v", patch["preset_color_schemes/Lumk_light"])
	}

	next, _ = got.handleThemePatchDefaultListInput(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)
	if got.State != ViewThemePatchDeployPrompt {
		t.Fatalf("state after default confirm = %v, want %v", got.State, ViewThemePatchDeployPrompt)
	}
	if got.ThemePatchDefaultKey != "Lumk_light" {
		t.Fatalf("ThemePatchDefaultKey = %q, want %q", got.ThemePatchDefaultKey, "Lumk_light")
	}

	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after default confirm error = %v", err)
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("yaml.Unmarshal() after default confirm error = %v", err)
	}

	patch, ok = normalizeStringMap(doc["patch"])
	if !ok {
		t.Fatalf("patch section missing or invalid after default confirm: %#v", doc["patch"])
	}
	if got := patch["style/color_scheme"]; got != "Lumk_light" {
		t.Fatalf("style/color_scheme = %#v, want %q", got, "Lumk_light")
	}
	if got := patch["style/color_scheme_dark"]; got != "Lumk_light" {
		t.Fatalf("style/color_scheme_dark = %#v, want %q", got, "Lumk_light")
	}

	next, _ = got.handleThemePatchDeployPromptInput(tea.KeyMsg{Type: tea.KeyEnter})
	got = next.(Model)
	if got.State != ViewResult {
		t.Fatalf("state after deploy confirm = %v, want %v", got.State, ViewResult)
	}
	if !got.ResultSuccess {
		t.Fatalf("ResultSuccess after deploy confirm = %v, want true", got.ResultSuccess)
	}
	if !deployCalled {
		t.Fatalf("deployThemePatch should run on deploy confirm")
	}
}

func TestHandleThemePatchListInputFiltersResultsAndEscClearsSearch(t *testing.T) {
	m := newThemePatchTestModel(t)
	m.InitThemePatchListView()

	next, _ := m.handleThemePatchListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
	got := next.(Model)
	if got.ThemePatchSearchQuery != "w" {
		t.Fatalf("ThemePatchSearchQuery = %q, want %q", got.ThemePatchSearchQuery, "w")
	}
	if len(got.themePatchFilteredDefinitions()) == 0 {
		t.Fatalf("themePatchFilteredDefinitions() should not be empty after search")
	}

	next, _ = got.handleThemePatchListInput(tea.KeyMsg{Type: tea.KeyEsc})
	got = next.(Model)
	if got.State != ViewThemePatchList {
		t.Fatalf("handleThemePatchListInput(esc with search) state = %v, want %v", got.State, ViewThemePatchList)
	}
	if got.ThemePatchSearchQuery != "" {
		t.Fatalf("ThemePatchSearchQuery after esc = %q, want empty", got.ThemePatchSearchQuery)
	}
}

func TestHandleThemePatchListInputSpaceSelectsFilteredTheme(t *testing.T) {
	m := newThemePatchTestModel(t)
	m.InitThemePatchListView()
	m.ThemePatchSearchQuery = "lumk"
	m.syncThemePatchFilterState()

	next, _ := m.handleThemePatchListInput(tea.KeyMsg{Type: tea.KeySpace})
	got := next.(Model)

	if !got.ThemePatchSelections["Lumk_light"] {
		t.Fatalf("ThemePatchSelections should contain filtered key %q: %#v", "Lumk_light", got.ThemePatchSelections)
	}
	if got.ThemePatchSelections["jianchun"] {
		t.Fatalf("ThemePatchSelections should not contain unfiltered first key: %#v", got.ThemePatchSelections)
	}
}

func TestHandleThemePatchListInputSpaceRuneSelectsFilteredTheme(t *testing.T) {
	m := newThemePatchTestModel(t)
	m.InitThemePatchListView()
	m.ThemePatchSearchQuery = "lumk"
	m.syncThemePatchFilterState()

	next, _ := m.handleThemePatchListInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	got := next.(Model)

	if !got.ThemePatchSelections["Lumk_light"] {
		t.Fatalf("ThemePatchSelections should contain filtered key %q: %#v", "Lumk_light", got.ThemePatchSelections)
	}
	if got.ThemePatchSearchQuery != "lumk" {
		t.Fatalf("ThemePatchSearchQuery = %q, want %q", got.ThemePatchSearchQuery, "lumk")
	}
}

func TestHandleThemePatchListInputClampsChoiceAfterFiltering(t *testing.T) {
	m := newThemePatchTestModel(t)
	m.InitThemePatchListView()
	m.ThemePatchChoice = len(themePatchDefinitions()) - 1
	m.ThemePatchSearchQuery = "jianchun"
	m.syncThemePatchFilterState()

	if m.ThemePatchChoice != 0 {
		t.Fatalf("ThemePatchChoice after filtering = %d, want 0", m.ThemePatchChoice)
	}
}

func TestRenderThemePatchListShowsOnlyCurrentPage(t *testing.T) {
	m := newThemePatchTestModel(t)
	m.InitThemePatchListView()
	m.Width = 80
	m.Height = 18

	rendered := m.renderThemePatchList()

	if !strings.Contains(rendered, "简纯") {
		t.Fatalf("renderThemePatchList() should include first page item: %q", rendered)
	}
	if strings.Contains(rendered, "高仿暗色 macOS 14") {
		t.Fatalf("renderThemePatchList() should not include items outside current page: %q", rendered)
	}
	if !strings.Contains(rendered, "第 1/") {
		t.Fatalf("renderThemePatchList() should include pagination summary: %q", rendered)
	}
}

func TestRenderThemePatchDefaultListShowsOnlyCurrentPage(t *testing.T) {
	m := newThemePatchTestModel(t)
	m.State = ViewThemePatchDefaultList
	m.Width = 80
	m.Height = 18
	m.ThemePatchSelections = map[string]bool{
		"jianchun":     true,
		"win11_light":  true,
		"win11_dark":   true,
		"mac_light":    true,
		"wechat":       true,
		"mac_dark":     true,
		"starcraft":    true,
	}

	rendered := m.renderThemePatchDefaultList()

	if !strings.Contains(rendered, "jianchun") {
		t.Fatalf("renderThemePatchDefaultList() should include first page item: %q", rendered)
	}
	if strings.Contains(rendered, "starcraft") {
		t.Fatalf("renderThemePatchDefaultList() should not include items outside current page: %q", rendered)
	}
	if !strings.Contains(rendered, "第 1/") {
		t.Fatalf("renderThemePatchDefaultList() should include pagination summary: %q", rendered)
	}
}

func newThemePatchTestModel(t *testing.T) Model {
	t.Helper()

	themeMgr := theme.NewManager()
	if err := themeMgr.SetTheme("one-dark"); err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}

	return Model{
		State:        ViewThemePatchList,
		ThemeManager: themeMgr,
		Styles:       DefaultStyles(themeMgr),
		Cfg: &config.Manager{
			Config: &types.Config{
				Language:         "zh-CN",
				InstalledEngines: []string{},
			},
		},
	}
}
