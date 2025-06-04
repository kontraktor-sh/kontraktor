# Task Dependencies

This guide explains how to work with task dependencies in Kontraktor.

## Basic Dependencies

Tasks can depend on other tasks using the `deps` field:

```yaml
version: "0.3"

tasks:
  build:
    desc: Build the application
    cmds:
      - echo "Building..."

  test:
    desc: Run tests
    deps: [build]
    cmds:
      - echo "Testing..."

  deploy:
    desc: Deploy the application
    deps: [test]
    cmds:
      - echo "Deploying..."
```

When you run `kontraktor run deploy`, it will execute tasks in this order:
1. `build`
2. `test`
3. `deploy`

## Multiple Dependencies

A task can depend on multiple other tasks:

```yaml
tasks:
  prepare:
    desc: Prepare environment
    cmds:
      - echo "Preparing..."

  build:
    desc: Build the application
    cmds:
      - echo "Building..."

  deploy:
    desc: Deploy the application
    deps: [prepare, build]
    cmds:
      - echo "Deploying..."
```

## Conditional Dependencies

You can make dependencies conditional using environment variables:

```yaml
tasks:
  build:
    desc: Build the application
    cmds:
      - echo "Building..."

  test:
    desc: Run tests
    deps: [build]
    cmds:
      - echo "Testing..."

  deploy:
    desc: Deploy the application
    deps: [test]
    cmds:
      - echo "Deploying..."

  deploy-skip-tests:
    desc: Deploy without running tests
    deps: [build]
    cmds:
      - echo "Deploying without tests..."
```

## Dependency Cycles

Kontraktor detects and prevents dependency cycles. For example, this would cause an error:

```yaml
tasks:
  task1:
    deps: [task2]
    cmds:
      - echo "Task 1"

  task2:
    deps: [task1]
    cmds:
      - echo "Task 2"
```

## Best Practices

1. **Keep Dependencies Simple**
   - Avoid deep dependency chains
   - Use clear, meaningful task names
   - Document complex dependency relationships

2. **Task Organization**
   - Group related tasks together
   - Use consistent naming conventions
   - Consider using task prefixes for organization

3. **Error Handling**
   - Handle errors in dependent tasks
   - Use appropriate exit codes
   - Consider using `--force` for specific cases

## Examples

### Development Workflow

```yaml
version: "0.3"

tasks:
  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf build/
      - rm -rf dist/

  install-deps:
    desc: Install dependencies
    cmds:
      - npm install

  lint:
    desc: Run linter
    deps: [install-deps]
    cmds:
      - npm run lint

  test:
    desc: Run tests
    deps: [install-deps]
    cmds:
      - npm test

  build:
    desc: Build the application
    deps: [install-deps, lint, test]
    cmds:
      - npm run build

  deploy:
    desc: Deploy to production
    deps: [build]
    cmds:
      - npm run deploy
```

### CI/CD Pipeline

```yaml
version: "0.3"

tasks:
  validate:
    desc: Validate code
    cmds:
      - echo "Validating code..."

  security-scan:
    desc: Run security scan
    cmds:
      - echo "Running security scan..."

  build:
    desc: Build artifacts
    deps: [validate, security-scan]
    cmds:
      - echo "Building artifacts..."

  test:
    desc: Run tests
    deps: [build]
    cmds:
      - echo "Running tests..."

  package:
    desc: Package artifacts
    deps: [test]
    cmds:
      - echo "Packaging artifacts..."

  deploy-staging:
    desc: Deploy to staging
    deps: [package]
    cmds:
      - echo "Deploying to staging..."

  deploy-prod:
    desc: Deploy to production
    deps: [deploy-staging]
    cmds:
      - echo "Deploying to production..."
```

### Multi-Environment Deployment

```yaml
version: "0.3"

tasks:
  build:
    desc: Build the application
    cmds:
      - echo "Building..."

  test:
    desc: Run tests
    deps: [build]
    cmds:
      - echo "Testing..."

  deploy-dev:
    desc: Deploy to development
    deps: [test]
    environment:
      ENV: development
    cmds:
      - echo "Deploying to ${ENV}..."

  deploy-staging:
    desc: Deploy to staging
    deps: [deploy-dev]
    environment:
      ENV: staging
    cmds:
      - echo "Deploying to ${ENV}..."

  deploy-prod:
    desc: Deploy to production
    deps: [deploy-staging]
    environment:
      ENV: production
    cmds:
      - echo "Deploying to ${ENV}..."
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