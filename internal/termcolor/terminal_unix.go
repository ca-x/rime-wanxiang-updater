//go:build !windows

package termcolor

import (
	"os"

	"golang.org/x/term"
)

func makeRaw() (*term.State, error) {
	fd := int(os.Stdin.Fd())
	return term.MakeRaw(fd)
}

func restore(state *term.State) {
	if state != nil {
		term.Restore(int(os.Stdin.Fd()), state)
	}
}
