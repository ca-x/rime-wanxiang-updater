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

// renderWizard 渲染向导
func (m Model) renderWizard() string {
	var b strings.Builder

	logo := m.Styles.Logo.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(m.Styles.ScanLine.Render(scanLine) + "\n\n")

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
		errorMsg := m.Styles.ErrorText.Render("⚠ 严重错误 ⚠ " + m.Err.Error())
		b.WriteString(errorMsg + "\n\n")
	}

	switch m.WizardStep {
	case WizardSchemeType:
		wizardTitle := RenderGradientTitle("⚡ 初始化向导 ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := m.Styles.InfoBox.Render("▸ 选择方案版本:")
		b.WriteString(question + "\n\n")

		b.WriteString(m.Styles.MenuItem.Render("  [1] ► 万象基础版") + "\n")
		b.WriteString(m.Styles.MenuItem.Render("  [2] ► 万象增强版（支持辅助码）") + "\n\n")

		b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")
		hint := m.Styles.Hint.Render("[>] Input: 1-2 | [Q] Quit")
		b.WriteString(hint)

	case WizardSchemeVariant:
		wizardTitle := RenderGradientTitle("⚡ 初始化向导 ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := m.Styles.InfoBox.Render("▸ 选择辅助码方案:")
		b.WriteString(question + "\n\n")

		for k, v := range types.SchemeMap {
			b.WriteString(m.Styles.MenuItem.Render(fmt.Sprintf("  [%s] ► %s", k, v)) + "\n")
		}

		b.WriteString("\n" + m.Styles.Grid.Render(gridLine) + "\n")
		hint := m.Styles.Hint.Render("[>] Input: 1-7 | [Q] Quit")
		b.WriteString(hint)

	case WizardDownloadSource:
		wizardTitle := RenderGradientTitle("⚡ 初始化向导 ⚡")
		b.WriteString(wizardTitle + "\n\n")

		question := m.Styles.InfoBox.Render("▸ 选择下载源:")
		b.WriteString(question + "\n\n")

		b.WriteString(m.Styles.MenuItem.Render("  [1] ► CNB 镜像（推荐，国内访问更快）") + "\n")
		b.WriteString(m.Styles.MenuItem.Render("  [2] ► GitHub 官方源") + "\n\n")

		b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")
		hint := m.Styles.Hint.Render("[>] Input: 1-2 | [Q] Quit")
		b.WriteString(hint)
	}

	return m.Styles.Container.Render(b.String())
}

// renderMenu 渲染菜单
func (m Model) renderMenu() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(m.Styles.ScanLine.Render(scanLine) + "\n\n")

	if !m.RimeInstallStatus.Installed {
		warningBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Styles.Error).
			Padding(1, 2).
			Width(60).
			Foreground(m.Styles.Error)
		b.WriteString(warningBox.Render(m.RimeInstallStatus.Message) + "\n\n")
	}

	menuTitle := RenderGradientTitle("⚡ 主控制面板 ⚡")
	b.WriteString(menuTitle + "\n\n")

	menuItems := []struct {
		icon string
		text string
	}{
		{termcolor.GetFallbackIcon("⚡", "⟳"), "自动更新"},                                        // ⚡ → ⟳ (循环箭头)
		{termcolor.GetFallbackIcon("📚", "≡"), "词库更新"},                                        // 📚 → ≡ (三横线，像书页)
		{termcolor.GetFallbackIcon("📦", "▢"), "方案更新"},                                        // 📦 → ▢ (空心方块)
		{termcolor.GetFallbackIcon("🤖", "◈"), "模型更新"},                                        // 🤖 → ◈ (菱形)
		{termcolor.GetFallbackIcon("⚙️", "⚙"), "查看配置"},                                       // ⚙️ → ⚙ (齿轮符号)
		{termcolor.GetFallbackIcon("🎨", "◐"), "切换主题 (" + m.ThemeManager.CurrentName() + ")"}, // 🎨 → ◐ (半圆)
		{termcolor.GetFallbackIcon("🧭", "◎"), "设置向导"},                                        // 🧭 → ◎ (双圆)
		{termcolor.GetFallbackIcon("🚪", "×"), "退出程序"},                                        // 🚪 → × (叉号)
	}

	for i, item := range menuItems {
		itemText := fmt.Sprintf(" %s  [%d] %s", item.icon, i+1, item.text)
		if i == m.MenuChoice {
			b.WriteString(m.Styles.SelectedMenuItem.Render("►"+itemText) + "\n")
		} else {
			b.WriteString(m.Styles.MenuItem.Render(" "+itemText) + "\n")
		}
	}

	b.WriteString("\n" + m.Styles.Grid.Render(gridLine) + "\n")

	if m.Cfg.Config.AutoUpdate && !m.AutoUpdateCancelled && m.AutoUpdateCountdown > 0 {
		countdownStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)
		countdownText := fmt.Sprintf("⏱  自动更新将在 %d 秒后开始... (按 ESC 取消)", m.AutoUpdateCountdown)
		b.WriteString(countdownStyle.Render(countdownText) + "\n\n")
	} else if m.Cfg.Config.AutoUpdate && m.AutoUpdateCancelled && m.AutoUpdateCountdown > 0 {
		cancelledStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
		b.WriteString(cancelledStyle.Render("✓ 已取消自动更新") + "\n\n")
	}

	hint := m.Styles.Hint.Render("[>] Input: 1-8 | Navigate: J/K or Arrow Keys | [Q] Quit")
	b.WriteString(hint + "\n\n")

	statusBar := RenderStatusBarThemed(
		m.Styles,
		version.GetVersion(),
		m.Cfg.GetEngineDisplayName(),
		func() string {
			if m.Cfg.Config.UseMirror {
				return "CNB镜像"
			}
			return "GitHub"
		}(),
		m.Cfg.GetSchemeDisplayName(),
	)
	b.WriteString(statusBar)

	return m.Styles.Container.Render(b.String())
}

