# Welcome to Kontraktor

Kontraktor is a powerful builder-helper that unifies environment configuration, secret aggregation, and developer task automation. Inspired by [taskfile.dev](https://taskfile.dev) but extended with a centralized configuration library and remote-vault integration.

## Key Features

### Central Configuration Library
- Version-controlled, queryable configuration across many projects and git repositories
- Stored in DuckDB (non-secret only)
- Easy to maintain and share configurations

### Secret Aggregation Layer
- Read-only connectors to various secret management systems:
  - Azure Key Vault
  - HashiCorp Vault
  - AWS Secrets Manager
- No secrets persisted in Kontraktor
- Secure secret management

### Task Runner
- Local execution of Bash (macOS/Linux) routines
- Expressive YAML syntax (`taskfile.ktr.yml`)
- Task dependencies and imports
- Environment variable management
- Argument passing and validation

## Quick Links

- [Installation Guide](getting-started/installation.md)
- [Quick Start Guide](getting-started/quickstart.md)
- [Taskfile Format](user-guide/taskfile-format.md)
- [Secret Management](user-guide/secret-management.md)

## Example Taskfile

```yaml
version: "0.3"

environment:
  GLOBAL_VAR: global-value

tasks:
  hello:
    desc: A simple hello world task
    args:
      - name: name
        type: string
        default: World
    cmds:
      - echo "Hello, ${name}!"
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](contributing/development.md) for more information.

## License

Kontraktor is open source software licensed under the [GNU 3.0](https://github.com/kontraktor-sh/kontraktor/blob/main/LICENSE). 