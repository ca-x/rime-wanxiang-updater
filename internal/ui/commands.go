package ui

import (
	"fmt"
	"time"

	"rime-wanxiang-updater/internal/updater"

	tea "github.com/charmbracelet/bubbletea"
)

// runDictUpdate 运行词库更新
func (m *Model) runDictUpdate() tea.Cmd {
	m.ProgressChan = make(chan UpdateMsg, 100)
	m.CompletionChan = make(chan UpdateCompleteMsg, 1)

	go func() {
		dictUpdater := updater.NewDictUpdater(m.Cfg)

		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.ProgressChan <- UpdateMsg{
				Message:      message,
				Percent:      percent,
				Source:       source,
				FileName:     fileName,
				Downloaded:   downloaded,
				Total:        total,
				Speed:        speed,
				DownloadMode: downloadMode,
			}:
			default:
			}
		}

		status, err := dictUpdater.GetStatus()
		if err != nil {
			m.CompletionChan <- UpdateCompleteMsg{Err: err, UpdateType: "词库", Skipped: false, StatusMessage: ""}
			close(m.ProgressChan)
			return
		}

		if !status.NeedsUpdate {
			progressFunc("词库已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			m.CompletionChan <- UpdateCompleteMsg{Err: nil, UpdateType: "词库", Skipped: true, StatusMessage: status.Message}
			close(m.ProgressChan)
			return
		}

		if err = dictUpdater.Run(progressFunc); err == nil {
			err = dictUpdater.Deploy()
		}

		m.CompletionChan <- UpdateCompleteMsg{Err: err, UpdateType: "词库", Skipped: false, StatusMessage: ""}
		close(m.ProgressChan)
	}()

	return listenForProgress(m.ProgressChan, m.CompletionChan)
}

// runSchemeUpdate 运行方案更新
func (m *Model) runSchemeUpdate() tea.Cmd {
	m.ProgressChan = make(chan UpdateMsg, 100)
	m.CompletionChan = make(chan UpdateCompleteMsg, 1)

	go func() {
		schemeUpdater := updater.NewSchemeUpdater(m.Cfg)

		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.ProgressChan <- UpdateMsg{
				Message:      message,
				Percent:      percent,
				Source:       source,
				FileName:     fileName,
				Downloaded:   downloaded,
				Total:        total,
				Speed:        speed,
				DownloadMode: downloadMode,
			}:
			default:
			}
		}

		status, err := schemeUpdater.GetStatus()
		if err != nil {
			m.CompletionChan <- UpdateCompleteMsg{Err: err, UpdateType: "方案", Skipped: false, StatusMessage: ""}
			close(m.ProgressChan)
			return
		}

		if !status.NeedsUpdate {
			progressFunc("方案已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			m.CompletionChan <- UpdateCompleteMsg{Err: nil, UpdateType: "方案", Skipped: true, StatusMessage: status.Message}
			close(m.ProgressChan)
			return
		}

		if err = schemeUpdater.Run(progressFunc); err == nil {
			err = schemeUpdater.Deploy()
		}

		m.CompletionChan <- UpdateCompleteMsg{Err: err, UpdateType: "方案", Skipped: false, StatusMessage: ""}
		close(m.ProgressChan)
	}()

	return listenForProgress(m.ProgressChan, m.CompletionChan)
}

// runModelUpdate 运行模型更新
func (m *Model) runModelUpdate() tea.Cmd {
	m.ProgressChan = make(chan UpdateMsg, 100)
	m.CompletionChan = make(chan UpdateCompleteMsg, 1)

	go func() {
		modelUpdater := updater.NewModelUpdater(m.Cfg)

		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.ProgressChan <- UpdateMsg{
				Message:      message,
				Percent:      percent,
				Source:       source,
				FileName:     fileName,
				Downloaded:   downloaded,
				Total:        total,
				Speed:        speed,
				DownloadMode: downloadMode,
			}:
			default:
			}
		}

		var err error
		if err = modelUpdater.Run(progressFunc); err == nil {
			err = modelUpdater.Deploy()
		}

		m.CompletionChan <- UpdateCompleteMsg{Err: err, UpdateType: "模型", Skipped: false, StatusMessage: ""}
		close(m.ProgressChan)
	}()

	return listenForProgress(m.ProgressChan, m.CompletionChan)
}

