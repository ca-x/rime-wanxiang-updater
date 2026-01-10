//go:build darwin

package config

import (
	"os"
	"path/filepath"

	"rime-wanxiang-updater/internal/types"
)

// EngineInfo 引擎信息
type EngineInfo struct {
	Name    string
	AppPath string
	DataDir string
}

// macOS 引擎定义
var macOSEngines = map[string]EngineInfo{
	"鼠须管": {
		Name:    "鼠须管",
		AppPath: "/Library/Input Methods/Squirrel.app",
		DataDir: "Library/Rime",
	},
	"小企鹅": {
		Name:    "小企鹅",
		// TODO: 确认小企鹅最终安装的应用名称
		// 安装器是 Fcitx5Installer.app，但最终应用可能是 Fcitx5.app
		AppPath: "/Library/Input Methods/Fcitx5.app",
		DataDir: ".local/share/fcitx5/rime",
	},
}

// DetectInstalledEngines 检测已安装的引擎
func DetectInstalledEngines() []string {
	homeDir, _ := os.UserHomeDir()
	var installed []string

	for _, engine := range macOSEngines {
		// 同时检查系统目录和用户目录
		systemPath := engine.AppPath
		userPath := filepath.Join(homeDir, engine.AppPath)

		if _, err := os.Stat(systemPath); err == nil {
			installed = append(installed, engine.Name)
		} else if _, err := os.Stat(userPath); err == nil {
			installed = append(installed, engine.Name)
		}
	}

	// 如果没有检测到任何引擎，返回鼠须管作为默认
	if len(installed) == 0 {
		return []string{"鼠须管"}
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
	if engine == "" {
		// 最后的默认值
		engine = "鼠须管"
	}

	if info, ok := macOSEngines[engine]; ok {
		return filepath.Join(homeDir, info.DataDir)
	}

	// 默认返回鼠须管路径
	return filepath.Join(homeDir, "Library", "Rime")
}

// GetEngineDataDir 获取指定引擎的数据目录
func GetEngineDataDir(engineName string) string {
	homeDir, _ := os.UserHomeDir()

	if info, ok := macOSEngines[engineName]; ok {
		return filepath.Join(homeDir, info.DataDir)
	}

	return ""
}

// GetAllEngineDataDirs 获取所有已安装引擎的数据目录
func GetAllEngineDataDirs(installedEngines []string) map[string]string {
	homeDir, _ := os.UserHomeDir()
	dirs := make(map[string]string)

	for _, engineName := range installedEngines {
		if info, ok := macOSEngines[engineName]; ok {
			dirs[engineName] = filepath.Join(homeDir, info.DataDir)
		}
	}

	return dirs
}

// DetectInstallationPaths 检测安装路径
func DetectInstallationPaths(engine string) map[string]string {
	detected := make(map[string]string)
	homeDir, _ := os.UserHomeDir()

	if info, ok := macOSEngines[engine]; ok {
		detected["rime_user_dir"] = filepath.Join(homeDir, info.DataDir)
		detected["app_path"] = info.AppPath
	} else {
		// 默认使用鼠须管
		detected["rime_user_dir"] = filepath.Join(homeDir, "Library", "Rime")
		detected["app_path"] = "/Library/Input Methods/Squirrel.app"
	}

	return detected
}
