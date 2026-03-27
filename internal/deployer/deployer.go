package deployer

import "rime-wanxiang-updater/internal/types"

// Deployer 部署接口
type Deployer interface {
	Deploy() error
	TerminateProcesses() error
}

// GetDeployer 获取当前平台的部署器
func GetDeployer(config *types.Config) Deployer {
	return newDeployer(config)
}
