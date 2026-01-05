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
func (s *SchemeUpdater) Run() error {
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
	if s.UpdateInfo.SHA256 != "" && s.CompareHash(s.UpdateInfo.SHA256, targetFile) {
		s.SaveRecord(recordPath, "scheme_file", s.Config.Config.SchemeFile, s.UpdateInfo)
		return nil
	}

	// 下载文件
	tempFile := filepath.Join(s.Config.CacheDir, fmt.Sprintf("temp_scheme_%s.zip", s.UpdateInfo.SHA256))
	if err := s.DownloadFile(s.UpdateInfo.URL, tempFile); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// 清理旧文件
	if fileutil.FileExists(targetFile) {
		s.CleanOldFiles(targetFile, tempFile, s.Config.GetExtractPath(), false)
	}

	// 应用更新
	if err := s.applyUpdate(tempFile, targetFile); err != nil {
		return err
	}

	// 清理 build 目录
	s.cleanBuild()

	return nil
}

// applyUpdate 应用更新
func (s *SchemeUpdater) applyUpdate(temp, target string) error {
	// 终止进程
	if err := s.TerminateProcesses(); err != nil {
		return fmt.Errorf("终止进程失败: %w", err)
	}

	// 解压文件
	if err := s.ExtractZip(temp, s.Config.GetExtractPath()); err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	// 重命名临时文件
	if fileutil.FileExists(target) {
		os.Remove(target)
	}
	if err := os.Rename(temp, target); err != nil {
		return fmt.Errorf("重命名失败: %w", err)
	}

	// 保存记录
	recordPath := s.Config.GetSchemeRecordPath()
	return s.SaveRecord(recordPath, "scheme_file", s.Config.Config.SchemeFile, s.UpdateInfo)
}

// cleanBuild 清理 build 目录
func (s *SchemeUpdater) cleanBuild() {
	buildDir := filepath.Join(s.Config.GetExtractPath(), "build")
	if fileutil.FileExists(buildDir) {
		os.RemoveAll(buildDir)
	}
}
