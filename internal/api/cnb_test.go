package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

type rewriteHostTransport struct {
	target *url.URL
	base   http.RoundTripper
}

func (t rewriteHostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.URL.Scheme = t.target.Scheme
	cloned.URL.Host = t.target.Host
	cloned.Host = t.target.Host

	return t.base.RoundTrip(cloned)
}

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

func TestFetchCNBReleasesPaginatesUntilTagFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-CNB-Total", "2")
		w.Header().Set("X-CNB-Page-Size", "1")

		switch page {
		case "1":
			fmt.Fprint(w, `{
				"releases": [
					{
						"tag_ref": "refs/tags/v1.0.0",
						"assets": [
							{
								"name": "base-dicts.zip",
								"path": "/assets/base-dicts.zip",
								"updated_at": "2026-04-01T00:00:00Z",
								"id": "dict-asset",
								"size_in_byte": 11
							}
						]
					}
				]
			}`)
		case "2":
			fmt.Fprint(w, `{
				"releases": [
					{
						"tag_ref": "refs/tags/model",
						"assets": [
							{
								"name": "wanxiang-lts-zh-hans.gram",
								"path": "/assets/wanxiang-lts-zh-hans.gram",
								"updated_at": "2026-04-02T07:44:10Z",
								"id": "model-asset",
								"size_in_byte": 210421804
							}
						]
					}
				]
			}`)
		default:
			t.Fatalf("unexpected page query: %q", page)
		}
	}))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse(server.URL) error = %v", err)
	}

	client := NewClient(getTestConfig())
	client.httpClient = &http.Client{
		Transport: rewriteHostTransport{
			target: targetURL,
			base:   http.DefaultTransport,
		},
	}

	releases, err := client.FetchCNBReleases("amzxyz", "rime-wanxiang", "model")
	if err != nil {
		t.Fatalf("FetchCNBReleases() error = %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("len(releases) = %d, want 1", len(releases))
	}

	release := releases[0]
	if release.TagName != "model" {
		t.Fatalf("release.TagName = %q, want %q", release.TagName, "model")
	}

	if len(release.Assets) != 1 {
		t.Fatalf("len(release.Assets) = %d, want 1", len(release.Assets))
	}

	asset := release.Assets[0]
	if asset.Name != types.MODEL_FILE {
		t.Fatalf("asset.Name = %q, want %q", asset.Name, types.MODEL_FILE)
	}

	if asset.Size != 210421804 {
		t.Fatalf("asset.Size = %d, want %d", asset.Size, 210421804)
	}
}

func TestFindLatestCNBAssetInfoUsesTagAndReleaseEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/amzxyz/rime-wanxiang/-/git/tags":
			w.Header().Set("X-CNB-Total", "2")
			w.Header().Set("X-CNB-Page-Size", "10")
			fmt.Fprint(w, `{
				"tags": [
					{"tag": "refs/tags/v15.5.0", "has_release": true},
					{"tag": "refs/tags/v15.4.5", "has_release": true}
				]
			}`)
		case "/amzxyz/rime-wanxiang/-/releases/tags/v15.5.0":
			fmt.Fprint(w, `{
				"release": {
					"tag_ref": "refs/tags/v15.5.0",
					"assets": [
						{
							"name": "rime-wanxiang-base.zip",
							"path": "/assets/rime-wanxiang-base.zip",
							"updated_at": "2026-04-03T08:00:00Z",
							"id": "scheme-asset",
							"size_in_byte": 33885033
						}
					]
				}
			}`)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse(server.URL) error = %v", err)
	}

	client := NewClient(getTestConfig())
	client.cnbBaseURL = server.URL
	client.httpClient = &http.Client{
		Transport: rewriteHostTransport{
			target: targetURL,
			base:   http.DefaultTransport,
		},
	}

	info, err := client.FindLatestCNBAssetInfo(
		"amzxyz",
		"rime-wanxiang",
		func(name string) bool { return name == "rime-wanxiang-base.zip" },
		"v1.0.0",
	)
	if err != nil {
		t.Fatalf("FindLatestCNBAssetInfo() error = %v", err)
	}

	if info.Tag != "v15.5.0" {
		t.Fatalf("info.Tag = %q, want %q", info.Tag, "v15.5.0")
	}

	if info.Name != "rime-wanxiang-base.zip" {
		t.Fatalf("info.Name = %q, want %q", info.Name, "rime-wanxiang-base.zip")
	}
}
