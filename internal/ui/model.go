package ui

import (
	"fmt"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// NewModel 创建新模型
func NewModel(cfg *config.Manager) Model {
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 60

	state := ViewMenu
	wizardStep := WizardSchemeType
	if cfg.Config.SchemeType == "" || cfg.Config.SchemeFile == "" || cfg.Config.DictFile == "" {
		state = ViewWizard
	}

	rimeStatus := detector.CheckRimeInstallation()

	countdown := cfg.Config.AutoUpdateCountdown
	if countdown <= 0 {
		countdown = 5
	}

	// 初始化主题管理器
	themeMgr := theme.NewManager()

	// 从配置加载主题设置
	if cfg.Config.ThemeAdaptive {
		light := cfg.Config.ThemeLight
		dark := cfg.Config.ThemeDark
		if light == "" {
			light = "cyberpunk-light"
		}
		if dark == "" {
			dark = "cyberpunk"
		}
		themeMgr.SetAdaptiveTheme(light, dark)
	} else if cfg.Config.ThemeFixed != "" {
		themeMgr.SetTheme(cfg.Config.ThemeFixed)
	}

	// 创建主题化样式
	styles := DefaultStyles(themeMgr)

	return Model{
		Cfg:                 cfg,
		ThemeManager:        themeMgr,
		Styles:              styles,
		State:               state,
		WizardStep:          wizardStep,
		Progress:            p,
		RimeInstallStatus:   rimeStatus,
		AutoUpdateCountdown: countdown,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	if m.State == ViewMenu && m.Cfg.Config.AutoUpdate {
		return countdownTick()
	}
	return nil
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Progress.Width = msg.Width - 4
		return m, nil

	case tea.KeyMsg:
		switch m.State {
		case ViewWizard:
			return m.handleWizardInput(msg)
		case ViewMenu:
			return m.handleMenuInput(msg)
		case ViewConfig:
			return m.handleConfigInput(msg)
		case ViewConfigEdit:
			return m.handleConfigEditInput(msg)
		case ViewExcludeList:
			return m.handleExcludeListInput(msg)
		case ViewExcludeEdit:
			return m.handleExcludeEditInput(msg)
		case ViewExcludeAdd:
			return m.handleExcludeAddInput(msg)
		case ViewFcitxConflict:
			return m.handleFcitxConflictInput(msg)
		case ViewThemeList:
			return m.handleThemeListInput(msg)
		case ViewResult:
			return m.handleResultInput(msg)
		case ViewUpdating:
			switch msg.String() {
			case "q", "esc":
				m.State = ViewMenu
				m.Updating = false
				m.ProgressChan = nil
				m.CompletionChan = nil
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

	case UpdateMsg:
		m.ProgressMsg = msg.Message

		if msg.DownloadMode {
			m.IsDownloading = true
			m.DownloadSource = msg.Source
			m.DownloadFileName = msg.FileName
			m.Downloaded = msg.Downloaded
			m.TotalSize = msg.Total
			m.DownloadSpeed = msg.Speed
		} else {
			m.IsDownloading = false
		}

		cmd := m.Progress.SetPercent(msg.Percent)
		if m.ProgressChan != nil && m.CompletionChan != nil {
			return m, tea.Batch(cmd, listenForProgress(m.ProgressChan, m.CompletionChan))
		}
		return m, cmd

	case UpdateCompleteMsg:
		m.Updating = false
		m.State = ViewResult

		m.ProgressChan = nil
		m.CompletionChan = nil

		m.ResultSkipped = msg.Skipped
		m.AutoUpdateResult = msg.AutoUpdateDetails

		if msg.Err != nil {
			m.ResultSuccess = false
			m.ResultMsg = fmt.Sprintf("%s更新失败: %v", msg.UpdateType, msg.Err)
		} else if msg.Skipped {
			m.ResultSuccess = true
			if msg.StatusMessage != "" {
				m.ResultMsg = fmt.Sprintf("%s%s", msg.UpdateType, msg.StatusMessage)
			} else {
				m.ResultMsg = fmt.Sprintf("%s已是最新版本，无需更新", msg.UpdateType)
			}
		} else {
			m.ResultSuccess = true
			m.ResultMsg = fmt.Sprintf("%s更新完成！", msg.UpdateType)
		}
		return m, nil

	case CountdownTickMsg:
		if m.State == ViewMenu && m.Cfg.Config.AutoUpdate && !m.AutoUpdateCancelled {
			m.AutoUpdateCountdown--
			if m.AutoUpdateCountdown <= 0 {
				m.State = ViewUpdating
				m.ProgressMsg = "检查所有更新..."
				m.AutoUpdateCancelled = true
				return m, m.runAutoUpdate()
			}
			return m, countdownTick()
		}
		return m, nil
	}

	return m, nil
}

// View 渲染视图
func (m Model) View() string {
	switch m.State {
	case ViewWizard:
		return m.renderWizard()
	case ViewMenu:
		return m.renderMenu()
	case ViewUpdating:
		return m.renderUpdating()
	case ViewConfig:
		return m.renderConfig()
	case ViewConfigEdit:
		return m.renderConfigEdit()
	case ViewExcludeList:
		return m.renderExcludeList()
	case ViewExcludeEdit:
		return m.renderExcludeEdit()
	case ViewExcludeAdd:
		return m.renderExcludeAdd()
	case ViewFcitxConflict:
		return m.renderFcitxConflict()
	case ViewThemeList:
		return m.renderThemeList()
	case ViewResult:
		return m.renderResult()
	}
	return ""
}
