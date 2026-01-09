package controller

// CommandType defines the type of command sent from UI to Controller
type CommandType int

const (
	// Update commands
	CmdAutoUpdate CommandType = iota
	CmdUpdateDict
	CmdUpdateScheme
	CmdUpdateModel

	// Configuration commands
	CmdConfigChange
	CmdConfigSave
	CmdConfigCancel

	// Wizard commands
	CmdWizardSetScheme
	CmdWizardSetMirror
	CmdWizardComplete

	// View commands
	CmdChangeView

	// Theme commands
	CmdThemeChange

	// Exclude file commands
	CmdExcludeAdd
	CmdExcludeRemove
	CmdExcludeEdit

	// System commands
	CmdShutdown
)

// Command is the message sent from UI to Controller
type Command struct {
	Type    CommandType
	Payload any
}

// ConfigChangePayload contains data for configuration changes
type ConfigChangePayload struct {
	Key   string
	Value any
}

// WizardSchemePayload contains data for wizard scheme selection
type WizardSchemePayload struct {
	SchemeType string
	Variant    string
}

// ThemeChangePayload contains data for theme changes
type ThemeChangePayload struct {
	ThemeName string
	IsLight   bool
}

// ExcludePayload contains data for exclude file operations
type ExcludePayload struct {
	Pattern     string
	Description string
	Index       int
}

// EventType defines the type of event sent from Controller to UI
type EventType int

const (
	// State update events
	EvtStateUpdate EventType = iota

	// Progress events
	EvtProgressUpdate
	EvtProgressComplete

	// Result events
	EvtUpdateSuccess
	EvtUpdateFailure
	EvtUpdateSkipped

	// Configuration events
	EvtConfigUpdated
	EvtConfigError

	// Wizard events
	EvtWizardComplete
	EvtWizardError

	// Error events
	EvtError

	// View navigation events
	EvtNavigateToView
)

// Event is the message sent from Controller to UI
type Event struct {
	Type    EventType
	Payload any
}

// ProgressUpdatePayload contains progress information
type ProgressUpdatePayload struct {
	Component  string
	Message    string
	Percent    float64
	Source     string
	FileName   string
	Downloaded int64
	TotalSize  int64
	Speed      float64
	IsDownload bool
}

// UpdateCompletePayload contains update completion information
type UpdateCompletePayload struct {
	UpdateType        string
	Success           bool
	Skipped           bool
	Message           string
	UpdatedComponents []string
	SkippedComponents []string
	ComponentVersions map[string]string
}

// ConfigUpdatedPayload contains updated configuration
type ConfigUpdatedPayload struct {
	Key   string
	Value any
}

// ErrorPayload contains error information
type ErrorPayload struct {
	Error   error
	Context string
}

// WizardCompletePayload contains wizard completion information
type WizardCompletePayload struct {
	SchemeType string
	SchemeFile string
	DictFile   string
	UseMirror  bool
}
