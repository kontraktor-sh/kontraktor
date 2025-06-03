# Quick Start Guide

This guide will help you get started with Kontraktor quickly.

## Your First Taskfile

Create a new file named `taskfile.ktr.yml` in your project directory:

```yaml
version: "0.3"

environment:
  GREETING: "Hello"

tasks:
  hello:
    desc: A simple hello world task
    args:
      - name: name
        type: string
        default: World
    cmds:
      - echo "${GREETING}, ${name}!"
```

## Running Tasks

To run the `hello` task with default arguments:

```bash
kontraktor run hello
```

To provide a custom argument:

```bash
kontraktor run hello name=John
```

## Task Dependencies

You can create tasks that depend on other tasks:

```yaml
version: "0.3"

tasks:
  setup:
    desc: Setup the environment
    cmds:
      - echo "Setting up environment..."
      - mkdir -p build

  build:
    desc: Build the project
    cmds:
      - task: setup
      - echo "Building project..."
      - touch build/output.txt

  test:
    desc: Run tests
    cmds:
      - task: build
      - echo "Running tests..."
```

## Using Environment Variables

You can define environment variables at different levels:

```yaml
version: "0.3"

environment:
  GLOBAL_VAR: global-value

tasks:
  env-test:
    desc: Test environment variables
    environment:
      TASK_VAR: task-value
    cmds:
      - echo "Global: ${GLOBAL_VAR}"
      - echo "Task: ${TASK_VAR}"
```

## Secret Management

To use Azure Key Vault secrets:

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
  secret-test:
    desc: Test secret access
    cmds:
      - echo "Using API key: ${API_KEY}"
      - echo "Using DB password: ${DB_PASSWORD}"
```

## Task Arguments

Define and validate task arguments:

```yaml
version: "0.3"

tasks:
  deploy:
    desc: Deploy to environment
    args:
      - name: env
        type: string
        default: dev
      - name: version
        type: string
    cmds:
      - echo "Deploying version ${version} to ${env}"
```

## Next Steps

- Learn more about the [Taskfile Format](user-guide/taskfile-format.md)
- Explore [Secret Management](user-guide/secret-management.md)
- Read about [Task Dependencies](advanced/task-dependencies.md)
- Check out [Azure Key Vault Integration](advanced/azure-keyvault.md) 