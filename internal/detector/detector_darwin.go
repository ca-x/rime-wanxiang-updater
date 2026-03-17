//go:build darwin

package detector

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Bundle ID 定义
const (
	squirrelBundleID = "im.rime.inputmethod.Squirrel"
	fcitx5BundleID   = "org.fcitx.inputmethod.Fcitx5"
)

// isAppValid 检查应用是否有效（存在 Info.plist）
func isAppValid(appPath string) bool {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if _, err := os.Stat(plistPath); err == nil {
		return true
	}
	return false
}

// isInputMethodRegistered 检查输入法是否已注册到系统
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
	return strings.Contains(string(output), bundleID)
}

// checkRimeInstallation 检测 macOS 上的 Rime 输入法是否已安装
// 同时检查应用是否存在和是否已注册到系统
func checkRimeInstallation() InstallationStatus {
	homeDir, _ := os.UserHomeDir()

	// 检查鼠须管
	squirrelAppValid := isAppValid("/Library/Input Methods/Squirrel.app") ||
		isAppValid(filepath.Join(homeDir, "Library", "Input Methods", "Squirrel.app"))
	squirrelInstalled := squirrelAppValid && isInputMethodRegistered(squirrelBundleID)

	// 检查小企鹅
	fcitx5AppValid := isAppValid("/Library/Input Methods/Fcitx5.app") ||
		isAppValid(filepath.Join(homeDir, "Library", "Input Methods", "Fcitx5.app"))
	fcitx5Installed := fcitx5AppValid && isInputMethodRegistered(fcitx5BundleID)

	// 只要有一个引擎安装就认为已安装
	if squirrelInstalled || fcitx5Installed {
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 全部未安装，返回安装建议
	message := "⚠️  未检测到任何 Rime 输入法引擎\n\n" +
		"支持的引擎：\n" +
		"  • 鼠须管 (Squirrel):\n" +
		"    brew install --cask squirrel\n\n" +
		"  • 小企鹅 (FCITX5):\n" +
		"    - 拼音版: brew install --cask tinypkg/tap/fcitx5-pinyin\n" +
		"    - 中州韵版: brew install --cask tinypkg/tap/fcitx5-rime\n" +
		"    - 原装版: brew install --cask tinypkg/tap/fcitx5\n" +
		"    - 安装器: https://github.com/fcitx-contrib/fcitx5-macos-installer\n\n" +
		"提示：程序仍可正常运行，但需要先安装至少一个引擎才能使用更新的配置。"

	return InstallationStatus{
		Installed: false,
		Message:   message,
	}
}