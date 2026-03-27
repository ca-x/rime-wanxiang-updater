//go:build linux

package config

import (
	"os"
	"path/filepath"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

func TestDetectInstalledEngines(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalPath := os.Getenv("PATH")
	t.Cleanup(func() {
		_ = os.Setenv("HOME", originalHome)
		_ = os.Setenv("PATH", originalPath)
	})

	tempHome := t.TempDir()
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("Setenv HOME failed: %v", err)
	}
	if err := os.Setenv("PATH", ""); err != nil {
		t.Fatalf("Setenv PATH failed: %v", err)
	}

	engines := DetectInstalledEngines()

	if len(engines) != 0 {
		t.Fatalf("DetectInstalledEngines() = %v, want no engines when nothing is installed", engines)
	}
}

func TestGetRimeUserDirNoInstalledEngine(t *testing.T) {
	result := getRimeUserDir(&types.Config{})

	if result != "" {
		t.Fatalf("getRimeUserDir() = %q, want empty path when no engine is configured or installed", result)
	}
}

func TestGetRimeUserDir(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		config   *types.Config
		wantPath string
	}{
		{
			name: "Primary engine is fcitx5",
			config: &types.Config{
				PrimaryEngine:    "fcitx5",
				InstalledEngines: []string{"fcitx5"},
			},
			wantPath: ".local/share/fcitx5/rime",
		},
		{
			name: "Primary engine is ibus",
			config: &types.Config{
				PrimaryEngine:    "ibus",
				InstalledEngines: []string{"ibus"},
			},
			wantPath: ".config/ibus/rime",
		},
		{
			name: "Primary engine is fcitx",
			config: &types.Config{
				PrimaryEngine:    "fcitx",
				InstalledEngines: []string{"fcitx"},
			},
			wantPath: ".config/fcitx/rime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRimeUserDir(tt.config)
			expectedPath := filepath.Join(homeDir, tt.wantPath)
			if result != expectedPath {
				t.Errorf("getRimeUserDir() = %v, want %v", result, expectedPath)
			}
		})
	}
}

func TestGetEngineDataDir(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name       string
		engineName string
		shouldFind bool
	}{
		{
			name:       "fcitx5 engine",
			engineName: "fcitx5",
			shouldFind: true,
		},
		{
			name:       "ibus engine",
			engineName: "ibus",
			shouldFind: true,
		},
		{
			name:       "fcitx engine",
			engineName: "fcitx",
			shouldFind: true,
		},
		{
			name:       "Unknown engine",
			engineName: "unknown",
			shouldFind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEngineDataDir(tt.engineName)
			if tt.shouldFind {
				if result == "" {
					t.Errorf("GetEngineDataDir() should return path for %s", tt.engineName)
				}
				if !filepath.IsAbs(result) {
					t.Errorf("GetEngineDataDir() should return absolute path, got %v", result)
				}
				if !filepath.HasPrefix(result, homeDir) {
					t.Errorf("GetEngineDataDir() should return path under home dir, got %v", result)
				}
			} else {
				if result != "" {
					t.Errorf("GetEngineDataDir() should return empty for unknown engine, got %v", result)
				}
			}
		})
	}
}

func TestGetAllEngineDataDirs(t *testing.T) {
	installedEngines := []string{"fcitx5", "ibus", "fcitx"}
	result := GetAllEngineDataDirs(installedEngines)

	if len(result) != len(installedEngines) {
		t.Errorf("GetAllEngineDataDirs() returned %d dirs, want %d", len(result), len(installedEngines))
	}

	for _, engine := range installedEngines {
		if _, ok := result[engine]; !ok {
			t.Errorf("GetAllEngineDataDirs() missing engine %s", engine)
		}
	}
}
