# CLAUDE.md

This file provides guidance to Claude Code when working with this Flappy Bird TUI project.

## Overview

Terminal-based Flappy Bird clone built with Go, Bubble Tea (TUI framework), and Lipgloss (styling). Implements the Elm Architecture (TEA) pattern with Model-Update-View separation.

## Quick Commands

```bash
# Build and run
go build && ./flappy-bird-tui

# Run directly
go run .

# Test
go test ./...

# Release (requires goreleaser)
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean

# Version info
./flappy-bird-tui --version
```

## Architecture

### Elm Architecture Pattern (TEA)

Built on Bubble Tea's TEA implementation:

1. **Model** (`model.go`) - Game state container
2. **Update** (`update.go`) - State transitions from messages
3. **View** (`view.go`) - Render ASCII representation

### Core Components

| File | Purpose | Key Types/Functions |
|------|---------|---------------------|
| `main.go` | Entry point, CLI flags | `main()`, version vars |
| `model.go` | Game state | `Model`, `GameState`, `initialModel()` |
| `update.go` | Input/tick handling | `Update()`, `Init()` |
| `view.go` | ASCII rendering | `View()`, theme-aware rendering |
| `bird.go` | Bird physics | `Bird`, gravity, jump |
| `pipe.go` | Pipe system | `Pipe`, generation, collision |
| `difficulty.go` | Difficulty settings | `Difficulty`, `DifficultySettings` |
| `theme.go` | Color themes | `Theme`, `ColorScheme` |
| `highscore.go` | Persistence | JSON save/load to `~/.flappy-bird-tui/` |
| `sound.go` | Terminal beeps | `PlaySound()` |

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
- `highscore.json` - All-time best
- `rankings.json` - Top 10 scores with stats

Stats tracked:
- Score, duration, date
- Jump count, min/max/avg height
- Difficulty level

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
- Physics calculations (`bird.go`)
- Collision detection (`pipe.go`)
- Difficulty settings (`difficulty.go`)
- High score logic (`highscore.go`)

Testing Bubble Tea components requires mocking tea.Msg flows.

## Release Process

Uses goreleaser for multi-platform builds:
1. Ensure `.goreleaser.yml` is current
2. Tag version: `git tag vX.Y.Z`
3. Push tag: `git push origin vX.Y.Z`
4. GoReleaser runs (manual or CI): `goreleaser release --clean`

Builds for macOS, Linux, Windows (see goreleaser config).

## Common Patterns

**Adding a new difficulty:**
1. Add constant to `difficulty.go`
2. Add case to `GetSettings()`
3. Add case to `String()`
4. Update title screen keybindings in `update.go`

**Adding a new theme:**
1. Add constant to `theme.go`
2. Add case to `GetColors()`
3. Add case to `String()`
4. Theme cycles with `T` key (handled in `update.go`)

**Adding new sound events:**
1. Add case to `sound.go` `PlaySound()`
2. Call `PlaySound("event_name")` in `update.go`

## Known Limitations

- Audio is terminal bell only (no external sounds)
- Fixed terminal size (80x24)
- Single player only
- No pause functionality
- No replay saving

## Performance Notes

- Target: ~16-50 FPS (60ms - 20ms tick intervals)
- Speed dynamically adjusts based on score
- Rendering is ASCII-based (low overhead)
- No goroutines - single event loop via Bubble Tea
