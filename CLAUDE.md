# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Download dependencies
go mod download

# Build the binary
go build -o procmap

# Run the application
./procmap

# Run during development
go run main.go

# Run tests
go test ./...

# Run tests for specific package
go test ./proc
go test ./ui

# Update dependencies
go mod tidy
```

## Code Architecture

This is a terminal UI application for visualizing process information in real-time.

### Package Structure

**main** (`main.go`):
- Entry point that initializes the Bubble Tea program
- Runs the TUI in alternate screen mode
- Handles top-level error handling

**proc** (`proc/process.go`):
- Abstracts process information gathering via gopsutil
- `ProcessInfo` struct contains normalized process data (PID, name, CPU%, memory%, status, etc.)
- `GetProcesses()` returns all accessible processes with error handling for permission-denied cases

**ui** (`ui/model.go` and `ui/bubble.go`):
- Implements the Bubble Tea Model-View-Update (MVU) pattern
- `model` struct holds application state (processes list, dimensions, errors, sortMode, maxBubbles)
- Auto-refreshes process list every 2 seconds via tick commands
- Renders bubble chart visualization where bubble size represents the selected metric
- Keyboard controls: c/m/t to switch metrics, +/- to adjust count, q to quit

### Key Design Patterns

**MVU (Model-View-Update)**:
The UI follows Bubble Tea's MVU architecture:
- `Init()` returns initial commands to execute
- `Update(msg)` handles messages (keyboard, ticks, process updates) and returns new model + commands
- `View()` renders current state as a string

**Message Types**:
- `tickMsg` - Periodic refresh trigger (every 2s)
- `processesMsg` - New process data loaded
- `errMsg` - Error occurred during process fetching
- `tea.KeyMsg` - Keyboard input
- `tea.WindowSizeMsg` - Terminal size changed

**Process Sorting**:
Processes are sorted based on `sortMode` (cpu/memory/threads) in the `sortProcesses()` method. Sorting happens when new process data arrives and when the user switches modes.

**Bubble Visualization** (`ui/bubble.go`):
- `packBubbles()` - Implements circle packing algorithm to position bubbles on screen
  - Takes top N processes (based on maxBubbles setting)
  - Calculates bubble sizes proportional to metric value
  - Positions largest bubble at center, packs remaining in expanding spiral pattern
  - Uses collision detection to prevent bubble overlap
- `renderBubbles()` - Creates 2D character canvas and draws all bubbles with colors
- `drawBubble()` - Draws individual bubble using Unicode circles (○) with process name + metric value
- Color coding: Red (top 25%), orange (25-50%), yellow (50-75%), green (bottom 25%)

## Key Dependencies

- **Bubble Tea**: TUI framework providing the MVU pattern and terminal handling
- **Lip Gloss**: Terminal styling (colors, bold, backgrounds) - see styles in `ui/model.go`
- **gopsutil/v3**: Cross-platform process information library

## Keyboard Controls

- `c` - Switch to CPU% mode (bubble size = CPU usage)
- `m` - Switch to Memory% mode (bubble size = memory usage)
- `t` - Switch to Threads mode (bubble size = thread count)
- `+` or `=` - Increase bubble count by 10
- `-` or `_` - Decrease bubble count by 10
- `q` or `Ctrl+C` - Quit

## Extension Points

To add new features:

- **New process fields**: Add to `ProcessInfo` struct in `proc/process.go` and populate in `getProcessInfo()`. Update bubble label in `drawBubble()`
- **New metrics for bubbles**: Add case to `getValue()` in `ui/bubble.go` and add keybinding in `Update()` in `ui/model.go`
- **Filtering**: Add filter state to model struct and filter processes before calling `packBubbles()`
- **Alternative layouts**: Modify `packBubbles()` algorithm (currently uses spiral packing, could add grid, random, force-directed)
- **Interactive selection**: Add cursor state to highlight specific bubble, show detailed info on selection
