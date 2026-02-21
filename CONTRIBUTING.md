# Contributing to gitcinema

Thank you for your interest in contributing! Here's how to get started.

## Development Setup

```bash
git clone https://github.com/meetsoni15/gitcinema
cd gitcinema
go mod tidy
go build .
go run . .
```

## Running the App

```bash
go run . [directory]          # scan a repo
go run . --branch feat/xyz .  # specific branch
go run . --max 100 .          # limit commits
```

## Code Structure

```
internal/git/
  log.go      — commit history parsing
  diff.go     — per-commit file change stats
  authors.go  — author color/symbol registry

internal/ui/
  model.go    — root Bubble Tea model + state machine
  styles.go   — Lipgloss color system
  timeline.go — timeline scrubber + legend + status bars
  filetree.go — left pane: file changes
  detail.go   — right pane: commit detail
```

## Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(ui): add branch selector
fix(git): handle empty repos gracefully
docs: update keybindings in README
```

## Pull Requests

1. Fork and create a branch: `git checkout -b feat/my-feature`
2. Make your changes and verify the build: `go build ./...`
3. Commit with a descriptive message
4. Open a PR with a clear description of what changed and why

## Reporting Issues

Please include:
- Your OS and terminal emulator
- Go version (`go version`)
- Steps to reproduce
- Expected vs actual behavior
