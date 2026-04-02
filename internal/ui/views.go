package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"rime-wanxiang-updater/internal/termcolor"
	"rime-wanxiang-updater/internal/types"
	"rime-wanxiang-updater/internal/version"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) contentWidth(defaultWidth int) int {
	if m.Width <= 0 {
		return defaultWidth
	}

	width := m.Width - 10
	if width < 48 {
		return 48
	}
	if width > defaultWidth {
		return defaultWidth
	}

	return width
}

func (m Model) pageWidth() int {
	return m.contentWidth(64)
}

func (m Model) renderPanel(content string, border lipgloss.Color) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Background(m.Styles.Surface).
		Padding(1, 2).
		Width(m.pageWidth()).
		Render(content)
}

func (m Model) renderScreen(content string) string {
	screenWidth := m.pageWidth() + 6
	if m.Width > screenWidth {
		screenWidth = m.Width
	}

	return m.Styles.Container.
		Width(screenWidth).
		Render(content)
}

func (m Model) renderLabeledChips(items [][2]string) string {
	return m.renderLabeledChipsWithWidth(m.pageWidth(), items)
}

func (m Model) renderLabeledChipsWithWidth(totalWidth int, items [][2]string) string {
	chips := make([]string, 0, len(items))
	if len(items) == 0 {
		return ""
	}

	gap := 2
	chipChrome := 4 // 2 columns of border + 2 columns of horizontal padding
	chipWidth := (totalWidth - (len(items)-1)*gap - len(items)*chipChrome) / len(items)
	if chipWidth < 12 {
		chipWidth = 12
	}

	for _, item := range items {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			m.Styles.StatusKey.Render(item[0]),
			m.Styles.StatusValue.Render(item[1]),
		)

		chip := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Styles.Muted).
			Background(m.Styles.Surface).
			Padding(0, 1).
			Width(chipWidth).
			Height(2).
			Render(content)
		chips = append(chips, chip)
	}

	row := make([]string, 0, len(chips)*2-1)
	for i, chip := range chips {
		if i > 0 {
			row = append(row, strings.Repeat(" ", gap))
		}
		row = append(row, chip)
	}

	return lipgloss.NewStyle().
		Width(totalWidth).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, row...))
}

func (m Model) renderHeaderBlock() string {
	var b strings.Builder

	brand := lipgloss.NewStyle().
		Foreground(m.Styles.Primary).
		Bold(true).
		Render("Rime Wanxiang Updater")

	build := lipgloss.NewStyle().
		Foreground(m.Styles.Muted).
		Render(version.GetVersion())

	line := lipgloss.JoinHorizontal(lipgloss.Center, brand, "  ", build)
	b.WriteString(lipgloss.NewStyle().Width(m.pageWidth()).Align(lipgloss.Center).Render(line))
	b.WriteString("\n")
	b.WriteString(m.Styles.ScanLine.Render(scanLine))
	b.WriteString("\n\n")

	return b.String()
}

func (m Model) configuredSourceLabel() string {
	if m.Cfg.Config.UseMirror {
		return m.sourceLabel("CNB 镜像")
	}

	return m.sourceLabel("GitHub 官方源")
}

func (m Model) autoUpdateStatusLabel() string {
	if !m.Cfg.Config.AutoUpdate {
		return m.t("menu.auto_update.disabled")
	}

	if !m.AutoUpdateCancelled && m.AutoUpdateCountdown > 0 {
		return m.t("menu.auto_update.in", m.AutoUpdateCountdown)
	}

	return m.t("menu.auto_update.enabled")
}

