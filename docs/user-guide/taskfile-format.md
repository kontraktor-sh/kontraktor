# Taskfile Format

The Kontraktor taskfile is written in YAML format and uses the `.ktr.yml` extension. This document describes the structure and available options.

## Basic Structure

```yaml
version: "0.3"  # Required: Version of the taskfile format

imports:        # Optional: Import tasks from other taskfiles
  - path/to/taskfile.ktr.yml
  - https://github.com/user/repo.git//path/to/taskfile.ktr.yml

environment:    # Optional: Global environment variables
  KEY: value

vaults:         # Optional: Secret vault configurations
  azure_keyvault:
    vault-name:
      keyvault_name: name
      secrets:
        ENV_VAR: secret-name

tasks:          # Required: Task definitions
  task-name:
    desc: "Task description"
    args:
      - name: arg-name
        type: string
        default: value
    environment:
      KEY: value
    cmds:
      - type: bash
        content:
          command: command-string
      - type: task
        content:
          name: task-name
          args:
            key: value
```

## Version

The `version` field is required and specifies the version of the taskfile format. Currently, only version "0.3" is supported.

## Imports

The `imports` field allows you to import tasks from other taskfiles. You can import from:
- Local files
- HTTP(S) URLs
- Git repositories (using the format `https://github.com/user/repo.git//path/to/file`)

Imported tasks are merged with the main taskfile, with tasks in the main file taking precedence.

## Environment Variables

Environment variables can be defined at three levels:

1. Global (taskfile level)
2. Task level
3. Command level (through variable substitution)

```yaml
version: "0.3"

environment:
  GLOBAL_VAR: global-value

tasks:
  env-test:
    environment:
      TASK_VAR: task-value
    cmds:
      - type: bash
        content:
          command: echo "${GLOBAL_VAR}"  # Uses global variable
      - type: bash
        content:
          command: echo "${TASK_VAR}"   # Uses task variable
```

## Tasks

Tasks are the main building blocks of a taskfile. Each task can have:

- Description
- Arguments
- Environment variables
- Commands

### Task Arguments

Arguments can be defined with name, type, and optional default value:

```yaml
tasks:
  deploy:
    args:
      - name: environment
        type: string
        default: dev
      - name: version
        type: string
```

### Task Commands

Commands can be of different types:

1. Bash commands:
   ```yaml
   cmds:
     - type: bash
       content:
         command: echo "Hello, World!"
     - type: bash
       content:
         command: ls -la
   ```

2. Task references:
   ```yaml
   cmds:
     - type: task
       content:
         name: setup
         args:
           key: value
   ```

3. Docker commands:
   ```yaml
   cmds:
     - type: docker
       content:
         image: node:latest
         command: ["npm", "start"]
         environment:
           NODE_ENV: production
         volumes:
           ./src:/app/src
   ```

## Secret Management

### Azure Key Vault

```yaml
vaults:
  azure_keyvault:
    my-vault:
      keyvault_name: my-keyvault
      secrets:
        API_KEY: api-secret
        DB_PASSWORD: db-secret
```

## Variable Substitution

Variables can be substituted in commands using `${VAR_NAME}` syntax:

```yaml
tasks:
  greet:
    args:
      - name: name
        default: World
    cmds:
      - type: bash
        content:
          command: echo "Hello, ${name}!"
```

## Best Practices

1. **Task Organization**
   - Group related tasks together
   - Use descriptive task names
   - Include clear descriptions

2. **Environment Variables**
   - Use global variables for project-wide settings
   - Use task variables for task-specific settings
   - Keep sensitive data in vaults

3. **Arguments**
   - Provide default values when possible
   - Use descriptive argument names
   - Document required arguments

4. **Secret Management**
   - Never store secrets in the taskfile
   - Use vaults for all sensitive data
   - Use descriptive secret names

5. **Command Structure**
   - Always use the explicit command type and content structure
   - Group related commands together
   - Use meaningful command descriptions

## Examples

### Complete Example

```yaml
version: "0.3"

imports:
  - https://github.com/kontraktor-sh/kontraktor.git//templates/docker.ktr.yml

environment:
  PROJECT_NAME: my-project
  BUILD_DIR: build

vaults:
  azure_keyvault:
    prod-vault:
      keyvault_name: prod-keyvault
      secrets:
        API_KEY: api-secret
        DB_PASSWORD: db-secret

tasks:
  setup:
    desc: Setup the build environment
    cmds:
      - type: bash
        content:
          command: mkdir -p ${BUILD_DIR}
      - type: bash
        content:
          command: echo "Setup complete"

  build:
    desc: Build the project
    args:
      - name: version
        type: string
        default: 1.0.0
    cmds:
      - type: task
        content:
          name: setup
      - type: bash
        content:
          command: echo "Building ${PROJECT_NAME} version ${version}"
      - type: bash
        content:
          command: touch ${BUILD_DIR}/output.txt

  deploy:
    desc: Deploy to environment
    args:
      - name: env
        type: string
        default: dev
    environment:
      DEPLOY_ENV: ${env}
    cmds:
      - type: task
        content:
          name: build
      - type: bash
        content:
          command: echo "Deploying to ${DEPLOY_ENV}"
      - type: bash
        content:
          command: echo "Using API key: ${API_KEY}"
``` 