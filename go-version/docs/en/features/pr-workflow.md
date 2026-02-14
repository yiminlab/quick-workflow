# PR Workflow Guide

Complete guide for managing pull requests with `qkflow` - from creation to approval to merge.

## 📋 Table of Contents

- [Overview](#overview)
- [Complete PR Lifecycle](#complete-pr-lifecycle)
- [1. Creating a PR](#1-creating-a-pr)
- [2. Approving a PR](#2-approving-a-pr)
- [3. Merging a PR](#3-merging-a-pr)
- [Common Workflows](#common-workflows)
- [URL Support](#url-support)
- [Related Documentation](#related-documentation)

## 🎯 Overview

`qkflow` provides three main PR commands that cover the entire pull request lifecycle:

| Command | Purpose | Key Features |
|---------|---------|--------------|
| `qkflow pr create` | Create a new PR | Web editor, Jira integration, auto-branching |
| `qkflow pr approve` | Approve a PR | URL support, custom comments, auto-merge |
| `qkflow pr merge` | Merge a PR | Auto-cleanup, Jira status update |

All commands support:
- ✅ **PR numbers** (e.g., `123`)
- ✅ **GitHub URLs** (e.g., `https://github.com/owner/repo/pull/123`)
- ✅ **Auto-detection** from current branch

## 🔄 Complete PR Lifecycle

Here's the typical flow from start to finish:

```
┌─────────────────────────────────────────────────────────────┐
│  1. CREATE PR                                               │
│  qkflow pr create PROJ-123                                  │
│  ├─ Create branch                                           │
│  ├─ Commit changes                                          │
│  ├─ Push to remote                                          │
│  ├─ [Optional] Add description with web editor              │
│  ├─ Create GitHub PR                                        │
│  ├─ Upload files to GitHub & Jira                          │
│  └─ Update Jira status                                     │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  2. CODE REVIEW                                             │
│  (Manual review process)                                     │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  3. APPROVE PR                                              │
│  qkflow pr approve 123 -c "LGTM!"                           │
│  ├─ Approve on GitHub                                       │
│  ├─ Add comment                                             │
│  └─ [Optional] Auto-merge with --merge flag                 │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  4. MERGE PR                                                │
│  qkflow pr merge 123                                        │
│  ├─ Merge on GitHub                                         │
│  ├─ Delete remote branch                                    │
│  ├─ Switch to main branch                                   │
│  ├─ Delete local branch                                     │
│  └─ Update Jira status to Done/Merged                      │
└─────────────────────────────────────────────────────────────┘
```

## 1. Creating a PR

### Quick Start

```bash
# Basic: Create PR with Jira ticket
qkflow pr create PROJ-123

# Without Jira ticket (press Enter when prompted)
qkflow pr create
```

### What Happens

1. ✅ Gets Jira ticket details (if provided)
2. ✅ Prompts for change types (feat, fix, etc.)
3. ✅ **[Optional]** Opens web editor for description & screenshots
4. ✅ Generates PR title
5. ✅ Creates git branch
6. ✅ Commits staged changes
7. ✅ Pushes to remote
8. ✅ Creates GitHub PR
9. ✅ Uploads files and adds comment to GitHub & Jira
10. ✅ Updates Jira status
11. ✅ Copies PR URL to clipboard

### Web Editor Feature

When creating a PR, you'll be prompted:

```
? Add detailed description with images/videos?
  > ⏭️  Skip (default)
    ✅ Yes, continue
```

If you choose "Yes, continue":
- Opens a beautiful web-based markdown editor in your browser
- Drag & drop images/videos
- Paste images from clipboard
- Real-time preview
- Auto-uploads to GitHub and Jira

**Example:**

```bash
$ qkflow pr create NA-9245

✓ Found Jira issue: Fix login button styling
📝 Select type(s) of changes:
  ✓ 🐛 Bug fix

? Add detailed description with images/videos?
  > ⏭️  Skip (default)
    ✅ Yes, continue

# Select "Yes, continue"
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

### Supported File Types

- **Images**: PNG, JPG, JPEG, GIF, WebP, SVG
- **Videos**: MP4, MOV, WebM, AVI

### Detailed Documentation

For complete details on the PR editor feature, see:
- **[PR Editor Feature](pr-editor.md)** - Technical implementation and advanced usage

## 2. Approving a PR

### Quick Start

```bash
# Approve by PR number (default 👍 comment)
qkflow pr approve 123

# Approve by URL (works from anywhere!)
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# Custom comment
qkflow pr approve 123 -c "LGTM! 🎉"

# Approve and auto-merge
qkflow pr approve 123 -c "Ship it! 🚀" -m

# Auto-detect from current branch
qkflow pr approve
```

### What Happens

1. ✅ Detects PR number or URL
2. ✅ Fetches PR details
3. ✅ Approves PR on GitHub
4. ✅ Adds comment (default: 👍, or custom with `-c`)
5. ✅ **[Optional]** Auto-merges if `--merge` flag is used

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--comment` | `-c` | Add a custom comment with approval |
| `--merge` | `-m` | Automatically merge after approval |
| `--help` | `-h` | Show help information |

### Examples

**Basic Approval:**
```bash
$ qkflow pr approve 123
✅ PR approved with comment: 👍
```

**Custom Comment:**
```bash
$ qkflow pr approve 123 -c "Great work! All tests passed."
✅ PR approved with comment: Great work! All tests passed.
```

**Approve and Merge:**
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

### Detailed Documentation

For complete details on approval features, see:
- **[PR Approve Guide](pr-approve.md)** - Complete guide with URL support, workflows, and use cases

## 3. Merging a PR

### Quick Start

```bash
# Merge by PR number
qkflow pr merge 123

# Merge by URL (works from anywhere!)
qkflow pr merge https://github.com/brain/planning-api/pull/2001

# Auto-detect from current branch
qkflow pr merge
```

### What Happens

1. ✅ Detects PR number or URL
2. ✅ Fetches PR details
3. ✅ Confirms merge with you
4. ✅ Merges PR on GitHub
5. ✅ Deletes remote branch
6. ✅ Switches to main branch
7. ✅ Pulls latest changes
8. ✅ Deletes local branch
9. ✅ Updates Jira status to Done/Merged
10. ✅ Adds merge comment to Jira

### Example

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

## 🌐 URL Support

All three PR commands support GitHub URLs, making it easy to work across repositories:

### Supported URL Formats

```bash
# HTTPS (most common)
https://github.com/owner/repo/pull/123

# With /files suffix (Files tab)
https://github.com/owner/repo/pull/123/files

# With /commits suffix (Commits tab)
https://github.com/owner/repo/pull/123/commits

# With /checks suffix (Checks tab)
https://github.com/owner/repo/pull/123/checks

# HTTP
http://github.com/owner/repo/pull/123

# Without protocol
github.com/owner/repo/pull/123

# With query params (parsed correctly)
https://github.com/owner/repo/pull/123?comments=all

# With fragments (parsed correctly)
https://github.com/owner/repo/pull/123#discussion_r123456
```

### Why Use URLs?

1. **Cross-Repository**: Work with PRs from different repos
2. **No Directory Change**: Work from anywhere
3. **Browser to CLI**: Copy URL from GitHub directly
4. **Batch Operations**: Script operations across multiple repos

### Examples

```bash
# Approve PR from different repo
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"

# Merge PR from shared link
qkflow pr merge https://github.com/company/backend/pull/200

# Copy URL from browser and paste
qkflow pr approve "$(pbpaste)" -m  # macOS
```

## 📋 Common Workflows

### Workflow 1: Complete Feature Development

```bash
# 1. Start feature
git checkout -b feature/new-feature
# ... make changes ...
git add .

# 2. Create PR
qkflow pr create PROJ-456
# → Select change types
# → [Optional] Add description with screenshots
# → PR created, Jira updated

# 3. Wait for review...

# 4. Reviewer approves
qkflow pr approve https://github.com/org/repo/pull/789 -c "LGTM!"

# 5. Merge PR
qkflow pr merge 789
# → PR merged, branches cleaned up, Jira updated to Done
```

### Workflow 2: Quick Bug Fix

```bash
# 1. Create hotfix branch
git checkout -b hotfix/critical-bug
# ... fix the bug ...
git add .

# 2. Create PR
qkflow pr create PROJ-789

# 3. Self-approve and merge immediately
qkflow pr approve 999 -c "Critical hotfix - merging immediately" -m
# → Approved and merged in one step!
```

### Workflow 3: Cross-Repository Review

```bash
# Review PRs from multiple repos without changing directories

# Frontend PR
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"

# Backend PR
qkflow pr approve https://github.com/team/backend/pull/200 -c "Looks good!"

# Mobile PR
qkflow pr approve https://github.com/team/mobile/pull/300 -c "LGTM!"
```

### Workflow 4: Browser to CLI

1. Open PR in GitHub web interface
2. Copy URL from address bar
3. Paste into terminal:

```bash
# Approve
qkflow pr approve "$(pbpaste)" -c "Reviewed and approved"

# Or merge
qkflow pr merge "$(pbpaste)"
```

### Workflow 5: Batch Operations

```bash
#!/bin/bash
# Approve multiple PRs
PR_URLS=(
  "https://github.com/org/repo1/pull/10"
  "https://github.com/org/repo2/pull/20"
  "https://github.com/org/repo3/pull/30"
)

for url in "${PR_URLS[@]}"; do
  qkflow pr approve "$url" -c "Auto-approved"
done
```

## 🎯 Pro Tips

### Tip 1: Use Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
alias approve='qkflow pr approve'
alias merge='qkflow pr merge'
alias prc='qkflow pr create'
alias gha='qkflow pr approve'  # GitHub Approve shortcut
```

### Tip 2: Combine with GitHub CLI

```bash
# List PRs with gh, approve with qkflow
gh pr list
qkflow pr approve https://github.com/owner/repo/pull/123 -m
```

### Tip 3: Comment Templates

```bash
# In your shell config
export APPROVE_LGTM="Looks good to me! 👍"
export APPROVE_EXCELLENT="Excellent work! 🎉"

# Usage:
qkflow pr approve 123 -c "$APPROVE_LGTM"
```

### Tip 4: Check Before Merge

```bash
# View PR details
gh pr view 123

# Check CI status
gh pr checks 123

# Then merge if all good
qkflow pr merge 123
```

## ⚠️ Error Handling

### PR Not Found

```bash
$ qkflow pr approve 999
❌ Failed to get PR: Pull request not found
```

**Solution:** Check PR number with `gh pr list`

### Invalid URL or Number

```bash
$ qkflow pr approve invalid-url
❌ Invalid PR number or URL: invalid-url
ℹ️  Expected: PR number (e.g., '123') or GitHub URL
```

**Solution:** Use a valid PR number or GitHub URL

### PR Already Closed/Merged

```bash
$ qkflow pr approve 123
❌ PR is not open (state: closed)
```

**Solution:** PR is already merged or closed

### Not Mergeable

```bash
$ qkflow pr approve 123 -m
✅ PR approved!
⚠️ Cannot merge PR: PR has conflicts and cannot be merged
```

**Solution:** Resolve conflicts first, then merge

## 🔗 Related Commands

- `qkflow pr create` - Create a new PR
- `qkflow pr approve` - Approve a PR
- `qkflow pr merge` - Merge a PR
- `gh pr view 123` - View PR details (GitHub CLI)
- `gh pr checks 123` - Check CI status (GitHub CLI)
- `gh pr list` - List all PRs (GitHub CLI)

## 📚 Related Documentation

### Detailed Feature Guides

- **[PR Approve Guide](pr-approve.md)** - Complete guide for `qkflow pr approve`
  - URL support details
  - Advanced workflows
  - Use cases and examples
  - Error handling

- **[PR Editor Feature](pr-editor.md)** - Technical details for PR creation editor
  - Web editor implementation
  - File upload mechanics
  - Technical architecture

### Other Related Docs

- **[Jira Integration](jira-integration.md)** - Complete Jira integration guide (Issue Reader & Status Configuration)
- **[Development Overview](../development/overview.md)** - Technical architecture ([中文](../../cn/development/overview.md))

## 🆚 Comparison with GitHub CLI

### GitHub CLI (`gh`)

```bash
# Create PR
gh pr create --title "Title" --body "Body"

# Approve
gh pr review 123 --approve --body "LGTM"

# Merge
gh pr merge 123

# Cleanup (manual)
git checkout main
git pull
git branch -D feature-branch
```

### qkflow

```bash
# Create PR (with Jira integration, web editor, auto-branching)
qkflow pr create PROJ-123

# Approve (with URL support, auto-merge option)
qkflow pr approve 123 -c "LGTM" -m

# Merge (with auto-cleanup, Jira update)
qkflow pr merge 123
```

**Benefits:**
- ✅ Fewer commands
- ✅ Auto-cleanup
- ✅ Jira integration
- ✅ Interactive prompts
- ✅ Branch auto-detection
- ✅ URL support
- ✅ Web-based editor

---

**Need Help?**

```bash
qkflow pr --help
qkflow pr create --help
qkflow pr approve --help
qkflow pr merge --help
```

**Found a Bug?**

Open an issue on GitHub with details!

---