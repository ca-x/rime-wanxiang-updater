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

// CheckUpdate 检查更新
func (d *DictUpdater) CheckUpdate() (*types.UpdateInfo, error) {
	var releases []types.GitHubRelease
	var err error

	if d.Config.Config.UseMirror {
		releases, err = d.APIClient.FetchCNBReleases(types.OWNER, types.CNB_REPO, types.DICT_TAG)
	} else {
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
func (d *DictUpdater) Run() error {
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
	if d.UpdateInfo.SHA256 != "" && d.CompareHash(d.UpdateInfo.SHA256, targetFile) {
		d.SaveRecord(recordPath, "dict_file", d.Config.Config.DictFile, d.UpdateInfo)
		return nil
	}

	// 下载文件
	tempFile := filepath.Join(d.Config.CacheDir, fmt.Sprintf("temp_dict_%s.zip", d.UpdateInfo.SHA256))
	if err := d.DownloadFile(d.UpdateInfo.URL, tempFile); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// 清理旧文件
	if fileutil.FileExists(targetFile) {
		d.CleanOldFiles(targetFile, tempFile, d.Config.GetDictExtractPath(), true)
	}

	// 应用更新
	return d.applyUpdate(tempFile, targetFile)
}

// applyUpdate 应用更新
func (d *DictUpdater) applyUpdate(temp, target string) error {
	// 终止进程
	if err := d.TerminateProcesses(); err != nil {
		return fmt.Errorf("终止进程失败: %w", err)
	}

	// 确保词库目录存在
	dictDir := d.Config.GetDictExtractPath()
	os.MkdirAll(dictDir, 0755)

	// 解压文件
	if err := d.ExtractZip(temp, dictDir); err != nil {
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
	recordPath := d.Config.GetDictRecordPath()
	return d.SaveRecord(recordPath, "dict_file", d.Config.Config.DictFile, d.UpdateInfo)
}
