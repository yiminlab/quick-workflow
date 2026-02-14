# iCloud Drive Configuration Sync Guide

## 🌟 New Feature

Starting from this version, `qkflow` will prioritize storing configurations in iCloud Drive on macOS, enabling automatic sync across devices!

## ☁️ What Gets Synced

The following configurations will automatically sync to all your Mac devices:

1. **Main Configuration File** - GitHub Token, Jira configuration, AI keys, etc.
2. **Jira Status Mappings** - Status configurations for each project

## 📍 Storage Locations

### macOS with iCloud Drive (Recommended)
```
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/config.yaml
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```

### Local Storage (Fallback)
```
~/.qkflow/config.yaml
~/.qkflow/jira-status.json
```

## 🔄 Migration Guide

### Automatic Migration

If you were previously using local configuration, you can manually migrate to iCloud:

```bash
# 1. Ensure iCloud Drive is enabled
# Open "System Settings" → "Apple ID" → "iCloud" → Ensure "iCloud Drive" is enabled

# 2. Migrate configuration files
# If you have old scattered configs, migrate to unified directory
if [ -f ~/.config/quick-workflow/config.yaml ]; then
  cp ~/.config/quick-workflow/config.yaml \
     ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/config.yaml
fi

if [ -f ~/.qkflow/jira-status.json ]; then
  cp ~/.qkflow/jira-status.json \
     ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/jira-status.json
fi

# 4. Verify migration
qkflow config
```

### Verify Sync Status

Run the following command to see the current storage location:

```bash
qkflow config
```

Example output:
```
💾 Storage:
  Location: iCloud Drive (synced across devices)
  Config: /Users/xxx/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/config.yaml
  Jira Status: /Users/xxx/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```

## 🎯 Multi-Device Usage

### Setting Up on a New Mac Device

On a new Mac device:

1. **Install qkflow**
   ```bash
   # Download and install
   curl -L https://github.com/Wangggym/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
   chmod +x qkflow
   sudo mv qkflow /usr/local/bin/
   ```

2. **Wait for iCloud Sync**
   - Open Finder → iCloud Drive
   - Ensure `.config` and `.qkflow` folders have finished syncing

3. **Verify Configuration**
   ```bash
   qkflow config
   ```

Configuration will automatically be read from iCloud Drive - no need to reconfigure!

## 🔒 Security Notes

- iCloud Drive storage is encrypted
- Configuration file permissions remain `0644` (user read/write only)
- Tokens and keys are securely stored in your iCloud account
- Only devices logged in with the same Apple ID can access

## ⚠️ Important Notes

### iCloud Sync Delay

iCloud Drive sync may take a few seconds to a few minutes, depending on:
- Network connection speed
- File size
- System load

### Offline Work

If there's no network connection:
- You can still read and modify configurations
- Changes will automatically sync when network is restored

### Fallback to Local Storage

If you don't want to use iCloud Drive:

```bash
# 1. Disable iCloud Drive (System Settings)
# 2. qkflow will automatically fallback to local storage

# Or manually move configs back to local
cp ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/config.yaml \
   ~/.qkflow/config.yaml

cp ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/jira-status.json \
   ~/.qkflow/jira-status.json
```

## 🐛 Troubleshooting

### Issue 1: "No configuration found" Error

**Cause**: iCloud Drive not enabled or files not synced yet

**Solution**:
```bash
# Check if iCloud Drive is available
ls -la ~/Library/Mobile\ Documents/com~apple~CloudDocs/

# If directory doesn't exist, enable iCloud Drive
# System Settings → Apple ID → iCloud → iCloud Drive

# Manually create configuration
qkflow init
```

### Issue 2: Configuration Not Syncing

**Cause**: iCloud sync delay or network issues

**Solution**:
```bash
# 1. Check iCloud sync status
# Open Finder → iCloud Drive → Check if files have cloud icon

# 2. Force sync
# Right-click file → "Download from iCloud"

# 3. Check network connection
ping icloud.com
```

### Issue 3: Multi-Device Configuration Conflicts

**Cause**: Modifying configuration on multiple devices simultaneously

**Solution**:
- iCloud will automatically handle conflicts
- If issues occur, choose one device's configuration as primary
- Run `qkflow init` on other devices to reinitialize

## 📚 More Information

- [Main README](../../../README.md)
- [Jira Integration Guide](jira-integration.md) - Complete Jira integration guide (includes status configuration)
- [Quick Start](../../../QUICKSTART.md)

## 💡 Tips

1. **Enable iCloud Drive (Recommended)**: Configurations will automatically sync to all your Mac devices
2. **Regular Backups**: Although iCloud is reliable, regularly backing up configuration files is still a good practice
3. **Team Usage**: Everyone should have their own configuration - don't share iCloud accounts

---

**Enjoy the convenience of cross-device sync!** 🎉

---