func (m Model) renderComponentList(
	title string,
	items []string,
	titleStyle lipgloss.Style,
	itemStyle lipgloss.Style,
	versions map[string]string,
) string {
	if len(items) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(title) + "\n")

	for _, item := range items {
		line := m.componentLabel(item)
		if version, ok := versions[item]; ok && version != "" {
			line = fmt.Sprintf("%s (%s)", line, m.localizedValue(version))
		}
		b.WriteString(itemStyle.Render("• "+line) + "\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}

func (m Model) updatingPanelWidth() int {
	width := m.pageWidth() - 10
	if width < 36 {
		return 36
	}

	return width
}

func hardWrapText(text string, width int) string {
	if width <= 0 || len(text) <= width {
		return text
	}

	var b strings.Builder
	for start := 0; start < len(text); start += width {
		end := start + width
		if end > len(text) {
			end = len(text)
		}
		if start > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(text[start:end])
	}

	return b.String()
}

// renderWizard 渲染向导
func (m Model) renderWizard() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	if !m.RimeInstallStatus.Installed {
		warningBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Styles.Error).
			Padding(1, 2).
			Width(60).
			Foreground(m.Styles.Error)
		b.WriteString(warningBox.Render(m.RimeInstallStatus.Message) + "\n\n")
	}

	if m.Err != nil {
		errorMsg := m.Styles.ErrorText.Render("⚠ " + m.runtimeText("严重错误") + " ⚠ " + m.Err.Error())
		b.WriteString(errorMsg + "\n\n")
	}

	switch m.WizardStep {
	case WizardSchemeType:
		wizardTitle := RenderGradientTitle("⚡ " + m.t("wizard.title") + " ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := m.Styles.InfoBox.Render("▸ " + m.t("wizard.scheme_type"))
		b.WriteString(question + "\n\n")

		b.WriteString(m.Styles.MenuItem.Render("  [1] ► "+m.t("wizard.scheme_base")) + "\n")
		b.WriteString(m.Styles.MenuItem.Render("  [2] ► "+m.t("wizard.scheme_pro")) + "\n\n")

		b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")
		hint := m.Styles.Hint.Render(m.t("wizard.hint.1_2"))
		b.WriteString(hint)

	case WizardSchemeVariant:
		wizardTitle := RenderGradientTitle("⚡ " + m.t("wizard.title") + " ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := m.Styles.InfoBox.Render("▸ " + m.t("wizard.variant"))
		b.WriteString(question + "\n\n")

		for k, v := range types.SchemeMap {
			b.WriteString(m.Styles.MenuItem.Render(fmt.Sprintf("  [%s] ► %s", k, m.schemeLabel(v))) + "\n")
		}

		b.WriteString("\n" + m.Styles.Grid.Render(gridLine) + "\n")
		hint := m.Styles.Hint.Render(m.t("wizard.hint.1_7"))
		b.WriteString(hint)

	case WizardDownloadSource:
		wizardTitle := RenderGradientTitle("⚡ " + m.t("wizard.title") + " ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := m.Styles.InfoBox.Render("▸ " + m.t("wizard.download_source"))
		b.WriteString(question + "\n\n")

		b.WriteString(m.Styles.MenuItem.Render("  [1] ► "+m.t("wizard.source.cnb")) + "\n")
		b.WriteString(m.Styles.MenuItem.Render("  [2] ► "+m.t("wizard.source.github")) + "\n\n")

		b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")
		hint := m.Styles.Hint.Render(m.t("wizard.hint.1_2"))
		b.WriteString(hint)
	}

	return m.renderScreen(b.String())
}

// renderMenu 渲染菜单
func (m Model) renderMenu() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	if !m.RimeInstallStatus.Installed {
		warningBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Styles.Error).
			Padding(1, 2).
			Width(60).
			Foreground(m.Styles.Error)
		b.WriteString(warningBox.Render(m.RimeInstallStatus.Message) + "\n\n")
	}

	menuTitle := RenderGradientTitle("⚡ " + m.t("menu.title") + " ⚡")
	b.WriteString(menuTitle + "\n\n")

	statusContent := strings.Join([]string{
		m.renderLabeledChips([][2]string{
			{m.t("menu.summary.scheme"), m.schemeLabel(m.Cfg.GetSchemeDisplayName())},
			{m.t("menu.summary.engine"), m.Cfg.GetEngineDisplayName()},
			{m.t("menu.summary.source"), m.configuredSourceLabel()},
		}),
		m.renderLabeledChips([][2]string{
			{m.t("menu.summary.theme"), m.ThemeManager.CurrentName()},
			{m.t("menu.summary.auto_update"), m.autoUpdateStatusLabel()},
		}),
	}, "\n")
	b.WriteString(statusContent + "\n\n")

	menuItems := []struct {
		icon string
		text string
		desc string
	}{
		{
			termcolor.GetFallbackIcon("⚡", "⟳"),
			m.t("menu.auto_update.title"),
			m.t("menu.auto_update.desc"),
		},
		{
			termcolor.GetFallbackIcon("📚", "≡"),
			m.t("menu.dict_update.title"),
			m.t("menu.dict_update.desc"),
		},
		{
			termcolor.GetFallbackIcon("📦", "▢"),
			m.t("menu.scheme_update.title"),
			m.t("menu.scheme_update.desc"),
		},
		{
			termcolor.GetFallbackIcon("🤖", "◈"),
			m.t("menu.model_update.title"),
			m.t("menu.model_update.desc"),
		},
		{
			termcolor.GetFallbackIcon("⚙️", "⚙"),
			m.t("menu.config.title"),
			m.t("menu.config.desc"),
		},
		{
			termcolor.GetFallbackIcon("🎨", "◐"),
			m.t("menu.theme.title", m.ThemeManager.CurrentName()),
			m.t("menu.theme.desc"),
		},
		{
			termcolor.GetFallbackIcon("🧭", "◎"),
			m.t("menu.wizard.title"),
			m.t("menu.wizard.desc"),
		},
		{
			termcolor.GetFallbackIcon("🚪", "×"),
			m.t("menu.quit.title"),
			m.t("menu.quit.desc"),
		},
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Primary).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Muted)
	selectedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Styles.Primary).
		Background(m.Styles.Surface).
		Padding(0, 1).
		Width(m.pageWidth())

	for i, item := range menuItems {
		if i == m.MenuChoice {
			entry := titleStyle.Render(fmt.Sprintf("► [%d] %s %s", i+1, item.icon, item.text)) +
				"\n" + descStyle.Render(item.desc)
			b.WriteString(selectedStyle.Render(entry) + "\n\n")
			continue
		}

		b.WriteString(titleStyle.Render(fmt.Sprintf("  [%d] %s %s", i+1, item.icon, item.text)) + "\n")
	}

	b.WriteString("\n" + m.Styles.Grid.Render(gridLine) + "\n")

	if m.Cfg.Config.AutoUpdate && !m.AutoUpdateCancelled && m.AutoUpdateCountdown > 0 {
		countdownStyle := lipgloss.NewStyle().
			Foreground(m.Styles.Warning).
			Bold(true)
		countdownText := "⏱  " + m.t("menu.auto_update.countdown", m.AutoUpdateCountdown)
		b.WriteString(countdownStyle.Render(countdownText) + "\n\n")
	} else if m.Cfg.Config.AutoUpdate && m.AutoUpdateCancelled && m.AutoUpdateCountdown > 0 {
		cancelledStyle := lipgloss.NewStyle().
			Foreground(m.Styles.Muted)
		b.WriteString(cancelledStyle.Render("✓ "+m.t("menu.auto_update.cancelled")) + "\n\n")
	}

	hint := m.Styles.Hint.Render(m.t("menu.hint"))
	b.WriteString(hint + "\n\n")

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
	b.WriteString(statusBar)

	return m.renderScreen(b.String())
}

