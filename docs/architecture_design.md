# Architecture Design: UI-Logic Separation with Channels

## Overview

This document defines the architecture for separating UI presentation from business logic using channel-based communication in the rime-wanxiang-updater project.

---

## Component Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Main Entry Point                        │
│                    (cmd/rime-wanxiang-updater)                  │
└────────────────┬────────────────────────┬───────────────────────┘
                 │                        │
                 v                        v
    ┌────────────────────┐    ┌──────────────────────┐
    │    UI Layer        │    │   Controller Layer   │
    │   (BubbleTea)      │    │   (Business Logic)   │
    │                    │    │                      │
    │ - Views            │    │ - State Management   │
    │ - Input Handling   │    │ - Update Orchestration│
    │ - Rendering        │    │ - Config Management  │
    │ - UI State         │    │ - Validation         │
    └─────────┬──────────┘    └──────────┬───────────┘
              │                          │
              │  CommandMsg (buffered)   │
              ├─────────────────────────>│
              │                          │
              │  EventMsg (buffered)     │
              │<─────────────────────────┤
              │                          │
              v                          v
       tea.Program.Run()          controller.Run()
                                   (goroutine)
```

---

## Package Structure

### New Structure
```
internal/
├── controller/              # NEW: Business logic layer
│   ├── controller.go        # Main controller with Run() loop
│   ├── types.go             # Controller state and types
│   ├── messages.go          # Command and Event definitions
│   ├── update_handlers.go   # Update operation handlers
│   ├── config_handlers.go   # Configuration handlers
│   └── wizard_handlers.go   # Wizard logic handlers
│
├── ui/                      # REFACTORED: UI presentation only
│   ├── model.go             # Minimal BubbleTea Model
│   ├── types.go             # UI-only state types
│   ├── handlers.go          # Input handlers (no business logic)
│   ├── views.go             # View rendering
│   ├── styles.go            # Styling (unchanged)
│   ├── themed_styles.go     # Themed styling (unchanged)
│   ├── exclude_manager.go   # Exclude UI (refactored)
│   └── theme_selector.go    # Theme UI (refactored)
│
└── (existing packages remain unchanged)
    ├── updater/
    ├── config/
    ├── detector/
    ├── deployer/
    ├── api/
    ├── theme/
    └── fileutil/
```

---

## Message Types

### Command Messages (UI → Controller)

Commands represent actions initiated by the user through the UI.

```go
// CommandType defines the type of command
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
    Payload interface{} // Command-specific data
}

// Specific command payloads
type ConfigChangePayload struct {
    Key   string
    Value interface{}
}

type WizardSchemePayload struct {
    SchemeType string
    Variant    string
}

type ThemeChangePayload struct {
    ThemeName string
    IsLight   bool
}
```

### Event Messages (Controller → UI)

Events represent state changes or results that the UI should display.

```go
// EventType defines the type of event
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
    Payload interface{} // Event-specific data
}

// Specific event payloads
type ProgressUpdatePayload struct {
    Component    string
    Message      string
    Percent      float64
    Source       string
    FileName     string
    Downloaded   int64
    TotalSize    int64
    Speed        float64
    IsDownload   bool
}

type UpdateCompletePayload struct {
    UpdateType        string
    Success           bool
    Skipped           bool
    Message           string
    UpdatedComponents []string
    SkippedComponents []string
    ComponentVersions map[string]string
}

type ConfigUpdatedPayload struct {
    Config *config.Config
}

type ErrorPayload struct {
    Error   error
    Context string
}
```

---

## Controller Design

### Controller Struct

```go
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
    wizardState      WizardState

    // Communication channels
    commandChan chan Command
    eventChan   chan Event

    // Shutdown
    done chan struct{}

    // Internal state
    mu sync.RWMutex // Protects shared state
}

// WizardState tracks wizard progress in controller
type WizardState struct {
    SchemeType   string
    Variant      string
    UseMirror    bool
    Completed    bool
}
```

### Controller Methods

```go
// NewController creates a new controller instance
func NewController(
    cfg *config.Manager,
    commandChan chan Command,
    eventChan chan Event,
) *Controller

// Run starts the controller's main loop (runs in goroutine)
func (c *Controller) Run()

// Stop gracefully stops the controller
func (c *Controller) Stop()

