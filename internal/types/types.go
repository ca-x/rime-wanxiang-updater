package types

import "time"

// 常量定义
const (
	VERSION      = "v1.0.0"
	OWNER        = "amzxyz"
	REPO         = "rime_wanxiang"
	CNB_REPO     = "rime-wanxiang"
	DICT_TAG     = "dict-nightly" // GitHub 词库 tag
	CNB_DICT_TAG = "v1.0.0"       // CNB 词库 tag
	MODEL_REPO   = "RIME-LMDG"
	MODEL_TAG    = "LTS"
	MODEL_FILE   = "wanxiang-lts-zh-hans.gram"
	ZH_DICTS     = "dicts"
)

// SchemeMap 方案映射
var SchemeMap = map[string]string{
	"1": "moqi",
	"2": "flypy",
	"3": "zrm",
	"4": "tiger",
	"5": "wubi",
	"6": "hanxin",
	"7": "shouyou",
}

// Config 配置结构
type Config struct {
	// 引擎配置 - 支持多引擎
	InstalledEngines []string `json:"installed_engines"` // 检测到的所有已安装引擎
	PrimaryEngine    string   `json:"primary_engine"`    // 用户选择的主引擎
	Engine           string   `json:"engine,omitempty"`  // 已弃用：保留用于配置迁移

	SchemeType          string   `json:"scheme_type"`
	SchemeFile          string   `json:"scheme_file"`
	DictFile            string   `json:"dict_file"`
	UseMirror           bool     `json:"use_mirror"`
	GithubToken         string   `json:"github_token"`
	ExcludeFiles        []string `json:"exclude_files"`
	AutoUpdate          bool     `json:"auto_update"`
	AutoUpdateCountdown int      `json:"auto_update_countdown"` // 自动更新倒计时（秒）
	ProxyEnabled        bool     `json:"proxy_enabled"`
	ProxyType           string   `json:"proxy_type"`
	ProxyAddress        string   `json:"proxy_address"`
	FcitxCompat         bool     `json:"fcitx_compat"`          // Linux 专用：兼容 ~/.config/fcitx/rime/
	FcitxUseLink        bool     `json:"fcitx_use_link"`        // Linux 专用：使用软链接（true）还是复制（false）
	FcitxConflictAction string   `json:"fcitx_conflict_action"` // Linux 专用：目录冲突处理方式 "delete" 或 "backup"，空表示未设置
	FcitxConflictPrompt bool     `json:"fcitx_conflict_prompt"` // Linux 专用：是否每次都提示（true）还是使用记忆的偏好（false）
	PreUpdateHook       string   `json:"pre_update_hook"`       // 更新前执行的脚本路径
	PostUpdateHook      string   `json:"post_update_hook"`      // 更新后执行的脚本路径

	// 主题配置
	ThemeAdaptive bool   `json:"theme_adaptive"` // 是否启用自适应主题（根据终端明暗自动切换）
	ThemeLight    string `json:"theme_light"`    // 浅色模式主题
	ThemeDark     string `json:"theme_dark"`     // 深色模式主题
	ThemeFixed    string `json:"theme_fixed"`    // 固定主题（非自适应模式时使用）
}

// UpdateInfo 更新信息
type UpdateInfo struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	UpdateTime  time.Time `json:"update_time"`
	Tag         string    `json:"tag"`
	Description string    `json:"description"`
	SHA256      string    `json:"sha256"`
	ID          string    `json:"id"`
	Size        int64     `json:"size"`
}

// UpdateRecord 更新记录
type UpdateRecord struct {
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"update_time"`
	Tag        string    `json:"tag"`
	ApplyTime  time.Time `json:"apply_time"`
	SHA256     string    `json:"sha256"`
	CnbID      string    `json:"cnb_id"`
}

// GitHubRelease GitHub Release 结构
type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	Body        string        `json:"body"`
	Assets      []GitHubAsset `json:"assets"`
	PublishedAt time.Time     `json:"published_at"`
}

// GitHubAsset GitHub Asset 结构
type GitHubAsset struct {
	Name               string    `json:"name"`
	BrowserDownloadURL string    `json:"browser_download_url"`
	UpdatedAt          time.Time `json:"updated_at"`
	Size               int64     `json:"size"`
}

// CNBRelease CNB Release 结构
type CNBRelease struct {
	Title  string     `json:"title"`
	TagRef string     `json:"tag_ref"`
	Body   string     `json:"body"`
	Assets []CNBAsset `json:"assets"`
}

// CNBAsset CNB Asset 结构
type CNBAsset struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	UpdatedAt  time.Time `json:"updated_at"`
	Digest     string    `json:"digest"`
	ID         string    `json:"id"`
	SizeInByte int64     `json:"sizeInByte"`
}

// ProgressFunc 进度回调函数
type ProgressFunc func(message string, percent float64, source string, fileName string, downloaded int64, total int64, speed float64, downloadMode bool)

// UpdateStatus 更新状态
type UpdateStatus struct {
	LocalVersion  string    // 本地版本
	RemoteVersion string    // 远程版本
	LocalTime     time.Time // 本地更新时间
	RemoteTime    time.Time // 远程更新时间
	NeedsUpdate   bool      // 是否需要更新
	Message       string    // 状态消息
}
