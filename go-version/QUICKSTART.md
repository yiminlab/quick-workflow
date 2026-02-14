# Quick Start Guide / 快速入门指南

[English](#english) | [中文](#中文)

---

## English {#english}

Get up and running with Quick Workflow Go version in 5 minutes!

### 📦 Installation (30 seconds)

#### macOS (Apple Silicon)
```bash
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow && \
chmod +x qkflow && \
sudo mv qkflow /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow && \
chmod +x qkflow && \
sudo mv qkflow /usr/local/bin/
```

#### Linux
```bash
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow && \
chmod +x qkflow && \
sudo mv qkflow /usr/local/bin/
```

#### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-windows-amd64.exe -OutFile qkflow.exe
# Move qkflow.exe to a directory in your PATH
```

### ⚙️ Setup (2 minutes)

#### 1. Ensure Prerequisites

```bash
# Install and authenticate GitHub CLI
brew install gh
gh auth login

# Get Jira API token
# Visit: https://id.atlassian.com/manage-profile/security/api-tokens
```

#### 2. Run Setup Wizard

```bash
qkflow init
```

Answer the prompts:
- **Email**: Your work email
- **GitHub Token**: Auto-detected from `gh` CLI
- **Jira Address**: `https://your-domain.atlassian.net`
- **Jira Token**: Paste the token from step 1
- **Branch Prefix**: Optional (e.g., `feature` or your username)

**Configuration Storage:**

✨ **NEW**: On macOS, all configs are automatically saved to iCloud Drive and synced across all your devices!

- **macOS with iCloud Drive**: Synced across devices ☁️
  - 📂 All configs in: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
- **Local Storage** (fallback):
  - 📂 All configs in: `~/.qkflow/`

Run `qkflow config` to see your actual storage location.

### 🎯 Your First PR (2 minutes)

#### Step 1: Make Your Changes

```bash
cd your-project
git checkout -b feature/test

# Make some changes
echo "# Test" >> README.md
git add README.md
```

#### Step 2: Create PR

```bash
qkflow pr create PROJ-123
```

Follow the prompts:
1. **Title**: Accept suggested or enter custom
2. **Description**: Optional short description
3. **Change Types**: Select applicable types (feat, fix, etc.)
4. **Jira Status**: Choose new status (optional)

**Done!** Your PR is created and Jira is updated! 🎉

### 🔄 Merge a PR (1 minute)

```bash
qkflow pr merge 123
```

Follow the prompts:
1. **Confirm merge**: Review PR details
2. **Delete branches**: Choose to clean up
3. **Update Jira**: Set final status

**Done!** PR merged and cleaned up! 🎉

### 💡 Pro Tips

#### Use Without Jira

```bash
# Skip Jira ticket (press Enter when prompted)
qkflow pr create
```

#### Keyboard Shortcuts in Prompts

- **Arrow keys**: Navigate options
- **Space**: Select/deselect (multi-select)
- **Enter**: Confirm selection
- **Ctrl+C**: Cancel operation

#### Quick Commands

```bash
# Show config
qkflow config

# Show version
qkflow version

# Get help
qkflow --help
qkflow pr --help
qkflow jira --help
```

### 🎨 Example Workflow

Here's a complete workflow example:

```bash
# 1. Start new feature
cd ~/projects/my-app
git checkout main
git pull

# 2. Make changes
git checkout -b feature/awesome-feature
# ... make your changes ...
git add .

# 3. Create PR
qkflow pr create PROJ-456
# Title: "Add awesome feature"
# Description: "This adds X, Y, Z"
# Types: [x] feat: New feature
# Jira Status: In Review

# Output:
# ✅ Branch created: feature/PROJ-456--Add-awesome-feature
# ✅ Changes committed
# ✅ Pushed to remote
# ✅ Pull request created: https://github.com/org/repo/pull/789
# ✅ Added PR link to Jira
# ✅ Updated Jira status to: In Review
# ✅ All done! 🎉

# 4. Get code review, make changes if needed
# ... after approval ...

# 5. Merge PR
qkflow pr merge 789
# Confirm merge: Yes
# Delete remote branch: Yes
# Delete local branch: Yes
# Update Jira: Yes
# New status: Done

# Output:
# ✅ Pull request merged!
# ✅ Remote branch deleted
# ✅ Local branch deleted
# ✅ Updated Jira status to: Done
# ✅ All done! 🎉
```

### 🐛 Common Issues

#### "Command not found: qkflow"

```bash
# Check if binary exists
ls -l /usr/local/bin/qkflow

# Check PATH
echo $PATH | grep -q "/usr/local/bin" && echo "OK" || echo "Add to PATH"

# Add to PATH if needed (add to ~/.zshrc)
export PATH="/usr/local/bin:$PATH"
```

#### "Failed to create GitHub client"

```bash
# Ensure gh is authenticated
gh auth status

# If not authenticated
gh auth login

# Re-run qkflow init
qkflow init
```

#### "Failed to get Jira issue"

```bash
# Verify Jira credentials
curl -u "your.email@example.com:your_jira_token" \
  https://your-domain.atlassian.net/rest/api/2/myself

# If fails, get new API token and re-run qkflow init
```

### 📚 Learn More

- **Full Documentation**: [README.md](README.md)
- **Migration Guide**: [Migration Guide](docs/en/migration/migration.md) ([中文](docs/cn/migration/migration.md))
- **GitHub Issues**: [Report bugs or request features](https://github.com/Wangggym/quick-workflow/issues)

### 🎉 You're Ready!

Congratulations! You're now set up with Quick Workflow. Enjoy your streamlined workflow!

**Common Commands to Remember:**
```bash
qkflow pr create      # Create PR
qkflow pr merge       # Merge PR
qkflow update         # Quick update
qkflow config         # Show config
qkflow --help         # Get help
```

---

## 中文 {#中文}

5 分钟内快速上手 Quick Workflow Go 版本！

### 📥 安装（30 秒）

#### 方式 1: 下载预编译二进制（推荐）

访问 [Releases 页面](https://github.com/Wangggym/quick-workflow/releases) 下载适合你系统的版本：

```bash
# macOS Apple Silicon (M1/M2/M3)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# macOS Intel
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Linux
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-windows-amd64.exe -OutFile qkflow.exe
# 将 qkflow.exe 移动到 PATH 中的目录
```

#### 方式 2: 从源码构建

```bash
git clone https://github.com/Wangggym/quick-workflow.git
cd quick-workflow/go-version
make gen      # 初始化依赖
make build    # 构建
make install  # 安装到 GOPATH/bin
```

### ⚙️ 初始化配置（2 分钟）

首次使用需要配置 GitHub 和 Jira 信息：

```bash
qkflow init
```

按提示输入：
- GitHub Token (从 https://github.com/settings/tokens 获取)
- GitHub Owner (你的用户名或组织名)
- GitHub Repo (仓库名)
- Jira URL (如 https://your-company.atlassian.net)
- Jira Email
- Jira API Token (从 Jira 账户设置获取)

#### 📱 iCloud 同步 (macOS)

配置会自动保存到 iCloud Drive（如果可用），实现多设备同步！

查看配置位置：
```bash
qkflow config
```

**配置存储位置：**

✨ **新特性**：在 macOS 上，所有配置会自动保存到 iCloud Drive，并在你的所有设备间同步！

- **macOS with iCloud Drive**：跨设备同步 ☁️
  - 📂 所有配置在：`~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
- **本地存储**（回退方案）：
  - 📂 所有配置在：`~/.qkflow/`

运行 `qkflow config` 查看实际的存储位置。

### 🎯 核心功能

#### 1. 创建 PR

```bash
# 在当前分支创建 PR
qkflow pr create

# 带 Jira ticket
qkflow pr create PROJ-123

# 交互式选择 AI 生成标题和描述
```

#### 2. 合并 PR

```bash
# 合并当前分支的 PR
qkflow pr merge

# 合并指定 PR
qkflow pr merge 123

# 自动更新 Jira 状态
```

#### 3. 快速更新 (qkupdate)

```bash
# 自动使用 PR 标题作为 commit message
# 等同于: git add --all && git commit -m "PR Title" && git push
qkflow update
```

#### 4. Jira 状态管理

```bash
# 查看已配置的项目
qkflow jira list

# 为项目配置状态映射
qkflow jira setup PROJECT-KEY

# 删除项目配置
qkflow jira delete PROJECT-KEY
```

### 📚 常用命令

```bash
# 查看版本
qkflow version

# 查看配置
qkflow config

# 查看帮助
qkflow --help
qkflow pr --help
qkflow jira --help
```

### 🔧 开发者命令

```bash
cd go-version

# 初始化依赖
make gen

# 格式化代码
make fix

# 运行测试
make test

# 构建
make build

# 安装到系统
make install

# 清理构建产物
make clean

# 查看所有命令
make help
```

### 🎓 工作流示例

#### 完整的功能开发流程

```bash
# 1. 创建并切换到功能分支
git checkout -b feature/add-login

# 2. 开发功能...
# (编写代码)

# 3. 快速提交和推送
qkflow update

# 4. 创建 PR
qkflow pr create

# 5. Code Review...
# (等待审核通过)

# 6. 合并 PR（自动更新 Jira）
qkflow pr merge
```

#### Bug 修复流程

```bash
# 1. 创建 bugfix 分支
git checkout -b bugfix/fix-login-error

# 2. 修复 bug...
# (修改代码)

# 3. 快速更新
qkflow update

# 4. 创建 PR
qkflow pr create

# 5. 合并
qkflow pr merge
```

### 🔐 环境变量（可选）

可以使用环境变量代替配置文件：

```bash
export GITHUB_TOKEN=your_token
export GITHUB_OWNER=your_username
export GITHUB_REPO=your_repo
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your@email.com
export JIRA_TOKEN=your_jira_token
export OPENAI_API_KEY=your_openai_key  # 可选
export DEEPSEEK_API_KEY=your_deepseek_key  # 可选
```

### 💡 小贴士

1. **第一次使用**: 运行 `qkflow init` 配置
2. **查看配置**: 运行 `qkflow config` 查看存储位置
3. **快速更新**: 使用 `qkflow update` 代替繁琐的 git 命令
4. **Jira 集成**: 配置后 PR 操作自动更新 Jira 状态
5. **iCloud 同步**: macOS 用户配置自动同步到所有设备

### 🐛 常见问题

#### "Command not found: qkflow"

```bash
# 检查二进制文件是否存在
ls -l /usr/local/bin/qkflow

# 检查 PATH
echo $PATH | grep -q "/usr/local/bin" && echo "OK" || echo "需要添加到 PATH"

# 如果需要，添加到 PATH（添加到 ~/.zshrc）
export PATH="/usr/local/bin:$PATH"
```

#### "Failed to create GitHub client"

```bash
# 确保 gh 已认证
gh auth status

# 如果未认证
gh auth login

# 重新运行 qkflow init
qkflow init
```

#### "Failed to get Jira issue"

```bash
# 验证 Jira 凭据
curl -u "your.email@example.com:your_jira_token" \
  https://your-domain.atlassian.net/rest/api/2/myself

# 如果失败，获取新的 API token 并重新运行 qkflow init
```

### 📖 更多文档

- [README.md](README.md) - 完整功能介绍
- [Release Guide](docs/en/release/release.md) ([中文](docs/cn/release/release.md)) - Release 详细指南
- [Release Guide](docs/en/release/release.md) ([中文](docs/cn/release/release.md)) - Release 详细指南（包含快速入门）
- [iCloud Migration](docs/en/features/icloud-migration.md) ([中文](docs/cn/features/icloud-migration.md)) - iCloud 同步指南
- [Jira Integration](docs/en/features/jira-integration.md) ([中文](docs/cn/features/jira-integration.md)) - Jira 集成指南（Issue Reader & 状态配置）
- [Contributing Guide](CONTRIBUTING.md) - 贡献指南

### 🆘 获取帮助

- 遇到问题？查看 [Issues](https://github.com/Wangggym/quick-workflow/issues)
- 想要新功能？提交 [Feature Request](https://github.com/Wangggym/quick-workflow/issues/new)
- 贡献代码？查看 [Contributing Guide](CONTRIBUTING.md)

### 🎉 开始使用

```bash
# 安装
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# 初始化
qkflow init

# 开始使用！
qkflow pr create
```

**常用命令：**
```bash
qkflow pr create      # 创建 PR
qkflow pr merge       # 合并 PR
qkflow update         # 快速更新
qkflow config         # 查看配置
qkflow --help         # 获取帮助
```

---

**Happy coding! 🚀**
