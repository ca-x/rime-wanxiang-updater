package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudflare/backoff"
	"rime-wanxiang-updater/internal/types"
)

// FetchGitHubReleases 获取 GitHub Releases
func (c *Client) FetchGitHubReleases(owner, repo, tag string) ([]types.GitHubRelease, error) {
	var url string
	if tag != "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	} else {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	}

	// 使用重试机制
	resp, err := c.fetchWithRetry(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if tag != "" {
		var release types.GitHubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return nil, fmt.Errorf("解析响应失败: %w", err)
		}
		return []types.GitHubRelease{release}, nil
	}

	var releases []types.GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return releases, nil
}

// fetchWithRetry 带重试的 HTTP 请求，使用 cloudflare backoff
func (c *Client) fetchWithRetry(url string) (*http.Response, error) {
	b := backoff.New(time.Second, 10*time.Second)
	var resp *http.Response
	var err error
	attempts := 0
	maxAttempts := 3

	for attempts < maxAttempts {
		attempts++
		resp, err = c.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempts < maxAttempts {
			time.Sleep(b.Duration())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("请求失败（尝试 %d 次）: %w", attempts, err)
	}
	return nil, fmt.Errorf("请求失败（尝试 %d 次），状态码: %d", attempts, resp.StatusCode)
}
