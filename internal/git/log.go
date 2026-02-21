package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Commit represents a single git commit.
type Commit struct {
	Hash      string
	ShortHash string
	Author    string
	Email     string
	Timestamp time.Time
	Subject   string
	Body      string
	Index     int // position in the full history (0-based)
}

// RelativeTime returns a human-friendly relative time string.
func (c *Commit) RelativeTime() string {
	d := time.Since(c.Timestamp)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case d < 30*24*time.Hour:
		day := int(d.Hours() / 24)
		if day == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", day)
	case d < 365*24*time.Hour:
		mo := int(d.Hours() / 24 / 30)
		if mo == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", mo)
	default:
		yr := int(d.Hours() / 24 / 365)
		if yr == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", yr)
	}
}

// FormattedDate returns a nicely formatted date string.
func (c *Commit) FormattedDate() string {
	return c.Timestamp.Format("Jan 02, 2006")
}

// IsGitRepo checks whether the given directory is inside a git repository.
func IsGitRepo(dir string) bool {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// DefaultBranch returns the current branch name or HEAD.
func DefaultBranch(dir string) string {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "HEAD"
	}
	return strings.TrimSpace(string(out))
}

// ListBranches returns all local branch names.
func ListBranches(dir string) ([]string, error) {
	out, err := exec.Command("git", "-C", dir, "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, fmt.Errorf("listing branches: %w", err)
	}
	var branches []string
	for _, b := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if b = strings.TrimSpace(b); b != "" {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

// LoadHistory parses the full git commit log for the given directory and branch.
// maxCount = 0 means no limit.
func LoadHistory(dir, branch string, maxCount int) ([]Commit, error) {
	args := []string{
		"-C", dir,
		"log",
		"--reverse",
		"--format=%H|%h|%an|%ae|%at|%s",
	}
	if maxCount > 0 {
		args = append(args, fmt.Sprintf("--max-count=%d", maxCount))
	}
	if branch != "" {
		args = append(args, branch)
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("git log: %w", err)
	}

	var commits []Commit
	for i, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 6)
		if len(parts) < 6 {
			continue
		}
		ts, _ := strconv.ParseInt(parts[4], 10, 64)
		commits = append(commits, Commit{
			Hash:      parts[0],
			ShortHash: parts[1],
			Author:    parts[2],
			Email:     parts[3],
			Timestamp: time.Unix(ts, 0),
			Subject:   parts[5],
			Index:     i,
		})
	}
	return commits, nil
}

// TotalCommits returns the total number of commits on a branch without loading them all.
func TotalCommits(dir, branch string) int {
	args := []string{"-C", dir, "rev-list", "--count"}
	if branch != "" {
		args = append(args, branch)
	} else {
		args = append(args, "HEAD")
	}
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return n
}
