# 迁移指南：从 Shell 版本到 Go 版本

本指南将帮助你从基于 Shell 的 quick-workflow 迁移到新的 Go 版本。

## 🎯 为什么迁移？

| 方面 | Shell 版本 | Go 版本 | 改进 |
|------|-----------|---------|------|
| **安装** | 克隆仓库 + 安装 4+ 个依赖 | 下载 1 个二进制文件 | ✅ 简化 90% |
| **配置** | 在 `.zshrc` 中手动设置环境变量 | 交互式 `qkflow init` | ✅ 更容易 |
| **启动时间** | ~1-2 秒 | <100ms | ✅ 快 10-20 倍 |
| **跨平台** | 仅 macOS/Linux | macOS/Linux/Windows | ✅ 通用 |
| **更新** | `git pull` + 重新安装 | 下载新二进制文件 | ✅ 更简单 |
| **错误消息** | Shell 错误 | 清晰、可操作的消息 | ✅ 更好的用户体验 |
| **类型安全** | 运行时错误 | 编译时检查 | ✅ 更可靠 |

## 📋 前提条件

在迁移之前，请确保：
- ✅ 你当前的 Shell 版本正在工作
- ✅ 可以访问你的 Jira 和 GitHub 凭据
- ✅ 记录你当前的配置（特别是环境变量）

## 🔄 迁移步骤

### 步骤 1：备份当前配置

```bash
# 保存你当前的环境变量
cat ~/.zshrc | grep -E "(EMAIL|JIRA|GH_|OPENAI|DEEPSEEK)" > ~/qk-backup.txt
```

你的变量应该类似：
```bash
export EMAIL=your.email@example.com
export JIRA_API_TOKEN=your_jira_api_token
export JIRA_SERVICE_ADDRESS=https://your-domain.atlassian.net
export GH_BRANCH_PREFIX=your_branch_prefix
export GITHUB_TOKEN=your_github_token  # 或 gh auth token
```

### 步骤 2：安装 Go 版本

选择以下安装方法之一：

#### 选项 A：下载二进制文件（推荐）

```bash
# macOS (Apple Silicon)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# 验证安装
qkflow version
```

#### 选项 B：从源码构建

```bash
cd /path/to/quick-workflow/go-version
make build
sudo cp bin/qkflow /usr/local/bin/
```

### 步骤 3：运行初始设置

```bash
qkflow init
```

这将提示你输入：

1. **Email**: 使用你 `EMAIL` 环境变量的值
2. **GitHub Token**: 将从 `gh auth token` 自动检测
3. **Jira Service Address**: 使用 `JIRA_SERVICE_ADDRESS` 的值
4. **Jira API Token**: 使用 `JIRA_API_TOKEN` 的值
5. **Branch Prefix**（可选）: 使用 `GH_BRANCH_PREFIX` 的值
6. **OpenAI Key**（可选）: 使用 `OPENAI_KEY` 或 `DEEPSEEK_KEY` 的值

### 步骤 4：测试新版本

创建一个测试分支以验证一切正常工作：

```bash
# 测试 PR 创建
cd your-project
git checkout -b test-qkflow-migration
echo "test" > test.txt
git add test.txt
qkflow pr create

# 如果成功，你会看到：
# ✅ Branch created
# ✅ Changes committed
# ✅ Pushed to remote
# ✅ Pull request created: https://github.com/...
```

### 步骤 5：验证 Jira 集成

如果你在步骤 4 中使用了 Jira ticket：
1. 检查 PR 链接是否已添加到你的 Jira issue
2. 验证状态是否已更新（如果你选择更新它）

### 步骤 6：更新 Shell 别名（可选）

如果你在 Shell 版本中有自定义别名，请更新它们：

**旧别名（Shell 版本）：**
```bash
alias prc='~/quick-workflow/pr-create.sh'
alias prm='~/quick-workflow/pr-merge.sh'
```

**新别名（Go 版本）：**
```bash
alias prc='qkflow pr create'
alias prm='qkflow pr merge'
```

或者直接使用更短的命令：
```bash
qkflow pr create    # 替代 pr-create.sh
qkflow pr merge     # 替代 pr-merge.sh
```

### 步骤 7：清理旧安装（可选）

一旦你验证了 Go 版本可以正常工作：

```bash
# 从 .zshrc 中删除旧的环境变量
# （如果其他地方使用 EMAIL，请保留）
# 编辑 ~/.zshrc 并删除：
# export JIRA_API_TOKEN=...
# export JIRA_SERVICE_ADDRESS=...
# export GH_BRANCH_PREFIX=...
# export OPENAI_KEY=...
# 等等

# 重新加载 shell
source ~/.zshrc

# 归档旧安装
mv ~/quick-workflow ~/quick-workflow-shell-backup
```

## 🔍 功能对比

### 命令映射

| Shell 版本 | Go 版本 | 备注 |
|-----------|---------|------|
| `pr-create.sh` | `qkflow pr create` | 相同功能，更快 |
| `pr-merge.sh` | `qkflow pr merge` | 相同功能，更快 |
| `qk.sh` | 即将推出 | 日志管理功能 |
| `qklogs.sh` | 即将推出 | 将被集成 |
| `qkfind.sh` | 即将推出 | 将被集成 |
| N/A | `qkflow init` | 新功能：设置向导 |
| N/A | `qkflow config` | 新功能：显示配置 |
| N/A | `qkflow version` | 新功能：版本信息 |

### 配置

| Shell 版本 | Go 版本 | 位置 |
|-----------|---------|------|
| `.zshrc` 环境变量 | YAML 配置 | `~/.qkflow/config.yaml` 或 iCloud Drive |
| `jira-status.txt` | 自动检测 | 由 API 处理 |
| 手动设置 | `qkflow init` | 交互式向导 |

