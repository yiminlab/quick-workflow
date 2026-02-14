# PR 批准命令指南

使用新的 `qkflow pr approve` 命令的快速指南。

> 📚 **寻找完整的 PR 工作流？** 请参阅 [PR Workflow Guide](pr-workflow.md) 了解从创建到合并的完整生命周期。

## 🚀 快速开始

### 基础批准

```bash
# 通过 PR 编号批准特定 PR（使用默认 👍 评论）
qkflow pr approve 123

# 通过 URL 批准 PR（可在任何地方使用！）
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# 也支持 /files、/commits、/checks URLs
qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

# 从当前分支自动检测
qkflow pr approve
```

### 带评论

默认情况下，所有批准都使用 👍 作为评论。使用 `-c` 标志自定义：

```bash
# 默认批准（带 👍）
qkflow pr approve 123

# 自定义评论
qkflow pr approve 123 -c "LGTM! 🎉"

# 通过 URL 批准并添加自定义评论
qkflow pr approve https://github.com/owner/repo/pull/456 -c "Great work!"

# 使用标志的长评论
qkflow pr approve 123 --comment "Great work! All tests passed. Approved for merge."
```

### 自动合并

```bash
# 一步完成批准和合并
qkflow pr approve 123 --merge

# 通过 URL 批准并合并
qkflow pr approve https://github.com/owner/repo/pull/789 --merge

# 短标志
qkflow pr approve 123 -m

# 带评论和合并
qkflow pr approve 123 -c "LGTM!" -m

# URL 带评论和合并
qkflow pr approve https://github.com/owner/repo/pull/789 -c "Ship it! 🚀" -m
```

## 🌐 URL 支持（新功能！）

现在你可以从**任何仓库**批准 PR，而无需在 git 目录中！

### 为什么使用 URL？

1. **跨仓库**: 批准不同仓库的 PR
2. **无需切换目录**: 可在任何地方工作
3. **浏览器到 CLI**: 直接从 GitHub 复制 URL
4. **批量操作**: 跨多个仓库编写脚本批准

### URL 示例

```bash
# 批准不同仓库的 PR
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# 同事分享 PR 链接，立即批准
qkflow pr approve https://github.com/company/frontend/pull/456 -c "Looks good!"

# 通过 URL 合并别人的 PR
qkflow pr merge https://github.com/team/backend/pull/789

# 通过 URL 批准并合并
qkflow pr approve https://github.com/org/project/pull/123 -c "LGTM! 🎉" -m
```

### 支持的 URL 格式

所有这些格式都可以使用：

```bash
# HTTPS（最常见）
https://github.com/owner/repo/pull/123

# 带 /files 后缀（Files 标签页）
https://github.com/brain/planning-api/pull/2001/files

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
https://github.com/owner/repo/pull/123/files?file-filters%5B%5D=.js

# 带片段（正确解析）
https://github.com/owner/repo/pull/123#discussion_r123456
```

**专业提示：** 只需从任何 PR 标签页（Overview、Files、Commits、Checks）复制 URL，它就能工作！🎉

### URL 解析详情

工具会自动：
- 检测参数是 URL 还是编号
- 从 URL 中提取 owner、repo 和 PR 编号
- 处理查询参数和片段
- 验证 URL 格式

### 何时使用 URL vs 编号

| 场景 | 使用 | 示例 |
|------|------|------|
| 你的仓库，在 git 目录中 | 编号 | `qkflow pr approve 123` |
| 你的仓库，当前分支 | 无 | `qkflow pr approve` |
| 不同仓库 | URL | `qkflow pr approve https://...` |
| 从浏览器共享链接 | URL | 复制并粘贴 URL |
| 跨多个仓库编写脚本 | URL | 循环遍历 URLs |

## 📋 常见工作流

### 工作流 1：快速代码审查

你正在审查同事的 PR：

```bash
# 1. 检出他们的分支（可选）
git fetch origin
git checkout feature-branch

# 2. 审查代码...

# 3. 批准
qkflow pr approve
# 从分支自动查找 PR
# 添加可选评论
# 完成！
```

### 工作流 2：批准并合并

你有批准权限并想立即合并：

```bash
# 一个命令完成批准和合并
qkflow pr approve 123 -c "Approved! Merging now." --merge

# 执行过程：
# ✅ 批准 PR #123
# ✅ 添加评论
# ✅ 检查是否可合并
# ✅ 与你确认
# ✅ 合并 PR
# ✅ 删除远程分支
# ✅ 切换到 main
# ✅ 删除本地分支
```

