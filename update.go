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

		case "1", "2", "3": // Difficulty selection (title screen only)
			if m.state == StateTitle {
				switch msg.String() {
				case "1":
					m.difficulty = DifficultyEasy
				case "2":
					m.difficulty = DifficultyNormal
				case "3":
					m.difficulty = DifficultyHard
				}
			}

		case "t": // Theme toggle (title screen only)
			if m.state == StateTitle {
				m.theme = m.theme.Next()
			}

		case " ": // Space bar
			switch m.state {
			case StateTitle:
				m = m.resetGame()
				return m, tick(m.gameSpeed)
			case StatePlaying:
				m.bird.Jump()
				m.jumpCount++ // Track jump count
				PlaySound("jump")
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

		// Track height statistics
		birdY := m.bird.GetY()
		if birdY < m.maxHeight {
			m.maxHeight = birdY
		}
		if birdY > m.minHeight {
			m.minHeight = birdY
		}
		m.totalHeight += birdY
		m.heightSamples++

		// Check ceiling/floor collision
		if m.bird.GetY() < 0 || m.bird.GetY() >= m.height {
			m.state = StateGameOver
			m = m.handleGameOver()
			return m, nil
		}

		// Update pipes
		for i := len(m.pipes) - 1; i >= 0; i-- {
			pipe := m.pipes[i]
			pipe.Update()

			// Check collision
			if pipe.CollidesWith(m.bird) {
				m.state = StateGameOver
				m = m.handleGameOver()
				return m, nil
			}

			// Check if passed
			if pipe.IsPassed(m.bird) {
				pipe.passed = true
				m.score++
				PlaySound("score")

				// Increase difficulty based on settings
				settings := m.difficulty.GetSettings()
				if m.score%settings.ScoreInterval == 0 {
					// Speed up game (reduce interval)
					if m.gameSpeed > settings.MinSpeed {
						m.gameSpeed -= settings.SpeedIncrement
					}
				}
			}

			// Remove off-screen pipes
			if pipe.IsOffScreen() {
				m.pipes = append(m.pipes[:i], m.pipes[i+1:]...)
			}
		}

		// Spawn new pipes
		if len(m.pipes) == 0 || m.pipes[len(m.pipes)-1].x < m.width-50 {
			settings := m.difficulty.GetSettings()
			m.pipes = append(m.pipes, NewPipe(m.width, m.height, settings.PipeGap))
		}

		return m, tick(m.gameSpeed)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// handleGameOver processes game over logic including high score checking
func (m Model) handleGameOver() Model {
	PlaySound("gameover")
	elapsed := time.Since(m.startTime)

	// Calculate average height
	avgHeight := 0.0
	if m.heightSamples > 0 {
		avgHeight = float64(m.totalHeight) / float64(m.heightSamples)
	}

	// Create high score entry with statistics
	newScore := HighScore{
		Score:      m.score,
		Duration:   elapsed,
		Date:       time.Now(),
		JumpCount:  m.jumpCount,
		MaxHeight:  m.maxHeight,
		MinHeight:  m.minHeight,
		AvgHeight:  avgHeight,
		Difficulty: m.difficulty.String(),
	}

	// Check if this is a new high score
	if IsNewHighScore(m.score, m.highScore) {
		m.isNewRecord = true
		// Save new high score
		if err := SaveHighScore(newScore); err == nil {
			// Reload high score from disk
			if hs, err := LoadHighScore(); err == nil {
				m.highScore = hs
			}
		}
	}

	// Add to rankings
	newRankings, _ := AddToRankings(m.rankings, newScore)
	if err := SaveRankings(newRankings); err == nil {
		m.rankings = newRankings
	}

	return m
}

// tick returns a command that waits for the specified duration and sends a TickMsg
func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
