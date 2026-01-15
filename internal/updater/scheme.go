package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/fileutil"
	"rime-wanxiang-updater/internal/types"
)

// SchemeUpdater 方案更新器
type SchemeUpdater struct {
	*BaseUpdater
	UpdateInfo *types.UpdateInfo
}

// NewSchemeUpdater 创建方案更新器
func NewSchemeUpdater(cfg *config.Manager) *SchemeUpdater {
	return &SchemeUpdater{
		BaseUpdater: NewBaseUpdater(cfg),
	}
}

// GetStatus 获取更新状态
func (s *SchemeUpdater) GetStatus() (*types.UpdateStatus, error) {
	// 获取远程版本信息
	remoteInfo, err := s.CheckUpdate()
	if err != nil {
		return nil, err
	}

	// 检查关键文件是否存在
	keyFile := filepath.Join(s.Config.GetExtractPath(), "lua", "wanxiang.lua")
	keyFileExists := fileutil.FileExists(keyFile)

	// 获取本地版本信息
	recordPath := s.Config.GetSchemeRecordPath()
	localRecord := s.GetLocalRecord(recordPath)

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
		// 检查文件名是否匹配，如果不匹配说明方案已切换
		if localRecord.Name != s.Config.Config.SchemeFile {
			status.LocalVersion = fmt.Sprintf("已切换方案 (从 %s)", localRecord.Name)
			status.Message = fmt.Sprintf("检测到可用版本: %s (方案已切换，需要更新)", remoteInfo.Tag)
			status.NeedsUpdate = true
			return status, nil
		}

		status.LocalVersion = localRecord.Tag
		status.LocalTime = localRecord.UpdateTime
		status.NeedsUpdate = remoteInfo.UpdateTime.After(localRecord.UpdateTime)

		if status.NeedsUpdate {
			status.Message = fmt.Sprintf("发现新版本: %s → %s", localRecord.Tag, remoteInfo.Tag)
		} else {
			status.Message = fmt.Sprintf("已是最新版本 (当前版本: %s)", remoteInfo.Tag)
		}
	} else {
		// 无版本记录但关键文件存在，说明记录丢失或首次管理已有安装
		if keyFileExists {
			status.LocalVersion = "未知版本"
			status.Message = fmt.Sprintf("检测到可用版本: %s (无版本记录，将重新安装)", remoteInfo.Tag)
		} else {
			status.LocalVersion = "未安装"
			status.Message = fmt.Sprintf("检测到可用版本: %s", remoteInfo.Tag)
		}
	}

	return status, nil
}

