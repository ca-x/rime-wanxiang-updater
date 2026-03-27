package controller

import (
	"fmt"

	"rime-wanxiang-updater/internal/updater"
)

func isFirstInstall(localVersion string) bool {
	return localVersion == "未安装" || localVersion == "未知版本"
}

func successMessageForSingleUpdate(updateType, localVersion string) string {
	if isFirstInstall(localVersion) {
		return fmt.Sprintf("%s安装完成！", updateType)
	}

	return fmt.Sprintf("%s更新完成！", updateType)
}

func successMessageForAutoUpdate(updatedComponents []string, previousVersions map[string]string) string {
	if len(updatedComponents) == 0 {
		return "更新完成！"
	}

	hasInstall := false
	hasUpdate := false

	for _, component := range updatedComponents {
		if isFirstInstall(previousVersions[component]) {
			hasInstall = true
		} else {
			hasUpdate = true
		}
	}

	switch {
	case hasInstall && hasUpdate:
		return "安装和更新完成！"
	case hasInstall:
		return "安装完成！"
	default:
		return "更新完成！"
	}
}

// handleAutoUpdate handles the auto update command
func (c *Controller) handleAutoUpdate(cmd Command) {
	c.mu.Lock()
	if c.updating {
		c.mu.Unlock()
		c.emitError(fmt.Errorf("update already in progress"), "auto update")
		return
	}
	c.updating = true
	c.currentOperation = "auto"
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			c.updating = false
			c.currentOperation = ""
			c.mu.Unlock()
		}()

		combined := updater.NewCombinedUpdater(c.cfg)

		progressFunc := func(component, message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			c.emitProgress(component, message, percent, source, fileName, downloaded, total, speed, downloadMode)
		}

		progressFunc("检查", "正在检查所有更新...", 0.0, "", "", 0, 0, 0, false)
		if err := combined.FetchAllUpdates(); err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "自动",
				Success:    false,
				Message:    fmt.Sprintf("检查更新失败: %v", err),
			})
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

			c.emitEvent(EvtUpdateSkipped, UpdateCompletePayload{
				UpdateType:        "自动",
				Success:           true,
				Skipped:           true,
				Message:           "所有组件已是最新版本",
				UpdatedComponents: []string{},
				SkippedComponents: []string{"方案", "词库", "模型"},
				ComponentVersions: componentVersions,
			})
			return
		}

		result, err := combined.RunAllWithProgress(progressFunc)

		if err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "自动",
				Success:    false,
				Message:    fmt.Sprintf("更新失败: %v", err),
			})
			return
		}

		var updatedComponents, skippedComponents []string
		componentVersions := make(map[string]string)
		previousVersions := make(map[string]string)
		if result != nil {
			updatedComponents = result.UpdatedComponents
			skippedComponents = result.SkippedComponents
			componentVersions = result.ComponentVersions
			previousVersions = result.PreviousVersions
		}

		c.emitEvent(EvtUpdateSuccess, UpdateCompletePayload{
			UpdateType:        "自动",
			Success:           true,
			Skipped:           len(updatedComponents) == 0,
			Message:           successMessageForAutoUpdate(updatedComponents, previousVersions),
			UpdatedComponents: updatedComponents,
			SkippedComponents: skippedComponents,
			ComponentVersions: componentVersions,
		})
	}()
}

// handleUpdateDict handles the dict update command
func (c *Controller) handleUpdateDict(cmd Command) {
	c.mu.Lock()
	if c.updating {
		c.mu.Unlock()
		c.emitError(fmt.Errorf("update already in progress"), "dict update")
		return
	}
	c.updating = true
	c.currentOperation = "dict"
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			c.updating = false
			c.currentOperation = ""
			c.mu.Unlock()
		}()

		dictUpdater := updater.NewDictUpdater(c.cfg)

		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			c.emitProgress("词库", message, percent, source, fileName, downloaded, total, speed, downloadMode)
		}

		status, err := dictUpdater.GetStatus()
		if err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "词库",
				Success:    false,
				Message:    fmt.Sprintf("获取状态失败: %v", err),
			})
			return
		}

		if !status.NeedsUpdate {
			progressFunc("词库已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			c.emitEvent(EvtUpdateSkipped, UpdateCompletePayload{
				UpdateType: "词库",
				Success:    true,
				Skipped:    true,
				Message:    status.Message,
			})
			return
		}

		if err = dictUpdater.Run(progressFunc); err == nil {
			err = dictUpdater.Deploy()
		}

		if err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "词库",
				Success:    false,
				Message:    fmt.Sprintf("更新失败: %v", err),
			})
			return
		}

		c.emitEvent(EvtUpdateSuccess, UpdateCompletePayload{
			UpdateType: "词库",
			Success:    true,
			Skipped:    false,
			Message:    successMessageForSingleUpdate("词库", status.LocalVersion),
		})
	}()
}

