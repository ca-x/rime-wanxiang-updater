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

// Deploy 部署 - 运行 rime_deployer 或重启输入法
func (d *linuxDeployer) Deploy() error {
	if d.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 确定 Rime 数据目录
	rimeDir, err := d.getRimeDataDir()
	if err != nil {
		return fmt.Errorf("无法确定 Rime 数据目录: %w", err)
	}

	// 尝试运行 rime_deployer
	if err := d.runRimeDeployer(rimeDir); err == nil {
		return nil
	}

	// 如果 rime_deployer 失败，尝试重启输入法
	return d.restartInputMethod()
}

// getRimeDataDir 获取 Rime 数据目录
func (d *linuxDeployer) getRimeDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// 按优先级尝试不同的目录
	candidates := []string{
		filepath.Join(homeDir, ".local/share/fcitx5/rime"),
		filepath.Join(homeDir, ".config/fcitx5/rime"),
		filepath.Join(homeDir, ".config/ibus/rime"),
		filepath.Join(homeDir, ".config/fcitx/rime"),
	}

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, nil
		}
	}

	// 如果都不存在，使用默认的 fcitx5 目录
	defaultDir := filepath.Join(homeDir, ".local/share/fcitx5/rime")
	return defaultDir, nil
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
func (d *linuxDeployer) restartIBus() error {
	// 检查 ibus 是否运行
	checkCmd := exec.Command("pgrep", "ibus-daemon")
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("ibus 未运行")
	}

	// 重启 ibus
	cmd := exec.Command("ibus", "restart")
	return cmd.Run()
}
