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
	isNewRecord  bool       // Flag to indicate if current game is a new high score
	difficulty   Difficulty // Current difficulty level
	theme        Theme      // Current color theme
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

	return Model{
		state:      StateTitle,
		bird:       NewBird(10, height/2),
		pipes:      []*Pipe{},
		score:      0,
		width:      width,
		height:     height,
		gameSpeed:  time.Millisecond * 45, // ~22 FPS - faster scroll speed
		highScore:  highScore,
		difficulty: DifficultyNormal, // Default difficulty
		theme:      ThemeClassic,     // Default theme
	}
}

// resetGame resets the game to initial playing state
func (m Model) resetGame() Model {
	settings := m.difficulty.GetSettings()

	return Model{
		state:       StatePlaying,
		bird:        NewBird(10, m.height/2),
		pipes:       []*Pipe{}, // Start with no pipes - gives player time to adjust
		score:       0,
		width:       m.width,
		height:      m.height,
		gameSpeed:   settings.InitialSpeed, // Use difficulty-based speed
		startTime:   time.Now(),            // Record game start time
		highScore:   m.highScore,
		isNewRecord: false,
		difficulty:  m.difficulty,
		theme:       m.theme,
	}
}
