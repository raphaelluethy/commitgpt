# CommitGPT

CommitGPT is a command-line tool that automatically generates meaningful Git commit messages using Claude AI (Anthropic's API). It analyzes your unstaged changes and creates a concise, descriptive commit message following conventional commit format.

## Features

- Automatically detects **unstaged** Git changes
- Generates semantic commit messages using AI
- Follows conventional commit format (e.g., `feat:`, `fix:`, etc.)
- Handles staging and commit creation automatically

## Prerequisites

Before using CommitGPT, ensure you have:

1. Go installed on your system
2. Git installed and configured
3. An Anthropic API key
4. A repository - lol

## Installation

1. Clone the repository:
```bash
git clone [your-repo-url]
cd commitgpt
```

2. Set up your Anthropic API key as an environment variable:
```bash
export ANTHROPIC_API_KEY='your-api-key-here'
```

3. Build the application:
```bash
go build -o commitgpt
```

4. (Optional) Add the binary to your PATH for system-wide access.

## Usage

Simply run CommitGPT in your Git repository when you have changes you want to commit:

```bash
./commitgpt
```

The tool will:
1. Check for unstaged changes
2. Generate a commit message using Claude AI
3. Stage all changes (`git add .`)
4. Create a commit with the generated message

### Options

- `--help` or `-h`: Display help information

## Example Output

```bash
$ ./commitgpt
Created commit: feat: add user authentication system with OAuth2 support
```

## Error Handling

The tool will exit with an error message if:
- No unstaged changes are found
- The Anthropic API call fails
- Git operations (staging/committing) fail

## Environment Variables

- `ANTHROPIC_API_KEY`: Your Anthropic API key (required)

## License

MIT License. See [LICENSE](LICENSE) for details.
