# Auto Update Feature

## Overview

qkflow includes intelligent auto-update functionality to ensure you're always using the latest version with the newest features and bug fixes.

## Features

- ✅ **Automatic Update Check** - Checks for new versions every 24 hours
- ✅ **Background Operation** - Runs silently without interrupting current commands
- ✅ **Smart Update** - Automatically downloads and installs when new version is found
- ✅ **Configurable** - Can be enabled or disabled
- ✅ **Manual Update** - Can trigger updates manually at any time
- ✅ **Safe Backup** - Automatically restores old version if update fails

## How It Works

### 1. Automatic Check

Every time you run a qkflow command (except `init`, `version`, `help`, `update-cli`), the system will:

1. Check if 24 hours have passed since last update check
2. If yes, check GitHub Releases for the latest version in the background
3. Compare current version with latest version

### 2. Automatic Update (Enabled by Default)

If a new version is found and auto-update is enabled:

```bash
🎉 New version available: v1.1.0 (current: v1.0.0)
⬇️  Downloading update...
✅ Successfully updated to version v1.1.0! Please restart qkflow.
```

The system will:
1. Download the binary for your platform from GitHub Releases
2. Backup current version
3. Replace with new version
4. Prompt to restart qkflow

### 3. Manual Update Mode

If auto-update is disabled, it will only show a prompt:

```bash
🎉 New version available: v1.1.0 (current: v1.0.0)
Run 'qkflow update-cli' to update, or visit: https://github.com/Wangggym/quick-workflow/releases
```

## Configuring Auto Update

### During Initialization

When running `qkflow init`, you'll be asked:

```bash
Enable automatic updates? (recommended) [Y/n]:
```

- Press `Y` or Enter: Enable auto-update (recommended)
- Press `n`: Disable auto-update

### Modify Configuration

Edit the config file `~/.qkflow/config.yaml` (or iCloud Drive path):

```yaml
# Enable auto-update (recommended)
auto_update: true

# Or disable auto-update
auto_update: false
```

Changes take effect on the next qkflow command run.

### Using Environment Variables

```bash
# Temporarily disable auto-update
export AUTO_UPDATE=false
qkflow pr create

# Temporarily enable auto-update
export AUTO_UPDATE=true
qkflow pr create
```

## Manual Update

You can manually check and update to the latest version at any time:

```bash
qkflow update-cli
```

This will:
1. Immediately check for the latest version
2. Download and install if a new version is available
3. Show confirmation if already on the latest version

Example output:

```bash
# New version available
🎉 New version available: v1.1.0 (current: v1.0.0)
⬇️  Downloading update...
✅ Successfully updated to version v1.1.0! Please restart qkflow.

# Already on latest version
✅ You are already running the latest version (v1.0.0)
```

## Update Frequency

- **Check Frequency**: Every 24 hours
- **Timestamp File**: `~/.qkflow/.last_update_check`
- **Timeout Settings**: 5 seconds for API requests, 5 minutes for downloads

You can manually delete the timestamp file to force an immediate check:

```bash
rm ~/.qkflow/.last_update_check
```

## Supported Platforms

Auto-update supports all platforms that qkflow supports:

- ✅ macOS (Intel) - `qkflow-darwin-amd64`
- ✅ macOS (Apple Silicon) - `qkflow-darwin-arm64`
- ✅ Linux (amd64) - `qkflow-linux-amd64`
- ✅ Windows (amd64) - `qkflow-windows-amd64.exe`

## Security

### Download Source

- All binaries are downloaded from official GitHub Releases
- URL: `https://api.github.com/repos/Wangggym/quick-workflow/releases/latest`
- Uses HTTPS encrypted transmission

### Update Process

1. Download to temporary file
2. Set executable permissions
3. Backup current version (`qkflow.backup`)
4. Replace with new version
5. Delete backup file
6. If failed, automatically restore backup

### Permission Requirements

Update requires write permissions to the directory where the binary is installed:

- If installed in `/usr/local/bin/`, requires sudo permissions
- If installed in `~/bin/` or `GOPATH/bin`, no special permissions needed

## FAQ

### Q: Will auto-update interrupt my work?

A: No. Update checks run in the background and don't block current commands. You'll only be prompted to restart after download completes.

### Q: What if update fails?

A: The system will automatically restore the backup. You can manually run `qkflow update-cli` to retry, or visit GitHub Releases to download manually.

### Q: How to check current version?

```bash
qkflow version
```

### Q: How to disable auto-update?

Edit `~/.qkflow/config.yaml`:

```yaml
auto_update: false
```

Or use environment variable:

```bash
export AUTO_UPDATE=false
```

### Q: Can I rollback to an older version?

Yes. Visit [GitHub Releases](https://github.com/Wangggym/quick-workflow/releases) to download a specific version binary and manually replace it.

### Q: Will update failures show errors?

A: No. To avoid disrupting user experience, update failures are handled silently. You can manually run `qkflow update-cli` to see detailed errors.

## Best Practices

1. **Enable Auto-Update (Recommended)** - Ensures you always use the latest version with newest features and bug fixes

2. **Regular Manual Checks** - If auto-update is disabled, regularly run:
   ```bash
   qkflow update-cli
   ```

3. **Update Before Important Operations** - Before important operations, ensure you're using the latest version:
   ```bash
   qkflow update-cli
   qkflow pr create
   ```

4. **CI/CD Environments** - In CI/CD, disable auto-update and use fixed versions:
   ```bash
   export AUTO_UPDATE=false
   ```

## Release Notes

View complete update history and version notes:

👉 [GitHub Releases](https://github.com/Wangggym/quick-workflow/releases)

## Feedback

If you encounter issues with auto-update:

- 🐛 [Report Bug](https://github.com/Wangggym/quick-workflow/issues/new?labels=bug)
- 💡 [Suggest Feature](https://github.com/Wangggym/quick-workflow/issues/new?labels=enhancement)

---

**Enjoy seamless auto-update experience!** 🚀

---