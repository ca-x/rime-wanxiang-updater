//go:build darwin

package detector

import "os"

// checkRimeInstallation 检测 macOS 上的 Rime（鼠须管）是否已安装
func checkRimeInstallation() InstallationStatus {
	// 检查鼠须管是否安装
	squirrelAppPath := "/Library/Input Methods/Squirrel.app"
	if _, err := os.Stat(squirrelAppPath); err == nil {
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 未安装，返回安装建议
	message := "⚠️  未检测到 Rime 输入法（鼠须管）\n\n" +
		"安装方式：\n" +
		"  1. 使用 Homebrew: brew install --cask squirrel\n" +
		"  2. 从官网下载: https://rime.im\n\n" +
		"提示：程序仍可正常运行，但需要先安装 Rime 才能使用更新的配置。"

	return InstallationStatus{
		Installed: false,
		Message:   message,
	}
}
