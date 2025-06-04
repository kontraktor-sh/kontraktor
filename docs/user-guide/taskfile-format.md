# Taskfile Format

The Kontraktor taskfile is written in YAML format and uses the `.ktr.yml` extension. This document describes the structure and available options.

## Basic Structure

```yaml
version: "0.3"  # Required: Version of the taskfile format

environment:     # Optional: Global environment variables
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
      - command or task reference
```

## Version

The `version` field is required and specifies the version of the taskfile format. Currently, only version "0.3" is supported.

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
      - echo "${GLOBAL_VAR}"  # Uses global variable
      - echo "${TASK_VAR}"   # Uses task variable
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

Commands can be:

1. Shell commands:
   ```yaml
   cmds:
     - echo "Hello, World!"
     - ls -la
   ```

2. Task references:
   ```yaml
   cmds:
     - task: setup
   ```

3. Uses references (with arguments):
   ```yaml
   cmds:
     - uses: docker-build
       args:
         image: my-image
         tag: latest
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

## Task Dependencies

Tasks can depend on other tasks through the `task` command:

```yaml
tasks:
  setup:
    cmds:
      - echo "Setting up..."

  build:
    cmds:
      - task: setup
      - echo "Building..."
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
      - echo "Hello, ${name}!"
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

## Examples

### Complete Example

```yaml
version: "0.3"

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
      - mkdir -p ${BUILD_DIR}
      - echo "Setup complete"

  build:
    desc: Build the project
    args:
      - name: version
        type: string
        default: 1.0.0
    cmds:
      - task: setup
      - echo "Building ${PROJECT_NAME} version ${version}"
      - touch ${BUILD_DIR}/output.txt

  deploy:
    desc: Deploy to environment
    args:
      - name: env
        type: string
        default: dev
    environment:
      DEPLOY_ENV: ${env}
    cmds:
      - task: build
      - echo "Deploying to ${DEPLOY_ENV}"
      - echo "Using API key: ${API_KEY}"
``` 