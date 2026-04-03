package updater

import (
	"strings"
	"testing"
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

// TestModelUpdaterCheckUpdate 测试模型更新检查
func TestModelUpdaterCheckUpdate(t *testing.T) {
	tests := []struct {
		name      string
		useMirror bool
		wantErr   bool
	}{
		{
			name:      "CNB 镜像模式",
			useMirror: true,
			wantErr:   false,
		},
		{
			name:      "GitHub 模式",
			useMirror: false,
			wantErr:   false, // 可能失败，取决于网络
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时配置
			cfg := &config.Manager{
				Config: &types.Config{
					Engine:       "测试",
					UseMirror:    tt.useMirror,
					GithubToken:  "",
					ProxyEnabled: false,
				},
			}

			// 创建模型更新器
			updater := NewModelUpdater(cfg)

			// 测试 CheckUpdate
			start := time.Now()
			info, err := updater.CheckUpdate()
			elapsed := time.Since(start)

			// 验证结果
			if tt.wantErr && err == nil {
				t.Errorf("期望出错，但成功了")
			}
			if tt.useMirror && err != nil && strings.Contains(err.Error(), "状态码: 429") {
				t.Skipf("CNB rate limited this live test: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("不期望出错，但失败了: %v", err)
			}

			// 验证响应时间（CNB 镜像模式应该很快）
			if tt.useMirror && elapsed > 5*time.Second {
				t.Errorf("CNB 镜像模式响应时间过长: %v", elapsed)
			}

			// 验证返回的信息
			if !tt.wantErr && info != nil {
				t.Logf("模式: %s", tt.name)
				t.Logf("文件名: %s", info.Name)
				t.Logf("URL: %s", info.URL)
				t.Logf("更新时间: %s", info.UpdateTime)
				t.Logf("Tag: %s", info.Tag)
				t.Logf("响应时间: %v", elapsed)

				if info.Name != types.MODEL_FILE {
					t.Errorf("文件名不匹配，期望 %s，得到 %s", types.MODEL_FILE, info.Name)
				}

				if info.URL == "" {
					t.Errorf("URL 不应为空")
				}

				if tt.useMirror {
					minExpectedTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
					if !info.UpdateTime.After(minExpectedTime) {
						t.Errorf("CNB 镜像模式时间过旧，期望晚于 %v，得到 %v", minExpectedTime, info.UpdateTime)
					}

					if info.Tag != "model" {
						t.Errorf("CNB 镜像模式 Tag = %q, want %q", info.Tag, "model")
					}

					if info.Size <= 0 {
						t.Errorf("CNB 镜像模式 Size = %d, want > 0", info.Size)
					}
				}
			}
		})
	}
}

// TestModelUpdaterCheckUpdateTimeout 测试超时行为
func TestModelUpdaterCheckUpdateTimeout(t *testing.T) {
	cfg := &config.Manager{
		Config: &types.Config{
			Engine:       "测试",
			UseMirror:    true,
			GithubToken:  "",
			ProxyEnabled: false,
		},
	}

	updater := NewModelUpdater(cfg)

	// 运行多次确保稳定性
	for i := 0; i < 3; i++ {
		start := time.Now()
		info, err := updater.CheckUpdate()
		elapsed := time.Since(start)

		if err != nil {
			if strings.Contains(err.Error(), "状态码: 429") {
				t.Skipf("CNB rate limited this live test: %v", err)
			}
			t.Errorf("第 %d 次检查失败: %v", i+1, err)
		}

		if info == nil {
			t.Errorf("第 %d 次检查返回 nil", i+1)
		}

		// CNB 镜像模式应该立即返回（不超过1秒）
		if elapsed > time.Second {
			t.Errorf("第 %d 次检查响应时间过长: %v", i+1, elapsed)
		}

		t.Logf("第 %d 次检查耗时: %v", i+1, elapsed)
	}
}

func TestFindModelReleaseUsesModelTagMetadata(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: "v1.0.0",
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/base-dicts.zip",
					UpdatedAt:          time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
					ID:                 "dict-asset",
					Size:               11,
				},
			},
		},
		{
			TagName: "model",
			Assets: []types.GitHubAsset{
				{
					Name:               types.MODEL_FILE,
					BrowserDownloadURL: "https://example.com/model.gram",
					UpdatedAt:          time.Date(2026, 4, 2, 7, 44, 10, 0, time.UTC),
					ID:                 "model-asset",
					SHA256:             "abc123",
					Size:               210421804,
				},
			},
		},
	}

	info, ok := findModelRelease(releases)
	if !ok {
		t.Fatal("findModelRelease() = no match, want match")
	}

	if info.Tag != "model" {
		t.Fatalf("info.Tag = %q, want %q", info.Tag, "model")
	}

	if info.Size != 210421804 {
		t.Fatalf("info.Size = %d, want %d", info.Size, 210421804)
	}

	if info.ID != "model-asset" {
		t.Fatalf("info.ID = %q, want %q", info.ID, "model-asset")
	}
}
