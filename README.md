# Kontraktor

[![Build and Test](https://github.com/kontraktor-sh/kontraktor/actions/workflows/build.yml/badge.svg)](https://github.com/kontraktor-sh/kontraktor/actions/workflows/build.yml)
[![Release](https://github.com/kontraktor-sh/kontraktor/actions/workflows/release.yml/badge.svg)](https://github.com/kontraktor-sh/kontraktor/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kontraktor-sh/kontraktor)](https://goreportcard.com/report/github.com/kontraktor-sh/kontraktor)
[![GoDoc](https://godoc.org/github.com/kontraktor-sh/kontraktor?status.svg)](https://godoc.org/github.com/kontraktor-sh/kontraktor)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Kontraktor is a powerful task runner and automation tool that helps you manage and execute tasks with ease. It provides a simple YAML-based configuration format and supports various features like task dependencies, environment variables, and secret management.

## Features

- **Simple YAML Configuration**: Define tasks in a clear, readable format
- **Task Dependencies**: Create complex workflows with task dependencies
- **Environment Variables**: Manage environment variables at global, task, and command levels
- **Secret Management**: Securely manage secrets using Azure Key Vault
- **Shell Command Execution**: Run shell commands with proper environment setup
- **Task Arguments**: Define and validate task arguments
- **Extensible**: Add support for more secret vaults and features

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/kontraktor-sh/kontraktor.git
cd kontraktor

# Build the CLI
go build -o kontraktor cmd/kontraktor/main.go

# Install the binary (optional)
sudo mv kontraktor /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/kontraktor-sh/kontraktor@latest
```

## Quick Start

1. Create a `taskfile.ktr.yml`:

```yaml
version: "0.3"

tasks:
  hello:
    desc: Say hello
    cmds:
      - echo "Hello, World!"
```

2. Run the task:

```bash
kontraktor run hello
```

## Documentation

- [Getting Started](docs/getting-started/installation.md)
- [Taskfile Format](docs/user-guide/taskfile-format.md)
- [Task Dependencies](docs/user-guide/task-dependencies.md)
- [Secret Management](docs/user-guide/secret-management.md)
- [Azure Key Vault Integration](docs/advanced/azure-keyvault.md)
- [Contributing](docs/contributing.md)

## Examples

### Basic Task

```yaml
version: "0.3"

tasks:
  build:
    desc: Build the application
    cmds:
      - go build -o app ./cmd/app
```

### Task with Dependencies

```yaml
version: "0.3"

tasks:
  test:
    desc: Run tests
    cmds:
      - go test ./...

  build:
    desc: Build the application
    deps: [test]
    cmds:
      - go build -o app ./cmd/app
```

### Task with Environment Variables

```yaml
version: "0.3"

environment:
  GO_ENV: development

tasks:
  deploy:
    desc: Deploy the application
    environment:
      GO_ENV: production
    cmds:
      - echo "Deploying to ${GO_ENV}..."
```

### Task with Azure Key Vault Secrets

```yaml
version: "0.3"

vaults:
  azure_keyvault:
    my-vault:
      keyvault_name: my-keyvault
      secrets:
        API_KEY: api-secret
        DB_PASSWORD: db-secret

tasks:
  deploy:
    desc: Deploy with secrets
    cmds:
      - echo "Using API key: ${API_KEY}"
      - echo "Using DB password: ${DB_PASSWORD}"
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](docs/contributing.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
