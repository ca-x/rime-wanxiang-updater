package deployer

// Deployer 部署接口
type Deployer interface {
	Deploy() error
	TerminateProcesses() error
}

// GetDeployer 获取当前平台的部署器
func GetDeployer(config interface{}) Deployer {
	return newDeployer(config)
}
