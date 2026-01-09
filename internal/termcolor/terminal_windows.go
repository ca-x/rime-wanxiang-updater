//go:build windows

package termcolor

type windowsState struct{}

func makeRaw() (*windowsState, error) {
	return &windowsState{}, nil
}

func restore(state *windowsState) {
}
