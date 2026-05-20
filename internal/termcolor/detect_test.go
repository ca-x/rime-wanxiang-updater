package termcolor

import (
	"os"
	"testing"
)

func TestDetectBackgroundDefaultsDarkWithoutEnv(t *testing.T) {
	clearBackgroundEnv(t)

	if got := DetectBackground(); got != BackgroundDark {
		t.Fatalf("DetectBackground() = %v, want %v", got, BackgroundDark)
	}
}

func TestDetectBackgroundUsesExplicitEnv(t *testing.T) {
	clearBackgroundEnv(t)
	t.Setenv("TERM_BACKGROUND", "light")

	if got := DetectBackground(); got != BackgroundLight {
		t.Fatalf("DetectBackground() = %v, want %v", got, BackgroundLight)
	}
}

func TestDetectBackgroundPrefersColorFgBg(t *testing.T) {
	clearBackgroundEnv(t)
	t.Setenv("COLORFGBG", "15;0")
	t.Setenv("TERM_BACKGROUND", "light")

	if got := DetectBackground(); got != BackgroundDark {
		t.Fatalf("DetectBackground() = %v, want %v", got, BackgroundDark)
	}
}

func TestInitLipglossPersistsDetectedBackground(t *testing.T) {
	clearBackgroundEnv(t)
	t.Setenv("TERMINAL_BACKGROUND", "light")

	if got := InitLipgloss(); got != BackgroundLight {
		t.Fatalf("InitLipgloss() = %v, want %v", got, BackgroundLight)
	}
	if got := os.Getenv("TERM_BACKGROUND"); got != "light" {
		t.Fatalf("TERM_BACKGROUND = %q, want %q", got, "light")
	}
}

func clearBackgroundEnv(t *testing.T) {
	t.Helper()
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERMINAL_BACKGROUND", "")
	t.Setenv("TERM_BACKGROUND", "")
	t.Setenv("COLORTHEME", "")
}