// CheckUpdate 检查更新
func (s *SchemeUpdater) CheckUpdate() (*types.UpdateInfo, error) {
	var releases []types.GitHubRelease
	var err error

	if s.Config.Config.UseMirror {
		releases, err = s.APIClient.FetchCNBReleases(types.OWNER, types.CNB_REPO, "")
	} else {
		releases, err = s.APIClient.FetchGitHubReleases(types.OWNER, types.REPO, "")
	}

	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("未找到任何发布版本")
	}

	for _, release := range releases {
		for _, asset := range release.Assets {
			if asset.Name == s.Config.Config.SchemeFile {
				return &types.UpdateInfo{
					Name:        asset.Name,
					URL:         asset.BrowserDownloadURL,
					UpdateTime:  asset.UpdatedAt,
					Tag:         release.TagName,
					Description: release.Body,
					Size:        asset.Size,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到匹配的方案文件: %s", s.Config.Config.SchemeFile)
}

// Run 执行更新
func (s *SchemeUpdater) Run(progress types.ProgressFunc) error {
	if progress == nil {
		progress = func(string, float64, string, string, int64, int64, float64, bool) {} // 空函数避免 nil 检查
	}

	// 执行更新前 hook
	if s.Config.Config.PreUpdateHook != "" {
		progress("执行更新前 hook...", 0.02, "", "", 0, 0, 0, false)
		if err := s.Config.ExecutePreUpdateHook(); err != nil {
			return fmt.Errorf("pre-update hook 失败，已取消更新: %w", err)
		}
	}

	// 显示下载源
	source := "GitHub"
	if s.Config.Config.UseMirror {
		source = "CNB 镜像"
	}
	progress(fmt.Sprintf("正在检查方案更新 [%s]...", source), 0.05, "", "", 0, 0, 0, false)

	if s.UpdateInfo == nil {
		info, err := s.CheckUpdate()
		if err != nil {
			return err
		}
		s.UpdateInfo = info
	}

	if s.UpdateInfo == nil {
		return fmt.Errorf("未找到方案更新")
	}

	recordPath := s.Config.GetSchemeRecordPath()
	targetFile := filepath.Join(s.Config.CacheDir, s.Config.Config.SchemeFile)

	// 校验本地文件
	progress("正在校验本地文件...", 0.1, "", "", 0, 0, 0, false)
	localRecord := s.GetLocalRecord(recordPath)
	if localRecord != nil && localRecord.SHA256 != "" && s.CompareHash(localRecord.SHA256, targetFile) {
		progress("本地文件已是最新版本", 1.0, "", "", 0, 0, 0, false)
		s.SaveRecord(recordPath, "scheme_file", s.Config.Config.SchemeFile, s.UpdateInfo)
		return nil
	}

	// 下载文件
	progress(fmt.Sprintf("准备从 %s 下载方案...", source), 0.15, "", "", 0, 0, 0, false)
	tempFile := filepath.Join(s.Config.CacheDir, fmt.Sprintf("temp_scheme_%d.zip", time.Now().Unix()))
	if err := s.DownloadFileWithValidation(s.UpdateInfo.URL, tempFile, s.Config.Config.SchemeFile, source, s.UpdateInfo.Size, progress); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// 计算下载文件的 SHA256
	progress("正在计算文件校验和...", 0.65, "", "", 0, 0, 0, false)
	if hash, err := fileutil.CalculateSHA256(tempFile); err == nil {
		s.UpdateInfo.SHA256 = hash
	}

	// 清理旧文件
	progress("正在清理旧文件...", 0.7, "", "", 0, 0, 0, false)
	if fileutil.FileExists(targetFile) {
		s.CleanOldFiles(targetFile, tempFile, s.Config.GetExtractPath(), false)
	}

	// 应用更新
	progress("正在应用更新...", 0.8, "", "", 0, 0, 0, false)
	if err := s.applyUpdate(tempFile, targetFile, progress); err != nil {
		return err
	}

	// 清理 build 目录
	progress("正在清理构建目录...", 0.95, "", "", 0, 0, 0, false)
	s.cleanBuild()

	return nil
}

// applyUpdate 应用更新
func (s *SchemeUpdater) applyUpdate(temp, target string, progress types.ProgressFunc) error {
	// 终止进程（组合更新时跳过）
	if !s.SkipTerminate {
		progress("正在终止相关进程...", 0.85, "", "", 0, 0, 0, false)
		if err := s.TerminateProcesses(); err != nil {
			return fmt.Errorf("终止进程失败: %w", err)
		}
	}

	// 解压文件到主引擎目录
	progress("正在解压方案文件...", 0.9, "", "", 0, 0, 0, false)
	if err := s.ExtractZip(temp, s.Config.GetExtractPath()); err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	// 处理 CNB 镜像的嵌套目录问题
	if s.Config.Config.UseMirror {
		if err := fileutil.HandleCNBNestedDir(s.Config.GetExtractPath(), s.Config.Config.SchemeFile); err != nil {
			return fmt.Errorf("处理嵌套目录失败: %w", err)
		}
	}

	// 同步到其他引擎目录
	if len(s.Config.Config.InstalledEngines) > 1 {
		progress("正在同步到其他引擎...", 0.92, "", "", 0, 0, 0, false)
		if err := s.syncToOtherEngines(); err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("同步到其他引擎失败: %v", err), 0.92, "", "", 0, 0, 0, false)
		}
	}

	// 重命名临时文件
	progress("正在保存文件...", 0.93, "", "", 0, 0, 0, false)
	if fileutil.FileExists(target) {
		os.Remove(target)
	}
	if err := fileutil.MoveFile(temp, target); err != nil {
		return fmt.Errorf("重命名失败: %w", err)
	}

	// 保存记录
	recordPath := s.Config.GetSchemeRecordPath()
	if err := s.SaveRecord(recordPath, "scheme_file", s.Config.Config.SchemeFile, s.UpdateInfo); err != nil {
		return err
	}

	// 执行更新后 hook（失败不影响更新结果）
	if s.Config.Config.PostUpdateHook != "" {
		progress("执行更新后 hook...", 1.0, "", "", 0, 0, 0, false)
		if err := s.Config.ExecutePostUpdateHook(); err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("post-update hook 失败: %v", err), 1.0, "", "", 0, 0, 0, false)
		}
	}

	// 同步到 fcitx 目录（如果启用）
	if s.Config.Config.FcitxCompat {
		_, _, err := s.Config.SyncToFcitxDir()
		if err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("fcitx 同步失败: %v", err), 1.0, "", "", 0, 0, 0, false)
		}
	}

	progress("更新完成！", 1.0, "", "", 0, 0, 0, false)
	return nil
}

// syncToOtherEngines 同步文件到其他引擎目录
func (s *SchemeUpdater) syncToOtherEngines() error {
	// 检查是否配置了要更新的引擎列表
	updateEngines := s.Config.Config.UpdateEngines
	if len(updateEngines) == 0 {
		// 未配置：默认更新所有已安装的引擎
		updateEngines = s.Config.Config.InstalledEngines
	}

	// 如果只有一个引擎需要更新，跳过同步
	if len(updateEngines) <= 1 {
		return nil
	}

	primaryEngine := s.Config.Config.PrimaryEngine
	if primaryEngine == "" && len(updateEngines) > 0 {
		primaryEngine = updateEngines[0]
	}

	sourceDir := s.Config.GetExtractPath()
	var errors []string

	// 遍历用户选择要更新的引擎
	for _, engine := range updateEngines {
		// 跳过主引擎（已经解压到主引擎目录）
		if engine == primaryEngine {
			continue
		}

		// 获取目标引擎的数据目录
		targetDir := config.GetEngineDataDir(engine)
		if targetDir == "" {
			errors = append(errors, fmt.Sprintf("无法获取引擎 %s 的数据目录", engine))
			continue
		}

		// 确保目标目录存在
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("创建引擎 %s 目录失败: %v", engine, err))
			continue
		}

		// 复制文件（排除 build 目录和用户配置）
		if err := fileutil.SyncDirectory(sourceDir, targetDir, s.Config.Config.ExcludeFiles); err != nil {
			errors = append(errors, fmt.Sprintf("同步到引擎 %s 失败: %v", engine, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分引擎同步失败: %v", errors)
	}

	return nil
}

// cleanBuild 清理 build 目录
func (s *SchemeUpdater) cleanBuild() {
	buildDir := filepath.Join(s.Config.GetExtractPath(), "build")
	if fileutil.FileExists(buildDir) {
		os.RemoveAll(buildDir)
	}
}
