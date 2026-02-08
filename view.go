package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	birdChar = "●" // Round bird - more visible
	pipeChar = "▓" // Block pattern - Mario-style
)

// getStyles returns the styles for the current theme
func (m Model) getStyles() (titleStyle, scoreStyle, gameOverStyle, newRecordStyle lipgloss.Style) {
	colors := m.theme.GetColors()

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colors.Title)

	scoreStyle = lipgloss.NewStyle().
		Foreground(colors.Score)

	gameOverStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colors.GameOver)

	newRecordStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colors.NewRecord)

	return
}

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

	// ASCII Art for FLAPPY
	flappyArt := `███████╗██╗      █████╗ ██████╗ ██████╗ ██╗   ██╗
██╔════╝██║     ██╔══██╗██╔══██╗██╔══██╗╚██╗ ██╔╝
█████╗  ██║     ███████║██████╔╝██████╔╝ ╚████╔╝
██╔══╝  ██║     ██╔══██║██╔═══╝ ██╔═══╝   ╚██╔╝
██║     ███████╗██║  ██║██║     ██║        ██║
╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝        ╚═╝   `

	// ASCII Art for BIRD
	birdArt := `██████╗ ██╗██████╗ ██████╗
██╔══██╗██║██╔══██╗██╔══██╗
██████╔╝██║██████╔╝██║  ██║
██╔══██╗██║██╔══██╗██║  ██║
██████╔╝██║██║  ██║██████╔╝
╚═════╝ ╚═╝╚═╝  ╚═╝╚═════╝ `

	subtitle := "T U I"

	// Difficulty selection with highlighted current difficulty
	difficultyText := "Difficulty: "
	if m.difficulty == DifficultyEasy {
		difficultyText += "[1: Easy*]  2: Normal  3: Hard"
	} else if m.difficulty == DifficultyHard {
		difficultyText += "1: Easy  2: Normal  [3: Hard*]"
	} else {
		difficultyText += "1: Easy  [2: Normal*]  3: Hard"
	}

	instructions := "Press SPACE to start  |  Press Q to quit"

	padding := (m.height - 16) / 2
	if padding < 0 {
		padding = 0
	}
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}

	// Get theme styles
	titleStyle, scoreStyle, _, _ := m.getStyles()

	// Display ASCII art title
	flappyStyled := titleStyle.Render(flappyArt)
	birdStyled := titleStyle.Render(birdArt)

	b.WriteString(centerText(flappyStyled, m.width))
	b.WriteString("\n")
	b.WriteString(centerText(birdStyled, m.width))
	b.WriteString("\n")
	b.WriteString(centerText(subtitle, m.width))
	b.WriteString("\n\n")
	b.WriteString(centerText(difficultyText, m.width))
	b.WriteString("\n")

	// Display theme selection
	themeText := fmt.Sprintf("Theme: %s (Press T to change)", m.theme.String())
	b.WriteString(centerText(themeText, m.width))
	b.WriteString("\n")

	b.WriteString(centerText(instructions, m.width))

	// Display high score if it exists
	if m.highScore.Score > 0 {
		b.WriteString("\n\n")
		minutes := int(m.highScore.Duration.Minutes())
		seconds := int(m.highScore.Duration.Seconds()) % 60
		milliseconds := int(m.highScore.Duration.Milliseconds()) % 1000
		highScoreText := fmt.Sprintf("High Score: %d  |  Time: %02d:%02d.%03d",
			m.highScore.Score, minutes, seconds, milliseconds)
		b.WriteString(centerText(scoreStyle.Render(highScoreText), m.width))
	}

	return b.String()
}

func (m Model) renderGame() string {
	// Get theme styles
	_, scoreStyle, _, _ := m.getStyles()

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
			for y := pipe.gapY + pipe.gapSize; y < m.height; y++ {
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

	// Add score and elapsed time
	elapsed := time.Since(m.startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	milliseconds := int(elapsed.Milliseconds()) % 1000

	score := scoreStyle.Render(fmt.Sprintf("Score: %d", m.score))
	timeDisplay := scoreStyle.Render(fmt.Sprintf("Time: %02d:%02d.%03d", minutes, seconds, milliseconds))

	b.WriteString(score)
	b.WriteString("  ")
	b.WriteString(timeDisplay)

	return b.String()
}

func (m Model) renderGameOver() string {
	// Get theme styles
	_, _, gameOverStyle, newRecordStyle := m.getStyles()

	var b strings.Builder

	// ASCII Art for GAME OVER
	asciiArt := ` ██████╗  █████╗ ███╗   ███╗███████╗     ██████╗ ██╗   ██╗███████╗██████╗
██╔════╝ ██╔══██╗████╗ ████║██╔════╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗
██║  ███╗███████║██╔████╔██║█████╗      ██║   ██║██║   ██║█████╗  ██████╔╝
██║   ██║██╔══██║██║╚██╔╝██║██╔══╝      ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗
╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗    ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝`

	padding := (m.height - 12) / 2
	if padding < 0 {
		padding = 0
	}
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}

	// Display ASCII art
	gameOverArt := gameOverStyle.Render(asciiArt)
	b.WriteString(centerText(gameOverArt, m.width))
	b.WriteString("\n\n")

	// Calculate elapsed time
	elapsed := time.Since(m.startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	milliseconds := int(elapsed.Milliseconds()) % 1000

	score := fmt.Sprintf("Final Score: %d  |  Time: %02d:%02d.%03d", m.score, minutes, seconds, milliseconds)
	instructions := "Press SPACE or R to restart  |  Press Q to quit"

	// Display new record message if applicable
	if m.isNewRecord {
		newRecord := newRecordStyle.Render("★ NEW RECORD! ★")
		b.WriteString(centerText(newRecord, m.width))
		b.WriteString("\n")
	}

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
