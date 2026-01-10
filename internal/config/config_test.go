package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

func TestGetActualFilenames(t *testing.T) {
	tests := []struct {
		name           string
		schemeKey      string
		useMirror      bool
		expectedScheme string
		expectedDict   string
	}{
		{
			name:           "base方案-CNB镜像",
			schemeKey:      "base",
			useMirror:      true,
			expectedScheme: "rime-wanxiang-base.zip",
			expectedDict:   "base-dicts.zip",
		},
		{
			name:           "base方案-GitHub",
			schemeKey:      "base",
			useMirror:      false,
			expectedScheme: "rime-wanxiang-base.zip",
			expectedDict:   "base-dicts.zip",
		},
		{
			name:           "墨奇方案-CNB镜像",
			schemeKey:      "moqi",
			useMirror:      true,
			expectedScheme: "rime-wanxiang-moqi-fuzhu.zip",
			expectedDict:   "pro-moqi-fuzhu-dicts.zip",
		},
		{
			name:           "小鹤方案-GitHub",
			schemeKey:      "flypy",
			useMirror:      false,
			expectedScheme: "rime-wanxiang-flypy-fuzhu.zip",
			expectedDict:   "pro-flypy-fuzhu-dicts.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &Manager{
				Config: &types.Config{
					UseMirror: tt.useMirror,
				},
			}

			schemeFile, dictFile, err := mgr.GetActualFilenames(tt.schemeKey)
			if err != nil {
				t.Logf("获取文件名失败（可能是网络问题）: %v", err)
				return
			}

			if schemeFile != tt.expectedScheme {
				t.Errorf("方案文件名错误：期望 %s，实际 %s", tt.expectedScheme, schemeFile)
			}

			if dictFile != tt.expectedDict {
				t.Errorf("词库文件名错误：期望 %s，实际 %s", tt.expectedDict, dictFile)
			}

			t.Logf("✓ 方案: %s, 词库: %s", schemeFile, dictFile)
		})
	}
}

func TestConfigMigration(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// 创建旧格式配置
	oldConfig := map[string]interface{}{
		"engine":       "鼠须管",
		"scheme_type":  "base",
		"use_mirror":   true,
		"github_token": "",
	}

	data, err := json.MarshalIndent(oldConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal old config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write old config: %v", err)
	}

	// 创建配置管理器
	m := &Manager{
		ConfigPath: configPath,
	}

	// 加载配置（应该触发迁移）
	config, err := m.loadOrCreateConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证迁移结果
	if len(config.InstalledEngines) == 0 {
		t.Error("Expected InstalledEngines to be populated after migration")
	}

	if config.PrimaryEngine == "" {
		t.Error("Expected PrimaryEngine to be set after migration")
	}

	// Engine 字段应该被清空（表示已迁移）
	if config.Engine != "" {
		t.Error("Expected Engine field to be cleared after migration")
	}
}

func TestGetEngineDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		config   *types.Config
		expected string
	}{
		{
			name: "Single engine",
			config: &types.Config{
				InstalledEngines: []string{"鼠须管"},
				PrimaryEngine:    "鼠须管",
			},
			expected: "鼠须管",
		},
		{
			name: "Multiple engines",
			config: &types.Config{
				InstalledEngines: []string{"鼠须管", "小企鹅"},
				PrimaryEngine:    "鼠须管",
			},
			expected: "鼠须管+小企鹅",
		},
		{
			name: "No installed engines",
			config: &types.Config{
				InstalledEngines: []string{},
				PrimaryEngine:    "鼠须管",
			},
			expected: "鼠须管",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				Config: tt.config,
			}
			result := m.GetEngineDisplayName()
			if result != tt.expected {
				t.Errorf("GetEngineDisplayName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRedetectEngines(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	m := &Manager{
		ConfigPath: configPath,
		Config: &types.Config{
			InstalledEngines: []string{"old_engine"},
			PrimaryEngine:    "old_engine",
		},
	}

	// 重新检测引擎
	err := m.RedetectEngines()
	if err != nil {
		t.Fatalf("RedetectEngines() error = %v", err)
	}

	// 验证引擎列表已更新
	if len(m.Config.InstalledEngines) == 0 {
		t.Error("Expected InstalledEngines to be populated")
	}

	// 验证主引擎已更新
	if m.Config.PrimaryEngine == "old_engine" {
		t.Error("Expected PrimaryEngine to be updated")
	}

	// 验证主引擎在已安装列表中
	found := false
	for _, engine := range m.Config.InstalledEngines {
		if engine == m.Config.PrimaryEngine {
			found = true
			break
		}
	}
	if !found {
		t.Error("PrimaryEngine should be in InstalledEngines")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	config := createDefaultConfig()

	if config == nil {
		t.Fatal("createDefaultConfig() returned nil")
	}

	if len(config.InstalledEngines) == 0 {
		t.Error("Expected InstalledEngines to be populated")
	}

	if config.PrimaryEngine == "" {
		t.Error("Expected PrimaryEngine to be set")
	}

	// 验证默认值
	if !config.UseMirror {
		t.Error("Expected UseMirror to be true by default")
	}

	if config.AutoUpdateCountdown != 5 {
		t.Error("Expected AutoUpdateCountdown to be 5 by default")
	}
}
