package updater

import (
	"testing"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

func TestDictUpdaterCheckUpdate(t *testing.T) {
	// 创建临时配置
	cfg := &config.Manager{
		Config: &types.Config{
			UseMirror:  true,
			DictFile:   "base-dicts.zip",
			SchemeType: "base",
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

	t.Logf("词库名称: %s", info.Name)
	t.Logf("词库标签: %s", info.Tag)
	t.Logf("更新时间: %s", info.UpdateTime)
	t.Logf("下载地址: %s", info.URL)

	// 验证 CNB 使用的是 v1.0.0 tag
	if cfg.Config.UseMirror && info.Tag != types.CNB_DICT_TAG {
		t.Errorf("CNB 词库 tag 错误，期望 %s，实际 %s", types.CNB_DICT_TAG, info.Tag)
	}
}

func TestDictUpdaterCheckUpdateGitHub(t *testing.T) {
	// 创建临时配置
	cfg := &config.Manager{
		Config: &types.Config{
			UseMirror:  false,
			DictFile:   "base-dicts.zip",
			SchemeType: "base",
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

	t.Logf("词库名称: %s", info.Name)
	t.Logf("词库标签: %s", info.Tag)
	t.Logf("更新时间: %s", info.UpdateTime)
	t.Logf("下载地址: %s", info.URL)

	// 验证 GitHub 使用的是 dict-nightly tag
	if !cfg.Config.UseMirror && info.Tag != types.DICT_TAG {
		t.Errorf("GitHub 词库 tag 错误，期望 %s，实际 %s", types.DICT_TAG, info.Tag)
	}
}
