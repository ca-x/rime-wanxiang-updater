//go:build !linux

package ui

import (
	"fmt"

	"rime-wanxiang-updater/internal/types"
)

var deployFcitxTheme = func(cfg *types.Config) error {
	return fmt.Errorf("fcitx5 themes are only supported on linux")
}

func availableFcitxThemes() ([]string, error) {
	return nil, fmt.Errorf("fcitx5 themes are only supported on linux")
}

func installAndSetFcitxTheme(themeName string) error {
	return fmt.Errorf("fcitx5 themes are only supported on linux")
}

func applyFcitxThemeDefault(themeName string) error {
	return fmt.Errorf("fcitx5 themes are only supported on linux")
}

func applyFcitxThemeConfig(cfg FcitxThemeConfig) error {
	return fmt.Errorf("fcitx5 themes are only supported on linux")
}

func currentFcitxThemeConfig() (FcitxThemeConfig, error) {
	return FcitxThemeConfig{}, fmt.Errorf("fcitx5 themes are only supported on linux")
}
