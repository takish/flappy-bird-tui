package main

import (
	"math/rand/v2"
)

const (
	pipeWidth = 5 // Width of pipes - increased from 4 for thicker blocks
	minPipeY  = 3 // Minimum gap from top
)

// Pipe represents an obstacle
type Pipe struct {
	x       int
	gapY    int // Y position of the gap's top
	gapSize int // Size of the gap
	passed  bool
}

// NewPipe creates a new pipe at the right edge of the screen
func NewPipe(screenWidth, screenHeight, gapSize int) *Pipe {
	// Random gap position, ensuring gap fits within screen
	maxGapY := screenHeight - gapSize - minPipeY
	gapY := rand.IntN(maxGapY-minPipeY) + minPipeY

	return &Pipe{
		x:       screenWidth,
		gapY:    gapY,
		gapSize: gapSize,
		passed:  false,
	}
}

// Update moves the pipe to the left
func (p *Pipe) Update() {
	p.x--
}

// IsOffScreen checks if the pipe has moved off the left edge
func (p *Pipe) IsOffScreen() bool {
	return p.x+pipeWidth < 0
}

// CollidesWith checks if the bird collides with this pipe
func (p *Pipe) CollidesWith(bird *Bird) bool {
	birdY := bird.GetY()

	// Check if bird is horizontally aligned with pipe
	if bird.x >= p.x && bird.x < p.x+pipeWidth {
		// Check if bird is outside the gap
		if birdY < p.gapY || birdY >= p.gapY+p.gapSize {
			return true
		}
	}

	return false
}

// IsPassed checks if the bird has passed this pipe
func (p *Pipe) IsPassed(bird *Bird) bool {
	return bird.x > p.x+pipeWidth && !p.passed
}
