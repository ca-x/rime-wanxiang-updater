package ui

import (
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/controller"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// NewModel 创建新模型
func NewModel(
	cfg *config.Manager,
	themeMgr *theme.Manager,
	commandChan chan<- controller.Command,
	eventChan <-chan controller.Event,
) Model {
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

	// 创建主题化样式
	styles := DefaultStyles(themeMgr)

	return Model{
		CommandChan:         commandChan,
		EventChan:           eventChan,
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
	// 启动事件监听和倒计时（如果需要）
	cmds := []tea.Cmd{listenForEvents(m.EventChan)}

	if m.State == ViewMenu && m.Cfg.Config.AutoUpdate {
		cmds = append(cmds, countdownTick())
	}

	return tea.Batch(cmds...)
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Progress.Width = msg.Width - 4
		return m, nil

	case controller.Event:
		return m.handleControllerEvent(msg)

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
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

	case CountdownTickMsg:
		if m.State == ViewMenu && m.Cfg.Config.AutoUpdate && !m.AutoUpdateCancelled {
			m.AutoUpdateCountdown--
			if m.AutoUpdateCountdown <= 0 {
				m.State = ViewUpdating
				m.ProgressMsg = "检查所有更新..."
				m.AutoUpdateCancelled = true
				return m, m.sendCommand(controller.Command{Type: controller.CmdAutoUpdate})
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

// handleControllerEvent handles events from the controller
func (m Model) handleControllerEvent(evt controller.Event) (tea.Model, tea.Cmd) {
	switch evt.Type {
	case controller.EvtProgressUpdate:
		payload := evt.Payload.(controller.ProgressUpdatePayload)
		m.ProgressMsg = payload.Message
		m.IsDownloading = payload.IsDownload
		m.DownloadSource = payload.Source
		m.DownloadFileName = payload.FileName
		m.Downloaded = payload.Downloaded
		m.TotalSize = payload.TotalSize
		m.DownloadSpeed = payload.Speed

		cmd := m.Progress.SetPercent(payload.Percent)
		return m, tea.Batch(cmd, listenForEvents(m.EventChan))

	case controller.EvtUpdateSuccess:
		payload := evt.Payload.(controller.UpdateCompletePayload)
		m.Updating = false
		m.State = ViewResult
		m.ResultSuccess = true
		m.ResultSkipped = payload.Skipped
		m.ResultMsg = payload.Message

		if payload.UpdatedComponents != nil {
			m.AutoUpdateResult = &AutoUpdateDetails{
				UpdatedComponents: payload.UpdatedComponents,
				SkippedComponents: payload.SkippedComponents,
				ComponentVersions: payload.ComponentVersions,
			}
		} else {
			m.AutoUpdateResult = nil
		}

		return m, listenForEvents(m.EventChan)

	case controller.EvtUpdateFailure:
		payload := evt.Payload.(controller.UpdateCompletePayload)
		m.Updating = false
		m.State = ViewResult
		m.ResultSuccess = false
		m.ResultMsg = payload.Message
		m.AutoUpdateResult = nil

		return m, listenForEvents(m.EventChan)

	case controller.EvtUpdateSkipped:
		payload := evt.Payload.(controller.UpdateCompletePayload)
		m.Updating = false
		m.State = ViewResult
		m.ResultSuccess = true
		m.ResultSkipped = true
		m.ResultMsg = payload.Message

		if payload.UpdatedComponents != nil {
			m.AutoUpdateResult = &AutoUpdateDetails{
				UpdatedComponents: payload.UpdatedComponents,
				SkippedComponents: payload.SkippedComponents,
				ComponentVersions: payload.ComponentVersions,
			}
		} else {
			m.AutoUpdateResult = nil
		}

		return m, listenForEvents(m.EventChan)

	case controller.EvtConfigUpdated:
		// Configuration updated successfully
		// Update is already in cfg, just continue listening
		return m, listenForEvents(m.EventChan)

	case controller.EvtWizardComplete:
		m.State = ViewMenu
		return m, listenForEvents(m.EventChan)

	case controller.EvtError:
		payload := evt.Payload.(controller.ErrorPayload)
		m.Err = payload.Error
		return m, listenForEvents(m.EventChan)
	}

	return m, listenForEvents(m.EventChan)
}

// sendCommand sends a command to the controller
func (m Model) sendCommand(cmd controller.Command) tea.Cmd {
	return func() tea.Msg {
		select {
		case m.CommandChan <- cmd:
		default:
			// Channel full, command dropped
		}
		return nil
	}
}

// listenForEvents listens for events from the controller
func listenForEvents(eventChan <-chan controller.Event) tea.Cmd {
	return func() tea.Msg {
		return <-eventChan
	}
}

// countdownTick 倒计时命令
func countdownTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return CountdownTickMsg{}
	})
}
