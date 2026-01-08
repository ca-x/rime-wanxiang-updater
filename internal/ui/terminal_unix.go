//go:build !windows
// +build !windows

package ui

import (
	"golang.org/x/term"
	"os"
)

// makeRaw 将终端设置为 raw 模式（用于读取终端响应）
func makeRaw() (*term.State, error) {
	return term.MakeRaw(int(os.Stdin.Fd()))
}

// restore 恢复终端到之前的状态
func restore(oldState *term.State) error {
	return term.Restore(int(os.Stdin.Fd()), oldState)
}
