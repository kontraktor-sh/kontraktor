# Installation Guide

This guide will help you install Kontraktor on your system.

## Prerequisites

- Go 1.22 or newer installed ([Download Go](https://go.dev/dl/))
- Git (for cloning the repository)

## Installation Methods

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/kontraktor-sh/kontraktor.git
   cd kontraktor
   ```

2. Build the CLI:
   ```bash
   go build -o kontraktor ./cmd/kontraktor
   ```

3. Install the binary:
   - System-wide installation (requires sudo):
     ```bash
     sudo mv kontraktor /usr/local/bin/
     ```
   - User-local installation (make sure `~/go/bin` is in your `$PATH`):
     ```bash
     mv kontraktor ~/go/bin/
     ```

### Using Go Install

You can also install Kontraktor directly using Go:

```bash
go install github.com/kontraktor-sh/kontraktor/cmd/kontraktor@latest
```

## Verify Installation

After installation, verify that Kontraktor is properly installed:

```bash
kontraktor --help
```

You should see the help message with available commands and options.

## Configuration

Kontraktor uses a configuration file named `taskfile.ktr.yml` in your project directory. See the [Taskfile Format](user-guide/taskfile-format.md) guide for details on how to configure your tasks.

## Next Steps

- Read the [Quick Start Guide](quickstart.md) to learn how to use Kontraktor
- Check out the [Taskfile Format](user-guide/taskfile-format.md) documentation
- Learn about [Secret Management](user-guide/secret-management.md)

## Troubleshooting

### Common Issues

1. **Command not found**
   - Ensure the installation directory is in your `$PATH`
   - Try running `which kontraktor` to verify the installation location

2. **Permission denied**
   - Check file permissions: `chmod +x /path/to/kontraktor`
   - Ensure you have the necessary permissions for the installation directory

3. **Go version issues**
   - Verify your Go version: `go version`
   - Update Go if necessary

If you encounter any other issues, please [open an issue](https://github.com/kontraktor-sh/kontraktor/issues) on GitHub. 