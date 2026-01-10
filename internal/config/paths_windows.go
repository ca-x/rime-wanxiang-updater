//go:build windows

package config

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
	"rime-wanxiang-updater/internal/types"
)

// EngineInfo 引擎信息
type EngineInfo struct {
	Name         string
	RegistryKeys []string // 注册表键路径
	DataDir      string   // 数据目录（相对于 APPDATA）
}

// Windows 引擎定义
var windowsEngines = []EngineInfo{
	{
		Name: "小狼毫",
		RegistryKeys: []string{
			`Software\Rime\Weasel`,                  // CURRENT_USER
			`SOFTWARE\WOW6432Node\Rime\Weasel`,      // LOCAL_MACHINE 64位
			`SOFTWARE\Rime\Weasel`,                  // LOCAL_MACHINE 32位
		},
		DataDir: "Rime",
	},
	{
		Name: "玉兔毫",
		RegistryKeys: []string{
			`Software\Rime\Rabbit`,                  // CURRENT_USER
			`SOFTWARE\WOW6432Node\Rime\Rabbit`,      // LOCAL_MACHINE 64位
			`SOFTWARE\Rime\Rabbit`,                  // LOCAL_MACHINE 32位
		},
		DataDir: "Rabbit",
	},
}

// DetectInstalledEngines 检测已安装的引擎
func DetectInstalledEngines() []string {
	var installed []string

	for _, engine := range windowsEngines {
		if isEngineInstalled(engine) {
			installed = append(installed, engine.Name)
		}
	}

	// 如果没有检测到任何引擎，返回小狼毫作为默认
	if len(installed) == 0 {
		return []string{"小狼毫"}
	}

	return installed
}

// isEngineInstalled 检查引擎是否已安装
func isEngineInstalled(engine EngineInfo) bool {
	// 方法1: 检查注册表 CURRENT_USER
	for _, keyPath := range engine.RegistryKeys {
		k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
		if err == nil {
			k.Close()
			return true
		}
	}

	// 方法2: 检查注册表 LOCAL_MACHINE
	for _, keyPath := range engine.RegistryKeys {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
		if err == nil {
			k.Close()
			return true
		}
	}

	// 方法3: 检查数据目录
	dataDir := filepath.Join(os.Getenv("APPDATA"), engine.DataDir)
	if info, err := os.Stat(dataDir); err == nil && info.IsDir() {
		return true
	}

	// 方法4: 对于特定引擎的额外检查
	if engine.Name == "玉兔毫" {
		// 检查常见的安装目录
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

		// 检查特征配置文件（玉兔毫特有）
		rabbitCustomYaml := filepath.Join(dataDir, "rabbit.custom.yaml")
		if _, err := os.Stat(rabbitCustomYaml); err == nil {
			return true
		}

		// 检查 default.custom.yaml + rabbit.yaml 组合
		defaultCustomYaml := filepath.Join(dataDir, "default.custom.yaml")
		if _, err := os.Stat(defaultCustomYaml); err == nil {
			rabbitYaml := filepath.Join(dataDir, "rabbit.yaml")
			if _, err := os.Stat(rabbitYaml); err == nil {
				return true
			}
		}
	}

	return false
}

// getRimeUserDir 获取 Rime 用户目录
func getRimeUserDir(config *types.Config) string {
	// 使用主引擎
	engine := config.PrimaryEngine
	if engine == "" && len(config.InstalledEngines) > 0 {
		engine = config.InstalledEngines[0]
	}
	if engine == "" {
		engine = config.Engine
	}
	if engine == "" {
		engine = "小狼毫"
	}

	// 查找引擎的数据目录
	for _, engineInfo := range windowsEngines {
		if engineInfo.Name == engine {
			// 尝试从注册表读取自定义路径
			if customPath := getEngineCustomPath(engineInfo); customPath != "" {
				return customPath
			}
			// 使用默认路径
			return filepath.Join(os.Getenv("APPDATA"), engineInfo.DataDir)
		}
	}

	// 默认返回小狼毫路径
	// 尝试从注册表读取自定义路径
	if customPath := getWindowsRimeDir(); customPath != "" {
		return customPath
	}
	return filepath.Join(os.Getenv("APPDATA"), "Rime")
}

// getEngineCustomPath 从注册表获取引擎的自定义路径
func getEngineCustomPath(engine EngineInfo) string {
	// 尝试从 CURRENT_USER 读取
	for _, keyPath := range engine.RegistryKeys {
		k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
		if err == nil {
			defer k.Close()
			if path, _, err := k.GetStringValue("RimeUserDir"); err == nil && path != "" {
				return path
			}
		}
	}

	// 尝试从 LOCAL_MACHINE 读取
	for _, keyPath := range engine.RegistryKeys {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
		if err == nil {
			defer k.Close()
			if path, _, err := k.GetStringValue("RimeUserDir"); err == nil && path != "" {
				return path
			}
		}
	}

	return ""
}

// getWindowsRimeDir 从注册表获取 Windows Rime 目录（向后兼容）
func getWindowsRimeDir() string {
	// 尝试读取用户自定义路径
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Rime\Weasel`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		if path, _, err := k.GetStringValue("RimeUserDir"); err == nil && path != "" {
			return path
		}
	}

	return ""
}

// GetEngineDataDir 获取指定引擎的数据目录
func GetEngineDataDir(engineName string) string {
	for _, engineInfo := range windowsEngines {
		if engineInfo.Name == engineName {
			// 尝试从注册表读取自定义路径
			if customPath := getEngineCustomPath(engineInfo); customPath != "" {
				return customPath
			}
			// 使用默认路径
			return filepath.Join(os.Getenv("APPDATA"), engineInfo.DataDir)
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

	// 读取 RimeUserDir
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Rime\Weasel`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		if path, _, err := k.GetStringValue("RimeUserDir"); err == nil {
			detected["rime_user_dir"] = path
		}
	}

	// 读取 WeaselRoot
	k2, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Rime\Weasel`, registry.QUERY_VALUE)
	if err == nil {
		defer k2.Close()
		if root, _, err := k2.GetStringValue("WeaselRoot"); err == nil {
			detected["weasel_root"] = root
		}
		if server, _, err := k2.GetStringValue("ServerExecutable"); err == nil {
			if root, ok := detected["weasel_root"]; ok {
				detected["server_exe"] = filepath.Join(root, server)
			}
		}
	}

	// 如果没有找到，使用默认路径
	if detected["rime_user_dir"] == "" {
		detected["rime_user_dir"] = filepath.Join(os.Getenv("APPDATA"), "Rime")
	}

	return detected
}
