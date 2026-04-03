//go:build linux

package ui

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	projectassets "rime-wanxiang-updater/assets"
	"rime-wanxiang-updater/internal/deployer"
	"rime-wanxiang-updater/internal/types"

	"github.com/godbus/dbus/v5"
)

var deployFcitxTheme = func(cfg *types.Config) error {
	restart := exec.Command("fcitx5-remote", "-r")
	if err := restart.Run(); err == nil {
		return nil
	}
	return deployer.GetDeployer(cfg).Deploy()
}

func availableFcitxThemes() ([]string, error) {
	themeFS := projectassets.Fcitx5Themes()
	if themeFS == nil {
		return nil, fmt.Errorf("fcitx5 themes are not embedded")
	}

	entries, err := fs.ReadDir(themeFS, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded fcitx5 themes: %w", err)
	}

	var themes []string
	for _, entry := range entries {
		if entry.IsDir() {
			themes = append(themes, entry.Name())
		}
	}

	slices.Sort(themes)
	if len(themes) == 0 {
		return nil, fmt.Errorf("no embedded fcitx5 themes found")
	}

	return themes, nil
}

func installAndSetFcitxTheme(themeName string) error {
	themeFS := projectassets.Fcitx5Themes()
	if themeFS == nil {
		return fmt.Errorf("fcitx5 themes are not embedded")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	themeRoot := filepath.Join(homeDir, ".local", "share", "fcitx5", "themes")
	if err := installFcitxTheme(themeFS, themeName, themeRoot); err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".config", "fcitx5", "conf", "classicui.conf")
	return setFcitxThemeWithFallback(themeName, configPath, setFcitxThemeViaDBus)
}

func setFcitxThemeViaDBus(themeName string) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("connect dbus: %w", err)
	}
	defer conn.Close()

	obj := conn.Object("org.fcitx.Fcitx5", "/controller")
	configValue := fmt.Sprintf(`{"Theme":"%s"}`, themeName)
	call := obj.Call(
		"org.fcitx.Fcitx.Controller1.SetConfig",
		0,
		"fcitx://config/addon/classicui/classicui",
		configValue,
	)
	if call.Err != nil {
		return fmt.Errorf("set fcitx theme via dbus: %w", call.Err)
	}

	return nil
}