// Command handlers (internal)
func (c *Controller) handleCommand(cmd Command)
func (c *Controller) handleAutoUpdate(cmd Command)
func (c *Controller) handleUpdateDict(cmd Command)
func (c *Controller) handleUpdateScheme(cmd Command)
func (c *Controller) handleUpdateModel(cmd Command)
func (c *Controller) handleConfigChange(cmd Command)
func (c *Controller) handleWizardCommand(cmd Command)
func (c *Controller) handleThemeChange(cmd Command)

// Event emission helpers
func (c *Controller) emitEvent(eventType EventType, payload interface{})
func (c *Controller) emitProgress(component, message string, percent float64, ...)
func (c *Controller) emitError(err error, context string)
```

---

## UI Model Design

### Minimal UI Model

```go
// Model represents the UI state only
type Model struct {
    // Communication with controller
    commandChan chan<- controller.Command
    eventChan   <-chan controller.Event

    // UI-only state
    state            ViewState
    wizardStep       WizardStep
    menuChoice       int
    configChoice     int
    editingKey       string
    editingValue     string

    // Visual state
    width  int
    height int

    // Theme and styling
    themeMgr *theme.Manager
    styles   *Styles

    // Progress display (received from controller)
    progress        progress.Model
    progressMsg     string
    isDownloading   bool
    downloadSource  string
    downloadFile    string
    downloaded      int64
    totalSize       int64
    downloadSpeed   float64

    // Result display (received from controller)
    resultMsg     string
    resultSuccess bool
    resultSkipped bool
    autoUpdateDetails *AutoUpdateDetails

    // Exclude manager state (UI only)
    excludeListChoice   int
    excludeEditInput    string
    excludeEditIndex    int
    excludeErrorMsg     string
    excludeDescriptions []string

    // Fcitx conflict dialog state
    fcitxConflictChoice   int
    fcitxConflictNoPrompt bool

    // Theme selector state
    themeListChoice int
    themeList       []string

    // Auto update countdown (UI display only)
    autoUpdateCountdown int
    autoUpdateCancelled bool

    // Error display (received from controller)
    err error

    // Config reference (read-only, for display)
    cfg *config.Manager
}
```

---

## Communication Patterns

### 1. Command Pattern (UI → Controller)

```go
// In UI handler
func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "1": // Auto update
        // Send command to controller
        cmd := controller.Command{
            Type: controller.CmdAutoUpdate,
        }
        return m, m.sendCommand(cmd)
    }
    return m, nil
}

// Helper method
func (m Model) sendCommand(cmd controller.Command) tea.Cmd {
    return func() tea.Msg {
        m.commandChan <- cmd
        return nil
    }
}
```

### 2. Event Pattern (Controller → UI)

```go
// In UI Update() method
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case controller.Event:
        return m.handleControllerEvent(msg)
    // ... other cases
    }
    return m, nil
}

// Event handler
func (m Model) handleControllerEvent(evt controller.Event) (tea.Model, tea.Cmd) {
    switch evt.Type {
    case controller.EvtProgressUpdate:
        payload := evt.Payload.(controller.ProgressUpdatePayload)
        m.progressMsg = payload.Message
        m.isDownloading = payload.IsDownload
        // Update progress display
        return m, m.progress.SetPercent(payload.Percent)

    case controller.EvtUpdateSuccess:
        payload := evt.Payload.(controller.UpdateCompletePayload)
        m.resultSuccess = true
        m.resultMsg = payload.Message
        m.state = ViewResult
        return m, nil
    }
    return m, nil
}
```

### 3. Event Listening Pattern

```go
// In UI Init() or as continuous command
func listenForEvents(eventChan <-chan controller.Event) tea.Cmd {
    return func() tea.Msg {
        return <-eventChan // Blocks until event arrives
    }
}