// renderUpdating 渲染更新中
func (m Model) renderUpdating() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚡ " + m.t("updating.title") + " ⚡")
	b.WriteString(title + "\n\n")

	component := m.CurrentComponent
	if component == "" {
		component = m.t("updating.stage.preparing")
	} else {
		component = m.componentLabel(component)
	}

	statusChip := lipgloss.NewStyle().
		Foreground(m.Styles.Background).
		Background(m.Styles.Accent).
		Padding(0, 1).
		Bold(true).
		Render(m.t("updating.stage", component))
	b.WriteString(statusChip + "\n\n")

	panelWidth := m.updatingPanelWidth()
	var progressContent strings.Builder
	progressContent.WriteString(
		m.Styles.ConfigKey.Render(m.t("updating.state")) + " " + m.Styles.ConfigValue.Render(m.ProgressMsg) + "\n",
	)

	if m.DownloadSource != "" {
		progressContent.WriteString(
			m.Styles.ConfigKey.Render(m.t("updating.source")) + " " + m.Styles.ConfigValue.Render(m.DownloadSource) + "\n",
		)
	}
	if m.DownloadFileName != "" {
		progressContent.WriteString(
			m.Styles.ConfigKey.Render(m.t("updating.file")) + " " + m.Styles.ConfigValue.Render(m.DownloadFileName) + "\n",
		)
	}
	if m.TotalSize > 0 {
		progressContent.WriteString(
			m.Styles.ConfigKey.Render(m.t("updating.progress")) + " " +
				m.Styles.ConfigValue.Render(
					fmt.Sprintf("%.2f MB / %.2f MB", float64(m.Downloaded)/1024/1024, float64(m.TotalSize)/1024/1024),
				) + "\n",
		)
	}
	if m.DownloadSpeed > 0 {
		progressContent.WriteString(
			m.Styles.ConfigKey.Render(m.t("updating.speed")) + " " +
				m.Styles.ConfigValue.Render(fmt.Sprintf("%.2f MB/s", m.DownloadSpeed)) + "\n",
		)
	}

	progressBar := m.Progress
	progressBar.Width = panelWidth
	progressContent.WriteString("\n" + progressBar.View())

	if m.DownloadURL != "" {
		urlBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Styles.Muted).
			Background(m.Styles.Background).
			Padding(0, 1).
			Width(panelWidth).
			Render(hardWrapText(m.DownloadURL, panelWidth-2))

		progressContent.WriteString(
			"\n\n" + m.Styles.ConfigKey.Render(m.t("updating.url")) + "\n" + urlBox,
		)
	}

	b.WriteString(m.renderPanel(progressContent.String(), m.Styles.Primary) + "\n\n")

	notice := lipgloss.NewStyle().
		Foreground(m.Styles.Warning).
		Render(m.t("updating.notice"))
	b.WriteString(notice + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")

	hint := m.Styles.Hint.Render(m.t("updating.hint"))
	b.WriteString(hint)

	return m.renderScreen(b.String())
}

