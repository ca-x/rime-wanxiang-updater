# Chocolatey 自动发布功能添加完成 ✅

## 概述

已成功为 rime-wanxiang-updater 项目添加 Chocolatey.org 自动发布功能。当您推送新的版本 tag 时，GitHub Actions 会自动构建并发布 Chocolatey 包。

## 添加的文件

### 1. Chocolatey 包配置文件

```
chocolatey/
├── rime-wanxiang-updater.nuspec           # 包元数据和描述
└── tools/
    ├── chocolateyInstall.ps1              # 安装脚本
    └── chocolateyUninstall.ps1            # 卸载脚本
```

#### rime-wanxiang-updater.nuspec
- 定义包的元数据（名称、版本、作者、描述等）
- 包含详细的中文描述和功能说明
- 支持自动版本替换（`$VERSION$` 占位符）

#### chocolateyInstall.ps1
- 自动检测系统架构（AMD64/ARM64）
- 从 GitHub Releases 下载对应的可执行文件
- SHA256 校验和验证
- 创建命令行 shim，使 `rime-wanxiang-updater` 命令全局可用

#### chocolateyUninstall.ps1
- 清理安装时创建的 shim
- 完整卸载流程

### 2. GitHub Actions Workflow

在 `.github/workflows/release.yml` 中添加了新的 job：

```yaml
publish-chocolatey:
  name: Publish to Chocolatey
  needs: release
  runs-on: windows-latest
```

**功能说明**：
1. 下载 Windows 平台的构建产物（AMD64 和 ARM64）
2. 计算 SHA256 校验和
3. 自动更新包配置文件中的版本号和校验和
4. 使用 `crazy-max/ghaction-chocolatey@v3.4.0` 构建包
5. 发布到 Chocolatey.org

### 3. 文档

- **docs/CHOCOLATEY_SETUP.md**: 完整的配置和使用指南
  - 如何获取 Chocolatey API Key
  - 如何配置 GitHub Secrets
  - 用户安装方法
  - 发布流程说明
  - 常见问题解答

### 4. .gitignore

添加了忽略规则：
- `*.nupkg` - Chocolatey 包文件
- `*.nupkg.metadata` - 包元数据
- 其他构建产物和临时文件

## 配置步骤

### 1. 获取 Chocolatey API Key

1. 访问 [Chocolatey.org](https://community.chocolatey.org/) 并登录/注册
2. 进入 Account 页面
3. 找到并复制您的 API Key

### 2. 配置 GitHub Secret

在 GitHub 仓库中添加 Secret：

```
Settings > Secrets and variables > Actions > New repository secret
```

- **Name**: `CHOCOLATEY_API_KEY`
- **Value**: 您的 Chocolatey API Key

### 3. 推送 Tag 触发发布

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## 自动发布流程

当推送新的版本 tag 时：

1. **Build Job**: 构建所有平台的二进制文件
2. **Release Job**: 创建 GitHub Release
3. **publish-aur Job**: 发布到 AUR (Arch Linux)
4. **publish-chocolatey Job** (新增):
   - 下载 Windows 二进制文件
   - 计算 SHA256 校验和
   - 更新包配置文件
   - 构建 Chocolatey 包
   - 发布到 Chocolatey.org

## 用户安装方法

### Windows 用户通过 Chocolatey 安装

```powershell
# 安装
choco install rime-wanxiang-updater

# 升级
choco upgrade rime-wanxiang-updater

# 卸载
choco uninstall rime-wanxiang-updater
```

安装后，直接在命令行运行：

```powershell
rime-wanxiang-updater
```

## 技术特性

### 架构支持
- Windows AMD64 (x86-64)
- Windows ARM64 (Qualcomm Snapdragon PC)

### 安全性
- SHA256 校验和验证
- 从官方 GitHub Releases 下载
- 签名验证（Chocolatey 平台）

### 自动化
- 版本号自动替换
- 校验和自动计算和更新
- 无需手动编辑配置文件

## 注意事项

### 首次发布
- 首次发布到 Chocolatey 需要人工审核
- 审核时间通常为 2-5 个工作日
- 审核通过后，后续版本会自动审核

### 版本号格式
- 使用语义化版本：`v1.0.0`, `v2.1.3`
- 不要使用 `-rc`, `-beta` 等后缀（Chocolatey 不支持）

### 错误处理
- Workflow 中使用 `continue-on-error: true`
- 即使 Chocolatey 发布失败，不会影响其他发布流程
- 可以手动重试或在 Chocolatey.org 手动上传

## 包信息

### Chocolatey 包页面
发布后可在此查看：
```
https://community.chocolatey.org/packages/rime-wanxiang-updater
```

### 包内容
- 包名: `rime-wanxiang-updater`
- 作者: `czyt`
- 标签: `rime`, `wanxiang`, `input-method`, `updater`, `chinese`, `tui`, `cli`
- 许可证: 项目 LICENSE 文件
- 项目主页: GitHub 仓库

## 本地测试

如需在本地测试 Chocolatey 包：

```powershell
# 进入 chocolatey 目录
cd chocolatey

# 手动设置版本号（用于测试）
# 编辑 rime-wanxiang-updater.nuspec
# 编辑 tools/chocolateyInstall.ps1

# 构建包
choco pack

# 本地安装测试
choco install rime-wanxiang-updater -s . -y

# 测试卸载
choco uninstall rime-wanxiang-updater -y
```

## 维护建议

1. **定期更新**: 推送新版本时检查 GitHub Actions 日志
2. **监控审核**: 首次发布后关注 Chocolatey 审核状态
3. **用户反馈**: 在 Chocolatey 包页面回复用户评论
4. **更新文档**: 如有新功能，更新 nuspec 中的描述

## 状态检查

### GitHub Actions
```
https://github.com/czyt/rime-wanxiang-updater/actions
```

### Chocolatey 包统计
发布后可查看：
- 下载次数
- 版本历史
- 用户评分和评论

## 总结

✅ **完成的工作**:
- Chocolatey 包配置文件创建完成
- GitHub Actions workflow 已更新
- 完整的文档和指南
- .gitignore 规则已添加

🔧 **需要配置**:
- 在 GitHub 仓库添加 `CHOCOLATEY_API_KEY` Secret

🚀 **准备就绪**:
- 推送新的版本 tag 即可触发自动发布
- Windows 用户可通过 Chocolatey 安装您的工具

---

**创建日期**: 2026-01-10
**状态**: ✅ 配置完成，等待 API Key 设置
