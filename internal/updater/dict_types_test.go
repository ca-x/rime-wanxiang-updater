package updater

import (
	"testing"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

func TestDictUpdaterDifferentTypes(t *testing.T) {
	// 测试不同的辅助码类型
	dictTypes := map[string]string{
		"base":    "base-dicts.zip",
		"moqi":    "pro-moqi-fuzhu-dicts.zip",
		"flypy":   "pro-flypy-fuzhu-dicts.zip",
		"zrm":     "pro-zrm-fuzhu-dicts.zip",
		"tiger":   "pro-tiger-fuzhu-dicts.zip",
		"wubi":    "pro-wubi-fuzhu-dicts.zip",
		"hanxin":  "pro-hanxin-fuzhu-dicts.zip",
		"shouyou": "pro-shouyou-fuzhu-dicts.zip",
	}

	for schemeType, expectedFile := range dictTypes {
		t.Run(schemeType, func(t *testing.T) {
			cfg := &config.Manager{
				Config: &types.Config{
					UseMirror:  true,
					DictFile:   expectedFile,
					SchemeType: schemeType,
				},
			}

			updater := NewDictUpdater(cfg)

			// 测试获取更新信息
			info, err := updater.CheckUpdate()
			if err != nil {
				t.Logf("获取更新信息失败（这可能是网络问题）: %v", err)
				return
			}

			// 验证返回的信息
			if info == nil {
				t.Fatal("更新信息为空")
			}

			t.Logf("方案类型: %s", schemeType)
			t.Logf("词库名称: %s", info.Name)
			t.Logf("词库标签: %s", info.Tag)
			t.Logf("更新时间: %s", info.UpdateTime)

			// 验证文件名匹配
			if info.Name != expectedFile {
				t.Errorf("词库文件名错误，期望 %s，实际 %s", expectedFile, info.Name)
			}

			// 验证 tag 正确
			if info.Tag != types.CNB_DICT_TAG {
				t.Errorf("词库 tag 错误，期望 %s，实际 %s", types.CNB_DICT_TAG, info.Tag)
			}
		})
	}
}