// renderConfig 渲染配置
func (m Model) renderConfig() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚡ " + m.t("config.title") + " ⚡")
	b.WriteString(title + "\n\n")

	editableConfigs := []struct {
		key      string
		value    string
		editable bool
		index    int
	}{
		{m.t("config.field.engine"), m.Cfg.GetEngineDisplayName(), false, -1},
	}

	// 如果检测到多个引擎，显示"管理更新引擎"选项
	if len(m.Cfg.Config.InstalledEngines) > 1 {
		updateEnginesDisplay := m.t("config.value.all_engines")
		if len(m.Cfg.Config.UpdateEngines) > 0 {
			updateEnginesDisplay = strings.Join(m.Cfg.Config.UpdateEngines, "、")
		}
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"⚙ " + m.t("config.field.manage_engines"), updateEnginesDisplay, true, 0},
		)
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.scheme_type_name"), m.schemeLabel(m.Cfg.Config.SchemeType), false, -1},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.scheme_file"), m.localizedValue(m.Cfg.Config.SchemeFile), false, -1},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.dict_file"), m.localizedValue(m.Cfg.Config.DictFile), false, -1},
	)

	// 计算可编辑项的起始索引
	editIndex := 0
	if len(m.Cfg.Config.InstalledEngines) > 1 {
		editIndex = 1 // 管理更新引擎已经占用了索引 0
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.use_mirror"), fmt.Sprintf("%v", m.Cfg.Config.UseMirror), true, editIndex},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.auto_update"), fmt.Sprintf("%v", m.Cfg.Config.AutoUpdate), true, editIndex + 1},
	)

	editIndex += 2

	if m.Cfg.Config.AutoUpdate {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{m.t("config.field.auto_update_secs"), fmt.Sprintf("%d", m.Cfg.Config.AutoUpdateCountdown), true, editIndex},
		)
		editIndex++
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.language"), m.languageLabel(m.Cfg.Config.Language), true, editIndex},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.proxy_enabled"), fmt.Sprintf("%v", m.Cfg.Config.ProxyEnabled), true, editIndex + 1},
	)
	editIndex += 2

	if runtime.GOOS == "linux" {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{m.t("config.field.fcitx_compat"), fmt.Sprintf("%v", m.Cfg.Config.FcitxCompat), true, editIndex},
		)
		editIndex++

		if m.Cfg.Config.FcitxCompat {
			linkMethod := m.t("config.value.copy")
			if m.Cfg.Config.FcitxUseLink {
				linkMethod = m.t("config.value.link")
			}
			editableConfigs = append(editableConfigs,
				struct {
					key      string
					value    string
					editable bool
					index    int
				}{m.t("config.field.fcitx_use_link"), linkMethod, true, editIndex},
			)
			editIndex++
		}
	}

	if m.Cfg.Config.ProxyEnabled {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{m.t("config.field.proxy_type"), m.localizedValue(m.Cfg.Config.ProxyType), true, editIndex},
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{m.t("config.field.proxy_address"), m.localizedValue(m.Cfg.Config.ProxyAddress), true, editIndex + 1},
		)
		editIndex += 2
	}

	preHookDisplay := m.Cfg.Config.PreUpdateHook
	if preHookDisplay == "" {
		preHookDisplay = m.t("config.value.unset")
	}
	postHookDisplay := m.Cfg.Config.PostUpdateHook
	if postHookDisplay == "" {
		postHookDisplay = m.t("config.value.unset")
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.pre_hook"), preHookDisplay, true, editIndex},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{m.t("config.field.post_hook"), postHookDisplay, true, editIndex + 1},
	)
	editIndex += 2

	excludeCount := fmt.Sprintf("(%d个模式)", len(m.Cfg.Config.ExcludeFiles))
	if string(m.locale()) == "en" {
		excludeCount = fmt.Sprintf("(%d patterns)", len(m.Cfg.Config.ExcludeFiles))
	}
	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"📋 " + m.t("config.field.exclude"), excludeCount, true, editIndex},
	)
	editIndex++

	adaptiveText := m.t("config.value.disabled")
	if m.Cfg.Config.ThemeAdaptive {
		adaptiveText = m.t("config.value.enabled")
	}
	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"🎨 " + m.t("config.field.theme_adaptive"), adaptiveText, true, editIndex},
	)
	editIndex++

	if m.Cfg.Config.ThemeAdaptive {
		lightTheme := m.Cfg.Config.ThemeLight
		if lightTheme == "" {
			lightTheme = "cyberpunk-light"
		}
		darkTheme := m.Cfg.Config.ThemeDark
		if darkTheme == "" {
			darkTheme = "cyberpunk"
		}
		bg := m.ThemeManager.Background()
		bgNote := m.t("theme.current_marker")
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"  ☀️ " + m.t("config.field.theme_light"), lightTheme + func() string {
				if !bg.IsDark() {
					return bgNote
				}
				return ""
			}(), true, editIndex},
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"  🌙 " + m.t("config.field.theme_dark"), darkTheme + func() string {
				if bg.IsDark() {
					return bgNote
				}
				return ""
			}(), true, editIndex + 1},
		)
		editIndex += 2
	} else {
		fixedTheme := m.Cfg.Config.ThemeFixed
		if fixedTheme == "" {
			fixedTheme = m.ThemeManager.CurrentName()
		}
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"  🎨 " + m.t("config.field.theme_fixed"), fixedTheme, true, editIndex},
		)
		editIndex++
	}

	var configContent strings.Builder
	selectedRowStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Styles.Primary).
		Background(m.Styles.Surface).
		Padding(0, 1)

	for _, cfg := range editableConfigs {
		key := m.Styles.ConfigKey.Render(cfg.key + ":")
		value := m.Styles.ConfigValue.Render(m.localizedValue(cfg.value))

		if cfg.editable && cfg.index == m.ConfigChoice {
			line := selectedRowStyle.Render("› " + key + " " + value)
			configContent.WriteString(line + "\n")
		} else {
			line := lipgloss.NewStyle().PaddingLeft(1).Render("  " + key + " " + value)
			configContent.WriteString(line + "\n")
		}
	}

	configBox := m.Styles.InfoBox.Render(configContent.String())
	b.WriteString(configBox + "\n\n")

	pathBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Styles.Secondary).
		Background(m.Styles.Surface).
		Padding(0, 1).
		Foreground(m.Styles.Secondary)

	pathInfo := pathBox.Render(m.t("config.path", m.Cfg.ConfigPath))
	b.WriteString(pathInfo + "\n\n")

	hint1 := m.Styles.WarningText.Render("[!] " + m.t("config.help"))
	b.WriteString(hint1 + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")

	hint2 := m.Styles.Hint.Render(m.t("config.hint"))
	b.WriteString(hint2)

	return m.renderScreen(b.String())
}

