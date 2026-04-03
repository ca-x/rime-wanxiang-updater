package releaseutil

import (
	"testing"
	"time"

	"rime-wanxiang-updater/internal/types"
)

func TestFindPreferredAssetInfoPrefersVersionedRelease(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/fallback.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 8, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			TagName: "v15.6.0",
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/v15.6.0.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	info, ok := FindPreferredAssetInfo(releases, func(name string) bool {
		return name == "base-dicts.zip"
	}, types.CNB_DICT_TAG)
	if !ok {
		t.Fatal("FindPreferredAssetInfo() = no match, want match")
	}

	if info.Tag != "v15.6.0" {
		t.Fatalf("FindPreferredAssetInfo().Tag = %q, want %q", info.Tag, "v15.6.0")
	}
}

func TestFindPreferredAssetInfoFallsBackToRollingRelease(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/fallback.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 8, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	info, ok := FindPreferredAssetInfo(releases, func(name string) bool {
		return name == "base-dicts.zip"
	}, types.CNB_DICT_TAG)
	if !ok {
		t.Fatal("FindPreferredAssetInfo() = no match, want match")
	}

	if info.Tag != types.CNB_DICT_TAG {
		t.Fatalf("FindPreferredAssetInfo().Tag = %q, want %q", info.Tag, types.CNB_DICT_TAG)
	}
}

func TestFindAssetInfoByTagMatchesExactTag(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: "v15.6.0",
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/v15.6.0.zip",
				},
			},
		},
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/v1.0.0.zip",
				},
			},
		},
	}

	info, ok := FindAssetInfoByTag(releases, func(name string) bool {
		return name == "base-dicts.zip"
	}, "v15.6.0")
	if !ok {
		t.Fatal("FindAssetInfoByTag() = no match, want match")
	}

	if info.Tag != "v15.6.0" {
		t.Fatalf("FindAssetInfoByTag().Tag = %q, want %q", info.Tag, "v15.6.0")
	}
}
