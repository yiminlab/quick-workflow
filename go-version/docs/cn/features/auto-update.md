# 自动更新功能

## 概述

qkflow 内置了智能的自动更新功能，可以确保你始终使用最新版本，享受最新特性和 bug 修复。

## 特性

- ✅ **自动检查更新** - 每24小时自动检查一次新版本
- ✅ **后台静默运行** - 不影响当前命令的执行
- ✅ **智能更新** - 发现新版本时自动下载并安装
- ✅ **可配置** - 可以选择启用或禁用自动更新
- ✅ **手动更新** - 随时可以手动触发更新
- ✅ **安全备份** - 更新失败时自动恢复旧版本

## 工作原理

### 1. 自动检查

每次运行 qkflow 命令时（除了 `init`、`version`、`help`、`update-cli`），系统会：

1. 检查距离上次检查更新是否已超过 24 小时
2. 如果是，则在后台检查 GitHub Releases 的最新版本
3. 比较当前版本和最新版本

### 2. 自动更新（默认开启）

如果发现新版本且启用了自动更新：

```bash
🎉 New version available: v1.1.0 (current: v1.0.0)
⬇️  Downloading update...
✅ Successfully updated to version v1.1.0! Please restart qkflow.
```

系统会：
1. 从 GitHub Releases 下载对应平台的二进制文件
2. 备份当前版本
3. 替换为新版本
4. 提示重启 qkflow

### 3. 手动更新模式

如果禁用了自动更新，只会提示：

```bash
🎉 New version available: v1.1.0 (current: v1.0.0)
Run 'qkflow update-cli' to update, or visit: https://github.com/Wangggym/quick-workflow/releases
```

## 配置自动更新

### 初始化时配置

运行 `qkflow init` 时会询问：

```bash
Enable automatic updates? (recommended) [Y/n]:
```

- 按 `Y` 或 直接回车：启用自动更新（推荐）
- 按 `n`：禁用自动更新

### 修改配置

编辑配置文件 `~/.qkflow/config.yaml`（或 iCloud Drive 路径）：

```yaml
# 启用自动更新（推荐）
auto_update: true

# 或禁用自动更新
auto_update: false
```

修改后，下次运行 qkflow 命令时生效。

### 使用环境变量

```bash
# 临时禁用自动更新
export AUTO_UPDATE=false
qkflow pr create

# 临时启用自动更新
export AUTO_UPDATE=true
qkflow pr create
```

## 手动更新

随时可以手动检查并更新到最新版本：

```bash
qkflow update-cli
```

这会：
1. 立即检查最新版本
2. 如果有新版本，下载并安装
3. 如果已是最新版本，显示确认信息

示例输出：

```bash
# 有新版本
🎉 New version available: v1.1.0 (current: v1.0.0)
⬇️  Downloading update...
✅ Successfully updated to version v1.1.0! Please restart qkflow.

# 已是最新版本
✅ You are already running the latest version (v1.0.0)
```

## 更新频率

- **检查频率**: 每 24 小时检查一次
- **时间戳文件**: `~/.qkflow/.last_update_check`
- **超时设置**: API 请求 5 秒超时，下载 5 分钟超时

可以手动删除时间戳文件来强制立即检查：

```bash
rm ~/.qkflow/.last_update_check
```

## 支持的平台

自动更新支持所有 qkflow 支持的平台：

- ✅ macOS (Intel) - `qkflow-darwin-amd64`
- ✅ macOS (Apple Silicon) - `qkflow-darwin-arm64`
- ✅ Linux (amd64) - `qkflow-linux-amd64`
- ✅ Windows (amd64) - `qkflow-windows-amd64.exe`

## 安全性

### 下载来源

- 所有二进制文件从 GitHub Releases 官方下载
- URL: `https://api.github.com/repos/Wangggym/quick-workflow/releases/latest`
- 使用 HTTPS 加密传输

### 更新流程

1. 下载到临时文件
2. 设置可执行权限
3. 备份当前版本 (`qkflow.backup`)
4. 替换为新版本
5. 删除备份文件
6. 如果失败，自动恢复备份

### 权限要求

更新需要对二进制文件所在目录有写权限：

- 如果安装在 `/usr/local/bin/`，需要 sudo 权限
- 如果安装在 `~/bin/` 或 `GOPATH/bin`，不需要特殊权限

## 常见问题

### Q: 自动更新会中断我的工作吗？

A: 不会。更新检查在后台运行，不阻塞当前命令。只有在下载完成后才会提示重启。

### Q: 更新失败怎么办？

A: 系统会自动恢复备份。你可以手动运行 `qkflow update-cli` 重试，或访问 GitHub Releases 手动下载。

### Q: 如何查看当前版本？

```bash
qkflow version
```

### Q: 如何禁用自动更新？

编辑 `~/.qkflow/config.yaml`：

```yaml
auto_update: false
```

或使用环境变量：

```bash
export AUTO_UPDATE=false
```

### Q: 能否回滚到旧版本？

可以。访问 [GitHub Releases](https://github.com/Wangggym/quick-workflow/releases) 下载特定版本的二进制文件，手动替换即可。

### Q: 自动更新失败会显示错误吗？

A: 不会。为了不影响用户体验，更新失败时静默处理。可以手动运行 `qkflow update-cli` 查看详细错误。

## 最佳实践

1. **推荐启用自动更新** - 确保始终使用最新版本，享受最新特性和 bug 修复

2. **定期手动检查** - 如果禁用了自动更新，建议定期运行：
   ```bash
   qkflow update-cli
   ```

3. **在重要操作前更新** - 在执行重要操作前，确保使用最新版本：
   ```bash
   qkflow update-cli
   qkflow pr create
   ```

4. **CI/CD 环境** - 在 CI/CD 中，建议禁用自动更新，使用固定版本：
   ```bash
   export AUTO_UPDATE=false
   ```

## 更新日志

查看完整的更新历史和版本说明：

👉 [GitHub Releases](https://github.com/Wangggym/quick-workflow/releases)

## 反馈

如果在使用自动更新功能时遇到问题，请：

- 🐛 [报告 Bug](https://github.com/Wangggym/quick-workflow/issues/new?labels=bug)
- 💡 [提出建议](https://github.com/Wangggym/quick-workflow/issues/new?labels=enhancement)

---

**享受无缝的自动更新体验！** 🚀