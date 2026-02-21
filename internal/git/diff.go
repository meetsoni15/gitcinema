package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ChangeStatus describes how a file changed in a commit.
type ChangeStatus int

const (
	StatusAdded ChangeStatus = iota
	StatusModified
	StatusDeleted
	StatusRenamed
)

func (s ChangeStatus) String() string {
	switch s {
	case StatusAdded:
		return "added"
	case StatusModified:
		return "modified"
	case StatusDeleted:
		return "deleted"
	case StatusRenamed:
		return "renamed"
	}
	return "unknown"
}

// Prefix returns the display prefix character for a file change.
func (s ChangeStatus) Prefix() string {
	switch s {
	case StatusAdded:
		return "+"
	case StatusModified:
		return "~"
	case StatusDeleted:
		return "-"
	case StatusRenamed:
		return "→"
	}
	return "?"
}

// FileChange holds change details for one file in a commit.
type FileChange struct {
	Path      string
	OldPath   string // non-empty for renames
	Status    ChangeStatus
	Additions int
	Deletions int
}

// CommitStats holds aggregate stats for a commit.
type CommitStats struct {
	Files     int
	Additions int
	Deletions int
	Changes   []FileChange
}

// cache avoids re-running git for the same commit.
var diffCache = map[string]*CommitStats{}

// LoadDiff returns the file changes for a given commit hash.
// Results are cached in memory.
func LoadDiff(dir, hash string) (*CommitStats, error) {
	if cached, ok := diffCache[hash]; ok {
		return cached, nil
	}

	stats, err := loadDiff(dir, hash)
	if err != nil {
		return nil, err
	}
	diffCache[hash] = stats
	return stats, nil
}

func loadDiff(dir, hash string) (*CommitStats, error) {
	// Get numstat for additions/deletions per file
	numstatOut, err := exec.Command(
		"git", "-C", dir, "show", "--numstat", "--format=", hash,
	).Output()
	if err != nil {
		return nil, fmt.Errorf("git show --numstat: %w", err)
	}

	// Get name-status for change type (A/M/D/R)
	namestatOut, err := exec.Command(
		"git", "-C", dir, "show", "--name-status", "--format=", hash,
	).Output()
	if err != nil {
		return nil, fmt.Errorf("git show --name-status: %w", err)
	}

	// Parse name-status to get file statuses
	statusMap := map[string]ChangeStatus{}
	renameMap := map[string]string{}

	for _, line := range strings.Split(strings.TrimSpace(string(namestatOut)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		code := parts[0]
		switch {
		case strings.HasPrefix(code, "A"):
			statusMap[parts[1]] = StatusAdded
		case strings.HasPrefix(code, "M"):
			statusMap[parts[1]] = StatusModified
		case strings.HasPrefix(code, "D"):
			statusMap[parts[1]] = StatusDeleted
		case strings.HasPrefix(code, "R") && len(parts) >= 3:
			statusMap[parts[2]] = StatusRenamed
			renameMap[parts[2]] = parts[1]
		}
	}

	// Parse numstat
	var changes []FileChange
	totalAdd, totalDel := 0, 0

	for _, line := range strings.Split(strings.TrimSpace(string(numstatOut)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		add, _ := strconv.Atoi(parts[0])
		del, _ := strconv.Atoi(parts[1])
		path := parts[2]

		// Handle rename "old => new" format in numstat
		if strings.Contains(path, " => ") {
			// already handled via name-status
		}

		status := StatusModified
		if s, ok := statusMap[path]; ok {
			status = s
		}

		totalAdd += add
		totalDel += del

		fc := FileChange{
			Path:      path,
			Status:    status,
			Additions: add,
			Deletions: del,
		}
		if old, ok := renameMap[path]; ok {
			fc.OldPath = old
		}
		changes = append(changes, fc)
	}

	// Sort: Added → Modified → Renamed → Deleted
	sortChanges(changes)

	return &CommitStats{
		Files:     len(changes),
		Additions: totalAdd,
		Deletions: totalDel,
		Changes:   changes,
	}, nil
}

// sortChanges orders file changes: Added, Modified, Renamed, Deleted.
func sortChanges(changes []FileChange) {
	order := func(s ChangeStatus) int {
		switch s {
		case StatusAdded:
			return 0
		case StatusModified:
			return 1
		case StatusRenamed:
			return 2
		case StatusDeleted:
			return 3
		}
		return 4
	}
	for i := 0; i < len(changes)-1; i++ {
		for j := i + 1; j < len(changes); j++ {
			if order(changes[j].Status) < order(changes[i].Status) {
				changes[i], changes[j] = changes[j], changes[i]
			}
		}
	}
}
