# PR Approve Command Guide

Quick guide for using the new `qkflow pr approve` command.

> 📚 **Looking for the complete PR workflow?** See [PR Workflow Guide](pr-workflow.md) for the full lifecycle from creation to merge.

## 🚀 Quick Start

### Basic Approval

```bash
# Approve a specific PR by number (uses default 👍 comment)
qkflow pr approve 123

# Approve a PR by URL (works from anywhere!)
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# Also works with /files, /commits, /checks URLs
qkflow pr approve https://github.com/brain/planning-api/pull/2001/files

# Auto-detect PR from current branch
qkflow pr approve
```

### With Comment

By default, all approvals use 👍 as the comment. Customize with `-c` flag:

```bash
# Default approval (with 👍)
qkflow pr approve 123

# Custom comment
qkflow pr approve 123 -c "LGTM! 🎉"

# Approve by URL with custom comment
qkflow pr approve https://github.com/owner/repo/pull/456 -c "Great work!"

# Long comment with flag
qkflow pr approve 123 --comment "Great work! All tests passed. Approved for merge."
```

### Auto-Merge

```bash
# Approve and merge in one step
qkflow pr approve 123 --merge

# Approve by URL and merge
qkflow pr approve https://github.com/owner/repo/pull/789 --merge

# Short flag
qkflow pr approve 123 -m

# With comment and merge
qkflow pr approve 123 -c "LGTM!" -m

# URL with comment and merge
qkflow pr approve https://github.com/owner/repo/pull/789 -c "Ship it! 🚀" -m
```

## 🌐 URL Support (NEW!)

Now you can approve PRs from **any repository** without being in the git directory!

### Why Use URLs?

1. **Cross-Repository**: Approve PRs from different repos
2. **No Directory Change**: Work from anywhere
3. **Browser to CLI**: Copy URL from GitHub directly
4. **Batch Operations**: Script approvals across multiple repos

### URL Examples

```bash
# Approve a PR from a different repo
qkflow pr approve https://github.com/brain/planning-api/pull/2001

# Your colleague shares a PR link, approve it instantly
qkflow pr approve https://github.com/company/frontend/pull/456 -c "Looks good!"

# Merge someone else's PR by URL
qkflow pr merge https://github.com/team/backend/pull/789

# Approve and merge with URL
qkflow pr approve https://github.com/org/project/pull/123 -c "LGTM! 🎉" -m
```

### Supported URL Formats

All these formats work:

```bash
# HTTPS (most common)
https://github.com/owner/repo/pull/123

# With /files suffix (Files tab)
https://github.com/brain/planning-api/pull/2001/files

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
https://github.com/owner/repo/pull/123/files?file-filters%5B%5D=.js

# With fragments (parsed correctly)
https://github.com/owner/repo/pull/123#discussion_r123456
```

**Pro Tip:** Just copy the URL from any PR tab (Overview, Files, Commits, Checks) and it will work! 🎉

### URL Parsing Details

The tool automatically:
- Detects if argument is a URL or number
- Extracts owner, repo, and PR number from URL
- Handles query parameters and fragments
- Validates URL format

### When to Use URL vs Number

| Scenario | Use | Example |
|----------|-----|---------|
| Your repo, in git dir | Number | `qkflow pr approve 123` |
| Your repo, current branch | None | `qkflow pr approve` |
| Different repo | URL | `qkflow pr approve https://...` |
| Shared link from browser | URL | Copy & paste URL |
| Scripting multiple repos | URL | Loop through URLs |

## 📋 Common Workflows

### Workflow 1: Quick Code Review

You're reviewing a colleague's PR:

```bash
# 1. Check out their branch (optional)
git fetch origin
git checkout feature-branch

# 2. Review the code...

# 3. Approve
qkflow pr approve
# Finds PR automatically from branch
# Adds optional comment
# Done!
```

### Workflow 2: Approve and Merge

You have approval rights and want to merge immediately:

```bash
# One command to approve and merge
qkflow pr approve 123 -c "Approved! Merging now." --merge

# What happens:
# ✅ Approves PR #123
# ✅ Adds comment
# ✅ Checks if mergeable
# ✅ Confirms with you
# ✅ Merges PR
# ✅ Deletes remote branch
# ✅ Switches to main
# ✅ Deletes local branch
```

### Workflow 3: Batch Approvals

Multiple PRs to review:

```bash
# List all open PRs first
gh pr list

# Approve them one by one
qkflow pr approve 101 -c "Approved"
qkflow pr approve 102 -c "Approved"
qkflow pr approve 103 -c "Approved"
```

### Workflow 4: Interactive Mode

Don't know the PR number? Let the tool help:

```bash
# Run without arguments
qkflow pr approve

# What happens:
# 1. Tries to find PR from current branch
# 2. If not found, shows list of all open PRs
# 3. You select one
# 4. Asks for optional comment
# 5. Approves!
```

## 🔍 Use Cases

### Use Case 1: Cross-Repository Reviews

Review PRs from multiple repositories without changing directories:

```bash
qkflow pr approve https://github.com/team/frontend/pull/100 -c "Approved"
qkflow pr approve https://github.com/team/backend/pull/200 -c "Approved"
qkflow pr approve https://github.com/team/mobile/pull/300 -c "Approved"
```

### Use Case 2: Browser to CLI Workflow

