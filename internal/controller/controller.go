package controller

import (
	"rime-wanxiang-updater/internal/config"
	"rime-wanxiang-updater/internal/detector"
	"rime-wanxiang-updater/internal/theme"
)

// NewController creates a new controller instance
func NewController(
	cfg *config.Manager,
	commandChan <-chan Command,
	eventChan chan<- Event,
) *Controller {
	themeMgr := theme.NewManager()

	// Load theme settings from config
	if cfg.Config.ThemeAdaptive {
		light := cfg.Config.ThemeLight
		dark := cfg.Config.ThemeDark
		if light == "" {
			light = "cyberpunk-light"
		}
		if dark == "" {
			dark = "cyberpunk"
		}
		themeMgr.SetAdaptiveTheme(light, dark)
	} else if cfg.Config.ThemeFixed != "" {
		themeMgr.SetTheme(cfg.Config.ThemeFixed)
	}

	rimeStatus := detector.CheckRimeInstallation()

	return &Controller{
		cfg:         cfg,
		themeMgr:    themeMgr,
		rimeStatus:  rimeStatus,
		commandChan: commandChan,
		eventChan:   eventChan,
		done:        make(chan struct{}),
	}
}

// Run starts the controller's main loop (runs in goroutine)
func (c *Controller) Run() {
	for {
		select {
		case cmd, ok := <-c.commandChan:
			if !ok {
				// Command channel closed, exit
				return
			}
			c.handleCommand(cmd)

		case <-c.done:
			// Graceful shutdown
			return
		}
	}
}

// Stop gracefully stops the controller
func (c *Controller) Stop() {
	close(c.done)
}

// handleCommand dispatches commands to appropriate handlers
func (c *Controller) handleCommand(cmd Command) {
	switch cmd.Type {
	case CmdAutoUpdate:
		c.handleAutoUpdate(cmd)
	case CmdUpdateDict:
		c.handleUpdateDict(cmd)
	case CmdUpdateScheme:
		c.handleUpdateScheme(cmd)
	case CmdUpdateModel:
		c.handleUpdateModel(cmd)
	case CmdConfigChange:
		c.handleConfigChange(cmd)
	case CmdConfigSave:
		c.handleConfigSave(cmd)
	case CmdWizardSetScheme:
		c.handleWizardSetScheme(cmd)
	case CmdWizardSetMirror:
		c.handleWizardSetMirror(cmd)
	case CmdWizardComplete:
		c.handleWizardComplete(cmd)
	case CmdThemeChange:
		c.handleThemeChange(cmd)
	case CmdExcludeAdd:
		c.handleExcludeAdd(cmd)
	case CmdExcludeRemove:
		c.handleExcludeRemove(cmd)
	case CmdExcludeEdit:
		c.handleExcludeEdit(cmd)
	case CmdShutdown:
		c.Stop()
	}
}

// emitEvent sends an event to the UI
func (c *Controller) emitEvent(eventType EventType, payload any) {
	select {
	case c.eventChan <- Event{Type: eventType, Payload: payload}:
	default:
		// Channel full or closed, skip
	}
}

// emitProgress sends a progress update event
func (c *Controller) emitProgress(component, message string, percent float64, source, fileName string, downloaded, totalSize int64, speed float64, isDownload bool) {
	c.emitEvent(EvtProgressUpdate, ProgressUpdatePayload{
		Component:  component,
		Message:    message,
		Percent:    percent,
		Source:     source,
		FileName:   fileName,
		Downloaded: downloaded,
		TotalSize:  totalSize,
		Speed:      speed,
		IsDownload: isDownload,
	})
}

// emitError sends an error event
func (c *Controller) emitError(err error, context string) {
	c.emitEvent(EvtError, ErrorPayload{
		Error:   err,
		Context: context,
	})
}
