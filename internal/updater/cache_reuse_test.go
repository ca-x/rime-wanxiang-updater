package updater

import (
	"path/filepath"
	"testing"
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

func TestCanReuseCachedAsset(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFile := filepath.Join(tmpDir, "asset.zip")

	cfg := &config.Manager{
		ConfigPath: filepath.Join(tmpDir, "config.json"),
		Config: &types.Config{
			SchemeType: "base",
			UseMirror:  true,
		},
	}

	updater := NewBaseUpdater(cfg)

	oldTime := time.Date(2026, 1, 17, 0, 40, 48, 0, time.UTC)
	newTime := time.Date(2026, 4, 2, 13, 31, 10, 0, time.UTC)

	record := &types.UpdateRecord{
		Name:       "base-dicts.zip",
		UpdateTime: oldTime,
		Tag:        "v1.0.0",
		SHA256:     "same-hash",
	}

	info := &types.UpdateInfo{
		Name:       "base-dicts.zip",
		UpdateTime: newTime,
		Tag:        "v1.0.0",
		SHA256:     "same-hash",
	}

	t.Run("same tag with newer remote time must not reuse cache", func(t *testing.T) {
		info.SHA256 = ""
		info.ID = ""
		record.CnbID = ""

		if updater.canReuseCachedAsset(record, info, cacheFile, func(string, string) bool { return true }) {
			t.Fatal("canReuseCachedAsset() = true, want false when remote time is newer and asset identity is unknown")
		}
	})

	t.Run("same tag with newer remote time may reuse cache when the asset identity matches", func(t *testing.T) {
		info.UpdateTime = newTime
		info.SHA256 = "same-hash"
		info.ID = "asset-123"
		record.CnbID = "asset-123"

		if !updater.canReuseCachedAsset(record, info, cacheFile, func(string, string) bool { return true }) {
			t.Fatal("canReuseCachedAsset() = false, want true when CNB reports the same asset with a newer timestamp")
		}
	})

	t.Run("same tag and same remote time may reuse cache", func(t *testing.T) {
		info.UpdateTime = oldTime
		info.ID = ""
		record.CnbID = ""

		if !updater.canReuseCachedAsset(record, info, cacheFile, func(string, string) bool { return true }) {
			t.Fatal("canReuseCachedAsset() = false, want true when remote time is not newer")
		}
	})

	t.Run("missing record hash must not reuse cache", func(t *testing.T) {
		info.UpdateTime = oldTime
		info.SHA256 = ""
		record.SHA256 = ""

		if updater.canReuseCachedAsset(record, info, cacheFile, func(string, string) bool { return true }) {
			t.Fatal("canReuseCachedAsset() = true, want false when neither local nor remote hash is available")
		}
	})
}
