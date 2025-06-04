# Contributing to Kontraktor

Thank you for your interest in contributing to Kontraktor! This guide will help you get started with the development process.

## Development Setup

### Prerequisites

1. **Go 1.22 or newer**
   - Install from [golang.org](https://golang.org/dl/)
   - Verify installation:
     ```bash
     go version
     ```

2. **Git**
   - Install from [git-scm.com](https://git-scm.com/downloads)
   - Verify installation:
     ```bash
     git --version
     ```

3. **Development Tools**
   - [golangci-lint](https://golangci-lint.run/usage/install/) for linting
   - [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) for import formatting
   - [gofmt](https://golang.org/cmd/gofmt/) for code formatting

### Getting Started

1. **Fork the Repository**
   - Visit [github.com/kontraktor-sh/kontraktor](https://github.com/kontraktor-sh/kontraktor)
   - Click "Fork" to create your copy

2. **Clone Your Fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/kontraktor.git
   cd kontraktor
   ```

3. **Add Upstream Remote**
   ```bash
   git remote add upstream https://github.com/kontraktor-sh/kontraktor.git
   ```

4. **Install Dependencies**
   ```bash
   go mod download
   ```

## Development Workflow

### Branching Strategy

1. **Main Branch**
   - `main`: Production-ready code
   - `develop`: Development branch

2. **Feature Branches**
   - Format: `feature/description`
   - Example: `feature/azure-keyvault`

3. **Bug Fix Branches**
   - Format: `fix/description`
   - Example: `fix/secret-rotation`

### Making Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature
   ```

2. **Make Changes**
   - Write code
   - Add tests
   - Update documentation

3. **Run Tests**
   ```bash
   go test ./...
   ```

4. **Run Linter**
   ```bash
   golangci-lint run
   ```

5. **Format Code**
   ```bash
   go fmt ./...
   goimports -w .
   ```

6. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

7. **Push Changes**
   ```bash
   git push origin feature/your-feature
   ```

### Pull Request Process

1. **Create Pull Request**
   - Visit your fork on GitHub
   - Click "New Pull Request"
   - Select `develop` as the base branch

2. **Pull Request Template**
   - Fill out the PR template
   - Describe changes
   - Link related issues
   - Add screenshots if applicable

3. **Code Review**
   - Address review comments
   - Update PR as needed
   - Ensure CI passes

4. **Merge**
   - Squash and merge
   - Delete feature branch

## Code Style

### Go Code

1. **Formatting**
   - Use `gofmt`
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Use `goimports` for import organization

2. **Naming**
   - Use camelCase for variables
   - Use PascalCase for exported names
   - Use short names in small scopes

3. **Comments**
   - Document exported functions
   - Use complete sentences
   - Follow [godoc](https://blog.golang.org/godoc) style

### YAML Files

1. **Formatting**
   - Use 2 spaces for indentation
   - Use consistent quotes
   - Sort keys alphabetically

2. **Comments**
   - Use `#` for comments
   - Add section headers
   - Document complex configurations

## Testing

### Unit Tests

1. **Test Files**
   - Name: `*_test.go`
   - Location: Same directory as source
   - Package: Same as source

2. **Test Functions**
   ```go
   func TestFunction(t *testing.T) {
       // Arrange
       // Act
       // Assert
   }
   ```

3. **Table-Driven Tests**
   ```go
   func TestFunction(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected string
       }{
           // Test cases
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test logic
           })
       }
   }
   ```

### Integration Tests

1. **Test Files**
   - Name: `*_integration_test.go`
   - Location: `tests/` directory
   - Package: `tests`

2. **Test Setup**
   ```go
   func TestMain(m *testing.M) {
       // Setup
       code := m.Run()
       // Teardown
       os.Exit(code)
   }
   ```

## Documentation

### Code Documentation

1. **Package Documentation**
   ```go
   // Package name provides functionality for...
   package name
   ```

2. **Function Documentation**
   ```go
   // FunctionName does something.
   // It takes parameters and returns results.
   func FunctionName() error {
       // Implementation
   }
   ```

### User Documentation

1. **Markdown Files**
   - Location: `docs/` directory
   - Format: Markdown
   - Include code examples

2. **API Documentation**
   - Location: `docs/api/` directory
   - Format: OpenAPI/Swagger
   - Include request/response examples

## Release Process

1. **Version Bumping**
   - Update version in `VERSION` file
   - Update changelog
   - Tag release

2. **Release Steps**
   ```bash
   # Update version
   echo "1.0.0" > VERSION
   
   # Update changelog
   # Commit changes
   git commit -am "chore: release v1.0.0"
   
   # Tag release
   git tag -a v1.0.0 -m "Release v1.0.0"
   
   # Push changes
   git push origin main --tags
   ```

## Community Guidelines

### Code of Conduct

1. **Be Respectful**
   - Use welcoming language
   - Be patient with others
   - Accept constructive criticism

2. **Be Professional**
   - Focus on the code
   - Avoid personal attacks
   - Keep discussions constructive

3. **Be Helpful**
   - Answer questions
   - Share knowledge
   - Mentor others

### Communication

#### Issues
   - Use templates
   - Provide context
   - Be specific

### Discussions
   - Use appropriate channels
   - Stay on topic
   - Be constructive

### Pull Requests
   - Follow guidelines
   - Be responsive
   - Keep changes focused

## Getting Help

### Documentation
   - Read the docs
   - Check examples
   - Search issues

### Community
   - Join discussions
   - Ask questions
   - Share ideas

### Support
   - Open issues
   - Contact maintainers
   - Use chat channels 