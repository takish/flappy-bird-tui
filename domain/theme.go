package domain

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme
type Theme int

const (
	ThemeClassic Theme = iota
	ThemeRetro
	ThemeNeon
)

// ColorScheme holds the colors for a theme
type ColorScheme struct {
	Title     lipgloss.Color
	Score     lipgloss.Color
	GameOver  lipgloss.Color
	NewRecord lipgloss.Color
}

// GetColors returns the color scheme for a theme
func (t Theme) GetColors() ColorScheme {
	switch t {
	case ThemeRetro:
		return ColorScheme{
			Title:     lipgloss.Color("10"), // Green
			Score:     lipgloss.Color("10"), // Green
			GameOver:  lipgloss.Color("9"),  // Red
			NewRecord: lipgloss.Color("11"), // Yellow
		}
	case ThemeNeon:
		return ColorScheme{
			Title:     lipgloss.Color("13"), // Magenta
			Score:     lipgloss.Color("14"), // Cyan
			GameOver:  lipgloss.Color("9"),  // Red
			NewRecord: lipgloss.Color("11"), // Yellow
		}
	default: // ThemeClassic
		return ColorScheme{
			Title:     lipgloss.Color("12"), // Blue
			Score:     lipgloss.Color("10"), // Green
			GameOver:  lipgloss.Color("9"),  // Red
			NewRecord: lipgloss.Color("11"), // Yellow
		}
	}
}

// String returns the string representation of the theme
func (t Theme) String() string {
	switch t {
	case ThemeRetro:
		return "Retro"
	case ThemeNeon:
		return "Neon"
	default:
		return "Classic"
	}
}

// Next returns the next theme in rotation
func (t Theme) Next() Theme {
	return (t + 1) % 3
}
