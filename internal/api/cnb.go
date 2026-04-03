package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cloudflare/backoff"
	"rime-wanxiang-updater/internal/types"
)

// FetchCNBReleases 获取 CNB Releases
func (c *Client) FetchCNBReleases(owner, repo, tag string) ([]types.GitHubRelease, error) {
	var releases []types.GitHubRelease

	baseURL := fmt.Sprintf(
		"%s/%s/%s/-/releases",
		strings.TrimRight(c.cnbBaseURL, "/"),
		owner,
		repo,
	)

	totalPages := 1
	for page := 1; page <= totalPages; page++ {
		pageURL, err := c.cnbReleasesPageURL(baseURL, page)
		if err != nil {
			return nil, fmt.Errorf("构建 CNB 分页链接失败: %w", err)
		}

		resp, err := c.fetchCNBPageWithRetry(pageURL)
		if err != nil {
			return nil, err
		}

		if page == 1 {
			totalPages = cnbTotalPages(resp.Header)
		}

		var result struct {
			Releases []types.CNBRelease `json:"releases"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("解析响应失败: %w", err)
		}
		resp.Body.Close()

		releases = append(releases, convertCNBReleases(result.Releases, tag, c.cnbBaseURL)...)
		if tag != "" && len(releases) > 0 {
			break
		}
	}

	return releases, nil
}

// FetchCNBReleaseByTag 获取指定 tag 的单个 CNB release。
func (c *Client) FetchCNBReleaseByTag(owner, repo, tag string) (*types.GitHubRelease, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag 不能为空")
	}

	cacheKey := fmt.Sprintf("%s/%s:%s", owner, repo, tag)
	c.mu.Lock()
	if cached, ok := c.cnbReleaseTagCache[cacheKey]; ok {
		c.mu.Unlock()
		releaseCopy := cached
		return &releaseCopy, nil
	}
	c.mu.Unlock()

	rawURL := fmt.Sprintf(
		"%s/%s/%s/-/releases/tags/%s",
		strings.TrimRight(c.cnbBaseURL, "/"),
		owner,
		repo,
		url.PathEscape(tag),
	)

	resp, err := c.fetchCNBPageWithRetry(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Release types.CNBRelease `json:"release"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	releases := convertCNBReleases([]types.CNBRelease{result.Release}, tag, c.cnbBaseURL)
	if len(releases) == 0 {
		return nil, fmt.Errorf("未找到 tag 为 %s 的 CNB release", tag)
	}

	c.mu.Lock()
	c.cnbReleaseTagCache[cacheKey] = releases[0]
	c.mu.Unlock()

	releaseCopy := releases[0]
	return &releaseCopy, nil
}

// FetchCNBReleaseTagsPage 获取一页带 release 的 CNB tag 列表。
func (c *Client) FetchCNBReleaseTagsPage(owner, repo string, page int) ([]string, int, error) {
	if page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("%s/%s:%d", owner, repo, page)
	c.mu.Lock()
	if cached, ok := c.cnbTagsPageCache[cacheKey]; ok {
		c.mu.Unlock()
		return append([]string(nil), cached.tags...), cached.totalPages, nil
	}
	c.mu.Unlock()

	rawURL := fmt.Sprintf(
		"%s/%s/%s/-/git/tags",
		strings.TrimRight(c.cnbBaseURL, "/"),
		owner,
		repo,
	)
	if page > 1 {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			return nil, 0, fmt.Errorf("解析 tags 链接失败: %w", err)
		}

		query := parsedURL.Query()
		query.Set("page", strconv.Itoa(page))
		parsedURL.RawQuery = query.Encode()
		rawURL = parsedURL.String()
	}

	resp, err := c.fetchCNBPageWithRetry(rawURL)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Tags []struct {
			Tag        string `json:"tag"`
			HasRelease bool   `json:"has_release"`
		} `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("解析响应失败: %w", err)
	}

	var tags []string
	for _, item := range result.Tags {
		if !item.HasRelease {
			continue
		}
		tags = append(tags, normalizeCNBTag(item.Tag))
	}

	totalPages := cnbTotalPages(resp.Header)

	c.mu.Lock()
	c.cnbTagsPageCache[cacheKey] = cnbTagsPageResult{
		tags:       append([]string(nil), tags...),
		totalPages: totalPages,
	}
	c.mu.Unlock()

	return tags, totalPages, nil
}

// FindLatestCNBAssetInfo 根据 tag 列表查找最新匹配资源，必要时回退到指定 tag。
func (c *Client) FindLatestCNBAssetInfo(
	owner string,
	repo string,
	match func(string) bool,
	fallbackTag string,
) (*types.UpdateInfo, error) {
	totalPages := 1
	for page := 1; page <= totalPages; page++ {
		tags, pages, err := c.FetchCNBReleaseTagsPage(owner, repo, page)
		if err != nil {
			return nil, err
		}
		if page == 1 {
			totalPages = pages
		}

		for _, tag := range tags {
			if fallbackTag != "" && tag == fallbackTag {
				continue
			}

			release, err := c.FetchCNBReleaseByTag(owner, repo, tag)
			if err != nil {
				return nil, err
			}

			for _, asset := range release.Assets {
				if !match(asset.Name) {
					continue
				}

				return &types.UpdateInfo{
					Name:        asset.Name,
					URL:         asset.BrowserDownloadURL,
					UpdateTime:  asset.UpdatedAt,
					Tag:         release.TagName,
					Description: release.Body,
					SHA256:      asset.SHA256,
					ID:          asset.ID,
					Size:        asset.Size,
				}, nil
			}
		}
	}

	if fallbackTag != "" {
		release, err := c.FetchCNBReleaseByTag(owner, repo, fallbackTag)
		if err != nil {
			return nil, err
		}

		for _, asset := range release.Assets {
			if !match(asset.Name) {
				continue
			}

			return &types.UpdateInfo{
				Name:        asset.Name,
				URL:         asset.BrowserDownloadURL,
				UpdateTime:  asset.UpdatedAt,
				Tag:         release.TagName,
				Description: release.Body,
				SHA256:      asset.SHA256,
				ID:          asset.ID,
				Size:        asset.Size,
			}, nil
		}
	}

	return nil, fmt.Errorf("未找到匹配的 CNB 资源")
}

// FetchLatestCNBReleaseTag 获取最新的带 release 的 tag。
func (c *Client) FetchLatestCNBReleaseTag(owner, repo string) (string, error) {
	tags, _, err := c.FetchCNBReleaseTagsPage(owner, repo, 1)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("未找到任何带 release 的 CNB tag")
	}

	return tags[0], nil
}

func (c *Client) fetchCNBPageWithRetry(rawURL string) (*http.Response, error) {
	b := backoff.New(500*time.Millisecond, 2*time.Second)
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= 2; attempt++ {
		req, reqErr := http.NewRequest(http.MethodGet, rawURL, nil)
		if reqErr != nil {
			return nil, fmt.Errorf("创建 CNB 请求失败: %w", reqErr)
		}

		req.Header.Set("Accept", "application/vnd.cnb.web+json")

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt < 2 {
			time.Sleep(b.Duration())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("CNB 请求失败: %w", err)
	}

	return nil, fmt.Errorf("CNB 请求失败，状态码: %d", resp.StatusCode)
}

func (c *Client) cnbReleasesPageURL(baseURL string, page int) (string, error) {
	if page <= 1 {
		return baseURL, nil
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	query.Set("page", strconv.Itoa(page))
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func cnbTotalPages(header http.Header) int {
	total, err := strconv.Atoi(header.Get("X-CNB-Total"))
	if err != nil || total <= 0 {
		return 1
	}

	pageSize, err := strconv.Atoi(header.Get("X-CNB-Page-Size"))
	if err != nil || pageSize <= 0 {
		return 1
	}

	pages := (total + pageSize - 1) / pageSize
	if pages < 1 {
		return 1
	}

	return pages
}

func convertCNBReleases(
	cnbReleases []types.CNBRelease,
	tag string,
	baseURL string,
) []types.GitHubRelease {
	var releases []types.GitHubRelease

	for _, cnbRelease := range cnbReleases {
		tagName := cnbRelease.TagRef
		if strings.Contains(tagName, "/") {
			parts := strings.Split(tagName, "/")
			tagName = parts[len(parts)-1]
		}

		if tag != "" && tagName != tag {
			continue
		}

		var assets []types.GitHubAsset
		for _, cnbAsset := range cnbRelease.Assets {
			sha256 := ""
			if strings.EqualFold(cnbAsset.HashAlgo, "sha256") {
				sha256 = cnbAsset.HashValue
			}

			assets = append(assets, types.GitHubAsset{
				Name:               cnbAsset.Name,
				BrowserDownloadURL: strings.TrimRight(baseURL, "/") + cnbAsset.Path,
				UpdatedAt:          cnbAsset.UpdatedAt,
				ID:                 cnbAsset.ID,
				SHA256:             sha256,
				Size:               cnbAsset.SizeInByte,
			})
		}

		releases = append(releases, types.GitHubRelease{
			TagName: tagName,
			Body:    cnbRelease.Body,
			Assets:  assets,
		})
	}

	return releases
}

func normalizeCNBTag(tag string) string {
	if strings.Contains(tag, "/") {
		parts := strings.Split(tag, "/")
		return parts[len(parts)-1]
	}

	return tag
}
