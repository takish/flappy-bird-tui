package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// HighScore represents a saved high score with statistics
type HighScore struct {
	Score      int           `json:"score"`
	Duration   time.Duration `json:"duration"`
	Date       time.Time     `json:"date"`
	JumpCount  int           `json:"jump_count"`
	MaxHeight  int           `json:"max_height"`
	MinHeight  int           `json:"min_height"`
	AvgHeight  float64       `json:"avg_height"`
	Difficulty string        `json:"difficulty"`
}

const configDir = ".flappy-bird-tui"
const highScoreFile = "highscore.json"
const rankingsFile = "rankings.json"
const maxRankings = 10

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
func SaveHighScore(hs HighScore) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
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

// LoadRankings loads the top 10 rankings from disk
func LoadRankings() ([]HighScore, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(configPath, rankingsFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No rankings yet
			return []HighScore{}, nil
		}
		return nil, err
	}

	var rankings []HighScore
	if err := json.Unmarshal(data, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

// SaveRankings saves the rankings to disk
func SaveRankings(rankings []HighScore) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(rankings, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(configPath, rankingsFile)
	return os.WriteFile(filePath, data, 0644)
}

// AddToRankings adds a new score to the rankings and returns the rank (1-based)
// Returns 0 if the score didn't make it to the rankings
func AddToRankings(rankings []HighScore, newScore HighScore) ([]HighScore, int) {
	// Add new score
	rankings = append(rankings, newScore)

	// Sort by score (descending), then by duration (ascending for same score)
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Score == rankings[j].Score {
			return rankings[i].Duration < rankings[j].Duration
		}
		return rankings[i].Score > rankings[j].Score
	})

	// Find rank of new score
	rank := 0
	for i, hs := range rankings {
		if hs.Date == newScore.Date && hs.Score == newScore.Score {
			rank = i + 1
			break
		}
	}

	// Keep only top 10
	if len(rankings) > maxRankings {
		rankings = rankings[:maxRankings]
		// If rank is beyond top 10, return 0
		if rank > maxRankings {
			rank = 0
		}
	}

	return rankings, rank
}