// renderUpdating 渲染更新中
func (m Model) renderUpdating() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	bootSeq := RenderBootSequence(version.GetVersion())
	b.WriteString(bootSeq + "\n")

	status := statusProcessingStyle.Render("⬢ 处理中 ⬢")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Render(status) + "\n\n")

	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚡ 正在更新 ⚡")
	b.WriteString(title + "\n\n")

	msgBox := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(neonGreen).
		Padding(1, 2).
		Width(60)

	var msgContent strings.Builder

	if m.IsDownloading {
		if m.DownloadSource != "" && m.DownloadFileName != "" {
			msgContent.WriteString(configKeyStyle.Render("▸ ") +
				configValueStyle.Render(m.DownloadSource) +
				configKeyStyle.Render(" > ") +
				configValueStyle.Render(m.DownloadFileName) + "\n\n")
		}

		if m.TotalSize > 0 {
			downloadedMB := float64(m.Downloaded) / 1024 / 1024
			totalMB := float64(m.TotalSize) / 1024 / 1024

			progressLine := successStyle.Render(fmt.Sprintf("%.2f MB / %.2f MB", downloadedMB, totalMB))
			if m.DownloadSpeed > 0 {
				progressLine += configKeyStyle.Render("  |  ") +
					neonGreenStyle.Render(fmt.Sprintf("%.2f MB/s", m.DownloadSpeed))
			}
			msgContent.WriteString(progressLine)
		} else {
			msgContent.WriteString(progressMsgStyle.Render("▸ " + m.ProgressMsg))
		}
	} else {
		msgContent.WriteString(progressMsgStyle.Render("▸ " + m.ProgressMsg))
	}

	b.WriteString(msgBox.Render(msgContent.String()) + "\n\n")

	if m.IsDownloading && m.TotalSize > 0 {
		progressBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(neonCyan).
			Padding(0, 1)

		percent := float64(m.Downloaded) / float64(m.TotalSize)
		progressBar := progressBox.Render(m.Progress.ViewAs(percent))
		b.WriteString(progressBar + "\n\n")
	}

	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	hint := hintStyle.Render("[...] Please wait... System is updating... | [Q]/[ESC] Cancel | [Ctrl+C] Quit")
	b.WriteString(hint)

	return containerStyle.Render(b.String())
}

