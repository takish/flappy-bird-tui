package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is sent on every game tick
type TickMsg time.Time

// Init initializes the game
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case " ": // Space bar
			switch m.state {
			case StateTitle:
				m = m.resetGame()
				return m, tick(m.gameSpeed)
			case StatePlaying:
				m.bird.Jump()
			case StateGameOver:
				m = m.resetGame()
				return m, tick(m.gameSpeed)
			}

		case "r":
			if m.state == StateGameOver {
				m = m.resetGame()
				return m, tick(m.gameSpeed)
			}
		}

	case TickMsg:
		if m.state != StatePlaying {
			return m, nil
		}

		// Update bird physics
		m.bird.Update()

		// Check ceiling/floor collision
		if m.bird.GetY() < 0 || m.bird.GetY() >= m.height {
			m.state = StateGameOver
			return m, nil
		}

		// Update pipes
		for i := len(m.pipes) - 1; i >= 0; i-- {
			pipe := m.pipes[i]
			pipe.Update()

			// Check collision
			if pipe.CollidesWith(m.bird) {
				m.state = StateGameOver
				return m, nil
			}

			// Check if passed
			if pipe.IsPassed(m.bird) {
				pipe.passed = true
				m.score++
			}

			// Remove off-screen pipes
			if pipe.IsOffScreen() {
				m.pipes = append(m.pipes[:i], m.pipes[i+1:]...)
			}
		}

		// Spawn new pipes
		if len(m.pipes) == 0 || m.pipes[len(m.pipes)-1].x < m.width-30 {
			m.pipes = append(m.pipes, NewPipe(m.width, m.height))
		}

		return m, tick(m.gameSpeed)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// tick returns a command that waits for the specified duration and sends a TickMsg
func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
