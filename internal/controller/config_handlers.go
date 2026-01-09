package controller

import (
	"fmt"
	"strconv"
)

// handleConfigChange handles configuration changes
func (c *Controller) handleConfigChange(cmd Command) {
	payload, ok := cmd.Payload.(ConfigChangePayload)
	if !ok {
		c.emitError(fmt.Errorf("invalid config change payload"), "config change")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	switch payload.Key {
	case "use_mirror":
		if val, ok := payload.Value.(bool); ok {
			c.cfg.Config.UseMirror = val
		}
	case "auto_update":
		if val, ok := payload.Value.(bool); ok {
			c.cfg.Config.AutoUpdate = val
		}
	case "auto_update_countdown":
		if val, ok := payload.Value.(int); ok {
			if val < 1 {
				val = 1
			} else if val > 60 {
				val = 60
			}
			c.cfg.Config.AutoUpdateCountdown = val
		} else if val, ok := payload.Value.(string); ok {
			if countdown, err := strconv.Atoi(val); err == nil {
				if countdown < 1 {
					countdown = 1
				} else if countdown > 60 {
					countdown = 60
				}
				c.cfg.Config.AutoUpdateCountdown = countdown
			}
		}
	case "proxy_enabled":
		if val, ok := payload.Value.(bool); ok {
			c.cfg.Config.ProxyEnabled = val
		}
	case "proxy_type":
		if val, ok := payload.Value.(string); ok {
			c.cfg.Config.ProxyType = val
		}
	case "proxy_address":
		if val, ok := payload.Value.(string); ok {
			c.cfg.Config.ProxyAddress = val
		}
	case "fcitx_compat":
		if val, ok := payload.Value.(bool); ok {
			c.cfg.Config.FcitxCompat = val
		}
	case "fcitx_use_link":
		if val, ok := payload.Value.(bool); ok {
			c.cfg.Config.FcitxUseLink = val
		}
	case "pre_update_hook":
		if val, ok := payload.Value.(string); ok {
			c.cfg.Config.PreUpdateHook = val
		}
	case "post_update_hook":
		if val, ok := payload.Value.(string); ok {
			c.cfg.Config.PostUpdateHook = val
		}
	case "theme_adaptive":
		if val, ok := payload.Value.(bool); ok {
			c.cfg.Config.ThemeAdaptive = val
			// Update theme manager
			if c.cfg.Config.ThemeAdaptive {
				light := c.cfg.Config.ThemeLight
				dark := c.cfg.Config.ThemeDark
				if light == "" {
					light = "cyberpunk-light"
				}
				if dark == "" {
					dark = "cyberpunk"
				}
				c.themeMgr.SetAdaptiveTheme(light, dark)
			} else if c.cfg.Config.ThemeFixed != "" {
				c.themeMgr.SetTheme(c.cfg.Config.ThemeFixed)
			}
		}
	default:
		c.emitError(fmt.Errorf("unknown config key: %s", payload.Key), "config change")
		return
	}

	c.emitEvent(EvtConfigUpdated, ConfigUpdatedPayload{
		Key:   payload.Key,
		Value: payload.Value,
	})
}

// handleConfigSave handles saving configuration
func (c *Controller) handleConfigSave(cmd Command) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.cfg.SaveConfig(); err != nil {
		c.emitError(err, "save config")
		c.emitEvent(EvtConfigError, ErrorPayload{
			Error:   err,
			Context: "Failed to save configuration",
		})
		return
	}

	c.emitEvent(EvtConfigUpdated, ConfigUpdatedPayload{
		Key:   "_saved",
		Value: true,
	})
}