### 工作流对比

**创建 PR - Shell 版本：**
```bash
# 1. 暂存更改
git add .

# 2. 运行脚本
./pr-create.sh PROJ-123

# 3. 等待 ~1-2 秒启动
# 4. 回答提示
# 5. 脚本调用 gh、jira、git、jq、python
# 6. 完成（如果没有错误）
```

**创建 PR - Go 版本：**
```bash
# 1. 暂存更改
git add .

# 2. 运行命令
qkflow pr create PROJ-123

# 3. 即时启动（<100ms）
# 4. 回答提示（与之前相同）
# 5. 所有操作在一个快速二进制文件中
# 6. 完成，错误消息更好
```

## 🐛 故障排除

### 问题："Command not found: qkflow"

**解决方案：**
```bash
# 检查二进制文件是否在 PATH 中
which qkflow

# 如果未找到，确保 /usr/local/bin 在 PATH 中
echo $PATH

# 如果需要，添加到 PATH（添加到 ~/.zshrc）
export PATH="/usr/local/bin:$PATH"
```

### 问题："Config not found" 或 "Please run qkflow init"

**解决方案：**
```bash
# 运行设置向导
qkflow init

# 或手动创建配置
mkdir -p ~/.qkflow
cat > ~/.qkflow/config.yaml << EOF
email: your.email@example.com
jira_api_token: your_token
jira_service_address: https://your-domain.atlassian.net
github_token: your_github_token
branch_prefix: feature
EOF
```

### 问题："Failed to create GitHub client"

**解决方案：**
```bash
# 确保 gh CLI 已认证
gh auth status

# 如果未认证，登录
gh auth login

# 重新运行 qkflow init 以获取 token
qkflow init
```

### 问题："Failed to get Jira issue"

**解决方案：**
1. 验证 Jira API token 是否正确
2. 检查 Jira service address 格式：`https://your-domain.atlassian.net`
3. 确保你有权访问 Jira 项目
4. 测试 Jira 凭据：
   ```bash
   curl -u your.email@example.com:your_jira_token \
     https://your-domain.atlassian.net/rest/api/2/myself
   ```

### 问题：Go 版本使用的状态与 Shell 版本不同

**解决方案：**
Go 版本动态查询 Jira 的可用状态。如果你在 `jira-status.txt` 中有自定义状态映射，你需要在提示时选择它们。Go 版本会记住你的选择。

## 📊 性能对比

实际基准测试（在 M1 MacBook Pro 上）：

| 操作 | Shell 版本 | Go 版本 | 速度提升 |
|------|-----------|---------|---------|
| 启动时间 | ~1.5s | ~50ms | 快 30 倍 |
| PR 创建（总计） | ~8-10s | ~6-7s | 快 1.5 倍 |
| PR 合并（总计） | ~5-7s | ~4-5s | 快 1.4 倍 |
| 配置加载 | ~200ms | ~5ms | 快 40 倍 |

## 🎓 学习新命令

### 快速参考卡

打印此卡片并在迁移期间随身携带：

```
┌─────────────────────────────────────────────────┐
│         Quick Workflow Go Version               │
├─────────────────────────────────────────────────┤
│ Setup                                           │
│   qkflow init          Initialize configuration │
│   qkflow config        Show current config      │
│                                                 │
│ Pull Requests                                   │
│   qkflow pr create     Create PR + update Jira  │
│   qkflow pr merge      Merge PR + update Jira  │
│                                                 │
│ Help                                            │
│   qkflow --help        Show all commands        │
│   qkflow pr --help     Show PR commands        │
│   qkflow version       Show version info       │
└─────────────────────────────────────────────────┘
```

## 🎉 迁移清单

使用此清单确保顺利迁移：

- [ ] 已备份当前配置
- [ ] 已安装 Go 版本
- [ ] 已运行 `qkflow init` 并配置
- [ ] 已测试 PR 创建
- [ ] 已验证 Jira 集成工作正常
- [ ] 已更新 shell 别名（如果有）
- [ ] 已删除旧的环境变量
- [ ] 已归档旧 Shell 安装
- [ ] 已验证所有团队成员已迁移
- [ ] 已更新团队文档
- [ ] 庆祝更快的工作流！🎊

## 🆘 需要帮助？

如果在迁移过程中遇到问题：

1. **查看本指南** - 上面涵盖了最常见的问题
2. **检查配置** - 运行 `qkflow config` 验证设置
3. **启用调试模式** - 在运行命令之前设置 `export QK_DEBUG=1`
4. **创建 issue** - [在 GitHub 上创建 issue](https://github.com/Wangggym/quick-workflow/issues)
5. **寻求帮助** - 联系维护者

## 📝 回滚计划

如果你需要回滚到 Shell 版本：

```bash
# 恢复 shell 版本
mv ~/quick-workflow-shell-backup ~/quick-workflow

# 恢复环境变量
# （编辑 ~/.zshrc 并添加回去）

# 重新加载 shell
source ~/.zshrc

# 删除 Go 版本
sudo rm /usr/local/bin/qkflow
```

## 🚀 下一步

成功迁移后：

1. **探索新功能** - 尝试 `qkflow --help` 查看所有命令
2. **自定义配置** - 编辑 `~/.qkflow/config.yaml` 或检查 iCloud Drive 位置
3. **与团队分享** - 帮助团队成员迁移
4. **提供反馈** - 让我们知道如何改进

---

**欢迎使用更快、更好的 Quick Workflow！🎉**

