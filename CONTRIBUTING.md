# Contributing to go-thalex

Thank you for considering contributing to go-thalex! This document outlines the development workflow and guidelines.

## Getting Started

### Prerequisites

- Go 1.21 or later
- [golangci-lint](https://golangci-lint.run/welcome/install-locally/)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) (`go install golang.org/x/tools/cmd/goimports@latest`)
- Make

### Setup

```bash
git clone https://github.com/amiwrpremium/go-thalex.git
cd go-thalex
go mod download
```

### Development Workflow

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Run all checks (format + vet + lint + test)
make check

# View coverage report
make coverage
```

## Making Changes

### Branch Naming

- `feature/description` — new features
- `fix/description` — bug fixes
- `docs/description` — documentation changes
- `refactor/description` — code refactoring
- `test/description` — test additions or improvements

### Commit Messages

Write clear, concise commit messages:

- Use the imperative mood ("Add feature" not "Added feature")
- Keep the first line under 72 characters
- Reference issues when applicable (`Fix #123`)

### Code Style

- Follow standard Go conventions (`gofmt`, `goimports`)
- All exported types, functions, and methods must have doc comments
- Use the existing patterns in the codebase as reference:
  - Builder pattern for request params (see `types/order.go`)
  - Functional options for client configuration (see `option.go`)
  - Enum types as `type X string` with constants (see `enums/`)

### Testing

- Write tests for all new code
- Aim for high coverage on all branches and edge cases
- Unit tests go in `*_test.go` files next to the code they test
- Integration tests use the `//go:build integration` build tag
- Run integration tests: `go test -tags=integration ./...`

```bash
# Run unit tests only
make test

# Run with verbose output
go test -v ./...

# Run a specific package
go test -v ./enums/...

# Run integration tests (requires THALEX_PEM_PATH and THALEX_KEY_ID)
go test -tags=integration ./...
```

### Pull Request Process

1. Fork the repository
2. Create your branch from `master`
3. Make your changes
4. Ensure all checks pass: `make check`
5. Push to your fork and open a Pull Request against `master`
6. Fill in the PR template with a description of your changes

### What We Look For

- Code correctness and safety
- Test coverage for new code
- Consistent style with the existing codebase
- Clear documentation for public APIs
- No unnecessary dependencies

## Project Structure

```
go-thalex/
  thalex.go, auth.go, option.go, errors.go  — Root package
  enums/          — Enum types (one file per type)
  types/          — Request/response types (grouped by domain)
  rest/           — REST API client
  ws/             — WebSocket JSON-RPC client
  internal/       — Internal packages (transport, jsonrpc)
  examples/       — Runnable example programs
  docs/           — Documentation
```

## Reporting Issues

- Use [GitHub Issues](https://github.com/amiwrpremium/go-thalex/issues)
- Search existing issues before creating a new one
- Use the provided issue templates
- Include Go version, OS, and SDK version in bug reports

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
