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

// GetStatus 获取更新状态
func (m *ModelUpdater) GetStatus() (*types.UpdateStatus, error) {
	// 获取远程版本信息
	remoteInfo, err := m.CheckUpdate()
	if err != nil {
		return nil, err
	}

	// 获取本地版本信息
	recordPath := m.Config.GetModelRecordPath()
	localRecord := m.GetLocalRecord(recordPath)

	status := &types.UpdateStatus{
		RemoteVersion: remoteInfo.Tag,
		RemoteTime:    remoteInfo.UpdateTime,
		NeedsUpdate:   true,
	}

	if localRecord != nil {
		// 对于模型文件，根据 Tag 内容决定显示格式
		// CNB 镜像使用日期格式（如 "2025-01-08"）
		// GitHub 使用 Tag 名称（如 "LTS"）
		if localRecord.Tag != "" && localRecord.Tag != "model" {
			status.LocalVersion = localRecord.Tag
		} else {
			// 如果 Tag 为空或为旧的 "model"，使用日期格式
			status.LocalVersion = localRecord.UpdateTime.Format("2006-01-02")
		}

		status.LocalTime = localRecord.UpdateTime
		status.NeedsUpdate = remoteInfo.UpdateTime.After(localRecord.UpdateTime)

		if status.NeedsUpdate {
			remoteVersion := remoteInfo.Tag
			if remoteVersion == "" || remoteVersion == "model" {
				remoteVersion = remoteInfo.UpdateTime.Format("2006-01-02")
			}
			status.Message = fmt.Sprintf("发现新版本: %s → %s", status.LocalVersion, remoteVersion)
		} else {
			status.Message = fmt.Sprintf("已是最新版本 (当前版本: %s)", status.LocalVersion)
		}
	} else {
		status.LocalVersion = "未安装"
		status.Message = fmt.Sprintf("检测到可用模型: %s", remoteInfo.Tag)
	}

	return status, nil
}

// CheckUpdate 检查更新
func (m *ModelUpdater) CheckUpdate() (*types.UpdateInfo, error) {
	var releases []types.GitHubRelease
	var err error

	if m.Config.Config.UseMirror {
		// CNB 镜像：模型文件不在版本列表中，使用静态下载地址
		modelURL := fmt.Sprintf("https://cnb.cool/%s/%s/-/releases/download/model/%s", types.OWNER, types.CNB_REPO, types.MODEL_FILE)

		// 对于 CNB 镜像，使用 HTTP Last-Modified 作为版本标识
		// 如果获取失败，使用一个固定的时间戳，避免每次都认为有更新
		updateTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // 默认固定时间
		if resp, err := m.APIClient.Head(modelURL); err == nil {
			if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
				if t, err := time.Parse(time.RFC1123, lastModified); err == nil {
					updateTime = t
				}
			}
		}

		return &types.UpdateInfo{
			Name:       types.MODEL_FILE,
			URL:        modelURL,
			UpdateTime: updateTime,
			Tag:        updateTime.Format("2006-01-02"), // 使用日期作为 Tag
			Size:       0,                               // CNB 不提供文件大小信息
		}, nil
	} else {
		// GitHub：从 RIME-LMDG 仓库获取，使用 tag "LTS"
		releases, err = m.APIClient.FetchGitHubReleases(types.OWNER, types.MODEL_REPO, types.MODEL_TAG)
		if err != nil {
			return nil, fmt.Errorf("获取版本信息失败: %w", err)
		}

		if len(releases) == 0 {
			return nil, fmt.Errorf("未找到任何发布版本")
		}

		// 遍历所有 release 查找模型文件
		for i := len(releases) - 1; i >= 0; i-- {
			release := releases[i]

			for _, asset := range release.Assets {
				// 查找 .gram 文件（不包含 mini）
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
		}

		return nil, fmt.Errorf("未找到匹配的模型文件: %s", types.MODEL_FILE)
	}
}

