//go:build windows

package deployer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/registry"
)

type windowsDeployer struct {
	weaselServer string
}

func newDeployer(config interface{}) Deployer {
	return &windowsDeployer{
		weaselServer: getWeaselServerPath(),
	}
}

// getWeaselServerPath 获取 WeaselServer.exe 路径
func getWeaselServerPath() string {
	// 尝试从注册表读取
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Rime\Weasel`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		if root, _, err := k.GetStringValue("WeaselRoot"); err == nil {
			if server, _, err := k.GetStringValue("ServerExecutable"); err == nil {
				return filepath.Join(root, server)
			}
		}
	}

	// 回退到默认路径
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Rime", "weasel-x64", "WeaselServer.exe")
}

// TerminateProcesses 终止进程
func (d *windowsDeployer) TerminateProcesses() error {
	if !d.gracefulStop() {
		d.hardStop()
	}
	return nil
}

// gracefulStop 优雅停止服务
func (d *windowsDeployer) gracefulStop() bool {
	cmd := exec.Command(d.weaselServer, "/q")
	if err := cmd.Run(); err != nil {
		return false
	}
	time.Sleep(500 * time.Millisecond)
	return true
}

// hardStop 强制终止
func (d *windowsDeployer) hardStop() {
	for i := 0; i < 3; i++ {
		exec.Command("taskkill", "/IM", "WeaselServer.exe", "/F").Run()
		exec.Command("taskkill", "/IM", "WeaselDeployer.exe", "/F").Run()
		time.Sleep(500 * time.Millisecond)
	}
}

// Deploy 部署
func (d *windowsDeployer) Deploy() error {
	// 启动服务（带重试）
	for retry := 0; retry < 3; retry++ {
		cmd := exec.Command(d.weaselServer)
		if err := cmd.Start(); err == nil {
			time.Sleep(2 * time.Second)
			break
		}
		if retry == 2 {
			return fmt.Errorf("启动 WeaselServer 失败")
		}
		time.Sleep(1 * time.Second)
	}

	// 执行部署
	deployer := filepath.Join(filepath.Dir(d.weaselServer), "WeaselDeployer.exe")
	return exec.Command(deployer, "/deploy").Run()
}
