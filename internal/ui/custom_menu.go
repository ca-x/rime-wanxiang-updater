package ui

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"

	projectassets "rime-wanxiang-updater/assets"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type customMenuItem struct {
	key  string
	icon string
	text string
	desc string
}

var (
	listAvailableFcitxThemes          = availableFcitxThemes
	loadInstalledFcitxThemeSelections = func(themeNames []string) (map[string]bool, error) {
		root, err := fcitxThemeRootPath()
		if err != nil {
			return nil, err
		}
		return installedFcitxThemeSelections(root, themeNames)
	}
	syncInstalledFcitxThemeSelections = func(themeNames []string, selections map[string]bool) error {
		themeFS := projectassets.Fcitx5Themes()
		if themeFS == nil {
			return fmt.Errorf("fcitx5 themes are not embedded")
		}
		root, err := fcitxThemeRootPath()
		if err != nil {
			return err
		}
		return syncInstalledFcitxThemes(themeFS, root, themeNames, selections)
	}
	loadCurrentFcitxThemeConfig = currentFcitxThemeConfig
	setFcitxThemeDefault        = applyFcitxThemeConfig
)

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
	m.ThemePatchSearchQuery = ""
	m.ThemePatchSelections = make(map[string]bool)
	m.ThemePatchDefaultChoice = 0
	m.ThemePatchDefaultKey = ""
	if filePath, err := m.themePatchFilePath(); err == nil {
		if selections, readErr := readThemePatchSelections(filePath); readErr == nil {
			m.ThemePatchSelections = selections
		}
	}
	m.syncThemePatchFilterState()
}

func (m Model) handleThemePatchListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		m.ThemePatchSearchQuery += string(msg.Runes)
		m.syncThemePatchFilterState()
		return m, nil
	case tea.KeyBackspace:
		if m.ThemePatchSearchQuery == "" {
			return m, nil
		}
		queryRunes := []rune(m.ThemePatchSearchQuery)
		m.ThemePatchSearchQuery = string(queryRunes[:len(queryRunes)-1])
		m.syncThemePatchFilterState()
		return m, nil
	case tea.KeyEsc:
		if m.ThemePatchSearchQuery != "" {
			m.ThemePatchSearchQuery = ""
			m.syncThemePatchFilterState()
			return m, nil
		}
		m.State = ViewCustomMenu
		return m, nil
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.ThemePatchChoice > 0 {
			m.ThemePatchChoice--
		}
	case tea.KeyDown:
		if m.ThemePatchChoice < len(m.themePatchFilteredDefinitions())-1 {
			m.ThemePatchChoice++
		}
	case tea.KeySpace:
		m.toggleThemePatchSelection()
	case tea.KeyEnter:
		return m.applyThemePatchPresetChoice()
	}

	return m, nil
}

func (m *Model) syncThemePatchFilterState() {
	filtered := m.themePatchFilteredDefinitions()
	switch {
	case len(filtered) == 0:
		m.ThemePatchChoice = 0
	case m.ThemePatchChoice < 0:
		m.ThemePatchChoice = 0
	case m.ThemePatchChoice >= len(filtered):
		m.ThemePatchChoice = len(filtered) - 1
	}
}

func (m Model) themePatchFilteredDefinitions() []themePatchDefinition {
	query := strings.ToLower(strings.TrimSpace(m.ThemePatchSearchQuery))
	if query == "" {
		return themePatchDefinitions()
	}

	filtered := make([]themePatchDefinition, 0)
	for _, definition := range themePatchDefinitions() {
		if strings.Contains(strings.ToLower(definition.Key), query) ||
			strings.Contains(strings.ToLower(definition.DisplayName), query) {
			filtered = append(filtered, definition)
		}
	}

	return filtered
}

func (m Model) pagedListPageSize() int {
	if m.Height <= 0 {
		return 8
	}

	pageSize := m.Height - 16
	if pageSize < 6 {
		return 6
	}
	if pageSize > 10 {
		return 10
	}

	return pageSize
}

func (m Model) pagedListPageWindow(choice, total int) (start, end, currentPage, totalPages int) {
	if total == 0 {
		return 0, 0, 0, 0
	}

	pageSize := m.pagedListPageSize()
	totalPages = (total + pageSize - 1) / pageSize
	currentPage = (choice / pageSize) + 1
	if currentPage > totalPages {
		currentPage = totalPages
	}

	start = (currentPage - 1) * pageSize
	end = start + pageSize
	if end > total {
		end = total
	}

	return start, end, currentPage, totalPages
}

