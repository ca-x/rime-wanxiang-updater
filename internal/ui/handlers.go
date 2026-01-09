package ui

import (
	"fmt"
	"runtime"

	"rime-wanxiang-updater/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

// handleWizardInput 处理向导输入
func (m Model) handleWizardInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.WizardStep {
	case WizardSchemeType:
		switch msg.String() {
		case "1":
			m.Cfg.Config.SchemeType = "base"
			m.SchemeChoice = "base"
			m.WizardStep = WizardDownloadSource
			return m, nil
		case "2":
			m.Cfg.Config.SchemeType = "pro"
			m.WizardStep = WizardSchemeVariant
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case WizardSchemeVariant:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if variant, ok := types.SchemeMap[key]; ok {
			m.SchemeChoice = variant
			m.WizardStep = WizardDownloadSource
			return m, nil
		}

	case WizardDownloadSource:
		switch msg.String() {
		case "1":
			m.MirrorChoice = true
			return m.completeWizard()
		case "2":
			m.MirrorChoice = false
			return m.completeWizard()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// completeWizard 完成向导
func (m Model) completeWizard() (tea.Model, tea.Cmd) {
	m.Cfg.Config.UseMirror = m.MirrorChoice

	schemeFile, dictFile, err := m.Cfg.GetActualFilenames(m.SchemeChoice)
	if err != nil {
		m.Err = err
		return m, nil
	}

	m.Cfg.Config.SchemeFile = schemeFile
	m.Cfg.Config.DictFile = dictFile

	if err := m.Cfg.SaveConfig(); err != nil {
		m.Err = err
		return m, nil
	}

	m.WizardStep = WizardComplete
	m.State = ViewMenu
	return m, nil
}

// handleMenuInput 处理菜单输入
func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.Cfg.Config.AutoUpdate && !m.AutoUpdateCancelled && m.AutoUpdateCountdown > 0 {
			m.AutoUpdateCancelled = true
			return m, nil
		}
		return m, nil
	case "1":
		m.State = ViewUpdating
		m.ProgressMsg = "检查所有更新..."
		return m, m.runAutoUpdate()
	case "2":
		m.State = ViewUpdating
		m.ProgressMsg = "检查词库更新..."
		return m, m.runDictUpdate()
	case "3":
		m.State = ViewUpdating
		m.ProgressMsg = "检查方案更新..."
		return m, m.runSchemeUpdate()
	case "4":
		m.State = ViewUpdating
		m.ProgressMsg = "检查模型更新..."
		return m, m.runModelUpdate()
	case "5":
		m.State = ViewConfig
		return m, nil
	case "6":
		// 切换主题 - 进入主题选择
		m.InitThemeListView("theme_quick")
		m.State = ViewThemeList
		return m, nil
	case "7":
		m.State = ViewWizard
		m.WizardStep = WizardSchemeType
		return m, nil
	case "8", "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.MenuChoice > 0 {
			m.MenuChoice--
		}
	case "down", "j":
		if m.MenuChoice < 7 {
			m.MenuChoice++
		}
	case "enter":
		switch m.MenuChoice {
		case 0:
			m.State = ViewUpdating
			m.ProgressMsg = "检查所有更新..."
			return m, m.runAutoUpdate()
		case 1:
			m.State = ViewUpdating
			m.ProgressMsg = "检查词库更新..."
			return m, m.runDictUpdate()
		case 2:
			m.State = ViewUpdating
			m.ProgressMsg = "检查方案更新..."
			return m, m.runSchemeUpdate()
		case 3:
			m.State = ViewUpdating
			m.ProgressMsg = "检查模型更新..."
			return m, m.runModelUpdate()
		case 4:
			m.State = ViewConfig
			return m, nil
		case 5:
			// 切换主题
			m.InitThemeListView("theme_quick")
			m.State = ViewThemeList
			return m, nil
		case 6:
			m.State = ViewWizard
			m.WizardStep = WizardSchemeType
			return m, nil
		case 7:
			return m, tea.Quit
		}
	}
	return m, nil
}

// handleConfigInput 处理配置输入
func (m Model) handleConfigInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.State = ViewMenu
		m.ConfigChoice = 0
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.ConfigChoice > 0 {
			m.ConfigChoice--
		}
	case "down", "j":
		maxChoice := 2 // UseMirror, AutoUpdate

		if m.Cfg.Config.AutoUpdate {
			maxChoice++ // AutoUpdateCountdown
		}

		maxChoice++ // ProxyEnabled

		if runtime.GOOS == "linux" {
			maxChoice++ // FcitxCompat
			if m.Cfg.Config.FcitxCompat {
				maxChoice++ // FcitxUseLink
			}
		}

		if m.Cfg.Config.ProxyEnabled {
			maxChoice += 2 // ProxyType, ProxyAddress
		}

		maxChoice += 2 // PreUpdateHook, PostUpdateHook
		maxChoice++    // ExcludeFileManager

		// 主题配置
		maxChoice++ // ThemeAdaptive
		if m.Cfg.Config.ThemeAdaptive {
			maxChoice += 2 // ThemeLight, ThemeDark
		} else {
			maxChoice++ // ThemeFixed
		}

		if m.ConfigChoice < maxChoice {
			m.ConfigChoice++
		}
	case "enter":
		return m.startConfigEdit()
	}
	return m, nil
}

