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

const (
	fcitxThemeSelectionLight = "light"
	fcitxThemeSelectionDark  = "dark"
)

type FcitxThemeConfig struct {
	Theme                string `json:"Theme,omitzero"`
	DarkTheme            string `json:"DarkTheme,omitzero"`
	UseDarkTheme         *bool  `json:"UseDarkTheme,omitzero"`
	FollowSystemDarkMode *bool  `json:"FollowSystemDarkMode,omitzero"`
}

func fcitxThemeSupportedForPlatform(platform string, installedEngines []string) bool {
	return platform == "linux" && slices.Contains(installedEngines, "fcitx5")
}

func boolPtr(value bool) *bool {
	return &value
}

func boolToFcitxString(value bool) string {
	if value {
		return "True"
	}
	return "False"
}

func parseFcitxBool(raw string) (*bool, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true", "1", "yes", "on":
		return boolPtr(true), true
	case "false", "0", "no", "off":
		return boolPtr(false), true
	default:
		return nil, false
	}
}

func (cfg FcitxThemeConfig) values() map[string]string {
	values := make(map[string]string)
	if cfg.Theme != "" {
		values["Theme"] = cfg.Theme
	}
	if cfg.DarkTheme != "" {
		values["DarkTheme"] = cfg.DarkTheme
	}
	if cfg.UseDarkTheme != nil {
		values["UseDarkTheme"] = boolToFcitxString(*cfg.UseDarkTheme)
	}
	if cfg.FollowSystemDarkMode != nil {
		values["FollowSystemDarkMode"] = boolToFcitxString(*cfg.FollowSystemDarkMode)
	}
	return values
}

func writeFcitxClassicUIConfig(configPath string, cfg FcitxThemeConfig) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("create classicui config dir: %w", err)
	}

	updates := cfg.values()
	var lines []string
	found := make(map[string]bool)

	file, err := os.Open(configPath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				if value, ok := updates[key]; ok {
					lines = append(lines, key+"="+value)
					found[key] = true
					continue
				}
			}
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read classicui config: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("open classicui config: %w", err)
	}

	for key, value := range updates {
		if found[key] {
			continue
		}
		lines = append(lines, key+"="+value)
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

func setFcitxThemeWithFallback(cfg FcitxThemeConfig, configPath string, dbusSetter func(FcitxThemeConfig) error) error {
	if dbusSetter == nil {
		return fmt.Errorf("dbus setter is nil")
	}
	if err := dbusSetter(cfg); err == nil {
		return nil
	}
	return writeFcitxClassicUIConfig(configPath, cfg)
}

func readFcitxClassicUIConfig(configPath string) (FcitxThemeConfig, error) {
	var cfg FcitxThemeConfig

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("open classicui config: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "Theme":
			cfg.Theme = value
		case "DarkTheme":
			cfg.DarkTheme = value
		case "UseDarkTheme":
			if parsed, ok := parseFcitxBool(value); ok {
				cfg.UseDarkTheme = parsed
			}
		case "FollowSystemDarkMode":
			if parsed, ok := parseFcitxBool(value); ok {
				cfg.FollowSystemDarkMode = parsed
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return cfg, fmt.Errorf("read classicui config: %w", err)
	}

	return cfg, nil
}

func readFcitxThemeConfigWithFallback(configPath string, dbusGetter func() (FcitxThemeConfig, error)) (FcitxThemeConfig, error) {
	if dbusGetter == nil {
		return readFcitxClassicUIConfig(configPath)
	}

	cfg, err := dbusGetter()
	if err == nil {
		return cfg, nil
	}

	return readFcitxClassicUIConfig(configPath)
}

func readFcitxClassicUITheme(configPath string) (string, error) {
	cfg, err := readFcitxClassicUIConfig(configPath)
	if err != nil {
		return "", err
	}
	return cfg.Theme, nil
}

func installedFcitxThemeSelections(destRoot string, builtinThemeNames []string) (map[string]bool, error) {
	selections := make(map[string]bool)
	for _, themeName := range builtinThemeNames {
		info, err := os.Stat(filepath.Join(destRoot, themeName))
		if err == nil && info.IsDir() {
			selections[themeName] = true
			continue
		}
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat installed theme %q: %w", themeName, err)
		}
	}

	return selections, nil
}

func syncInstalledFcitxThemes(themeFS fs.FS, destRoot string, builtinThemeNames []string, selections map[string]bool) error {
	for _, themeName := range builtinThemeNames {
		targetDir := filepath.Join(destRoot, themeName)
		if selections[themeName] {
			if err := installFcitxTheme(themeFS, themeName, destRoot); err != nil {
				return err
			}
			continue
		}
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("remove theme dir %q: %w", targetDir, err)
		}
	}

	return nil
}

func fcitxThemeRootPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	return filepath.Join(homeDir, ".local", "share", "fcitx5", "themes"), nil
}

func fcitxClassicUIConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	return filepath.Join(homeDir, ".config", "fcitx5", "conf", "classicui.conf"), nil
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
