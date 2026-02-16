# Notify - Windows CLI Notification Utility

A lightweight CLI utility for sending beautiful Windows toast notifications directly from the command line.

## Features

- **Native Windows Toast Notifications**: Real desktop notifications, not console output
- **Multiple Types**: Success (green), Error (red), Info (blue), Warning (yellow)
- **Configurable**: Set timeout and auto-close behavior
- **Sound Support**: Different sounds for different notification types
- **Optimized**: Small binary with fast startup

## Installation

1. Download `notify.exe`
2. Place it in your PATH (e.g., `C:\Windows\` or add folder to PATH)

## Building from Source

```bash
# Install dependencies
go mod tidy

# Build optimized executable
go build -ldflags="-s -w" -trimpath -o notify.exe main.go
```

Or use the provided build script:

```bash
build.bat
```

## Usage

```bash
notify MESSAGE [OPTIONS]
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `--type` | Type: success, error, info, warning | info |
| `--timeout` | Timeout in seconds | 5 |
| `--autoclose` | Auto close after timeout (true/false) | true |
| `--help` | Show help message | - |

### Examples

```bash
# Success notification
notify "Build completed successfully!" --type success

# Error notification
notify "Compilation failed!" --type error

# Warning with custom timeout
notify "Low disk space" --type warning --timeout 10

# Info without auto-close
notify "Download started" --type info --autoclose false

# Quick notification (uses defaults)
notify "Task done"
```

## Notification Types

| Type | Title | Use Case |
|------|-------|----------|
| `success` | Success | Operations completed successfully |
| `error` | Error | Errors or failures |
| `info` | Info | Informational messages |
| `warning` | Warning | Warnings or cautionary messages |

## Requirements

- Windows 10/11
- Go 1.16+ (for building from source)

## License

MIT
