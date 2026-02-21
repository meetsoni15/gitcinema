package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meetsoni15/gitcinema/internal/git"
)

// renderFileTree renders the left pane showing files changed in the current commit.
func renderFileTree(m *Model) string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("📂 Changed Files") + "\n")
	sb.WriteString(strings.Repeat("─", m.leftWidth-2) + "\n")

	if m.currentDiff == nil {
		if m.loadingDiff {
			sb.WriteString(HelpStyle.Render("  loading…"))
		} else {
			sb.WriteString(HelpStyle.Render("  select a commit"))
		}
		return sb.String()
	}

	diff := m.currentDiff
	if len(diff.Changes) == 0 {
		sb.WriteString(HelpStyle.Render("  no file changes"))
		return sb.String()
	}

	// Visible height
	visH := m.height - 14
	if visH < 1 {
		visH = 1
	}

	// Scroll to keep fileScroll in view
	start := m.fileScroll
	end := start + visH
	if end > len(diff.Changes) {
		end = len(diff.Changes)
	}

	for _, fc := range diff.Changes[start:end] {
		prefix := styledChangePrefix(fc.Status)
		name := truncate(fc.Path, m.leftWidth-12)

		if fc.Status == git.StatusRenamed && fc.OldPath != "" {
			name = truncate(fc.OldPath, m.leftWidth/2-4) + " → " + truncate(fc.Path, m.leftWidth/2-4)
		}

		statStr := ""
		if fc.Additions > 0 || fc.Deletions > 0 {
			statStr = " " + lipgloss.NewStyle().Foreground(ColorAdded).Render(fmt.Sprintf("+%d", fc.Additions)) +
				lipgloss.NewStyle().Foreground(ColorDeleted).Render(fmt.Sprintf("-%d", fc.Deletions))
		}

		line := fmt.Sprintf(" %s %s%s", prefix, name, statStr)
		sb.WriteString(line + "\n")
	}

	// Summary footer
	sb.WriteString("\n")
	sb.WriteString(SubtitleStyle.Render(fmt.Sprintf(
		"  %d file(s)  ", diff.Files,
	)))
	sb.WriteString(renderFileStats(diff.Additions, diff.Deletions))

	return sb.String()
}
