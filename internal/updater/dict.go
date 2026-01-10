package updater

import (
	"fmt"
	"os"
	"path/filepath"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/fileutil"
	"rime-wanxiang-updater/internal/types"
)

// DictUpdater 词库更新器
type DictUpdater struct {
	*BaseUpdater
	UpdateInfo *types.UpdateInfo
}

// NewDictUpdater 创建词库更新器
func NewDictUpdater(cfg *config.Manager) *DictUpdater {
	return &DictUpdater{
		BaseUpdater: NewBaseUpdater(cfg),
	}
}

// GetStatus 获取更新状态
func (d *DictUpdater) GetStatus() (*types.UpdateStatus, error) {
	// 获取远程版本信息
	remoteInfo, err := d.CheckUpdate()
	if err != nil {
		return nil, err
	}

	// 检查关键文件是否存在
	keyFile := filepath.Join(d.Config.GetDictExtractPath(), "chengyu.txt")
	keyFileExists := fileutil.FileExists(keyFile)

	// 获取本地版本信息
	recordPath := d.Config.GetDictRecordPath()
	localRecord := d.GetLocalRecord(recordPath)

	status := &types.UpdateStatus{
		RemoteVersion: remoteInfo.Tag,
		RemoteTime:    remoteInfo.UpdateTime,
		NeedsUpdate:   true,
	}

	// 如果关键文件不存在，强制更新
	if !keyFileExists {
		status.LocalVersion = "未安装"
		status.Message = fmt.Sprintf("检测到可用版本: %s (关键文件缺失)", remoteInfo.Tag)
		return status, nil
	}

	if localRecord != nil {
		status.LocalVersion = localRecord.Tag
		status.LocalTime = localRecord.UpdateTime
		status.NeedsUpdate = remoteInfo.UpdateTime.After(localRecord.UpdateTime)

		if status.NeedsUpdate {
			status.Message = fmt.Sprintf("发现新版本: %s → %s", localRecord.Tag, remoteInfo.Tag)
		} else {
			status.Message = fmt.Sprintf("已是最新版本 (当前版本: %s)", remoteInfo.Tag)
		}
	} else {
		status.LocalVersion = "未安装"
		status.Message = fmt.Sprintf("检测到可用版本: %s", remoteInfo.Tag)
	}

	return status, nil
}

