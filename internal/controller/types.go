package controller

import (
	"sync"

	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"
)

// Controller manages business logic and state
type Controller struct {
	// Configuration
	cfg *config.Manager

	// Theme management
	themeMgr *theme.Manager

	// Rime detection
	rimeStatus detector.InstallationStatus

	// Update state
	updating         bool
	currentOperation string

	// Wizard state
	wizardState WizardState

	// Communication channels
	commandChan <-chan Command
	eventChan   chan<- Event

	// Shutdown
	done chan struct{}

	// Internal state protection
	mu sync.RWMutex
}

// WizardState tracks wizard progress in controller
type WizardState struct {
	SchemeType string
	Variant    string
	UseMirror  bool
	Completed  bool
}
