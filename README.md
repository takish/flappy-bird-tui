# Flappy Bird TUI

A terminal user interface (TUI) version of Flappy Bird built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

### Gameplay
- Classic Flappy Bird gameplay in your terminal
- Progressive difficulty - game speeds up as you score
- ASCII art graphics with retro charm
- Smooth ~16-50 FPS gameplay
- Cross-platform (macOS, Linux, Windows)

### Customization
- **3 Difficulty Levels**
  - Easy: Slower pace, wider gaps, gentle acceleration
  - Normal: Balanced challenge (default)
  - Hard: Fast pace, narrow gaps, aggressive acceleration
- **3 Color Themes**
  - Classic: Blue and green (default)
  - Retro: Green terminal aesthetic
  - Neon: Magenta and cyan cyberpunk style

### Progress Tracking
- High score persistence
- Elapsed time tracking (MM:SS.mmm)
- New record celebration

### Audio
- Terminal beep sounds for:
  - Jump action
  - Score increments
  - Game over

## Installation

### Homebrew (Recommended)

```bash
brew tap takish/tap
brew install flappy-bird-tui
```

### From Source

```bash
git clone https://github.com/takish/flappy-bird-tui.git
cd flappy-bird-tui
go build
./flappy-bird-tui
```

### Go Install

```bash
go install github.com/takish/flappy-bird-tui@latest
```

## How to Play

### Controls
- **Space** - Jump / Start game
- **1, 2, 3** - Select difficulty (title screen)
- **T** - Change theme (title screen)
- **Q** - Quit game

### Tips
- Start slow to get familiar with the physics
- The game speeds up every few points - stay focused!
- Try different difficulties to find your sweet spot
- Your high score is saved automatically

## Architecture

Built using the Elm Architecture (Model-Update-View) pattern via Bubble Tea:

- `main.go` - Entry point
- `model.go` - Game state management
- `update.go` - Input and tick handling
- `view.go` - ASCII rendering with theming
- `bird.go` - Bird physics (gravity, jump)
- `pipe.go` - Pipe generation, movement, and collision detection
- `difficulty.go` - Difficulty level settings
- `theme.go` - Color theme definitions
- `highscore.go` - High score persistence
- `sound.go` - Audio feedback

## Development

### Build
```bash
go build -o flappy-bird-tui
```

### Test
```bash
go test ./...
```

### Release
```bash
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

## License

MIT
