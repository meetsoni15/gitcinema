package ui

import "github.com/charmbracelet/lipgloss"

// ── Base Colors (Tokyo Night) ───────────────────────────────────────────────

var (
	ColorBg       = lipgloss.AdaptiveColor{Light: "#f5f5f5", Dark: "#1a1b26"}
	ColorSurface  = lipgloss.AdaptiveColor{Light: "#e8e8e8", Dark: "#24283b"}
	ColorBorder   = lipgloss.AdaptiveColor{Light: "#9898a6", Dark: "#414868"}
	ColorText     = lipgloss.AdaptiveColor{Light: "#1a1b26", Dark: "#c0caf5"}
	ColorSubtle   = lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#565f89"}
	ColorSelected = lipgloss.Color("#7aa2f7")
	ColorAccent   = lipgloss.Color("#7aa2f7")
	ColorDim      = lipgloss.Color("#414868")
	ColorMuted    = lipgloss.Color("#565f89")
)

// ── File Change Colors ──────────────────────────────────────────────────────

var (
	ColorAdded    = lipgloss.Color("#9ece6a") // green
	ColorModified = lipgloss.Color("#e0af68") // yellow
	ColorDeleted  = lipgloss.Color("#f7768e") // red
	ColorRenamed  = lipgloss.Color("#2ac3de") // cyan
)

// ── Pane Styles ──────────────────────────────────────────────────────────────

var (
	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	ActivePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(0, 1)

	HeaderBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e2030")).
			Foreground(ColorText).
			Padding(0, 2)

	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e2030")).
			Foreground(ColorMuted).
			Padding(0, 1)

	TimelineBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1e2030")).
				Foreground(ColorText).
				Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Italic(true)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e2030")).
			Foreground(ColorSelected).
			Bold(true)

	StatAddStyle = lipgloss.NewStyle().
			Foreground(ColorAdded).Bold(true)

	StatDelStyle = lipgloss.NewStyle().
			Foreground(ColorDeleted).Bold(true)

	HashStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bb9af7")).
			Bold(true)

	DateStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle)
)

// ── File Change Status Styles ─────────────────────────────────────────────

func ChangeStyle(prefix string) lipgloss.Style {
	switch prefix {
	case "+":
		return lipgloss.NewStyle().Foreground(ColorAdded).Bold(true)
	case "~":
		return lipgloss.NewStyle().Foreground(ColorModified).Bold(true)
	case "-":
		return lipgloss.NewStyle().Foreground(ColorDeleted).Bold(true)
	case "→":
		return lipgloss.NewStyle().Foreground(ColorRenamed).Bold(true)
	}
	return lipgloss.NewStyle()
}
