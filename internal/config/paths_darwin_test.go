//go:build darwin

package config

import (
	"os"
	"path/filepath"
	"testing"

	"rime-wanxiang-updater/internal/types"
)

func TestDetectInstalledEngines(t *testing.T) {
	// 这个测试会检测实际系统中的引擎
	engines := DetectInstalledEngines()

	if len(engines) == 0 {
		t.Error("Expected at least one engine (default)")
	}

	// 至少应该返回一个引擎名称
	for _, engine := range engines {
		if engine == "" {
			t.Error("Engine name should not be empty")
		}
	}
}

func TestGetRimeUserDir(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		config   *types.Config
		expected string
	}{
		{
			name: "Primary engine is 鼠须管",
			config: &types.Config{
				PrimaryEngine:    "鼠须管",
				InstalledEngines: []string{"鼠须管"},
			},
			expected: filepath.Join(homeDir, "Library", "Rime"),
		},
		{
			name: "Primary engine is 小企鹅",
			config: &types.Config{
				PrimaryEngine:    "小企鹅",
				InstalledEngines: []string{"小企鹅"},
			},
			expected: filepath.Join(homeDir, ".local", "share", "fcitx5", "rime"),
		},
		{
			name: "No primary engine but has installed engines",
			config: &types.Config{
				PrimaryEngine:    "",
				InstalledEngines: []string{"鼠须管", "小企鹅"},
			},
			expected: filepath.Join(homeDir, "Library", "Rime"),
		},
		{
			name: "Legacy config with Engine field",
			config: &types.Config{
				Engine:           "小企鹅",
				PrimaryEngine:    "",
				InstalledEngines: []string{},
			},
			expected: filepath.Join(homeDir, ".local", "share", "fcitx5", "rime"),
		},
		{
			name: "Empty config defaults to 鼠须管",
			config: &types.Config{
				PrimaryEngine:    "",
				InstalledEngines: []string{},
			},
			expected: filepath.Join(homeDir, "Library", "Rime"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRimeUserDir(tt.config)
			if result != tt.expected {
				t.Errorf("getRimeUserDir() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEngineDataDir(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name       string
		engineName string
		expected   string
	}{
		{
			name:       "鼠须管 engine",
			engineName: "鼠须管",
			expected:   filepath.Join(homeDir, "Library", "Rime"),
		},
		{
			name:       "小企鹅 engine",
			engineName: "小企鹅",
			expected:   filepath.Join(homeDir, ".local", "share", "fcitx5", "rime"),
		},
		{
			name:       "Unknown engine",
			engineName: "unknown",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEngineDataDir(tt.engineName)
			if result != tt.expected {
				t.Errorf("GetEngineDataDir() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetAllEngineDataDirs(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	installedEngines := []string{"鼠须管", "小企鹅"}
	result := GetAllEngineDataDirs(installedEngines)

	expected := map[string]string{
		"鼠须管": filepath.Join(homeDir, "Library", "Rime"),
		"小企鹅": filepath.Join(homeDir, ".local", "share", "fcitx5", "rime"),
	}

	if len(result) != len(expected) {
		t.Errorf("GetAllEngineDataDirs() returned %d dirs, want %d", len(result), len(expected))
	}

	for engine, expectedPath := range expected {
		if result[engine] != expectedPath {
			t.Errorf("GetAllEngineDataDirs()[%s] = %v, want %v", engine, result[engine], expectedPath)
		}
	}
}

func TestDetectInstallationPaths(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name           string
		engine         string
		expectedRime   string
		expectedAppKey string
	}{
		{
			name:           "鼠须管 installation",
			engine:         "鼠须管",
			expectedRime:   filepath.Join(homeDir, "Library", "Rime"),
			expectedAppKey: "app_path",
		},
		{
			name:           "小企鹅 installation",
			engine:         "小企鹅",
			expectedRime:   filepath.Join(homeDir, ".local", "share", "fcitx5", "rime"),
			expectedAppKey: "app_path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectInstallationPaths(tt.engine)

			if result["rime_user_dir"] != tt.expectedRime {
				t.Errorf("DetectInstallationPaths()[rime_user_dir] = %v, want %v",
					result["rime_user_dir"], tt.expectedRime)
			}

			if _, ok := result[tt.expectedAppKey]; !ok {
				t.Errorf("DetectInstallationPaths() missing key %s", tt.expectedAppKey)
			}
		})
	}
}
