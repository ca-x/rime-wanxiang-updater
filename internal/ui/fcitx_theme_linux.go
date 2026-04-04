//go:build linux

package ui

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os/exec"
	"slices"

	projectassets "rime-wanxiang-updater/assets"
	"rime-wanxiang-updater/internal/deployer"
	"rime-wanxiang-updater/internal/types"

	"github.com/godbus/dbus/v5"
)

var deployFcitxTheme = func(cfg *types.Config) error {
	return reloadFcitxTheme(cfg, reloadFcitxClassicUIAddonConfig, reloadFcitxRemoteConfig, func(cfg *types.Config) error {
		return deployer.GetDeployer(cfg).Deploy()
	})
}

func reloadFcitxTheme(
	cfg *types.Config,
	addonReloader func() error,
	remoteReloader func() error,
	deployFallback func(*types.Config) error,
) error {
	if addonReloader == nil {
		return fmt.Errorf("addon reloader is nil")
	}
	if remoteReloader == nil {
		return fmt.Errorf("remote reloader is nil")
	}
	if deployFallback == nil {
		return fmt.Errorf("deploy fallback is nil")
	}

	if err := addonReloader(); err == nil {
		return nil
	}
	if err := remoteReloader(); err == nil {
		return nil
	}
	return deployFallback(cfg)
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

	themeRoot, err := fcitxThemeRootPath()
	if err != nil {
		return err
	}
	if err := installFcitxTheme(themeFS, themeName, themeRoot); err != nil {
		return err
	}

	return applyFcitxThemeConfig(FcitxThemeConfig{Theme: themeName})
}

func applyFcitxThemeDefault(themeName string) error {
	return applyFcitxThemeConfig(FcitxThemeConfig{Theme: themeName})
}

func applyFcitxThemeConfig(cfg FcitxThemeConfig) error {
	configPath, err := fcitxClassicUIConfigPath()
	if err != nil {
		return err
	}

	return setFcitxThemeWithFallback(cfg, configPath, setFcitxThemeViaDBus)
}

func currentFcitxThemeConfig() (FcitxThemeConfig, error) {
	configPath, err := fcitxClassicUIConfigPath()
	if err != nil {
		return FcitxThemeConfig{}, err
	}

	return readFcitxThemeConfigWithFallback(configPath, getFcitxThemeConfigViaDBus)
}

func setFcitxThemeViaDBus(cfg FcitxThemeConfig) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("connect dbus: %w", err)
	}
	defer conn.Close()

	payload, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal fcitx theme config: %w", err)
	}

	obj := conn.Object("org.fcitx.Fcitx5", "/controller")
	call := obj.Call(
		"org.fcitx.Fcitx.Controller1.SetConfig",
		0,
		"fcitx://config/addon/classicui/classicui",
		string(payload),
	)
	if call.Err != nil {
		return fmt.Errorf("set fcitx theme via dbus: %w", call.Err)
	}

	return nil
}

func getFcitxThemeConfigViaDBus() (FcitxThemeConfig, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return FcitxThemeConfig{}, fmt.Errorf("connect dbus: %w", err)
	}
	defer conn.Close()

	obj := conn.Object("org.fcitx.Fcitx5", "/controller")
	call := obj.Call(
		"org.fcitx.Fcitx.Controller1.GetConfig",
		0,
		"fcitx://config/addon/classicui/classicui",
	)
	if call.Err != nil {
		return FcitxThemeConfig{}, fmt.Errorf("get fcitx theme via dbus: %w", call.Err)
	}

	var result string
	if err := call.Store(&result); err != nil {
		return FcitxThemeConfig{}, fmt.Errorf("store fcitx theme config: %w", err)
	}

	var cfg FcitxThemeConfig
	if err := json.Unmarshal([]byte(result), &cfg); err != nil {
		return FcitxThemeConfig{}, fmt.Errorf("unmarshal fcitx theme config: %w", err)
	}

	return cfg, nil
}

func reloadFcitxClassicUIAddonConfig() error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("connect dbus: %w", err)
	}
	defer conn.Close()

	obj := conn.Object("org.fcitx.Fcitx5", "/controller")
	call := obj.Call(
		"org.fcitx.Fcitx.Controller1.ReloadAddonConfig",
		0,
		"classicui",
	)
	if call.Err != nil {
		return fmt.Errorf("reload classicui addon config via dbus: %w", call.Err)
	}

	return nil
}

func reloadFcitxRemoteConfig() error {
	restart := exec.Command("fcitx5-remote", "-r")
	if err := restart.Run(); err != nil {
		return fmt.Errorf("reload fcitx config via fcitx5-remote: %w", err)
	}
	return nil
}
