package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// handleThemeListInput 处理主题列表输入
func (m Model) handleThemeListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.State = ViewConfig
		return m, nil

	case "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.ThemeListChoice > 0 {
			m.ThemeListChoice--
		}

	case "down", "j":
		maxChoice := len(m.ThemeList) - 1
		if m.ThemeListChoice < maxChoice {
			m.ThemeListChoice++
		}

	case "enter", " ":
		return m.applyThemeChoice()
	}

	return m, nil
}

// applyThemeChoice 应用选中的主题
func (m Model) applyThemeChoice() (tea.Model, tea.Cmd) {
	if m.ThemeListChoice >= 0 && m.ThemeListChoice < len(m.ThemeList) {
		themeName := m.ThemeList[m.ThemeListChoice]

		// 根据编辑的是哪个配置来应用主题
		switch m.EditingKey {
		case "theme_dark":
			m.Cfg.Config.ThemeDark = themeName
			if m.Cfg.Config.ThemeAdaptive {
				m.ThemeManager.SetAdaptiveTheme(m.Cfg.Config.ThemeLight, themeName)
			}
		case "theme_light":
			m.Cfg.Config.ThemeLight = themeName
			if m.Cfg.Config.ThemeAdaptive {
				m.ThemeManager.SetAdaptiveTheme(themeName, m.Cfg.Config.ThemeDark)
			}
		case "theme_fixed":
			m.Cfg.Config.ThemeFixed = themeName
			if !m.Cfg.Config.ThemeAdaptive {
				m.ThemeManager.SetTheme(themeName)
			}
		case "theme_quick":
			// 快速切换主题：关闭自适应模式并设置固定主题
			m.Cfg.Config.ThemeAdaptive = false
			m.Cfg.Config.ThemeFixed = themeName
			m.ThemeManager.SetTheme(themeName)
		}

		// 保存配置
		if err := m.Cfg.SaveConfig(); err != nil {
			m.Err = err
		}

		// 刷新样式
		m.Styles = DefaultStyles(m.ThemeManager)
	}

	// 快速切换后返回主菜单，其他情况返回配置页面
	if m.EditingKey == "theme_quick" {
		m.State = ViewMenu
	} else {
		m.State = ViewConfig
	}
	m.EditingKey = ""
	return m, nil
}

// InitThemeListView 初始化主题列表视图
func (m *Model) InitThemeListView(editingKey string) {
	m.EditingKey = editingKey
	m.ThemeListChoice = 0

	// 根据编辑的配置项选择显示哪些主题
	switch editingKey {
	case "theme_dark":
		m.ThemeList = m.ThemeManager.ListDark()
		// 找到当前选中的主题
		for i, name := range m.ThemeList {
			if name == m.Cfg.Config.ThemeDark {
				m.ThemeListChoice = i
				break
			}
		}
	case "theme_light":
		m.ThemeList = m.ThemeManager.ListLight()
		for i, name := range m.ThemeList {
			if name == m.Cfg.Config.ThemeLight {
				m.ThemeListChoice = i
				break
			}
		}
	case "theme_fixed":
		m.ThemeList = m.ThemeManager.List()
		for i, name := range m.ThemeList {
			if name == m.Cfg.Config.ThemeFixed {
				m.ThemeListChoice = i
				break
			}
		}
	default:
		m.ThemeList = m.ThemeManager.List()
	}
}

// renderThemeList 渲染主题列表
func (m Model) renderThemeList() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	var titleText string
	switch m.EditingKey {
	case "theme_dark":
		titleText = "🌙 " + m.t("theme.select.dark") + " 🌙"
	case "theme_light":
		titleText = "☀️ " + m.t("theme.select.light") + " ☀️"
	default:
		titleText = "🎨 " + m.t("theme.select.title") + " 🎨"
	}
	b.WriteString(m.renderTitle(titleText) + "\n\n")

	currentInfo := lipgloss.NewStyle().
		Foreground(m.Styles.Primary).
		Render(m.t("theme.current", m.ThemeManager.CurrentName()))

	bgInfo := ""
	if m.ThemeManager.IsAdaptive() {
		bg := m.ThemeManager.Background()
		bgType := m.t("theme.bg.dark")
		if !bg.IsDark() {
			bgType = m.t("theme.bg.light")
		}
		bgInfo = lipgloss.NewStyle().
			Foreground(m.Styles.Warning).
			Render(m.t("theme.adaptive.current", bgType))
	}

	b.WriteString(currentInfo + bgInfo + "\n\n")

	var listContent strings.Builder
	for i, themeName := range m.ThemeList {
		scheme, _ := m.ThemeManager.GetScheme(themeName)
		var displayName string
		if scheme != nil {
			displayName = fmt.Sprintf("%s (%s)", themeName, scheme.Scheme)
		} else {
			displayName = themeName
		}

		if scheme != nil {
			preview := renderThemePreview(scheme.Base08, scheme.Base0B, scheme.Base0D, scheme.Base0E)
			displayName = displayName + " " + preview
		}

		cursor := "  "
		style := lipgloss.NewStyle().
			Foreground(m.Styles.Foreground).
			PaddingLeft(1)

		if m.ThemeListChoice == i {
			cursor = "› "
			style = m.Styles.SelectedMenuItem
		}

		listContent.WriteString(style.Render(cursor+displayName) + "\n")
	}

	b.WriteString(m.renderPanel(strings.TrimSuffix(listContent.String(), "\n"), m.Styles.Secondary) + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")

	var hint string
	if m.EditingKey == "theme_quick" {
		hint = m.Styles.Hint.Render(m.t("theme.quick_hint"))
	} else {
		hint = m.Styles.Hint.Render(m.t("theme.hint"))
	}
	b.WriteString(hint + "\n\n")
	b.WriteString(m.renderHintStrip(m.t("ui.hint.nav"), m.t("ui.hint.apply_theme"), m.t("ui.hint.back")))

	return m.renderScreen(b.String())
}

// renderThemePreview 渲染主题颜色预览
func renderThemePreview(red, green, blue, magenta string) string {
	r := lipgloss.NewStyle().Foreground(lipgloss.Color("#" + red)).Render("●")
	g := lipgloss.NewStyle().Foreground(lipgloss.Color("#" + green)).Render("●")
	bl := lipgloss.NewStyle().Foreground(lipgloss.Color("#" + blue)).Render("●")
	m := lipgloss.NewStyle().Foreground(lipgloss.Color("#" + magenta)).Render("●")
	return r + g + bl + m
}
