package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Status represents the git status of a repository
type Status struct {
	Branch      string    // Current branch name
	IsClean     bool      // True if working tree is clean
	HasChanges  bool      // True if there are uncommitted changes
	Ahead       int       // Number of commits ahead of remote
	Behind      int       // Number of commits behind remote
	LastChecked time.Time // When status was last checked
}

// GetStatus reads the current git status of a repository using the git CLI
func GetStatus(repoPath string) (*Status, error) {
	// Verify this is a git repository
	_, err := gitCommand(repoPath, "rev-parse", "--git-dir")
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	status := &Status{
		LastChecked: time.Now(),
	}

	// Get current branch name
	branch, err := gitCommand(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		// Repository might not have any commits yet
		status.Branch = ""
		status.IsClean = true
		return status, nil
	}

	if branch == "HEAD" {
		status.Branch = "HEAD detached"
	} else {
		status.Branch = branch
	}

	// Get working tree status
	porcelain, err := gitCommand(repoPath, "status", "--porcelain")
	if err != nil {
		return status, nil
	}

	status.IsClean = porcelain == ""
	status.HasChanges = porcelain != ""

	// Get ahead/behind info (only for branches with an upstream)
	if status.Branch != "HEAD detached" {
		ahead, behind := getAheadBehind(repoPath, status.Branch)
		status.Ahead = ahead
		status.Behind = behind
	}

	return status, nil
}

// getAheadBehind calculates commits ahead and behind the remote branch
func getAheadBehind(repoPath, branch string) (ahead int, behind int) {
	// Use rev-list to count commits ahead and behind in one call
	output, err := gitCommand(repoPath, "rev-list", "--count", "--left-right", branch+"...origin/"+branch)
	if err != nil {
		// No remote tracking branch or other error
		return 0, 0
	}

	parts := strings.Fields(output)
	if len(parts) != 2 {
		return 0, 0
	}

	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0
	}

	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0
	}

	return a, b
}

// gitCommand runs a git command in the given directory and returns trimmed stdout
func gitCommand(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			return "", fmt.Errorf("git %s failed: %s", strings.Join(args, " "), stderr)
		}
		return "", fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// IsStale checks if the status is older than the given duration
func (s *Status) IsStale(ttl time.Duration) bool {
	return time.Since(s.LastChecked) > ttl
}

// String returns a human-readable status string
func (s *Status) String() string {
	if s.Branch == "" {
		return "no commits"
	}

	status := s.Branch
	if s.HasChanges {
		status += " (dirty)"
	} else {
		status += " (clean)"
	}

	if s.Ahead > 0 && s.Behind > 0 {
		status += fmt.Sprintf(" ↑%d ↓%d", s.Ahead, s.Behind)
	} else if s.Ahead > 0 {
		status += fmt.Sprintf(" ↑%d", s.Ahead)
	} else if s.Behind > 0 {
		status += fmt.Sprintf(" ↓%d", s.Behind)
	}

	return status
}
