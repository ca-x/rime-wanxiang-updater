package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

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

	// 配置迁移：从旧的 Engine 字段迁移到新的多引擎结构
	if config.Engine != "" {
		config.PrimaryEngine = config.Engine
		config.Engine = ""
		// 迁移后保存
		if err := m.saveConfig(&config); err != nil {
			fmt.Printf("警告：配置迁移失败: %v\n", err)
		}
	}

	// 每次启动都重新检测已安装的引擎
	config.InstalledEngines = DetectInstalledEngines()

	// 验证并清理用户配置
	needsSave := false

	// 验证 PrimaryEngine 是否仍然存在
	if config.PrimaryEngine != "" {
		found := false
		for _, engine := range config.InstalledEngines {
			if engine == config.PrimaryEngine {
				found = true
				break
			}
		}
		if !found {
			// 主引擎已不存在
			if len(config.InstalledEngines) > 0 {
				oldPrimary := config.PrimaryEngine
				config.PrimaryEngine = config.InstalledEngines[0]
				fmt.Printf("⚠️  主引擎 %s 未检测到，已切换到 %s\n", oldPrimary, config.PrimaryEngine)
				needsSave = true
			} else {
				config.PrimaryEngine = ""
				fmt.Println("⚠️  未检测到任何已安装的引擎")
			}
		}
	}

	// 如果没有主引擎但有已安装的引擎，设置第一个为主引擎
	if config.PrimaryEngine == "" && len(config.InstalledEngines) > 0 {
		config.PrimaryEngine = config.InstalledEngines[0]
		needsSave = true
	}

	// 清理 UpdateEngines 中已卸载的引擎
	if len(config.UpdateEngines) > 0 {
		var validUpdateEngines []string
		for _, engine := range config.UpdateEngines {
			for _, installed := range config.InstalledEngines {
				if engine == installed {
					validUpdateEngines = append(validUpdateEngines, engine)
					break
				}
			}
		}
		if len(validUpdateEngines) != len(config.UpdateEngines) {
			config.UpdateEngines = validUpdateEngines
			needsSave = true
		}
	}

	// 修复 UpdateEngines：如果为空且有多个引擎，默认设置为所有已安装引擎
	if len(config.UpdateEngines) == 0 && len(config.InstalledEngines) > 1 {
		config.UpdateEngines = make([]string, len(config.InstalledEngines))
		copy(config.UpdateEngines, config.InstalledEngines)
		needsSave = true
	}

	// 保存更新后的配置
	if needsSave {
		if err := m.saveConfig(&config); err != nil {
			fmt.Printf("警告：保存配置失败: %v\n", err)
		} else {
			fmt.Println("✓ 配置已自动更新")
		}
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

// RedetectEngines 重新检测已安装的引擎
func (m *Manager) RedetectEngines() error {
	m.Config.InstalledEngines = DetectInstalledEngines()

	// 如果主引擎不在新的已安装列表中，重新设置
	if m.Config.PrimaryEngine != "" {
		found := false
		for _, engine := range m.Config.InstalledEngines {
			if engine == m.Config.PrimaryEngine {
				found = true
				break
			}
		}
		if !found && len(m.Config.InstalledEngines) > 0 {
			m.Config.PrimaryEngine = m.Config.InstalledEngines[0]
		}
	} else if len(m.Config.InstalledEngines) > 0 {
		m.Config.PrimaryEngine = m.Config.InstalledEngines[0]
	}

	return m.SaveConfig()
}

// GetEngineDisplayName 获取引擎显示名称
// 多引擎时用 + 连接
func (m *Manager) GetEngineDisplayName() string {
	if len(m.Config.InstalledEngines) == 0 {
		return m.Config.PrimaryEngine
	}

	if len(m.Config.InstalledEngines) == 1 {
		return m.Config.InstalledEngines[0]
	}

	// 多个引擎，用 + 连接
	return strings.Join(m.Config.InstalledEngines, "+")
}

// AddExcludePattern 添加排除模式
func (m *Manager) AddExcludePattern(pattern string) error {
	// 验证模式
	_, err := ParseExcludePattern(pattern)
	if err != nil {
		return fmt.Errorf("无效的排除模式: %w", err)
	}

	// 检查是否已存在
	for _, existing := range m.Config.ExcludeFiles {
		if existing == pattern {
			return fmt.Errorf("模式已存在")
		}
	}

	m.Config.ExcludeFiles = append(m.Config.ExcludeFiles, pattern)
	return m.SaveConfig()
}

// RemoveExcludePattern 删除排除模式
func (m *Manager) RemoveExcludePattern(index int) error {
	if index < 0 || index >= len(m.Config.ExcludeFiles) {
		return fmt.Errorf("索引超出范围")
	}

	m.Config.ExcludeFiles = append(
		m.Config.ExcludeFiles[:index],
		m.Config.ExcludeFiles[index+1:]...,
	)
	return m.SaveConfig()
}

// ResetExcludePatterns 重置为默认排除模式
func (m *Manager) ResetExcludePatterns() error {
	m.Config.ExcludeFiles = make([]string, len(DefaultExcludePatterns))
	copy(m.Config.ExcludeFiles, DefaultExcludePatterns)
	return m.SaveConfig()
}

// GetExcludePatternDescriptions 获取所有排除模式的描述
func (m *Manager) GetExcludePatternDescriptions() ([]string, error) {
	patterns, errs := ParseExcludePatterns(m.Config.ExcludeFiles)
	if len(errs) > 0 {
		return nil, fmt.Errorf("部分模式解析失败: %v", errs[0])
	}

	descriptions := make([]string, len(patterns))
	for i, p := range patterns {
		descriptions[i] = p.GetPatternDescription()
	}
	return descriptions, nil
}

// createDefaultConfig 创建默认配置
func createDefaultConfig() *types.Config {
	// 检测已安装的引擎
	installedEngines := DetectInstalledEngines()
	var primaryEngine string
	if len(installedEngines) > 0 {
		primaryEngine = installedEngines[0]
	}

	return &types.Config{
		InstalledEngines:    installedEngines,
		PrimaryEngine:       primaryEngine,
		SchemeType:          "",
		SchemeFile:          "",
		DictFile:            "",
		UseMirror:           true,
		GithubToken:         "",
		ExcludeFiles:        DefaultExcludePatterns, // 使用默认排除模式
		AutoUpdate:          false,
		AutoUpdateCountdown: 5, // 默认倒计时 5 秒
		ProxyEnabled:        false,
		ProxyType:           "http",
		ProxyAddress:        "127.0.0.1:7890",
		FcitxCompat:         false,
		FcitxUseLink:        true, // 默认使用软链接
		FcitxConflictAction: "",   // 默认未设置，首次会询问
		FcitxConflictPrompt: true, // 默认每次都提示
		PreUpdateHook:       "",
		PostUpdateHook:      "",
		ThemeAdaptive:       true,              // 默认启用自适应主题
		ThemeLight:          "cyberpunk-light", // 默认浅色主题
		ThemeDark:           "cyberpunk",       // 默认深色主题
		ThemeFixed:          "",                // 固定主题为空
	}
}

// detectEngine 检测输入法引擎（已弃用，保留向后兼容）
// 使用 DetectInstalledEngines() 替代
func detectEngine() string {
	installed := DetectInstalledEngines()
	if len(installed) > 0 {
		return installed[0]
	}
	return ""
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
// 返回详细的错误信息和建议
func ValidateExcludeFiles(patterns []string) error {
	var errors []string

	for i, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// 尝试解析模式
		_, err := ParseExcludePattern(pattern)
		if err != nil {
			// 提供友好的错误信息
			suggestion := getSuggestionForPattern(pattern)
			errors = append(errors, fmt.Sprintf(
				"[行 %d] 模式 '%s' 无效: %v\n提示: %s",
				i+1, pattern, err, suggestion,
			))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("排除文件配置验证失败:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// getSuggestionForPattern 为无效的模式提供修正建议
func getSuggestionForPattern(pattern string) string {
	// 常见错误模式和建议
	suggestions := map[string]string{
		".*userdb":  "可能想写: *.userdb 或 .*\\.userdb$",
		"sync/*":    "匹配 sync 目录下所有文件",
		"*.yaml":    "匹配所有 yaml 文件",
		"^sync/.*$": "正则表达式：匹配 sync/ 开头的所有文件",
	}

	// 检查是否有相似的有效模式
	for valid, desc := range suggestions {
		if strings.Contains(pattern, strings.TrimSuffix(strings.TrimPrefix(valid, "^"), "$")) {
			return fmt.Sprintf("您可能想使用: %s (%s)", valid, desc)
		}
	}

	// 根据模式内容提供通用建议
	if strings.Contains(pattern, ".") && !strings.Contains(pattern, "\\") {
		return "如果要匹配点号(.)，在正则表达式中需要转义: \\."
	}

	if strings.HasPrefix(pattern, "/") || strings.HasPrefix(pattern, "\\") {
		return "路径不需要以斜杠开头，例如: sync/*.txt 而不是 /sync/*.txt"
	}

	return "参考示例: *.userdb (通配符) 或 ^sync/.*$ (正则) 或 user.yaml (精确匹配)"
}

// SyncToFcitxDir 同步到 fcitx 兼容目录
// 仅在 Linux 平台且启用 FcitxCompat 时生效
// 返回 needsPrompt 表示需要用户确认，conflictExists 表示目录已存在
func (m *Manager) SyncToFcitxDir() (needsPrompt bool, conflictExists bool, err error) {
	// 仅 Linux 平台支持
	if runtime.GOOS != "linux" {
		return false, false, nil
	}

	// 未启用 fcitx 兼容
	if !m.Config.FcitxCompat {
		return false, false, nil
	}

	// 获取源目录（fcitx5 配置目录）
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, false, fmt.Errorf("获取用户目录失败: %w", err)
	}

	sourceDir := m.RimeDir
	targetDir := filepath.Join(homeDir, ".config", "fcitx", "rime")

	// 检查源目录是否存在
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return false, false, fmt.Errorf("源目录不存在: %s", sourceDir)
	}

	// 创建目标父目录
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return false, false, fmt.Errorf("创建目标父目录失败: %w", err)
	}

	// 检查目标是否存在
	targetInfo, err := os.Lstat(targetDir)
	if err == nil {
		// 目标存在，检查是否是指向正确位置的软链接
		if targetInfo.Mode()&os.ModeSymlink != 0 {
			// 是软链接，检查指向
			link, err := os.Readlink(targetDir)
			if err == nil && link == sourceDir {
				// 已经是正确的软链接，无需操作
				return false, false, nil
			}
		}

		// 目标存在但不是正确的软链接，需要处理冲突
		conflictExists = true

		// 检查是否需要提示用户
		if m.Config.FcitxConflictPrompt || m.Config.FcitxConflictAction == "" {
			// 需要用户确认
			return true, true, nil
		}

		// 使用已保存的偏好处理冲突
		if err := m.handleFcitxConflict(targetDir); err != nil {
			return false, true, err
		}
	}

	// 根据配置选择软链接或复制
	if m.Config.FcitxUseLink {
		// 创建软链接
		if err := os.Symlink(sourceDir, targetDir); err != nil {
			return false, false, fmt.Errorf("创建软链接失败: %w", err)
		}
	} else {
		// 复制目录
		if err := copyDir(sourceDir, targetDir); err != nil {
			return false, false, fmt.Errorf("复制目录失败: %w", err)
		}
	}

	return false, false, nil
}

// handleFcitxConflict 处理 fcitx 目录冲突（使用已保存的偏好）
func (m *Manager) handleFcitxConflict(targetDir string) error {
	if m.Config.FcitxConflictAction == "backup" {
		// 备份现有目录，使用时间戳格式：rime.backup.20060102_150405
		timestamp := time.Now().Format("20060102_150405")
		backupName := "rime.backup." + timestamp
		backupDir := filepath.Join(filepath.Dir(targetDir), backupName)

		if err := os.Rename(targetDir, backupDir); err != nil {
			return fmt.Errorf("备份目录失败: %w", err)
		}
	} else {
		// 直接删除
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("删除旧目标失败: %w", err)
		}
	}
	return nil
}

// ResolveFcitxConflict 解决 fcitx 目录冲突（在用户确认后调用）
func (m *Manager) ResolveFcitxConflict() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %w", err)
	}

	targetDir := filepath.Join(homeDir, ".config", "fcitx", "rime")

	// 处理冲突
	if err := m.handleFcitxConflict(targetDir); err != nil {
		return err
	}

	// 执行同步
	_, _, err = m.SyncToFcitxDir()
	return err
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

