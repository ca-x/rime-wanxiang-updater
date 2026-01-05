//go:build darwin

package config

import (
	"os"
	"path/filepath"

	"rime-wanxiang-updater/internal/types"
)

// getRimeUserDir 获取 Rime 用户目录
func getRimeUserDir(config *types.Config) string {
	homeDir, _ := os.UserHomeDir()

	if config.Engine == "鼠须管" {
		return filepath.Join(homeDir, "Library", "Rime")
	}

	// 小企鹅或其他
	return filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")
}

// DetectInstallationPaths 检测安装路径
func DetectInstallationPaths(engine string) map[string]string {
	detected := make(map[string]string)
	homeDir, _ := os.UserHomeDir()

	if engine == "鼠须管" {
		detected["rime_user_dir"] = filepath.Join(homeDir, "Library", "Rime")
	} else {
		detected["rime_user_dir"] = filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")
	}

	return detected
}
