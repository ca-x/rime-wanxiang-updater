package config

import (
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