func (m Model) pagedListPageSummary(choice, filteredCount, totalCount int, pageKey, emptyKey string) string {
	if filteredCount == 0 {
		return m.t(emptyKey, totalCount)
	}

	_, _, currentPage, totalPages := m.pagedListPageWindow(choice, filteredCount)
	return m.t(pageKey, currentPage, totalPages, filteredCount, totalCount)
}

func (m Model) themePatchPageWindow(total int) (start, end, currentPage, totalPages int) {
	return m.pagedListPageWindow(m.ThemePatchChoice, total)
}

func (m Model) themePatchPageSummary(filteredCount int) string {
	totalCount := len(themePatchDefinitions())
	return m.pagedListPageSummary(
		m.ThemePatchChoice,
		filteredCount,
		totalCount,
		"custom.theme_patch.page",
		"custom.theme_patch.page_empty",
	)
}

func (m *Model) toggleThemePatchSelection() {
	definitions := m.themePatchFilteredDefinitions()
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

	if err := syncThemePatchPresets(filePath, m.ThemePatchSelections); err != nil {
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

	searchValue := m.ThemePatchSearchQuery
	searchStyle := lipgloss.NewStyle().Foreground(m.Styles.Foreground)
	if strings.TrimSpace(searchValue) == "" {
		searchValue = m.t("custom.theme_patch.search_placeholder")
		searchStyle = lipgloss.NewStyle().Foreground(m.Styles.Muted)
	}
	searchLabel := lipgloss.NewStyle().
		Foreground(m.Styles.Secondary).
		Bold(true).
		Render(m.t("custom.theme_patch.search_label"))
	searchLine := lipgloss.JoinHorizontal(lipgloss.Left, searchLabel, searchStyle.Render(searchValue))
	b.WriteString(m.renderPanel(searchLine, m.Styles.Border) + "\n\n")

	filtered := m.themePatchFilteredDefinitions()
	var listContent strings.Builder
	if len(filtered) == 0 {
		listContent.WriteString(m.Styles.WarningText.Render(m.t("custom.theme_patch.empty")))
	} else {
		start, end, _, _ := m.themePatchPageWindow(len(filtered))
		for i := start; i < end; i++ {
			definition := filtered[i]
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
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.themePatchPageSummary(len(filtered))) + "\n\n")
	b.WriteString(m.renderHintStrip(
		m.t("custom.theme_patch.hint.nav"),
		m.t("custom.theme_patch.hint.search"),
		m.t("custom.theme_patch.hint.clear"),
		m.t("custom.theme_patch.hint.toggle"),
		m.t("custom.theme_patch.hint.next"),
		m.t("ui.hint.back"),
	))

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
	start, end, _, _ := m.pagedListPageWindow(m.ThemePatchDefaultChoice, len(selected))
	for i := start; i < end; i++ {
		definition := selected[i]
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
	b.WriteString(m.Styles.Hint.Render(m.pagedListPageSummary(
		m.ThemePatchDefaultChoice,
		len(selected),
		len(selected),
		"custom.theme_patch.page",
		"custom.theme_patch.page_empty",
	)) + "\n\n")
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
	themes, err := listAvailableFcitxThemes()
	if err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_error", err)
		m.State = ViewResult
		return m, nil
	}

	selections, err := loadInstalledFcitxThemeSelections(themes)
	if err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_error", err)
		m.State = ViewResult
		return m, nil
	}

	currentConfig, err := loadCurrentFcitxThemeConfig()
	if err != nil {
		currentConfig = FcitxThemeConfig{}
	}

	m.FcitxThemeList = themes
	m.FcitxThemeChoice = 0
	m.FcitxThemeSearchQuery = ""
	m.FcitxThemeSelections = selections
	m.FcitxThemeDefaultChoice = 0
	m.FcitxThemeDefaultKey = fcitxThemeSelectionLight
	m.FcitxThemeLightSelected = ""
	m.FcitxThemeDarkSelected = ""
	m.FcitxThemeCurrent = currentConfig
	m.syncFcitxThemeFilterState()
	m.State = ViewFcitxThemeList
	return m, nil
}