// handleUpdateScheme handles the scheme update command
func (c *Controller) handleUpdateScheme(cmd Command) {
	c.mu.Lock()
	if c.updating {
		c.mu.Unlock()
		c.emitError(fmt.Errorf("update already in progress"), "scheme update")
		return
	}
	c.updating = true
	c.currentOperation = "scheme"
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			c.updating = false
			c.currentOperation = ""
			c.mu.Unlock()
		}()

		schemeUpdater := updater.NewSchemeUpdater(c.cfg)

		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			c.emitProgress("方案", message, percent, source, fileName, downloaded, total, speed, downloadMode)
		}

		status, err := schemeUpdater.GetStatus()
		if err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "方案",
				Success:    false,
				Message:    fmt.Sprintf("获取状态失败: %v", err),
			})
			return
		}

		if !status.NeedsUpdate {
			progressFunc("方案已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			c.emitEvent(EvtUpdateSkipped, UpdateCompletePayload{
				UpdateType: "方案",
				Success:    true,
				Skipped:    true,
				Message:    status.Message,
			})
			return
		}

		if err = schemeUpdater.Run(progressFunc); err == nil {
			err = schemeUpdater.Deploy()
		}

		if err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "方案",
				Success:    false,
				Message:    fmt.Sprintf("更新失败: %v", err),
			})
			return
		}

		c.emitEvent(EvtUpdateSuccess, UpdateCompletePayload{
			UpdateType: "方案",
			Success:    true,
			Skipped:    false,
			Message:    successMessageForSingleUpdate("方案", status.LocalVersion),
		})
	}()
}

// handleUpdateModel handles the model update command
func (c *Controller) handleUpdateModel(cmd Command) {
	c.mu.Lock()
	if c.updating {
		c.mu.Unlock()
		c.emitError(fmt.Errorf("update already in progress"), "model update")
		return
	}
	c.updating = true
	c.currentOperation = "model"
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			c.updating = false
			c.currentOperation = ""
			c.mu.Unlock()
		}()

		modelUpdater := updater.NewModelUpdater(c.cfg)

		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			c.emitProgress("模型", message, percent, source, fileName, downloaded, total, speed, downloadMode)
		}

		status, err := modelUpdater.GetStatus()
		if err != nil {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "模型",
				Success:    false,
				Message:    fmt.Sprintf("获取状态失败: %v", err),
			})
			return
		}

		if !status.NeedsUpdate {
			progressFunc("模型已是最新版本，跳过更新", 1.0, "", "", 0, 0, 0, false)
			c.emitEvent(EvtUpdateSkipped, UpdateCompletePayload{
				UpdateType: "模型",
				Success:    true,
				Skipped:    true,
				Message:    status.Message,
			})
			return
		}

		if err := modelUpdater.Run(progressFunc); err == nil {
			err = modelUpdater.Deploy()
			if err != nil {
				c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
					UpdateType: "模型",
					Success:    false,
					Message:    fmt.Sprintf("更新失败: %v", err),
				})
				return
			}
		} else {
			c.emitEvent(EvtUpdateFailure, UpdateCompletePayload{
				UpdateType: "模型",
				Success:    false,
				Message:    fmt.Sprintf("更新失败: %v", err),
			})
			return
		}

		c.emitEvent(EvtUpdateSuccess, UpdateCompletePayload{
			UpdateType: "模型",
			Success:    true,
			Skipped:    false,
			Message:    successMessageForSingleUpdate("模型", status.LocalVersion),
		})
	}()
}
