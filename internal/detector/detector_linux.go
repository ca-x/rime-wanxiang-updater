//go:build linux

package detector

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Linux 引擎检测信息
type linuxEngineInfo struct {
	Name     string
	DataDirs []string // 数据目录（相对于用户目录）
	Commands []string // 可执行命令
}

var linuxEngines = []linuxEngineInfo{
	{
		Name: "fcitx5",
		DataDirs: []string{
			".local/share/fcitx5/rime",
			".config/fcitx5/rime",
		},
		Commands: []string{"fcitx5"},
	},
	{
		Name: "ibus",
		DataDirs: []string{
			".config/ibus/rime",
		},
		Commands: []string{"ibus-daemon"},
	},
	{
		Name: "fcitx",
		DataDirs: []string{
			".config/fcitx/rime",
		},
		Commands: []string{"fcitx"},
	},
}

// checkRimeInstallation 检测 Linux 上的 Rime 是否已安装
func checkRimeInstallation() InstallationStatus {
	homeDir, _ := os.UserHomeDir()
	installed := false

	// 检查各个引擎的数据目录
	for _, engine := range linuxEngines {
		for _, dataDir := range engine.DataDirs {
			fullPath := filepath.Join(homeDir, dataDir)
			// 使用 Lstat 避免跟随符号链接，确保检测到的是真实的引擎目录
			if info, err := os.Lstat(fullPath); err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
				installed = true
				break
			}
		}
		if installed {
			break
		}
	}

	// 检查公共的 rime 数据目录
	if !installed {
		if _, err := os.Stat("/usr/share/rime-data"); err == nil {
			installed = true
		}
	}

	// 检查是否安装了 rime 相关的可执行文件
	if !installed {
		for _, engine := range linuxEngines {
			for _, cmd := range engine.Commands {
				if _, err := exec.LookPath(cmd); err == nil {
					installed = true
					break
				}
			}
			if installed {
				break
			}
		}
	}

	// 特殊检查：ibus-rime
	if !installed {
		if _, err := exec.LookPath("ibus-daemon"); err == nil {
			if _, err := os.Stat("/usr/share/ibus/component/rime.xml"); err == nil {
				installed = true
			}
		}
	}

	if installed {
		return InstallationStatus{
			Installed: true,
			Message:   "",
		}
	}

	// 未安装，返回安装建议
	message := "⚠️  未检测到任何 Rime 输入法引擎\n\n" +
		"支持的引擎：\n" +
		"  • FCITX5-Rime (推荐)\n" +
		"    - Debian/Ubuntu: sudo apt install fcitx5-rime\n" +
		"    - Fedora: sudo dnf install fcitx5-rime\n" +
		"    - Arch Linux: sudo pacman -S fcitx5-rime\n" +
		"  • IBus-Rime\n" +
		"    - Debian/Ubuntu: sudo apt install ibus-rime\n" +
		"    - Fedora: sudo dnf install ibus-rime\n" +
		"    - Arch Linux: sudo pacman -S ibus-rime\n" +
		"  • FCITX-Rime (旧版)\n" +
		"    - 不推荐，建议使用 FCITX5\n\n" +
		"提示：程序仍可正常运行，但需要先安装至少一个引擎才能使用更新的配置。"

	return InstallationStatus{
		Installed: false,
		Message:   message,
	}
}
