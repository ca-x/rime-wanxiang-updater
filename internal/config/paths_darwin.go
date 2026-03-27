//go:build darwin

package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"rime-wanxiang-updater/internal/types"
)

// EngineInfo 引擎信息
type EngineInfo struct {
	Name     string
	AppPath  string
	DataDir  string
	BundleID string // macOS 输入法 Bundle ID
}

// macOS 引擎定义
var macOSEngines = map[string]EngineInfo{
	"鼠须管": {
		Name:     "鼠须管",
		AppPath:  "/Library/Input Methods/Squirrel.app",
		DataDir:  "Library/Rime",
		BundleID: "im.rime.inputmethod.Squirrel",
	},
	"小企鹅": {
		Name:     "小企鹅",
		AppPath:  "/Library/Input Methods/Fcitx5.app",
		DataDir:  ".local/share/fcitx5/rime",
		BundleID: "org.fcitx.inputmethod.Fcitx5",
	},
}

// isAppValid 检查应用是否有效（存在 Info.plist）
func isAppValid(appPath string) bool {
	// 检查 Contents/Info.plist 是否存在
	// 这是 macOS 应用的标准文件，用于验证应用完整性
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if _, err := os.Stat(plistPath); err == nil {
		return true
	}
	return false
}

// isInputMethodRegistered 检查输入法是否已注册到系统
// 通过检查 AppleEnabledInputSources 和 AppleInputSourceHistory
func isInputMethodRegistered(bundleID string) bool {
	// 检查启用的输入法
	if checkInputSourceList("AppleEnabledInputSources", bundleID) {
		return true
	}
	// 检查输入法历史（用户可能禁用了但仍安装）
	if checkInputSourceList("AppleInputSourceHistory", bundleID) {
		return true
	}
	return false
}

// checkInputSourceList 检查指定的输入法列表是否包含 bundleID
func checkInputSourceList(domain, bundleID string) bool {
	cmd := exec.Command("defaults", "read", "com.apple.HIToolbox", domain)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	// 在输出中搜索 Bundle ID
	return strings.Contains(string(output), bundleID)
}

// DetectInstalledEngines 检测已安装的引擎
// 同时检查应用是否存在和是否已注册到系统
func DetectInstalledEngines() []string {
	homeDir, _ := os.UserHomeDir()
	var installed []string

	for _, engine := range macOSEngines {
		appValid := false

		// 检查系统目录
		if isAppValid(engine.AppPath) {
			appValid = true
		} else {
			// 检查用户目录
			userAppPath := filepath.Join(homeDir, "Library", "Input Methods", filepath.Base(engine.AppPath))
			if isAppValid(userAppPath) {
				appValid = true
			}
		}

		// 应用存在且已注册到系统才算真正安装
		if appValid && isInputMethodRegistered(engine.BundleID) {
			installed = append(installed, engine.Name)
		}
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
		return ""
	}

	if info, ok := macOSEngines[engine]; ok {
		return filepath.Join(homeDir, info.DataDir)
	}

	return ""
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
