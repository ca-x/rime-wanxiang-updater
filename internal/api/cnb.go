package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudflare/backoff"
	"rime-wanxiang-updater/internal/types"
)

// FetchCNBReleases 获取 CNB Releases
func (c *Client) FetchCNBReleases(owner, repo, tag string) ([]types.GitHubRelease, error) {
	url := fmt.Sprintf("https://cnb.cool/%s/%s/-/releases", owner, repo)

	b := backoff.New(time.Second, 10*time.Second)
	var resp *http.Response
	var err error
	attempts := 0
	maxAttempts := 3

	// 使用重试机制
	for attempts < maxAttempts {
		attempts++
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/vnd.cnb.web+json")

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempts < maxAttempts {
			time.Sleep(b.Duration())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("CNB 请求失败（尝试 %d 次）: %w", attempts, err)
	}
	defer resp.Body.Close()

	var result struct {
		Releases []types.CNBRelease `json:"releases"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 转换 CNB Release 为 GitHub Release 格式
	var releases []types.GitHubRelease
	for _, cnbRelease := range result.Releases {
		tagName := cnbRelease.TagRef
		if strings.Contains(tagName, "/") {
			parts := strings.Split(tagName, "/")
			tagName = parts[len(parts)-1]
		}

		// 如果指定了 tag，则过滤
		if tag != "" && tagName != tag {
			continue
		}

		var assets []types.GitHubAsset
		for _, cnbAsset := range cnbRelease.Assets {
			assets = append(assets, types.GitHubAsset{
				Name:               cnbAsset.Name,
				BrowserDownloadURL: "https://cnb.cool" + cnbAsset.Path,
				UpdatedAt:          cnbAsset.UpdatedAt,
				Size:               cnbAsset.SizeInByte,
			})
		}

		releases = append(releases, types.GitHubRelease{
			TagName: tagName,
			Body:    cnbRelease.Body,
			Assets:  assets,
		})
	}

	return releases, nil
}
