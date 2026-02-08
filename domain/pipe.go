package domain

import (
	"math/rand/v2"
)

const (
	PipeWidth = 8 // Width of pipes - doubled for thicker dokan
	minPipeY  = 3 // Minimum gap from top
)

// Pipe represents an obstacle
type Pipe struct {
	X       int
	GapY    int // Y position of the gap's top
	GapSize int // Size of the gap
	Passed  bool
}

// NewPipe creates a new pipe at the right edge of the screen
func NewPipe(screenWidth, screenHeight, gapSize int) *Pipe {
	// Random gap position, ensuring gap fits within screen
	maxGapY := screenHeight - gapSize - minPipeY
	gapY := rand.IntN(maxGapY-minPipeY) + minPipeY

	return &Pipe{
		X:       screenWidth,
		GapY:    gapY,
		GapSize: gapSize,
		Passed:  false,
	}
}

// Update moves the pipe to the left
func (p *Pipe) Update() {
	p.X--
}

// IsOffScreen checks if the pipe has moved off the left edge
func (p *Pipe) IsOffScreen() bool {
	return p.X+PipeWidth < 0
}

// CollidesWith checks if the bird collides with this pipe
func (p *Pipe) CollidesWith(bird *Bird) bool {
	birdY := bird.GetY()

	// Bird is 2 characters wide (mO or wO), check both positions
	for birdX := bird.X; birdX < bird.X+2; birdX++ {
		// Check if bird is horizontally aligned with pipe
		if birdX >= p.X && birdX < p.X+PipeWidth {
			// Check if bird is outside the gap
			if birdY < p.GapY || birdY >= p.GapY+p.GapSize {
				return true
			}
		}
	}

	return false
}

// IsPassed checks if the bird has passed this pipe
func (p *Pipe) IsPassed(bird *Bird) bool {
	// Bird is 2 characters wide, check the right edge (bird.x+1)
	return bird.X+1 > p.X+PipeWidth && !p.Passed
}
