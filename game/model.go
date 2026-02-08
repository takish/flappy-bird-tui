package game

import (
	"time"

	"github.com/takish/flappy-bird-tui/domain"
	"github.com/takish/flappy-bird-tui/storage"
)

// GameState represents the current state of the game
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateGameOver
)

// Stats holds game statistics
type Stats struct {
	JumpCount     int // Number of jumps
	MaxHeight     int // Highest Y position reached
	MinHeight     int // Lowest Y position reached
	TotalHeight   int // Sum of heights for average calculation
	HeightSamples int // Number of height samples
}

// Model holds the entire game state
type Model struct {
	State       GameState
	Bird        *domain.Bird
	Pipes       []*domain.Pipe
	Score       int
	Width       int
	Height      int
	GameSpeed   time.Duration
	StartTime   time.Time // Game start time for elapsed time display
	HighScore   *storage.HighScore
	Rankings    []storage.HighScore // Top 10 rankings
	IsNewRecord bool                // Flag to indicate if current game is a new high score
	Difficulty  domain.Difficulty   // Current difficulty level
	Theme       domain.Theme        // Current color theme
	Stats       Stats               // Game statistics
	Err         error
}

// NewModel creates a new game with default values
func NewModel() Model {
	width := 80
	height := 24

	// Load high score
	highScore, err := storage.LoadHighScore()
	if err != nil {
		highScore = &storage.HighScore{} // Use empty high score on error
	}

	// Load rankings
	rankings, err := storage.LoadRankings()
	if err != nil {
		rankings = []storage.HighScore{} // Use empty rankings on error
	}

	return Model{
		State:      StateTitle,
		Bird:       domain.NewBird(10, height/2),
		Pipes:      []*domain.Pipe{},
		Score:      0,
		Width:      width,
		Height:     height,
		GameSpeed:  time.Millisecond * 45, // ~22 FPS - faster scroll speed
		HighScore:  highScore,
		Rankings:   rankings,
		Difficulty: domain.DifficultyNormal, // Default difficulty
		Theme:      domain.ThemeClassic,     // Default theme
	}
}

// AvgHeight calculates the average height from statistics
func (m Model) AvgHeight() float64 {
	if m.Stats.HeightSamples > 0 {
		return float64(m.Stats.TotalHeight) / float64(m.Stats.HeightSamples)
	}
	return 0.0
}

// resetGame resets the game to initial playing state
func (m Model) resetGame() Model {
	settings := m.Difficulty.GetSettings()

	return Model{
		State:       StatePlaying,
		Bird:        domain.NewBird(10, m.Height/2),
		Pipes:       []*domain.Pipe{}, // Start with no pipes - gives player time to adjust
		Score:       0,
		Width:       m.Width,
		Height:      m.Height,
		GameSpeed:   settings.InitialSpeed, // Use difficulty-based speed
		StartTime:   time.Now(),            // Record game start time
		HighScore:   m.HighScore,
		Rankings:    m.Rankings,
		IsNewRecord: false,
		Difficulty:  m.Difficulty,
		Theme:       m.Theme,
		Stats: Stats{
			JumpCount:     0,
			MaxHeight:     m.Height / 2,
			MinHeight:     m.Height / 2,
			TotalHeight:   0,
			HeightSamples: 0,
		},
	}
}
