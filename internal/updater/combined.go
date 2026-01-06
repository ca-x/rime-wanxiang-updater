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
	hasScheme := c.SchemeUpdater.UpdateInfo != nil &&
		c.SchemeUpdater.HasUpdate(c.SchemeUpdater.UpdateInfo, c.Config.GetSchemeRecordPath())

	hasDict := c.DictUpdater.UpdateInfo != nil &&
		c.DictUpdater.HasUpdate(c.DictUpdater.UpdateInfo, c.Config.GetDictRecordPath())

	hasModel := c.ModelUpdater.UpdateInfo != nil &&
		c.ModelUpdater.HasUpdate(c.ModelUpdater.UpdateInfo, c.Config.GetModelRecordPath())

	return hasScheme || hasDict || hasModel
}

// RunAll 执行所有更新
func (c *CombinedUpdater) RunAll() error {
	var errors []string

	// 更新方案
	if c.SchemeUpdater.UpdateInfo != nil &&
		c.SchemeUpdater.HasUpdate(c.SchemeUpdater.UpdateInfo, c.Config.GetSchemeRecordPath()) {
		if err := c.SchemeUpdater.Run(nil); err != nil {
			errors = append(errors, fmt.Sprintf("方案更新失败: %v", err))
		}
	}

	// 更新词库
	if c.DictUpdater.UpdateInfo != nil &&
		c.DictUpdater.HasUpdate(c.DictUpdater.UpdateInfo, c.Config.GetDictRecordPath()) {
		if err := c.DictUpdater.Run(nil); err != nil {
			errors = append(errors, fmt.Sprintf("词库更新失败: %v", err))
		}
	}

	// 更新模型
	if c.ModelUpdater.UpdateInfo != nil &&
		c.ModelUpdater.HasUpdate(c.ModelUpdater.UpdateInfo, c.Config.GetModelRecordPath()) {
		if err := c.ModelUpdater.Run(nil); err != nil {
			errors = append(errors, fmt.Sprintf("模型更新失败: %v", err))
		}
	}

	// 如果没有错误，执行部署
	if len(errors) == 0 {
		if err := c.SchemeUpdater.Deploy(); err != nil {
			errors = append(errors, fmt.Sprintf("部署失败: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("更新过程中出现错误: %v", errors)
	}

	return nil
}