// ExecutePreUpdateHook 执行更新前 hook
// 如果 hook 执行失败，返回错误以取消更新
func (m *Manager) ExecutePreUpdateHook() error {
	if m.Config.PreUpdateHook == "" {
		// Hook 未设置，直接返回
		return nil
	}

	// 展开路径中的 ~ 为用户目录
	hookPath := expandPath(m.Config.PreUpdateHook)

	// 检查脚本是否存在
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return fmt.Errorf("pre-update hook 脚本不存在: %s", hookPath)
	}

	// 执行脚本
	cmd := exec.Command(hookPath)
	cmd.Env = append(os.Environ(),
		"RIME_DIR="+m.RimeDir,
		"RIME_CACHE_DIR="+m.CacheDir,
		"HOOK_TYPE=pre_update",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pre-update hook 执行失败: %w\n输出: %s", err, string(output))
	}

	return nil
}

// ExecutePostUpdateHook 执行更新后 hook
// 即使失败也不影响更新结果，只记录错误
func (m *Manager) ExecutePostUpdateHook() error {
	if m.Config.PostUpdateHook == "" {
		return nil
	}

	// 展开路径中的 ~ 为用户目录
	hookPath := expandPath(m.Config.PostUpdateHook)

	// 检查脚本是否存在
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return fmt.Errorf("post-update hook 脚本不存在: %s", hookPath)
	}

	// 执行脚本
	cmd := exec.Command(hookPath)
	cmd.Env = append(os.Environ(),
		"RIME_DIR="+m.RimeDir,
		"RIME_CACHE_DIR="+m.CacheDir,
		"HOOK_TYPE=post_update",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("post-update hook 执行失败: %w\n输出: %s", err, string(output))
	}

	return nil
}

// expandPath 展开路径中的 ~ 为用户目录
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
