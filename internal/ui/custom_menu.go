package ui

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type customMenuItem struct {
	key  string
	icon string
	text string
	desc string
}

func themePatchTargetForPlatform(platform string, installedEngines []string) (string, string, bool) {
	switch platform {
	case "windows":
		if slices.Contains(installedEngines, "小狼毫") {
			return "小狼毫", "weasel.custom.yaml", true
		}
	case "darwin":
		if slices.Contains(installedEngines, "鼠须管") {
			return "鼠须管", "squirrel.custom.yaml", true
		}
	}

	return "", "", false
}

func (m Model) customMenuItems() []customMenuItem {
	items := []customMenuItem{
		{
			key:  "program_tui",
			icon: "◐",
			text: m.t("custom.program_tui.title"),
			desc: m.t("custom.program_tui.desc"),
		},
	}

	if _, _, ok := themePatchTargetForPlatform(runtime.GOOS, m.Cfg.Config.InstalledEngines); ok {
		items = append(items, customMenuItem{
			key:  "theme_patch",
			icon: "◈",
			text: m.t("custom.theme_patch.title"),
			desc: m.t("custom.theme_patch.desc"),
		})
	}
	if fcitxThemeSupportedForPlatform(runtime.GOOS, m.Cfg.Config.InstalledEngines) {
		items = append(items, customMenuItem{
			key:  "fcitx_theme",
			icon: "◆",
			text: m.t("custom.fcitx_theme.title"),
			desc: m.t("custom.fcitx_theme.desc"),
		})
	}

	return items
}

func (m Model) handleCustomMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := m.customMenuItems()

	switch msg.String() {
	case "q", "esc":
		m.State = ViewMenu
		m.CustomMenuChoice = 0
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.CustomMenuChoice > 0 {
			m.CustomMenuChoice--
		}
	case "down", "j":
		if m.CustomMenuChoice < len(items)-1 {
			m.CustomMenuChoice++
		}
	case "enter":
		return m.applyCustomMenuChoice()
	default:
		if n, err := strconv.Atoi(msg.String()); err == nil && n >= 1 && n <= len(items) {
			m.CustomMenuChoice = n - 1
			return m.applyCustomMenuChoice()
		}
	}

	return m, nil
}

func (m Model) applyCustomMenuChoice() (tea.Model, tea.Cmd) {
	items := m.customMenuItems()
	if m.CustomMenuChoice < 0 || m.CustomMenuChoice >= len(items) {
		return m, nil
	}

	switch items[m.CustomMenuChoice].key {
	case "program_tui":
		m.InitThemeListView("theme_quick")
		m.State = ViewThemeList
	case "theme_patch":
		m.InitThemePatchListView()
		m.State = ViewThemePatchList
	case "fcitx_theme":
		return m.openFcitxThemeList()
	}

	return m, nil
}

func (m Model) renderCustomMenu() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("⚙ "+m.t("custom.menu.title")+" ⚙") + "\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(m.Styles.Muted).Render(m.t("custom.menu.subtitle")) + "\n\n")

	items := m.customMenuItems()
	menuItems := make([]menuEntry, 0, len(items))
	for _, item := range items {
		menuItems = append(menuItems, menuEntry{
			icon: item.icon,
			text: item.text,
			desc: item.desc,
		})
	}

	b.WriteString(m.renderChoiceList(menuItems, m.CustomMenuChoice) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("ui.hint.nav"), m.t("ui.hint.select"), m.t("ui.hint.back")))

	return m.renderScreen(b.String())
}

func (m *Model) InitThemePatchListView() {
	m.ThemePatchChoice = 0
	m.ThemePatchSelections = make(map[string]bool)
	m.ThemePatchDefaultChoice = 0
	m.ThemePatchDefaultKey = ""
}

func (m Model) handleThemePatchListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.State = ViewCustomMenu
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.ThemePatchChoice > 0 {
			m.ThemePatchChoice--
		}
	case "down", "j":
		if m.ThemePatchChoice < len(themePatchDefinitions())-1 {
			m.ThemePatchChoice++
		}
	case " ":
		m.toggleThemePatchSelection()
	case "enter":
		return m.applyThemePatchPresetChoice()
	}

	return m, nil
}

func (m *Model) toggleThemePatchSelection() {
	definitions := themePatchDefinitions()
	if m.ThemePatchChoice < 0 || m.ThemePatchChoice >= len(definitions) {
		return
	}
	if m.ThemePatchSelections == nil {
		m.ThemePatchSelections = make(map[string]bool)
	}
	key := definitions[m.ThemePatchChoice].Key
	m.ThemePatchSelections[key] = !m.ThemePatchSelections[key]
	if !m.ThemePatchSelections[key] {
		delete(m.ThemePatchSelections, key)
	}
}

