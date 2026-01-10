//go:build windows

package detector

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// checkRimeInstallation 检测 Windows 上的 Rime 输入法是否已安装
// 检查所有支持的引擎，只要有一个安装即认为已安装
func checkRimeInstallation() InstallationStatus {
	// 检查小狼毫 (Weasel)
	weaselInstalled := checkWeaselInstalled()

	// 检查玉兔毫 (Rabbit)
	rabbitInstalled := checkRabbitInstalled()

	// 只要有一个引擎安装就认为已安装
	if weaselInstalled || rabbitInstalled {
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 全部未安装，返回安装建议
	message := "⚠️  未检测到任何 Rime 输入法引擎\n\n" +
		"支持的引擎：\n" +
		"  • 小狼毫 (Weasel) - 官方版本\n" +
		"    - 官网下载: https://rime.im\n" +
		"    - Scoop: scoop install weasel\n" +
		"    - Chocolatey: choco install weasel\n" +
		"  • 玉兔毫 (Rabbit) - 社区版本\n" +
		"    - GitHub: https://github.com/amorphobia/rabbit\n\n" +
		"提示：程序仍可正常运行，但需要先安装至少一个引擎才能使用更新的配置。"

	return InstallationStatus{
		Installed: false,
		Message:   message,
	}
}

// checkWeaselInstalled 检查小狼毫是否已安装
func checkWeaselInstalled() bool {
	registryKeys := []string{
		`Software\Rime\Weasel`,             // CURRENT_USER
		`SOFTWARE\WOW6432Node\Rime\Weasel`, // LOCAL_MACHINE 64位
		`SOFTWARE\Rime\Weasel`,             // LOCAL_MACHINE 32位
	}

	// 检查 CURRENT_USER 注册表
	for _, keyPath := range registryKeys {
		k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
		if err == nil {
			k.Close()
			return true
		}
	}

	// 检查 LOCAL_MACHINE 注册表
	for _, keyPath := range registryKeys {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
		if err == nil {
			k.Close()
			return true
		}
	}

	// 检查默认的用户数据目录
	rimeDir := filepath.Join(os.Getenv("APPDATA"), "Rime")
	if _, err := os.Stat(rimeDir); err == nil {
		return true
	}

	return false
}

// checkRabbitInstalled 检查玉兔毫是否已安装
func checkRabbitInstalled() bool {
	registryKeys := []string{
		`Software\Rime\Rabbit`,             // CURRENT_USER
		`SOFTWARE\WOW6432Node\Rime\Rabbit`, // LOCAL_MACHINE 64位
		`SOFTWARE\Rime\Rabbit`,             // LOCAL_MACHINE 32位
	}

	// 方法1: 检查 CURRENT_USER 注册表
	for _, keyPath := range registryKeys {
		k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
		if err == nil {
			k.Close()
			return true
		}
	}

	// 方法2: 检查 LOCAL_MACHINE 注册表
	for _, keyPath := range registryKeys {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
		if err == nil {
			k.Close()
			return true
		}
	}

	// 方法3: 检查数据目录
	rabbitDir := filepath.Join(os.Getenv("APPDATA"), "Rabbit")
	if info, err := os.Stat(rabbitDir); err == nil && info.IsDir() {
		return true
	}

	// 方法4: 检查常见的安装目录
	possibleInstallDirs := []string{
		filepath.Join(os.Getenv("ProgramFiles"), "Rabbit"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Rabbit"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Rabbit"),
		filepath.Join(os.Getenv("USERPROFILE"), "Rabbit"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Rabbit"),
	}

	for _, dir := range possibleInstallDirs {
		// 检查是否存在 Rabbit.exe
		exePath := filepath.Join(dir, "Rabbit.exe")
		if _, err := os.Stat(exePath); err == nil {
			return true
		}
		// 检查是否存在 RabbitDeployer.exe
		deployerPath := filepath.Join(dir, "RabbitDeployer.exe")
		if _, err := os.Stat(deployerPath); err == nil {
			return true
		}
		// 检查是否存在 Rabbit.ahk (便携版)
		ahkPath := filepath.Join(dir, "Rabbit.ahk")
		if _, err := os.Stat(ahkPath); err == nil {
			return true
		}
	}

	// 方法5: 检查用户数据目录中的特征文件（玉兔毫特有）
	rabbitCustomYaml := filepath.Join(rabbitDir, "rabbit.custom.yaml")
	if _, err := os.Stat(rabbitCustomYaml); err == nil {
		return true
	}

	// 检查 default.custom.yaml (Rime 通用，但如果 Rabbit 目录存在就很可能是玉兔毫)
	defaultCustomYaml := filepath.Join(rabbitDir, "default.custom.yaml")
	if _, err := os.Stat(defaultCustomYaml); err == nil {
		// 再检查是否有 Rabbit 相关的其他文件
		rabbitYaml := filepath.Join(rabbitDir, "rabbit.yaml")
		if _, err := os.Stat(rabbitYaml); err == nil {
			return true
		}
	}

	return false
}
