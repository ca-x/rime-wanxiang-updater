package version

// Version 版本号，由编译时 ldflags 注入
// 使用方式: go build -ldflags "-X rime-wanxiang-updater/internal/version.Version=v1.0.0"
var Version = "dev"

// GetVersion 获取版本号
func GetVersion() string {
	return Version
}
