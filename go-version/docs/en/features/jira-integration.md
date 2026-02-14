# Jira Integration Guide

Complete guide for Jira integration features in Quick Workflow.

## 📚 Table of Contents

1. [Jira Issue Reader](#jira-issue-reader) - Read and export Jira issues
2. [Jira Status Configuration](#jira-status-configuration) - Configure status mappings for PR workflow

---

## Part 1: Jira Issue Reader {#jira-issue-reader}

A powerful tool to read and export Jira issues, optimized for use with Cursor AI.

### 🎯 Quick Start

#### For Cursor Users (Recommended)

The simplest way to use this in Cursor:

```bash
# In Cursor terminal, run:
qkflow jira read NA-9245

# Then in Cursor chat, simply say:
"Summarize the Jira ticket content I just read"
```

Cursor will automatically read the exported files and provide you with a comprehensive summary!

### 📚 Available Commands

#### 1. `show` - Quick Terminal View

Display issue information directly in the terminal.

```bash
# Basic view (metadata only)
qkflow jira show NA-9245

# Full view (includes description and comments)
qkflow jira show NA-9245 --full
```

**Use when:**
- You need a quick peek at the issue
- You only need text content
- You want the fastest response

#### 2. `export` - Complete Export with Files

Export issue to local files with optional images.

```bash
# Export text only
qkflow jira export NA-9245

# Export with all images and attachments
qkflow jira export NA-9245 --with-images

# Export to custom directory
qkflow jira export NA-9245 -o ~/jira-exports/ --with-images
```

**Output structure:**
```
/tmp/qkflow/jira/NA-9245/
├── README.md           # How to use in Cursor
├── content.md          # Main content (Markdown)
└── attachments/        # Downloaded files (if --with-images)
    ├── screenshot.png
    └── diagram.jpg
```

**Use when:**
- You need images/attachments
- You want to keep a local copy
- You need formatted Markdown

#### 3. `read` - Intelligent Mode ⭐️ **RECOMMENDED**

Automatically decides the best way to present the issue.

```bash
# Auto mode (smart decision)
qkflow jira read NA-9245
```

**How it works:**
- ✅ **Has images?** → Exports to files with images
- ✅ **Text only?** → Displays directly in terminal
- ✅ Automatically optimized for Cursor

**Use when:**
- Working with Cursor AI (best experience)
- You want the tool to decide the best format
- You're not sure if the issue has images

#### 4. `clean` - Clean Up Exports

Remove exported files to free up disk space.

```bash
# Clean specific issue
qkflow jira clean NA-9245

# Clean all exports
qkflow jira clean --all

# Preview what would be deleted (dry run)
qkflow jira clean --all --dry-run

# Force delete without confirmation
qkflow jira clean --all --force
```

### 🎨 Usage Examples in Cursor

#### Example 1: Simple Text Analysis

```
You in Cursor: "Read NA-9245 using qkflow and summarize"

Cursor executes: qkflow jira read NA-9245
Cursor responds: "This ticket (NA-9245) is about..."
```

#### Example 2: With Images

```
You in Cursor: "Use qkflow to read NA-9245 with all images, analyze the architecture design"

Cursor executes: qkflow jira export NA-9245 --with-images
Cursor reads: content.md + all images in attachments/
Cursor responds: "Based on the architecture diagram, this system includes..."
```

#### Example 3: Manual Control

```bash
# Step 1: Export (you run this)
qkflow jira export NA-9245 --with-images

# Step 2: Tell Cursor what to read
"Read /tmp/qkflow/jira/NA-9245/content.md and analyze the architecture diagram"

# Step 3: Clean up when done
qkflow jira clean NA-9245
```

### 💡 Tips & Best Practices

#### For Cursor Users

1. **Use `read` command by default** - It's optimized for AI consumption
2. **Be specific in prompts** - Tell Cursor what you want to know
3. **Clean up regularly** - Use `clean --all` to free up space

#### Command Comparison

| Command | Speed | Images | Output | Best For |
|---------|-------|--------|--------|----------|
| `show` | ⚡️ Fastest | ❌ | Terminal | Quick peek |
| `show --full` | ⚡️ Fast | ❌ | Terminal | Full text |
| `export` | 🐌 Slower | ❌ | Files | Text archive |
| `export --with-images` | 🐌 Slowest | ✅ | Files | Complete archive |
| `read` ⭐️ | ⚡️ Smart | ✅ Smart | Smart | **Cursor AI** |

#### Cursor Prompt Templates

```bash
# General summary
"Read <ISSUE-KEY> using qkflow and summarize the content"

# Specific analysis
"Use qkflow to read <ISSUE-KEY> and analyze the technical solution"

# With context
"Read <ISSUE-KEY> and compare with our current implementation"

# With images
"qkflow read <ISSUE-KEY> including all images, analyze the architecture design"
```

### 📊 Output Formats

#### Terminal Output (show command)

```
╔═══════════════════════════════════════════════════════╗
║ 🎫 NA-9245: Implement user authentication            ║
╚═══════════════════════════════════════════════════════╝

📋 Type:        Story
📊 Status:      In Progress
🏷️  Priority:    High
👤 Assignee:    John Doe

🔗 View in Jira: https://brain-ai.atlassian.net/browse/NA-9245
```

#### Markdown Output (export command)

```markdown
---
issue_key: NA-9245
title: Implement user authentication
type: Story
status: In Progress
priority: High
---

# NA-9245: Implement user authentication

## 📊 Metadata
...

## 📝 Description
...

## 📎 Attachments (3)
1. **screenshot.png** (245 KB)
   ![screenshot.png](./attachments/screenshot.png)
...
```

#### Cursor-Optimized Output (read command)

The `read` command provides special output markers that Cursor recognizes:

```
✅ Exported to: /tmp/qkflow/jira/NA-9245/

Main file: /tmp/qkflow/jira/NA-9245/content.md
Images: /tmp/qkflow/jira/NA-9245/attachments/ (3 files)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💡 CURSOR: Please read the following files:
1. /tmp/qkflow/jira/NA-9245/content.md
2. All images in /tmp/qkflow/jira/NA-9245/attachments/
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 🚀 Advanced Usage

#### Batch Processing

```bash
# Export multiple issues
for issue in NA-9245 NA-9246 NA-9247; do
  qkflow jira export $issue --with-images
done

# Tell Cursor to analyze all
"Read all Jira exports in /tmp/qkflow/jira/ and summarize"
```

#### Custom Workflows

```bash
# Create a script for your team
#!/bin/bash
ISSUE_KEY=$1
qkflow jira read "$ISSUE_KEY"
echo "Ready for Cursor to analyze!"
```

#### Integration with Other Tools

```bash
# Export and open in VS Code/Cursor
qkflow jira export NA-9245 --with-images
code /tmp/qkflow/jira/NA-9245/content.md
```

### 🐛 Troubleshooting

#### "Failed to create Jira client"
- Check your config: `qkflow config`
- Verify your API token is valid
- Ensure `jira_service_address` is correct

#### "Failed to get issue"
- Verify the issue key is correct (e.g., NA-9245)
- Check you have permission to view the issue
- Try accessing the issue in your browser first

#### "Failed to download attachment"
- Check your network connection
- Verify your API token has attachment download permissions
- Some files may be restricted

#### Cursor not reading files
- Make sure the export command completed successfully
- Check the file paths in the output
- Try manually attaching the file: `@/tmp/qkflow/jira/NA-9245/content.md`

### 📝 Notes

- Exports are temporary and stored in `/tmp/qkflow/jira/` by default
- Use `clean` command regularly to free up space
- Images are only downloaded with `--with-images` flag
- The `read` command is specifically designed for Cursor AI integration

---

## Part 2: Jira Status Configuration {#jira-status-configuration}

Configure Jira status mappings for automatic status updates during PR workflow.

### 📋 Configuration Storage Location

Jira status configurations are intelligently stored based on your system:

**macOS with iCloud Drive** (Recommended):
```
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```
Configurations automatically sync across all your Mac devices ☁️

**Local Storage** (Fallback):
```
~/.qkflow/jira-status.json
```

Run `qkflow config` to see your actual storage location.

**Configuration File Structure:**
```json
{
  "mappings": {
    "PROJ-123": {
      "project_key": "PROJ",
      "pr_created_status": "In Progress",
      "pr_merged_status": "Done"
    },
    "TEAM-456": {
      "project_key": "TEAM",
      "pr_created_status": "进行中",
      "pr_merged_status": "已完成"
    }
  }
}
```

### 🎯 Configuration Overview

#### 1. Basic Configuration (Required)

Configure basic Jira information in `config.yaml`:

```yaml
email: your.email@example.com
jira_api_token: your_jira_api_token
jira_service_address: https://your-domain.atlassian.net
github_token: ghp_your_github_token
```

Run `qkflow init` to initialize configuration.

**Tip**: If you use macOS with iCloud Drive enabled, all configurations automatically sync to your other Mac devices!

#### 2. Project Status Mapping (Per Project)

Each Jira project needs two status configurations:

- **PR Created Status**: When creating a PR, the Jira issue will be updated to this status
  - Common values: `In Progress`, `进行中`, `开发中`, etc.

- **PR Merged Status**: When a PR is merged, the Jira issue will be updated to this status
  - Common values: `Done`, `已完成`, `Resolved`, etc.

### 🛠️ How to Configure

#### Method 1: Auto-Configuration on First Use (Recommended)

When you first create a PR for a project, the system will automatically prompt you to configure status mappings:

```bash
# Create PR
qkflow pr create PROJ-123

# If this is the first time using this project, interactive configuration will appear:
# 1. Fetch all available statuses for this project from Jira
# 2. Let you select "PR Created" status (e.g., In Progress)
# 3. Let you select "PR Merged" status (e.g., Done)
# 4. Automatically save to ~/.qkflow/jira-status.json
```

#### Method 2: Manual Setup/Update Project Configuration

```bash
# Set status mapping for a specific project
qkflow jira setup PROJ

# The system will:
# 1. Connect to Jira and fetch all available statuses for this project
# 2. Display interactive selection interface
# 3. Save your selection
```

#### Method 3: View Configured Projects

```bash
# List all configured project status mappings
qkflow jira list

# Example output:
# 📋 Jira Status Mappings:
#
# Project: PROJ
#   PR Created → In Progress
#   PR Merged  → Done
#
# Project: TEAM
#   PR Created → 进行中
#   PR Merged  → 已完成
```

#### Method 4: Delete Project Configuration

```bash
# Delete status mapping for a specific project
qkflow jira delete PROJ

# Will require confirmation before deletion
```

#### Method 5: Manual Configuration File Editing

You can also edit the configuration file directly:

```bash
# Edit configuration
vim ~/.qkflow/jira-status.json

# Or
code ~/.qkflow/jira-status.json
```

### 🔄 Workflow

#### When Creating PR (`qkflow pr create`)

1. Check if project has status mapping
2. If not, automatically trigger configuration flow
3. Create PR
4. Update Jira issue to `PR Created Status` (e.g., In Progress)
5. Add PR link to Jira issue

#### When Merging PR (`qkflow pr merge`)

1. Read project status mapping
2. Merge PR
3. Update Jira issue to `PR Merged Status` (e.g., Done)
4. Add merge comment to Jira issue

### 📝 Examples

#### Complete Configuration Example

```bash
# 1. Initialize basic configuration
qkflow init

# 2. View current configuration
qkflow config

# 3. Set status mapping for project PROJ
qkflow jira setup PROJ
# Select PR Created: In Progress
# Select PR Merged: Done

# 4. View all status mappings
qkflow jira list

# 5. Create PR (will automatically use configured status)
qkflow pr create PROJ-123

# 6. Merge PR (will automatically use configured status)
qkflow pr merge 456
```

#### Multi-Project Configuration Example

If you work with multiple Jira projects:

```bash
# Configure project A
qkflow jira setup PROJA
# Select: In Progress / Done

# Configure project B (may use Chinese statuses)
qkflow jira setup PROJB
# Select: 进行中 / 已完成

# Configure project C (may use custom statuses)
qkflow jira setup PROJC
# Select: Development / Resolved

# View all configurations
qkflow jira list
```

### ⚙️ Technical Implementation

#### Status Retrieval

The system automatically fetches all available statuses for a project via Jira REST API:

```
GET /rest/api/2/project/{projectKey}/statuses
```

This ensures you can only select statuses actually supported by the project, avoiding configuration errors.

#### Caching Mechanism

- Configuration persists once saved, unless manually updated or deleted
- Configuration is automatically read for the corresponding project on each operation
- If configuration is deleted, you'll be prompted to reconfigure on next use

### 🔍 Troubleshooting

#### Issue 1: Status Configuration Not Found

```bash
# Check if configuration file exists
ls -la ~/.qkflow/jira-status.json

# If it doesn't exist, reconfigure
qkflow jira setup YOUR_PROJECT_KEY
```

#### Issue 2: Status Name Mismatch

If status names change in Jira:

```bash
# Reconfigure the project
qkflow jira setup YOUR_PROJECT_KEY

# Or manually edit configuration file
vim ~/.qkflow/jira-status.json
```

#### Issue 3: Cannot Fetch Project Statuses

Ensure:
1. Jira API Token is valid
2. You have access permission for the project
3. Jira Service Address is correct

```bash
# Check basic configuration
qkflow config

# Reinitialize configuration
qkflow init
```

### 🎨 Best Practices

1. **Configure on first use**: You'll be prompted when first creating a PR for a project, recommended to complete configuration then
2. **Unified naming**: If your team has multiple projects, try to use unified status names
3. **Regular checks**: Use `qkflow jira list` regularly to verify configurations are correct
4. **Backup configuration**: You can backup the `~/.qkflow/jira-status.json` file

### 📚 Related Commands

```bash
# Jira commands
qkflow jira list           # List all status mappings
qkflow jira setup [key]    # Set/update project status mapping
qkflow jira delete [key]   # Delete project status mapping

# Configuration commands
qkflow init               # Initialize configuration
qkflow config              # View current configuration

# Commands that use configuration
qkflow pr create [ticket] # Create PR (uses status configuration)
qkflow pr merge [number]  # Merge PR (uses status configuration)
```

### 🔗 Related Files

#### Configuration File Locations

**macOS with iCloud Drive**:
- Config directory: `~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/`
  - Basic config: `config.yaml`
  - Status mapping: `jira-status.json`

**Local Storage**:
- Config directory: `~/.qkflow/`
  - Basic config: `config.yaml`
  - Status mapping: `jira-status.json`

#### Source Code
- `internal/utils/paths.go` - Path management and iCloud detection
- `internal/jira/status_cache.go` - Status cache management
- `cmd/qkflow/commands/jira.go` - Jira commands
- `cmd/qkflow/commands/pr_create.go` - PR creation logic
- `cmd/qkflow/commands/pr_merge.go` - PR merge logic

---

## 🔧 Common Configuration

Both features require Jira credentials to be configured:

```bash
qkflow init
```

Required settings:
- `jira_service_address`: Your Jira instance URL (e.g., https://your-domain.atlassian.net)
- `jira_api_token`: Your Jira API token
- `email`: Your Jira email

## 🔗 Related Commands

- `qkflow init` - Configure Jira credentials
- `qkflow config` - View current configuration
- `qkflow pr create` - Create PR (can auto-link to Jira and update status)
- `qkflow pr merge` - Merge PR (can update Jira status)

---