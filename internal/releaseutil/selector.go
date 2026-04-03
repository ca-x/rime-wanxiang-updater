package releaseutil

import "rime-wanxiang-updater/internal/types"

// FindPreferredAssetInfo 优先从版本化发布中选择匹配资源，
// 找不到时再回退到指定的兜底 tag。
func FindPreferredAssetInfo(
	releases []types.GitHubRelease,
	match func(name string) bool,
	fallbackTag string,
) (*types.UpdateInfo, bool) {
	if info, ok := findAssetInfoWithTagFilter(releases, match, func(tag string) bool {
		return fallbackTag == "" || tag != fallbackTag
	}); ok {
		return info, true
	}

	if fallbackTag == "" {
		return nil, false
	}

	return findAssetInfoWithTagFilter(releases, match, func(string) bool {
		return true
	})
}

// FindPreferredAssetName 使用相同的 tag 优先级规则选择资源文件名。
func FindPreferredAssetName(
	releases []types.GitHubRelease,
	match func(name string) bool,
	fallbackTag string,
) (string, bool) {
	info, ok := FindPreferredAssetInfo(releases, match, fallbackTag)
	if !ok {
		return "", false
	}

	return info.Name, true
}

// FindAssetInfoByTag 按指定 tag 精确匹配资源。
func FindAssetInfoByTag(
	releases []types.GitHubRelease,
	match func(name string) bool,
	tag string,
) (*types.UpdateInfo, bool) {
	if tag == "" {
		return nil, false
	}

	return findAssetInfoWithTagFilter(releases, match, func(currentTag string) bool {
		return currentTag == tag
	})
}

func findAssetInfoWithTagFilter(
	releases []types.GitHubRelease,
	match func(name string) bool,
	allowTag func(tag string) bool,
) (*types.UpdateInfo, bool) {
	for _, release := range releases {
		if !allowTag(release.TagName) {
			continue
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
			}, true
		}
	}

	return nil, false
}
