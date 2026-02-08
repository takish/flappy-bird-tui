package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	scoreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	gameOverStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))

	birdChar = ">"
	pipeChar = "█"
)

// View renders the current game state
func (m Model) View() string {
	if m.width < 40 || m.height < 10 {
		return "Terminal too small! Please resize to at least 40x10"
	}

	switch m.state {
	case StateTitle:
		return m.renderTitle()
	case StatePlaying:
		return m.renderGame()
	case StateGameOver:
		return m.renderGameOver()
	}

	return ""
}

func (m Model) renderTitle() string {
	var b strings.Builder

	// Center title
	title := titleStyle.Render("FLAPPY BIRD TUI")
	instructions := "Press SPACE to start\nPress Q to quit"

	padding := (m.height - 5) / 2
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}

	b.WriteString(centerText(title, m.width))
	b.WriteString("\n\n")
	b.WriteString(centerText(instructions, m.width))

	return b.String()
}

func (m Model) renderGame() string {
	// Create empty canvas
	canvas := make([][]rune, m.height)
	for i := range canvas {
		canvas[i] = make([]rune, m.width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	// Draw pipes
	for _, pipe := range m.pipes {
		for x := pipe.x; x < pipe.x+pipeWidth && x < m.width; x++ {
			if x < 0 {
				continue
			}
			// Top pipe
			for y := 0; y < pipe.gapY; y++ {
				canvas[y][x] = []rune(pipeChar)[0]
			}
			// Bottom pipe
			for y := pipe.gapY + pipeGap; y < m.height; y++ {
				canvas[y][x] = []rune(pipeChar)[0]
			}
		}
	}

	// Draw bird
	birdY := m.bird.GetY()
	if birdY >= 0 && birdY < m.height && m.bird.x >= 0 && m.bird.x < m.width {
		canvas[birdY][m.bird.x] = []rune(birdChar)[0]
	}

	// Convert canvas to string
	var b strings.Builder
	for _, row := range canvas {
		b.WriteString(string(row))
		b.WriteString("\n")
	}

	// Add score
	score := scoreStyle.Render(fmt.Sprintf("Score: %d", m.score))
	b.WriteString(score)

	return b.String()
}

func (m Model) renderGameOver() string {
	var b strings.Builder

	padding := (m.height - 6) / 2
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}

	gameOver := gameOverStyle.Render("GAME OVER")
	score := fmt.Sprintf("Final Score: %d", m.score)
	instructions := "Press SPACE or R to restart\nPress Q to quit"

	b.WriteString(centerText(gameOver, m.width))
	b.WriteString("\n\n")
	b.WriteString(centerText(score, m.width))
	b.WriteString("\n\n")
	b.WriteString(centerText(instructions, m.width))

	return b.String()
}

func centerText(text string, width int) string {
	lines := strings.Split(text, "\n")
	var centered strings.Builder

	for i, line := range lines {
		// Remove ANSI codes for length calculation
		plainLen := lipgloss.Width(line)
		padding := (width - plainLen) / 2
		if padding < 0 {
			padding = 0
		}

		centered.WriteString(strings.Repeat(" ", padding))
		centered.WriteString(line)
		if i < len(lines)-1 {
			centered.WriteString("\n")
		}
	}

	return centered.String()
}
