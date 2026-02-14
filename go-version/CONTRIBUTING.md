# Contributing to Quick Workflow

Thank you for your interest in contributing to Quick Workflow! 🎉

## 🚀 Quick Start

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/quick-workflow.git`
3. Create a branch: `git checkout -b feature/amazing-feature`
4. Make your changes
5. Test your changes: `make test`
6. Commit: `git commit -m "Add amazing feature"`
7. Push: `git push origin feature/amazing-feature`
8. Open a Pull Request

## 🏗️ Development Setup

### Prerequisites

- Go 1.21 or higher
- Make
- Git
- GitHub CLI (`gh`) for testing
- Jira account for integration testing (optional)

### Local Setup

```bash
# Clone the repository
git clone https://github.com/Wangggym/quick-workflow.git
cd quick-workflow/go-version

# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Run linters
make lint
```

## 📝 Code Style

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format code (run `make fmt`)
- Keep functions small and focused
- Write meaningful variable names
- Add comments for exported functions and types

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

Examples:
```
feat: add support for GitLab integration
fix: handle empty Jira ticket gracefully
docs: update installation instructions
test: add tests for GitHub client
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/github/...
go test ./internal/jira/...

# Run with coverage
make coverage

# Run with race detector
go test -race ./...
```

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests for multiple test cases
- Mock external dependencies (GitHub API, Jira API)
- Test both success and error cases

Example:
```go
func TestSanitizeBranchName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"simple", "hello world", "hello-world"},
        {"special chars", "hello@world!", "hello-world"},
        {"multiple dashes", "hello--world", "hello-world"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := SanitizeBranchName(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

## 📚 Documentation

- Update README.md for user-facing changes
- Update docs/en/migration/migration.md and docs/cn/migration/migration.md for migration-related changes
- Add godoc comments for exported types and functions
- Update CHANGELOG.md for all changes (follow [Keep a Changelog](https://keepachangelog.com/) format)

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Description**: Clear description of the issue
2. **Steps to Reproduce**: Detailed steps
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happens
5. **Environment**:
   - OS and version
   - Go version
   - qk version (`qk version`)
6. **Logs**: Any relevant error messages

Use the bug report template when creating an issue.

## 💡 Feature Requests

When requesting features:

1. **Use Case**: Describe the problem you're trying to solve
2. **Proposed Solution**: Your idea for solving it
3. **Alternatives**: Other solutions you've considered
4. **Additional Context**: Any other relevant information

## 🔍 Code Review

All submissions require review. We use GitHub Pull Requests for this purpose.

### PR Checklist

Before submitting a PR, ensure:

- [ ] Tests pass (`make test`)
- [ ] Linters pass (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] PR description explains the changes
- [ ] No sensitive information (tokens, passwords) in code

### PR Review Process

1. Automated checks run (tests, linters)
2. Maintainer reviews code
3. Requested changes are addressed
4. PR is approved and merged

## 🏷️ Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Documentation improvements
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention needed
- `question`: Further information requested

## 📋 Project Structure

```
go-version/
├── cmd/qkflow/              # Main application entry point
│   ├── main.go
│   └── commands/        # CLI commands
├── internal/            # Internal packages
│   ├── github/          # GitHub client
│   ├── jira/            # Jira client
│   ├── git/             # Git operations
│   └── ui/              # User interface
├── pkg/                 # Public packages
│   └── config/          # Configuration management
├── scripts/             # Build and release scripts
└── docs/                # Additional documentation
```

## 🎯 Areas for Contribution

We welcome contributions in these areas:

### High Priority
- [ ] Windows support improvements
- [ ] GitLab support
- [ ] Bitbucket support
- [ ] More comprehensive tests
- [ ] Performance optimizations

### Medium Priority
- [ ] Configuration validation
- [ ] Better error messages
- [ ] Template support for PR bodies
- [ ] Custom workflows
- [ ] Integration with other issue trackers

### Documentation
- [ ] Video tutorials
- [ ] More examples
- [ ] Troubleshooting guide
- [ ] API documentation

## 🤝 Community

- **Discussions**: Use GitHub Discussions for questions and ideas
- **Issues**: Report bugs and request features
- **PRs**: Submit code contributions
- **Code of Conduct**: Be respectful and inclusive

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You!

Your contributions make Quick Workflow better for everyone. Thank you for taking the time to contribute!

---

If you have questions, feel free to open an issue or start a discussion.