// renderConfigEdit 渲染配置编辑
func (m Model) renderConfigEdit() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚡ " + m.t("config.edit.title") + " ⚡")
	b.WriteString(title + "\n\n")

	var configName string
	var inputHint string
	isBooleanField := false
	switch m.EditingKey {
	case "use_mirror":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.bool", m.t("config.edit.option.on"), m.t("config.edit.option.off"))
		isBooleanField = true
	case "auto_update":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.bool", m.t("config.edit.option.on"), m.t("config.edit.option.off"))
		isBooleanField = true
	case "auto_update_countdown":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.countdown")
	case "language":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.language")
	case "proxy_enabled":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.bool", m.t("config.edit.option.on"), m.t("config.edit.option.off"))
		isBooleanField = true
	case "fcitx_compat":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.fcitx_compat")
		isBooleanField = true
	case "fcitx_use_link":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.fcitx_link")
		isBooleanField = true
	case "proxy_type":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.proxy_type")
	case "proxy_address":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.proxy_addr")
	case "pre_update_hook":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.pre_hook")
	case "post_update_hook":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.post_hook")
	case "theme_adaptive":
		configName = m.configFieldLabel(m.EditingKey)
		inputHint = m.t("config.edit.hint.theme")
		isBooleanField = true
	}

	editBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Styles.Secondary).
		Background(m.Styles.Surface).
		Padding(1, 2).
		Width(m.contentWidth(56))

	var editContent strings.Builder
	editContent.WriteString(m.Styles.ConfigKey.Render(m.t("config.edit.item")) + " " +
		m.Styles.ConfigValue.Render(configName) + "\n\n")

	if isBooleanField {
		trueSelected := m.EditingValue == "true"
		falseSelected := m.EditingValue == "false"

		var trueOption, falseOption string
		if trueSelected {
			trueOption = m.Styles.SelectedMenuItem.Render("► [1] " + m.t("config.edit.option.on") + " (true)")
		} else {
			trueOption = m.Styles.MenuItem.Render("  [1] " + m.t("config.edit.option.on") + " (true)")
		}

		if falseSelected {
			falseOption = m.Styles.SelectedMenuItem.Render("► [2] " + m.t("config.edit.option.off") + " (false)")
		} else {
			falseOption = m.Styles.MenuItem.Render("  [2] " + m.t("config.edit.option.off") + " (false)")
		}

		editContent.WriteString(trueOption + "\n")
		editContent.WriteString(falseOption + "\n\n")
	} else if m.EditingKey == "language" {
		zhOption := "  [1] " + m.languageLabel("zh-CN")
		enOption := "  [2] " + m.languageLabel("en")
		if m.EditingValue == "zh-CN" {
			zhOption = "► [1] " + m.languageLabel("zh-CN")
		}
		if m.EditingValue == "en" {
			enOption = "► [2] " + m.languageLabel("en")
		}

		if m.EditingValue == "zh-CN" {
			editContent.WriteString(m.Styles.SelectedMenuItem.Render(zhOption) + "\n")
			editContent.WriteString(m.Styles.MenuItem.Render(enOption) + "\n\n")
		} else {
			editContent.WriteString(m.Styles.MenuItem.Render(zhOption) + "\n")
			editContent.WriteString(m.Styles.SelectedMenuItem.Render(enOption) + "\n\n")
		}
	} else {
		editContent.WriteString(m.Styles.ConfigKey.Render(m.t("config.edit.current")) + " ")
		valueWithCursor := m.EditingValue + m.Styles.Blink.Render("_")
		editContent.WriteString(m.Styles.SuccessText.Render(valueWithCursor) + "\n\n")
	}

	editContent.WriteString(m.Styles.Hint.Render(inputHint))

	editBoxRendered := editBox.Render(editContent.String())
	b.WriteString(editBoxRendered + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")

	hint := m.Styles.Hint.Render(m.t("config.edit.hint.save"))
	b.WriteString(hint)

	return m.renderScreen(b.String())
}

