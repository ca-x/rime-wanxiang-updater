package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"rime-wanxiang-updater/internal/api"
	"rime-wanxiang-updater/internal/types"
)

// Manager 配置管理器
type Manager struct {
	ConfigPath string
	Config     *types.Config
	RimeDir    string
	ZhDictsDir string
	CacheDir   string
}

// NewManager 创建配置管理器
func NewManager() (*Manager, error) {
	m := &Manager{
		ConfigPath: getConfigPath(),
	}

	// 加载或创建配置
	config, err := m.loadOrCreateConfig()
	if err != nil {
		return nil, err
	}
	m.Config = config

	// 设置目录
	m.RimeDir = getRimeUserDir(config)
	if config.SchemeType == "base" {
		m.ZhDictsDir = types.ZH_DICTS
	} else {
		m.ZhDictsDir = types.ZH_DICTS
	}
	m.CacheDir = getCacheDir()

	return m, nil
}

// loadOrCreateConfig 加载或创建配置
func (m *Manager) loadOrCreateConfig() (*types.Config, error) {
	if _, err := os.Stat(m.ConfigPath); os.IsNotExist(err) {
		// 创建默认配置
		config := createDefaultConfig()
		if err := m.saveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// 加载现有配置
	data, err := os.ReadFile(m.ConfigPath)
	if err != nil {
		return createDefaultConfig(), nil
	}

	var config types.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return createDefaultConfig(), nil
	}

	return &config, nil
}

// saveConfig 保存配置
func (m *Manager) saveConfig(config *types.Config) error {
	os.MkdirAll(filepath.Dir(m.ConfigPath), 0755)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	return os.WriteFile(m.ConfigPath, data, 0644)
}

// SaveConfig 保存当前配置
func (m *Manager) SaveConfig() error {
	return m.saveConfig(m.Config)
}

// createDefaultConfig 创建默认配置
func createDefaultConfig() *types.Config {
	return &types.Config{
		Engine:       detectEngine(),
		SchemeType:   "",
		SchemeFile:   "",
		DictFile:     "",
		UseMirror:    true,
		GithubToken:  "",
		ExcludeFiles: []string{},
		AutoUpdate:   false,
		ProxyEnabled: false,
		ProxyType:    "http",
		ProxyAddress: "127.0.0.1:7890",
		FcitxCompat:  false,
		FcitxUseLink: true, // 默认使用软链接
	}
}

// detectEngine 检测输入法引擎
func detectEngine() string {
	switch runtime.GOOS {
	case "windows":
		return "小狼毫"
	case "darwin":
		return "鼠须管"
	default:
		return "fcitx5"
	}
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".rime-updater", "config.json")
}

// getCacheDir 获取缓存目录
func getCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".rime-updater", "cache")
	os.MkdirAll(cacheDir, 0755)
	return cacheDir
}

