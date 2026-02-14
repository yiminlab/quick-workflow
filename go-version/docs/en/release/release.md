# Release Guide

## 🚀 Quick Start

### First Release

```bash
# 1. Ensure code is pushed
git push origin main

# 2. Run release check (runs tests and build)
make release VERSION=v1.0.0

# 3. If check passes, create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. Wait 2-3 minutes, GitHub Actions will automatically build and release
```

### Regular Release

```bash
# Patch version (bug fixes)
make release VERSION=v1.0.1
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1

# Minor version (new features)
make release VERSION=v1.1.0
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# Major version (breaking changes)
make release VERSION=v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0
```

### Quick Checklist

```bash
# ✅ All changes committed
git status

# ✅ Run release check
make release VERSION=vX.Y.Z

# ✅ Create tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# ✅ Push tag
git push origin vX.Y.Z

# ✅ Check GitHub Actions
# https://github.com/Wangggym/quick-workflow/actions

# ✅ Verify Release
# https://github.com/Wangggym/quick-workflow/releases
```

---

## 📦 Detailed Release Process

### 1. Prepare for Release

Ensure all changes are committed and pushed to main branch:

```bash
# Ensure on main branch
git checkout main
git pull origin main

# Check current status
git status

# If there are uncommitted changes, commit them first
git add .
git commit -m "chore: prepare for release v1.0.0"
git push origin main
```

### 2. Create and Push Tag

```bash
# Create tag (follow semantic versioning)
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag to remote repository (this triggers GitHub Actions)
git push origin v1.0.0
```

**Version Number Convention** (Semantic Versioning SemVer):
- `v1.0.0` - Major.Minor.Patch
- `v1.0.0-beta.1` - Pre-release version
- `v1.0.0-rc.1` - Release candidate

### 3. Automatic Build and Release

After pushing the tag, GitHub Actions will automatically:

1. ✅ Run tests
2. ✅ Code linting
3. ✅ Build multi-platform binaries:
   - `qkflow-darwin-amd64` (macOS Intel)
   - `qkflow-darwin-arm64` (macOS Apple Silicon)
   - `qkflow-linux-amd64` (Linux)
   - `qkflow-windows-amd64.exe` (Windows)
4. ✅ Create GitHub Release
5. ✅ Upload all binary files to Release

### 4. Verify Release

1. Visit GitHub Releases page:
   ```
   https://github.com/Wangggym/quick-workflow/releases
   ```

2. Confirm the following:
   - ✅ Release title and description are correct
   - ✅ All platform binaries are uploaded
   - ✅ File sizes are reasonable (typically 10-20MB)
   - ✅ Download links are available

3. Test downloaded binary locally:
   ```bash
   # Download and test
   curl -L https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-darwin-arm64 -o qkflow
   chmod +x qkflow
   ./qkflow version
   ```

## 🔄 Version Management

### Version Number Increment Rules

- **MAJOR**: Incompatible API changes
  - Example: `v1.0.0` → `v2.0.0`
  - Examples: Major refactoring, removing old features

- **MINOR**: Backward-compatible feature additions
  - Example: `v1.0.0` → `v1.1.0`
  - Examples: Adding new commands, new features

- **PATCH**: Backward-compatible bug fixes
  - Example: `v1.0.0` → `v1.0.1`
  - Examples: Bug fixes, documentation updates

### Pre-release Versions

Use pre-release versions when testing new features:

```bash
# Beta version
git tag -a v1.1.0-beta.1 -m "Beta release for testing"
git push origin v1.1.0-beta.1

# Release Candidate
git tag -a v1.1.0-rc.1 -m "Release candidate"
git push origin v1.1.0-rc.1
```

## 📝 Release Notes Template

When creating a tag, it's recommended to use detailed release notes:

