package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meetsoni15/gitcinema/internal/git"
)

// ── Speed presets ────────────────────────────────────────────────────────────

var speedPresets = []float64{0.25, 0.5, 1.0, 2.0, 4.0}
var defaultSpeedIdx = 2 // 1.0x

const defaultInterval = 800 * time.Millisecond

// ── Model state ──────────────────────────────────────────────────────────────

type AppState int

const (
	StateLoading   AppState = iota
	StateReady              // normal interactive mode
	StatePlaying            // auto-advancing
	StateSearching          // "/" search input active
	StateFiltering          // "f" filter author input active
)

// ActivePane tracks which pane has focus.
type ActivePane int

const (
	PaneFiles  ActivePane = iota
	PaneDetail ActivePane = iota
)

// ── Messages ─────────────────────────────────────────────────────────────────

type loadDoneMsg struct {
	commits  []git.Commit
	registry *git.Registry
	err      error
}

type diffLoadedMsg struct {
	hash  string
	stats *git.CommitStats
	err   error
}

type playTickMsg struct{}
type spinnerTickMsg struct{}

// ── Model ─────────────────────────────────────────────────────────────────────

// Model is the root Bubble Tea model for gitcinema.
type Model struct {
	// data
	root     string
	branch   string
	maxCount int
	commits  []git.Commit
	registry *git.Registry

	// navigation
	cursor     int
	fileScroll int
	activePane ActivePane

	// playback
	state    AppState
	playing  bool
	speedIdx int
	speed    float64

	// diff
	currentDiff *git.CommitStats
	loadingDiff bool

	// search
	searchQuery   string
	searchResults []int // indices into commits

	// filter
	filterQuery     string
	filterAuthor    string
	filteredCommits []git.Commit // nil = no filter

	// layout
	width      int
	height     int
	leftWidth  int
	rightWidth int

	// spinner
	spinnerFrame int

	// error
	err error
}

// New creates the initial model.
func New(root, branch string, maxCount int) Model {
	return Model{
		root:     root,
		branch:   branch,
		maxCount: maxCount,
		state:    StateLoading,
		speedIdx: defaultSpeedIdx,
		speed:    speedPresets[defaultSpeedIdx],
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadHistory(), spinnerTick())
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (m Model) loadHistory() tea.Cmd {
	return func() tea.Msg {
		commits, err := git.LoadHistory(m.root, m.branch, m.maxCount)
		if err != nil {
			return loadDoneMsg{err: err}
		}
		registry := git.BuildRegistry(commits)
		return loadDoneMsg{commits: commits, registry: registry}
	}
}

func (m Model) loadDiff(hash string) tea.Cmd {
	return func() tea.Msg {
		stats, err := git.LoadDiff(m.root, hash)
		return diffLoadedMsg{hash: hash, stats: stats, err: err}
	}
}

func playTick(speed float64) tea.Cmd {
	dur := time.Duration(float64(defaultInterval) / speed)
	return tea.Tick(dur, func(t time.Time) tea.Msg { return playTickMsg{} })
}

func spinnerTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg { return spinnerTickMsg{} })
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()

	case spinnerTickMsg:
		m.spinnerFrame = (m.spinnerFrame + 1) % len(spinnerFrames)
		if m.state == StateLoading {
			return m, spinnerTick()
		}

	case loadDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StateReady
			return m, nil
		}
		m.commits = msg.commits
		m.registry = msg.registry
		m.state = StateReady
		if len(m.commits) > 0 {
			return m, m.loadDiff(m.commits[0].Hash)
		}

	case diffLoadedMsg:
		m.loadingDiff = false
		if msg.err == nil && msg.hash == m.currentCommit().Hash {
			m.currentDiff = msg.stats
		}

	case playTickMsg:
		if m.playing {
			return m.stepForward()
		}

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ── Search mode ──────────────────────────────────────────────────────────
	if m.state == StateSearching {
		return m.handleSearchKey(msg)
	}
	// ── Filter mode ──────────────────────────────────────────────────────────
	if m.state == StateFiltering {
		return m.handleFilterKey(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "j", "down":
		m, cmd := m.stepForward()
		return m, cmd

	case "k", "up":
		m, cmd := m.stepBackward()
		return m, cmd

	case "g":
		m.cursor = 0
		m.currentDiff = nil
		m.loadingDiff = true
		return m, m.loadDiff(m.currentCommit().Hash)

	case "G":
		m.cursor = len(m.activeCommits()) - 1
		m.currentDiff = nil
		m.loadingDiff = true
		return m, m.loadDiff(m.currentCommit().Hash)

	case " ":
		m.playing = !m.playing
		if m.playing {
			m.state = StatePlaying
			return m, playTick(m.speed)
		}
		m.state = StateReady

	case "+", "=":
		if m.speedIdx < len(speedPresets)-1 {
			m.speedIdx++
			m.speed = speedPresets[m.speedIdx]
		}

	case "-":
		if m.speedIdx > 0 {
			m.speedIdx--
			m.speed = speedPresets[m.speedIdx]
		}

	case "tab":
		if m.activePane == PaneFiles {
			m.activePane = PaneDetail
		} else {
			m.activePane = PaneFiles
		}

	case "/":
		m.stopPlaying()
		m.state = StateSearching
		m.searchQuery = ""
		m.searchResults = nil

	case "f":
		m.stopPlaying()
		m.state = StateFiltering
		m.filterQuery = ""

	case "esc":
		m.stopPlaying()
		m.state = StateReady
		m.filterAuthor = ""
		m.filteredCommits = nil
		m.searchQuery = ""
		m.searchResults = nil
		m.cursor = 0
		m.currentDiff = nil
		m.loadingDiff = true
		if len(m.commits) > 0 {
			return m, m.loadDiff(m.commits[0].Hash)
		}
	}

	return m, nil
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateReady
		m.searchQuery = ""
		m.searchResults = nil
	case "enter":
		if len(m.searchResults) > 0 {
			m.cursor = m.searchResults[0]
			m.state = StateReady
			m.currentDiff = nil
			m.loadingDiff = true
			return m, m.loadDiff(m.currentCommit().Hash)
		}
		m.state = StateReady
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.runSearch()
		}
	default:
		if len(msg.Runes) > 0 {
			m.searchQuery += string(msg.Runes)
			m.runSearch()
		}
	}
	return m, nil
}

