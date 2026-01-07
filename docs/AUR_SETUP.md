# AUR 发布配置指南

本项目支持自动发布到 Arch User Repository (AUR)。

## 前置要求

1. 拥有 AUR 账号
2. 在 AUR 上注册包名 `rime-wanxiang-updater`

## 配置步骤

### 1. 生成 SSH 密钥对

```bash
ssh-keygen -t ed25519 -C "your_email@example.com" -f ~/.ssh/aur
```

### 2. 添加公钥到 AUR

1. 登录 [AUR 官网](https://aur.archlinux.org/)
2. 进入 Account -> My Account
3. 在 "SSH Public Key" 部分粘贴 `~/.ssh/aur.pub` 的内容
4. 点击 "Update"

### 3. 配置 GitHub Secrets

在 GitHub 仓库的 Settings -> Secrets and variables -> Actions 中添加以下 secrets：

| Secret 名称 | 说明 | 示例 |
|------------|------|------|
| `AUR_USERNAME` | AUR 用户名 | `yourname` |
| `AUR_EMAIL` | AUR 注册邮箱 | `your_email@example.com` |
| `AUR_SSH_PRIVATE_KEY` | SSH 私钥内容 | `~/.ssh/aur` 的完整内容 |

#### 获取 SSH 私钥内容

```bash
cat ~/.ssh/aur
```

将完整输出（包括 `-----BEGIN OPENSSH PRIVATE KEY-----` 和 `-----END OPENSSH PRIVATE KEY-----`）复制到 `AUR_SSH_PRIVATE_KEY`。

### 4. 首次发布到 AUR

首次发布需要手动操作：

```bash
# 克隆 AUR 仓库
git clone ssh://aur@aur.archlinux.org/rime-wanxiang-updater.git aur-repo
cd aur-repo

# 复制 PKGBUILD
cp ../PKGBUILD .

# 生成 .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# 提交并推送
git add PKGBUILD .SRCINFO
git commit -m "Initial commit"
git push origin master
```

### 5. 自动发布

配置完成后，每次推送新 tag（如 `v0.5.7`）时，GitHub Actions 会自动：

1. 构建多平台二进制文件
2. 创建 GitHub Release
3. 更新 PKGBUILD 版本号
4. 推送到 AUR 仓库

## 工作流程

```mermaid
graph LR
    A[推送 Tag] --> B[构建二进制]
    B --> C[创建 Release]
    C --> D[更新 PKGBUILD]
    D --> E[推送到 AUR]
```

## 验证发布

发布成功后，可以在以下位置验证：

- GitHub Releases: https://github.com/ca-x/rime-wanxiang-updater/releases
- AUR 包页面: https://aur.archlinux.org/packages/rime-wanxiang-updater

## 用户安装

Arch Linux 用户可以通过以下方式安装：

```bash
# 使用 yay
yay -S rime-wanxiang-updater

# 使用 paru
paru -S rime-wanxiang-updater

# 手动安装
git clone https://aur.archlinux.org/rime-wanxiang-updater.git
cd rime-wanxiang-updater
makepkg -si
```

## 故障排查

### 1. SSH 认证失败

确认：
- SSH 公钥已正确添加到 AUR 账号
- GitHub Secret `AUR_SSH_PRIVATE_KEY` 包含完整的私钥内容（包括头尾标识）

### 2. 包名冲突

确认：
- 包名 `rime-wanxiang-updater` 在 AUR 上已注册
- 如果是首次发布，需要先手动完成第 4 步

### 3. 版本号问题

- 确认 PKGBUILD 中的 `pkgver` 格式正确（不带 'v' 前缀）
- GitHub Actions 会自动去除 tag 中的 'v' 前缀

## 参考资料

- [AUR 提交指南](https://wiki.archlinux.org/title/AUR_submission_guidelines)
- [PKGBUILD 文档](https://wiki.archlinux.org/title/PKGBUILD)
- [github-actions-deploy-aur](https://github.com/KSXGitHub/github-actions-deploy-aur)
