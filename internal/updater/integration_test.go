package updater

import (
	"testing"
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

// TestFetchAllUpdatesPerformance 测试完整的检查流程性能
func TestFetchAllUpdatesPerformance(t *testing.T) {
	// 创建临时配置
	cfg := &config.Manager{
		Config: &types.Config{
			Engine:       "测试",
			SchemeType:   "base",
			SchemeFile:   "rime-wanxiang-base.zip",
			DictFile:     "base-dicts.zip",
			UseMirror:    true,
			GithubToken:  "",
			ProxyEnabled: false,
		},
	}

	// 创建组合更新器
	combined := NewCombinedUpdater(cfg)

	// 测试完整的检查流程
	t.Log("开始检查所有更新...")
	start := time.Now()

	// 单独测试每个组件
	t.Log("\n=== 测试方案检查 ===")
	schemeStart := time.Now()
	schemeInfo, schemeErr := combined.SchemeUpdater.CheckUpdate()
	schemeElapsed := time.Since(schemeStart)
	t.Logf("方案检查耗时: %v", schemeElapsed)
	if schemeErr != nil {
		t.Logf("方案检查错误: %v", schemeErr)
	} else if schemeInfo != nil {
		t.Logf("方案信息: %s, Tag: %s", schemeInfo.Name, schemeInfo.Tag)
	}

	t.Log("\n=== 测试词库检查 ===")
	dictStart := time.Now()
	dictInfo, dictErr := combined.DictUpdater.CheckUpdate()
	dictElapsed := time.Since(dictStart)
	t.Logf("词库检查耗时: %v", dictElapsed)
	if dictErr != nil {
		t.Logf("词库检查错误: %v", dictErr)
	} else if dictInfo != nil {
		t.Logf("词库信息: %s, Tag: %s", dictInfo.Name, dictInfo.Tag)
	}

	t.Log("\n=== 测试模型检查 ===")
	modelStart := time.Now()
	modelInfo, modelErr := combined.ModelUpdater.CheckUpdate()
	modelElapsed := time.Since(modelStart)
	t.Logf("模型检查耗时: %v", modelElapsed)
	if modelErr != nil {
		t.Logf("模型检查错误: %v", modelErr)
	} else if modelInfo != nil {
		t.Logf("模型信息: %s, Tag: %s", modelInfo.Name, modelInfo.Tag)
	}

	t.Log("\n=== 测试完整流程 ===")
	fullStart := time.Now()
	err := combined.FetchAllUpdates()
	fullElapsed := time.Since(fullStart)

	totalElapsed := time.Since(start)

	// 输出结果
	t.Logf("\n完整流程耗时: %v", fullElapsed)
	t.Logf("总耗时: %v", totalElapsed)
	t.Logf("方案: %v | 词库: %v | 模型: %v", schemeElapsed, dictElapsed, modelElapsed)

	if err != nil {
		t.Logf("检查更新出现错误: %v", err)
	} else {
		t.Log("检查更新成功")
	}

	// 验证总耗时不应该太长（超过30秒就有问题）
	if totalElapsed > 30*time.Second {
		t.Errorf("总耗时过长: %v，可能有网络问题或卡住", totalElapsed)
	}
}

// TestIndividualComponentPerformance 测试各组件独立性能
func TestIndividualComponentPerformance(t *testing.T) {
	cfg := &config.Manager{
		Config: &types.Config{
			Engine:       "测试",
			SchemeType:   "base",
			SchemeFile:   "rime-wanxiang-base.zip",
			DictFile:     "base-dicts.zip",
			UseMirror:    true,
			GithubToken:  "",
			ProxyEnabled: false,
		},
	}

	tests := []struct {
		name    string
		timeout time.Duration
		test    func() (interface{}, error)
	}{
		{
			name:    "方案检查",
			timeout: 15 * time.Second,
			test: func() (interface{}, error) {
				return NewSchemeUpdater(cfg).CheckUpdate()
			},
		},
		{
			name:    "词库检查",
			timeout: 15 * time.Second,
			test: func() (interface{}, error) {
				return NewDictUpdater(cfg).CheckUpdate()
			},
		},
		{
			name:    "模型检查",
			timeout: 1 * time.Second, // 模型应该很快
			test: func() (interface{}, error) {
				return NewModelUpdater(cfg).CheckUpdate()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool)
			var elapsed time.Duration

			go func() {
				start := time.Now()
				result, err := tt.test()
				elapsed = time.Since(start)

				if err != nil {
					t.Logf("%s 错误: %v", tt.name, err)
				} else {
					t.Logf("%s 成功，耗时: %v", tt.name, elapsed)
					if result != nil {
						t.Logf("%s 结果: %+v", tt.name, result)
					}
				}
				done <- true
			}()

			select {
			case <-done:
				if elapsed > tt.timeout {
					t.Errorf("%s 超时: 耗时 %v，超过限制 %v", tt.name, elapsed, tt.timeout)
				}
			case <-time.After(tt.timeout + 5*time.Second):
				t.Errorf("%s 严重超时: 超过 %v + 5秒", tt.name, tt.timeout)
			}
		})
	}
}
