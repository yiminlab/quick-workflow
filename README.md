# quick-workflow

> Streamlined GitHub and Jira workflow automation tools

A modern CLI tool to automate your daily GitHub PR and Jira workflow, written in Go.

## 💡 Why use it

> Are you tired of the hassle of changing the status every time you create a PR, such as 'code review,' 'merged,' and dealing with repetitive tasks?

> Especially when you're busy and have opened multiple browser windows, waiting for pages to load just to complete these mundane tasks can be time-consuming, diverting our focus from more important matters.

**Highlighted Benefits:**

- **Efficiency**: Streamline the process by automating status updates and repetitive tasks.
- **Time-Saving**: Eliminate the need to open multiple browser tabs and wait for pages to load.
- **Enhanced Focus**: Allow users to concentrate on more crucial tasks.
- **Reduced Manual Work**: Minimize the need for manual status updates for each PR.
- **Standardized Workflow**: Ensures a consistent and error-free process.

## 🚀 Features

- **Automated PR Creation**: Create branch, commit, push, and open PR in one command
- **Jira Integration**: Automatically link PRs to Jira tickets and update status
- **Smart Branch Management**: Auto-detect default branch (main/master)
- **Interactive CLI**: User-friendly prompts with support for cancellation (Ctrl+C)
- **One-Command Merge**: Merge PR, delete branch, and update Jira status
- **Browser Integration**: Auto-open PR in browser and copy URL to clipboard
- **Quick Update**: Commit and push with PR title as commit message (`qkflow update`)
- **iCloud Sync**: Automatic config sync across Mac devices (macOS only) ☁️
- **Auto Update**: Automatically check and install updates (configurable) 🔄

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/yiminlab/quick-workflow?style=flat&logo=github)](https://github.com/yiminlab/quick-workflow/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/yiminlab/quick-workflow/build.yml?branch=master&style=flat&logo=github-actions)](https://github.com/yiminlab/quick-workflow/actions)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat)](https://github.com/yiminlab/quick-workflow)
[![iCloud Sync](https://img.shields.io/badge/iCloud-Sync%20Enabled-blue?style=flat&logo=icloud)](go-version/ICLOUD_MIGRATION.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](go-version/CONTRIBUTING.md)

## 📦 Installation

### Download Binary (Recommended)

```bash
# macOS (Apple Silicon)
curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Linux
curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/yiminlab/quick-workflow/releases/latest/download/qkflow-windows-amd64.exe -OutFile qkflow.exe
# Move qkflow.exe to a directory in your PATH
```

### Build from Source

```bash
git clone https://github.com/yiminlab/quick-workflow.git
cd quick-workflow/go-version
make gen      # Initialize dependencies
make build    # Build binary
make install  # Install to GOPATH/bin
```

## ⚙️ Initial Setup

Run the interactive setup wizard:

```bash
qkflow init
```

This will prompt you for:
- GitHub token (auto-detected from `gh` CLI if available)
- GitHub owner (username or organization)
- GitHub repository name
- Jira URL (e.g., https://your-company.atlassian.net)
- Jira email
- Jira API token ([Get one here](https://id.atlassian.com/manage-profile/security/api-tokens))
- Optional: OpenAI API key for AI features

**Configuration Storage:**

✨ **NEW**: On macOS, all configs are automatically saved to iCloud Drive and synced across your devices!

- **macOS with iCloud Drive**: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/` ☁️
- **Local Storage** (fallback): `~/.qkflow/`

Both locations contain:
- `config.yaml` - Main configuration
- `jira-status.json` - Jira status mappings

Run `qkflow config` to see your actual storage location.

## 🎯 Usage

### Create a Pull Request

```bash
# With Jira ticket
qkflow pr create PROJ-123

# Interactive mode (will prompt for ticket)
qkflow pr create
```

**What it does:**
1. ✅ Fetches Jira issue details (if ticket provided)
2. ✅ Prompts for PR title and description (or uses AI to generate)
3. ✅ Creates a new git branch
4. ✅ Commits your staged changes
5. ✅ Pushes to remote
6. ✅ Creates GitHub PR
7. ✅ Adds PR link to Jira
8. ✅ Updates Jira status
9. ✅ Copies PR URL to clipboard

### Merge a Pull Request

```bash
# Merge PR by number
qkflow pr merge 123

# Interactive mode (will show your open PRs)
qkflow pr merge
```

**What it does:**
1. ✅ Fetches PR details
2. ✅ Confirms merge with you
3. ✅ Merges the PR on GitHub
4. ✅ Deletes remote branch
5. ✅ Switches to main/master branch
6. ✅ Deletes local branch
7. ✅ Updates Jira status
8. ✅ Adds merge comment to Jira

### Quick Update (qkupdate equivalent)

```bash
# Quick commit and push with PR title as commit message
qkflow update
```

**What it does:**
1. ✅ Gets the current PR title from GitHub
2. ✅ Stages all changes (`git add --all`)
3. ✅ Commits with PR title as commit message
4. ✅ Pushes to origin
5. ✅ Falls back to "update" if no PR found

Perfect for quick updates to an existing PR!

### Jira Status Management

```bash
# List configured Jira status mappings
qkflow jira list

# Setup status mapping for a project
qkflow jira setup PROJ

# Delete status mapping
qkflow jira delete PROJ
```

### Other Commands

```bash
# Show current configuration and storage location
qkflow config

# Update qkflow to latest version
qkflow update-cli

# Show version
qkflow version

# Get help
qkflow help
qkflow pr --help
qkflow jira --help
```

## 🔧 Configuration

Configuration is stored in `~/.qkflow/config.yaml` (or iCloud on macOS):

```yaml
email: your.email@example.com
github_token: ghp_your_github_token
github_owner: your-username
github_repo: your-repo
jira_api_token: your_jira_token
jira_service_address: https://your-domain.atlassian.net
jira_email: your.email@example.com
branch_prefix: feature  # optional
openai_key: sk-your_openai_key  # optional
deepseek_key: sk-your_deepseek_key  # optional
```

You can also use environment variables:
- `GITHUB_TOKEN`
- `JIRA_API_TOKEN`
- `JIRA_URL`
- `OPENAI_API_KEY`
- `DEEPSEEK_API_KEY`

## 🛠️ Development

```bash
cd go-version

# Initialize dependencies
make gen

# Format and fix code
make fix

# Run linters
make check

# Run tests
make test

# Build
make build

# Build for all platforms
make build-all

# Install to system
make install

# Clean build artifacts
make clean

# Show all commands
make help
```

## 🚀 Release

```bash
# Check and prepare for release
make release VERSION=v1.0.1

# Create and push tag (triggers CI/CD)
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1
```

GitHub Actions will automatically:
- Run tests
- Build binaries for all platforms
- Create a GitHub Release
- Upload all binaries

## 📚 Documentation

- [Getting Started Guide](./go-version/GETTING_STARTED.md) - Quick start tutorial
- [Auto Update Guide](./go-version/AUTO_UPDATE.md) - Automatic update feature
- [Release Guide](./go-version/RELEASE.md) - Detailed release instructions
- [Release Quickstart](./go-version/RELEASE_QUICKSTART.md) - Quick release commands
- [iCloud Sync Guide](./go-version/ICLOUD_MIGRATION.md) - iCloud configuration sync
- [Jira Configuration](./go-version/JIRA_STATUS_CONFIG.md) - Jira setup details
- [Contributing Guide](./go-version/CONTRIBUTING.md) - How to contribute
- [Project Structure](./go-version/STRUCTURE.md) - Codebase overview

## ✨ Why qkflow?

**Before qkflow:**
```bash
git checkout -b feature/my-feature
git add .
git commit -m "Add feature"
git push origin feature/my-feature
gh pr create --title "Add feature" --body "Description..."
# Open Jira, find ticket, add PR link, update status...
```

**With qkflow:**
```bash
qkflow pr create PROJ-123
# Done! ✨
```

## 🔒 Security

- Tokens are stored securely with file permissions `0600`
- iCloud storage is encrypted and secure
- Never commit config files
- Use environment variables for CI/CD

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./go-version/CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts
- [go-github](https://github.com/google/go-github) - GitHub API client
- [go-jira](https://github.com/andygrunwald/go-jira) - Jira API client

## 📞 Support

- 🐛 [Report a bug](https://github.com/yiminlab/quick-workflow/issues/new?labels=bug)
- 💡 [Request a feature](https://github.com/yiminlab/quick-workflow/issues/new?labels=enhancement)
- 📖 [View Documentation](https://github.com/yiminlab/quick-workflow)

---

**Made with ❤️ by [@yiminlab](https://github.com/yiminlab)**

⭐ **Star this repo if you find it useful!**