// renderResult 渲染更新结果
func (m Model) renderResult() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚡ " + m.t("result.title") + " ⚡")
	b.WriteString(title + "\n\n")

	borderColor := m.Styles.Error
	headlineStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Error).
		Bold(true)
	itemStyle := lipgloss.NewStyle().
		Foreground(m.Styles.Muted)
	headline := m.t("result.failure")
	if m.ResultSuccess {
		borderColor = m.Styles.Success
		headlineStyle = lipgloss.NewStyle().
			Foreground(m.Styles.Success).
			Bold(true)
		itemStyle = m.Styles.ConfigValue
		headline = m.t("result.success")
		if m.ResultSkipped {
			borderColor = m.Styles.Warning
			headlineStyle = lipgloss.NewStyle().
				Foreground(m.Styles.Warning).
				Bold(true)
			itemStyle = lipgloss.NewStyle().
				Foreground(m.Styles.Muted)
			headline = m.t("result.skipped")
		}
	}

	var resultContent strings.Builder
	resultContent.WriteString(headlineStyle.Render(headline) + "\n")
	resultContent.WriteString(itemStyle.Render(m.ResultMsg))

	if m.AutoUpdateResult != nil {
		resultContent.WriteString("\n\n")
		resultContent.WriteString(
			m.renderLabeledChipsWithWidth(m.pageWidth()-6, [][2]string{
				{
					m.t("result.updated_count"),
					m.t("result.updated_count.value", len(m.AutoUpdateResult.UpdatedComponents)),
				},
				{
					m.t("result.skipped_count"),
					m.t("result.skipped_count.value", len(m.AutoUpdateResult.SkippedComponents)),
				},
			}),
		)

		if len(m.AutoUpdateResult.UpdatedComponents) > 0 {
			resultContent.WriteString("\n\n")
			resultContent.WriteString(
				m.renderComponentList(
					m.t("result.updated_components"),
					m.AutoUpdateResult.UpdatedComponents,
					lipgloss.NewStyle().Foreground(m.Styles.Success).Bold(true),
					m.Styles.ConfigValue,
					m.AutoUpdateResult.ComponentVersions,
				),
			)
		}
		if len(m.AutoUpdateResult.SkippedComponents) > 0 {
			resultContent.WriteString("\n\n")
			resultContent.WriteString(
				m.renderComponentList(
					m.t("result.unchanged_components"),
					m.AutoUpdateResult.SkippedComponents,
					lipgloss.NewStyle().Foreground(m.Styles.Warning).Bold(true),
					lipgloss.NewStyle().Foreground(m.Styles.Muted),
					m.AutoUpdateResult.ComponentVersions,
				),
			)
		}
	}

	b.WriteString(m.renderPanel(resultContent.String(), borderColor) + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")

	hint := m.Styles.Hint.Render(m.t("result.hint"))
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(m.pageWidth()).Render(hint))

	return m.renderScreen(b.String())
}

