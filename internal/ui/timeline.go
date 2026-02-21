package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meetsoni15/gitcinema/internal/git"
)

// renderTimeline renders the bottom timeline scrubber bar.
func renderTimeline(m *Model) string {
	if len(m.commits) == 0 {
		return TimelineBarStyle.Width(m.width).Render("  no commits")
	}

	c := m.commits[m.cursor]
	total := len(m.commits)

	// ── Progress bar ──────────────────────────────────────────────────────────
	barWidth := m.width - 20
	if barWidth < 10 {
		barWidth = 10
	}
	filled := int(math.Round(float64(m.cursor) / float64(total-1) * float64(barWidth)))
	if total == 1 {
		filled = barWidth
	}

	bar := lipgloss.NewStyle().Foreground(ColorAccent).Render(strings.Repeat("━", filled)) +
		lipgloss.NewStyle().Foreground(ColorDim).Render(strings.Repeat("─", barWidth-filled))

	// ── Playback indicator ───────────────────────────────────────────────────
	playIcon := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("▶")
	if m.playing {
		playIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e")).Bold(true).Render("⏸")
	}
	speedStr := lipgloss.NewStyle().Foreground(ColorSubtle).
		Render(fmt.Sprintf("%.2gx", m.speed))

	// ── Position label ───────────────────────────────────────────────────────
	posLabel := lipgloss.NewStyle().Foreground(ColorMuted).
		Render(fmt.Sprintf(" %d/%d ", m.cursor+1, total))

	// Author badge
	authorBadge := ""
	if a := m.registry.Get(c.Email); a != nil {
		authorBadge = " " + a.Badge() + " "
	}

	// Date
	dateStr := DateStyle.Render(c.FormattedDate())

	row1 := " " + playIcon + " " + speedStr + "  " + bar
	row2 := posLabel + authorBadge +
		HashStyle.Render(c.ShortHash) + "  " +
		lipgloss.NewStyle().Foreground(ColorText).Render(truncate(c.Subject, m.width-40)) +
		"  " + dateStr

	return TimelineBarStyle.Width(m.width).Render(row1) + "\n" +
		TimelineBarStyle.Width(m.width).Render(row2)
}

// renderLegend renders the top author legend strip.
func renderLegend(m *Model) string {
	authors := m.registry.All()
	if len(authors) == 0 {
		return ""
	}

	var parts []string
	for _, a := range authors {
		parts = append(parts, a.Tag())
	}

	// Filter indicator
	filterStr := ""
	if m.filterAuthor != "" {
		filterStr = "  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")).
			Bold(true).
			Render("filter: "+m.filterAuthor)
	}

	legend := strings.Join(parts, "  ")
	right := HelpStyle.Render(fmt.Sprintf("%d authors", len(authors))) + filterStr

	// Right-align the count
	pad := m.width - lipgloss.Width(legend) - lipgloss.Width(right) - 6
	if pad < 1 {
		pad = 1
	}

	return HeaderBarStyle.Width(m.width).Render(
		"  " + legend + strings.Repeat(" ", pad) + right,
	)
}

// renderStatusBar renders the bottom keybinding help bar.
func renderStatusBar(m *Model) string {
	bindings := []string{
		KeyStyle.Render("Space") + HelpStyle.Render(" play/pause"),
		KeyStyle.Render("j/k") + HelpStyle.Render(" step"),
		KeyStyle.Render("+/-") + HelpStyle.Render(" speed"),
		KeyStyle.Render("g/G") + HelpStyle.Render(" first/last"),
		KeyStyle.Render("f") + HelpStyle.Render(" filter"),
		KeyStyle.Render("/") + HelpStyle.Render(" search"),
		KeyStyle.Render("Tab") + HelpStyle.Render(" pane"),
		KeyStyle.Render("q") + HelpStyle.Render(" quit"),
	}
	return StatusBarStyle.Width(m.width).Render("  " + strings.Join(bindings, "  "))
}

// renderSearchBar renders the search input when in search mode.
func renderSearchBar(m *Model) string {
	prompt := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("/")
	input := lipgloss.NewStyle().Foreground(ColorText).Render(m.searchQuery)
	cursor := lipgloss.NewStyle().Foreground(ColorAccent).Render("█")
	hits := ""
	if m.searchQuery != "" {
		hits = HelpStyle.Render(fmt.Sprintf("  %d results", len(m.searchResults)))
	}
	return StatusBarStyle.Width(m.width).Render("  " + prompt + " " + input + cursor + hits)
}

// renderFilterBar renders the author filter input.
func renderFilterBar(m *Model) string {
	prompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0af68")).Bold(true).Render("filter author:")
	input := lipgloss.NewStyle().Foreground(ColorText).Render(m.filterQuery)
	cursor := lipgloss.NewStyle().Foreground(lipgloss.Color("#e0af68")).Render("█")
	return StatusBarStyle.Width(m.width).Render("  " + prompt + " " + input + cursor)
}

// truncate truncates a string to maxLen runes, adding "…" if needed.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return "…"
	}
	return string(runes[:maxLen-1]) + "…"
}

// renderScanningHeader renders the top bar during/after load.
func renderHeader(m *Model) string {
	branch := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(m.branch)
	path := SubtitleStyle.Render(m.root)
	total := HelpStyle.Render(fmt.Sprintf("%d commits", len(m.commits)))

	title := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("🎬 gitcinema")

	return HeaderBarStyle.Width(m.width).Render(
		fmt.Sprintf("  %s  %s  branch: %s  %s", title, path, branch, total),
	)
}

// renderFileStats renders "+N -N" stat inline.
func renderFileStats(adds, dels int) string {
	return StatAddStyle.Render(fmt.Sprintf("+%d", adds)) + "  " +
		StatDelStyle.Render(fmt.Sprintf("-%d", dels))
}

// styledChangePrefix returns a styled prefix for a file change.
func styledChangePrefix(status git.ChangeStatus) string {
	prefix := status.Prefix()
	return ChangeStyle(prefix).Render(prefix)
}
