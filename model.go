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
	state      GameState
	bird       *Bird
	pipes      []*Pipe
	score      int
	width      int
	height     int
	gameSpeed  time.Duration
	err        error
}

// initialModel creates a new game with default values
func initialModel() Model {
	width := 80
	height := 24

	return Model{
		state:     StateTitle,
		bird:      NewBird(10, height/2),
		pipes:     []*Pipe{},
		score:     0,
		width:     width,
		height:    height,
		gameSpeed: time.Millisecond * 50, // 20 FPS for smooth gameplay
	}
}

// resetGame resets the game to initial playing state
func (m Model) resetGame() Model {
	return Model{
		state:     StatePlaying,
		bird:      NewBird(10, m.height/2),
		pipes:     []*Pipe{NewPipe(m.width, m.height)},
		score:     0,
		width:     m.width,
		height:    m.height,
		gameSpeed: m.gameSpeed,
	}
}