// startConfigEdit 开始编辑配置
func (m Model) startConfigEdit() (tea.Model, tea.Cmd) {
	configItems := []string{"use_mirror", "auto_update"}

	if m.Cfg.Config.AutoUpdate {
		configItems = append(configItems, "auto_update_countdown")
	}

	configItems = append(configItems, "proxy_enabled")

	if runtime.GOOS == "linux" {
		configItems = append(configItems, "fcitx_compat")
		if m.Cfg.Config.FcitxCompat {
			configItems = append(configItems, "fcitx_use_link")
		}
	}

	if m.Cfg.Config.ProxyEnabled {
		configItems = append(configItems, "proxy_type", "proxy_address")
	}

	configItems = append(configItems, "pre_update_hook", "post_update_hook")
	configItems = append(configItems, "exclude_file_manager")

	// 主题配置
	configItems = append(configItems, "theme_adaptive")
	if m.Cfg.Config.ThemeAdaptive {
		configItems = append(configItems, "theme_light", "theme_dark")
	} else {
		configItems = append(configItems, "theme_fixed")
	}

	if m.ConfigChoice < len(configItems) {
		selectedKey := configItems[m.ConfigChoice]

		if selectedKey == "exclude_file_manager" {
			m.InitExcludeView()
			m.State = ViewExcludeList
			return m, nil
		}

		// 主题选择器
		if selectedKey == "theme_dark" || selectedKey == "theme_light" || selectedKey == "theme_fixed" {
			m.InitThemeListView(selectedKey)
			m.State = ViewThemeList
			return m, nil
		}

		m.EditingKey = selectedKey

		switch m.EditingKey {
		case "use_mirror":
			if m.Cfg.Config.UseMirror {
				m.EditingValue = "true"
			} else {
				m.EditingValue = "false"
			}
		case "auto_update":
			if m.Cfg.Config.AutoUpdate {
				m.EditingValue = "true"
			} else {
				m.EditingValue = "false"
			}
		case "auto_update_countdown":
			m.EditingValue = fmt.Sprintf("%d", m.Cfg.Config.AutoUpdateCountdown)
		case "proxy_enabled":
			if m.Cfg.Config.ProxyEnabled {
				m.EditingValue = "true"
			} else {
				m.EditingValue = "false"
			}
		case "fcitx_compat":
			if m.Cfg.Config.FcitxCompat {
				m.EditingValue = "true"
			} else {
				m.EditingValue = "false"
			}
		case "fcitx_use_link":
			if m.Cfg.Config.FcitxUseLink {
				m.EditingValue = "true"
			} else {
				m.EditingValue = "false"
			}
		case "proxy_type":
			m.EditingValue = m.Cfg.Config.ProxyType
		case "proxy_address":
			m.EditingValue = m.Cfg.Config.ProxyAddress
		case "pre_update_hook":
			m.EditingValue = m.Cfg.Config.PreUpdateHook
		case "post_update_hook":
			m.EditingValue = m.Cfg.Config.PostUpdateHook
		case "theme_adaptive":
			if m.Cfg.Config.ThemeAdaptive {
				m.EditingValue = "true"
			} else {
				m.EditingValue = "false"
			}
		}

		m.State = ViewConfigEdit
	}
	return m, nil
}

