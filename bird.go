package main

const (
	gravity   = 0.3  // Reduced from 0.8 - gentler fall
	jumpForce = -2.3 // Increased from -1.8 - more responsive jump
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

// GetSprite returns the bird's sprite based on velocity
func (b *Bird) GetSprite() string {
	if b.velocity < 0 {
		// Going up
		return "^○"
	}
	// Going down
	return "v○"
}
