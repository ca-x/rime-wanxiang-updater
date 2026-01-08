//go:build windows

package detector

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// checkRimeInstallation 检测 Windows 上的 Rime（小狼毫）是否已安装
func checkRimeInstallation() InstallationStatus {
	// 方法1: 检查注册表
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Rime\Weasel`, registry.QUERY_VALUE)
	if err == nil {
		k.Close()
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 方法2: 检查注册表（本地机器）
	k2, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Rime\Weasel`, registry.QUERY_VALUE)
	if err == nil {
		k2.Close()
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 方法3: 检查默认的用户数据目录
	rimeDir := filepath.Join(os.Getenv("APPDATA"), "Rime")
	if _, err := os.Stat(rimeDir); err == nil {
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 未安装，返回安装建议
	message := "⚠️  未检测到 Rime 输入法（小狼毫）\n\n" +
		"安装方式：\n" +
		"  • 从官网下载安装程序: https://rime.im\n" +
		"  • 使用 Scoop: scoop install weasel\n\n" +
		"提示：程序仍可正常运行，但需要先安装 Rime 才能使用更新的配置。"

	return InstallationStatus{
		Installed: false,
		Message:   message,
	}
}
