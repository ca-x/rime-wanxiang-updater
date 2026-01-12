//go:build linux

package deployer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"rime-wanxiang-updater/internal/types"
)

type linuxDeployer struct {
	config *types.Config
}

func newDeployer(config interface{}) Deployer {
	cfg, ok := config.(*types.Config)
	if !ok {
		return &linuxDeployer{}
	}
	return &linuxDeployer{config: cfg}
}

// TerminateProcesses Linux 不需要终止进程
func (d *linuxDeployer) TerminateProcesses() error {
	// Linux 输入法通常不需要强制终止
	return nil
}

// Deploy 部署 - 优先使用 qdbus6 自动部署，其次尝试 rime_deployer，最后重启输入法
func (d *linuxDeployer) Deploy() error {
	if d.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 确定 Rime 数据目录
	rimeDir, err := d.getRimeDataDir()
	if err != nil {
		return fmt.Errorf("无法确定 Rime 数据目录: %w", err)
	}

	// 尝试 1: 使用 qdbus6 自动部署 Fcitx5（最推荐的方法）
	if err := d.deployWithQdbus6(); err == nil {
		return nil
	}

	// 尝试 2: 运行 rime_deployer
	if err := d.runRimeDeployer(rimeDir); err == nil {
		return nil
	}

	// 尝试 3: 重启输入法
	return d.restartInputMethod()
}

// getRimeDataDir 获取 Rime 数据目录
// 参考 Shell 脚本：按优先级检测多个可能的位置
func (d *linuxDeployer) getRimeDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// 按优先级尝试不同的目录（与 Shell 脚本保持一致）
	candidates := []string{
		filepath.Join(homeDir, ".local/share/fcitx5/rime"),
		filepath.Join(homeDir, ".config/fcitx5/rime"),
		filepath.Join(homeDir, ".config/ibus/rime"),
		filepath.Join(homeDir, ".config/fcitx/rime"),
	}

	for _, dir := range candidates {
		// 使用 Lstat 避免跟随符号链接
		if info, err := os.Lstat(dir); err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			return dir, nil
		}
	}

	// 如果都不存在，使用默认的 fcitx5 目录
	defaultDir := filepath.Join(homeDir, ".local/share/fcitx5/rime")
	return defaultDir, nil
}

// deployWithQdbus6 使用 qdbus6 自动部署 Fcitx5
// 这是最优雅的部署方式，参考 Shell 脚本
func (d *linuxDeployer) deployWithQdbus6() error {
	// 检查 qdbus6 是否可用
	if _, err := exec.LookPath("qdbus6"); err != nil {
		return fmt.Errorf("qdbus6 不可用")
	}

	// 使用 qdbus6 触发 Fcitx5 重新部署
	cmd := exec.Command("qdbus6",
		"org.fcitx.Fcitx5",
		"/controller",
		"org.fcitx.Fcitx.Controller1.SetConfig",
		"fcitx://config/addon/rime/deploy",
		"")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("qdbus6 部署失败: %w", err)
	}

	return nil
}

// runRimeDeployer 运行 rime_deployer
func (d *linuxDeployer) runRimeDeployer(rimeDir string) error {
	// 尝试查找 rime_deployer
	deployerPaths := []string{
		"/usr/lib/rime/rime_deployer",
		"/usr/lib64/rime/rime_deployer",
		"/usr/local/lib/rime/rime_deployer",
		"rime_deployer", // 系统 PATH 中
	}

	for _, deployerPath := range deployerPaths {
		cmd := exec.Command(deployerPath, "--build", rimeDir)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("未找到 rime_deployer 或执行失败")
}

// restartInputMethod 重启输入法
func (d *linuxDeployer) restartInputMethod() error {
	// 尝试重启 fcitx5
	if err := d.restartFcitx5(); err == nil {
		return nil
	}

	// 尝试重启 ibus
	if err := d.restartIBus(); err == nil {
		return nil
	}

	return fmt.Errorf("无法重启输入法，请手动重启")
}

// restartFcitx5 重启 fcitx5
// 参考 Shell 脚本使用 fcitx5-remote -r
func (d *linuxDeployer) restartFcitx5() error {
	// 检查 fcitx5 是否运行
	checkCmd := exec.Command("pgrep", "fcitx5")
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("fcitx5 未运行")
	}

	// 使用 fcitx5-remote 重启
	cmd := exec.Command("fcitx5-remote", "-r")
	return cmd.Run()
}

// restartIBus 重启 ibus
// 参考 Shell 脚本使用 ibus-daemon -drx
func (d *linuxDeployer) restartIBus() error {
	// 检查 ibus 是否运行
	checkCmd := exec.Command("pgrep", "ibus-daemon")
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("ibus 未运行")
	}

	// 重启 ibus
	cmd := exec.Command("ibus-daemon", "-drx")
	return cmd.Run()
}

// DeployToAllEnginesWithProgress 部署到所有已安装的引擎，并报告进度
func (d *linuxDeployer) DeployToAllEnginesWithProgress(progressFunc func(engine string, index, total int)) error {
	if d.config == nil {
		return d.Deploy() // 回退到主引擎部署
	}

	// 获取要部署的引擎列表
	deployEngines := d.config.UpdateEngines
	if len(deployEngines) == 0 {
		// 未配置：默认部署所有已安装的引擎
		deployEngines = d.config.InstalledEngines
	}

	if len(deployEngines) == 0 {
		return d.Deploy() // 回退到主引擎部署
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户目录: %w", err)
	}

	var errors []string
	total := len(deployEngines)

	// 为每个引擎执行部署
	for i, engine := range deployEngines {
		if progressFunc != nil {
			progressFunc(engine, i+1, total)
		}

		// 根据引擎确定数据目录并部署
		var rimeDir string
		switch engine {
		case "fcitx5":
			rimeDir = filepath.Join(homeDir, ".local/share/fcitx5/rime")
			// 使用 Lstat 避免跟随符号链接
			if info, err := os.Lstat(rimeDir); os.IsNotExist(err) || (err == nil && info.Mode()&os.ModeSymlink != 0) {
				rimeDir = filepath.Join(homeDir, ".config/fcitx5/rime")
			}
			// 尝试 qdbus6 部署
			if err := d.deployWithQdbus6(); err != nil {
				// 失败则尝试 rime_deployer
				if err := d.runRimeDeployer(rimeDir); err != nil {
					// 最后尝试重启
					if err := d.restartFcitx5(); err != nil {
						errors = append(errors, fmt.Sprintf("%s: 部署失败", engine))
					}
				}
			}
		case "ibus":
			rimeDir = filepath.Join(homeDir, ".config/ibus/rime")
			if err := d.runRimeDeployer(rimeDir); err != nil {
				if err := d.restartIBus(); err != nil {
					errors = append(errors, fmt.Sprintf("%s: 部署失败", engine))
				}
			}
		case "fcitx":
			rimeDir = filepath.Join(homeDir, ".config/fcitx/rime")
			if err := d.runRimeDeployer(rimeDir); err != nil {
				errors = append(errors, fmt.Sprintf("%s: 部署失败", engine))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分引擎部署失败: %v", errors)
	}

	return nil
}