func (m Model) selectedThemePatchDefinitions() []themePatchDefinition {
	var selected []themePatchDefinition
	for _, definition := range themePatchDefinitions() {
		if m.ThemePatchSelections[definition.Key] {
			selected = append(selected, definition)
		}
	}
	return selected
}

func (m Model) applyThemePatchPresetChoice() (tea.Model, tea.Cmd) {
	selected := m.selectedThemePatchDefinitions()
	if len(selected) == 0 {
		return m, nil
	}

	filePath, err := m.themePatchFilePath()
	if err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.patch_path_error", err)
		m.State = ViewResult
		return m, nil
	}

	if err := writeThemePatchPresets(filePath, selected); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.patch_write_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.ThemePatchDefaultChoice = 0
	m.ThemePatchDefaultKey = ""
	m.State = ViewThemePatchDefaultList
	return m, nil
}

func (m Model) renderThemePatchList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◈ "+m.t("custom.theme_patch.title")+" ◈") + "\n\n")

	if filePath, err := m.themePatchFilePath(); err == nil {
		b.WriteString(lipgloss.NewStyle().
			Foreground(m.Styles.Muted).
			Render(m.t("custom.theme_patch.target", filePath)) + "\n\n")
	}

	selectedCount := len(m.selectedThemePatchDefinitions())
	b.WriteString(lipgloss.NewStyle().
		Foreground(m.Styles.Muted).
		Render(m.t("custom.theme_patch.selected_count", selectedCount)) + "\n\n")

	var listContent strings.Builder
	for i, definition := range themePatchDefinitions() {
		cursor := "  "
		marker := "[ ]"
		if m.ThemePatchSelections[definition.Key] {
			marker = "[x]"
		}
		style := lipgloss.NewStyle().
			Foreground(m.Styles.Foreground).
			PaddingLeft(1)

		if m.ThemePatchChoice == i {
			cursor = "› "
			style = m.Styles.SelectedMenuItem
		}

		line := fmt.Sprintf("%s %s (%s)", marker, definition.Key, definition.DisplayName)
		listContent.WriteString(style.Render(cursor+line) + "\n")
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("custom.theme_patch.hint")) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("ui.hint.nav"), m.t("custom.theme_patch.hint.toggle"), m.t("custom.theme_patch.hint.next"), m.t("ui.hint.back")))

	return m.renderScreen(b.String())
}

func (m Model) handleThemePatchDefaultListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selected := m.selectedThemePatchDefinitions()
	switch msg.String() {
	case "q", "esc":
		m.State = ViewCustomMenu
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.ThemePatchDefaultChoice > 0 {
			m.ThemePatchDefaultChoice--
		}
	case "down", "j":
		if m.ThemePatchDefaultChoice < len(selected)-1 {
			m.ThemePatchDefaultChoice++
		}
	case "enter":
		return m.applyThemePatchDefaultChoice()
	}

	return m, nil
}

func (m Model) applyThemePatchDefaultChoice() (tea.Model, tea.Cmd) {
	selected := m.selectedThemePatchDefinitions()
	if m.ThemePatchDefaultChoice < 0 || m.ThemePatchDefaultChoice >= len(selected) {
		return m, nil
	}

	filePath, err := m.themePatchFilePath()
	if err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.patch_path_error", err)
		m.State = ViewResult
		return m, nil
	}

	defaultKey := selected[m.ThemePatchDefaultChoice].Key
	if err := writeThemePatchDefault(filePath, defaultKey); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.patch_write_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.ThemePatchDefaultKey = defaultKey
	m.State = ViewThemePatchDeployPrompt
	return m, nil
}

func (m Model) renderThemePatchDefaultList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◈ "+m.t("custom.theme_patch.default_title")+" ◈") + "\n\n")

	selected := m.selectedThemePatchDefinitions()
	if len(selected) == 0 {
		b.WriteString(m.renderPanel(m.t("custom.theme_patch.default_empty"), m.Styles.Warning) + "\n\n")
		b.WriteString(m.renderHintStrip(m.t("ui.hint.back")))
		return m.renderScreen(b.String())
	}

	var listContent strings.Builder
	for i, definition := range selected {
		cursor := "  "
		style := lipgloss.NewStyle().
			Foreground(m.Styles.Foreground).
			PaddingLeft(1)

		if m.ThemePatchDefaultChoice == i {
			cursor = "› "
			style = m.Styles.SelectedMenuItem
		}

		line := fmt.Sprintf("%s (%s)", definition.Key, definition.DisplayName)
		listContent.WriteString(style.Render(cursor+line) + "\n")
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("custom.theme_patch.default_hint")) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("ui.hint.nav"), m.t("ui.hint.apply_theme"), m.t("ui.hint.back")))

	return m.renderScreen(b.String())
}

func (m Model) handleThemePatchDeployPromptInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		return m.runThemePatchDeploy()
	default:
		m.State = ViewCustomMenu
		return m, nil
	}
}

func (m Model) runThemePatchDeploy() (tea.Model, tea.Cmd) {
	if m.Cfg == nil || m.Cfg.Config == nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.patch_deploy_error", fmt.Errorf("config unavailable"))
		m.State = ViewResult
		return m, nil
	}

	if err := deployThemePatch(m.Cfg.Config); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.patch_deploy_error", err)
		m.State = ViewResult
		return m, nil
	}

	filePath, pathErr := m.themePatchFilePath()
	if pathErr != nil {
		filePath = ""
	}
	m.ResultSuccess = true
	m.ResultSkipped = false
	m.ResultMsg = m.t("custom.result.patch_deploy_success", m.ThemePatchDefaultKey, filePath)
	m.State = ViewResult
	return m, nil
}

func (m Model) renderThemePatchDeployPrompt() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◈ "+m.t("custom.theme_patch.deploy_title")+" ◈") + "\n\n")

	content := m.t("custom.theme_patch.deploy_body", m.ThemePatchDefaultKey)
	b.WriteString(m.renderPanel(content, m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("custom.theme_patch.hint.deploy"), m.t("custom.theme_patch.hint.return")))

	return m.renderScreen(b.String())
}

func (m Model) openFcitxThemeList() (tea.Model, tea.Cmd) {
	themes, err := availableFcitxThemes()
	if err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.FcitxThemeList = themes
	m.FcitxThemeChoice = 0
	m.FcitxThemeSelected = ""
	m.State = ViewFcitxThemeList
	return m, nil
}

func (m Model) handleFcitxThemeListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.State = ViewCustomMenu
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.FcitxThemeChoice > 0 {
			m.FcitxThemeChoice--
		}
	case "down", "j":
		if m.FcitxThemeChoice < len(m.FcitxThemeList)-1 {
			m.FcitxThemeChoice++
		}
	case "enter":
		return m.applyFcitxThemeChoice()
	}

	return m, nil
}

func (m Model) applyFcitxThemeChoice() (tea.Model, tea.Cmd) {
	if m.FcitxThemeChoice < 0 || m.FcitxThemeChoice >= len(m.FcitxThemeList) {
		return m, nil
	}

	themeName := m.FcitxThemeList[m.FcitxThemeChoice]
	if err := installAndSetFcitxTheme(themeName); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.FcitxThemeSelected = themeName
	m.State = ViewFcitxThemeDeployPrompt
	return m, nil
}

func (m Model) renderFcitxThemeList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◆ "+m.t("custom.fcitx_theme.title")+" ◆") + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("custom.fcitx_theme.hint")) + "\n\n")

	var listContent strings.Builder
	for i, themeName := range m.FcitxThemeList {
		cursor := "  "
		style := lipgloss.NewStyle().
			Foreground(m.Styles.Foreground).
			PaddingLeft(1)

		if m.FcitxThemeChoice == i {
			cursor = "› "
			style = m.Styles.SelectedMenuItem
		}

		listContent.WriteString(style.Render(cursor+themeName) + "\n")
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("ui.hint.nav"), m.t("ui.hint.apply_theme"), m.t("ui.hint.back")))

	return m.renderScreen(b.String())
}

func (m Model) handleFcitxThemeDeployPromptInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		return m.runFcitxThemeDeploy()
	default:
		m.State = ViewCustomMenu
		return m, nil
	}
}

func (m Model) runFcitxThemeDeploy() (tea.Model, tea.Cmd) {
	if m.Cfg == nil || m.Cfg.Config == nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_deploy_error", fmt.Errorf("config unavailable"))
		m.State = ViewResult
		return m, nil
	}

	if err := deployFcitxTheme(m.Cfg.Config); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_deploy_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.ResultSuccess = true
	m.ResultSkipped = false
	m.ResultMsg = m.t("custom.result.fcitx_theme_deploy_success", m.FcitxThemeSelected)
	m.State = ViewResult
	return m, nil
}

func (m Model) renderFcitxThemeDeployPrompt() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◆ "+m.t("custom.fcitx_theme.deploy_title")+" ◆") + "\n\n")
	b.WriteString(m.renderPanel(m.t("custom.fcitx_theme.deploy_body", m.FcitxThemeSelected), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("custom.fcitx_theme.hint.deploy"), m.t("custom.fcitx_theme.hint.return")))

	return m.renderScreen(b.String())
}
