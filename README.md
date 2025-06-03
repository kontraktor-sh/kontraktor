# kontraktor-sh
Kontraktor.sh

## Build and Install from Source

### Prerequisites
- Go 1.22 or newer installed ([Download Go](https://go.dev/dl/))

### Clone the Repository
```sh
git clone https://github.com/<youruser>/kontraktor.git
cd kontraktor
```

### Build the CLI
```sh
cd kontraktor
go build -o kontraktor ./cmd/ktr
```

### Install the CLI
Move the binary to a directory in your `$PATH` (e.g., `/usr/local/bin`):

```sh
sudo mv kontraktor /usr/local/bin/
```

Or, for user-local install (make sure `~/go/bin` is in your `$PATH`):

```sh
mv kontraktor ~/go/bin/
```

### Verify Installation
```sh
kontraktor --help
```