func (m Model) handleFcitxThemeListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		m.FcitxThemeSearchQuery += string(msg.Runes)
		m.syncFcitxThemeFilterState()
		return m, nil
	case tea.KeyBackspace:
		if m.FcitxThemeSearchQuery == "" {
			return m, nil
		}
		queryRunes := []rune(m.FcitxThemeSearchQuery)
		m.FcitxThemeSearchQuery = string(queryRunes[:len(queryRunes)-1])
		m.syncFcitxThemeFilterState()
		return m, nil
	case tea.KeyEsc:
		if m.FcitxThemeSearchQuery != "" {
			m.FcitxThemeSearchQuery = ""
			m.syncFcitxThemeFilterState()
			return m, nil
		}
		m.State = ViewCustomMenu
		return m, nil
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.FcitxThemeChoice > 0 {
			m.FcitxThemeChoice--
		}
	case tea.KeyDown:
		if m.FcitxThemeChoice < len(m.fcitxThemeFilteredList())-1 {
			m.FcitxThemeChoice++
		}
	case tea.KeySpace:
		m.toggleFcitxThemeSelection()
	case tea.KeyEnter:
		return m.applyFcitxThemeChoice()
	}

	return m, nil
}

func (m *Model) syncFcitxThemeFilterState() {
	filtered := m.fcitxThemeFilteredList()
	switch {
	case len(filtered) == 0:
		m.FcitxThemeChoice = 0
	case m.FcitxThemeChoice < 0:
		m.FcitxThemeChoice = 0
	case m.FcitxThemeChoice >= len(filtered):
		m.FcitxThemeChoice = len(filtered) - 1
	}
}

func (m Model) fcitxThemeFilteredList() []string {
	query := strings.ToLower(strings.TrimSpace(m.FcitxThemeSearchQuery))
	if query == "" {
		return m.FcitxThemeList
	}

	filtered := make([]string, 0, len(m.FcitxThemeList))
	for _, themeName := range m.FcitxThemeList {
		if strings.Contains(strings.ToLower(themeName), query) {
			filtered = append(filtered, themeName)
		}
	}

	return filtered
}

func (m Model) fcitxThemePageWindow(total int) (start, end, currentPage, totalPages int) {
	return m.pagedListPageWindow(m.FcitxThemeChoice, total)
}

func (m Model) fcitxThemePageSummary(filteredCount int) string {
	totalCount := len(m.FcitxThemeList)
	return m.pagedListPageSummary(
		m.FcitxThemeChoice,
		filteredCount,
		totalCount,
		"custom.fcitx_theme.page",
		"custom.fcitx_theme.page_empty",
	)
}

func (m *Model) toggleFcitxThemeSelection() {
	filtered := m.fcitxThemeFilteredList()
	if m.FcitxThemeChoice < 0 || m.FcitxThemeChoice >= len(filtered) {
		return
	}
	if m.FcitxThemeSelections == nil {
		m.FcitxThemeSelections = make(map[string]bool)
	}
	themeName := filtered[m.FcitxThemeChoice]
	m.FcitxThemeSelections[themeName] = !m.FcitxThemeSelections[themeName]
	if !m.FcitxThemeSelections[themeName] {
		delete(m.FcitxThemeSelections, themeName)
	}
}

func (m Model) selectedFcitxThemes() []string {
	selected := make([]string, 0)
	for _, themeName := range m.FcitxThemeList {
		if m.FcitxThemeSelections[themeName] {
			selected = append(selected, themeName)
		}
	}
	return selected
}

func (m Model) applyFcitxThemeChoice() (tea.Model, tea.Cmd) {
	selected := m.selectedFcitxThemes()
	if len(selected) == 0 {
		return m, nil
	}

	if m.FcitxThemeCurrent == (FcitxThemeConfig{}) {
		currentConfig, err := loadCurrentFcitxThemeConfig()
		if err == nil {
			m.FcitxThemeCurrent = currentConfig
		}
	}

	if err := syncInstalledFcitxThemeSelections(m.FcitxThemeList, m.FcitxThemeSelections); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.FcitxThemeDefaultKey = fcitxThemeSelectionLight
	m.FcitxThemeLightSelected = ""
	m.FcitxThemeDarkSelected = ""
	m.FcitxThemeDefaultChoice = preferredFcitxThemeChoice(selected, m.FcitxThemeCurrent.Theme)
	m.State = ViewFcitxThemeDefaultList
	return m, nil
}

