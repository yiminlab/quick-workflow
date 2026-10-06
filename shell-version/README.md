# quick-workflow

> The shell command based on jira-cli and github-cli for quick work

## Two quick commands for streamlined workflow:

1. `pr-create`: Automatically create a Git branch, submit a PR (Pull Request), capture PR details, link the PR to Jira, and update its status.
2. `pr-merge`: Automatically merge a PR, delete the Git branch, and update the Jira status in one swift command.

### Why use it

> Are you tired of the hassle of changing the status every time you create a PR, such as 'code review,' 'merged,' and dealing with repetitive tasks?

> Especially when you're busy and have opened multiple browser windows, waiting for pages to load just to complete these mundane tasks can be time-consuming, diverting our focus from more important matters.

Highlighted Benefits:

- Efficiency: Streamline the process by automating status updates and repetitive tasks.
- Time-Saving: Eliminate the need to open multiple browser tabs and wait for pages to load.
- Enhanced Focus: Allow users to concentrate on more crucial tasks.
- Reduced Manual Work: Minimize the need for manual status updates for each PR.
- Standardized Workflow: Ensures a consistent and error-free process.

### Installation

1. Install required dependencies:

```shell
# Install github-cli
brew install gh
gh auth login

# Install jira-cli
brew tap ankitpokhrel/jira-cli
brew install jira-cli

# Install jq for JSON processing
brew install jq
```

2. Clone the repository:
```shell
gh repo clone yiminlab/quick-workflow
```

3. Navigate to the shell-version directory and run the installation script:
```shell
cd quick-workflow/shell-version
chmod +x install.sh
./install.sh
```

4. After installation, update your environment variables in your shell's rc file (`.zshrc` or `.bashrc`):
```shell
export EMAIL=your.email@example.com
export JIRA_API_TOKEN=your_jira_api_token
export JIRA_SERVICE_ADDRESS=your_jira_service_address
export GH_BRANCH_PREFIX=your_branch_prefix
export OPENAI_KEY=your_openai_api_key
export DEEPSEEK_KEY=your_deepseek_api_key
export OPENAI_PROXY_URL=your_openai_proxy_url
export OPENAI_PROXY_KEY=your_openai_proxy_key
```

5. Configure Jira CLI:
   - [Get a Jira API token](https://id.atlassian.com/manage-profile/security/api-tokens)
   - Run `jira init` and follow the prompts to complete the setup

6. Source your shell's rc file or restart your terminal:
```shell
source ~/.zshrc  # For Zsh users
# or
source ~/.bashrc  # For Bash users
```

### Available Commands

After installation, you'll have access to the following commands:

- `pr-create`: Create PR and update Jira status
- `pr-merge`: Merge PR and update Jira status
- `qkupdate`: Quick updates
- `proxy`: Check and set proxy settings
- `qklogs`: Quick access to logs
- `qkfind`: Quick find utility
- `qksearch`: Quick search utility
- `qk`: Quick workflow utility

### Python Quick Run Utility

- `qkpy.sh` or `qkpy`: Conveniently run Python scripts, automatically setting `PYTHONPATH` to the current directory, suitable for cases where you need to reference local packages.

**Usage Example:**
```shell
./qkpy.sh src/case/planning/travel_extractor.py
# or
qkpy src/case/planning/travel_extractor.py
```

### Updates

See `todolist.md` for the latest updates and other important information.

Thank you for your valuable feedback and PRs. Let's work together to make these commands better!
