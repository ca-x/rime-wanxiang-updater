# Go Code Review Notes

## Project Overview
- **Project**: rime-wanxiang-updater
- **Go Version**: 1.25.5 (modern, supports all latest features)
- **Type**: TUI application using Bubble Tea framework
- **Purpose**: RIME input method updater with multi-engine support

## Architecture
- `cmd/rime-wanxiang-updater/main.go` - Entry point
- `internal/` packages:
  - `api/` - HTTP client and GitHub/CNB API
  - `config/` - Configuration management
  - `controller/` - Business logic controller
  - `deployer/` - Platform-specific deployment
  - `detector/` - Engine detection
  - `fileutil/` - File operations
  - `types/` - Type definitions
  - `ui/` - Bubble Tea UI
  - `updater/` - Update logic
  - `theme/` - Theme management

## Issues Found

### 1. main.go Violations (Go Best Practices)
**File**: `cmd/rime-wanxiang-updater/main.go`

**Issue**: main.go is NOT a stub - contains 184 lines with boot sequence, styling, and business logic.

**Violation**: Go best practice #1 states "main.go is a stub. Call cmd.Execute() and nothing else."

**Impact**: Medium - Makes testing difficult, mixes presentation with initialization.

**Recommendation**: Move printBootSequence() and all styling to internal/ui or internal/boot package.

### 2. Missing Context Usage (Modern Go)
**Files**: Multiple files in `internal/api/`, `internal/fileutil/`, `internal/updater/`

**Issue**: HTTP requests and long-running operations don't use context.Context.

**Examples**:
- `api/client.go:72` - `Get()` method doesn't accept context
- `api/github.go:46` - `fetchWithRetry()` doesn't use context
- `fileutil/download.go:11` - `DownloadFile()` doesn't accept context
- `updater/base.go:98` - `DownloadFile()` doesn't use context

**Impact**: High - Cannot cancel operations, no timeout control, resource leaks.

**Recommendation**: Add context.Context as first parameter to all I/O operations.

### 3. Error Wrapping Issues (Golang Style)
**Files**: Multiple files

**Issue**: Some errors lack proper context wrapping.

**Examples**:
- `api/client.go:75` - Error message in Chinese, should use English for error types
- `api/github.go:32,39` - Error messages in Chinese
- `fileutil/download.go:20,30,46,56,64` - All error messages in Chinese

**Impact**: Medium - Makes debugging harder, violates Go conventions.

**Recommendation**: Use English for error messages, Chinese for user-facing UI only.

### 4. Channel Operations Without Non-Blocking Pattern (Go Best Practices)
**File**: `internal/controller/controller.go:102-108`

**Issue**: `emitEvent()` uses select with default but doesn't log drops.

**Violation**: Go best practice #4 states "Never block the sender. Log drops at debug level. Document the policy."

**Impact**: Low - Events silently dropped without logging.

**Recommendation**: Add debug logging when events are dropped.

### 5. No Ordered Shutdown (Go Best Practices)
**File**: `cmd/rime-wanxiang-updater/main.go:170`

**Issue**: `ctrl.Stop()` just closes a channel, no ordered shutdown of dependencies.

**Violation**: Go best practice #3 states "Ordered shutdown. Cancel dependents before their dependencies."

**Impact**: Medium - Potential resource leaks, goroutine leaks.

**Recommendation**: Implement proper shutdown with context cancellation and WaitGroup.

### 6. Goroutine Leak Risk (Go Best Practices)
**Files**: `internal/controller/update_handlers.go`

**Issue**: Multiple goroutines spawned (lines 21, 113, 182, 251) without leak detection.

**Examples**:
- `handleAutoUpdate()` - spawns goroutine without tracking
- `handleUpdateDict()` - spawns goroutine without tracking
- `handleUpdateScheme()` - spawns goroutine without tracking
- `handleUpdateModel()` - spawns goroutine without tracking

**Impact**: High - Potential goroutine leaks if controller stops while updates running.

**Recommendation**: Use sync.WaitGroup or context cancellation to track goroutines.

### 7. Modern Go Features Not Used

#### 7.1 sync.WaitGroup.Go() (Go 1.25+)
**Files**: None currently use WaitGroup, but should for goroutine management.

**Issue**: Not using `wg.Go(fn)` pattern for spawning goroutines.

**Recommendation**: When adding WaitGroup tracking, use `wg.Go()` instead of `wg.Add(1)` + `go func()`.

#### 7.2 errors.AsType[T]() (Go 1.26+)
**Files**: No current usage, but could benefit from it.

