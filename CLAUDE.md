# CLAUDE.md

This file provides guidance to Claude Code when working with this Flappy Bird TUI project.

## Overview

Terminal-based Flappy Bird clone built with Go, Bubble Tea (TUI framework), and Lipgloss (styling). Implements the Elm Architecture (TEA) pattern with Model-Update-View separation.

## Quick Commands

**Binary name: `flappy-bird-tui`** (統一ルール)

```bash
# Build (default target)
make
# or: make build

# Build and run
make build && ./flappy-bird-tui

# Run directly (without building binary)
make run
# or: go run .

# Test
make test
# or: go test ./...

# Lint (requires golangci-lint)
make lint

# Clean build artifacts
make clean

# Multi-platform release build
make release

# Install to GOPATH/bin
make install

# Help
make help

# Version info
./flappy-bird-tui --version

# Release (requires goreleaser)
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

## Architecture

### Elm Architecture Pattern (TEA)

Built on Bubble Tea's TEA implementation with a layered architecture:

1. **Model** (`game/model.go`) - Game state container
2. **Update** (`game/update.go`) - State transitions from messages
3. **View** (`ui/view.go`) - Render ASCII representation

### Package Structure

Follows a layered architecture pattern:

```
flappy-bird-tui/
├── main.go          # Entry point with modelWrapper
├── domain/          # Core domain models (entities)
│   ├── bird.go      # Bird physics and state
│   ├── pipe.go      # Pipe generation and collision
│   ├── difficulty.go # Difficulty level settings
│   └── theme.go     # Color theme definitions
├── game/            # Game logic (use cases)
│   ├── model.go     # Game state and stats
│   └── update.go    # Input/tick handling, state transitions
├── storage/         # Persistence layer
│   ├── highscore.go # JSON save/load to ~/.flappy-bird-tui/
│   └── sound.go     # Terminal beep initialization
└── ui/              # Presentation layer
    └── view.go      # ASCII rendering with theming
```

### Core Components

| Package/File | Purpose | Key Types/Functions |
|-------------|---------|---------------------|
| `main.go` | Entry point, CLI flags | `main()`, `modelWrapper`, version vars |
| `domain/bird.go` | Bird physics | `Bird`, `Update()`, `Jump()` |
| `domain/pipe.go` | Pipe system | `Pipe`, `NewPipe()`, `CheckCollision()` |
| `domain/difficulty.go` | Difficulty settings | `Difficulty`, `GetSettings()` |
| `domain/theme.go` | Color themes | `Theme`, `GetColors()` |
| `game/model.go` | Game state | `Model`, `GameState`, `Stats`, `NewModel()` |
| `game/update.go` | Input/tick handling | `Update()`, `Init()`, `resetGame()` |
| `storage/highscore.go` | Persistence | `HighScore`, `Load/SaveHighScore()`, `Load/SaveRankings()` |
| `storage/sound.go` | Audio feedback | `PlaySound()` (terminal bell) |
| `ui/view.go` | ASCII rendering | `View()`, theme-aware rendering |

### Game States

```go
const (
    StateTitle      // Title screen - select difficulty/theme
    StatePlaying    // Active gameplay
    StateGameOver   // Game over screen with stats
)
```

### Key Mechanics

**Physics:**
- Gravity applies constant downward velocity
- Jump applies fixed upward impulse
- Bird sprite animates (wings up/down)

**Difficulty Scaling:**
- Speed increases every N points (difficulty-dependent)
- Gap width varies by difficulty (9-15 units)
- Initial/min speed configurable per difficulty

**Collision Detection:**
- Bird Y position vs pipe gaps
- Screen boundaries (top/bottom)

## Data Persistence

High scores saved to `~/.flappy-bird-tui/`:
- `highscore.json` - All-time best score
- `rankings.json` - Top 10 scores with stats

Stats tracked (stored in `storage.HighScore`):
- Score, duration (elapsed time), date
- Jump count (total jumps during game)
- Max height (highest Y position reached)
- Min height (lowest Y position reached)
- Avg height (average Y position across samples)
- Difficulty level (Easy/Normal/Hard)

## Dependencies

Key packages (see `go.mod`):
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling

Bubble Tea handles:
- Keyboard input
- Terminal size detection
- Alt screen buffer
- Tick-based updates

## Code Conventions

- Standard Go formatting (`gofmt`)
- Exported types use doc comments
- Constants use iota for enums
- Error handling with explicit returns
- No external state - Model holds everything

## Testing Strategy

Currently minimal tests. Testable areas:
- Physics calculations (`domain/bird.go`)
- Collision detection (`domain/pipe.go`)
- Difficulty settings (`domain/difficulty.go`)
- High score logic (`storage/highscore.go`)
- Stats calculations (`game/model.go` - Stats struct)

Testing Bubble Tea components (`game/update.go`) requires mocking tea.Msg flows.

## Release Process

Uses goreleaser for multi-platform builds:
1. Ensure `.goreleaser.yml` is current
2. Tag version: `git tag vX.Y.Z`
3. Push tag: `git push origin vX.Y.Z`
4. GoReleaser runs (manual or CI): `goreleaser release --clean`

Builds for macOS, Linux, Windows (see goreleaser config).

## Common Patterns

**Adding a new difficulty:**
1. Add constant to `domain/difficulty.go`
2. Add case to `GetSettings()`
3. Add case to `String()`
4. Update title screen keybindings in `game/update.go`

**Adding a new theme:**
1. Add constant to `domain/theme.go`
2. Add case to `GetColors()`
3. Add case to `String()`
4. Theme cycles with `T` key (handled in `game/update.go`)

**Adding new sound events:**
1. Add case to `storage/sound.go` `PlaySound()`
2. Call `storage.PlaySound("event_name")` in `game/update.go`

**Adding a new statistic:**
1. Add field to `Stats` struct in `game/model.go`
2. Update calculation logic in `game/update.go`
3. Add corresponding field to `storage.HighScore` for persistence
4. Update display in `ui/view.go`

## Known Limitations

- Audio is terminal bell only (simple beeps via `\a` escape sequence)
- Fixed terminal size (80x24)
- Single player only
- No pause functionality
- No replay saving

## Performance Notes

- Target: ~16-50 FPS (60ms - 20ms tick intervals)
- Speed dynamically adjusts based on score
- Rendering is ASCII-based (low overhead)
- No goroutines - single event loop via Bubble Tea
