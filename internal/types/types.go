package types

import "time"

// 常量定义
const (
	VERSION    = "v1.0.0"
	OWNER      = "amzxyz"
	REPO       = "rime_wanxiang"
	CNB_REPO   = "rime-wanxiang"
	DICT_TAG   = "dict-nightly"
	MODEL_REPO = "RIME-LMDG"
	MODEL_TAG  = "LTS"
	MODEL_FILE = "wanxiang-lts-zh-hans.gram"
	ZH_DICTS   = "dicts"
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
	Engine       string   `json:"engine"`
	SchemeType   string   `json:"scheme_type"`
	SchemeFile   string   `json:"scheme_file"`
	DictFile     string   `json:"dict_file"`
	UseMirror    bool     `json:"use_mirror"`
	GithubToken  string   `json:"github_token"`
	ExcludeFiles []string `json:"exclude_files"`
	AutoUpdate   bool     `json:"auto_update"`
	ProxyEnabled bool     `json:"proxy_enabled"`
	ProxyType    string   `json:"proxy_type"`
	ProxyAddress string   `json:"proxy_address"`
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
