package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/takish/flappy-bird-tui/domain"
	"github.com/takish/flappy-bird-tui/game"
)

const (
	birdChar         = "●" // Round bird - more visible
	pipeBodyChar     = "▓" // Pipe body - dark pattern block
	pipeEdgeChar     = "█" // Pipe edge - full block
	titlePadding     = 16  // Vertical padding for title screen ASCII art
	gameOverPadding  = 12  // Vertical padding for game over screen ASCII art
)

// formatDuration formats a duration as MM:SS.mmm
func formatDuration(d time.Duration) (minutes, seconds, milliseconds int) {
	minutes = int(d.Minutes())
	seconds = int(d.Seconds()) % 60
	milliseconds = int(d.Milliseconds()) % 1000
	return
}

// View renders the current game state
func View(m game.Model) string {
	if m.Width < 40 || m.Height < 10 {
		return "Terminal too small! Please resize to at least 40x10"
	}

	switch m.State {
	case game.StateTitle:
		return renderTitle(m)
	case game.StatePlaying:
		return renderGame(m)
	case game.StateGameOver:
		return renderGameOver(m)
	}

	return ""
}

// getStyles returns the styles for the current theme
func getStyles(m game.Model) (titleStyle, scoreStyle, gameOverStyle, newRecordStyle lipgloss.Style) {
	colors := m.Theme.GetColors()

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

func renderTitle(m game.Model) string {
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
	if m.Difficulty == domain.DifficultyEasy {
		difficultyText += "[1: Easy*]  2: Normal  3: Hard"
	} else if m.Difficulty == domain.DifficultyHard {
		difficultyText += "1: Easy  2: Normal  [3: Hard*]"
	} else {
		difficultyText += "1: Easy  [2: Normal*]  3: Hard"
	}

	instructions := "Press SPACE to start  |  Press Q to quit"

	padding := (m.Height - titlePadding) / 2
	if padding < 0 {
		padding = 0
	}
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}

	// Get theme styles
	titleStyle, scoreStyle, _, _ := getStyles(m)

	// Display ASCII art title
	flappyStyled := titleStyle.Render(flappyArt)
	birdStyled := titleStyle.Render(birdArt)

	b.WriteString(centerText(flappyStyled, m.Width))
	b.WriteString("\n")
	b.WriteString(centerText(birdStyled, m.Width))
	b.WriteString("\n")
	b.WriteString(centerText(subtitle, m.Width))
	b.WriteString("\n\n")
	b.WriteString(centerText(difficultyText, m.Width))
	b.WriteString("\n")

	// Display theme selection
	themeText := fmt.Sprintf("Theme: %s (Press T to change)", m.Theme.String())
	b.WriteString(centerText(themeText, m.Width))
	b.WriteString("\n")

	b.WriteString(centerText(instructions, m.Width))

	// Display high score if it exists
	if m.HighScore.Score > 0 {
		b.WriteString("\n\n")
		minutes, seconds, milliseconds := formatDuration(m.HighScore.Duration)
		highScoreText := fmt.Sprintf("High Score: %d  |  Time: %02d:%02d.%03d",
			m.HighScore.Score, minutes, seconds, milliseconds)
		b.WriteString(centerText(scoreStyle.Render(highScoreText), m.Width))
	}

	return b.String()
}

