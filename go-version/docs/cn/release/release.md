# 发布指南

## 🚀 快速入门

### 第一次发布

```bash
# 1. 确保代码已推送
git push origin main

# 2. 运行发布检查（运行测试和构建）
make release VERSION=v1.0.0

# 3. 如果检查通过，创建并推送 tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. 等待 2-3 分钟，GitHub Actions 将自动构建和发布
```

### 日常发布

```bash
# 补丁版本（bug 修复）
make release VERSION=v1.0.1
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1

# 次版本（新功能）
make release VERSION=v1.1.0
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 主版本（破坏性更改）
make release VERSION=v2.0.0
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0
```

### 快速检查清单

```bash
# ✅ 所有更改已提交
git status

# ✅ 运行发布检查
make release VERSION=vX.Y.Z

# ✅ 创建 tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# ✅ 推送 tag
git push origin vX.Y.Z

# ✅ 检查 GitHub Actions
# https://github.com/Wangggym/quick-workflow/actions

# ✅ 验证 Release
# https://github.com/Wangggym/quick-workflow/releases
```

---

## 📦 详细发布流程

### 1. 准备发布

确保所有更改已提交并推送到 main 分支：

```bash
# 确保在 main 分支
git checkout main
git pull origin main

# 检查当前状态
git status

# 如果有未提交的更改，请先提交
git add .
git commit -m "chore: prepare for release v1.0.0"
git push origin main
```

### 2. 创建并推送 Tag

```bash
# 创建 tag（遵循语义化版本）
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送 tag 到远程仓库（这会触发 GitHub Actions）
git push origin v1.0.0
```

**版本号约定**（语义化版本 SemVer）：
- `v1.0.0` - 主版本.次版本.修订号
- `v1.0.0-beta.1` - 预发布版本
- `v1.0.0-rc.1` - 发布候选版本

### 3. 自动构建和发布

推送 tag 后，GitHub Actions 将自动：

1. ✅ 运行测试
2. ✅ 代码检查
3. ✅ 构建多平台二进制文件：
   - `qkflow-darwin-amd64` (macOS Intel)
   - `qkflow-darwin-arm64` (macOS Apple Silicon)
   - `qkflow-linux-amd64` (Linux)
   - `qkflow-windows-amd64.exe` (Windows)
4. ✅ 创建 GitHub Release
5. ✅ 上传所有二进制文件到 Release

### 4. 验证 Release

1. 访问 GitHub Releases 页面：
   ```
   https://github.com/Wangggym/quick-workflow/releases
   ```

2. 确认以下内容：
   - ✅ Release 标题和描述正确
   - ✅ 所有平台二进制文件已上传
   - ✅ 文件大小合理（通常为 10-20MB）
   - ✅ 下载链接可用

3. 在本地测试下载的二进制文件：
   ```bash
   # 下载并测试
   curl -L https://github.com/Wangggym/quick-workflow/releases/download/v1.0.0/qkflow-darwin-arm64 -o qkflow
   chmod +x qkflow
   ./qkflow version
   ```

## 🔄 版本管理

### 版本号递增规则

- **主版本 (MAJOR)**: 不兼容的 API 更改
  - 示例：`v1.0.0` → `v2.0.0`
  - 示例：重大重构、移除旧功能

- **次版本 (MINOR)**: 向后兼容的功能添加
  - 示例：`v1.0.0` → `v1.1.0`
  - 示例：添加新命令、新功能

- **修订号 (PATCH)**: 向后兼容的 bug 修复
  - 示例：`v1.0.0` → `v1.0.1`
  - 示例：bug 修复、文档更新

### 预发布版本

在测试新功能时使用预发布版本：

```bash
# Beta 版本
git tag -a v1.1.0-beta.1 -m "Beta release for testing"
git push origin v1.1.0-beta.1

# 发布候选版本
git tag -a v1.1.0-rc.1 -m "Release candidate"
git push origin v1.1.0-rc.1
```

## 📝 Release Notes 模板

创建 tag 时，建议使用详细的 release notes：

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

或在 GitHub Release 页面上手动编辑。

## 🛠️ 手动 Release（如果自动失败）

如果 GitHub Actions 失败，可以手动构建和发布：

```bash
# 1. 为所有平台构建
make build-all

# 2. 验证构建产物
ls -lh bin/

# 3. 在 GitHub 网页上创建 Release
# - 访问：https://github.com/Wangggym/quick-workflow/releases/new
# - 选择 tag：v1.0.0
# - 填写标题和描述
# - 上传 bin/ 目录中的所有文件
# - 点击 "Publish release"
```

## 🔍 检查 CI/CD 状态

查看 GitHub Actions 执行状态：

1. 访问 Actions 页面：
   ```
   https://github.com/Wangggym/quick-workflow/actions
   ```

2. 查看构建日志：
   - 点击最新的工作流运行
   - 查看每个作业的详细日志
   - 如有错误，根据日志修复

## 🚨 回滚 Release

如果发现严重问题需要回滚：

```bash
# 1. 在 GitHub 上删除 Release（标记为草稿或删除）

# 2. 删除远程 tag
git push --delete origin v1.0.0

# 3. 删除本地 tag
git tag -d v1.0.0

# 4. 修复问题并重新发布
git tag -a v1.0.1 -m "Fix critical bug in v1.0.0"
git push origin v1.0.1
```

## 📊 发布检查清单

发布前检查清单：

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

## 🎯 快速命令参考

```bash
# 查看所有 tag
git tag -l

# 查看 tag 详情
git show v1.0.0

# 删除本地 tag
git tag -d v1.0.0

# 删除远程 tag
git push --delete origin v1.0.0

# 推送所有 tag
git push origin --tags

# 获取最新 tag
git describe --tags --abbrev=0

# 基于特定提交创建 tag
git tag -a v1.0.0 <commit-hash> -m "Release v1.0.0"
```

## 🛠️ 常用命令

```bash
# 初始化依赖
make gen

# 运行测试
make test

# 构建本地版本
make build

# 安装到系统
make install

# 清理构建产物
make clean

# 查看帮助
make help
```

## 📦 发布后的产物

每次发布会生成以下文件：
- `qkflow-darwin-amd64` - macOS Intel
- `qkflow-darwin-arm64` - macOS Apple Silicon
- `qkflow-linux-amd64` - Linux
- `qkflow-windows-amd64.exe` - Windows

## 🔍 查看构建状态

- Actions: https://github.com/Wangggym/quick-workflow/actions
- Releases: https://github.com/Wangggym/quick-workflow/releases

## 💡 最佳实践

1. **定期发布**: 建议每 2-4 周发布一次
2. **语义化版本**: 严格遵循 SemVer 规范
3. **详细的 Release Notes**: 帮助用户了解更改
4. **充分测试**: 发布前在本地充分测试
5. **备份**: 保留旧版本的二进制文件
6. **通知用户**: 通过 GitHub、文档等方式通知用户

## 📞 遇到问题？

- 查看 [GitHub Actions 文档](https://docs.github.com/en/actions)
- 查看构建日志进行故障排除
- 联系维护者寻求帮助

---

**记住**: 每次推送 tag 都会触发自动构建和发布，请小心！

