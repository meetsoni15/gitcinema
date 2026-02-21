<div align="center">

```
   ___ _ _            _
  / _ (_) |_ ___(_)_ __   ___ _ __ ___   __ _
 / /_\/ | __/ __| | '_ \ / _ \ '_ ` _ \ / _` |
/ /_\\| | || (__| | | | |  __/ | | | | | (_| |
\____/_|\__\___|_|_| |_|\___|_| |_| |_|\__,_|
```

**Watch your codebase evolve — git history as a movie.**

Step through every commit like a film, with live file changes, author characters, and a timeline scrubber.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go CI](https://github.com/meetsoni15/gitcinema/actions/workflows/ci.yml/badge.svg)](https://github.com/meetsoni15/gitcinema/actions/workflows/ci.yml)
[![Go Release](https://github.com/meetsoni15/gitcinema/actions/workflows/release.yml/badge.svg)](https://github.com/meetsoni15/gitcinema/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/meetsoni15/gitcinema)](https://goreportcard.com/report/github.com/meetsoni15/gitcinema)
[![Downloads](https://img.shields.io/github/downloads/meetsoni15/gitcinema/total?color=blue)](https://github.com/meetsoni15/gitcinema/releases)
[![Built with Bubble Tea](https://img.shields.io/badge/built%20with-Bubble%20Tea-ff69b4)](https://github.com/charmbracelet/bubbletea)

![demo](demo.gif)

</div>

---

## What is gitcinema?

`gitcinema` turns any git repository's history into an interactive terminal experience. Every commit is a **frame**. Press `Space` to play and watch your codebase evolve — files appear, grow, and disappear as your team's work unfolds chronologically.

- 🎬 **Play / Pause** — auto-advance through commits at configurable speed
- 📂 **Live file tree** — see exactly which files changed, added, or deleted in each frame
- 🎭 **Author characters** — every contributor gets a unique color + symbol (`●◆▲■★`)
- 🔍 **Search** — fuzzy search through commit messages with `/`
- 🎛️ **Filter** — show only one author's commits with `f`
- 🌿 **Branch-aware** — inspect any branch with `--branch`

---

## Installation

### Using `go install`
```bash
go install github.com/meetsoni15/gitcinema@latest
```
> Requires Go 1.24+. Make sure `$GOPATH/bin` is in your `$PATH`.

### Download Binary
Grab a pre-built binary for your platform from the [Releases](https://github.com/meetsoni15/gitcinema/releases) page.

### Build from Source
```bash
git clone https://github.com/meetsoni15/gitcinema
cd gitcinema
go build -o gitcinema .
```

---

## Usage

```bash
# Watch the current repository
gitcinema

# Watch a specific repository
gitcinema ./path/to/repo

# Watch a specific branch
gitcinema --branch develop .

# Limit commit history (useful for very large repos)
gitcinema --max 200 .

# Show help
gitcinema --help
```

---

## Features

### 🎬 Playback Controls
- `Space` to **play/pause** — commits advance automatically like a movie
- Adjustable speed: `0.25x → 0.5x → 1x → 2x → 4x` via `+` / `-`
- Auto-stops at the last commit

### 📂 File Tree Pane (Left)
Each commit shows which files changed, with colored prefixes:

| Symbol | Color | Meaning |
|---|---|---|
| `+` | 🟢 Green | File added |
| `~` | 🟡 Yellow | File modified |
| `-` | 🔴 Red | File deleted |
| `→` | 🔵 Cyan | File renamed |

Per-file `+N -N` line counts shown inline.

### 📝 Commit Detail Pane (Right)
- Full + short commit hash
- Commit subject in bold
- Author with their unique color badge: `● meet soni`
- Absolute date + relative time (`2 hours ago`)
- Total insertions / deletions
- Full list of changed files with per-file stats

### 🎭 Author Legend
Every unique contributor appears in the top strip with their assigned color and symbol. Colors are deterministically assigned by the order authors first appear in history — consistent across runs.

### 🔍 Search (`/`)
Type to instantly filter commits by message, author, or hash. Press `Enter` to jump to the first result.

### 🎛️ Author Filter (`f`)
Type an author name to show only their commits. Press `Esc` to clear.

---

## Keyboard Shortcuts

### Playback
| Key | Action |
|---|---|
| `Space` | Play / Pause |
| `+` / `=` | Speed up |
| `-` | Slow down |

### Navigation
| Key | Action |
|---|---|
| `j` / `↓` | Next commit |
| `k` / `↑` | Previous commit |
| `g` | Jump to first commit |
| `G` | Jump to last commit |
| `Tab` | Switch pane focus |

### Search & Filter
| Key | Action |
|---|---|
| `/` | Open commit search |
| `f` | Filter by author |
| `Esc` | Clear search / filter |

### General
| Key | Action |
|---|---|
| `q` / `Ctrl+C` | Quit |

---

## Terminal Compatibility

For the best experience, use a modern terminal with true color support:

| Terminal | Platform | Recommended |
|---|---|---|
| [Ghostty](https://ghostty.org) | macOS / Linux | ✅ Excellent |
| [Kitty](https://sw.kovidgoyal.net/kitty/) | macOS / Linux | ✅ Excellent |
| [WezTerm](https://wezfurlong.org/wezterm/) | All | ✅ Excellent |
| [iTerm2](https://iterm2.com) | macOS | ✅ Great |
| [Alacritty](https://alacritty.org) | All | ✅ Great |

---

## Built With

| Library | Purpose |
|---|---|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Styles, borders, color rendering |
| [Bubbles](https://github.com/charmbracelet/bubbles) | UI components |

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, code structure, and PR guidelines.

---

## License

MIT — see [LICENSE](LICENSE) for details.
