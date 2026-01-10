package ui

import (
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/controller"
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
	ViewConfigEdit      // 配置编辑
	ViewResult          // 显示更新结果
	ViewExcludeList     // 排除文件列表
	ViewExcludeEdit     // 编辑排除模式
	ViewExcludeAdd      // 添加排除模式
	ViewFcitxConflict   // Fcitx 目录冲突对话框
	ViewThemeList       // 主题列表
	ViewEngineSelector  // 引擎选择界面
	ViewEnginePrompt    // 多引擎未配置提示对话框
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
	// Communication with controller
	CommandChan chan<- controller.Command
	EventChan   <-chan controller.Event

	// Config reference (read-only for display)
	Cfg *config.Manager

	// Theme and styling
	ThemeManager *theme.Manager
	Styles       *Styles

	// UI-only state
	State      ViewState
	WizardStep WizardStep
	MenuChoice int

	// Configuration UI state
	ConfigChoice int
	EditingKey   string
	EditingValue string

	// Wizard UI state
	SchemeChoice  string
	VariantChoice string
	MirrorChoice  bool

	// Progress display (received from controller)
	Updating         bool
	Progress         progress.Model
	ProgressMsg      string
	DownloadSource   string
	DownloadFileName string
	Downloaded       int64
	TotalSize        int64
	DownloadSpeed    float64
	IsDownloading    bool

	// Result display (received from controller)
	ResultMsg        string
	ResultSuccess    bool
	ResultSkipped    bool
	AutoUpdateResult *AutoUpdateDetails

	// Display state
	Width  int
	Height int

	// Exclude file management UI state
	ExcludeListChoice   int
	ExcludeEditInput    string
	ExcludeEditIndex    int
	ExcludeErrorMsg     string
	ExcludeDescriptions []string

	// Fcitx conflict dialog state
	FcitxConflictChoice   int
	FcitxConflictNoPrompt bool
	FcitxConflictCallback func()

	// Rime installation status (display only)
	RimeInstallStatus detector.InstallationStatus

	// Auto update countdown UI state
	AutoUpdateCountdown int
	AutoUpdateCancelled bool

	// Theme selector UI state
	ThemeListChoice int
	ThemeList       []string

	// Engine selector UI state
	EngineSelections map[string]bool // 引擎名 -> 是否选中
	EngineCursor     int             // 当前光标位置
	EngineList       []string        // 引擎列表（按顺序）

	// Error display
	Err error
}

// AutoUpdateDetails 自动更新的详细结果
type AutoUpdateDetails struct {
	UpdatedComponents []string          // 已更新的组件
	SkippedComponents []string          // 跳过的组件（已是最新版本）
	ComponentVersions map[string]string // 组件版本信息（组件名 -> 版本号）
}

// CountdownTickMsg 倒计时消息
type CountdownTickMsg struct{}