func renderGame(m game.Model) string {
	// Get theme styles
	_, scoreStyle, _, _ := getStyles(m)

	// Create empty canvas
	canvas := make([][]rune, m.Height)
	for i := range canvas {
		canvas[i] = make([]rune, m.Width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	// Draw pipes (▢▢▢▢ with ■■■■ edge)
	for _, pipe := range m.Pipes {
		for x := pipe.X; x < pipe.X+domain.PipeWidth && x < m.Width; x++ {
			if x < 0 {
				continue
			}

			// Top pipe - body is ▢, bottom edge is ■
			for y := 0; y < pipe.GapY; y++ {
				char := pipeBodyChar
				if y == pipe.GapY-1 {
					// Bottom edge of top pipe
					char = pipeEdgeChar
				}
				canvas[y][x] = []rune(char)[0]
			}

			// Bottom pipe - body is ▢, top edge is ■
			for y := pipe.GapY + pipe.GapSize; y < m.Height; y++ {
				char := pipeBodyChar
				if y == pipe.GapY+pipe.GapSize {
					// Top edge of bottom pipe
					char = pipeEdgeChar
				}
				canvas[y][x] = []rune(char)[0]
			}
		}
	}

	// Draw bird (2 characters: mO or wO)
	birdY := m.Bird.GetY()
	birdSprite := m.Bird.GetSprite()
	if birdY >= 0 && birdY < m.Height && m.Bird.X >= 0 && m.Bird.X+1 < m.Width {
		spriteRunes := []rune(birdSprite)
		canvas[birdY][m.Bird.X] = spriteRunes[0]   // First character (m or w)
		canvas[birdY][m.Bird.X+1] = spriteRunes[1] // Second character (O)
	}

	// Convert canvas to string
	var b strings.Builder
	for _, row := range canvas {
		b.WriteString(string(row))
		b.WriteString("\n")
	}

	// Add score and elapsed time
	elapsed := time.Since(m.StartTime)
	minutes, seconds, milliseconds := formatDuration(elapsed)

	score := scoreStyle.Render(fmt.Sprintf("Score: %d", m.Score))
	timeDisplay := scoreStyle.Render(fmt.Sprintf("Time: %02d:%02d.%03d", minutes, seconds, milliseconds))

	b.WriteString(score)
	b.WriteString("  ")
	b.WriteString(timeDisplay)

	return b.String()
}

func renderGameOver(m game.Model) string {
	// Get theme styles
	_, _, gameOverStyle, newRecordStyle := getStyles(m)

	var b strings.Builder

	// ASCII Art for GAME OVER
	asciiArt := ` ██████╗  █████╗ ███╗   ███╗███████╗     ██████╗ ██╗   ██╗███████╗██████╗
██╔════╝ ██╔══██╗████╗ ████║██╔════╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗
██║  ███╗███████║██╔████╔██║█████╗      ██║   ██║██║   ██║█████╗  ██████╔╝
██║   ██║██╔══██║██║╚██╔╝██║██╔══╝      ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗
╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗    ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝`

	padding := (m.Height - gameOverPadding) / 2
	if padding < 0 {
		padding = 0
	}
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}

	// Display ASCII art
	gameOverArt := gameOverStyle.Render(asciiArt)
	b.WriteString(centerText(gameOverArt, m.Width))
	b.WriteString("\n\n")

	// Calculate elapsed time
	elapsed := time.Since(m.StartTime)
	minutes, seconds, milliseconds := formatDuration(elapsed)

	// Display new record message if applicable
	if m.IsNewRecord {
		newRecord := newRecordStyle.Render("★ NEW RECORD! ★")
		b.WriteString(centerText(newRecord, m.Width))
		b.WriteString("\n")
	}

	// Display score and time
	score := fmt.Sprintf("Score: %d  |  Time: %02d:%02d.%03d", m.Score, minutes, seconds, milliseconds)
	b.WriteString(centerText(score, m.Width))
	b.WriteString("\n")

	// Display statistics
	stats := fmt.Sprintf("Jumps: %d  |  Max Height: %d  |  Avg: %.1f", m.Stats.JumpCount, m.Stats.MaxHeight, m.AvgHeight())
	b.WriteString(centerText(stats, m.Width))
	b.WriteString("\n\n")

	// Display rankings (top 5 for game over screen)
	if len(m.Rankings) > 0 {
		b.WriteString(centerText("=== TOP RANKINGS ===", m.Width))
		b.WriteString("\n")
		displayCount := 5
		if len(m.Rankings) < displayCount {
			displayCount = len(m.Rankings)
		}
		for i := 0; i < displayCount; i++ {
			rank := m.Rankings[i]
			rankMin := int(rank.Duration.Minutes())
			rankSec := int(rank.Duration.Seconds()) % 60
			rankMs := int(rank.Duration.Milliseconds()) % 1000
			rankText := fmt.Sprintf("%d. %d pts  %02d:%02d.%03d  [%s]",
				i+1, rank.Score, rankMin, rankSec, rankMs, rank.Difficulty)
			b.WriteString(centerText(rankText, m.Width))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	instructions := "Press SPACE or R to restart  |  Press Q to quit"
	b.WriteString(centerText(instructions, m.Width))

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
