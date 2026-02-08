package main

import (
	"time"
)

// GameState represents the current state of the game
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateGameOver
)

// Model holds the entire game state
type Model struct {
	state        GameState
	bird         *Bird
	pipes        []*Pipe
	score        int
	width        int
	height       int
	gameSpeed    time.Duration
	startTime    time.Time // Game start time for elapsed time display
	highScore    *HighScore
	rankings     []HighScore // Top 10 rankings
	isNewRecord  bool        // Flag to indicate if current game is a new high score
	difficulty   Difficulty  // Current difficulty level
	theme        Theme       // Current color theme
	// Statistics
	jumpCount    int // Number of jumps
	maxHeight    int // Highest Y position reached
	minHeight    int // Lowest Y position reached
	totalHeight  int // Sum of heights for average calculation
	heightSamples int // Number of height samples
	err          error
}

// initialModel creates a new game with default values
func initialModel() Model {
	width := 80
	height := 24

	// Load high score
	highScore, err := LoadHighScore()
	if err != nil {
		highScore = &HighScore{} // Use empty high score on error
	}

	// Load rankings
	rankings, err := LoadRankings()
	if err != nil {
		rankings = []HighScore{} // Use empty rankings on error
	}

	return Model{
		state:      StateTitle,
		bird:       NewBird(10, height/2),
		pipes:      []*Pipe{},
		score:      0,
		width:      width,
		height:     height,
		gameSpeed:  time.Millisecond * 45, // ~22 FPS - faster scroll speed
		highScore:  highScore,
		rankings:   rankings,
		difficulty: DifficultyNormal, // Default difficulty
		theme:      ThemeClassic,     // Default theme
	}
}

// resetGame resets the game to initial playing state
func (m Model) resetGame() Model {
	settings := m.difficulty.GetSettings()

	return Model{
		state:         StatePlaying,
		bird:          NewBird(10, m.height/2),
		pipes:         []*Pipe{}, // Start with no pipes - gives player time to adjust
		score:         0,
		width:         m.width,
		height:        m.height,
		gameSpeed:     settings.InitialSpeed, // Use difficulty-based speed
		startTime:     time.Now(),            // Record game start time
		highScore:     m.highScore,
		rankings:      m.rankings,
		isNewRecord:   false,
		difficulty:    m.difficulty,
		theme:         m.theme,
		// Initialize statistics
		jumpCount:     0,
		maxHeight:     m.height / 2,
		minHeight:     m.height / 2,
		totalHeight:   0,
		heightSamples: 0,
	}
}
