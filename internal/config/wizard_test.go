package config

import (
	"os"
	"path/filepath"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

// TestWizardConfigSave 测试向导保存配置的逻辑
func TestWizardConfigSave(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	tests := []struct {
		name           string
		schemeType     string
		useMirror      bool
		expectedScheme string
		expectedDict   string
	}{
		{
			name:           "base方案-CNB镜像",
			schemeType:     "base",
			useMirror:      true,
			expectedScheme: "rime-wanxiang-base.zip",
			expectedDict:   "base-dicts.zip",
		},
		{
			name:           "moqi方案-CNB镜像",
			schemeType:     "moqi",
			useMirror:      true,
			expectedScheme: "rime-wanxiang-moqi-fuzhu.zip",
			expectedDict:   "pro-moqi-fuzhu-dicts.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建配置管理器
			mgr := &Manager{
				ConfigPath: configPath,
				Config: &types.Config{
					Engine:       "fcitx5",
					SchemeType:   "",
					SchemeFile:   "",
					DictFile:     "",
					UseMirror:    tt.useMirror,
					GithubToken:  "",
					ExcludeFiles: []string{},
					AutoUpdate:   false,
					ProxyEnabled: false,
					ProxyType:    "http",
					ProxyAddress: "127.0.0.1:7890",
				},
			}

			// 模拟向导流程：获取文件名
			schemeFile, dictFile, err := mgr.GetActualFilenames(tt.schemeType)
			if err != nil {
				t.Logf("获取文件名失败（可能是网络问题）: %v", err)
				return
			}

			t.Logf("获取到的文件名:")
			t.Logf("  schemeFile = %s", schemeFile)
			t.Logf("  dictFile   = %s", dictFile)

			// 检查返回值
			if schemeFile != tt.expectedScheme {
				t.Errorf("方案文件名错误：期望 %s，实际 %s", tt.expectedScheme, schemeFile)
			}
			if dictFile != tt.expectedDict {
				t.Errorf("词库文件名错误：期望 %s，实际 %s", tt.expectedDict, dictFile)
			}

			// 模拟向导保存配置
			mgr.Config.SchemeType = tt.schemeType
			mgr.Config.SchemeFile = schemeFile
			mgr.Config.DictFile = dictFile

			// 保存配置
			if err := mgr.saveConfig(mgr.Config); err != nil {
				t.Fatalf("保存配置失败: %v", err)
			}

			// 重新加载配置
			loaded, err := mgr.loadOrCreateConfig()
			if err != nil {
				t.Fatalf("加载配置失败: %v", err)
			}

			t.Logf("保存并重新加载后的配置:")
			t.Logf("  SchemeType = %s", loaded.SchemeType)
			t.Logf("  SchemeFile = %s", loaded.SchemeFile)
			t.Logf("  DictFile   = %s", loaded.DictFile)

			// 验证保存的配置
			if loaded.SchemeType != tt.schemeType {
				t.Errorf("SchemeType 错误：期望 %s，实际 %s", tt.schemeType, loaded.SchemeType)
			}
			if loaded.SchemeFile != tt.expectedScheme {
				t.Errorf("SchemeFile 错误：期望 %s，实际 %s", tt.expectedScheme, loaded.SchemeFile)
			}
			if loaded.DictFile != tt.expectedDict {
				t.Errorf("DictFile 错误：期望 %s，实际 %s", tt.expectedDict, loaded.DictFile)
			}

			// 最关键的检查：确保 DictFile 不等于 SchemeFile
			if loaded.DictFile == loaded.SchemeFile {
				t.Errorf("❌ 严重错误：DictFile 和 SchemeFile 相同！DictFile=%s, SchemeFile=%s",
					loaded.DictFile, loaded.SchemeFile)
			} else {
				t.Logf("✓ DictFile 和 SchemeFile 不同（正确）")
			}

			// 清理
			os.Remove(configPath)
		})
	}
}
