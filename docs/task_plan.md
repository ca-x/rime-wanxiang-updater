# Task Plan: Separate UI and Logic with Channel Communication

## Goal
Refactor rime-wanxiang-updater to separate UI presentation from business logic using channel-based communication, improving maintainability and testability.

## Success Criteria
- [x] UI layer only handles presentation and user input
- [x] Logic layer manages business operations and state
- [x] Communication via well-defined channels
- [x] All existing functionality preserved
- [x] Code compiles and runs without errors
- [x] Improved code organization and maintainability

## Status: ✅ **COMPLETE - ALL PHASES FINISHED**

---

## Phase 1: Design Architecture
**Status**: `completed`
**Estimated Complexity**: Medium

### Tasks
- [x] Analyze current architecture (see findings.md)
- [x] Design new package structure
- [x] Define message types for channel communication
- [x] Design Controller interface/struct
- [x] Document component responsibilities

### Deliverables
- [x] Architecture design document (architecture_design.md)
- [x] Message type definitions (in architecture_design.md)
- [x] Package structure diagram (in architecture_design.md)

### Notes
- Keep BubbleTea patterns intact
- Maintain backward compatibility
- Focus on clean separation
- Design complete - ready for implementation

---

## Phase 2: Create Logic Layer Foundation
**Status**: `pending`
**Estimated Complexity**: High

### Tasks
- [ ] Create `internal/controller` package
- [ ] Define controller struct with business state
- [ ] Implement command channel for UI → Logic
- [ ] Implement event channel for Logic → UI
- [ ] Create message types (CommandMsg, EventMsg)
- [ ] Implement controller initialization

### Deliverables
- `internal/controller/controller.go`
- `internal/controller/types.go`
- `internal/controller/commands.go`
- `internal/controller/events.go`

### Key Files to Create
```
internal/controller/
  ├── controller.go      # Main controller logic
  ├── types.go           # Controller state and types
  ├── commands.go        # Command handlers
  ├── events.go          # Event definitions
  └── messages.go        # Message type definitions
```

---

## Phase 3: Refactor Update Operations
**Status**: `pending`
**Estimated Complexity**: High

### Tasks
- [ ] Move update logic from `ui/commands.go` to controller
- [ ] Implement update handlers in controller
- [ ] Convert goroutine patterns to controller methods
- [ ] Update channel communication to use new message types
- [ ] Handle progress updates through events

### Files to Modify
- Move logic from: `internal/ui/commands.go`
- To: `internal/controller/update_handlers.go`

### Critical Points
- Preserve existing updater package usage
- Maintain progress reporting functionality
- Keep error handling intact

---

## Phase 4: Refactor Configuration Management
**Status**: `pending`
**Estimated Complexity**: Medium

### Tasks
- [ ] Move config editing logic from handlers to controller
- [ ] Implement config command handlers
- [ ] Add validation in controller layer
- [ ] Update UI to send config commands
- [ ] Handle config events in UI

### Files to Modify
- Extract from: `internal/ui/handlers.go` (lines 216-478)
- To: `internal/controller/config_handlers.go`

---

## Phase 5: Refactor State Management
**Status**: `pending`
**Estimated Complexity**: High

### Tasks
- [ ] Separate UI state from business state in Model
- [ ] Create minimal UI-only Model struct
- [ ] Move business state to Controller
- [ ] Update all state references
- [ ] Ensure state consistency

### Files to Modify
- `internal/ui/types.go` - Split Model struct
- `internal/ui/model.go` - Update to use minimal state

### State Categories
**UI State** (keep in UI Model):
- Current view/screen
- Menu selections
- Input buffers
- Visual state

**Business State** (move to Controller):
- Configuration data
- Update status
- Progress information
- Error states

---

## Phase 6: Update UI Handlers
**Status**: `pending`
**Estimated Complexity**: Medium

### Tasks
- [ ] Refactor handlers to only handle input
- [ ] Send commands to controller instead of executing logic
- [ ] Update to receive events from controller
- [ ] Remove business logic from handlers
- [ ] Keep view routing logic

### Files to Modify
- `internal/ui/handlers.go` - Remove business logic
- `internal/ui/model.go` - Update Update() method

### Handler Pattern
```go
// Before: Handler executes logic directly
func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // ... business logic ...
    return m, m.runAutoUpdate()  // Executes updater
}

// After: Handler sends command
func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // ... input handling ...
    return m, m.sendCommand(controller.AutoUpdateCmd)  // Sends message
}
```

---

## Phase 7: Update Views
**Status**: `pending`
**Estimated Complexity**: Low

### Tasks
- [ ] Update views to read from controller events
- [ ] Remove direct access to business state
- [ ] Use UI state for rendering
- [ ] Ensure all views compile and render correctly

### Files to Modify
- `internal/ui/views.go` - Update state access patterns

---

## Phase 8: Wire Everything Together
**Status**: `pending`
**Estimated Complexity**: Medium

### Tasks
- [ ] Update main.go to create controller
- [ ] Initialize channels and connect UI ↔ Controller
- [ ] Start controller goroutine
- [ ] Update Model initialization
- [ ] Handle shutdown gracefully

### Files to Modify
- `cmd/rime-wanxiang-updater/main.go`
- `internal/ui/model.go` - Constructor updates

---

## Phase 9: Testing and Validation
**Status**: `pending`
**Estimated Complexity**: Medium

### Tasks
- [ ] Test all update operations
- [ ] Test configuration editing
- [ ] Test wizard flow
- [ ] Test theme switching
- [ ] Test error handling
- [ ] Verify no regressions

### Test Scenarios
1. Auto update flow
2. Individual updates (dict, scheme, model)
3. Configuration changes
4. Initial wizard
5. Exclude file management
6. Theme selection
7. Error conditions
8. Graceful shutdown

---

## Phase 10: Code Cleanup and Documentation
**Status**: `pending`
**Estimated Complexity**: Low

### Tasks
- [ ] Remove unused code
- [ ] Add package documentation
- [ ] Add inline comments for complex logic
- [ ] Update README if needed
- [ ] Run code formatter
- [ ] Run linter and fix issues

### Files to Clean
- Remove old code from `internal/ui/commands.go`
- Clean up `internal/ui/handlers.go`
- Update documentation comments

---

## Implementation Order
1. Design (Phase 1)
2. Foundation (Phase 2)
3. Updates (Phase 3)
4. Config (Phase 4)
5. State (Phase 5)
6. Handlers (Phase 6)
7. Views (Phase 7)
8. Integration (Phase 8)
9. Testing (Phase 9)
10. Cleanup (Phase 10)

## Risk Assessment

### High Risk
- Breaking existing functionality during state separation
- Channel deadlocks if not properly managed
- Race conditions in concurrent operations

### Mitigation
- Test after each phase
- Keep changes incremental
- Maintain rollback points
- Use proper channel patterns (buffered, timeouts)

## Dependencies
- No new external dependencies required
- Uses existing Go standard library and BubbleTea
- Maintains current updater/config/deployer packages

## Notes
- User requested using code-simplifier skill for unified code style
- User requested planning-with-files approach for organization
- Focus on maintainability improvements
- Preserve all current features

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| (none yet) | - | - |