// handleConfigEditInput 处理配置编辑输入
func (m Model) handleConfigEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	isBooleanField := m.EditingKey == "use_mirror" || m.EditingKey == "auto_update" || m.EditingKey == "proxy_enabled" ||
		m.EditingKey == "fcitx_compat" || m.EditingKey == "fcitx_use_link" || m.EditingKey == "theme_adaptive"

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = ViewConfig
		m.EditingKey = ""
		m.EditingValue = ""
		return m, nil
	case "enter":
		return m.saveConfigEdit()
	case "backspace":
		if !isBooleanField && len(m.EditingValue) > 0 {
			m.EditingValue = m.EditingValue[:len(m.EditingValue)-1]
		}
	default:
		if isBooleanField {
			key := msg.String()
			switch key {
			case "1":
				m.EditingValue = "true"
			case "2":
				m.EditingValue = "false"
			case "left", "right", "up", "down":
				if m.EditingValue == "true" {
					m.EditingValue = "false"
				} else {
					m.EditingValue = "true"
				}
			}
		} else {
			if len(msg.String()) == 1 {
				m.EditingValue += msg.String()
			}
		}
	}
	return m, nil
}

// saveConfigEdit 保存配置编辑
func (m Model) saveConfigEdit() (tea.Model, tea.Cmd) {
	switch m.EditingKey {
	case "use_mirror":
		m.Cfg.Config.UseMirror = m.EditingValue == "true"
	case "auto_update":
		oldValue := m.Cfg.Config.AutoUpdate
		m.Cfg.Config.AutoUpdate = m.EditingValue == "true"
		if oldValue && !m.Cfg.Config.AutoUpdate {
			m.AutoUpdateCancelled = false
		}
	case "auto_update_countdown":
		var countdown int
		if _, err := fmt.Sscanf(m.EditingValue, "%d", &countdown); err == nil {
			if countdown < 1 {
				countdown = 1
			} else if countdown > 60 {
				countdown = 60
			}
			m.Cfg.Config.AutoUpdateCountdown = countdown
			m.AutoUpdateCountdown = countdown
		}
	case "proxy_enabled":
		m.Cfg.Config.ProxyEnabled = m.EditingValue == "true"
	case "fcitx_compat":
		oldValue := m.Cfg.Config.FcitxCompat
		newValue := m.EditingValue == "true"
		m.Cfg.Config.FcitxCompat = newValue

		if newValue != oldValue {
			if newValue {
				needsPrompt, conflictExists, err := m.Cfg.SyncToFcitxDir()
				if err != nil {
					m.Err = err
				} else if needsPrompt && conflictExists {
					m.FcitxConflictChoice = 0
					m.FcitxConflictNoPrompt = false
					m.FcitxConflictCallback = func() {
						if err := m.Cfg.ResolveFcitxConflict(); err != nil {
							m.Err = err
						}
					}
					if err := m.Cfg.SaveConfig(); err != nil {
						m.Err = err
					}
					m.State = ViewFcitxConflict
					m.EditingKey = ""
					m.EditingValue = ""
					return m, nil
				}
			} else {
				m.ConfigChoice = 3
			}
		}
	case "fcitx_use_link":
		m.Cfg.Config.FcitxUseLink = m.EditingValue == "true"
		if m.Cfg.Config.FcitxCompat {
			needsPrompt, conflictExists, err := m.Cfg.SyncToFcitxDir()
			if err != nil {
				m.Err = err
			} else if needsPrompt && conflictExists {
				m.FcitxConflictChoice = 0
				m.FcitxConflictNoPrompt = false
				m.FcitxConflictCallback = func() {
					if err := m.Cfg.ResolveFcitxConflict(); err != nil {
						m.Err = err
					}
				}
				if err := m.Cfg.SaveConfig(); err != nil {
					m.Err = err
				}
				m.State = ViewFcitxConflict
				m.EditingKey = ""
				m.EditingValue = ""
				return m, nil
			}
		}
	case "proxy_type":
		m.Cfg.Config.ProxyType = m.EditingValue
	case "proxy_address":
		m.Cfg.Config.ProxyAddress = m.EditingValue
	case "pre_update_hook":
		m.Cfg.Config.PreUpdateHook = m.EditingValue
	case "post_update_hook":
		m.Cfg.Config.PostUpdateHook = m.EditingValue
	case "theme_adaptive":
		m.Cfg.Config.ThemeAdaptive = m.EditingValue == "true"
		// 更新主题管理器
		if m.Cfg.Config.ThemeAdaptive {
			light := m.Cfg.Config.ThemeLight
			dark := m.Cfg.Config.ThemeDark
			if light == "" {
				light = "cyberpunk-light"
			}
			if dark == "" {
				dark = "cyberpunk"
			}
			m.ThemeManager.SetAdaptiveTheme(light, dark)
		} else if m.Cfg.Config.ThemeFixed != "" {
			m.ThemeManager.SetTheme(m.Cfg.Config.ThemeFixed)
		}
		// 刷新样式
		m.Styles = DefaultStyles(m.ThemeManager)
	}

	if err := m.Cfg.SaveConfig(); err != nil {
		m.Err = err
		m.State = ViewConfig
		return m, nil
	}

	m.State = ViewConfig
	m.EditingKey = ""
	m.EditingValue = ""
	return m, nil
}

