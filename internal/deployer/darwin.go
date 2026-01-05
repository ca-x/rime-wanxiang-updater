//go:build darwin

package deployer

import (
	"fmt"
	"os/exec"
)

type darwinDeployer struct {
	engine string
}

func newDeployer(config interface{}) Deployer {
	// 这里应该从 config 中获取 engine 信息
	// 为了简化，这里硬编码为鼠须管
	return &darwinDeployer{
		engine: "鼠须管",
	}
}

// TerminateProcesses macOS 不需要终止进程
func (d *darwinDeployer) TerminateProcesses() error {
	return nil
}

// Deploy 部署
func (d *darwinDeployer) Deploy() error {
	var executable string
	var args []string

	if d.engine == "鼠须管" {
		executable = "/Library/Input Methods/Squirrel.app/Contents/MacOS/Squirrel"
		args = []string{"--reload"}
	} else {
		executable = "/Library/Input Methods/Fcitx5.app/Contents/bin/fcitx5-curl"
		args = []string{"/config/addon/rime/deploy", "-X", "POST", "-d", "{}"}
	}

	cmd := exec.Command(executable, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("部署失败: %w", err)
	}

	return nil
}
