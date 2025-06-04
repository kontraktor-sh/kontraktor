# Kontraktor

[![Build and Test](https://github.com/kontraktor-sh/kontraktor/actions/workflows/build.yml/badge.svg)](https://github.com/kontraktor-sh/kontraktor/actions/workflows/build.yml)
[![Release](https://github.com/kontraktor-sh/kontraktor/actions/workflows/release.yml/badge.svg)](https://github.com/kontraktor-sh/kontraktor/actions/workflows/release.yml)
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
- **Task Imports**: Import tasks from other taskfiles or Git repositories

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

environment:
  GREETING: "Hello"

tasks:
  hello:
    desc: Say hello
    cmds:
      - type: bash
        content:
          command: echo "${GREETING}, World!"
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
      - type: bash
        content:
          command: go build -o app ./cmd/app
```

### Task with Dependencies

```yaml
version: "0.3"

tasks:
  test:
    desc: Run tests
    cmds:
      - type: bash
        content:
          command: go test ./...

  build:
    desc: Build the application
    cmds:
      - type: task
        content:
          name: test
      - type: bash
        content:
          command: go build -o app ./cmd/app
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
      - type: bash
        content:
          command: echo "Deploying to ${GO_ENV}..."
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
      - type: bash
        content:
          command: echo "Using API key: ${API_KEY}"
      - type: bash
        content:
          command: echo "Using DB password: ${DB_PASSWORD}"
```

### Task with Imports

```yaml
version: "0.3"

imports:
  - https://github.com/kontraktor-sh/kontraktor.git//templates/docker.ktr.yml

tasks:
  build:
    desc: Build using imported tasks
    cmds:
      - type: task
        content:
          name: docker-build
          args:
            image: my-app
            tag: latest
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](docs/contributing.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