// handleResultInput 处理结果页面输入
func (m Model) handleResultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	m.State = ViewMenu
	return m, nil
}

// handleFcitxConflictInput 处理 Fcitx 冲突对话框输入
func (m Model) handleFcitxConflictInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = ViewConfig
		return m, nil
	case "up", "left", "k":
		if m.FcitxConflictChoice > 0 {
			m.FcitxConflictChoice--
		}
	case "down", "right", "j":
		if m.FcitxConflictChoice < 2 {
			m.FcitxConflictChoice++
		}
	case "1":
		m.FcitxConflictChoice = 0
	case "2":
		m.FcitxConflictChoice = 1
	case " ":
		if m.FcitxConflictChoice == 2 {
			m.FcitxConflictNoPrompt = !m.FcitxConflictNoPrompt
		}
	case "enter":
		if m.FcitxConflictChoice == 2 {
			m.FcitxConflictNoPrompt = !m.FcitxConflictNoPrompt
		} else {
			return m.applyFcitxConflictChoice()
		}
	}
	return m, nil
}

// applyFcitxConflictChoice 应用 Fcitx 冲突选择
func (m Model) applyFcitxConflictChoice() (tea.Model, tea.Cmd) {
	if m.FcitxConflictChoice == 0 {
		m.Cfg.Config.FcitxConflictAction = "delete"
	} else {
		m.Cfg.Config.FcitxConflictAction = "backup"
	}

	if m.FcitxConflictNoPrompt {
		m.Cfg.Config.FcitxConflictPrompt = false
		if err := m.Cfg.SaveConfig(); err != nil {
			m.Err = err
		}
	}

	if m.FcitxConflictCallback != nil {
		m.FcitxConflictCallback()
	}

	m.State = ViewConfig
	return m, nil
}


