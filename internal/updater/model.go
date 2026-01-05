package updater

import (
	"fmt"
	"os"
	"path/filepath"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/fileutil"
	"rime-wanxiang-updater/internal/types"
)

// ModelUpdater 模型更新器
type ModelUpdater struct {
	*BaseUpdater
	UpdateInfo *types.UpdateInfo
}

// NewModelUpdater 创建模型更新器
func NewModelUpdater(cfg *config.Manager) *ModelUpdater {
	return &ModelUpdater{
		BaseUpdater: NewBaseUpdater(cfg),
	}
}

// CheckUpdate 检查更新
func (m *ModelUpdater) CheckUpdate() (*types.UpdateInfo, error) {
	var releases []types.GitHubRelease
	var err error

	if m.Config.Config.UseMirror {
		releases, err = m.APIClient.FetchCNBReleases(types.OWNER, types.CNB_REPO, types.MODEL_TAG)
	} else {
		releases, err = m.APIClient.FetchGitHubReleases(types.OWNER, types.MODEL_REPO, types.MODEL_TAG)
	}

	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("未找到任何发布版本")
	}

	// 获取最新的 release
	release := releases[len(releases)-1]
	for _, asset := range release.Assets {
		if asset.Name == types.MODEL_FILE {
			return &types.UpdateInfo{
				Name:       asset.Name,
				URL:        asset.BrowserDownloadURL,
				UpdateTime: asset.UpdatedAt,
				Tag:        release.TagName,
				Size:       asset.Size,
			}, nil
		}
	}

	return nil, fmt.Errorf("未找到匹配的模型文件: %s", types.MODEL_FILE)
}

// Run 执行更新
func (m *ModelUpdater) Run() error {
	if m.UpdateInfo == nil {
		info, err := m.CheckUpdate()
		if err != nil {
			return err
		}
		m.UpdateInfo = info
	}

	if m.UpdateInfo == nil {
		return fmt.Errorf("未找到模型更新")
	}

	recordPath := m.Config.GetModelRecordPath()
	targetPath := filepath.Join(m.Config.GetExtractPath(), types.MODEL_FILE)

	// 校验本地文件
	if m.UpdateInfo.SHA256 != "" && m.CompareHash(m.UpdateInfo.SHA256, targetPath) {
		m.SaveRecord(recordPath, "model_name", types.MODEL_FILE, m.UpdateInfo)
		return nil
	}

	// 下载文件
	tempFile := filepath.Join(m.Config.CacheDir, fmt.Sprintf("%s_%s.tmp", types.MODEL_FILE, m.UpdateInfo.SHA256))
	if err := m.DownloadFile(m.UpdateInfo.URL, tempFile); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// 应用更新
	return m.applyUpdate(tempFile, targetPath)
}

// applyUpdate 应用更新
func (m *ModelUpdater) applyUpdate(temp, target string) error {
	// 终止进程
	if err := m.TerminateProcesses(); err != nil {
		return fmt.Errorf("终止进程失败: %w", err)
	}

	// 覆盖目标文件
	if fileutil.FileExists(target) {
		os.Remove(target)
	}
	if err := os.Rename(temp, target); err != nil {
		return fmt.Errorf("替换文件失败: %w", err)
	}

	// 保存记录
	recordPath := m.Config.GetModelRecordPath()
	return m.SaveRecord(recordPath, "model_name", types.MODEL_FILE, m.UpdateInfo)
}