### 工作流 3：批量批准

多个 PR 需要审查：

```bash
# 首先列出所有打开的 PR
gh pr list

# 逐个批准它们
qkflow pr approve 101 -c "Approved"
qkflow pr approve 102 -c "Approved"
qkflow pr approve 103 -c "Approved"
```

### 工作流 4：交互模式

不知道 PR 编号？让工具帮助你：

```bash
# 不带参数运行
qkflow pr approve

# 执行过程：
# 1. 尝试从当前分支查找 PR
# 2. 如果未找到，显示所有打开的 PR 列表
# 3. 你选择一个
# 4. 询问可选评论
# 5. 批准！
```

## 🔍 用例

### 用例 1：跨仓库审查

无需切换目录即可审查多个仓库的 PR：

```bash
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"
qkflow pr approve https://github.com/team/backend/pull/200 -c "Approved"
qkflow pr approve https://github.com/team/mobile/pull/300 -c "Approved"
```

### 用例 2：浏览器到 CLI 工作流

1. 在 GitHub Web 中打开 PR
2. 从地址栏复制 URL
3. 粘贴到终端：

```bash
qkflow pr approve https://github.com/company/project/pull/1234 -c "Looks good!"
```

### 用例 3：Slack/Email 集成

队友在 Slack 中分享 PR 链接？立即批准：

```bash
# 从 Slack 复制链接
qkflow pr approve <paste-url-here> -c "Reviewed and approved"
```

### 用例 4：跨仓库编写脚本

自动化跨多个仓库的批准：

```bash
#!/bin/bash
PR_URLS=(
  "https://github.com/org/repo1/pull/10"
  "https://github.com/org/repo2/pull/20"
  "https://github.com/org/repo3/pull/30"
)

for url in "${PR_URLS[@]}"; do
  qkflow pr approve "$url" -c "Auto-approved by bot"
done
```

### 用例 5：团队负责人批准

作为团队负责人，你需要每天批准 PR：

```bash
# 早晨例行：批准所有准备好的 PR
for pr in 121 122 123; do
  qkflow pr approve $pr -c "Reviewed and approved"
done
```

### 用例 6：CI/CD 集成

添加到你的 CI 流水线：

```bash
#!/bin/bash
# 测试通过后自动批准 dependabot PR
if [[ "$PR_AUTHOR" == "dependabot" ]] && [[ "$TESTS_PASSED" == "true" ]]; then
  qkflow pr approve $PR_NUMBER -c "Auto-approved: Tests passed" -m
fi
```

### 用例 7：热修复工作流

快速处理紧急修复：

```bash
# 创建热修复
git checkout -b hotfix/critical-bug
# ... 修复 bug ...
git add .
qkflow pr create

# 尽快批准并合并
qkflow pr approve 999 -c "Critical hotfix - merging immediately" -m
```

## 🎯 专业技巧

### 技巧 1：别名

添加到 `.bashrc` 或 `.zshrc`：

```bash
# 快速批准
alias approve='qkflow pr approve'
alias merge='qkflow pr approve --merge'
alias gha='qkflow pr approve'  # GitHub Approve 快捷方式

# 使用：
approve 123 -c "LGTM"
merge 123
gha https://github.com/owner/repo/pull/123 -c "LGTM!"
```

### 技巧 1b：与 `pbpaste` 一起使用（macOS）

```bash
# 在浏览器中复制 URL，然后：
qkflow pr approve "$(pbpaste)" -c "Approved!"
```

### 技巧 1c：与 GitHub CLI 集成

与 `gh` CLI 结合使用：

```bash
# 使用 gh 列出 PR，使用 qkflow 批准
gh pr list
qkflow pr approve https://github.com/owner/repo/pull/123 -m
```

### 技巧 2：评论模板

保存常用评论：

```bash
# 在 shell 配置中
export APPROVE_LGTM="Looks good to me! 👍"
export APPROVE_MINOR="Approved with minor comments. Please address in follow-up."
export APPROVE_EXCELLENT="Excellent work! 🎉"

# 使用：
qkflow pr approve 123 -c "$APPROVE_LGTM"
```

### 技巧 3：合并前检查

在使用 `--merge` 之前，验证 PR：

```bash
# 查看 PR 详情
gh pr view 123

# 检查 CI 状态
gh pr checks 123

# 如果一切正常，批准并合并
qkflow pr approve 123 -m
```

### 技巧 4：分支保护

如果启用了分支保护：

