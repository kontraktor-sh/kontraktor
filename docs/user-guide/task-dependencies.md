# Task Dependencies

Kontraktor allows you to create complex workflows by defining dependencies between tasks. This document explains how to use task dependencies effectively.

## Basic Task References

Tasks can reference other tasks using the `task` command type:

```yaml
version: "0.3"

tasks:
  setup:
    desc: Setup the environment
    cmds:
      - type: bash
        content:
          command: echo "Setting up environment..."
      - type: bash
        content:
          command: mkdir -p build

  build:
    desc: Build the project
    cmds:
      - type: task
        content:
          name: setup
      - type: bash
        content:
          command: echo "Building project..."
      - type: bash
        content:
          command: touch build/output.txt

  test:
    desc: Run tests
    cmds:
      - type: task
        content:
          name: build
      - type: bash
        content:
          command: echo "Running tests..."
```

## Passing Arguments to Dependent Tasks

You can pass arguments to dependent tasks:

```yaml
version: "0.3"

tasks:
  build:
    desc: Build with version
    args:
      - name: version
        type: string
        default: 1.0.0
    cmds:
      - type: bash
        content:
          command: echo "Building version ${version}"

  deploy:
    desc: Deploy the build
    cmds:
      - type: task
        content:
          name: build
          args:
            version: 2.0.0
      - type: bash
        content:
          command: echo "Deploying..."
```

## Environment Variable Inheritance

Dependent tasks inherit environment variables from their parent tasks:

```yaml
version: "0.3"

environment:
  GLOBAL_VAR: global

tasks:
  setup:
    desc: Setup with environment
    environment:
      SETUP_VAR: setup-value
    cmds:
      - type: bash
        content:
          command: echo "Setup: ${SETUP_VAR}"

  build:
    desc: Build with inherited environment
    environment:
      BUILD_VAR: build-value
    cmds:
      - type: task
        content:
          name: setup
      - type: bash
        content:
          command: echo "Build: ${BUILD_VAR}, ${SETUP_VAR}, ${GLOBAL_VAR}"
```

## Best Practices

1. **Task Organization**
   - Group related tasks together
   - Use consistent naming conventions
   - Consider using task prefixes for organization

2. **Error Handling**
   - Handle errors in dependent tasks
   - Use appropriate exit codes
   - Consider using `--force` for specific cases

3. **Task Dependencies**
   - Keep dependency chains short and clear
   - Avoid circular dependencies
   - Document task dependencies in descriptions

## Examples

### Development Workflow

```yaml
version: "0.3"

tasks:
  clean:
    desc: Clean build artifacts
    cmds:
      - type: bash
        content:
          command: rm -rf build/
      - type: bash
        content:
          command: rm -rf dist/

  install-deps:
    desc: Install dependencies
    cmds:
      - type: bash
        content:
          command: npm install

  lint:
    desc: Run linter
    cmds:
      - type: task
        content:
          name: install-deps
      - type: bash
        content:
          command: npm run lint

  test:
    desc: Run tests
    cmds:
      - type: task
        content:
          name: install-deps
      - type: bash
        content:
          command: npm test

  build:
    desc: Build the application
    cmds:
      - type: task
        content:
          name: install-deps
      - type: task
        content:
          name: lint
      - type: task
        content:
          name: test
      - type: bash
        content:
          command: npm run build

  deploy:
    desc: Deploy to production
    cmds:
      - type: task
        content:
          name: build
      - type: bash
        content:
          command: npm run deploy
```

### CI/CD Pipeline

```yaml
version: "0.3"

tasks:
  validate:
    desc: Validate code
    cmds:
      - type: bash
        content:
          command: echo "Validating code..."

  security-scan:
    desc: Run security scan
    cmds:
      - type: bash
        content:
          command: echo "Running security scan..."

  build:
    desc: Build artifacts
    cmds:
      - type: task
        content:
          name: validate
      - type: task
        content:
          name: security-scan
      - type: bash
        content:
          command: echo "Building artifacts..."

  test:
    desc: Run tests
    cmds:
      - type: task
        content:
          name: build
      - type: bash
        content:
          command: echo "Running tests..."

  package:
    desc: Package artifacts
    cmds:
      - type: task
        content:
          name: test
      - type: bash
        content:
          command: echo "Packaging artifacts..."

  deploy-staging:
    desc: Deploy to staging
    cmds:
      - type: task
        content:
          name: package
      - type: bash
        content:
          command: echo "Deploying to staging..."

  deploy-prod:
    desc: Deploy to production
    cmds:
      - type: task
        content:
          name: deploy-staging
      - type: bash
        content:
          command: echo "Deploying to production..."
```

### Multi-Environment Deployment

```yaml
version: "0.3"

tasks:
  build:
    desc: Build the application
    args:
      - name: env
        type: string
        default: dev
    cmds:
      - type: bash
        content:
          command: echo "Building for ${env}"

  deploy:
    desc: Deploy to environment
    args:
      - name: env
        type: string
        default: dev
    cmds:
      - type: task
        content:
          name: build
          args:
            env: ${env}
      - type: bash
        content:
          command: echo "Deploying to ${env}"
```

## Troubleshooting

### Common Issues

1. **Circular Dependencies**
   - Error: "Circular dependency detected"
   - Solution: Review and restructure task dependencies

2. **Missing Dependencies**
   - Error: "Task not found"
   - Solution: Check task names and ensure they exist

3. **Execution Order**
   - Issue: Tasks not executing in expected order
   - Solution: Review dependency chain and task definitions

### Debugging

Enable debug logging to see dependency resolution:

```bash
KONTRAKTOR_DEBUG=1 kontraktor run task-name
``` 