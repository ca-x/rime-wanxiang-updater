package detector

import (
	"testing"
)

// TestCheckRimeInstallation 测试 Rime 安装检测
func TestCheckRimeInstallation(t *testing.T) {
	status := CheckRimeInstallation()

	// 测试返回的状态结构是否正确
	if !status.Installed && status.Message == "" {
		t.Error("当 Rime 未安装时，应该提供安装建议消息")
	}

	if status.Installed && status.Message != "" {
		t.Error("当 Rime 已安装时，不应该有消息")
	}

	// 打印结果供调试
	t.Logf("Rime 安装状态: %v", status.Installed)
	if !status.Installed {
		t.Logf("安装建议消息:\n%s", status.Message)
	}
}