// GetActualFilenames 获取实际文件名
func (m *Manager) GetActualFilenames(schemeKey string) (string, string, error) {
	var schemePattern, dictPattern string

	if schemeKey == "base" {
		schemePattern = `rime-wanxiang-base\.zip`
		dictPattern = `base-dicts\.zip`
	} else {
		schemePattern = fmt.Sprintf(`.*%s.*fuzhu\.zip`, schemeKey)
		dictPattern = fmt.Sprintf(`pro-%s-fuzhu-dicts\.zip`, schemeKey)
	}

	schemeRegex := regexp.MustCompile(schemePattern)
	dictRegex := regexp.MustCompile(dictPattern)

	// 创建 API 客户端
	client := api.NewClient(m.Config)

	// 获取方案文件
	var releases []types.GitHubRelease
	var err error

	if m.Config.UseMirror {
		releases, err = client.FetchCNBReleases(types.OWNER, types.CNB_REPO, "")
	} else {
		releases, err = client.FetchGitHubReleases(types.OWNER, types.REPO, "")
	}

	if err != nil {
		return "", "", fmt.Errorf("获取版本信息失败: %w", err)
	}

	var schemeFile, dictFile string

	// 查找方案文件
	for _, release := range releases {
		for _, asset := range release.Assets {
			if schemeRegex.MatchString(asset.Name) {
				schemeFile = asset.Name
				break
			}
		}
		if schemeFile != "" {
			break
		}
	}

	// 获取词库文件
	if m.Config.UseMirror {
		// CNB 使用 v1.0.0 tag
		releases, err = client.FetchCNBReleases(types.OWNER, types.CNB_REPO, types.CNB_DICT_TAG)
	} else {
		// GitHub 使用 dict-nightly tag
		releases, err = client.FetchGitHubReleases(types.OWNER, types.REPO, types.DICT_TAG)
	}

	if err != nil {
		return "", "", fmt.Errorf("获取词库信息失败: %w", err)
	}

	// 查找词库文件
	for _, release := range releases {
		for _, asset := range release.Assets {
			if dictRegex.MatchString(asset.Name) {
				dictFile = asset.Name
				break
			}
		}
		if dictFile != "" {
			break
		}
	}

	if schemeFile == "" {
		return "", "", fmt.Errorf("未找到匹配的方案文件")
	}
	if dictFile == "" {
		return "", "", fmt.Errorf("未找到匹配的词库文件")
	}

	return schemeFile, dictFile, nil
}

// GetExtractPath 获取解压路径
func (m *Manager) GetExtractPath() string {
	return m.RimeDir
}

// GetDictExtractPath 获取词库解压路径
func (m *Manager) GetDictExtractPath() string {
	return filepath.Join(m.RimeDir, m.ZhDictsDir)
}

// GetSchemeRecordPath 获取方案记录文件路径
func (m *Manager) GetSchemeRecordPath() string {
	return filepath.Join(m.CacheDir, "scheme_record.json")
}

// GetDictRecordPath 获取词库记录文件路径
func (m *Manager) GetDictRecordPath() string {
	return filepath.Join(m.CacheDir, "dict_record.json")
}

// GetModelRecordPath 获取模型记录文件路径
func (m *Manager) GetModelRecordPath() string {
	return filepath.Join(m.CacheDir, "model_record.json")
}

// ValidateExcludeFiles 验证排除文件配置
func ValidateExcludeFiles(patterns []string) error {
	for _, pattern := range patterns {
		if strings.TrimSpace(pattern) == "" {
			continue
		}
		// 验证正则表达式
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("无效的排除模式 %s: %w", pattern, err)
		}
	}
	return nil
}

// SyncToFcitxDir 同步到 fcitx 兼容目录
// 仅在 Linux 平台且启用 FcitxCompat 时生效
func (m *Manager) SyncToFcitxDir() error {
	// 仅 Linux 平台支持
	if runtime.GOOS != "linux" {
		return nil
	}

	// 未启用 fcitx 兼容
	if !m.Config.FcitxCompat {
		return nil
	}

	// 获取源目录（fcitx5 配置目录）
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %w", err)
	}

	sourceDir := m.RimeDir
	targetDir := filepath.Join(homeDir, ".config", "fcitx", "rime")

	// 检查源目录是否存在
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("源目录不存在: %s", sourceDir)
	}

	// 创建目标父目录
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("创建目标父目录失败: %w", err)
	}

	// 如果目标已存在，先删除
	if _, err := os.Lstat(targetDir); err == nil {
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("删除旧目标失败: %w", err)
		}
	}

	// 根据配置选择软链接或复制
	if m.Config.FcitxUseLink {
		// 创建软链接
		if err := os.Symlink(sourceDir, targetDir); err != nil {
			return fmt.Errorf("创建软链接失败: %w", err)
		}
	} else {
		// 复制目录
		if err := copyDir(sourceDir, targetDir); err != nil {
			return fmt.Errorf("复制目录失败: %w", err)
		}
	}

	return nil
}

// copyDir 递归复制目录
func copyDir(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 递归复制每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 读取源文件
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// 获取源文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 写入目标文件
	return os.WriteFile(dst, data, srcInfo.Mode())
}
