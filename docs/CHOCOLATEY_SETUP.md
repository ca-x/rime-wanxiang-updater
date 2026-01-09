# Chocolatey 发布配置指南

本文档说明如何配置 GitHub Actions 自动发布到 Chocolatey.org。

## 前置条件

1. **Chocolatey 账户**: 在 [Chocolatey.org](https://community.chocolatey.org/) 注册账户
2. **API Key**: 从您的 Chocolatey 账户获取 API Key

## 获取 Chocolatey API Key

1. 登录 [Chocolatey.org](https://community.chocolatey.org/)
2. 点击右上角的用户名，进入 "Account" 页面
3. 找到 "API Key" 部分
4. 复制您的 API Key

## 配置 GitHub Secrets

在您的 GitHub 仓库中添加以下 Secret：

1. 进入仓库的 Settings > Secrets and variables > Actions
2. 点击 "New repository secret"
3. 添加以下 Secret：

   - **Name**: `CHOCOLATEY_API_KEY`
   - **Value**: 您的 Chocolatey API Key

## Workflow 说明

当您推送新的 tag (如 `v1.0.0`) 时，GitHub Actions 会自动：

1. 构建 Windows 平台的可执行文件（AMD64 和 ARM64）
2. 计算文件的 SHA256 校验和
3. 更新 Chocolatey 包配置文件中的版本号和校验和
4. 构建 Chocolatey 包（.nupkg 文件）
5. 发布到 Chocolatey.org

## 包文件结构

```
chocolatey/
├── rime-wanxiang-updater.nuspec     # 包元数据
└── tools/
    ├── chocolateyInstall.ps1        # 安装脚本
    └── chocolateyUninstall.ps1      # 卸载脚本
```

## 安装脚本说明

### chocolateyInstall.ps1

安装脚本会：
1. 自动检测系统架构（AMD64 或 ARM64）
2. 从 GitHub Releases 下载对应的可执行文件
3. 验证 SHA256 校验和
4. 创建 shim，使 `rime-wanxiang-updater` 命令在 PATH 中可用

### chocolateyUninstall.ps1

卸载脚本会：
1. 移除创建的 shim
2. 清理相关文件

## 用户安装方法

### 通过 Chocolatey 安装

```powershell
choco install rime-wanxiang-updater
```

### 升级

```powershell
choco upgrade rime-wanxiang-updater
```

### 卸载

```powershell
choco uninstall rime-wanxiang-updater
```

## 发布流程

1. **创建新的 Git Tag**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions 自动执行**:
   - 构建多平台二进制文件
   - 创建 GitHub Release
   - 发布到 AUR (Arch Linux)
   - 发布到 Chocolatey (Windows)

3. **验证发布**:
   - 检查 [GitHub Releases](https://github.com/czyt/rime-wanxiang-updater/releases)
   - 检查 [Chocolatey Package](https://community.chocolatey.org/packages/rime-wanxiang-updater)
   - GitHub Actions 日志

## 注意事项

1. **首次发布**: 首次发布到 Chocolatey 时，包需要经过审核，通常需要几个工作日
2. **后续更新**: 审核通过后，后续版本通常会自动审核通过
3. **版本号**: 确保 Git tag 使用语义化版本 (如 `v1.0.0`)
4. **校验和**: 工作流会自动计算并更新校验和，无需手动操作
5. **错误处理**: 如果发布失败，检查 GitHub Actions 日志和 Chocolatey API Key

## 包元数据更新

如果需要更新包的元数据（如描述、标签等），编辑以下文件：

- `chocolatey/rime-wanxiang-updater.nuspec`

然后推送新的 tag 触发发布。

## 常见问题

### Q: 如何测试 Chocolatey 包？

A: 可以在本地构建和测试：

```powershell
# 进入 chocolatey 目录
cd chocolatey

# 手动更新版本号和校验和
# 编辑 rime-wanxiang-updater.nuspec
# 编辑 tools/chocolateyInstall.ps1

# 构建包
choco pack

# 本地安装测试
choco install rime-wanxiang-updater -s . -y
```

### Q: 发布失败怎么办？

A: 检查以下几点：
1. `CHOCOLATEY_API_KEY` Secret 是否正确设置
2. API Key 是否有效
3. 包名是否已被占用（首次发布时）
4. 查看 GitHub Actions 日志获取详细错误信息

### Q: 如何撤回已发布的版本？

A: 登录 Chocolatey.org，在包管理页面可以撤回（unlist）特定版本。注意：撤回不是删除，包仍然可以通过直接版本号安装。

## 参考资料

- [Chocolatey 官方文档](https://docs.chocolatey.org/)
- [Chocolatey 包创建指南](https://docs.chocolatey.org/en-us/create/create-packages)
- [GitHub Actions Chocolatey Action](https://github.com/crazy-max/ghaction-chocolatey)
