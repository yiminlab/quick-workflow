# PR 工作流指南

使用 `qkflow` 管理 Pull Request 的完整指南 - 从创建到批准到合并。

## 📋 目录

- [概述](#概述)
- [完整 PR 生命周期](#完整-pr-生命周期)
- [1. 创建 PR](#1-创建-pr)
- [2. 批准 PR](#2-批准-pr)
- [3. 合并 PR](#3-合并-pr)
- [常见工作流](#常见工作流)
- [URL 支持](#url-支持)
- [相关文档](#相关文档)

## 🎯 概述

`qkflow` 提供三个主要的 PR 命令，覆盖整个 Pull Request 生命周期：

| 命令 | 用途 | 主要特性 |
|------|------|----------|
| `qkflow pr create` | 创建新 PR | Web 编辑器、Jira 集成、自动分支创建 |
| `qkflow pr approve` | 批准 PR | URL 支持、自定义评论、自动合并 |
| `qkflow pr merge` | 合并 PR | 自动清理、Jira 状态更新 |

所有命令都支持：
- ✅ **PR 编号**（例如：`123`）
- ✅ **GitHub URLs**（例如：`https://github.com/owner/repo/pull/123`）
- ✅ **自动检测**当前分支

## 🔄 完整 PR 生命周期

这是从开始到结束的典型流程：

```
┌─────────────────────────────────────────────────────────────┐
│  1. 创建 PR                                                 │
│  qkflow pr create PROJ-123                                 │
│  ├─ 创建分支                                                │
│  ├─ 提交更改                                                │
│  ├─ 推送到远程                                              │
│  ├─ [可选] 使用 Web 编辑器添加描述                          │
│  ├─ 创建 GitHub PR                                          │
│  ├─ 上传文件到 GitHub & Jira                               │
│  └─ 更新 Jira 状态                                         │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  2. 代码审查                                                │
│  (手动审查过程)                                              │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  3. 批准 PR                                                 │
│  qkflow pr approve 123 -c "LGTM!"                           │
│  ├─ 在 GitHub 上批准                                        │
│  ├─ 添加评论                                                │
│  └─ [可选] 使用 --merge 标志自动合并                        │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  4. 合并 PR                                                 │
│  qkflow pr merge 123                                        │
│  ├─ 在 GitHub 上合并                                        │
│  ├─ 删除远程分支                                            │
│  ├─ 切换到主分支                                            │
│  ├─ 删除本地分支                                            │
│  └─ 更新 Jira 状态为 Done/Merged                           │
└─────────────────────────────────────────────────────────────┘
```

## 1. 创建 PR

### 快速开始

```bash
# 基础：使用 Jira ticket 创建 PR
qkflow pr create PROJ-123

# 不使用 Jira ticket（提示时按 Enter）
qkflow pr create
```

### 执行过程

1. ✅ 获取 Jira ticket 详情（如果提供）
2. ✅ 提示选择变更类型（feat、fix 等）
3. ✅ **[可选]** 打开 Web 编辑器添加描述和截图
4. ✅ 生成 PR 标题
5. ✅ 创建 git 分支
6. ✅ 提交暂存的更改
7. ✅ 推送到远程
8. ✅ 创建 GitHub PR
9. ✅ 上传文件并添加评论到 GitHub & Jira
10. ✅ 更新 Jira 状态
11. ✅ 复制 PR URL 到剪贴板

### Web 编辑器功能

创建 PR 时，会提示：

```
? Add detailed description with images/videos?
  > ⏭️  Skip (default)
    ✅ Yes, continue
```

如果选择 "Yes, continue"：
- 在浏览器中打开基于 Web 的 Markdown 编辑器
- 拖放图片/视频
- 从剪贴板粘贴图片
- 实时预览
- 自动上传到 GitHub 和 Jira

**示例：**

```bash
$ qkflow pr create NA-9245

✓ Found Jira issue: Fix login button styling
📝 Select type(s) of changes:
  ✓ 🐛 Bug fix

? Add detailed description with images/videos?
  > ⏭️  Skip (default)
    ✅ Yes, continue

# 选择 "Yes, continue"
🌐 Opening editor in your browser: http://localhost:54321
📝 Please edit your content in the browser and click 'Save and Continue'

✅ Content saved! (245 characters, 2 files)
✅ Generated title: fix: Update login button hover state
✅ Creating branch: NA-9245--fix-update-login-button-hover-state
...
✅ Pull request created: https://github.com/owner/repo/pull/123
📤 Uploading 2 file(s)...
✅ Description added to GitHub PR
✅ Description added to Jira
✅ All done! 🎉
```

### 支持的文件类型

- **图片**: PNG, JPG, JPEG, GIF, WebP, SVG
- **视频**: MP4, MOV, WebM, AVI

### 详细文档

有关 PR 编辑器功能的完整详情，请参阅：
- **[PR Editor Feature](pr-editor.md)** - 技术实现和高级用法

## 2. 批准 PR

### 快速开始

```bash
# 通过 PR 编号批准（默认 👍 评论）
qkflow pr approve 123

# 通过 URL 批准（可在任何地方使用！）
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# 自定义评论
qkflow pr approve 123 -c "LGTM! 🎉"

# 批准并自动合并
qkflow pr approve 123 -c "Ship it! 🚀" -m

# 从当前分支自动检测
qkflow pr approve
```

### 执行过程

1. ✅ 检测 PR 编号或 URL
2. ✅ 获取 PR 详情
3. ✅ 在 GitHub 上批准 PR
4. ✅ 添加评论（默认：👍，或使用 `-c` 自定义）
5. ✅ **[可选]** 如果使用 `--merge` 标志则自动合并

### 标志

| 标志 | 简写 | 描述 |
|------|------|------|
| `--comment` | `-c` | 添加自定义评论 |
| `--merge` | `-m` | 批准后自动合并 |
| `--help` | `-h` | 显示帮助信息 |

### 示例

**基础批准：**
```bash
$ qkflow pr approve 123
✅ PR approved with comment: 👍
```

**自定义评论：**
```bash
$ qkflow pr approve 123 -c "Great work! All tests passed."
✅ PR approved with comment: Great work! All tests passed.
```

**批准并合并：**
```bash
$ qkflow pr approve 123 -c "Ship it! 🚀" -m
✅ PR approved with comment: Ship it! 🚀
✅ Checking if PR is mergeable...
✅ PR merged successfully!
✅ Remote branch deleted
✅ Switched to default branch
✅ Local branch deleted
✅ All done! 🎉
```

### 详细文档

有关批准功能的完整详情，请参阅：
- **[PR Approve Guide](pr-approve.md)** - 包含 URL 支持、工作流和用例的完整指南

## 3. 合并 PR

### 快速开始

```bash
# 通过 PR 编号合并
qkflow pr merge 123

# 通过 URL 合并（可在任何地方使用！）
qkflow pr merge https://github.com/brain/planning-api/pull/2001

# 从当前分支自动检测
qkflow pr merge
```

### 执行过程

1. ✅ 检测 PR 编号或 URL
2. ✅ 获取 PR 详情
3. ✅ 与你确认合并
4. ✅ 在 GitHub 上合并 PR
5. ✅ 删除远程分支
6. ✅ 切换到主分支
7. ✅ 拉取最新更改
8. ✅ 删除本地分支
9. ✅ 更新 Jira 状态为 Done/Merged
10. ✅ 添加合并评论到 Jira

### 示例

```bash
$ qkflow pr merge 123

ℹ️  Fetching PR #123...
ℹ️  PR: feat: Add user authentication
ℹ️  Branch: feature/auth -> main
ℹ️  State: open
❓ Proceed with merging the PR? (Y/n) y
ℹ️  Merging PR #123...
✅ Pull request merged!
ℹ️  Deleting remote branch feature/auth...
✅ Remote branch deleted
ℹ️  Switching to main branch...
ℹ️  Pulling latest changes from main...
✅ Updated to latest changes
ℹ️  Deleting local branch feature/auth...
✅ Local branch deleted
ℹ️  Found Jira ticket: PROJ-123
ℹ️  Updating Jira status to: Done
✅ Updated Jira status to: Done
✅ All done! 🎉
```

## 🌐 URL 支持

所有三个 PR 命令都支持 GitHub URLs，方便跨仓库工作：

### 支持的 URL 格式

```bash
# HTTPS（最常见）
https://github.com/owner/repo/pull/123

# 带 /files 后缀（Files 标签页）
https://github.com/owner/repo/pull/123/files

# 带 /commits 后缀（Commits 标签页）
https://github.com/owner/repo/pull/123/commits

# 带 /checks 后缀（Checks 标签页）
https://github.com/owner/repo/pull/123/checks

# HTTP
http://github.com/owner/repo/pull/123

# 无协议
github.com/owner/repo/pull/123

# 带查询参数（正确解析）
https://github.com/owner/repo/pull/123?comments=all

# 带片段（正确解析）
https://github.com/owner/repo/pull/123#discussion_r123456
```

### 为什么使用 URL？

1. **跨仓库**: 处理不同仓库的 PR
2. **无需切换目录**: 可在任何地方工作
3. **浏览器到 CLI**: 直接从 GitHub 复制 URL
4. **批量操作**: 跨多个仓库编写脚本操作

### 示例

```bash
# 批准不同仓库的 PR
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"

# 合并共享链接的 PR
qkflow pr merge https://github.com/company/backend/pull/200

# 从浏览器复制 URL 并粘贴
qkflow pr approve "$(pbpaste)" -m  # macOS
```

## 📋 常见工作流

### 工作流 1：完整功能开发

```bash
# 1. 开始功能
git checkout -b feature/new-feature
# ... 进行更改 ...
git add .

# 2. 创建 PR
qkflow pr create PROJ-456
# → 选择变更类型
# → [可选] 添加带截图的描述
# → PR 已创建，Jira 已更新

# 3. 等待审查...

# 4. 审查者批准
qkflow pr approve https://github.com/org/repo/pull/789 -c "LGTM!"

# 5. 合并 PR
qkflow pr merge 789
# → PR 已合并，分支已清理，Jira 已更新为 Done
```

### 工作流 2：快速 Bug 修复

```bash
# 1. 创建热修复分支
git checkout -b hotfix/critical-bug
# ... 修复 bug ...
git add .

# 2. 创建 PR
qkflow pr create PROJ-789

# 3. 自我批准并立即合并
qkflow pr approve 999 -c "Critical hotfix - merging immediately" -m
# → 一步完成批准和合并！
```

### 工作流 3：跨仓库审查

```bash
# 无需切换目录即可审查多个仓库的 PR

# Frontend PR
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"

# Backend PR
qkflow pr approve https://github.com/team/backend/pull/200 -c "Looks good!"

# Mobile PR
qkflow pr approve https://github.com/team/mobile/pull/300 -c "LGTM!"
```

### 工作流 4：浏览器到 CLI

1. 在 GitHub Web 界面中打开 PR
2. 从地址栏复制 URL
3. 粘贴到终端：

```bash
# 批准
qkflow pr approve "$(pbpaste)" -c "Reviewed and approved"

# 或合并
qkflow pr merge "$(pbpaste)"
```

### 工作流 5：批量操作

```bash
#!/bin/bash
# 批准多个 PR
PR_URLS=(
  "https://github.com/org/repo1/pull/10"
  "https://github.com/org/repo2/pull/20"
  "https://github.com/org/repo3/pull/30"
)

for url in "${PR_URLS[@]}"; do
  qkflow pr approve "$url" -c "Auto-approved"
done
```

## 🎯 专业技巧

### 技巧 1：使用别名

添加到 `.bashrc` 或 `.zshrc`：

```bash
alias approve='qkflow pr approve'
alias merge='qkflow pr merge'
alias prc='qkflow pr create'
alias gha='qkflow pr approve'  # GitHub Approve 快捷方式
```

### 技巧 2：与 GitHub CLI 结合使用

```bash
# 使用 gh 列出 PR，使用 qkflow 批准
gh pr list
qkflow pr approve https://github.com/owner/repo/pull/123 -m
```

### 技巧 3：评论模板

```bash
# 在 shell 配置中
export APPROVE_LGTM="Looks good to me! 👍"
export APPROVE_EXCELLENT="Excellent work! 🎉"

# 使用：
qkflow pr approve 123 -c "$APPROVE_LGTM"
```

### 技巧 4：合并前检查

```bash
# 查看 PR 详情
gh pr view 123

# 检查 CI 状态
gh pr checks 123

# 如果一切正常，然后合并
qkflow pr merge 123
```

## ⚠️ 错误处理

### PR 未找到

```bash
$ qkflow pr approve 999
❌ Failed to get PR: Pull request not found
```

**解决方案：** 使用 `gh pr list` 检查 PR 编号

### 无效的 URL 或编号

```bash
$ qkflow pr approve invalid-url
❌ Invalid PR number or URL: invalid-url
ℹ️  Expected: PR number (e.g., '123') or GitHub URL
```

**解决方案：** 使用有效的 PR 编号或 GitHub URL

### PR 已关闭/合并

```bash
$ qkflow pr approve 123
❌ PR is not open (state: closed)
```

**解决方案：** PR 已经合并或关闭

### 无法合并

```bash
$ qkflow pr approve 123 -m
✅ PR approved!
⚠️ Cannot merge PR: PR has conflicts and cannot be merged
```

**解决方案：** 先解决冲突，然后合并

## 🔗 相关命令

- `qkflow pr create` - 创建新 PR
- `qkflow pr approve` - 批准 PR
- `qkflow pr merge` - 合并 PR
- `gh pr view 123` - 查看 PR 详情（GitHub CLI）
- `gh pr checks 123` - 检查 CI 状态（GitHub CLI）
- `gh pr list` - 列出所有 PR（GitHub CLI）

## 📚 相关文档

### 详细功能指南

- **[PR Approve Guide](pr-approve.md)** - `qkflow pr approve` 的完整指南
  - URL 支持详情
  - 高级工作流
  - 用例和示例
  - 错误处理

- **[PR Editor Feature](pr-editor.md)** - PR 创建编辑器的技术详情
  - Web 编辑器实现
  - 文件上传机制
  - 技术架构

### 其他相关文档

- **[Jira Integration](jira-integration.md)** - 完整的 Jira 集成指南（Issue Reader & 状态配置）
- **[Development Overview](../development/overview.md)** - 技术架构 ([English](../../en/development/overview.md))

## 🆚 与 GitHub CLI 对比

### GitHub CLI (`gh`)

```bash
# 创建 PR
gh pr create --title "Title" --body "Body"

# 批准
gh pr review 123 --approve --body "LGTM"

# 合并
gh pr merge 123

# 清理（手动）
git checkout main
git pull
git branch -D feature-branch
```

### qkflow

```bash
# 创建 PR（带 Jira 集成、Web 编辑器、自动分支创建）
qkflow pr create PROJ-123

# 批准（带 URL 支持、自动合并选项）
qkflow pr approve 123 -c "LGTM" -m

# 合并（带自动清理、Jira 更新）
qkflow pr merge 123
```

**优势：**
- ✅ 更少的命令
- ✅ 自动清理
- ✅ Jira 集成
- ✅ 交互式提示
- ✅ 分支自动检测
- ✅ URL 支持
- ✅ 基于 Web 的编辑器

---

**需要帮助？**

```bash
qkflow pr --help
qkflow pr create --help
qkflow pr approve --help
qkflow pr merge --help
```

**发现 Bug？**

在 GitHub 上提交 issue 并提供详细信息！