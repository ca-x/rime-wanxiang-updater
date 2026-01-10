# 研究发现 - 多输入法引擎支持

## 任务背景
支持 Mac 下 FCITX5 (https://github.com/fcitx-contrib/fcitx5-macos-installer) 和鼠须管 (Squirrel) 并存，需要调整安装部署和检测架构。Linux 和 Windows 也可能有类似情况。

---

## 当前架构分析

### 1. Engine 检测逻辑 (`internal/config/config.go:175-185`)
```go
func detectEngine() string {
	switch runtime.GOOS {
	case "windows":
		return "小狼毫"
	case "darwin":
		return "鼠须管"
	default:
		return "fcitx5"
	}
}
```
**问题：**
- 基于操作系统简单假设，无实际检测
- 不支持同一系统多个输入法并存
- 返回单个引擎字符串

### 2. macOS 路径映射 (`internal/config/paths_darwin.go:13-22`)
```go
func getRimeUserDir(config *types.Config) string {
	homeDir, _ := os.UserHomeDir()

	if config.Engine == "鼠须管" {
		return filepath.Join(homeDir, "Library", "Rime")
	}

	// 小企鹅或其他
	return filepath.Join(homeDir, ".local", "share", "fcitx5", "rime")
}
```
**发现：**
- ✅ 已经知道两个路径位置
- ❌ 基于 config.Engine 单选，不支持多引擎
- 鼠须管: `~/Library/Rime`
- FCITX5: `~/.local/share/fcitx5/rime`

### 3. Linux 已有多路径检测 (`internal/deployer/linux.go:58-83`)
```go
func (d *linuxDeployer) getRimeDataDir() (string, error) {
	candidates := []string{
		filepath.Join(homeDir, ".local/share/fcitx5/rime"),
		filepath.Join(homeDir, ".config/fcitx5/rime"),
		filepath.Join(homeDir, ".config/ibus/rime"),
		filepath.Join(homeDir, ".config/fcitx/rime"),
	}

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, nil
		}
	}
	// ...
}
```
**启示：**
- ✅ Linux deployer 已经实现了多路径检测和优先级
- 支持 fcitx5, ibus, fcitx (v4) 三个引擎
- 可以作为多引擎检测的参考实现

### 4. macOS 部署器 (`internal/deployer/darwin.go:28-46`)
```go
func (d *darwinDeployer) Deploy() error {
	var executable string
	var args []string

	if d.engine == "鼠须管" {
		executable = "/Library/Input Methods/Squirrel.app/Contents/MacOS/Squirrel"
		args = []string{"--reload"}
	} else {
		executable = "/Library/Input Methods/Fcitx5.app/Contents/bin/fcitx5-curl"
		args = []string{"/config/addon/rime/deploy", "-X", "POST", "-d", "{}"}
	}
	// ...
}
```
**发现：**
- ✅ 已有 FCITX5 部署逻辑框架
- ❌ 但构造函数硬编码为"鼠须管" (line 17-18)
- ✅ 知道两个引擎的部署命令
- ⚠️ **路径错误**: 代码写的是 `Fcitx5.app` 但实际应该是 `Fcitx5Installer.app`

### 5. 配置结构 (`internal/types/types.go:31-56`)
```go
type Config struct {
	Engine              string   `json:"engine"`  // 单个引擎字符串
	// ...
}
```
**问题：**
- `Engine` 字段是单个字符串
- 无法表示多个已安装引擎
- 无法指定主次引擎或优先级

---

## 各平台可能的多引擎情况

### macOS
- ✅ **已确认**: Squirrel (鼠须管) + FCITX5 (小企鹅) 可以并存
- **小企鹅版本**: 有三个发行版（拼音版、中州韵版、原装版），但安装后位置相同
- **系统要求**: macOS >= 13
- **路径**:
  - Squirrel (鼠须管): `/Library/Input Methods/Squirrel.app`, 数据: `~/Library/Rime`
  - FCITX5 (小企鹅): `/Library/Input Methods/Fcitx5Installer.app`, 数据: `~/.local/share/fcitx5/rime` ✅ 已确认
  - ⚠️ 注意: 代码中 `darwin.go:36` 路径需要修正

### Linux
- ✅ **已知**: fcitx5, ibus, fcitx(v4) 可能并存
- 部署器已支持多路径检测
- 路径优先级（deployer/linux.go:67-72）:
  1. `~/.local/share/fcitx5/rime`
  2. `~/.config/fcitx5/rime`
  3. `~/.config/ibus/rime`
  4. `~/.config/fcitx/rime`

### Windows
- ⚠️ **待确认**: 小狼毫是否可能与其他 Rime 实现并存
- 当前只检测小狼毫 (Weasel)
- 路径通过注册表读取: `HKCU\Software\Rime\Weasel\RimeUserDir`

---

## 关键文件列表

| 文件 | 功能 | 改动需求 |
|------|------|---------|
| `internal/types/types.go` | 配置结构定义 | 🔴 需修改 Engine 字段 |
| `internal/config/config.go` | 配置管理，detectEngine() | 🔴 需重构检测逻辑 |
| `internal/config/paths_darwin.go` | macOS 路径检测 | 🔴 需支持多引擎 |
| `internal/config/paths_linux.go` | Linux 路径检测 | 🟡 参考 deployer 实现 |
| `internal/config/paths_windows.go` | Windows 路径检测 | 🟢 评估是否需要 |
| `internal/deployer/deployer.go` | 部署器接口 | 🟡 可能需调整 |
| `internal/deployer/darwin.go` | macOS 部署实现 | 🔴 需支持多引擎部署 |
| `internal/deployer/linux.go` | Linux 部署实现 | 🟢 已有多路径支持 |
| `internal/deployer/windows.go` | Windows 部署实现 | 🟢 评估是否需要 |
| `internal/ui/*` | UI 层 | 🔴 需显示/选择引擎 |

---

## 技术决策点

### Q1: 配置中如何表示多引擎？
**选项 A**: 保持单个 `Engine` 字符串，改用引擎数组
```go
InstalledEngines []string `json:"installed_engines"` // 检测到的所有引擎
PrimaryEngine    string   `json:"primary_engine"`    // 用户选择的主引擎
```

**选项 B**: 更复杂的引擎结构
```go
type EngineInfo struct {
    Name        string `json:"name"`
    Path        string `json:"path"`
    DataDir     string `json:"data_dir"`
    IsInstalled bool   `json:"is_installed"`
}
Engines map[string]EngineInfo `json:"engines"`
PrimaryEngine string `json:"primary_engine"`
```

**建议**: 先用选项 A，简单且够用

### Q2: 更新/部署时如何处理多引擎？
**选项 A**: 只更新/部署到主引擎
**选项 B**: 更新/部署到所有已安装引擎
**选项 C**: 让用户选择要部署到哪些引擎

**建议**: 默认选项 A（只更新主引擎），UI 提供选项让用户选择是否同步到其他引擎

### Q3: 引擎检测的时机？
**选项 A**: 每次启动时检测
**选项 B**: 只在配置向导时检测
**选项 C**: 启动时检测 + 提供手动重新检测按钮

**建议**: 选项 C，兼顾性能和灵活性

---

## 参考资料

1. **FCITX5 macOS**: https://github.com/fcitx-contrib/fcitx5-macos-installer
2. **Linux deployer 多路径检测实现**: `internal/deployer/linux.go:58-83`
3. **macOS FCITX5 部署命令**: `/Library/Input Methods/Fcitx5.app/Contents/bin/fcitx5-curl /config/addon/rime/deploy -X POST -d {}`

---

## 下一步调查

- [ ] ⚠️ **关键确认**: 小企鹅安装器 (`Fcitx5Installer.app`) 运行后，最终安装到 `/Library/Input Methods/` 的应用名称是什么？
  - 安装器下载文件：`Fcitx5-Pinyin.zip` / `Fcitx5-Rime.zip` / `Fcitx5Installer.zip`
  - 安装器程序：`Fcitx5Installer.app`（用户打开的）
  - 最终输入法：`?` （需要确认）
- [ ] 验证小企鹅的 fcitx5-curl 部署命令的完整路径
- [ ] 确认 Windows 是否可能有多个 Rime 引擎并存
- [x] ~~测试 FCITX5 macOS 的实际安装路径和部署命令~~ **部分确认，等待最终应用名称**
- [ ] 检查是否有其他平台特定的 Rime 实现
