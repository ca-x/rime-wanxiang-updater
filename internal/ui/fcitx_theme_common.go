package ui

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func fcitxThemeSupportedForPlatform(platform string, installedEngines []string) bool {
	return platform == "linux" && slices.Contains(installedEngines, "fcitx5")
}

func writeFcitxClassicUIConfig(configPath, themeName string) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("create classicui config dir: %w", err)
	}

	var lines []string
	found := false

	file, err := os.Open(configPath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Theme=") {
				lines = append(lines, "Theme="+themeName)
				found = true
				continue
			}
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read classicui config: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("open classicui config: %w", err)
	}

	if !found {
		lines = append(lines, "Theme="+themeName)
	}

	content := strings.Join(lines, "\n")
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write classicui config: %w", err)
	}

	return nil
}

func setFcitxThemeWithFallback(themeName, configPath string, dbusSetter func(string) error) error {
	if dbusSetter == nil {
		return fmt.Errorf("dbus setter is nil")
	}
	if err := dbusSetter(themeName); err == nil {
		return nil
	}
	return writeFcitxClassicUIConfig(configPath, themeName)
}

func installFcitxTheme(themeFS fs.FS, themeName, destRoot string) error {
	if themeFS == nil {
		return fmt.Errorf("theme fs is nil")
	}

	targetDir := filepath.Join(destRoot, themeName)
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("remove old theme dir: %w", err)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create theme dir: %w", err)
	}

	return fs.WalkDir(themeFS, themeName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(themeName, path)
		if err != nil {
			return fmt.Errorf("resolve theme path: %w", err)
		}

		destPath := filepath.Join(targetDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		srcFile, err := themeFS.Open(path)
		if err != nil {
			return fmt.Errorf("open embedded theme file: %w", err)
		}
		defer srcFile.Close()

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("create theme parent dir: %w", err)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("read embedded theme file info: %w", err)
		}

		dstFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			return fmt.Errorf("create target theme file: %w", err)
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return fmt.Errorf("copy theme file: %w", err)
		}

		return nil
	})
}