// handleWizardSetScheme handles wizard scheme selection
func (c *Controller) handleWizardSetScheme(cmd Command) {
	payload, ok := cmd.Payload.(WizardSchemePayload)
	if !ok {
		c.emitError(fmt.Errorf("invalid wizard scheme payload"), "wizard scheme")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.wizardState.SchemeType = payload.SchemeType
	c.wizardState.Variant = payload.Variant

	c.cfg.Config.SchemeType = payload.SchemeType

	c.emitEvent(EvtStateUpdate, c.wizardState)
}

// handleWizardSetMirror handles wizard mirror selection
func (c *Controller) handleWizardSetMirror(cmd Command) {
	payload, ok := cmd.Payload.(bool)
	if !ok {
		c.emitError(fmt.Errorf("invalid wizard mirror payload"), "wizard mirror")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.wizardState.UseMirror = payload
	c.cfg.Config.UseMirror = payload

	c.emitEvent(EvtStateUpdate, c.wizardState)
}

// handleWizardComplete handles wizard completion
func (c *Controller) handleWizardComplete(cmd Command) {
	c.mu.Lock()
	defer c.mu.Unlock()

	schemeChoice := c.wizardState.SchemeType
	if c.wizardState.Variant != "" {
		schemeChoice = c.wizardState.Variant
	}

	schemeFile, dictFile, err := c.cfg.GetActualFilenames(schemeChoice)
	if err != nil {
		c.emitError(err, "wizard complete")
		c.emitEvent(EvtWizardError, ErrorPayload{
			Error:   err,
			Context: "Failed to get filenames",
		})
		return
	}

	c.cfg.Config.SchemeFile = schemeFile
	c.cfg.Config.DictFile = dictFile

	if err := c.cfg.SaveConfig(); err != nil {
		c.emitError(err, "wizard complete")
		c.emitEvent(EvtWizardError, ErrorPayload{
			Error:   err,
			Context: "Failed to save configuration",
		})
		return
	}

	c.wizardState.Completed = true

	c.emitEvent(EvtWizardComplete, WizardCompletePayload{
		SchemeType: c.cfg.Config.SchemeType,
		SchemeFile: schemeFile,
		DictFile:   dictFile,
		UseMirror:  c.cfg.Config.UseMirror,
	})
}

// handleThemeChange handles theme changes
func (c *Controller) handleThemeChange(cmd Command) {
	payload, ok := cmd.Payload.(ThemeChangePayload)
	if !ok {
		c.emitError(fmt.Errorf("invalid theme change payload"), "theme change")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cfg.Config.ThemeAdaptive {
		if payload.IsLight {
			c.cfg.Config.ThemeLight = payload.ThemeName
		} else {
			c.cfg.Config.ThemeDark = payload.ThemeName
		}
		c.themeMgr.SetAdaptiveTheme(c.cfg.Config.ThemeLight, c.cfg.Config.ThemeDark)
	} else {
		c.cfg.Config.ThemeFixed = payload.ThemeName
		c.themeMgr.SetTheme(payload.ThemeName)
	}

	if err := c.cfg.SaveConfig(); err != nil {
		c.emitError(err, "theme change")
		return
	}

	c.emitEvent(EvtConfigUpdated, ConfigUpdatedPayload{
		Key:   "theme",
		Value: payload.ThemeName,
	})
}

// handleExcludeAdd handles adding an exclude pattern
func (c *Controller) handleExcludeAdd(cmd Command) {
	payload, ok := cmd.Payload.(ExcludePayload)
	if !ok {
		c.emitError(fmt.Errorf("invalid exclude add payload"), "exclude add")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cfg.Config.ExcludeFiles = append(c.cfg.Config.ExcludeFiles, payload.Pattern)

	if err := c.cfg.SaveConfig(); err != nil {
		c.emitError(err, "exclude add")
		return
	}

	c.emitEvent(EvtConfigUpdated, ConfigUpdatedPayload{
		Key:   "exclude_files",
		Value: c.cfg.Config.ExcludeFiles,
	})
}

// handleExcludeRemove handles removing an exclude pattern
func (c *Controller) handleExcludeRemove(cmd Command) {
	payload, ok := cmd.Payload.(ExcludePayload)
	if !ok {
		c.emitError(fmt.Errorf("invalid exclude remove payload"), "exclude remove")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if payload.Index >= 0 && payload.Index < len(c.cfg.Config.ExcludeFiles) {
		c.cfg.Config.ExcludeFiles = append(
			c.cfg.Config.ExcludeFiles[:payload.Index],
			c.cfg.Config.ExcludeFiles[payload.Index+1:]...,
		)

		if err := c.cfg.SaveConfig(); err != nil {
			c.emitError(err, "exclude remove")
			return
		}

		c.emitEvent(EvtConfigUpdated, ConfigUpdatedPayload{
			Key:   "exclude_files",
			Value: c.cfg.Config.ExcludeFiles,
		})
	}
}

// handleExcludeEdit handles editing an exclude pattern
func (c *Controller) handleExcludeEdit(cmd Command) {
	payload, ok := cmd.Payload.(ExcludePayload)
	if !ok {
		c.emitError(fmt.Errorf("invalid exclude edit payload"), "exclude edit")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if payload.Index >= 0 && payload.Index < len(c.cfg.Config.ExcludeFiles) {
		c.cfg.Config.ExcludeFiles[payload.Index] = payload.Pattern

		if err := c.cfg.SaveConfig(); err != nil {
			c.emitError(err, "exclude edit")
			return
		}

		c.emitEvent(EvtConfigUpdated, ConfigUpdatedPayload{
			Key:   "exclude_files",
			Value: c.cfg.Config.ExcludeFiles,
		})
	}
}
