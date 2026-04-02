package updater

import (
	"testing"
	"time"

	"rime-wanxiang-updater/internal/types"
)

func TestFindSchemeReleasePrefersVersionedCNBRelease(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/rolling.zip",
					UpdatedAt:          time.Date(2026, 4, 2, 13, 37, 22, 0, time.UTC),
				},
			},
		},
		{
			TagName: "v15.5.0",
			Assets: []types.GitHubAsset{
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/v15.5.0.zip",
					UpdatedAt:          time.Date(2026, 4, 2, 13, 13, 5, 0, time.UTC),
				},
			},
		},
	}

	info, ok := findSchemeRelease(releases, "rime-wanxiang-base.zip")
	if !ok {
		t.Fatal("findSchemeRelease() = no match, want match")
	}

	if info.Tag != "v15.5.0" {
		t.Fatalf("findSchemeRelease().Tag = %q, want %q", info.Tag, "v15.5.0")
	}
}

func TestFindSchemeReleaseFallsBackToRollingPreview(t *testing.T) {
	releases := []types.GitHubRelease{
		{
			TagName: types.CNB_DICT_TAG,
			Assets: []types.GitHubAsset{
				{
					Name:               "rime-wanxiang-base.zip",
					BrowserDownloadURL: "https://example.com/rolling.zip",
					UpdatedAt:          time.Date(2026, 4, 2, 13, 37, 22, 0, time.UTC),
				},
			},
		},
	}

	info, ok := findSchemeRelease(releases, "rime-wanxiang-base.zip")
	if !ok {
		t.Fatal("findSchemeRelease() = no match, want match")
	}

	if info.Tag != types.CNB_DICT_TAG {
		t.Fatalf("findSchemeRelease().Tag = %q, want %q", info.Tag, types.CNB_DICT_TAG)
	}
}
