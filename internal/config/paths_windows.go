//go:build windows

package config

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
	"rime-wanxiang-updater/internal/types"
)

// getRimeUserDir 获取 Rime 用户目录
func getRimeUserDir(config *types.Config) string {
	// 尝试从注册表读取自定义路径
	if customPath := getWindowsRimeDir(); customPath != "" {
		return customPath
	}
	// 使用默认路径
	return filepath.Join(os.Getenv("APPDATA"), "Rime")
}

// getWindowsRimeDir 从注册表获取 Windows Rime 目录
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