func (m *Model) runSearch() {
	m.searchResults = nil
	if m.searchQuery == "" {
		return
	}
	q := strings.ToLower(m.searchQuery)
	for i, c := range m.commits {
		if strings.Contains(strings.ToLower(c.Subject), q) ||
			strings.Contains(strings.ToLower(c.Author), q) ||
			strings.Contains(c.ShortHash, q) {
			m.searchResults = append(m.searchResults, i)
		}
	}
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateReady
		m.filterQuery = ""
	case "enter":
		m.filterAuthor = m.filterQuery
		m.filteredCommits = nil
		if m.filterAuthor != "" {
			q := strings.ToLower(m.filterAuthor)
			for _, c := range m.commits {
				if strings.Contains(strings.ToLower(c.Author), q) ||
					strings.Contains(strings.ToLower(c.Email), q) {
					m.filteredCommits = append(m.filteredCommits, c)
				}
			}
		}
		m.cursor = 0
		m.currentDiff = nil
		m.loadingDiff = true
		m.state = StateReady
		if len(m.activeCommits()) > 0 {
			return m, m.loadDiff(m.activeCommits()[0].Hash)
		}
	case "backspace":
		if len(m.filterQuery) > 0 {
			m.filterQuery = m.filterQuery[:len(m.filterQuery)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.filterQuery += string(msg.Runes)
		}
	}
	return m, nil
}

// activeCommits returns filtered commits if a filter is active, otherwise all commits.
func (m *Model) activeCommits() []git.Commit {
	if m.filteredCommits != nil {
		return m.filteredCommits
	}
	return m.commits
}

func (m *Model) currentCommit() git.Commit {
	ac := m.activeCommits()
	if len(ac) == 0 || m.cursor >= len(ac) {
		return git.Commit{}
	}
	return ac[m.cursor]
}

func (m Model) stepForward() (Model, tea.Cmd) {
	ac := m.activeCommits()
	if m.cursor < len(ac)-1 {
		m.cursor++
		m.currentDiff = nil
		m.loadingDiff = true
		var cmds []tea.Cmd
		cmds = append(cmds, m.loadDiff(m.currentCommit().Hash))
		if m.playing {
			cmds = append(cmds, playTick(m.speed))
		}
		return m, tea.Batch(cmds...)
	}
	// Reached end — stop playback
	m.playing = false
	m.state = StateReady
	return m, nil
}

func (m Model) stepBackward() (Model, tea.Cmd) {
	if m.cursor > 0 {
		m.cursor--
		m.currentDiff = nil
		m.loadingDiff = true
		return m, m.loadDiff(m.currentCommit().Hash)
	}
	return m, nil
}

func (m *Model) stopPlaying() {
	m.playing = false
	m.state = StateReady
}

func (m *Model) recalcLayout() {
	m.leftWidth = m.width * 38 / 100
	m.rightWidth = m.width - m.leftWidth - 3
}

// ── View ──────────────────────────────────────────────────────────────────────

var spinnerFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

func (m Model) View() string {
	if m.width == 0 {
		return "initializing…"
	}

	if m.state == StateLoading {
		return m.renderLoading()
	}

	if m.err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e")).Bold(true).
			Render(fmt.Sprintf("\n  Error: %v\n\n  Press q to quit.", m.err))
	}

	if len(m.commits) == 0 {
		return lipgloss.NewStyle().Foreground(ColorSubtle).
			Render("\n  No commits found in this branch.\n\n  Press q to quit.")
	}

	m.recalcLayout()

	header := renderHeader(&m)
	legend := renderLegend(&m)

	// ── Main body ─────────────────────────────────────────────────────────────
	var body string
	if m.state == StateSearching {
		body = PaneStyle.Width(m.width - 4).Height(m.height - 9).
			Render(renderSearchResults(&m))
	} else {
		leftStyle := PaneStyle
		rightStyle := PaneStyle
		if m.activePane == PaneFiles {
			leftStyle = ActivePaneStyle
		} else {
			rightStyle = ActivePaneStyle
		}
		left := leftStyle.Width(m.leftWidth).Height(m.height - 9).Render(renderFileTree(&m))
		right := rightStyle.Width(m.rightWidth).Height(m.height - 9).Render(renderDetail(&m))
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)
	}

	// ── Timeline ──────────────────────────────────────────────────────────────
	timeline := renderTimeline(&m)

	// ── Status/search/filter bar ─────────────────────────────────────────────
	var statusBar string
	switch m.state {
	case StateSearching:
		statusBar = renderSearchBar(&m)
	case StateFiltering:
		statusBar = renderFilterBar(&m)
	default:
		statusBar = renderStatusBar(&m)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, legend, body, timeline, statusBar)
}

func (m Model) renderLoading() string {
	frame := spinnerFrames[m.spinnerFrame]
	msg := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).
		Render(fmt.Sprintf("\n  %s  Loading git history for %s…\n", frame, m.root))
	hint := HelpStyle.Render("  Reading commits, building author registry…")
	return msg + "\n" + hint
}
