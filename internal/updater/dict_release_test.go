package updater

import (
	"testing"
	"time"

	"rime-wanxiang-updater/internal/types"
)

func TestFindDictReleasePrefersSchemeTag(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/v1.0.0-dicts.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 8, 0, 0, 0, time.UTC),
				},
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/v1.0.0-scheme.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 8, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			TagName: "v15.6.0",
			Assets: []types.GitHubAsset{
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/v15.6.0-scheme.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC),
				},
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/v15.6.0-dicts.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	info, ok := findDictRelease(releases, "rime-wanxiang-base.zip", "base-dicts.zip")
	if !ok {
		t.Fatal("findDictRelease() = no match, want match")
	}

	if info.Tag != "v15.6.0" {
		t.Fatalf("findDictRelease().Tag = %q, want %q", info.Tag, "v15.6.0")
	}
}

func TestFindDictReleaseFallsBackToCNBDefaultTag(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "base-dicts.zip",
					BrowserDownloadURL: "https://example.com/v1.0.0-dicts.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 8, 0, 0, 0, time.UTC),
				},
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/v1.0.0-scheme.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 8, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			TagName: "v15.6.0",
			Assets: []types.GitHubAsset{
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/v15.6.0-scheme.zip",
					UpdatedAt:          time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	info, ok := findDictRelease(releases, "rime-wanxiang-base.zip", "base-dicts.zip")
	if !ok {
		t.Fatal("findDictRelease() = no match, want match")
	}

	if info.Tag != types.CNB_DICT_TAG {
		t.Fatalf("findDictRelease().Tag = %q, want %q", info.Tag, types.CNB_DICT_TAG)
	}
}
