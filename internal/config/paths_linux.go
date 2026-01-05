//go:build linux

package config

import (
	"os"
	"path/filepath"

	"rime-wanxiang-updater/internal/types"
)

// getRimeUserDir 获取 Rime 用户目录
func getRimeUserDir(config *types.Config) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")
}

// DetectInstallationPaths 检测安装路径
func DetectInstallationPaths() map[string]string {
	detected := make(map[string]string)
	homeDir, _ := os.UserHomeDir()

	detected["rime_user_dir"] = filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")

	return detected
}
