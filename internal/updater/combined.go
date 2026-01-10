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
	_, err := c.RunAllWithProgress(nil)
	return err
}

// UpdateResult 更新结果
type UpdateResult struct {
	UpdatedComponents []string          // 已更新的组件
	SkippedComponents []string          // 跳过的组件（已是最新版本）
	ComponentVersions map[string]string // 组件版本信息（组件名 -> 版本号）
}

// RunAllWithProgress 执行所有更新并报告进度
func (c *CombinedUpdater) RunAllWithProgress(progress func(component, message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool)) (*UpdateResult, error) {
	var errors []string
	result := &UpdateResult{
		UpdatedComponents: []string{},
		SkippedComponents: []string{},
		ComponentVersions: make(map[string]string),
	}

	// 如果没有提供进度回调，使用空函数
	if progress == nil {
		progress = func(string, string, float64, string, string, int64, int64, float64, bool) {}
	}

	// 收集需要更新的项
	needsSchemeUpdate := false
	needsDictUpdate := false
	needsModelUpdate := false

	if c.SchemeUpdater.UpdateInfo != nil {
		if schemeStatus, err := c.SchemeUpdater.GetStatus(); err == nil && schemeStatus.NeedsUpdate {
			needsSchemeUpdate = true
		} else if err == nil {
			result.SkippedComponents = append(result.SkippedComponents, "方案")
			result.ComponentVersions["方案"] = schemeStatus.LocalVersion
		}
	}

	if c.DictUpdater.UpdateInfo != nil {
		if dictStatus, err := c.DictUpdater.GetStatus(); err == nil && dictStatus.NeedsUpdate {
			needsDictUpdate = true
		} else if err == nil {
			result.SkippedComponents = append(result.SkippedComponents, "词库")
			result.ComponentVersions["词库"] = dictStatus.LocalVersion
		}
	}

	if c.ModelUpdater.UpdateInfo != nil &&
		c.ModelUpdater.HasUpdate(c.ModelUpdater.UpdateInfo, c.Config.GetModelRecordPath()) {
		needsModelUpdate = true
	} else if c.ModelUpdater.UpdateInfo != nil {
		if modelStatus, err := c.ModelUpdater.GetStatus(); err == nil {
			result.SkippedComponents = append(result.SkippedComponents, "模型")
			result.ComponentVersions["模型"] = modelStatus.LocalVersion
		}
	}

	// 如果没有任何更新，直接返回
	if !needsSchemeUpdate && !needsDictUpdate && !needsModelUpdate {
		progress("完成", "已是最新版本", 1.0, "", "", 0, 0, 0, false)
		return result, nil
	}

	// 统一在开始前终止进程（只终止一次）
	progress("准备", "正在终止相关进程...", 0.0, "", "", 0, 0, 0, false)
	if err := c.SchemeUpdater.TerminateProcesses(); err != nil {
		return result, fmt.Errorf("终止进程失败: %w", err)
	}

	// 标记为组合更新模式，让子更新器跳过终止进程步骤
	c.SchemeUpdater.SkipTerminate = true
	c.DictUpdater.SkipTerminate = true
	c.ModelUpdater.SkipTerminate = true

	defer func() {
		// 恢复默认设置
		c.SchemeUpdater.SkipTerminate = false
		c.DictUpdater.SkipTerminate = false
		c.ModelUpdater.SkipTerminate = false
	}()

	// 更新方案
	if needsSchemeUpdate {
		progress("方案", "正在更新方案...", 0.05, "", "", 0, 0, 0, false)
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			progress("方案", message, 0.05+percent*0.30, source, fileName, downloaded, total, speed, downloadMode) // 方案占 30%
		}
		if err := c.SchemeUpdater.Run(progressFunc); err != nil {
			errors = append(errors, fmt.Sprintf("方案更新失败: %v", err))
		} else {
			result.UpdatedComponents = append(result.UpdatedComponents, "方案")
		}
	}

	// 更新词库
	if needsDictUpdate {
		progress("词库", "正在更新词库...", 0.35, "", "", 0, 0, 0, false)
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			progress("词库", message, 0.35+percent*0.30, source, fileName, downloaded, total, speed, downloadMode) // 词库占 30%
		}
		if err := c.DictUpdater.Run(progressFunc); err != nil {
			errors = append(errors, fmt.Sprintf("词库更新失败: %v", err))
		} else {
			result.UpdatedComponents = append(result.UpdatedComponents, "词库")
		}
	}

	// 更新模型
	if needsModelUpdate {
		progress("模型", "正在更新模型...", 0.65, "", "", 0, 0, 0, false)
		progressFunc := func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool) {
			progress("模型", message, 0.65+percent*0.25, source, fileName, downloaded, total, speed, downloadMode) // 模型占 25%
		}
		if err := c.ModelUpdater.Run(progressFunc); err != nil {
			errors = append(errors, fmt.Sprintf("模型更新失败: %v", err))
		} else {
			result.UpdatedComponents = append(result.UpdatedComponents, "模型")
		}
	}

	// 如果没有错误，执行部署（会重启服务）
	if len(errors) == 0 {
		// 获取要部署的引擎列表
		deployEngines := c.Config.Config.UpdateEngines
		if len(deployEngines) == 0 {
			// 未配置：默认部署所有已安装的引擎
			deployEngines = c.Config.Config.InstalledEngines
		}

		// 检查是否有多个引擎需要部署
		if len(deployEngines) > 1 {
			// 定义接口用于类型断言
			type multiEngineDeployer interface {
				DeployToAllEnginesWithProgress(progressFunc func(engine string, index, total int)) error
			}

			if med, ok := c.SchemeUpdater.Deployer.(multiEngineDeployer); ok {
				// 支持多引擎进度回调的 deployer
				err := med.DeployToAllEnginesWithProgress(func(engine string, index, total int) {
					deployMsg := fmt.Sprintf("正在部署到 %s (%d/%d)...", engine, index, total)
					deployPercent := 0.90 + float64(index-1)/float64(total)*0.09 // 0.90-0.99
					progress("部署", deployMsg, deployPercent, "", "", 0, 0, 0, false)
				})
				if err != nil {
					errors = append(errors, fmt.Sprintf("部署失败: %v", err))
				}
			} else {
				// 回退到单引擎部署
				engineName := c.Config.GetEngineDisplayName()
				progress("部署", fmt.Sprintf("正在部署到 %s...", engineName), 0.90, "", "", 0, 0, 0, false)
				if err := c.SchemeUpdater.Deploy(); err != nil {
					errors = append(errors, fmt.Sprintf("部署失败: %v", err))
				}
			}
		} else {
			// 单引擎部署
			engineName := "输入法"
			if len(deployEngines) == 1 {
				engineName = deployEngines[0]
			} else if c.Config.Config.PrimaryEngine != "" {
				engineName = c.Config.Config.PrimaryEngine
			}
			progress("部署", fmt.Sprintf("正在部署到 %s...", engineName), 0.90, "", "", 0, 0, 0, false)
			if err := c.SchemeUpdater.Deploy(); err != nil {
				errors = append(errors, fmt.Sprintf("部署失败: %v", err))
			}
		}
	} else {
		// 即使有错误，也尝试重启服务，让用户能继续使用输入法
		progress("恢复", "尝试重启服务...", 0.90, "", "", 0, 0, 0, false)
		_ = c.SchemeUpdater.Deploy() // 忽略错误
	}

	if len(errors) > 0 {
		return result, fmt.Errorf("更新过程中出现错误: %v", errors)
	}

	progress("完成", "所有更新已完成", 1.0, "", "", 0, 0, 0, false)
	return result, nil
}