func (m Model) renderFcitxThemeList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◆ "+m.t("custom.fcitx_theme.title")+" ◆") + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t("custom.fcitx_theme.hint")) + "\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(m.Styles.Muted).
		Render(m.t("custom.fcitx_theme.selected_count", len(m.selectedFcitxThemes()))) + "\n\n")
	b.WriteString(m.renderSummaryCard([][2]string{
		{m.t("custom.fcitx_theme.current_light"), fcitxThemeNameOrUnset(m, m.FcitxThemeCurrent.Theme)},
		{m.t("custom.fcitx_theme.current_dark"), fcitxThemeNameOrUnset(m, m.FcitxThemeCurrent.DarkTheme)},
	}) + "\n\n")

	searchValue := m.FcitxThemeSearchQuery
	searchStyle := lipgloss.NewStyle().Foreground(m.Styles.Foreground)
	if strings.TrimSpace(searchValue) == "" {
		searchValue = m.t("custom.fcitx_theme.search_placeholder")
		searchStyle = lipgloss.NewStyle().Foreground(m.Styles.Muted)
	}
	searchLabel := lipgloss.NewStyle().
		Foreground(m.Styles.Secondary).
		Bold(true).
		Render(m.t("custom.fcitx_theme.search_label"))
	searchLine := lipgloss.JoinHorizontal(lipgloss.Left, searchLabel, searchStyle.Render(searchValue))
	b.WriteString(m.renderPanel(searchLine, m.Styles.Border) + "\n\n")

	filtered := m.fcitxThemeFilteredList()
	var listContent strings.Builder
	if len(filtered) == 0 {
		listContent.WriteString(m.Styles.WarningText.Render(m.t("custom.fcitx_theme.empty")))
	} else {
		start, end, _, _ := m.fcitxThemePageWindow(len(filtered))
		for i := start; i < end; i++ {
			themeName := filtered[i]
			cursor := "  "
			marker := "[ ]"
			if m.FcitxThemeSelections[themeName] {
				marker = "[x]"
			}
			style := lipgloss.NewStyle().
				Foreground(m.Styles.Foreground).
				PaddingLeft(1)

			if m.FcitxThemeChoice == i {
				cursor = "› "
				style = m.Styles.SelectedMenuItem
			}

			listContent.WriteString(style.Render(cursor+marker+" "+themeName) + "\n")
		}
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.fcitxThemePageSummary(len(filtered))) + "\n\n")
	b.WriteString(m.renderHintStrip(
		m.t("custom.fcitx_theme.hint.nav"),
		m.t("custom.fcitx_theme.hint.search"),
		m.t("custom.fcitx_theme.hint.clear"),
		m.t("custom.fcitx_theme.hint.toggle"),
		m.t("custom.fcitx_theme.hint.next"),
		m.t("ui.hint.back"),
	))

	return m.renderScreen(b.String())
}

func (m Model) handleFcitxThemeDefaultListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selected := m.selectedFcitxThemes()
	switch msg.String() {
	case "q", "esc":
		m.State = ViewCustomMenu
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.FcitxThemeDefaultChoice > 0 {
			m.FcitxThemeDefaultChoice--
		}
	case "down", "j":
		if m.FcitxThemeDefaultChoice < len(selected)-1 {
			m.FcitxThemeDefaultChoice++
		}
	case "enter":
		return m.applyFcitxThemeDefaultChoice()
	}

	return m, nil
}

func (m Model) applyFcitxThemeDefaultChoice() (tea.Model, tea.Cmd) {
	selected := m.selectedFcitxThemes()
	if m.FcitxThemeDefaultChoice < 0 || m.FcitxThemeDefaultChoice >= len(selected) {
		return m, nil
	}

	themeName := selected[m.FcitxThemeDefaultChoice]
	if m.FcitxThemeDefaultKey == fcitxThemeSelectionLight {
		m.FcitxThemeLightSelected = themeName
		m.FcitxThemeDefaultKey = fcitxThemeSelectionDark
		m.FcitxThemeDefaultChoice = preferredFcitxThemeChoice(selected, m.FcitxThemeCurrent.DarkTheme, themeName)
		return m, nil
	}

	cfg := FcitxThemeConfig{
		Theme:                m.FcitxThemeLightSelected,
		DarkTheme:            themeName,
		UseDarkTheme:         m.FcitxThemeCurrent.UseDarkTheme,
		FollowSystemDarkMode: boolPtr(true),
	}
	if err := setFcitxThemeDefault(cfg); err != nil {
		m.ResultSuccess = false
		m.ResultSkipped = false
		m.ResultMsg = m.t("custom.result.fcitx_theme_error", err)
		m.State = ViewResult
		return m, nil
	}

	m.FcitxThemeDarkSelected = themeName
	m.State = ViewFcitxThemeDeployPrompt
	return m, nil
}

