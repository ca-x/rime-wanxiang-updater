//go:build darwin

package deployer

import (
	"fmt"
	"os/exec"
	"strings"

	"rime-wanxiang-updater/internal/types"
)

type darwinDeployer struct {
	config *types.Config
}

func newDeployer(config interface{}) Deployer {
	cfg, ok := config.(*types.Config)
	if !ok {
		return &darwinDeployer{}
	}
	return &darwinDeployer{config: cfg}
}

// TerminateProcesses macOS 不需要终止进程
func (d *darwinDeployer) TerminateProcesses() error {
	return nil
}

// Deploy 部署到主引擎
func (d *darwinDeployer) Deploy() error {
	if d.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	engine := d.config.PrimaryEngine
	if engine == "" && len(d.config.InstalledEngines) > 0 {
		engine = d.config.InstalledEngines[0]
	}
	if engine == "" {
		// 向后兼容
		engine = d.config.Engine
	}
	if engine == "" {
		engine = "鼠须管"
	}

	return d.deployToEngine(engine)
}

// deployToEngine 部署到指定引擎
func (d *darwinDeployer) deployToEngine(engine string) error {
	var executable string
	var args []string

	switch engine {
	case "鼠须管":
		executable = "/Library/Input Methods/Squirrel.app/Contents/MacOS/Squirrel"
		args = []string{"--reload"}
	case "小企鹅":
		// TODO: 确认小企鹅的实际安装路径
		// 可能是 Fcitx5.app 或 Fcitx5Installer.app
		executable = "/Library/Input Methods/Fcitx5.app/Contents/bin/fcitx5-curl"
		args = []string{"/config/addon/rime/deploy", "-X", "POST", "-d", "{}"}
	default:
		return fmt.Errorf("不支持的引擎: %s", engine)
	}

	cmd := exec.Command(executable, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("部署到 %s 失败: %w", engine, err)
	}

	return nil
}

// DeployToAllEngines 部署到所有已安装的引擎
func (d *darwinDeployer) DeployToAllEngines() error {
	if d.config == nil || len(d.config.InstalledEngines) == 0 {
		return d.Deploy() // 回退到部署主引擎
	}

	var errors []string
	for _, engine := range d.config.InstalledEngines {
		if err := d.deployToEngine(engine); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", engine, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分引擎部署失败: %s", strings.Join(errors, "; "))
	}

	return nil
}
