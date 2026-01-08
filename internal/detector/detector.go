package detector

// InstallationStatus Rime 安装状态
type InstallationStatus struct {
	Installed bool   // 是否已安装
	Message   string // 提示信息（未安装时显示）
}

// CheckRimeInstallation 检测 Rime 是否已安装
// 该函数由各平台的特定文件实现
func CheckRimeInstallation() InstallationStatus {
	return checkRimeInstallation()
}
