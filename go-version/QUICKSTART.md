# Quick Start Guide

Get up and running with Quick Workflow Go version in 5 minutes!

## 📦 Installation (30 seconds)

### macOS (Apple Silicon)
```bash
curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qk-darwin-arm64 -o qk && \
chmod +x qk && \
sudo mv qk /usr/local/bin/
```

### macOS (Intel)
```bash
curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qk-darwin-amd64 -o qk && \
chmod +x qk && \
sudo mv qk /usr/local/bin/
```

### Linux
```bash
curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qk-linux-amd64 -o qk && \
chmod +x qk && \
sudo mv qk /usr/local/bin/
```

## ⚙️ Setup (2 minutes)

### 1. Ensure Prerequisites

```bash
# Install and authenticate GitHub CLI
brew install gh
gh auth login

# Get Jira API token
# Visit: https://id.atlassian.com/manage-profile/security/api-tokens
```

### 2. Run Setup Wizard

```bash
qk init
```

Answer the prompts:
- **Email**: Your work email
- **GitHub Token**: Auto-detected from `gh` CLI
- **Jira Address**: `https://your-domain.atlassian.net`
- **Jira Token**: Paste the token from step 1
- **Branch Prefix**: Optional (e.g., `feature` or your username)

## 🎯 Your First PR (2 minutes)

### Step 1: Make Your Changes

```bash
cd your-project
git checkout -b feature/test

# Make some changes
echo "# Test" >> README.md
git add README.md
```

### Step 2: Create PR

```bash
qk pr create PROJ-123
```

Follow the prompts:
1. **Title**: Accept suggested or enter custom
2. **Description**: Optional short description
3. **Change Types**: Select applicable types (feat, fix, etc.)
4. **Jira Status**: Choose new status (optional)

**Done!** Your PR is created and Jira is updated! 🎉

## 🔄 Merge a PR (1 minute)

```bash
qk pr merge 123
```

Follow the prompts:
1. **Confirm merge**: Review PR details
2. **Delete branches**: Choose to clean up
3. **Update Jira**: Set final status

**Done!** PR merged and cleaned up! 🎉

## 💡 Pro Tips

### Use Without Jira

```bash
# Skip Jira ticket (press Enter when prompted)
qk pr create
```

### Keyboard Shortcuts in Prompts

- **Arrow keys**: Navigate options
- **Space**: Select/deselect (multi-select)
- **Enter**: Confirm selection
- **Ctrl+C**: Cancel operation

### Quick Commands

```bash
# Show config
qk config

# Show version
qk version

# Get help
qk --help
qk pr --help
```

## 🎨 Example Workflow

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
qk pr create PROJ-456
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
qk pr merge 789
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

## 🐛 Common Issues

### "Command not found: qk"

```bash
# Check if binary exists
ls -l /usr/local/bin/qk

# Check PATH
echo $PATH | grep -q "/usr/local/bin" && echo "OK" || echo "Add to PATH"

# Add to PATH if needed (add to ~/.zshrc)
export PATH="/usr/local/bin:$PATH"
```

### "Failed to create GitHub client"

```bash
# Ensure gh is authenticated
gh auth status

# If not authenticated
gh auth login

# Re-run qk init
qk init
```

### "Failed to get Jira issue"

```bash
# Verify Jira credentials
curl -u "your.email@example.com:your_jira_token" \
  https://your-domain.atlassian.net/rest/api/2/myself

# If fails, get new API token and re-run qk init
```

## 📚 Learn More

- **Full Documentation**: [README.md](README.md)
- **Migration Guide**: [MIGRATION.md](MIGRATION.md)
- **GitHub Issues**: [Report bugs or request features](https://github.com/yiminlab/quick-workflow/issues)

## 🎉 You're Ready!

Congratulations! You're now set up with Quick Workflow. Enjoy your streamlined workflow!

**Common Commands to Remember:**
```bash
qk pr create      # Create PR
qk pr merge       # Merge PR
qk config         # Show config
qk --help         # Get help
```

---

**Happy coding! 🚀**