**Recommendation**: If error type checking is needed, use `errors.AsType[T]()` instead of `errors.As()`.

#### 7.3 new(val) for Pointer Creation (Go 1.26+)
**Files**: Could be used in struct initialization.

**Recommendation**: Use `new(30)` instead of `x := 30; &x` pattern.

### 8. HTTP Client Timeout Issues
**File**: `internal/updater/base.go:101`

**Issue**: Download client has `Timeout: 0` (no timeout).

**Impact**: High - Downloads can hang indefinitely.

**Recommendation**: Use context.Context for cancellation instead of disabling timeout.

### 9. Missing godoc Comments (Effective Go)
**Files**: Multiple files

**Issue**: Some exported functions lack proper godoc comments.

**Examples**:
- `api/client.go:29` - `getHTTPClient` is unexported but called from exported function
- Many exported types and functions have Chinese comments instead of English

**Impact**: Low - Documentation quality.

**Recommendation**: Add English godoc comments for all exported symbols.

### 10. Potential Race Conditions
**File**: `internal/controller/update_handlers.go`

**Issue**: `c.updating` and `c.currentOperation` accessed with mutex, but goroutines may continue after Stop().

**Impact**: Medium - Race conditions if Stop() called during update.

**Recommendation**: Use context cancellation to signal goroutines to stop.

## Positive Findings

1. **Good use of interfaces**: Deployer interface for platform-specific code
2. **Proper error wrapping**: Most errors use `fmt.Errorf` with `%w`
3. **Test coverage**: Test files present for key packages
4. **Modern Go version**: Using Go 1.25.5 enables all modern features
5. **Clean package structure**: Well-organized internal packages
6. **Progress reporting**: Good progress callback pattern

## Priority Improvements

### High Priority
1. Add context.Context to all I/O operations
2. Fix goroutine leak risks with proper tracking
3. Implement ordered shutdown
4. Fix HTTP client timeout issues

### Medium Priority
5. Refactor main.go to be a stub
6. Add logging for dropped events
7. Convert error messages to English
8. Add godoc comments

### Low Priority
9. Use modern Go 1.25/1.26 features where applicable
10. Improve test coverage

## Files Requiring Changes

### Critical
- `cmd/rime-wanxiang-updater/main.go` - Refactor to stub
- `internal/api/client.go` - Add context support
- `internal/api/github.go` - Add context support
- `internal/fileutil/download.go` - Add context support
- `internal/updater/base.go` - Add context support, fix timeout
- `internal/controller/controller.go` - Implement ordered shutdown
- `internal/controller/update_handlers.go` - Fix goroutine leaks
- `internal/ui/model.go` - ✓ FIXED: Wizard not showing when RIME uninstalled

### Important
- `internal/controller/controller.go` - Add event drop logging
- Multiple files - Convert error messages to English

### Nice to Have
- Multiple files - Add/improve godoc comments
- Multiple files - Use modern Go features

## Bug Fixes Applied

### Bug #1: Wizard Not Showing When RIME Data Directory Deleted
**File**: `internal/ui/model.go:25-34`

**Issue**: When user has RIME engine installed (fcitx5/ibus) but deletes the wanxiang data directory, the wizard doesn't show. The app goes straight to menu and runs "updates" that do nothing useful.

**Root Cause**:
1. `detector.CheckRimeInstallation()` returns `Installed: true` because it detects the engine command (fcitx5, ibus-daemon)
2. `DetectInstalledEngines()` in config falls back to `["fcitx5"]` default
3. App config still has `SchemeType/SchemeFile/DictFile` populated
4. Wizard only triggered when those config fields are empty

**Fix**: Added two additional conditions to wizard trigger:
1. `!rimeStatus.Installed` - RIME engine not installed
2. `rimeDirMissing` - RIME user data directory doesn't exist

Now wizard shows when ANY of these are true:
- Config fields are empty (SchemeType/SchemeFile/DictFile)
- RIME engine is not installed
- RIME user data directory is missing

**Code Change**:
```go
_, statErr := os.Stat(cfg.RimeDir)
rimeDirMissing := cfg.RimeDir == "" || os.IsNotExist(statErr)
if cfg.Config.SchemeType == "" || cfg.Config.SchemeFile == "" || cfg.Config.DictFile == "" || !rimeStatus.Installed || rimeDirMissing {
    state = ViewWizard
}
```

**Impact**: Users who delete their RIME data directory will now see the setup wizard instead of the app silently "succeeding" with no actual setup.