1. Open PR in GitHub web
2. Copy URL from address bar
3. Paste into terminal:

```bash
qkflow pr approve https://github.com/company/project/pull/1234 -c "Looks good!"
```

### Use Case 3: Slack/Email Integration

Teammate shares a PR link in Slack? Approve it instantly:

```bash
# Copy link from Slack
qkflow pr approve <paste-url-here> -c "Reviewed and approved"
```

### Use Case 4: Scripting Across Repos

Automate approvals across multiple repositories:

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

### Use Case 5: Team Lead Approval

As a team lead, you need to approve PRs daily:

```bash
# Morning routine: approve all ready PRs
for pr in 121 122 123; do
  qkflow pr approve $pr -c "Reviewed and approved"
done
```

### Use Case 6: CI/CD Integration

Add to your CI pipeline:

```bash
#!/bin/bash
# Auto-approve dependabot PRs after tests pass
if [[ "$PR_AUTHOR" == "dependabot" ]] && [[ "$TESTS_PASSED" == "true" ]]; then
  qkflow pr approve $PR_NUMBER -c "Auto-approved: Tests passed" -m
fi
```

### Use Case 7: Hotfix Workflow

Fast-track urgent fixes:

```bash
# Create hotfix
git checkout -b hotfix/critical-bug
# ... fix the bug ...
git add .
qkflow pr create

# Get it approved and merged ASAP
qkflow pr approve 999 -c "Critical hotfix - merging immediately" -m
```

## 🎯 Pro Tips

### Tip 1: Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
# Quick approve
alias approve='qkflow pr approve'
alias merge='qkflow pr approve --merge'
alias gha='qkflow pr approve'  # GitHub Approve shortcut

# Usage:
approve 123 -c "LGTM"
merge 123
gha https://github.com/owner/repo/pull/123 -c "LGTM!"
```

### Tip 1b: Use with `pbpaste` (macOS)

```bash
# Copy URL in browser, then:
qkflow pr approve "$(pbpaste)" -c "Approved!"
```

### Tip 1c: Integration with GitHub CLI

Combine with `gh` CLI:

```bash
# List PRs with gh, approve with qkflow
gh pr list
qkflow pr approve https://github.com/owner/repo/pull/123 -m
```

### Tip 2: Comment Templates

Save common comments:

```bash
# In your shell config
export APPROVE_LGTM="Looks good to me! 👍"
export APPROVE_MINOR="Approved with minor comments. Please address in follow-up."
export APPROVE_EXCELLENT="Excellent work! 🎉"

# Usage:
qkflow pr approve 123 -c "$APPROVE_LGTM"
```

### Tip 3: Check Before Merge

Before using `--merge`, verify the PR:

```bash
# View PR details
gh pr view 123

# Check CI status
gh pr checks 123

# Approve and merge if all good
qkflow pr approve 123 -m
```

### Tip 4: Branch Protection

If branch protection is enabled:

```bash
# Just approve - let GitHub merge rules handle the rest
qkflow pr approve 123 -c "Approved"

# Don't use --merge if you need multiple approvals
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
ℹ️  Expected: PR number (e.g., '123') or GitHub URL (e.g., 'https://github.com/owner/repo/pull/123')
```

**Solution:** Use a valid PR number or GitHub URL

### PR Already Closed

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

### No PR for Branch

```bash
$ qkflow pr approve
⚠️ No PR found for branch: feature-xyz
Do you want to select a PR from the list? (Y/n)
```

**Solution:** Either select from list or specify PR number

## 🆚 Comparison

### GitHub CLI (`gh`)

```bash
# Approve
gh pr review 123 --approve --body "LGTM"

# Then merge separately
gh pr merge 123

# Then cleanup
git checkout main
git pull
git branch -D feature-branch
```

### qkflow (New!)

```bash
# All in one!
qkflow pr approve 123 -c "LGTM" -m
```

**Benefits:**
- ✅ Fewer commands
- ✅ Auto-cleanup
- ✅ Interactive prompts
- ✅ Branch auto-detection
- ✅ Merge validation

## 🔗 Related Commands

- `qkflow pr create` - Create a new PR
- `qkflow pr merge` - Merge without approving first
- `gh pr view 123` - View PR details
- `gh pr checks 123` - Check CI status

## 📚 Full Reference

### Command Syntax

```
qkflow pr approve [pr-number] [flags]
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--comment` | `-c` | Add a comment with the approval |
| `--merge` | `-m` | Automatically merge after approval |
| `--help` | `-h` | Show help information |

### Exit Codes

- `0`: Success
- `1`: Error (PR not found, not mergeable, etc.)

### Environment

Requires:
- GitHub token configured (`qkflow config`)
- Git repository with remote
- Write access to repository

## 📝 Examples in Action

### Example 1: Quick Review (Default Comment)

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

### Example 2: Quick Review (Custom Comment)

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

### Example 3: Approve and Merge

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

## 🤝 Backwards Compatibility

All existing workflows still work:

```bash
# PR number (requires being in repo)
qkflow pr approve 123

# Auto-detect from current branch
qkflow pr approve

# Interactive selection
qkflow pr approve
# → Shows list of PRs to choose from
```

---

**Need Help?**

```bash
qkflow pr approve --help
qkflow help
```

**Found a Bug?**

Open an issue on GitHub with details!

---