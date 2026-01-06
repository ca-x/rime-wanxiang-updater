package updater

import (
	"fmt"
	"os"
	"path/filepath"

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

	// 获取本地版本信息
	recordPath := s.Config.GetSchemeRecordPath()
	localRecord := s.GetLocalRecord(recordPath)

	status := &types.UpdateStatus{
		RemoteVersion: remoteInfo.Tag,
		RemoteTime:    remoteInfo.UpdateTime,
		NeedsUpdate:   true,
	}

	if localRecord != nil {
		status.LocalVersion = localRecord.Tag
		status.LocalTime = localRecord.UpdateTime
		status.NeedsUpdate = remoteInfo.UpdateTime.After(localRecord.UpdateTime)

		if status.NeedsUpdate {
			status.Message = fmt.Sprintf("发现新版本: %s → %s", localRecord.Tag, remoteInfo.Tag)
		} else {
			status.Message = "已是最新版本"
		}
	} else {
		status.LocalVersion = "未安装"
		status.Message = fmt.Sprintf("检测到可用版本: %s", remoteInfo.Tag)
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
	if s.UpdateInfo.SHA256 != "" && s.CompareHash(s.UpdateInfo.SHA256, targetFile) {
		progress("本地文件已是最新版本", 1.0, "", "", 0, 0, 0, false)
		s.SaveRecord(recordPath, "scheme_file", s.Config.Config.SchemeFile, s.UpdateInfo)
		return nil
	}

	// 下载文件
	progress(fmt.Sprintf("准备从 %s 下载方案...", source), 0.15, "", "", 0, 0, 0, false)
	tempFile := filepath.Join(s.Config.CacheDir, fmt.Sprintf("temp_scheme_%s.zip", s.UpdateInfo.SHA256))
	if err := s.DownloadFileWithValidation(s.UpdateInfo.URL, tempFile, s.Config.Config.SchemeFile, source, s.UpdateInfo.Size, progress); err != nil {
		return fmt.Errorf("下载失败: %w", err)
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
	// 终止进程
	progress("正在终止相关进程...", 0.85, "", "", 0, 0, 0, false)
	if err := s.TerminateProcesses(); err != nil {
		return fmt.Errorf("终止进程失败: %w", err)
	}

	// 解压文件
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
	progress("方案更新完成！", 1.0, "", "", 0, 0, 0, false)
	return s.SaveRecord(recordPath, "scheme_file", s.Config.Config.SchemeFile, s.UpdateInfo)
}

// cleanBuild 清理 build 目录
func (s *SchemeUpdater) cleanBuild() {
	buildDir := filepath.Join(s.Config.GetExtractPath(), "build")
	if fileutil.FileExists(buildDir) {
		os.RemoveAll(buildDir)
	}
}
