//go:build linux

package detector

import (
	"os"
	"os/exec"
	"path/filepath"
)

// checkRimeInstallation 检测 Linux 上的 Rime 是否已安装
func checkRimeInstallation() InstallationStatus {
	homeDir, _ := os.UserHomeDir()

	// 检查常见的 Rime 相关路径
	possiblePaths := []string{
		// fcitx5-rime
		filepath.Join(homeDir, ".local", "share", "fcitx5", "rime"),
		"/usr/share/rime-data",
		// ibus-rime
		filepath.Join(homeDir, ".config", "ibus", "rime"),
		// fcitx-rime (旧版)
		filepath.Join(homeDir, ".config", "fcitx", "rime"),
	}

	// 检查是否存在任意一个路径
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return InstallationStatus{
				Installed: true,
				Message:   "",
			}
		}
	}

	// 检查是否安装了 rime 相关的包
	// fcitx5-rime
	if _, err := exec.LookPath("fcitx5"); err == nil {
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// ibus-rime
	if _, err := exec.LookPath("ibus-daemon"); err == nil {
		if _, err := os.Stat("/usr/share/ibus/component/rime.xml"); err == nil {
			return InstallationStatus{
				Installed: true,
				Message:   "",
			}
		}
	}

	// 未安装，返回安装建议
	message := "⚠️  未检测到 Rime 输入法\n\n" +
		"安装方式：\n" +
		"  • Debian/Ubuntu: sudo apt install fcitx5-rime 或 ibus-rime\n" +
		"  • Fedora: sudo dnf install fcitx5-rime 或 ibus-rime\n" +
		"  • Arch Linux: sudo pacman -S fcitx5-rime 或 ibus-rime\n" +
		"  • 从官网下载: https://rime.im\n\n" +
		"提示：程序仍可正常运行，但需要先安装 Rime 才能使用更新的配置。"

	return InstallationStatus{
		Installed: false,
		Message:   message,
	}
}