// Recursive listening (continuous)
func (m Model) Init() tea.Cmd {
    return listenForEvents(m.eventChan)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case controller.Event:
        // Handle event
        newModel, cmd := m.handleControllerEvent(msg)
        // Continue listening
        return newModel, tea.Batch(cmd, listenForEvents(m.eventChan))
    }
    return m, nil
}
```

---

## Controller Main Loop

```go
func (c *Controller) Run() {
    for {
        select {
        case cmd := <-c.commandChan:
            // Handle commands from UI
            c.handleCommand(cmd)

        case <-c.done:
            // Graceful shutdown
            return
        }
    }
}
```

---

## Initialization Sequence

```go
// In main.go
func main() {
    // 1. Load configuration
    cfg, err := config.NewManager()
    if err != nil {
        log.Fatal(err)
    }

    // 2. Create communication channels
    commandChan := make(chan controller.Command, 10) // Buffered
    eventChan := make(chan controller.Event, 100)    // Buffered for progress

    // 3. Create controller
    ctrl := controller.NewController(cfg, commandChan, eventChan)

    // 4. Start controller in goroutine
    go ctrl.Run()

    // 5. Create UI model
    model := ui.NewModel(cfg, commandChan, eventChan)

    // 6. Run BubbleTea program
    p := tea.NewProgram(model)
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }

    // 7. Stop controller
    ctrl.Stop()
}
```

---

## State Flow Examples

### Example 1: Auto Update Flow

```
User Input (Press 1)
    │
    ├─> UI: handleMenuInput()
    │       ├─> Create Command{Type: CmdAutoUpdate}
    │       └─> Send to commandChan
    │
    ├─> Controller: Receive command
    │       ├─> handleAutoUpdate()
    │       ├─> Emit Event{Type: EvtProgressUpdate, "Starting..."}
    │       ├─> Run updaters
    │       ├─> Emit Event{Type: EvtProgressUpdate, "50%..."}
    │       └─> Emit Event{Type: EvtUpdateSuccess, result}
    │
    └─> UI: Receive events
            ├─> Update progress display
            ├─> Update state to ViewResult
            └─> Render result
```

### Example 2: Configuration Change Flow

```
User Input (Edit config)
    │
    ├─> UI: handleConfigEditInput()
    │       ├─> User enters value
    │       ├─> Press Enter
    │       ├─> Create Command{Type: CmdConfigChange, key, value}
    │       └─> Send to commandChan
    │
    ├─> Controller: Receive command
    │       ├─> Validate value
    │       ├─> Update config
    │       ├─> Save config
    │       └─> Emit Event{Type: EvtConfigUpdated, config}
    │
    └─> UI: Receive event
            ├─> Update display
            ├─> Return to config view
            └─> Show success
```

---

## Channel Buffer Sizes

Based on current usage patterns:

- **Command Channel**: Buffered with size 10
  - Low frequency (user input only)
  - No blocking concerns

- **Event Channel**: Buffered with size 100
  - High frequency during downloads (progress updates)
  - Must not block controller
  - UI processes as fast as possible

---

## Error Handling Strategy

1. **Controller Errors**:
   - Emit `EvtError` event with error details
   - Continue running (don't crash)
   - UI displays error to user

2. **Channel Errors**:
   - Use timeouts for sends if needed
   - Graceful degradation
   - Log but don't crash

3. **Shutdown**:
   - Controller drains command channel before stopping
   - Close event channel on shutdown
   - UI handles closed channel gracefully

---

## Concurrency Safety

1. **Controller State**:
   - Protected by `sync.RWMutex`
   - Read lock for queries
   - Write lock for modifications

2. **Channels**:
   - Buffered to prevent blocking
   - Single writer per channel
   - Multiple readers okay (if needed)

3. **Config Access**:
   - Config manager is thread-safe
   - Controller owns config modifications
   - UI has read-only reference

---

## Benefits of This Design

1. **Separation of Concerns**:
   - UI only knows about presentation
   - Controller only knows about business logic
   - Clear boundaries

2. **Testability**:
   - Controller can be tested without UI
   - Mock channels for unit tests
   - Integration tests easy to write

3. **Maintainability**:
   - Each component has single responsibility
   - Easy to modify UI without touching logic
   - Easy to add new commands/events

4. **Flexibility**:
   - UI can be swapped (TUI → GUI)
   - Controller can be reused
   - Easy to add new features

5. **Robustness**:
   - Graceful error handling
   - No shared mutable state
   - Clear ownership of data

---

## Migration Strategy

The implementation will follow the phased approach in task_plan.md:

1. Create controller package with message types
2. Move update operations to controller
3. Move config operations to controller
4. Refactor UI to send commands/receive events
5. Test each component independently
6. Integrate and test end-to-end

---

## Next Steps

- [x] Architecture design complete
- [ ] Proceed to Phase 2: Implementation
- [ ] Create controller package structure
- [ ] Define message types in code
- [ ] Implement controller foundation
