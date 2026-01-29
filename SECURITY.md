# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability within setup-mac, please send an email to [stepan.kutaj@tldr-it.com](mailto:stepan.kutaj@tldr-it.com).

Please do not report security vulnerabilities through public GitHub issues.

### What to include

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Any suggested fixes (optional)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Depending on complexity, typically within 30 days

### Process

1. Your report will be acknowledged within 48 hours
2. We will investigate and validate the issue
3. We will work on a fix and coordinate disclosure
4. A security advisory will be published after the fix is released

## Security Best Practices

When using setup-mac:

- Always review the configuration file before running
- Use `--dry-run` flag to preview changes before applying
- Keep the tool updated to the latest version
- Review shell scripts and commands that will be executed