// renderConfig 渲染配置
func (m Model) renderConfig() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(m.Styles.ScanLine.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚡ 系统配置 ⚡")
	b.WriteString(title + "\n\n")

	editableConfigs := []struct {
		key      string
		value    string
		editable bool
		index    int
	}{
		{"引擎", m.Cfg.GetEngineDisplayName(), false, -1},
	}

	// 如果检测到多个引擎，显示"管理更新引擎"选项
	if len(m.Cfg.Config.InstalledEngines) > 1 {
		updateEnginesDisplay := "全部引擎"
		if len(m.Cfg.Config.UpdateEngines) > 0 {
			updateEnginesDisplay = strings.Join(m.Cfg.Config.UpdateEngines, "、")
		}
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"⚙ 管理更新引擎", updateEnginesDisplay, true, 0},
		)
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"方案类型", m.Cfg.Config.SchemeType, false, -1},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"方案文件", m.Cfg.Config.SchemeFile, false, -1},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"词库文件", m.Cfg.Config.DictFile, false, -1},
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
		}{"使用镜像", fmt.Sprintf("%v", m.Cfg.Config.UseMirror), true, editIndex},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"自动更新", fmt.Sprintf("%v", m.Cfg.Config.AutoUpdate), true, editIndex + 1},
	)

	editIndex += 2

	if m.Cfg.Config.AutoUpdate {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"自动更新倒计时(秒)", fmt.Sprintf("%d", m.Cfg.Config.AutoUpdateCountdown), true, editIndex},
		)
		editIndex++
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"代理启用", fmt.Sprintf("%v", m.Cfg.Config.ProxyEnabled), true, editIndex},
	)
	editIndex++

	if runtime.GOOS == "linux" {
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"Fcitx兼容(同步到~/.config/fcitx/rime)", fmt.Sprintf("%v", m.Cfg.Config.FcitxCompat), true, editIndex},
		)
		editIndex++

		if m.Cfg.Config.FcitxCompat {
			linkMethod := "复制文件"
			if m.Cfg.Config.FcitxUseLink {
				linkMethod = "软链接"
			}
			editableConfigs = append(editableConfigs,
				struct {
					key      string
					value    string
					editable bool
					index    int
				}{"同步方式", linkMethod, true, editIndex},
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
			}{"代理类型", m.Cfg.Config.ProxyType, true, editIndex},
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"代理地址", m.Cfg.Config.ProxyAddress, true, editIndex + 1},
		)
		editIndex += 2
	}

	preHookDisplay := m.Cfg.Config.PreUpdateHook
	if preHookDisplay == "" {
		preHookDisplay = "(未设置)"
	}
	postHookDisplay := m.Cfg.Config.PostUpdateHook
	if postHookDisplay == "" {
		postHookDisplay = "(未设置)"
	}

	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"更新前Hook", preHookDisplay, true, editIndex},
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"更新后Hook", postHookDisplay, true, editIndex + 1},
	)
	editIndex += 2

	excludeCount := fmt.Sprintf("(%d个模式)", len(m.Cfg.Config.ExcludeFiles))
	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"📋 管理排除文件", excludeCount, true, editIndex},
	)
	editIndex++

	// 主题配置
	adaptiveText := "禁用"
	if m.Cfg.Config.ThemeAdaptive {
		adaptiveText = "启用"
	}
	editableConfigs = append(editableConfigs,
		struct {
			key      string
			value    string
			editable bool
			index    int
		}{"🎨 自适应主题", adaptiveText, true, editIndex},
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
		// 显示检测到的背景
		bg := m.ThemeManager.Background()
		bgNote := ""
		if bg.IsDark() {
			bgNote = " (当前使用↓)"
		} else {
			bgNote = " (当前使用↓)"
		}
		editableConfigs = append(editableConfigs,
			struct {
				key      string
				value    string
				editable bool
				index    int
			}{"  ☀️ 浅色主题", lightTheme + func() string {
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
			}{"  🌙 深色主题", darkTheme + func() string {
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
			}{"  🎨 固定主题", fixedTheme, true, editIndex},
		)
		editIndex++
	}

	var configContent strings.Builder
	for _, cfg := range editableConfigs {
		key := m.Styles.ConfigKey.Render(cfg.key + ":")
		value := m.Styles.ConfigValue.Render(cfg.value)

		if cfg.editable && cfg.index == m.ConfigChoice {
			line := m.Styles.SelectedMenuItem.Render("►") + "  ▸ " + key + " " + value
			configContent.WriteString(line + "\n")
		} else {
			line := " " + "  ▸ " + key + " " + value
			configContent.WriteString(line + "\n")
		}
	}

	configBox := m.Styles.InfoBox.Render(configContent.String())
	b.WriteString(configBox + "\n\n")

	pathBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.Styles.Secondary).
		Padding(0, 1).
		Foreground(m.Styles.Secondary)

	pathInfo := pathBox.Render("配置路径: " + m.Cfg.ConfigPath)
	b.WriteString(pathInfo + "\n\n")

	hint1 := m.Styles.WarningText.Render("[!] Use Arrow Keys to select, Enter to edit")
	b.WriteString(hint1 + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")

	hint2 := m.Styles.Hint.Render("[>] Navigate: J/K or Arrow Keys | [Enter] Edit | [Q]/[ESC] Back")
	b.WriteString(hint2)

	return m.Styles.Container.Render(b.String())
}

// renderConfigEdit 渲染配置编辑
func (m Model) renderConfigEdit() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(m.Styles.ScanLine.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚡ 编辑配置 ⚡")
	b.WriteString(title + "\n\n")

	var configName string
	var inputHint string
	isBooleanField := false
	switch m.EditingKey {
	case "use_mirror":
		configName = "使用镜像"
		inputHint = "Select: [1] Enable  [2] Disable | Arrow keys to toggle"
		isBooleanField = true
	case "auto_update":
		configName = "自动更新"
		inputHint = "Select: [1] Enable  [2] Disable | Arrow keys to toggle"
		isBooleanField = true
	case "auto_update_countdown":
		configName = "自动更新倒计时(秒)"
		inputHint = "输入倒计时秒数 (1-60秒)"
	case "proxy_enabled":
		configName = "代理启用"
		inputHint = "Select: [1] Enable  [2] Disable | Arrow keys to toggle"
		isBooleanField = true
	case "fcitx_compat":
		configName = "Fcitx兼容"
		inputHint = "启用后将同步配置到 ~/.config/fcitx/rime/ 以兼容外部插件 | [1] Enable  [2] Disable"
		isBooleanField = true
	case "fcitx_use_link":
		configName = "同步方式"
		inputHint = "[1] 软链接(推荐,自动同步,节省空间)  [2] 复制文件(独立,更安全)"
		isBooleanField = true
	case "proxy_type":
		configName = "代理类型"
		inputHint = "Input proxy type: http/https/socks5"
	case "proxy_address":
		configName = "代理地址"
		inputHint = "Input proxy address (e.g. 127.0.0.1:7890)"
	case "pre_update_hook":
		configName = "更新前Hook"
		inputHint = "脚本路径(如~/backup.sh),更新前执行,失败将取消更新"
	case "post_update_hook":
		configName = "更新后Hook"
		inputHint = "脚本路径(如~/notify.sh),更新后执行,失败不影响更新结果"
	case "theme_adaptive":
		configName = "自适应主题"
		inputHint = "启用后根据终端明暗自动切换主题 | [1] Enable  [2] Disable"
		isBooleanField = true
	}

	editBox := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(m.Styles.Secondary).
		Padding(1, 2).
		Width(60)

	var editContent strings.Builder
	editContent.WriteString(m.Styles.ConfigKey.Render("配置项: ") + m.Styles.ConfigValue.Render(configName) + "\n\n")

	if isBooleanField {
		trueSelected := m.EditingValue == "true"
		falseSelected := m.EditingValue == "false"

		var trueOption, falseOption string
		if trueSelected {
			trueOption = m.Styles.SelectedMenuItem.Render("► [1] Enable (true)")
		} else {
			trueOption = m.Styles.MenuItem.Render("  [1] Enable (true)")
		}

		if falseSelected {
			falseOption = m.Styles.SelectedMenuItem.Render("► [2] Disable (false)")
		} else {
			falseOption = m.Styles.MenuItem.Render("  [2] Disable (false)")
		}

		editContent.WriteString(trueOption + "\n")
		editContent.WriteString(falseOption + "\n\n")
	} else {
		editContent.WriteString(m.Styles.ConfigKey.Render("当前值: "))
		valueWithCursor := m.EditingValue + m.Styles.Blink.Render("_")
		editContent.WriteString(m.Styles.SuccessText.Render(valueWithCursor) + "\n\n")
	}

	editContent.WriteString(m.Styles.Hint.Render(inputHint))

	editBoxRendered := editBox.Render(editContent.String())
	b.WriteString(editBoxRendered + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n\n")

	hint := m.Styles.Hint.Render("[>] [Enter] Save | [ESC] Cancel | [Backspace] Delete")
	b.WriteString(hint)

	return m.Styles.Container.Render(b.String())
}

// renderResult 渲染更新结果
func (m Model) renderResult() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚡ 更新结果 ⚡")
	b.WriteString(title + "\n\n")

	var resultBox lipgloss.Style
	var icon string

	if m.ResultSuccess {
		resultBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(neonGreen).
			Padding(2, 3).
			Width(60)
		icon = "✓"
	} else {
		resultBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(glitchRed).
			Padding(2, 3).
			Width(60)
		icon = "✗"
	}

	var msgContent strings.Builder
	if m.ResultSuccess {
		msgContent.WriteString(successStyle.Render(fmt.Sprintf("%s %s", icon, m.ResultMsg)))

		if m.AutoUpdateResult != nil {
			msgContent.WriteString("\n\n")

			if len(m.AutoUpdateResult.UpdatedComponents) > 0 {
				msgContent.WriteString(RenderCheckList("Updated", m.AutoUpdateResult.UpdatedComponents, true, m.AutoUpdateResult.ComponentVersions))
			}

			if len(m.AutoUpdateResult.SkippedComponents) > 0 {
				if len(m.AutoUpdateResult.UpdatedComponents) > 0 {
					msgContent.WriteString("\n")
				}
				msgContent.WriteString(RenderCheckList("Up-to-date", m.AutoUpdateResult.SkippedComponents, false, m.AutoUpdateResult.ComponentVersions))
			}
		}

		if !m.ResultSkipped && m.AutoUpdateResult != nil && len(m.AutoUpdateResult.UpdatedComponents) > 0 {
			msgContent.WriteString("\n")
			msgContent.WriteString(configValueStyle.Render("System update completed | 更新已成功应用到系统"))
		}
	} else {
		msgContent.WriteString(errorStyle.Render(fmt.Sprintf("%s %s", icon, m.ResultMsg)))
		msgContent.WriteString("\n\n")
		msgContent.WriteString(configValueStyle.Render("Please check error and retry | 请检查错误信息并重试"))
	}

	resultMessage := resultBox.Render(msgContent.String())
	b.WriteString(resultMessage + "\n\n")

	b.WriteString(gridStyle.Render(gridLine) + "\n\n")

	hint := blinkStyle.Render("[>] Press any key to return to main menu...")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(65).Render(hint))

	return containerStyle.Render(b.String())
}