func (m Model) renderFcitxThemeDefaultList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	titleKey := "custom.fcitx_theme.default_title_light"
	hintKey := "custom.fcitx_theme.default_hint_light"
	summary := [][2]string{
		{m.t("custom.fcitx_theme.current_light"), fcitxThemeNameOrUnset(m, m.FcitxThemeCurrent.Theme)},
		{m.t("custom.fcitx_theme.current_dark"), fcitxThemeNameOrUnset(m, m.FcitxThemeCurrent.DarkTheme)},
	}
	if m.FcitxThemeDefaultKey == fcitxThemeSelectionDark {
		titleKey = "custom.fcitx_theme.default_title_dark"
		hintKey = "custom.fcitx_theme.default_hint_dark"
		summary = [][2]string{
			{m.t("custom.fcitx_theme.selected_light"), fcitxThemeNameOrUnset(m, m.FcitxThemeLightSelected)},
			{m.t("custom.fcitx_theme.current_dark"), fcitxThemeNameOrUnset(m, m.FcitxThemeCurrent.DarkTheme)},
		}
	}

	b.WriteString(m.renderTitle("◆ "+m.t(titleKey)+" ◆") + "\n\n")

	selected := m.selectedFcitxThemes()
	if len(selected) == 0 {
		b.WriteString(m.renderPanel(m.t("custom.fcitx_theme.default_empty"), m.Styles.Warning) + "\n\n")
		b.WriteString(m.renderHintStrip(m.t("ui.hint.back")))
		return m.renderScreen(b.String())
	}

	b.WriteString(m.renderSummaryCard(summary) + "\n\n")

	var listContent strings.Builder
	start, end, _, _ := m.pagedListPageWindow(m.FcitxThemeDefaultChoice, len(selected))
	for i := start; i < end; i++ {
		themeName := selected[i]
		cursor := "  "
		style := lipgloss.NewStyle().
			Foreground(m.Styles.Foreground).
			PaddingLeft(1)

		if m.FcitxThemeDefaultChoice == i {
			cursor = "› "
			style = m.Styles.SelectedMenuItem
		}

		listContent.WriteString(style.Render(cursor+themeName) + "\n")
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.pagedListPageSummary(
		m.FcitxThemeDefaultChoice,
		len(selected),
		len(selected),
		"custom.fcitx_theme.page",
		"custom.fcitx_theme.page_empty",
	)) + "\n\n")
	b.WriteString(m.Styles.Hint.Render(m.t(hintKey)) + "\n\n")
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
	m.ResultMsg = m.t("custom.result.fcitx_theme_deploy_success", m.FcitxThemeLightSelected, m.FcitxThemeDarkSelected)
	m.State = ViewResult
	return m, nil
}

func (m Model) renderFcitxThemeDeployPrompt() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())
	b.WriteString(m.renderTitle("◆ "+m.t("custom.fcitx_theme.deploy_title")+" ◆") + "\n\n")
	b.WriteString(m.renderPanel(m.t("custom.fcitx_theme.deploy_body", m.FcitxThemeLightSelected, m.FcitxThemeDarkSelected), m.Styles.Secondary) + "\n\n")
	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("custom.fcitx_theme.hint.deploy"), m.t("custom.fcitx_theme.hint.return")))

	return m.renderScreen(b.String())
}

func preferredFcitxThemeChoice(selected []string, candidates ...string) int {
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		for index, themeName := range selected {
			if themeName == candidate {
				return index
			}
		}
	}
	return 0
}

func fcitxThemeNameOrUnset(m Model, themeName string) string {
	if strings.TrimSpace(themeName) == "" {
		return m.t("config.value.unset")
	}
	return themeName
}
