package main

const (
	gravity  = 0.8
	jumpForce = -3.0
)

// Bird represents the player character
type Bird struct {
	x        int
	y        float64
	velocity float64
}

// NewBird creates a new bird at the given position
func NewBird(x int, y int) *Bird {
	return &Bird{
		x:        x,
		y:        float64(y),
		velocity: 0,
	}
}

// Jump makes the bird flap upward
func (b *Bird) Jump() {
	b.velocity = jumpForce
}

// Update applies gravity and updates position
func (b *Bird) Update() {
	b.velocity += gravity
	b.y += b.velocity
}

// GetY returns the bird's current Y position as an integer
func (b *Bird) GetY() int {
	return int(b.y)
}
