# 任务计划 - 支持多输入法引擎架构

**目标**: 重构代码以支持同一系统上多个 Rime 输入法引擎并存，用户可选择主引擎和部署目标。

**复杂度**: 中等 - 涉及配置结构、检测逻辑、部署流程和 UI 的多层改动

**预计阶段数**: 7 个阶段

---

## Phase 1: 设计新的配置结构 ⏸️ pending
**目标**: 设计支持多引擎的配置结构

**任务**:
- [ ] 设计 `InstalledEngines` 数组和 `PrimaryEngine` 字段
- [ ] 决定是否需要完整的 `EngineInfo` 结构
- [ ] 设计向后兼容策略（旧配置迁移）
- [ ] 更新 `internal/types/types.go` 中的 Config 结构

**输出**:
- 更新的 `types.Config` 结构定义
- 配置迁移逻辑设计文档

**依赖**: 无

**风险**: 配置迁移可能影响现有用户

---

## Phase 2: 实现引擎自动检测 ⏸️ pending
**目标**: 为每个平台实现真实的引擎检测逻辑

**任务**:
- [ ] **macOS**: 检测 `/Library/Input Methods/` 下的 Squirrel.app 和 Fcitx5Installer.app
  - [ ] 🔧 修正 `darwin.go:36` 中的路径错误 (`Fcitx5.app` → `Fcitx5Installer.app`)
- [ ] **Linux**: 参考 `deployer/linux.go` 实现多路径检测
- [ ] **Windows**: 检测注册表和默认路径
- [ ] 实现 `detectInstalledEngines()` 函数返回引擎列表
- [ ] 更新 `config.go` 的 `detectEngine()` 或替换为新函数

**输出**:
- 跨平台的引擎检测函数
- 更新的 `paths_*.go` 文件
- 修正的 macOS 小企鹅路径

**依赖**: Phase 1 完成

**风险**: 不同环境下路径可能有变化

---

## Phase 3: 更新路径获取逻辑 ⏸️ pending
**目标**: 让路径获取支持多引擎

**任务**:
- [ ] 修改 `getRimeUserDir()` 支持引擎参数
- [ ] 实现 `getAllEngineDataDirs()` 返回所有引擎的数据目录
- [ ] 更新 `paths_darwin.go` 支持鼠须管和小企鹅
- [ ] 更新 `paths_linux.go` 支持多个引擎（参考 deployer）
- [ ] 评估 `paths_windows.go` 是否需要改动

**输出**:
- 支持多引擎的路径获取函数
- 单元测试

**依赖**: Phase 2 完成

**风险**: 路径逻辑改动可能影响现有功能

---

## Phase 4: 重构部署器 ⏸️ pending
**目标**: 支持部署到多个引擎

**任务**:
- [ ] 修改 `Deployer` 接口支持引擎选择
- [ ] 更新 `darwin.go` deployer 支持鼠须管和小企鹅
- [ ] 更新 `linux.go` deployer 支持多引擎部署
- [ ] 评估 `windows.go` 是否需要改动
- [ ] 实现"部署到所有引擎"和"部署到主引擎"两种模式

**输出**:
- 更新的 deployer 实现
- 部署到多引擎的逻辑

**依赖**: Phase 3 完成

**风险**: 部署失败处理需要更细致

---

## Phase 5: 更新配置管理逻辑 ⏸️ pending
**目标**: 配置加载/保存支持新结构

**任务**:
- [ ] 实现配置迁移逻辑（旧 `Engine` 字符串 → 新结构）
- [ ] 更新 `createDefaultConfig()` 调用引擎检测
- [ ] 更新 `loadOrCreateConfig()` 处理迁移
- [ ] 实现配置验证（主引擎必须在已安装列表中）
- [ ] 添加"重新检测引擎"功能

**输出**:
- 配置迁移代码
- 配置验证逻辑

**依赖**: Phase 1, 2 完成

**风险**: 迁移逻辑必须稳定可靠

---

## Phase 6: UI 层适配 ⏸️ pending
**目标**: UI 支持显示和选择引擎

**任务**:
- [ ] 配置向导显示检测到的所有引擎
- [ ] 让用户选择主引擎
- [ ] 设置界面显示已安装引擎列表
- [ ] **引擎显示格式**: 多引擎时用 `+` 连接（例如：`鼠须管+小企鹅`）
- [ ] 更新时提供"仅主引擎"/"所有引擎"选项
- [ ] 添加"重新检测引擎"按钮
- [ ] 引擎未检测到时的友好提示

**输出**:
- 更新的 UI 视图和交互
- 引擎选择流程
- 多引擎显示格式

**依赖**: Phase 2, 5 完成

**风险**: UI 复杂度增加

---

## Phase 7: 测试和文档 ⏸️ pending
**目标**: 完整测试和更新文档

**任务**:
- [ ] 单元测试：引擎检测、配置迁移、路径获取
- [ ] 集成测试：多引擎部署流程
- [ ] 手动测试：macOS (鼠须管+小企鹅)
- [ ] 手动测试：Linux (fcitx5/ibus 等)
- [ ] 手动测试：Windows
- [ ] 更新 README.md 说明多引擎支持
- [ ] 更新配置文件文档

**输出**:
- 测试覆盖率报告
- 更新的文档

**依赖**: Phase 1-6 完成

**风险**: 不同环境测试成本高

---

## 错误记录

| 错误 | 尝试次数 | 解决方案 |
|------|---------|---------|
| - | - | - |

---

## 关键决策

| 决策点 | 选择 | 理由 | 日期 |
|--------|------|------|------|
| 配置结构 | 待定 | 需要在简单性和扩展性间权衡 | - |
| 部署模式 | 待定 | 默认只部署主引擎 vs 部署所有 | - |
| 检测时机 | 待定 | 启动检测 + 手动重新检测 | - |

---

## 文件修改清单

### 核心修改（必须）
- [ ] `internal/types/types.go` - Config 结构
- [ ] `internal/config/config.go` - detectEngine(), createDefaultConfig()
- [ ] `internal/config/paths_darwin.go` - getRimeUserDir()
- [ ] `internal/config/paths_linux.go` - getRimeUserDir()
- [ ] `internal/deployer/darwin.go` - 多引擎部署
- [ ] `internal/ui/model.go` - 引擎选择状态
- [ ] `internal/ui/views.go` - 引擎显示
- [ ] `internal/ui/handlers.go` - 引擎选择逻辑

### 可选修改
- [ ] `internal/config/paths_windows.go` - 如果 Windows 需要
- [ ] `internal/deployer/deployer.go` - 接口调整
- [ ] `internal/deployer/linux.go` - 优化多引擎支持
- [ ] `internal/deployer/windows.go` - 如果 Windows 需要

---

## 注意事项

1. **向后兼容**: 旧配置文件必须能无缝迁移
2. **用户体验**: 单引擎用户不应感到额外复杂
3. **错误处理**: 引擎未检测到时的降级策略
4. **性能**: 引擎检测不应显著增加启动时间
5. **测试覆盖**: 必须在真实多引擎环境测试

---

## 当前状态

**活跃阶段**: 无（规划中）
**已完成阶段**: 0/7
**总体进度**: 0%

**下一步行动**: 完成 Phase 1 的配置结构设计
