package domain

import "time"

// Difficulty represents the game difficulty level
type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyNormal
	DifficultyHard
)

// DifficultySettings holds the parameters for each difficulty level
type DifficultySettings struct {
	InitialSpeed   time.Duration // Initial game speed
	SpeedIncrement time.Duration // How much to speed up
	ScoreInterval  int           // Score interval for speed increase
	MinSpeed       time.Duration // Minimum (maximum) speed
	PipeGap        int           // Gap between top and bottom pipes
}

// GetSettings returns the settings for a difficulty level
func (d Difficulty) GetSettings() DifficultySettings {
	switch d {
	case DifficultyEasy:
		return DifficultySettings{
			InitialSpeed:   time.Millisecond * 60, // Slower
			SpeedIncrement: time.Millisecond * 5,  // Gentler acceleration
			ScoreInterval:  5,                     // Every 5 points
			MinSpeed:       time.Millisecond * 30, // Not too fast
			PipeGap:        15,                    // Wider gap
		}
	case DifficultyHard:
		return DifficultySettings{
			InitialSpeed:   time.Millisecond * 30, // Faster
			SpeedIncrement: time.Millisecond * 10, // Aggressive acceleration
			ScoreInterval:  2,                     // Every 2 points
			MinSpeed:       time.Millisecond * 10, // Very fast
			PipeGap:        9,                     // Narrower gap
		}
	default: // DifficultyNormal
		return DifficultySettings{
			InitialSpeed:   time.Millisecond * 45, // Current default
			SpeedIncrement: time.Millisecond * 8,
			ScoreInterval:  3,
			MinSpeed:       time.Millisecond * 20,
			PipeGap:        12,
		}
	}
}

// String returns the string representation of the difficulty
func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Easy"
	case DifficultyHard:
		return "Hard"
	default:
		return "Normal"
	}
}
