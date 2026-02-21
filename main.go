package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meetsoni15/gitcinema/internal/git"
	"github.com/meetsoni15/gitcinema/internal/ui"
)

const version = "0.1.0"

const banner = `
   ___ _ _            _
  / _ (_) |_ ___(_)_ __   ___ _ __ ___   __ _
 / /_\/ | __/ __| | '_ \ / _ \ '_ ` + "`" + ` _ \ / _` + "`" + ` |
/ /_\\| | || (__| | | | |  __/ | | | | | (_| |
\____/_|\__\___|_|_| |_|\___|_| |_| |_|\__,_|

  Watch your codebase evolve — git history as a movie.
  Version ` + version

func main() {
	args := os.Args[1:]

	// ── Flags ────────────────────────────────────────────────────────────────
	var (
		root     = "."
		branch   = ""
		maxCount = 500
		author   = ""
	)

	positionals := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			printHelp()
			return
		case "--version", "-v":
			fmt.Println("gitcinema v" + version)
			return
		case "--branch", "-b":
			if i+1 < len(args) {
				i++
				branch = args[i]
			}
		case "--max":
			if i+1 < len(args) {
				i++
				maxCount, _ = strconv.Atoi(args[i])
			}
		case "--author":
			if i+1 < len(args) {
				i++
				author = args[i]
			}
		default:
			if len(args[i]) > 0 && args[i][0] != '-' {
				positionals = append(positionals, args[i])
			}
		}
	}

	if len(positionals) > 0 {
		root = positionals[0]
	}

	// ── Validate ─────────────────────────────────────────────────────────────
	absRoot, err := filepath.Abs(root)
	if err != nil || !isDir(absRoot) {
		fmt.Fprintf(os.Stderr, "Error: %q is not a valid directory\n", root)
		os.Exit(1)
	}

	if !git.IsGitRepo(absRoot) {
		fmt.Fprintf(os.Stderr, "Error: %q is not inside a git repository\n", absRoot)
		os.Exit(1)
	}

	if branch == "" {
		branch = git.DefaultBranch(absRoot)
	}

	// ── Launch ────────────────────────────────────────────────────────────────
	m := ui.New(absRoot, branch, maxCount)
	if author != "" {
		// Pre-set author filter (passed as CLI flag)
		_ = author // model init will handle it in a future enhancement
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func printHelp() {
	fmt.Println(banner)
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  gitcinema [flags] [directory]")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  directory    Path to git repository (default: current directory)")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  -b, --branch string   Branch to walk (default: current branch)")
	fmt.Println("  --max int             Max commits to load (default: 500)")
	fmt.Println("  --author string       Pre-filter by author name")
	fmt.Println("  -v, --version         Show version")
	fmt.Println("  -h, --help            Show this help")
	fmt.Println()
	fmt.Println("KEYBINDINGS:")
	fmt.Println("  Space        Play / Pause")
	fmt.Println("  j / ↓        Next commit")
	fmt.Println("  k / ↑        Previous commit")
	fmt.Println("  g            First commit")
	fmt.Println("  G            Last commit")
	fmt.Println("  + / -        Speed up / slow down")
	fmt.Println("  /            Search commit messages")
	fmt.Println("  f            Filter by author")
	fmt.Println("  Tab          Switch pane focus")
	fmt.Println("  Esc          Clear filter / search")
	fmt.Println("  q / Ctrl+C   Quit")
}
