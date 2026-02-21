package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meetsoni15/gitcinema/internal/git"
)

// renderDetail renders the right pane showing commit details.
func renderDetail(m *Model) string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("📝 Commit") + "\n")
	sb.WriteString(strings.Repeat("─", m.rightWidth-4) + "\n")

	if len(m.commits) == 0 || m.cursor >= len(m.commits) {
		sb.WriteString(HelpStyle.Render("  No commits loaded."))
		return sb.String()
	}

	c := m.commits[m.cursor]

	// ── Hash ─────────────────────────────────────────────────────────────────
	sb.WriteString(
		HashStyle.Render("  "+c.ShortHash) + "  " +
			SubtitleStyle.Render(c.Hash[:16]+"…") + "\n\n",
	)

	// ── Subject ───────────────────────────────────────────────────────────────
	sb.WriteString(
		lipgloss.NewStyle().Foreground(ColorText).Bold(true).
			Width(m.rightWidth-4).
			Render("  "+c.Subject) + "\n\n",
	)

	// ── Author ────────────────────────────────────────────────────────────────
	var authorTag string
	if a := m.registry.Get(c.Email); a != nil {
		authorTag = a.Tag()
	} else {
		authorTag = lipgloss.NewStyle().Foreground(ColorAccent).Render("● " + c.Author)
	}
	sb.WriteString("  " + authorTag + "\n")

	// ── Dates ─────────────────────────────────────────────────────────────────
	sb.WriteString(
		"  " + DateStyle.Render(c.FormattedDate()) +
			HelpStyle.Render("  ·  "+c.RelativeTime()) + "\n\n",
	)

	// ── Divider ───────────────────────────────────────────────────────────────
	sb.WriteString(
		lipgloss.NewStyle().Foreground(ColorDim).Render("  "+strings.Repeat("─", m.rightWidth-6)) + "\n\n",
	)

	// ── Diff stats ────────────────────────────────────────────────────────────
	if m.currentDiff != nil {
		d := m.currentDiff
		sb.WriteString(
			"  " +
				StatAddStyle.Render(fmt.Sprintf("+%d", d.Additions)) + "  " +
				StatDelStyle.Render(fmt.Sprintf("-%d", d.Deletions)) + "  " +
				HelpStyle.Render(fmt.Sprintf("%d file(s) changed", d.Files)) + "\n\n",
		)

		// ── File change list ──────────────────────────────────────────────────
		sb.WriteString(
			lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("  Files Changed") + "\n",
		)

		limit := 20
		for i, fc := range d.Changes {
			if i >= limit {
				sb.WriteString(HelpStyle.Render(fmt.Sprintf("  … and %d more", len(d.Changes)-limit)) + "\n")
				break
			}
			prefix := styledChangePrefix(fc.Status)
			name := truncate(fc.Path, m.rightWidth-18)
			stat := ""
			if fc.Additions > 0 || fc.Deletions > 0 {
				stat = " " +
					lipgloss.NewStyle().Foreground(ColorAdded).Render(fmt.Sprintf("+%d", fc.Additions)) + " " +
					lipgloss.NewStyle().Foreground(ColorDeleted).Render(fmt.Sprintf("-%d", fc.Deletions))
			}
			sb.WriteString(fmt.Sprintf("  %s %s%s\n", prefix, name, stat))
		}
	} else if m.loadingDiff {
		sb.WriteString(HelpStyle.Render("  Loading diff…") + "\n")
	}

	// ── Body ──────────────────────────────────────────────────────────────────
	if c.Body != "" {
		sb.WriteString("\n" + lipgloss.NewStyle().Foreground(ColorSubtle).
			Width(m.rightWidth-4).Render("  "+c.Body) + "\n")
	}

	return sb.String()
}

// renderSearchResults renders commit list filtered by search.
func renderSearchResults(m *Model) string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("🔍 Search Results") + "\n")
	sb.WriteString(strings.Repeat("─", m.width-4) + "\n")

	if len(m.searchResults) == 0 {
		sb.WriteString(HelpStyle.Render("  No commits match \"" + m.searchQuery + "\""))
		return sb.String()
	}

	visH := m.height - 10
	for i, idx := range m.searchResults {
		if i >= visH {
			break
		}
		c := m.commits[idx]
		var authorTag string
		if reg := m.registry; reg != nil {
			if a := reg.Get(c.Email); a != nil {
				authorTag = a.Badge() + " "
			}
		}
		line := fmt.Sprintf("  %s %s  %s  %s",
			HashStyle.Render(c.ShortHash),
			authorTag,
			truncate(c.Subject, m.width-40),
			DateStyle.Render(c.FormattedDate()),
		)
		if idx == m.cursor {
			line = SelectedStyle.Width(m.width - 6).Render(line)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

// miniStatBar renders a small insertions/deletions bar for the timeline.
func miniStatBar(adds, dels int) string {
	total := adds + dels
	if total == 0 {
		return ""
	}
	width := 10
	addW := int(float64(adds) / float64(total) * float64(width))
	delW := width - addW
	return lipgloss.NewStyle().Foreground(ColorAdded).Render(strings.Repeat("█", addW)) +
		lipgloss.NewStyle().Foreground(ColorDeleted).Render(strings.Repeat("█", delW))
}

// Ensure git import is used
var _ = git.StatusAdded
