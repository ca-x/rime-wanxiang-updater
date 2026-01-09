package ui

import (
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"

	"github.com/charmbracelet/bubbles/progress"
)

// ViewState 视图状态
type ViewState int

const (
	ViewWizard ViewState = iota
	ViewMenu
	ViewUpdating
	ViewConfig
	ViewConfigEdit    // 配置编辑
	ViewResult        // 显示更新结果
	ViewExcludeList   // 排除文件列表
	ViewExcludeEdit   // 编辑排除模式
	ViewExcludeAdd    // 添加排除模式
	ViewFcitxConflict // Fcitx 目录冲突对话框
	ViewThemeList     // 主题列表
)

// WizardStep 向导步骤
type WizardStep int

const (
	WizardSchemeType WizardStep = iota
	WizardSchemeVariant
	WizardDownloadSource
	WizardComplete
)

// Model Bubble Tea 模型
type Model struct {
	Cfg              *config.Manager
	ThemeManager     *theme.Manager // 主题管理器
	Styles           *Styles        // 主题化样式
	State            ViewState
	WizardStep       WizardStep
	MenuChoice       int
	ConfigChoice     int    // 配置菜单选择
	EditingKey       string // 正在编辑的配置键
	EditingValue     string // 编辑中的值
	SchemeChoice     string
	VariantChoice    string
	MirrorChoice     bool // 是否使用镜像
	Updating         bool
	Progress         progress.Model
	ProgressMsg      string
	DownloadSource   string                 // 下载源
	DownloadFileName string                 // 下载文件名
	Downloaded       int64                  // 已下载字节
	TotalSize        int64                  // 总大小字节
	DownloadSpeed    float64                // 下载速度
	IsDownloading    bool                   // 是否在下载中
	ProgressChan     chan UpdateMsg         // 进度通道
	CompletionChan   chan UpdateCompleteMsg // 完成通道
	Err              error
	ResultMsg        string             // 结果消息
	ResultSuccess    bool               // 是否成功
	ResultSkipped    bool               // 是否跳过更新（已是最新版本）
	AutoUpdateResult *AutoUpdateDetails // 自动更新的详细结果
	Width            int
	Height           int

	// 排除文件管理相关
	ExcludeListChoice   int      // 排除列表光标位置
	ExcludeEditInput    string   // 编辑/添加排除模式的输入
	ExcludeEditIndex    int      // 正在编辑的模式索引
	ExcludeErrorMsg     string   // 排除模式错误消息
	ExcludeDescriptions []string // 排除模式的描述

	// Fcitx 冲突处理相关
	FcitxConflictChoice   int    // 对话框按钮选择 (0=删除, 1=备份, 2=不再提示复选框)
	FcitxConflictNoPrompt bool   // 是否选中"不再提示"
	FcitxConflictCallback func() // 冲突解决后的回调函数

	// Rime 安装检测
	RimeInstallStatus detector.InstallationStatus // Rime 安装状态

	// 自动更新倒计时相关
	AutoUpdateCountdown int  // 自动更新倒计时（秒）
	AutoUpdateCancelled bool // 是否已取消自动更新

	// 主题选择相关
	ThemeListChoice int      // 主题列表光标位置
	ThemeList       []string // 当前显示的主题列表
}

// UpdateMsg 更新消息类型
type UpdateMsg struct {
	Message      string
	Percent      float64
	Source       string  // 下载源
	FileName     string  // 文件名
	Downloaded   int64   // 已下载字节
	Total        int64   // 总大小字节
	Speed        float64 // 下载速度 MB/s
	DownloadMode bool    // 是否在下载模式
}

// UpdateCompleteMsg 更新完成消息
type UpdateCompleteMsg struct {
	Err           error
	UpdateType    string // 更新类型：词库、方案、模型、自动
	Skipped       bool   // 是否跳过更新（已是最新版本）
	StatusMessage string // 状态消息（包含版本信息）
	// 自动更新的详细结果
	AutoUpdateDetails *AutoUpdateDetails
}

// AutoUpdateDetails 自动更新的详细结果
type AutoUpdateDetails struct {
	UpdatedComponents []string          // 已更新的组件
	SkippedComponents []string          // 跳过的组件（已是最新版本）
	ComponentVersions map[string]string // 组件版本信息（组件名 -> 版本号）
}

// CountdownTickMsg 倒计时消息
type CountdownTickMsg struct{}
