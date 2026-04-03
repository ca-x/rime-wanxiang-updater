//go:build linux

package assets

import (
	"embed"
	"io/fs"
)

//go:embed fcitx5-themes/*
var fcitx5Themes embed.FS

func Fcitx5Themes() fs.FS {
	sub, err := fs.Sub(fcitx5Themes, "fcitx5-themes")
	if err != nil {
		return nil
	}
	return sub
}
