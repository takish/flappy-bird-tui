# Flappy Bird TUI

A terminal user interface (TUI) version of Flappy Bird built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Classic Flappy Bird gameplay in your terminal
- ASCII art graphics
- 60 FPS gameplay
- Cross-platform (macOS, Linux, Windows)

## Installation

```bash
go install github.com/takish/flappy-bird-tui@latest
```

Or build from source:

```bash
git clone https://github.com/takish/flappy-bird-tui.git
cd flappy-bird-tui
go build
./flappy-bird-tui
```

## How to Play

- Press **Space** to flap
- Press **q** to quit
- Avoid the pipes and stay alive!

## Architecture

Built using the Elm Architecture (Model-Update-View) pattern via Bubble Tea:

- `model.go` - Game state
- `update.go` - Input and tick handling
- `view.go` - ASCII rendering
- `bird.go` - Bird physics (gravity, jump)
- `pipe.go` - Pipe generation, movement, and collision detection

## License

MIT