// renderFcitxConflict 渲染 Fcitx 目录冲突对话框
func (m Model) renderFcitxConflict() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(scanLineStyle.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚠ Fcitx 目录冲突 ⚠")
	b.WriteString(title + "\n\n")

	homeDir, _ := os.UserHomeDir()
	targetDir := filepath.Join(homeDir, ".config", "fcitx", "rime")

	question := warningStyle.Render(fmt.Sprintf("检测到目录已存在: %s", targetDir))
	question += "\n\n" + configValueStyle.Render("请选择如何处理:")

	deleteButton := dialogButtonStyle.Render("[1] 直接删除")
	backupButton := dialogButtonStyle.Render("[2] 备份后删除")

	if m.FcitxConflictChoice == 0 {
		deleteButton = dialogActiveButtonStyle.Render("► [1] 直接删除")
	} else if m.FcitxConflictChoice == 1 {
		backupButton = dialogActiveButtonStyle.Render("► [2] 备份后删除")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, deleteButton, backupButton)

	checkbox := "[ ] 不再提示，记住我的选择"
	if m.FcitxConflictNoPrompt {
		checkbox = "[✓] 不再提示，记住我的选择"
	}

	checkboxRendered := dialogCheckboxStyle.Render(checkbox)
	if m.FcitxConflictNoPrompt {
		checkboxRendered = dialogCheckboxCheckedStyle.Render(checkbox)
	}
	if m.FcitxConflictChoice == 2 {
		checkboxRendered = dialogActiveButtonStyle.Render("► " + checkbox)
	}

	ui := lipgloss.JoinVertical(lipgloss.Left, question, buttons, checkboxRendered)

	dialog := lipgloss.Place(65, 12,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	b.WriteString(dialog + "\n\n")

	b.WriteString(gridStyle.Render(gridLine) + "\n\n")

	hint := hintStyle.Render("[>] Navigate: 1-2 or Arrow Keys | [Space/Enter] Toggle/Confirm | [ESC] Cancel")
	b.WriteString(hint)

	return containerStyle.Render(b.String())
}

// renderEngineSelector 渲染引擎选择界面
func (m Model) renderEngineSelector() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(m.Styles.ScanLine.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚙ 选择要更新的引擎 ⚙")
	b.WriteString(title + "\n\n")

	info := m.Styles.InfoBox.Render("使用 空格 或 回车 切换选择，按 S 保存")
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
	hint := m.Styles.Hint.Render("[Space/Enter] Toggle | [S] Save | [Q/ESC] Cancel")
	b.WriteString(hint)

	return m.Styles.Container.Render(b.String())
}

// renderEnginePrompt 渲染多引擎未配置提示
func (m Model) renderEnginePrompt() string {
	var b strings.Builder

	logo := logoStyle.Render(asciiLogo)
	b.WriteString(logo + "\n")

	header := RenderHeader(version.GetVersion())
	b.WriteString(header + "\n")

	b.WriteString(m.Styles.ScanLine.Render(scanLine) + "\n\n")

	title := RenderGradientTitle("⚡ 多引擎检测 ⚡")
	b.WriteString(title + "\n\n")

	// 显示检测到的引擎
	engineList := strings.Join(m.Cfg.Config.InstalledEngines, "、")
	message := fmt.Sprintf("检测到您安装了多个输入法引擎：%s", engineList)
	info := m.Styles.InfoBox.Render(message)
	b.WriteString(info + "\n\n")

	question := m.Styles.InfoBox.Render("您希望如何处理更新？")
	b.WriteString(question + "\n\n")

	b.WriteString(m.Styles.MenuItem.Render("  [1] ► 进入设置选择要更新的引擎") + "\n")
	b.WriteString(m.Styles.MenuItem.Render("  [2] ► 更新所有已安装的引擎") + "\n\n")

	b.WriteString(m.Styles.Grid.Render(gridLine) + "\n")
	hint := m.Styles.Hint.Render("[>] Input: 1-2 | [Q/ESC] Cancel")
	b.WriteString(hint)

	return m.Styles.Container.Render(b.String())
}
