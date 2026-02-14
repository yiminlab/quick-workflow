# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Documentation reorganization: All documentation moved to `docs/` directory
- Unified CHANGELOG.md in root directory
- CONTRIBUTING.md moved to root directory for better GitHub integration
- **Documentation Structure Refactor** - Split bilingual documents into separate language directories
  - Created `docs/en/` and `docs/cn/` directories for English and Chinese documentation
  - Split 6 feature documents from bilingual format to separate language files
  - Improved file organization: smaller files (~50% reduction), clearer structure
  - Better maintainability: languages can be updated independently
  - See [Documentation Index](docs/README.md) for new structure

### Changed
- Documentation structure: Migrated from bilingual single-file format to separate `en/` and `cn/` directories
  - All feature documentation now in `docs/en/features/` and `docs/cn/features/`
  - Updated all cross-references in README.md, CHANGELOG.md, and QUICKSTART.md
  - Old bilingual files backed up to `docs/features.old.bilingual/`

## [1.4.0] - 2024-11-18

### Added
- **PR Approve Command** - New `qkflow pr approve` command with URL support
  - Approve PRs by number or GitHub URL
  - Default 👍 comment for quick approvals
  - Custom comment support with `-c` flag
  - Auto-merge option with `--merge` flag
  - Cross-repository support via URLs
  - Works with `/files`, `/commits`, `/checks` URL suffixes
  - See [PR Approve Guide](docs/en/features/pr-approve.md) ([中文](docs/cn/features/pr-approve.md)) for details

- **PR URL Support** - Enhanced `pr approve` and `pr merge` commands
  - Support for GitHub PR URLs in addition to PR numbers
  - Copy-paste friendly from browser
  - Works from any directory
  - See [PR Approve Guide](docs/en/features/pr-approve.md) ([中文](docs/cn/features/pr-approve.md)) for details

- **PR Editor Interaction Improvement** - Better UX for PR description editor
  - Changed from Yes/No prompt to visual selection interface
  - Skip option clearly marked as default
  - More intuitive interaction with icons
  - See [PR Editor Feature](docs/en/features/pr-editor.md) ([中文](docs/cn/features/pr-editor.md)) for details

### Changed
- Improved PR editor prompt UX (visual selection instead of y/n)

## [1.3.0] - 2024-11-18

### Added
- **Jira Issue Reader** - Comprehensive Jira issue reading and exporting
  - `qkflow jira show` - Quick terminal view
  - `qkflow jira export` - Complete export to files with Markdown
  - `qkflow jira read` - Intelligent mode (recommended for Cursor AI)
  - `qkflow jira clean` - Cleanup utility for exports
  - Image and attachment support with `--with-images` flag
  - Cursor AI optimized output
  - See [Jira Integration Guide](docs/en/features/jira-integration.md) ([中文](docs/cn/features/jira-integration.md)) for details

### Technical Details
- New packages: `internal/jira/formatter.go`, `internal/jira/exporter.go`, `internal/jira/cleaner.go`
- Extended `internal/jira/client.go` with detailed issue fetching
- New commands: `jira_show.go`, `jira_export.go`, `jira_read.go`, `jira_clean.go`
- ~1,200 lines of code added

## [1.2.0] - 2024-11-18

### Added
- **PR Editor Feature** - Web-based editor for PR descriptions
  - Beautiful GitHub-style web editor
  - Markdown editor with live preview
  - Drag & drop images and videos
  - Paste images from clipboard
  - Automatic upload to GitHub and Jira
  - See [PR Editor Feature](docs/en/features/pr-editor.md) ([中文](docs/cn/features/pr-editor.md)) for details

### Technical Details
- New package: `internal/editor/` with `server.go`, `html.go`, `uploader.go`
- Enhanced `internal/github/client.go` with `AddPRComment()` method
- Enhanced `internal/jira/client.go` with `AddAttachment()` method

## [1.1.0] - 2024-11-18

### Added
- **Auto Update** - Automatic update checking and installation
  - Checks for updates every 24 hours
  - Automatic download and installation
  - Manual update with `qkflow update-cli`
  - Configurable via `auto_update` setting
  - See [Auto Update Guide](docs/en/features/auto-update.md) ([中文](docs/cn/features/auto-update.md)) for details

- **iCloud Drive Sync** - Automatic config sync on macOS
  - Configs automatically synced across Mac devices
  - Stored in iCloud Drive when available
  - Fallback to local storage
  - See [iCloud Migration Guide](docs/en/features/icloud-migration.md) ([中文](docs/cn/features/icloud-migration.md)) for details

- **Watch Daemon** - Automatic PR monitoring
  - Monitors PRs every 15 minutes
  - Auto-updates Jira when PR is merged
  - Desktop notifications (macOS)
  - Auto-start on login
  - See README.md for details

### Changed
- Configuration storage: Now uses iCloud Drive on macOS (when available)
- Update mechanism: Automatic update checking added

## [1.0.0] - 2024-11-04

### Added
- Initial release of Go version
- PR creation with automatic branch management
- PR merging with cleanup
- Jira integration (status updates, comments, links)
- GitHub API integration (PR CRUD operations)
- Git operations (branch, commit, push, merge)
- Interactive CLI with beautiful prompts
- Configuration management with `qkflow init`
- Cross-platform support (macOS, Linux, Windows)
- Single binary distribution

### Changed
- Complete rewrite from Shell version to Go
- Improved performance (15x faster startup)
- Better error handling
- Type-safe code

---

## Release Notes Format

Each release includes:
- **Added** - New features
- **Changed** - Changes in existing functionality
- **Deprecated** - Soon-to-be removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Security fixes

## Links

- [Full Documentation](docs/README.md)
- [Migration Guide](docs/en/migration/migration.md) ([中文](docs/cn/migration/migration.md))
- [Contributing Guide](CONTRIBUTING.md)
- [GitHub Releases](https://github.com/Wangggym/quick-workflow/releases)