// runAutoUpdate 运行自动更新
func (m *Model) runAutoUpdate() tea.Cmd {
	m.ProgressChan = make(chan UpdateMsg, 100)
	m.CompletionChan = make(chan UpdateCompleteMsg, 1)

	go func() {
		combined := updater.NewCombinedUpdater(m.Cfg)

		progressFunc := func(component, message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			select {
			case m.ProgressChan <- UpdateMsg{
				Message:      fmt.Sprintf("[%s] %s", component, message),
				Percent:      percent,
				Source:       source,
				FileName:     fileName,
				Downloaded:   downloaded,
				Total:        total,
				Speed:        speed,
				DownloadMode: downloadMode,
			}:
			default:
			}
		}

		progressFunc("检查", "正在检查所有更新...", 0.0, "", "", 0, 0, 0, false)
		if err := combined.FetchAllUpdates(); err != nil {
			m.CompletionChan <- UpdateCompleteMsg{Err: err, UpdateType: "自动", Skipped: false, StatusMessage: "", AutoUpdateDetails: nil}
			close(m.ProgressChan)
			return
		}

		if !combined.HasAnyUpdate() {
			progressFunc("完成", "所有组件已是最新版本", 1.0, "", "", 0, 0, 0, false)
			componentVersions := make(map[string]string)
			if schemeStatus, err := combined.SchemeUpdater.GetStatus(); err == nil {
				componentVersions["方案"] = schemeStatus.LocalVersion
			}
			if dictStatus, err := combined.DictUpdater.GetStatus(); err == nil {
				componentVersions["词库"] = dictStatus.LocalVersion
			}
			if modelStatus, err := combined.ModelUpdater.GetStatus(); err == nil {
				componentVersions["模型"] = modelStatus.LocalVersion
			}

			details := &AutoUpdateDetails{
				UpdatedComponents: []string{},
				SkippedComponents: []string{"方案", "词库", "模型"},
				ComponentVersions: componentVersions,
			}
			m.CompletionChan <- UpdateCompleteMsg{Err: nil, UpdateType: "自动", Skipped: true, StatusMessage: "所有组件已是最新版本", AutoUpdateDetails: details}
			close(m.ProgressChan)
			return
		}

		result, err := combined.RunAllWithProgress(progressFunc)

		var details *AutoUpdateDetails
		if result != nil {
			details = &AutoUpdateDetails{
				UpdatedComponents: result.UpdatedComponents,
				SkippedComponents: result.SkippedComponents,
				ComponentVersions: result.ComponentVersions,
			}
		}

		m.CompletionChan <- UpdateCompleteMsg{
			Err:               err,
			UpdateType:        "自动",
			Skipped:           err == nil && len(result.UpdatedComponents) == 0,
			StatusMessage:     "",
			AutoUpdateDetails: details,
		}
		close(m.ProgressChan)
	}()

	return listenForProgress(m.ProgressChan, m.CompletionChan)
}

// listenForProgress 持续监听进度更新
func listenForProgress(progressChan chan UpdateMsg, completeChan chan UpdateCompleteMsg) tea.Cmd {
	return func() tea.Msg {
		select {
		case msg, ok := <-progressChan:
			if ok {
				return msg
			}
			select {
			case msg := <-completeChan:
				return msg
			default:
				return <-completeChan
			}
		case msg := <-completeChan:
			return msg
		}
	}
}

// countdownTick 倒计时命令
func countdownTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return CountdownTickMsg{}
	})
}
