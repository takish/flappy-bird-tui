package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/takish/flappy-bird-tui/game"
	"github.com/takish/flappy-bird-tui/storage"
	"github.com/takish/flappy-bird-tui/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// modelWrapper wraps game.Model to provide the View() method
type modelWrapper struct {
	game.Model
}

// Init initializes the game
func (m modelWrapper) Init() tea.Cmd {
	return m.Model.Init()
}

// Update handles messages and updates the model
func (m modelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := m.Model.Update(msg)
	return modelWrapper{newModel}, cmd
}

// View renders the current game state
func (m modelWrapper) View() string {
	return ui.View(m.Model)
}

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("flappy-bird-tui %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	// Initialize sound system
	if err := storage.InitSound(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize sound: %v\n", err)
		// Continue without sound - not a fatal error
	}

	p := tea.NewProgram(modelWrapper{game.NewModel()}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
