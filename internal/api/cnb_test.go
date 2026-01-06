package api

import (
	"fmt"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

// getTestConfig 获取测试用配置
func getTestConfig() *types.Config {
	return &types.Config{
		Engine:       "fcitx5",
		SchemeType:   "base",
		UseMirror:    true,
		ProxyEnabled: false,
		GithubToken:  "",
	}
}

// TestFetchCNBReleases 测试 CNB API 获取 releases
func TestFetchCNBReleases(t *testing.T) {
	client := NewClient(getTestConfig())

	// 测试获取 rime-wanxiang 仓库的所有 releases
	releases, err := client.FetchCNBReleases("amzxyz", "rime-wanxiang", "")
	if err != nil {
		t.Fatalf("获取 CNB releases 失败: %v", err)
	}

	if len(releases) == 0 {
		t.Fatal("未获取到任何 releases")
	}

	fmt.Printf("总共获取到 %d 个 releases\n", len(releases))

	// 打印所有 releases 的信息
	for i, release := range releases {
		fmt.Printf("\n=== Release %d ===\n", i+1)
		fmt.Printf("Tag: %s\n", release.TagName)
		fmt.Printf("Assets count: %d\n", len(release.Assets))

		// 打印所有 assets
		for j, asset := range release.Assets {
			fmt.Printf("  Asset %d: %s (size: %d)\n", j+1, asset.Name, asset.Size)
		}
	}

	// 查找包含 "model" tag 的 release
	var modelRelease *types.GitHubRelease
	for i := range releases {
		if releases[i].TagName == "model" {
			modelRelease = &releases[i]
			break
		}
	}

	if modelRelease == nil {
		t.Error("未找到 tag 为 'model' 的 release")
		return
	}

	fmt.Printf("\n=== Model Release ===\n")
	fmt.Printf("Tag: %s\n", modelRelease.TagName)
	fmt.Printf("Assets: %d\n", len(modelRelease.Assets))

	// 查找 .gram 文件
	var foundGram bool
	for _, asset := range modelRelease.Assets {
		if asset.Name == "wanxiang-lts-zh-hans.gram" {
			foundGram = true
			fmt.Printf("找到模型文件: %s\n", asset.Name)
			fmt.Printf("URL: %s\n", asset.BrowserDownloadURL)
			fmt.Printf("Size: %d bytes\n", asset.Size)
			break
		}
	}

	if !foundGram {
		t.Error("在 model release 中未找到 wanxiang-lts-zh-hans.gram 文件")
		fmt.Println("\n可用的文件列表:")
		for _, asset := range modelRelease.Assets {
			fmt.Printf("  - %s\n", asset.Name)
		}
	}
}

// TestFetchCNBReleasesDebug 调试用测试，打印详细信息
func TestFetchCNBReleasesDebug(t *testing.T) {
	client := NewClient(getTestConfig())

	releases, err := client.FetchCNBReleases("amzxyz", "rime-wanxiang", "v1.0.0")
	if err != nil {
		t.Fatalf("获取 CNB releases 失败: %v", err)
	}

	fmt.Printf("\n=== 所有 Releases 详细信息 ===\n")
	for i, release := range releases {
		fmt.Printf("\nRelease %d:\n", i+1)
		fmt.Printf("  TagName: '%s'\n", release.TagName)
		fmt.Printf("  Assets count: %d\n", len(release.Assets))

		if len(release.Assets) > 0 {
			fmt.Printf("  Assets:\n")
			for _, asset := range release.Assets {
				fmt.Printf("    - Name: '%s'\n", asset.Name)
				fmt.Printf("      URL: '%s'\n", asset.BrowserDownloadURL)
				fmt.Printf("      Size: %d\n", asset.Size)
			}
		}
	}
}
