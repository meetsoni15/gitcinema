package git

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

// authorSymbols is the ordered set of symbols assigned to each new author.
var authorSymbols = []string{"●", "◆", "▲", "■", "★", "✦", "◉", "⬟", "◈", "⬡", "◇", "▼"}

// authorPalette is a curated set of vibrant, distinct colors for authors.
var authorPalette = []lipgloss.Color{
	"#7aa2f7", // blue
	"#9ece6a", // green
	"#f7768e", // red/pink
	"#e0af68", // yellow
	"#bb9af7", // purple
	"#2ac3de", // cyan
	"#ff9e64", // orange
	"#73daca", // teal
	"#b4f9f8", // light cyan
	"#c0caf5", // lavender
	"#f7c67f", // light orange
	"#db4b4b", // deep red
}

// Author represents a unique contributor with display attributes.
type Author struct {
	Name   string
	Email  string
	Color  lipgloss.Color
	Symbol string
}

// Badge returns the colored symbol badge for this author.
func (a *Author) Badge() string {
	return lipgloss.NewStyle().Foreground(a.Color).Bold(true).Render(a.Symbol)
}

// Tag returns the full colored "● Name" tag.
func (a *Author) Tag() string {
	return lipgloss.NewStyle().Foreground(a.Color).Bold(true).
		Render(a.Symbol + " " + a.Name)
}

// Registry tracks all unique authors encountered in the history.
type Registry struct {
	authors map[string]*Author // keyed by email
	order   []string           // insertion order (emails)
}

// NewRegistry creates a fresh author registry.
func NewRegistry() *Registry {
	return &Registry{authors: make(map[string]*Author)}
}

// Register ensures an author is tracked, assigning color+symbol on first encounter.
func (r *Registry) Register(name, email string) *Author {
	if a, ok := r.authors[email]; ok {
		return a
	}
	idx := len(r.order)
	a := &Author{
		Name:   name,
		Email:  email,
		Color:  paletteColorForEmail(email, idx),
		Symbol: symbolForIndex(idx),
	}
	r.authors[email] = a
	r.order = append(r.order, email)
	return a
}

// Get returns an author by email, or nil.
func (r *Registry) Get(email string) *Author {
	return r.authors[email]
}

// All returns all authors in registration order.
func (r *Registry) All() []*Author {
	out := make([]*Author, len(r.order))
	for i, email := range r.order {
		out[i] = r.authors[email]
	}
	return out
}

// Len returns the number of unique authors.
func (r *Registry) Len() int {
	return len(r.order)
}

// BuildRegistry walks the full commit history and registers all authors.
func BuildRegistry(commits []Commit) *Registry {
	r := NewRegistry()
	for _, c := range commits {
		r.Register(c.Author, c.Email)
	}
	return r
}

// paletteColorForEmail picks a color — first by stable index, with FNV fallback.
func paletteColorForEmail(email string, idx int) lipgloss.Color {
	if idx < len(authorPalette) {
		return authorPalette[idx]
	}
	// Deterministic fallback via hash
	h := fnv.New32a()
	h.Write([]byte(email))
	return authorPalette[h.Sum32()%uint32(len(authorPalette))]
}

func symbolForIndex(idx int) string {
	return authorSymbols[idx%len(authorSymbols)]
}
