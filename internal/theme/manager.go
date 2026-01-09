package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"rime-wanxiang-updater/internal/base16"
	"rime-wanxiang-updater/internal/termcolor"

	"github.com/charmbracelet/lipgloss"
)

// Manager 主题管理器
type Manager struct {
	schemes      map[string]*base16.Scheme
	currentName  string
	currentTheme *base16.Theme
	adaptiveMode bool
	lightScheme  string
	darkScheme   string
	background   termcolor.Background
}

// NewManager 创建主题管理器
func NewManager() *Manager {
	m := &Manager{
		schemes:      make(map[string]*base16.Scheme),
		adaptiveMode: true,
		lightScheme:  "cyberpunk-light",
		darkScheme:   "cyberpunk",
	}

	// 注册内置主题
	for name, scheme := range base16.BuiltinSchemes() {
		m.schemes[name] = scheme
	}

	// 检测终端背景并设置默认主题
	m.background = termcolor.InitLipgloss()
	m.applyAdaptiveTheme()

	return m
}

// LoadFromDirectory 从目录加载自定义主题
func (m *Manager) LoadFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read themes directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		scheme, err := base16.LoadFromYAML(path)
		if err != nil {
			continue
		}

		name := strings.ToLower(strings.TrimSuffix(entry.Name(), ext))
		m.schemes[name] = scheme
	}

	return nil
}

// SetTheme 设置当前主题
func (m *Manager) SetTheme(name string) error {
	name = strings.ToLower(name)
	scheme, ok := m.schemes[name]
	if !ok {
		return fmt.Errorf("theme '%s' not found", name)
	}

	m.currentName = name
	m.currentTheme = base16.NewTheme(scheme)
	m.adaptiveMode = false

	return nil
}

// SetAdaptiveTheme 设置自适应主题（明暗模式自动切换）
func (m *Manager) SetAdaptiveTheme(lightName, darkName string) error {
	lightName = strings.ToLower(lightName)
	darkName = strings.ToLower(darkName)

	if _, ok := m.schemes[lightName]; !ok {
		return fmt.Errorf("light theme '%s' not found", lightName)
	}
	if _, ok := m.schemes[darkName]; !ok {
		return fmt.Errorf("dark theme '%s' not found", darkName)
	}

	m.lightScheme = lightName
	m.darkScheme = darkName
	m.adaptiveMode = true
	m.applyAdaptiveTheme()

	return nil
}

// applyAdaptiveTheme 应用自适应主题
func (m *Manager) applyAdaptiveTheme() {
	if m.background.IsDark() {
		m.currentName = m.darkScheme
	} else {
		m.currentName = m.lightScheme
	}

	if scheme, ok := m.schemes[m.currentName]; ok {
		m.currentTheme = base16.NewTheme(scheme)
	}
}

// Current 获取当前主题
func (m *Manager) Current() *base16.Theme {
	return m.currentTheme
}

// CurrentScheme 获取当前配色方案
func (m *Manager) CurrentScheme() *base16.Scheme {
	if m.currentTheme != nil {
		return m.currentTheme.Scheme
	}
	return nil
}

// CurrentName 获取当前主题名称
func (m *Manager) CurrentName() string {
	return m.currentName
}

// IsAdaptive 是否为自适应模式
func (m *Manager) IsAdaptive() bool {
	return m.adaptiveMode
}

// GetAdaptiveSchemes 获取自适应模式的明暗主题名称
func (m *Manager) GetAdaptiveSchemes() (light, dark string) {
	return m.lightScheme, m.darkScheme
}

// List 列出所有可用主题
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.schemes))
	for name := range m.schemes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListDark 列出所有深色主题
func (m *Manager) ListDark() []string {
	darkThemes := []string{
		"cyberpunk", "nord", "dracula", "gruvbox-dark", "monokai",
		"tokyo-night", "catppuccin-mocha", "one-dark",
		"solarized-dark", "classic-dark",
	}

	result := make([]string, 0)
	for _, name := range darkThemes {
		if _, ok := m.schemes[name]; ok {
			result = append(result, name)
		}
	}
	return result
}

// ListLight 列出所有浅色主题
func (m *Manager) ListLight() []string {
	lightThemes := []string{
		"cyberpunk-light", "gruvbox-light", "tokyo-night-light",
		"catppuccin-latte", "one-light", "solarized-light", "classic-light",
	}

	result := make([]string, 0)
	for _, name := range lightThemes {
		if _, ok := m.schemes[name]; ok {
			result = append(result, name)
		}
	}
	return result
}

// Get 获取指定主题
func (m *Manager) Get(name string) (*base16.Theme, error) {
	name = strings.ToLower(name)
	scheme, ok := m.schemes[name]
	if !ok {
		return nil, fmt.Errorf("theme '%s' not found", name)
	}
	return base16.NewTheme(scheme), nil
}

// GetScheme 获取指定配色方案
func (m *Manager) GetScheme(name string) (*base16.Scheme, error) {
	name = strings.ToLower(name)
	scheme, ok := m.schemes[name]
	if !ok {
		return nil, fmt.Errorf("scheme '%s' not found", name)
	}
	return scheme, nil
}

// Register 注册自定义主题
func (m *Manager) Register(name string, scheme *base16.Scheme) {
	m.schemes[strings.ToLower(name)] = scheme
}

// Background 获取当前终端背景
func (m *Manager) Background() termcolor.Background {
	return m.background
}

// RefreshBackground 重新检测终端背景并更新主题
func (m *Manager) RefreshBackground() {
	m.background = termcolor.DetectBackground()
	if m.adaptiveMode {
		m.applyAdaptiveTheme()
	}
}

// CreateAdaptiveColor 创建自适应颜色
func (m *Manager) CreateAdaptiveColor(lightBase, darkBase int) lipgloss.AdaptiveColor {
	lightScheme := m.schemes[m.lightScheme]
	darkScheme := m.schemes[m.darkScheme]

	if lightScheme == nil || darkScheme == nil {
		return lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	}

	return base16.AdaptiveColor(lightScheme, darkScheme, lightBase)
}

// Config 主题配置（用于保存到配置文件）
type Config struct {
	AdaptiveMode bool   `yaml:"adaptive_mode"`
	LightTheme   string `yaml:"light_theme"`
	DarkTheme    string `yaml:"dark_theme"`
	FixedTheme   string `yaml:"fixed_theme,omitempty"`
}

// GetConfig 获取当前配置
func (m *Manager) GetConfig() Config {
	cfg := Config{
		AdaptiveMode: m.adaptiveMode,
		LightTheme:   m.lightScheme,
		DarkTheme:    m.darkScheme,
	}

	if !m.adaptiveMode {
		cfg.FixedTheme = m.currentName
	}

	return cfg
}

// ApplyConfig 应用配置
func (m *Manager) ApplyConfig(cfg Config) error {
	if cfg.AdaptiveMode {
		light := cfg.LightTheme
		dark := cfg.DarkTheme

		if light == "" {
			light = "cyberpunk-light"
		}
		if dark == "" {
			dark = "cyberpunk"
		}

		return m.SetAdaptiveTheme(light, dark)
	}

	if cfg.FixedTheme != "" {
		return m.SetTheme(cfg.FixedTheme)
	}

	return nil
}
