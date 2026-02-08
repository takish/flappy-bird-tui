package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// HighScore represents a saved high score
type HighScore struct {
	Score    int           `json:"score"`
	Duration time.Duration `json:"duration"`
	Date     time.Time     `json:"date"`
}

const configDir = ".flappy-bird-tui"
const highScoreFile = "highscore.json"

// getConfigPath returns the path to the config directory
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir), nil
}

// LoadHighScore loads the high score from disk
func LoadHighScore() (*HighScore, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(configPath, highScoreFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No high score yet
			return &HighScore{}, nil
		}
		return nil, err
	}

	var hs HighScore
	if err := json.Unmarshal(data, &hs); err != nil {
		return nil, err
	}

	return &hs, nil
}

// SaveHighScore saves the high score to disk
func SaveHighScore(score int, duration time.Duration) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	hs := HighScore{
		Score:    score,
		Duration: duration,
		Date:     time.Now(),
	}

	data, err := json.MarshalIndent(hs, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(configPath, highScoreFile)
	return os.WriteFile(filePath, data, 0644)
}

// IsNewHighScore checks if the current score is a new high score
func IsNewHighScore(currentScore int, highScore *HighScore) bool {
	return currentScore > highScore.Score
}
