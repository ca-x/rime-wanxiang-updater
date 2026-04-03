package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestFcitxThemeSupportedForPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		engines  []string
		want     bool
	}{
		{
			name:     "linux with fcitx5",
			platform: "linux",
			engines:  []string{"fcitx5", "ibus"},
			want:     true,
		},
		{
			name:     "linux without fcitx5",
			platform: "linux",
			engines:  []string{"ibus"},
			want:     false,
		},
		{
			name:     "darwin with fcitx5 is not supported",
			platform: "darwin",
			engines:  []string{"fcitx5"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fcitxThemeSupportedForPlatform(tt.platform, tt.engines); got != tt.want {
				t.Fatalf("fcitxThemeSupportedForPlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFcitxClassicUIConfigPreservesOtherLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "classicui.conf")

	initial := "[Groups/0]\nName=ClassicUI\nTheme=old-theme\nVertical Candidate List=False\n"
	if err := os.WriteFile(path, []byte(initial), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := writeFcitxClassicUIConfig(path, "new-theme"); err != nil {
		t.Fatalf("writeFcitxClassicUIConfig() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	content := string(data)
	for _, want := range []string{
		"Name=ClassicUI",
		"Theme=new-theme",
		"Vertical Candidate List=False",
	} {
		if !containsLine(content, want) {
			t.Fatalf("classicui.conf missing line %q in %q", want, content)
		}
	}
}

func TestInstallFcitxThemeCopiesThemeDirectory(t *testing.T) {
	source := fstest.MapFS{
		"demo/theme.conf":  &fstest.MapFile{Data: []byte("Theme=demo\n")},
		"demo/panel.svg":   &fstest.MapFile{Data: []byte("<svg></svg>")},
		"demo/arrow.png":   &fstest.MapFile{Data: []byte("png")},
		"other/theme.conf": &fstest.MapFile{Data: []byte("Theme=other\n")},
	}

	destRoot := t.TempDir()
	if err := installFcitxTheme(source, "demo", destRoot); err != nil {
		t.Fatalf("installFcitxTheme() error = %v", err)
	}

	for _, rel := range []string{
		"demo/theme.conf",
		"demo/panel.svg",
		"demo/arrow.png",
	} {
		path := filepath.Join(destRoot, filepath.FromSlash(rel))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat(%q) error = %v", path, err)
		}
		if info.IsDir() {
			t.Fatalf("%q should be a file", path)
		}
	}

	if _, err := os.Stat(filepath.Join(destRoot, "other")); !os.IsNotExist(err) {
		t.Fatalf("unexpected unrelated theme copied")
	}
}

func TestSetFcitxThemePrefersDBusAndFallsBackToConfig(t *testing.T) {
	t.Run("dbus success skips file fallback", func(t *testing.T) {
		configPath := filepath.Join(t.TempDir(), "classicui.conf")
		called := false

		err := setFcitxThemeWithFallback("demo", configPath, func(themeName string) error {
			called = true
			if themeName != "demo" {
				t.Fatalf("themeName = %q, want %q", themeName, "demo")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("setFcitxThemeWithFallback() error = %v", err)
		}
		if !called {
			t.Fatalf("dbus setter was not called")
		}
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			t.Fatalf("config fallback should not run when dbus succeeds")
		}
	})

	t.Run("dbus failure falls back to config", func(t *testing.T) {
		configPath := filepath.Join(t.TempDir(), "classicui.conf")

		err := setFcitxThemeWithFallback("fallback-theme", configPath, func(string) error {
			return os.ErrPermission
		})
		if err != nil {
			t.Fatalf("setFcitxThemeWithFallback() error = %v", err)
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if !containsLine(string(data), "Theme=fallback-theme") {
			t.Fatalf("config fallback did not write theme: %q", string(data))
		}
	})
}

func containsLine(content, line string) bool {
	for _, item := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if item == line {
			return true
		}
	}
	return false
}
