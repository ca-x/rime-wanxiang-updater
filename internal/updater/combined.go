package updater

import (
	"fmt"

	"rime-wanxiang-updater/internal/config"
)

// CombinedUpdater 组合更新器
type CombinedUpdater struct {
	Config        *config.Manager
	SchemeUpdater *SchemeUpdater
	DictUpdater   *DictUpdater
	ModelUpdater  *ModelUpdater
}

// NewCombinedUpdater 创建组合更新器
func NewCombinedUpdater(cfg *config.Manager) *CombinedUpdater {
	return &CombinedUpdater{
		Config:        cfg,
		SchemeUpdater: NewSchemeUpdater(cfg),
		DictUpdater:   NewDictUpdater(cfg),
		ModelUpdater:  NewModelUpdater(cfg),
	}
}

// FetchAllUpdates 获取所有更新信息
func (c *CombinedUpdater) FetchAllUpdates() error {
	var errors []string

	// 检查方案更新
	schemeInfo, err := c.SchemeUpdater.CheckUpdate()
	if err != nil {
		errors = append(errors, fmt.Sprintf("方案: %v", err))
	}
	c.SchemeUpdater.UpdateInfo = schemeInfo

	// 检查词库更新
	dictInfo, err := c.DictUpdater.CheckUpdate()
	if err != nil {
		errors = append(errors, fmt.Sprintf("词库: %v", err))
	}
	c.DictUpdater.UpdateInfo = dictInfo

	// 检查模型更新
	modelInfo, err := c.ModelUpdater.CheckUpdate()
	if err != nil {
		errors = append(errors, fmt.Sprintf("模型: %v", err))
	}
	c.ModelUpdater.UpdateInfo = modelInfo

	if len(errors) > 0 {
		return fmt.Errorf("检查更新失败: %v", errors)
	}

	return nil
}

// HasAnyUpdate 检查是否有任何更新
func (c *CombinedUpdater) HasAnyUpdate() bool {
	// 检查方案更新 - 使用 GetStatus 来检查实际文件
	if schemeStatus, err := c.SchemeUpdater.GetStatus(); err == nil && schemeStatus.NeedsUpdate {
		return true
	}

	// 检查词库更新 - 使用 GetStatus 来检查实际文件
	if dictStatus, err := c.DictUpdater.GetStatus(); err == nil && dictStatus.NeedsUpdate {
		return true
	}

	// 检查模型更新 - 保持原有逻辑
	hasModel := c.ModelUpdater.UpdateInfo != nil &&
		c.ModelUpdater.HasUpdate(c.ModelUpdater.UpdateInfo, c.Config.GetModelRecordPath())

	return hasModel
}

// RunAll 执行所有更新
func (c *CombinedUpdater) RunAll() error {
	return c.RunAllWithProgress(nil)
}

// RunAllWithProgress 执行所有更新并报告进度
func (c *CombinedUpdater) RunAllWithProgress(progress func(component, message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool)) error {
	var errors []string

	// 如果没有提供进度回调，使用空函数
	if progress == nil {
		progress = func(string, string, float64, string, string, int64, int64, float64, bool) {}
	}

	// 更新方案 - 使用 GetStatus 检查是否需要更新
	if c.SchemeUpdater.UpdateInfo != nil {
		if schemeStatus, err := c.SchemeUpdater.GetStatus(); err == nil && schemeStatus.NeedsUpdate {
			progress("方案", "正在更新方案...", 0.0, "", "", 0, 0, 0, false)
			progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
				progress("方案", message, percent*0.33, source, fileName, downloaded, total, speed, downloadMode) // 方案占 33%
			}
			if err := c.SchemeUpdater.Run(progressFunc); err != nil {
				errors = append(errors, fmt.Sprintf("方案更新失败: %v", err))
			}
		}
	}

	// 更新词库 - 使用 GetStatus 检查是否需要更新
	if c.DictUpdater.UpdateInfo != nil {
		if dictStatus, err := c.DictUpdater.GetStatus(); err == nil && dictStatus.NeedsUpdate {
			progress("词库", "正在更新词库...", 0.33, "", "", 0, 0, 0, false)
			progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
				progress("词库", message, 0.33+percent*0.33, source, fileName, downloaded, total, speed, downloadMode) // 词库占 33%
			}
			if err := c.DictUpdater.Run(progressFunc); err != nil {
				errors = append(errors, fmt.Sprintf("词库更新失败: %v", err))
			}
		}
	}

	// 更新模型 - 保持原有逻辑
	if c.ModelUpdater.UpdateInfo != nil &&
		c.ModelUpdater.HasUpdate(c.ModelUpdater.UpdateInfo, c.Config.GetModelRecordPath()) {
		progress("模型", "正在更新模型...", 0.66, "", "", 0, 0, 0, false)
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			progress("模型", message, 0.66+percent*0.34, source, fileName, downloaded, total, speed, downloadMode) // 模型占 34%
		}
		if err := c.ModelUpdater.Run(progressFunc); err != nil {
			errors = append(errors, fmt.Sprintf("模型更新失败: %v", err))
		}
	}

	// 如果没有错误，执行部署
	if len(errors) == 0 {
		progress("部署", "正在部署...", 0.95, "", "", 0, 0, 0, false)
		if err := c.SchemeUpdater.Deploy(); err != nil {
			errors = append(errors, fmt.Sprintf("部署失败: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("更新过程中出现错误: %v", errors)
	}

	progress("完成", "所有更新已完成", 1.0, "", "", 0, 0, 0, false)
	return nil
}