```bash
git tag -a v1.0.0 -m "Release v1.0.0

## 🚀 New Features
- Add iCloud Drive sync support for configs
- Implement qkflow update command
- Add jira status mapping management

## 🐛 Bug Fixes
- Fix config path resolution on macOS
- Improve error handling in PR creation

## 📚 Documentation
- Add iCloud migration guide
- Update README with new features

## 🔧 Changes
- Rename binary from qkg to qkflow
- Unify config directory structure
"
```

Or edit manually on the GitHub Release page.

## 🛠️ Manual Release (If Auto Fails)

If GitHub Actions fails, you can manually build and release:

```bash
# 1. Build for all platforms
make build-all

# 2. Verify build artifacts
ls -lh bin/

# 3. Create Release on GitHub web
# - Visit: https://github.com/Wangggym/quick-workflow/releases/new
# - Select tag: v1.0.0
# - Fill in title and description
# - Upload all files from bin/ directory
# - Click "Publish release"
```

## 🔍 Check CI/CD Status

View GitHub Actions execution status:

1. Visit Actions page:
   ```
   https://github.com/Wangggym/quick-workflow/actions
   ```

2. View build logs:
   - Click on the latest workflow run
   - View detailed logs for each job
   - Fix errors based on logs if any

## 🚨 Rollback Release

If critical issues are found and rollback is needed:

```bash
# 1. Delete Release on GitHub (mark as draft or delete)

# 2. Delete remote tag
git push --delete origin v1.0.0

# 3. Delete local tag
git tag -d v1.0.0

# 4. Fix issues and re-release
git tag -a v1.0.1 -m "Fix critical bug in v1.0.0"
git push origin v1.0.1
```

## 📊 Release Checklist

Pre-release checklist:

```markdown
## Pre-release Checklist
- [ ] All tests pass (`make test`)
- [ ] Code linting passes (`make lint`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version number follows convention
- [ ] Tested build locally (`make build-all`)
- [ ] Version number in README updated

## Release Checklist
- [ ] Create git tag
- [ ] Push tag to remote
- [ ] GitHub Actions build succeeds
- [ ] All binary files uploaded
- [ ] Release notes filled
- [ ] Download links available

## Post-release Checklist
- [ ] Announce on social media (if needed)
- [ ] Update documentation website (if needed)
- [ ] Notify team members
- [ ] Close related issues
```

## 🎯 Quick Command Reference

```bash
# View all tags
git tag -l

# View tag details
git show v1.0.0

# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push --delete origin v1.0.0

# Push all tags
git push origin --tags

# Get latest tag
git describe --tags --abbrev=0

# Create tag based on specific commit
git tag -a v1.0.0 <commit-hash> -m "Release v1.0.0"
```

## 🛠️ Common Commands

```bash
# Initialize dependencies
make gen

# Run tests
make test

# Build local version
make build

# Install to system
make install

# Clean build artifacts
make clean

# View help
make help
```

## 📦 Release Artifacts

Each release generates the following files:
- `qkflow-darwin-amd64` - macOS Intel
- `qkflow-darwin-arm64` - macOS Apple Silicon
- `qkflow-linux-amd64` - Linux
- `qkflow-windows-amd64.exe` - Windows

## 🔍 View Build Status

- Actions: https://github.com/Wangggym/quick-workflow/actions
- Releases: https://github.com/Wangggym/quick-workflow/releases

## 💡 Best Practices

1. **Regular Releases**: Recommend releasing every 2-4 weeks
2. **Semantic Versioning**: Strictly follow SemVer specification
3. **Detailed Release Notes**: Help users understand changes
4. **Thorough Testing**: Test thoroughly locally before release
5. **Backup**: Keep old version binaries
6. **Notify Users**: Notify users via GitHub, documentation, etc.

## 📞 Need Help?

- Check [GitHub Actions Documentation](https://docs.github.com/en/actions)
- Check build logs to troubleshoot
- Contact maintainer for help

---

**Remember**: Every tag push triggers automatic build and release, please be careful!

