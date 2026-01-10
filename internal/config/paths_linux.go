//go:build linux

package config

import (
	"os"
	"path/filepath"

	"rime-wanxiang-updater/internal/types"
)

// EngineInfo 引擎信息
type EngineInfo struct {
	Name     string
	DataDirs []string // 可能的数据目录（按优先级排序）
}

// Linux 引擎定义
var linuxEngines = []EngineInfo{
	{
		Name: "fcitx5",
		DataDirs: []string{
			".local/share/fcitx5/rime",
			".config/fcitx5/rime",
		},
	},
	{
		Name: "ibus",
		DataDirs: []string{
			".config/ibus/rime",
		},
	},
	{
		Name: "fcitx",
		DataDirs: []string{
			".config/fcitx/rime",
		},
	},
}

// DetectInstalledEngines 检测已安装的引擎
// 通过检测数据目录是否存在来判断
func DetectInstalledEngines() []string {
	homeDir, _ := os.UserHomeDir()
	var installed []string

	for _, engine := range linuxEngines {
		for _, dataDir := range engine.DataDirs {
			fullPath := filepath.Join(homeDir, dataDir)
			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				installed = append(installed, engine.Name)
				break // 找到一个目录就行，不需要检查其他目录
			}
		}
	}

	// 如果没有检测到任何引擎，返回 fcitx5 作为默认
	if len(installed) == 0 {
		return []string{"fcitx5"}
	}

	return installed
}

// getRimeUserDir 获取 Rime 用户目录
func getRimeUserDir(config *types.Config) string {
	homeDir, _ := os.UserHomeDir()

	// 使用主引擎
	engine := config.PrimaryEngine
	if engine == "" && len(config.InstalledEngines) > 0 {
		// 向后兼容：如果没有主引擎但有已安装列表，使用第一个
		engine = config.InstalledEngines[0]
	}
	if engine == "" {
		// 向后兼容：使用旧的 Engine 字段
		engine = config.Engine
	}

	// 查找引擎的数据目录
	for _, engineInfo := range linuxEngines {
		if engineInfo.Name == engine {
			// 查找第一个存在的目录
			for _, dataDir := range engineInfo.DataDirs {
				fullPath := filepath.Join(homeDir, dataDir)
				if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
					return fullPath
				}
			}
			// 如果都不存在，返回第一个作为默认
			return filepath.Join(homeDir, engineInfo.DataDirs[0])
		}
	}

	// 默认返回 fcitx5 路径
	return filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")
}

// GetEngineDataDir 获取指定引擎的数据目录
func GetEngineDataDir(engineName string) string {
	homeDir, _ := os.UserHomeDir()

	for _, engineInfo := range linuxEngines {
		if engineInfo.Name == engineName {
			// 返回第一个存在的目录
			for _, dataDir := range engineInfo.DataDirs {
				fullPath := filepath.Join(homeDir, dataDir)
				if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
					return fullPath
				}
			}
			// 如果都不存在，返回第一个作为默认
			return filepath.Join(homeDir, engineInfo.DataDirs[0])
		}
	}

	return ""
}

// GetAllEngineDataDirs 获取所有已安装引擎的数据目录
func GetAllEngineDataDirs(installedEngines []string) map[string]string {
	dirs := make(map[string]string)

	for _, engineName := range installedEngines {
		dir := GetEngineDataDir(engineName)
		if dir != "" {
			dirs[engineName] = dir
		}
	}

	return dirs
}

// DetectInstallationPaths 检测安装路径
func DetectInstallationPaths() map[string]string {
	detected := make(map[string]string)
	homeDir, _ := os.UserHomeDir()

	// 按优先级检测所有可能的路径
	candidates := []string{
		filepath.Join(homeDir, ".local/share/fcitx5/rime"),
		filepath.Join(homeDir, ".config/fcitx5/rime"),
		filepath.Join(homeDir, ".config/ibus/rime"),
		filepath.Join(homeDir, ".config/fcitx/rime"),
	}

	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			detected["rime_user_dir"] = path
			return detected
		}
	}

	// 默认使用 fcitx5 路径
	detected["rime_user_dir"] = filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")

	return detected
}
