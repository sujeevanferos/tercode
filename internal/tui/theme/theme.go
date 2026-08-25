// Package theme defines colors, typography, borders, and styles for the TUI engine.
package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme encapsulates styled Lip Gloss definitions for rendering.
type Theme struct {
	Name    string
	Primary lipgloss.Style
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Muted   lipgloss.Style
	Border  lipgloss.Style
	Header  lipgloss.Style
	Prompt  lipgloss.Style
}

// Default returns the default Tercode dark theme.
func Default() *Theme {
	purple := lipgloss.Color("#7C3AED")
	green := lipgloss.Color("#10B981")
	yellow := lipgloss.Color("#F59E0B")
	red := lipgloss.Color("#EF4444")
	gray := lipgloss.Color("#6B7280")
	darkGray := lipgloss.Color("#374151")

	return &Theme{
		Name:    "default",
		Primary: lipgloss.NewStyle().Foreground(purple).Bold(true),
		Success: lipgloss.NewStyle().Foreground(green),
		Warning: lipgloss.NewStyle().Foreground(yellow),
		Error:   lipgloss.NewStyle().Foreground(red).Bold(true),
		Muted:   lipgloss.NewStyle().Foreground(gray),
		Border:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(darkGray),
		Header:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(purple).Padding(0, 1).Bold(true),
		Prompt:  lipgloss.NewStyle().Foreground(purple).Bold(true),
	}
}
