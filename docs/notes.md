# Notes: Go Code Review

## Code Quality Analysis (Go 1.25.5)

### 1. Modern Go Features Not Utilized

#### slices package (Go 1.21+)
- `config/config.go:93-97`: Manual loop to check if engine exists - should use `slices.Contains`
- `config/config.go:178-183`: Same pattern in `RedetectEngines()`
- `config/config.go:121-128`: Manual loop to filter slice - could use `slices.DeleteFunc` pattern

#### cmp package (Go 1.22+)
- `updater/base.go:37-38`: Can use `cmp.Or` for nil checks with defaults
- `ui/model.go:36-39`: Countdown default can use `cmp.Or`

#### omitzero tag (Go 1.24+)
- `types/types.go`: JSON tags use `omitempty` for `time.Time` fields - should use `omitzero`
  - `UpdateInfo.UpdateTime`, `UpdateRecord.UpdateTime`, `UpdateRecord.ApplyTime`
  - `GitHubRelease.PublishedAt`, `GitHubAsset.UpdatedAt`, `CNBAsset.UpdatedAt`

### 2. Error Handling Issues

#### Missing error wrapping
- `config/config.go:157`: `os.MkdirAll` error not wrapped
- `config/config.go:346`: `os.MkdirAll` error ignored (should be checked)
- `updater/base.go:126`: `resp.Body.Close()` error not checked

#### Error context missing
- `updater/base.go:65-72`: Returns nil without context when errors occur
- `fileutil/` files need review for error wrapping

### 3. Code Style Issues

#### Happy path violations
- `config/config.go:54-61`: Nested if for file not exists - should invert condition
- `updater/base.go:138-149`: Nested conditions for file handling

#### Line length > 120 chars
- `types/types.go:120`: ProgressFunc signature too long
- `config/config.go:474-476`: Error formatting too long

#### Comments not ending with period
- Various Chinese comments throughout (acceptable but inconsistent)

### 4. Interface Design Issues

- `deployer/deployer.go:10`: `GetDeployer(config interface{})` - uses `any`/`interface{}` instead of typed parameter
  - Should be `*types.Config` for type safety

### 5. Potential Sentinel Errors

Package-level errors that should be defined:
```go
// config/config.go
var ErrNoEngineDetected = errors.New("no input method engine detected")

// updater/base.go
var ErrFileNotFound = errors.New("file not found")
```

### 6. Test Improvements (Go 1.24+)

- Tests should use `t.Context()` instead of `context.WithCancel(context.Background())`
- Benchmarks should use `b.Loop()` pattern

### 7. Concurrency Patterns

- `ui/model.go:257-263`: Non-blocking channel send is correct, but lacks documentation of drop policy
- `updater/base.go:99-110`: HTTP client created per download call - could be cached

### 8. Deprecated Patterns

- `config/config.go:326-334`: `detectEngine()` marked deprecated but still exists - should be removed if truly unused

## Files Requiring Changes

| Priority | File | Issue | Status |
|----------|------|-------|--------|
| High | `internal/types/types.go` | `omitzero` for time fields | ✅ Done |
| High | `internal/config/config.go` | `slices.Contains`, error wrapping | ✅ Done |
| Medium | `internal/updater/base.go` | Error handling, modern patterns | Pending |
| Medium | `internal/deployer/deployer.go` | Typed interface parameter | ✅ Done |
| Low | `internal/ui/model.go` | `cmp.Or` for defaults | ✅ Done |
| Low | `internal/api/client.go` | Minor style improvements | Pending |

---

## Implementation Summary (2026-03-27)

### Changes Applied

1. **types/types.go** - Updated JSON tags from `omitempty` to `omitzero` for `time.Time` fields (Go 1.24+ feature)

2. **config/config.go**
   - Added `slices` and `errors` imports
   - Added sentinel errors: `ErrNoEngineDetected`, `ErrPatternExists`, `ErrIndexOutOfRange`
   - Replaced manual loops with `slices.Contains` in 3 locations
   - Fixed error wrapping in `saveConfig()` and `AddExcludePattern()`

3. **deployer/deployer.go + platform files**
   - Changed `GetDeployer(config interface{})` to `GetDeployer(config *types.Config)`
   - Removed type assertion boilerplate in `darwin.go`, `linux.go`, `windows.go`

4. **ui/model.go**
   - Added `cmp` import
   - Used `cmp.Or(cfg.Config.AutoUpdateCountdown, 5)` for default value

### Build Status
- `go build ./...` - ✅ Passes
- `go test ./...` - ⚠️ API tests fail due to external API response changes (not related to our changes)
- `config` package tests - ✅ Pass