// renderFcitxConflict 渲染 Fcitx 目录冲突对话框
func (m Model) renderFcitxConflict() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚠ " + m.t("fcitx.title") + " ⚠")
	b.WriteString(title + "\n\n")

	homeDir, _ := os.UserHomeDir()
	targetDir := filepath.Join(homeDir, ".config", "fcitx", "rime")

	question := m.Styles.WarningText.Render(m.t("fcitx.detected", targetDir))
	question += "\n\n" + m.Styles.ConfigValue.Render(m.t("fcitx.question"))

	deleteButton := m.Styles.DialogButton.Render("[1] " + m.t("fcitx.delete"))
	backupButton := m.Styles.DialogButton.Render("[2] " + m.t("fcitx.backup"))

	if m.FcitxConflictChoice == 0 {
		deleteButton = m.Styles.DialogActiveButton.Render("► [1] " + m.t("fcitx.delete"))
	} else if m.FcitxConflictChoice == 1 {
		backupButton = m.Styles.DialogActiveButton.Render("► [2] " + m.t("fcitx.backup"))
	}

	buttons := lipgloss.NewStyle().
		Width(m.contentWidth(48)).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, deleteButton, backupButton))

	checkbox := "[ ] " + m.t("fcitx.no_prompt")
	if m.FcitxConflictNoPrompt {
		checkbox = "[✓] " + m.t("fcitx.no_prompt")
	}

	checkboxRendered := m.Styles.DialogCheckbox.Render(checkbox)
	if m.FcitxConflictNoPrompt {
		checkboxRendered = m.Styles.NeonGreen.Render(checkbox)
	}
	if m.FcitxConflictChoice == 2 {
		checkboxRendered = m.Styles.DialogActiveButton.Render("► " + checkbox)
	}

	questionBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Styles.Warning).
		Background(m.Styles.Surface).
		Padding(1, 2).
		Width(m.contentWidth(46)).
		Render(question)
	ui := lipgloss.JoinVertical(
		lipgloss.Center,
		questionBox,
		"",
		buttons,
		"",
		checkboxRendered,
	)

	dialog := lipgloss.Place(m.pageWidth(), 12,
		lipgloss.Center, lipgloss.Center,
		m.Styles.DialogBox.Width(m.contentWidth(56)).Render(ui),
	)

	backdrop := lipgloss.NewStyle().
		Width(m.pageWidth()).
		Align(lipgloss.Center).
		Foreground(m.Styles.Muted).
		Render(dialog)

	b.WriteString(backdrop + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")

	hint := m.Styles.Hint.Render(m.t("fcitx.hint"))
	b.WriteString(hint)

	return m.renderScreen(b.String())
}

