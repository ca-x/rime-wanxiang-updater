//go:build !linux

package assets

import "io/fs"

func Fcitx5Themes() fs.FS {
	return nil
}