// CheckUpdate 检查更新
func (d *DictUpdater) CheckUpdate() (*types.UpdateInfo, error) {
	var releases []types.GitHubRelease
	var err error

	if d.Config.Config.UseMirror {
		// CNB 使用 v1.0.0 tag
		releases, err = d.APIClient.FetchCNBReleases(types.OWNER, types.CNB_REPO, types.CNB_DICT_TAG)
	} else {
		// GitHub 使用 dict-nightly tag
		releases, err = d.APIClient.FetchGitHubReleases(types.OWNER, types.REPO, types.DICT_TAG)
	}

	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("未找到任何发布版本")
	}

	for _, release := range releases {
		for _, asset := range release.Assets {
			if asset.Name == d.Config.Config.DictFile {
				return &types.UpdateInfo{
					Name:       asset.Name,
					URL:        asset.BrowserDownloadURL,
					UpdateTime: asset.UpdatedAt,
					Tag:        release.TagName,
					Size:       asset.Size,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到匹配的词库文件: %s", d.Config.Config.DictFile)
}

// Run 执行更新
func (d *DictUpdater) Run(progress types.ProgressFunc) error {
	if progress == nil {
		progress = func(string, float64, string, string, int64, int64, float64, bool) {} // 空函数避免 nil 检查
	}

	// 执行更新前 hook
	if d.Config.Config.PreUpdateHook != "" {
		progress("执行更新前 hook...", 0.02, "", "", 0, 0, 0, false)
		if err := d.Config.ExecutePreUpdateHook(); err != nil {
			return fmt.Errorf("pre-update hook 失败，已取消更新: %w", err)
		}
	}

	// 显示下载源
	source := "GitHub"
	if d.Config.Config.UseMirror {
		source = "CNB 镜像"
	}
	progress(fmt.Sprintf("正在检查词库更新 [%s]...", source), 0.05, "", "", 0, 0, 0, false)

	if d.UpdateInfo == nil {
		info, err := d.CheckUpdate()
		if err != nil {
			return err
		}
		d.UpdateInfo = info
	}

	if d.UpdateInfo == nil {
		return fmt.Errorf("未找到词库更新")
	}

	recordPath := d.Config.GetDictRecordPath()
	targetFile := filepath.Join(d.Config.CacheDir, d.Config.Config.DictFile)

	// 校验本地文件
	progress("正在校验本地文件...", 0.1, "", "", 0, 0, 0, false)
	if d.UpdateInfo.SHA256 != "" && d.CompareHash(d.UpdateInfo.SHA256, targetFile) {
		progress("本地文件已是最新版本", 1.0, "", "", 0, 0, 0, false)
		d.SaveRecord(recordPath, "dict_file", d.Config.Config.DictFile, d.UpdateInfo)
		return nil
	}

	// 下载文件
	progress(fmt.Sprintf("准备从 %s 下载词库...", source), 0.15, "", "", 0, 0, 0, false)
	tempFile := filepath.Join(d.Config.CacheDir, fmt.Sprintf("temp_dict_%s.zip", d.UpdateInfo.SHA256))
	if err := d.DownloadFileWithValidation(d.UpdateInfo.URL, tempFile, d.Config.Config.DictFile, source, d.UpdateInfo.Size, progress); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// 清理旧文件
	progress("正在清理旧文件...", 0.7, "", "", 0, 0, 0, false)
	if fileutil.FileExists(targetFile) {
		d.CleanOldFiles(targetFile, tempFile, d.Config.GetDictExtractPath(), true)
	}

	// 应用更新
	progress("正在应用更新...", 0.8, "", "", 0, 0, 0, false)
	return d.applyUpdate(tempFile, targetFile, progress)
}

// applyUpdate 应用更新
func (d *DictUpdater) applyUpdate(temp, target string, progress types.ProgressFunc) error {
	// 终止进程（组合更新时跳过）
	if !d.SkipTerminate {
		progress("正在终止相关进程...", 0.85, "", "", 0, 0, 0, false)
		if err := d.TerminateProcesses(); err != nil {
			return fmt.Errorf("终止进程失败: %w", err)
		}
	}

	// 确保词库目录存在
	dictDir := d.Config.GetDictExtractPath()
	os.MkdirAll(dictDir, 0755)

	// 解压文件到主引擎目录
	progress("正在解压词库文件...", 0.9, "", "", 0, 0, 0, false)
	if err := d.ExtractZip(temp, dictDir); err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	// 处理 CNB 镜像的嵌套目录问题
	if d.Config.Config.UseMirror {
		if err := fileutil.HandleCNBNestedDir(dictDir, d.Config.Config.DictFile); err != nil {
			return fmt.Errorf("处理嵌套目录失败: %w", err)
		}
	}

	// 同步到其他引擎目录
	if len(d.Config.Config.InstalledEngines) > 1 {
		progress("正在同步词库到其他引擎...", 0.92, "", "", 0, 0, 0, false)
		if err := d.syncDictToOtherEngines(dictDir); err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("同步词库到其他引擎失败: %v", err), 0.92, "", "", 0, 0, 0, false)
		}
	}

	// 重命名临时文件
	progress("正在保存文件...", 0.95, "", "", 0, 0, 0, false)
	if fileutil.FileExists(target) {
		os.Remove(target)
	}
	if err := fileutil.MoveFile(temp, target); err != nil {
		return fmt.Errorf("重命名失败: %w", err)
	}

	// 保存记录
	recordPath := d.Config.GetDictRecordPath()
	if err := d.SaveRecord(recordPath, "dict_file", d.Config.Config.DictFile, d.UpdateInfo); err != nil {
		return err
	}

	// 执行更新后 hook（失败不影响更新结果）
	if d.Config.Config.PostUpdateHook != "" {
		progress("执行更新后 hook...", 1.0, "", "", 0, 0, 0, false)
		if err := d.Config.ExecutePostUpdateHook(); err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("post-update hook 失败: %v", err), 1.0, "", "", 0, 0, 0, false)
		}
	}

	// 同步到 fcitx 目录（如果启用）
	if d.Config.Config.FcitxCompat {
		_, _, err := d.Config.SyncToFcitxDir()
		if err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("fcitx 同步失败: %v", err), 1.0, "", "", 0, 0, 0, false)
		}
	}

	progress("更新完成！", 1.0, "", "", 0, 0, 0, false)
	return nil
}

// syncDictToOtherEngines 同步词库到其他引擎目录
func (d *DictUpdater) syncDictToOtherEngines(sourceDictDir string) error {
	// 检查是否配置了要更新的引擎列表
	updateEngines := d.Config.Config.UpdateEngines
	if len(updateEngines) == 0 {
		// 未配置：默认更新所有已安装的引擎
		updateEngines = d.Config.Config.InstalledEngines
	}

	// 如果只有一个引擎需要更新，跳过同步
	if len(updateEngines) <= 1 {
		return nil
	}

	primaryEngine := d.Config.Config.PrimaryEngine
	if primaryEngine == "" && len(updateEngines) > 0 {
		primaryEngine = updateEngines[0]
	}

	var errors []string

	// 遍历用户选择要更新的引擎
	for _, engine := range updateEngines {
		// 跳过主引擎（已经解压到主引擎目录）
		if engine == primaryEngine {
			continue
		}

		// 获取目标引擎的数据目录
		targetRimeDir := config.GetEngineDataDir(engine)
		if targetRimeDir == "" {
			errors = append(errors, fmt.Sprintf("无法获取引擎 %s 的数据目录", engine))
			continue
		}

		// 词库子目录
		targetDictDir := filepath.Join(targetRimeDir, d.Config.ZhDictsDir)

		// 确保目标目录存在
		if err := os.MkdirAll(targetDictDir, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("创建引擎 %s 词库目录失败: %v", engine, err))
			continue
		}

		// 复制词库文件
		if err := fileutil.SyncDirectory(sourceDictDir, targetDictDir, d.Config.Config.ExcludeFiles); err != nil {
			errors = append(errors, fmt.Sprintf("同步词库到引擎 %s 失败: %v", engine, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分引擎词库同步失败: %v", errors)
	}

	return nil
}