// renderEngineSelector 渲染引擎选择界面
func (m Model) renderEngineSelector() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚙ " + m.t("engine.title") + " ⚙")
	b.WriteString(title + "\n\n")

	info := m.Styles.InfoBox.Render(m.t("engine.help"))
	b.WriteString(info + "\n\n")

	// 显示引擎列表
	for i, engine := range m.EngineList {
		checked := " "
		if m.EngineSelections[engine] {
			checked = "✓"
		}

		cursor := "  "
		if i == m.EngineCursor {
			cursor = "► "
		}

		style := m.Styles.MenuItem
		if i == m.EngineCursor {
			style = m.Styles.SelectedMenuItem
		}

		line := fmt.Sprintf("%s[%s] %s", cursor, checked, engine)
		b.WriteString(style.Render(line) + "\n")
	}

	b.WriteString("\n" + m.Styles.Grid.Render(gridLine) + "\n")
	hint := m.Styles.Hint.Render(m.t("engine.hint"))
	b.WriteString(hint)

	return m.renderScreen(b.String())
}

// renderEnginePrompt 渲染多引擎未配置提示
func (m Model) renderEnginePrompt() string {
	var b strings.Builder

	b.WriteString(m.renderHeaderBlock())

	title := RenderGradientTitle("⚡ " + m.t("engine.prompt.title") + " ⚡")
	b.WriteString(title + "\n\n")

	separator := "、"
	if string(m.locale()) == "en" {
		separator = ", "
	}
	engineList := strings.Join(m.Cfg.Config.InstalledEngines, separator)
	message := m.t("engine.prompt.message", engineList)
	info := m.Styles.InfoBox.Render(message)
	b.WriteString(info + "\n\n")

	question := m.Styles.InfoBox.Render(m.t("engine.prompt.question"))
	b.WriteString(question + "\n\n")

	b.WriteString(m.Styles.MenuItem.Render("  [1] ► "+m.t("engine.prompt.manage")) + "\n")
	b.WriteString(m.Styles.MenuItem.Render("  [2] ► "+m.t("engine.prompt.all")) + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")
	hint := m.Styles.Hint.Render(m.t("engine.prompt.hint"))
	b.WriteString(hint)

	return m.renderScreen(b.String())
}
