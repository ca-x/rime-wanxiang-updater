package updater

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/types"
)

func TestModelStatusRequiresActualInstalledFile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("test simulates a Linux engine install path")
	}

	tmpDir := t.TempDir()
	tempBin := t.TempDir()

	originalHome := os.Getenv("HOME")
	originalPath := os.Getenv("PATH")
	t.Cleanup(func() {
		_ = os.Setenv("HOME", originalHome)
		_ = os.Setenv("PATH", originalPath)
	})

	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Setenv HOME failed: %v", err)
	}
	if err := os.Setenv("PATH", tempBin); err != nil {
		t.Fatalf("Setenv PATH failed: %v", err)
	}

	fakeBinary := filepath.Join(tempBin, "fcitx5")
	if err := os.WriteFile(fakeBinary, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatalf("WriteFile(fakeBinary) error = %v", err)
	}

	rimeDir := filepath.Join(tmpDir, ".local", "share", "fcitx5", "rime")
	if err := os.MkdirAll(rimeDir, 0755); err != nil {
		t.Fatalf("MkdirAll(rimeDir) error = %v", err)
	}

	cfg := &config.Manager{
		ConfigPath: filepath.Join(tmpDir, "config.json"),
		RimeDir:    rimeDir,
		CacheDir:   tmpDir,
		ZhDictsDir: types.ZH_DICTS,
		Config: &types.Config{
			InstalledEngines: []string{"test-engine"},
			PrimaryEngine:    "test-engine",
			UseMirror:        true,
		},
	}

	updater := NewModelUpdater(cfg)
	updater.UpdateInfo = &types.UpdateInfo{
		Name:       types.MODEL_FILE,
		Tag:        "2026-01-08",
		UpdateTime: time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC),
	}

	recordPath := filepath.Join(tmpDir, "model_record.json")
	if err := updater.SaveRecord(recordPath, "model_name", types.MODEL_FILE, updater.UpdateInfo); err != nil {
		t.Fatalf("SaveRecord() error = %v", err)
	}

	status, err := updater.GetStatus()
	if err != nil {
		if strings.Contains(err.Error(), "状态码: 429") {
			t.Skipf("CNB rate limited this live test: %v", err)
		}
		t.Fatalf("GetStatus() error = %v", err)
	}

	if status.LocalVersion != "未安装" {
		t.Fatalf("GetStatus().LocalVersion = %q, want %q when model file is missing", status.LocalVersion, "未安装")
	}
	if !status.NeedsUpdate {
		t.Fatal("GetStatus().NeedsUpdate = false, want true when model file is missing")
	}
}
