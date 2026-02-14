# Quick Workflow (Go Version)

> A modern, blazing-fast CLI tool for streamlined GitHub and Jira workflows

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Wangggym/quick-workflow?style=flat&logo=github)](https://github.com/Wangggym/quick-workflow/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/Wangggym/quick-workflow/build.yml?branch=master&style=flat&logo=github-actions)](https://github.com/Wangggym/quick-workflow/actions)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat)](https://github.com/Wangggym/quick-workflow)
[![iCloud Sync](https://img.shields.io/badge/iCloud-Sync%20Enabled-blue?style=flat&logo=icloud)](docs/en/features/icloud-migration.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

## 🚀 What's New in Go Version

This is a complete rewrite of the original Shell-based quick-workflow tool in Go, bringing:

- **📦 Single Binary** - No dependencies, just download and run
- **⚡ Faster** - Native performance, instant startup
- **🔒 Type Safe** - Catch errors at compile time
- **🎨 Better UX** - Interactive prompts, colored output
- **🧪 Testable** - Comprehensive test coverage
- **🌍 Cross-platform** - Works on macOS, Linux, and Windows
- **☁️ iCloud Sync** - Automatic config sync across Mac devices (macOS only)

## ✨ Features

- **PR Creation** - Create PRs with automatic branch creation, commit, and push
- **PR Editor** - 🆕 Web-based editor for adding descriptions with images/videos 🎨
- **PR Merging** - Merge PRs and clean up branches automatically
- **Quick Update** - Commit and push with PR title as commit message
- **Jira Integration** - Automatically update Jira status and add PR links
- **Jira Reader** - 🆕 Read and export Jira issues (optimized for Cursor AI) 🤖
- **Watch Daemon** - 🆕 Automatically monitor PRs and update Jira when merged ⚡
- **Interactive CLI** - Beautiful prompts and progress indicators
- **Configuration Management** - Simple setup with `qkflow init`
- **iCloud Sync** - Seamlessly sync configs across all your Mac devices ☁️
- **Auto Update** - Automatically check and install updates (24h interval) 🔄

## 📦 Installation

### Option 1: Download Binary (Recommended)

#### macOS Installation

```bash
# macOS (Apple Silicon - M1/M2/M3)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow

# macOS (Intel)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-amd64 -o qkflow

# 🔑 Important: Remove macOS quarantine attribute to bypass Gatekeeper (if needed)
# Note: If you get "No such xattr: com.apple.quarantine", that's fine - skip this step
xattr -d com.apple.quarantine qkflow 2>/dev/null || echo "No quarantine attribute found (this is fine)"

# Make executable and install
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Verify installation
qkflow version
```

> **⚠️ macOS Security Notice**: If you see "qkflow-darwin-arm64 cannot be opened because Apple cannot verify that it is free of malware", this is normal for unsigned binaries. The `xattr -d com.apple.quarantine` command above will resolve this issue safely.

#### Linux Installation

```bash
# Linux (x86_64)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-linux-amd64 -o qkflow
chmod +x qkflow
sudo mv qkflow /usr/local/bin/

# Verify installation
qkflow version
```

### Option 2: Install with Go

```bash
go install github.com/Wangggym/quick-workflow/cmd/qkflow@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/Wangggym/quick-workflow.git
cd quick-workflow/go-version
make build
sudo cp bin/qkflow /usr/local/bin/
```

## ⚙️ Setup

### Prerequisites

- Git installed and configured
- GitHub CLI (`gh`) installed and authenticated: `gh auth login`
- Jira API token: [Get one here](https://id.atlassian.com/manage-profile/security/api-tokens)

### Initial Configuration

Run the interactive setup:

```bash
qkflow init
```

This will prompt you for:
- Email address
- GitHub token (auto-detected from `gh` CLI)
- Jira service address (e.g., https://your-domain.atlassian.net)
- Jira API token
- Optional: Branch prefix
- Optional: OpenAI API key for AI features

**Configuration Storage:**

✨ **NEW**: On macOS, all configs are automatically saved to iCloud Drive in a single directory and synced across all your devices!

- **macOS with iCloud Drive**: Synced across devices ☁️
  - 📂 All configs in: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
- **Local Storage** (fallback):
  - 📂 All configs in: `~/.qkflow/`

Both locations contain:
- `config.yaml` - Main configuration
- `jira-status.json` - Jira status mappings

Run `qkflow config` to see your actual storage location.

📖 See [iCloud Migration Guide](docs/en/features/icloud-migration.md) ([中文](docs/cn/features/icloud-migration.md)) for more details.

> 📚 **Documentation**: All documentation is now organized in the [`docs/`](docs/) directory. See [Documentation Index](docs/README.md) for a complete list.
>
> 🚀 **Quick Start**: See [QUICKSTART.md](QUICKSTART.md) to get started in 5 minutes (English & 中文).

## 🎯 Usage

### Create a Pull Request

```bash
# With Jira ticket
qkflow pr create PROJ-123

# Without Jira ticket (will prompt)
qkflow pr create

# Interactive mode (no arguments)
qkflow pr create
```

**What it does:**
1. ✅ Fetches Jira issue details (if ticket provided)
2. ✅ Prompts for PR title and description
3. ✅ Lets you select change types (feat, fix, docs, etc.)
4. ✅ ⭐ **NEW**: Optionally add detailed description with images/videos (web editor)
5. ✅ Creates a new git branch
6. ✅ Commits your staged changes
7. ✅ Pushes to remote
8. ✅ Creates GitHub PR
9. ✅ ⭐ **NEW**: Uploads files and adds comment to GitHub & Jira
10. ✅ Adds PR link to Jira
11. ✅ Updates Jira status (optional)
12. ✅ Copies PR URL to clipboard

### Merge a Pull Request

```bash
# Merge PR by number
qkflow pr merge 123

# Merge PR by URL (works from anywhere!)
qkflow pr merge https://github.com/brain/planning-api/pull/2001

# Interactive mode (auto-detect from current branch)
qkflow pr merge
```

**What it does:**
1. ✅ Supports PR number OR full GitHub URL
2. ✅ Fetches PR details
3. ✅ Confirms merge with you
4. ✅ Merges the PR on GitHub
5. ✅ Deletes remote branch (optional)
6. ✅ Switches to main branch
7. ✅ Deletes local branch (optional)
8. ✅ Updates Jira status (optional)
9. ✅ Adds merge comment to Jira

### Approve a Pull Request

```bash
# Approve PR by number (with default 👍 comment)
qkflow pr approve 123

# Approve PR by URL (works from anywhere!)
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# URL also works with /files, /commits, /checks suffixes
qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

# Custom comment
qkflow pr approve 123 --comment "LGTM! 🎉"
qkflow pr approve 123 -c "Looks good!"

# Approve and auto-merge (with default 👍 comment)
qkflow pr approve 123 --merge
qkflow pr approve 123 -m

# Approve by URL with custom comment and merge
qkflow pr approve https://github.com/owner/repo/pull/456 -c "Ship it! 🚀" -m

# Interactive mode (auto-detect PR from current branch)
qkflow pr approve
```

**What it does:**
1. ✅ Supports PR number OR full GitHub URL (including /files, /commits, /checks paths)
2. ✅ Auto-detects PR from current branch (if no argument provided)
3. ✅ Fetches PR details
4. ✅ Approves the PR on GitHub
5. ✅ Adds a comment (default: 👍, customize with -c flag)
6. ✅ Optionally auto-merges after approval (with --merge flag)
7. ✅ Checks if PR is mergeable before merging
8. ✅ Cleans up branches after merge (if merged)

**Examples:**
```bash
# Simple approve (uses default 👍 comment)
qkflow pr approve 123

# Approve from Files tab URL
qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

# Custom comment
qkflow pr approve 123 -c "LGTM!"

# Approve PR from another repo by URL
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# Approve and merge in one command (with 👍)
qkflow pr approve 123 --merge

# Approve current branch's PR
qkflow pr approve  # Will find PR automatically
```

### Quick Update (qkupdate)

```bash
# Quick commit and push with PR title as commit message
qkflow update
```

**What it does:**
1. ✅ Gets the current PR title from GitHub
2. ✅ Stages all changes (git add --all)
3. ✅ Commits with PR title as commit message
4. ✅ Pushes to origin
5. ✅ Falls back to "update" if no PR found

This is perfect for quick updates to an existing PR!

### PR Editor (Add Rich Descriptions)

**NEW!** Add detailed descriptions with images and videos to your PRs using a beautiful web-based editor.

```bash
# During pr create, you'll be prompted:
? Add detailed description with images/videos?
  > ⏭️  Skip (default)    # Press Enter to skip
    ✅ Yes, continue       # Press Space to toggle, then Enter

# If you select "Yes, continue":
🌐 Opening editor in your browser...
📝 Please edit your content in the browser and click 'Save and Continue'
```

**The web editor provides:**

- 📝 **Markdown Editor** with live preview and formatting toolbar
- 🖼️ **Drag & Drop** images and videos from Finder/Explorer
- 📋 **Paste** images directly from clipboard (Cmd+V / Ctrl+V)
- 🎨 **GitHub-style UI** with dark theme
- ⚡ **Instant Upload** to both GitHub PR and Jira issue
- 🔄 **Auto-conversion** of local paths to online URLs

**Supported formats:**
- Images: PNG, JPG, JPEG, GIF, WebP, SVG
- Videos: MP4, MOV, WebM, AVI

**What happens after you save:**

1. ✅ Files are uploaded to GitHub and Jira
2. ✅ Markdown paths are replaced with actual URLs
3. ✅ Description is added as a comment to the GitHub PR
4. ✅ Same description is added as a comment to the Jira issue
5. ✅ Temporary files are cleaned up

**Perfect for:**
- Bug fixes with before/after screenshots
- Features with demo videos
- Visual documentation of UI changes
- Architecture diagrams
- Test results and screenshots

### Jira Reader (Cursor AI Integration)

**NEW!** Read and export Jira issues, optimized for Cursor AI.

```bash
# Intelligent read (recommended for Cursor AI)
qkflow jira read NA-9245

# Quick terminal view
qkflow jira show NA-9245
qkflow jira show NA-9245 --full    # Full details

# Export to files
qkflow jira export NA-9245
qkflow jira export NA-9245 --with-images    # Include images

# Clean up exports
qkflow jira clean NA-9245
qkflow jira clean --all
```

**For Cursor AI users:**

Simply tell Cursor in chat:
```
"通过 qkflow 工具读取 NA-9245 所有内容并总结"
```

Cursor will automatically:
1. ✅ Run the qkflow command
2. ✅ Read exported files (including images)
3. ✅ Provide comprehensive analysis

📖 See [Jira Integration Guide](docs/en/features/jira-integration.md) ([中文](docs/cn/features/jira-integration.md)) for detailed documentation.

### Watch Daemon (Auto-update Jira)

**NEW!** Automatically monitor your PRs and update Jira when they're merged.

```bash
# Install and start watch daemon (with auto-start on login)
qkflow watch install

# Check daemon status
qkflow watch status

# View processing history
qkflow watch history

# View logs
qkflow watch log
qkflow watch log --follow

# Manual check (without daemon)
qkflow watch check

# Stop/Start daemon
qkflow watch stop
qkflow watch start

# Uninstall daemon
qkflow watch uninstall
```

**What it does:**
- ✅ Monitors YOUR PRs every 15 minutes (8:30-24:00)
- ✅ Night mode: checks at 2:00 and 6:00 only
- ✅ Auto-updates Jira status when PR is merged
- ✅ Desktop notifications (macOS)
- ✅ Auto-start on login (launchd on macOS)
- ✅ Logs all activities
- ✅ No manual intervention needed!

**Prerequisites:**
1. Run `qkflow jira setup` first to configure Jira status mappings
2. Make sure "PR Merged" status is configured (default: "In Review")

📖 See [Jira Integration Guide](docs/en/features/jira-integration.md) ([中文](docs/cn/features/jira-integration.md)) for status configuration details.

### Other Commands

```bash
# Show current configuration
qkflow config

# Update qkflow to latest version
qkflow update-cli

# Show version
qkflow version

# Get help
qkflow help
qkflow pr --help
```

## 🏗️ Project Structure

```
go-version/
├── cmd/
│   └── qkflow/
│       ├── main.go              # Entry point
│       └── commands/
│           ├── root.go          # Root command
│           ├── init.go          # Setup wizard
│           ├── pr.go            # PR commands
│           ├── pr_create.go     # Create PR
│           ├── pr_merge.go      # Merge PR
│           ├── update.go        # Quick update
│           └── jira.go          # Jira commands
├── internal/
│   ├── github/
│   │   └── client.go           # GitHub API client
│   ├── jira/
│   │   ├── client.go           # Jira API client
│   │   └── status_cache.go     # Status cache
│   ├── git/
│   │   └── operations.go       # Git operations
│   ├── ai/
│   │   └── client.go           # AI client
│   └── ui/
│       └── prompt.go           # User interface
├── pkg/
│   └── config/
│       └── config.go           # Configuration management
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make (optional but recommended)

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to GOPATH/bin
make install

# Run tests
make test

# Run linters
make lint

# Format code
make fmt
```

### Quick Development

```bash
# Run without building
go run ./cmd/qkflow pr create

# Or use Makefile
make run-pr-create
make run-pr-merge
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test ./internal/github/...
go test ./internal/jira/...
```

## 📝 Configuration Files

### Storage Location

qkflow intelligently stores configurations based on your system:

**macOS with iCloud Drive** (Recommended):
- Configs are stored in iCloud Drive and automatically synced across your Mac devices
- All configs in: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
  - `config.yaml` - Main configuration
  - `jira-status.json` - Jira status mappings

**Local Storage** (Fallback):
- All configs in: `~/.qkflow/`
  - `config.yaml` - Main configuration
  - `jira-status.json` - Jira status mappings

Run `qkflow config` to see your actual storage location.

### Configuration Format

```yaml
email: your.email@example.com
jira_api_token: your_jira_token
jira_service_address: https://your-domain.atlassian.net
github_token: ghp_your_github_token
branch_prefix: feature  # optional
openai_key: sk-your_openai_key  # optional
```

## 🔒 Security

- Tokens are stored securely in your config directory (local or iCloud)
- File permissions are set to `0600` (user read/write only)
- iCloud storage is encrypted and secure
- Never commit the config file or share it
- Use environment variables for CI/CD:
  ```bash
  export QK_GITHUB_TOKEN=xxx
  export QK_JIRA_API_TOKEN=xxx
  ```

**Note**: If using iCloud Drive, your configurations will sync across your Mac devices automatically, providing a seamless experience.

## 🚧 Migration from Shell Version

See [Migration Guide](docs/en/migration/migration.md) ([中文](docs/cn/migration/migration.md)) for detailed migration guide.

**Quick comparison:**

| Feature | Shell Version | Go Version |
|---------|--------------|------------|
| Installation | Clone + dependencies | Single binary |
| Configuration | `.zshrc` / `.bashrc` | `qk init` |
| Dependencies | `gh`, `jira`, `jq`, etc. | None (self-contained) |
| Speed | ~1-2s startup | <100ms startup |
| Platform | macOS/Linux | macOS/Linux/Windows |
| Updates | `git pull` | Download new binary |

## 🔧 Troubleshooting

### macOS Installation Issues

#### "qkflow-darwin-arm64 cannot be opened" Error

If you encounter this error:
```
"qkflow-darwin-arm64" cannot be opened because Apple cannot verify that it is free of malware that may harm your Mac or compromise your privacy.
```

**Solution 1: Remove Quarantine Attribute (Recommended)**
```bash
# Remove quarantine attribute (if it exists)
xattr -d com.apple.quarantine qkflow-darwin-arm64 2>/dev/null || echo "No quarantine attribute (this is fine)"
chmod +x qkflow-darwin-arm64
```

> **Note**: If you see "No such xattr: com.apple.quarantine", that means the file wasn't quarantined and you can skip the xattr step.

**Solution 2: System Settings**
1. Try to run the binary (it will show the security warning)
2. Go to **System Settings** → **Privacy & Security**
3. Scroll down to find the blocked app
4. Click **Open Anyway**

**Solution 3: Temporary Gatekeeper Disable**
```bash
# Disable Gatekeeper temporarily (requires admin)
sudo spctl --master-disable

# Run your binary, then re-enable
sudo spctl --master-enable
```

#### "Permission Denied" Error

```bash
# Make sure the file is executable
chmod +x qkflow

# Check if /usr/local/bin is in your PATH
echo $PATH | grep -q "/usr/local/bin" && echo "✅ PATH is correct" || echo "❌ Add /usr/local/bin to PATH"

# Add to PATH if needed (add to ~/.zshrc or ~/.bash_profile)
export PATH="/usr/local/bin:$PATH"
```

### General Issues

#### Command Not Found

```bash
# Check if qkflow is installed
which qkflow

# If not found, ensure it's in your PATH
ls -la /usr/local/bin/qkflow

# Reload shell configuration
source ~/.zshrc  # or ~/.bash_profile
```

#### Update Issues

```bash
# Check current version
qkflow version

# Manual update (download latest binary and replace)
curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o /tmp/qkflow-new
xattr -d com.apple.quarantine /tmp/qkflow-new
chmod +x /tmp/qkflow-new
sudo mv /tmp/qkflow-new /usr/local/bin/qkflow
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on how to contribute.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Original Shell version: [quick-workflow](https://github.com/Wangggym/quick-workflow)
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts
- [go-github](https://github.com/google/go-github) - GitHub API client
- [go-jira](https://github.com/andygrunwald/go-jira) - Jira API client

## 📞 Support

- 🐛 [Report a bug](https://github.com/Wangggym/quick-workflow/issues/new?labels=bug)
- 💡 [Request a feature](https://github.com/Wangggym/quick-workflow/issues/new?labels=enhancement)
- 📖 [Documentation](https://github.com/Wangggym/quick-workflow/wiki)

---

Made with ❤️ by [Wangggym](https://github.com/Wangggym)

