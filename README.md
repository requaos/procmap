# procmap

A terminal-based process visualization tool written in Go. Displays real-time process information using an interactive bubble chart visualization.

## Features

- Real-time process monitoring with auto-refresh (every 2 seconds)
- **Bubble chart visualization** - bubble size proportional to selected metric
- Multiple metric modes: CPU%, Memory%, or Thread count
- Configurable process count (adjust how many bubbles to show)
- Color-coded bubbles based on metric intensity
- Circle packing algorithm for optimal space usage
- Clean terminal UI with keyboard controls

## Installation

```bash
go mod download
go build -o procmap
```

## Usage

Run the process monitor:

```bash
./procmap
```

### Keyboard Controls

- `c` - Switch to CPU% mode (bubble size = CPU usage)
- `m` - Switch to Memory% mode (bubble size = memory usage)
- `t` - Switch to Threads mode (bubble size = thread count)
- `+` or `=` - Increase number of bubbles shown (by 10)
- `-` or `_` - Decrease number of bubbles shown (by 10)
- `q` or `Ctrl+C` - Quit

## Architecture

The project is organized into three main packages:

- **main**: Entry point that initializes the TUI
- **proc**: Process information gathering using gopsutil
- **ui**: Bubble Tea TUI implementation with Lip Gloss styling

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform process/system info

## Development

Run with auto-reload during development:

```bash
go run main.go
```

Run tests:

```bash
go test ./...
```

## Platform Support

Currently optimized for Linux systems. Should work on macOS and Windows but may require additional testing.