```bash
# 只批准 - 让 GitHub 合并规则处理其余部分
qkflow pr approve 123 -c "Approved"

# 如果需要多个批准，不要使用 --merge
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
ℹ️  Expected: PR number (e.g., '123') or GitHub URL (e.g., 'https://github.com/owner/repo/pull/123')
```

**解决方案：** 使用有效的 PR 编号或 GitHub URL

### PR 已关闭

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

### 分支没有 PR

```bash
$ qkflow pr approve
⚠️ No PR found for branch: feature-xyz
Do you want to select a PR from the list? (Y/n)
```

**解决方案：** 从列表中选择或指定 PR 编号

## 🆚 对比

### GitHub CLI (`gh`)

```bash
# 批准
gh pr review 123 --approve --body "LGTM"

# 然后单独合并
gh pr merge 123

# 然后清理
git checkout main
git pull
git branch -D feature-branch
```

### qkflow（新功能！）

```bash
# 一步完成！
qkflow pr approve 123 -c "LGTM" -m
```

**优势：**
- ✅ 更少的命令
- ✅ 自动清理
- ✅ 交互式提示
- ✅ 分支自动检测
- ✅ 合并验证

## 🔗 相关命令

- `qkflow pr create` - 创建新 PR
- `qkflow pr merge` - 不先批准就合并
- `gh pr view 123` - 查看 PR 详情
- `gh pr checks 123` - 检查 CI 状态

## 📚 完整参考

### 命令语法

```
qkflow pr approve [pr-number] [flags]
```

### 标志

| 标志 | 简写 | 描述 |
|------|------|------|
| `--comment` | `-c` | 添加评论 |
| `--merge` | `-m` | 批准后自动合并 |
| `--help` | `-h` | 显示帮助信息 |

### 退出代码

- `0`: 成功
- `1`: 错误（PR 未找到、无法合并等）

### 环境

需要：
- 配置 GitHub token（`qkflow config`）
- 带远程的 Git 仓库
- 对仓库的写权限

## 📝 实际示例

### 示例 1：快速审查（默认评论）

```bash
$ qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

ℹ️  Detected GitHub PR URL, parsing...
✅ Parsed: brain/planning-api PR #2001
ℹ️  Fetching PR #2001...
ℹ️  PR: feat: Add user authentication
ℹ️  Branch: feature/auth -> main
ℹ️  State: open
ℹ️  Using default comment: 👍 (use -c flag to customize)
ℹ️  Approving PR #2001...
✅ PR approved with comment: 👍

ℹ️  PR approved. Use 'qkflow pr merge' to merge it later, or run with --merge flag to auto-merge.
```

### 示例 2：快速审查（自定义评论）

```bash
$ qkflow pr approve https://github.com/brain/planning-api/pull/2001 -c "LGTM!"

ℹ️  Detected GitHub PR URL, parsing...
✅ Parsed: brain/planning-api PR #2001
ℹ️  Fetching PR #2001...
ℹ️  PR: feat: Add user authentication
ℹ️  Branch: feature/auth -> main
ℹ️  State: open
ℹ️  Approving PR #2001...
✅ PR approved with comment: LGTM!

ℹ️  PR approved. Use 'qkflow pr merge' to merge it later, or run with --merge flag to auto-merge.
```

### 示例 3：批准并合并

```bash
$ qkflow pr approve https://github.com/team/backend/pull/456 -c "Ship it! 🚀" -m

ℹ️  Detected GitHub PR URL, parsing...
✅ Parsed: team/backend PR #456
ℹ️  Fetching PR #456...
ℹ️  PR: fix: Database connection timeout
ℹ️  Branch: fix/db-timeout -> main
ℹ️  State: open
ℹ️  Approving PR #456...
✅ PR approved with comment: Ship it! 🚀

ℹ️  Checking if PR is mergeable...
❓ Proceed with merging the PR? (Y/n) y
ℹ️  Merging PR #456...
✅ 🎉 PR merged successfully!
ℹ️  Deleting remote branch fix/db-timeout...
✅ Remote branch deleted

✅ All done! 🎉
```

## 🤝 向后兼容性

所有现有工作流仍然有效：

```bash
# PR 编号（需要在仓库中）
qkflow pr approve 123

# 从当前分支自动检测
qkflow pr approve

# 交互式选择
qkflow pr approve
# → 显示要选择的 PR 列表
```

---

**需要帮助？**

```bash
qkflow pr approve --help
qkflow help
```

**发现 Bug？**

在 GitHub 上提交 issue 并提供详细信息！