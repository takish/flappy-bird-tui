package game

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/takish/flappy-bird-tui/domain"
	"github.com/takish/flappy-bird-tui/storage"
)

// TickMsg is sent on every game tick
type TickMsg time.Time

const (
	pipeSpawnGap = 50 // Horizontal gap between pipes
)

// Init initializes the game
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "1", "2", "3": // Difficulty selection (title screen only)
			if m.State == StateTitle {
				switch msg.String() {
				case "1":
					m.Difficulty = domain.DifficultyEasy
				case "2":
					m.Difficulty = domain.DifficultyNormal
				case "3":
					m.Difficulty = domain.DifficultyHard
				}
			}

		case "t": // Theme toggle (title screen only)
			if m.State == StateTitle {
				m.Theme = m.Theme.Next()
			}

		case " ": // Space bar
			switch m.State {
			case StateTitle:
				m = m.resetGame()
				return m, tick(m.GameSpeed)
			case StatePlaying:
				m.Bird.Jump()
				m.Stats.JumpCount++ // Track jump count
				storage.PlaySound("jump")
			case StateGameOver:
				m = m.resetGame()
				return m, tick(m.GameSpeed)
			}

		case "r":
			if m.State == StateGameOver {
				m = m.resetGame()
				return m, tick(m.GameSpeed)
			}
		}

	case TickMsg:
		if m.State != StatePlaying {
			return m, nil
		}

		// Update bird physics
		m.Bird.Update()

		// Track height statistics
		birdY := m.Bird.GetY()
		if birdY < m.Stats.MaxHeight {
			m.Stats.MaxHeight = birdY
		}
		if birdY > m.Stats.MinHeight {
			m.Stats.MinHeight = birdY
		}
		m.Stats.TotalHeight += birdY
		m.Stats.HeightSamples++

		// Check ceiling/floor collision
		if m.Bird.GetY() < 0 || m.Bird.GetY() >= m.Height {
			m.State = StateGameOver
			m = m.handleGameOver()
			return m, nil
		}

		// Update pipes
		for i := len(m.Pipes) - 1; i >= 0; i-- {
			pipe := m.Pipes[i]
			pipe.Update()

			// Check collision
			if pipe.CollidesWith(m.Bird) {
				m.State = StateGameOver
				m = m.handleGameOver()
				return m, nil
			}

			// Check if passed
			if pipe.IsPassed(m.Bird) {
				pipe.Passed = true
				m.Score++
				storage.PlaySound("score")

				// Increase difficulty based on settings
				settings := m.Difficulty.GetSettings()
				if m.Score%settings.ScoreInterval == 0 {
					// Speed up game (reduce interval)
					if m.GameSpeed > settings.MinSpeed {
						m.GameSpeed -= settings.SpeedIncrement
					}
				}
			}

			// Remove off-screen pipes
			if pipe.IsOffScreen() {
				m.Pipes = append(m.Pipes[:i], m.Pipes[i+1:]...)
			}
		}

		// Spawn new pipes
		if len(m.Pipes) == 0 || m.Pipes[len(m.Pipes)-1].X < m.Width-pipeSpawnGap {
			settings := m.Difficulty.GetSettings()
			m.Pipes = append(m.Pipes, domain.NewPipe(m.Width, m.Height, settings.PipeGap))
		}

		return m, tick(m.GameSpeed)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}

// handleGameOver processes game over logic including high score checking
func (m Model) handleGameOver() Model {
	storage.PlaySound("gameover")
	elapsed := time.Since(m.StartTime)

	// Create high score entry with statistics
	newScore := storage.HighScore{
		Score:      m.Score,
		Duration:   elapsed,
		Date:       time.Now(),
		JumpCount:  m.Stats.JumpCount,
		MaxHeight:  m.Stats.MaxHeight,
		MinHeight:  m.Stats.MinHeight,
		AvgHeight:  m.AvgHeight(),
		Difficulty: m.Difficulty.String(),
	}

	// Check if this is a new high score
	if storage.IsNewHighScore(m.Score, m.HighScore) {
		m.IsNewRecord = true
		// Save new high score
		if err := storage.SaveHighScore(newScore); err == nil {
			// Reload high score from disk
			if hs, err := storage.LoadHighScore(); err == nil {
				m.HighScore = hs
			}
		}
	}

	// Add to rankings
	newRankings, _ := storage.AddToRankings(m.Rankings, newScore)
	if err := storage.SaveRankings(newRankings); err == nil {
		m.Rankings = newRankings
	}

	return m
}

// tick returns a command that waits for the specified duration and sends a TickMsg
func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
