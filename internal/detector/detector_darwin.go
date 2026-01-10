//go:build darwin

package detector

import "os"

// checkRimeInstallation 检测 macOS 上的 Rime 输入法是否已安装
// 检查所有支持的引擎，只要有一个安装即认为已安装
func checkRimeInstallation() InstallationStatus {
	homeDir, _ := os.UserHomeDir()

	// 检查鼠须管（同时检查系统目录和用户目录）
	squirrelSystemPath := "/Library/Input Methods/Squirrel.app"
	squirrelUserPath := homeDir + "/Library/Input Methods/Squirrel.app"
	squirrelInstalled := false
	if _, err := os.Stat(squirrelSystemPath); err == nil {
		squirrelInstalled = true
	} else if _, err := os.Stat(squirrelUserPath); err == nil {
		squirrelInstalled = true
	}

	// 检查小企鹅（同时检查系统目录和用户目录）
	fcitx5SystemPath := "/Library/Input Methods/Fcitx5.app"
	fcitx5UserPath := homeDir + "/Library/Input Methods/Fcitx5.app"
	fcitx5Installed := false
	if _, err := os.Stat(fcitx5SystemPath); err == nil {
		fcitx5Installed = true
	} else if _, err := os.Stat(fcitx5UserPath); err == nil {
		fcitx5Installed = true
	}

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