// Run 执行更新
func (m *ModelUpdater) Run(progress types.ProgressFunc) error {
	if progress == nil {
		progress = func(string, float64, string, string, int64, int64, float64, bool) {} // 空函数避免 nil 检查
	}

	// 执行更新前 hook
	if m.Config.Config.PreUpdateHook != "" {
		progress("执行更新前 hook...", 0.02, "", "", 0, 0, 0, false)
		if err := m.Config.ExecutePreUpdateHook(); err != nil {
			return fmt.Errorf("pre-update hook 失败，已取消更新: %w", err)
		}
	}

	// 显示下载源
	source := "GitHub"
	if m.Config.Config.UseMirror {
		source = "CNB 镜像"
	}
	progress(fmt.Sprintf("正在检查模型更新 [%s]...", source), 0.05, "", "", 0, 0, 0, false)

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
	progress("正在校验本地文件...", 0.1, "", "", 0, 0, 0, false)

	// 优先使用 SHA256 校验（如果有）
	if m.UpdateInfo.SHA256 != "" && m.CompareHash(m.UpdateInfo.SHA256, targetPath) {
		progress("本地文件已是最新版本", 1.0, "", "", 0, 0, 0, false)
		m.SaveRecord(recordPath, "model_name", types.MODEL_FILE, m.UpdateInfo)
		return nil
	}

	// 如果没有 SHA256（如 CNB 镜像），检查文件是否存在且有本地记录
	if m.UpdateInfo.SHA256 == "" && fileutil.FileExists(targetPath) {
		localRecord := m.GetLocalRecord(recordPath)
		if localRecord != nil && localRecord.Name == types.MODEL_FILE {
			// 文件存在且有记录，认为已是最新版本（除非 UpdateTime 更新）
			if !m.UpdateInfo.UpdateTime.After(localRecord.UpdateTime) {
				progress("本地文件已是最新版本", 1.0, "", "", 0, 0, 0, false)
				return nil
			}
		}
	}

	// 下载文件
	progress(fmt.Sprintf("准备从 %s 下载模型...", source), 0.15, "", "", 0, 0, 0, false)
	tempFile := filepath.Join(m.Config.CacheDir, fmt.Sprintf("%s_%s.tmp", types.MODEL_FILE, m.UpdateInfo.SHA256))
	if err := m.DownloadFile(m.UpdateInfo.URL, tempFile, types.MODEL_FILE, source, progress); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// 应用更新
	progress("正在应用更新...", 0.8, "", "", 0, 0, 0, false)
	return m.applyUpdate(tempFile, targetPath, progress)
}

// applyUpdate 应用更新
func (m *ModelUpdater) applyUpdate(temp, target string, progress types.ProgressFunc) error {
	// 终止进程（组合更新时跳过）
	if !m.SkipTerminate {
		progress("正在终止相关进程...", 0.85, "", "", 0, 0, 0, false)
		if err := m.TerminateProcesses(); err != nil {
			return fmt.Errorf("终止进程失败: %w", err)
		}
	}

	// 覆盖目标文件
	progress("正在保存模型文件...", 0.9, "", "", 0, 0, 0, false)
	if fileutil.FileExists(target) {
		os.Remove(target)
	}
	if err := fileutil.MoveFile(temp, target); err != nil {
		return fmt.Errorf("替换文件失败: %w", err)
	}

	// 保存记录
	recordPath := m.Config.GetModelRecordPath()
	if err := m.SaveRecord(recordPath, "model_name", types.MODEL_FILE, m.UpdateInfo); err != nil {
		return err
	}

	// 执行更新后 hook（失败不影响更新结果）
	if m.Config.Config.PostUpdateHook != "" {
		progress("执行更新后 hook...", 1.0, "", "", 0, 0, 0, false)
		if err := m.Config.ExecutePostUpdateHook(); err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("post-update hook 失败: %v", err), 1.0, "", "", 0, 0, 0, false)
		}
	}

	// 同步到 fcitx 目录（如果启用）
	if m.Config.Config.FcitxCompat {
		_, _, err := m.Config.SyncToFcitxDir()
		if err != nil {
			// 只记录错误，不返回失败
			progress(fmt.Sprintf("fcitx 同步失败: %v", err), 1.0, "", "", 0, 0, 0, false)
		}
	}

	progress("更新完成！", 1.0, "", "", 0, 0, 0, false)
	return nil
}
