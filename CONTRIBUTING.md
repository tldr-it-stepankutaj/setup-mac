# Contributing to setup-mac

Thank you for your interest in contributing to setup-mac!

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/setup-mac.git
   cd setup-mac
   ```
3. Set up development environment:
   ```bash
   make deps
   make install-tools
   ```
4. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run Linter

```bash
make lint
```

### Format Code

```bash
make fmt
```

### Test Your Changes

```bash
./bin/setup-mac install --all --dry-run
```

## Pull Request Process

1. Ensure your code passes all tests and linting
2. Update documentation if needed
3. Add tests for new functionality
4. Create a pull request with a clear description

### PR Title Format

Use conventional commit format:

- `feat: add new feature`
- `fix: resolve bug`
- `docs: update documentation`
- `refactor: improve code structure`
- `test: add tests`
- `chore: maintenance tasks`

### PR Description

- Describe what changes you made
- Explain why these changes are needed
- Reference related issues (e.g., `Fixes #123`)

## Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Keep functions small and focused
- Add comments for non-obvious code
- Use meaningful variable and function names

## Adding New Installers

1. Create a new file in `internal/installer/`
2. Implement the `Installer` interface:
   ```go
   type Installer interface {
       Name() string
       Description() string
       IsInstalled(ctx context.Context) bool
       Install(ctx context.Context) error
   }
   ```
3. Register the installer in `internal/installer/installer.go`
4. Add configuration schema in `internal/config/schema.go`
5. Update default config in `internal/config/defaults.yaml`

## Reporting Issues

- Check existing issues before creating a new one
- Use issue templates when available
- Provide clear reproduction steps
- Include system information (macOS version, Go version)

## Questions?

Feel free to open an issue for questions or discussions.
