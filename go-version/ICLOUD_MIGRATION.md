# iCloud Drive 配置同步指南

## 🌟 新特性

从此版本开始，`qkflow` 在 macOS 上会优先将配置存储到 iCloud Drive，实现多设备自动同步！

## ☁️ 自动同步的内容

以下配置会自动同步到你的所有 Mac 设备：

1. **主配置文件** - GitHub Token、Jira 配置、AI 密钥等
2. **Jira 状态映射** - 每个项目的状态配置

## 📍 存储位置

### macOS with iCloud Drive (推荐)
```
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/config.yaml
~/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```

### 本地存储 (回退方案)
```
~/.qkflow/config.yaml
~/.qkflow/jira-status.json
```

## 🔄 迁移指南

### 自动迁移

如果你之前使用的是本地配置，可以手动迁移到 iCloud：

```bash
# 1. 确保 iCloud Drive 已启用
# 打开"系统设置" → "Apple ID" → "iCloud" → 确保"iCloud Drive"已开启

# 2. 迁移配置文件
# 如果有旧的分散配置，迁移到统一目录
if [ -f ~/.config/quick-workflow/config.yaml ]; then
  cp ~/.config/quick-workflow/config.yaml \
     ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/config.yaml
fi

if [ -f ~/.qkflow/jira-status.json ]; then
  cp ~/.qkflow/jira-status.json \
     ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/jira-status.json
fi

# 4. 验证迁移
qkflow config
```

### 验证同步状态

运行以下命令查看当前的存储位置：

```bash
qkflow config
```

输出示例：
```
💾 Storage:
  Location: iCloud Drive (synced across devices)
  Config: /Users/xxx/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/config.yaml
  Jira Status: /Users/xxx/Library/Mobile Documents/com~apple~CloudDocs/.qkflow/jira-status.json
```

## 🎯 多设备使用

### 新 Mac 设备设置

在新的 Mac 设备上：

1. **安装 qkflow**
   ```bash
   # 下载并安装
   curl -L https://github.com/yiminlab/quick-workflow/releases/latest/download/qkflow-darwin-arm64 -o qkflow
   chmod +x qkflow
   sudo mv qkflow /usr/local/bin/
   ```

2. **等待 iCloud 同步**
   - 打开 Finder → iCloud Drive
   - 确保 `.config` 和 `.qkflow` 文件夹已同步完成

3. **验证配置**
   ```bash
   qkflow config
   ```

配置会自动从 iCloud Drive 读取，无需重新配置！

## 🔒 安全说明

- iCloud Drive 存储是加密的
- 配置文件权限仍然是 `0644`（仅用户可读写）
- Token 和密钥安全地存储在你的 iCloud 账户中
- 只有登录同一 Apple ID 的设备才能访问

## ⚠️ 注意事项

### iCloud 同步延迟

iCloud Drive 同步可能需要几秒到几分钟，取决于：
- 网络连接速度
- 文件大小
- 系统负载

### 离线工作

如果没有网络连接：
- 仍然可以读取和修改配置
- 更改会在网络恢复后自动同步

### 回退到本地存储

如果你不想使用 iCloud Drive：

```bash
# 1. 禁用 iCloud Drive (系统设置)
# 2. qkflow 会自动回退到本地存储

# 或者手动移动配置回本地
cp ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/config.yaml \
   ~/.qkflow/config.yaml

cp ~/Library/Mobile\ Documents/com~apple~CloudDocs/.qkflow/jira-status.json \
   ~/.qkflow/jira-status.json
```

## 🐛 故障排除

### 问题 1: "No configuration found" 错误

**原因**: iCloud Drive 未启用或文件未同步完成

**解决方案**:
```bash
# 检查 iCloud Drive 是否可用
ls -la ~/Library/Mobile\ Documents/com~apple~CloudDocs/

# 如果目录不存在，启用 iCloud Drive
# 系统设置 → Apple ID → iCloud → iCloud Drive

# 手动创建配置
qkflow init
```

### 问题 2: 配置不同步

**原因**: iCloud 同步延迟或网络问题

**解决方案**:
```bash
# 1. 检查 iCloud 同步状态
# 打开 Finder → iCloud Drive → 检查文件是否有云图标

# 2. 强制同步
# 右键点击文件 → "从 iCloud 下载"

# 3. 检查网络连接
ping icloud.com
```

### 问题 3: 多设备配置冲突

**原因**: 同时在多台设备上修改配置

**解决方案**:
- iCloud 会自动处理冲突
- 如果出现问题，选择一个设备上的配置作为主配置
- 在其他设备上运行 `qkflow init` 重新初始化

## 📚 更多信息

- [主 README](README.md)
- [Jira 状态配置指南](JIRA_STATUS_CONFIG.md)
- [快速开始](QUICKSTART.md)

## 💡 提示

1. **建议启用 iCloud Drive**: 配置会自动同步到你的所有 Mac 设备
2. **定期备份**: 虽然 iCloud 很可靠，但定期备份配置文件仍然是好习惯
3. **团队使用**: 每个人都应该有自己的配置，不要共享 iCloud 账户

---

享受跨设备同步的便利吧！🎉

