# Findings: Current Architecture Analysis

## Discovery Date
2026-01-10

## Current Architecture Overview

### UI Package Structure (`internal/ui/`)
- **model.go**: BubbleTea Model with Init(), Update(), View() methods
- **types.go**: Model struct definition with extensive state fields
- **handlers.go**: Input handlers mixed with business logic
- **commands.go**: Updater execution functions (runDictUpdate, runSchemeUpdate, etc.)
- **views.go**: View rendering functions
- **styles.go**: Styling definitions
- **exclude_manager.go**: Exclude file management UI
- **theme_selector.go**: Theme selection UI

### Key Issues Identified

#### 1. **Tight Coupling Between UI and Logic**
- The `Model` struct (types.go:39-93) contains:
  - UI state (ViewState, WizardStep, MenuChoice, etc.)
  - Business data (Config, Progress, RimeInstallStatus)
  - Channel communication (ProgressChan, CompletionChan)
  - Update execution logic

#### 2. **Mixed Responsibilities**
- `handlers.go` contains both:
  - Input handling (keyboard events)
  - Business logic (config editing, validation)
  - State transitions
- `commands.go` contains update execution logic directly in UI package

#### 3. **Channel Usage**
Current channel usage (types.go:61-62):
```go
ProgressChan     chan UpdateMsg         // 进度通道
CompletionChan   chan UpdateCompleteMsg // 完成通道
```
These are created in command functions (commands.go:14-15) and passed through the Model

#### 4. **Update Flow**
Current flow for updates (e.g., commands.go:13-58):
1. User selects menu option → Handler called
2. Handler creates channels and stores in Model
3. Handler launches goroutine with updater logic
4. Goroutine sends progress via channel
5. BubbleTea listens and updates UI via Update() method

### Domain Packages (Well-Structured)
- **internal/updater/**: Update logic (dict, scheme, model, combined)
- **internal/config/**: Configuration management
- **internal/detector/**: Rime installation detection
- **internal/deployer/**: Deployment logic
- **internal/api/**: GitHub and CNB API clients
- **internal/theme/**: Theme management
- **internal/fileutil/**: File operations

### Entry Point
- **cmd/rime-wanxiang-updater/main.go**:
  - Boot sequence animation
  - Config loading
  - BubbleTea program initialization

## Refactoring Goals

### Separation Strategy
1. **UI Layer** (presentation only)
   - View rendering
   - Input handling
   - State display

2. **Logic Layer** (business logic)
   - Update orchestration
   - Configuration management
   - Validation
   - State management

3. **Communication** (channels)
   - Command channel: UI → Logic
   - Event channel: Logic → UI
   - Clean message passing

### Benefits
- **Maintainability**: Clear separation of concerns
- **Testability**: Logic can be tested independently
- **Flexibility**: UI can be swapped (TUI, GUI, CLI)
- **Clarity**: Each component has single responsibility

## File Structure Analysis

### Files to Refactor
1. **internal/ui/model.go** - Split into UI model and logic controller
2. **internal/ui/handlers.go** - Extract business logic
3. **internal/ui/commands.go** - Move to logic layer
4. **internal/ui/types.go** - Separate UI state from business state

### Files to Keep (Mostly Unchanged)
- **internal/ui/views.go** - Pure rendering, minor changes
- **internal/ui/styles.go** - No changes
- **internal/updater/** - Already well-structured
- **internal/config/** - Already well-structured

## Technical Notes

### BubbleTea Architecture
- Uses Elm Architecture (Model-Update-View)
- Commands return `tea.Cmd` functions
- Messages flow through Update() method
- Current implementation mixes UI and logic in Update()

### Channel Patterns
- Buffered channels for progress (100 buffer)
- Unbuffered channels for completion (1 buffer)
- listenForProgress() wraps channel reads as tea.Cmd

### Go Patterns Observed
- goroutines for async operations
- select statements for channel multiplexing
- defer close() for channel cleanup
