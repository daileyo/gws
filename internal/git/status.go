package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

// GetStatus reads the current git status of a repository
func GetStatus(repoPath string) (*Status, error) {
	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	status := &Status{
		LastChecked: time.Now(),
	}

	// Get current branch
	head, err := repo.Head()
	if err != nil {
		// Repository might not have any commits yet
		status.Branch = ""
		status.IsClean = true
		return status, nil
	}

	// Extract branch name
	if head.Name().IsBranch() {
		status.Branch = head.Name().Short()
	} else {
		status.Branch = "HEAD detached"
	}

	// Get working tree status
	worktree, err := repo.Worktree()
	if err != nil {
		return status, nil // Return what we have
	}

	wStatus, err := worktree.Status()
	if err != nil {
		return status, nil // Return what we have
	}

	// Check if there are uncommitted changes
	status.IsClean = wStatus.IsClean()
	status.HasChanges = !wStatus.IsClean()

	// Get ahead/behind info
	ahead, behind, err := getAheadBehind(repo, head)
	if err == nil {
		status.Ahead = ahead
		status.Behind = behind
	}

	return status, nil
}

// getAheadBehind calculates commits ahead and behind the remote branch
func getAheadBehind(repo *git.Repository, head *plumbing.Reference) (ahead int, behind int, error error) {
	// Get the remote reference for the current branch
	branchName := head.Name().Short()

	// Try to get the remote tracking branch
	remoteBranch := plumbing.NewRemoteReferenceName("origin", branchName)
	remoteRef, err := repo.Reference(remoteBranch, true)
	if err != nil {
		// No remote tracking branch
		return 0, 0, nil
	}

	// Get commit iterators
	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return 0, 0, err
	}

	remoteCommit, err := repo.CommitObject(remoteRef.Hash())
	if err != nil {
		return 0, 0, err
	}

	// Count commits ahead (local commits not in remote)
	ahead, err = countCommitsAhead(repo, headCommit.Hash, remoteCommit.Hash)
	if err != nil {
		return 0, 0, err
	}

	// Count commits behind (remote commits not in local)
	behind, err = countCommitsAhead(repo, remoteCommit.Hash, headCommit.Hash)
	if err != nil {
		return 0, 0, err
	}

	return ahead, behind, nil
}

// countCommitsAhead counts commits reachable from 'from' but not from 'to'
func countCommitsAhead(repo *git.Repository, from, to plumbing.Hash) (int, error) {
	if from == to {
		return 0, nil
	}

	// Get all commits reachable from 'from'
	fromIter, err := repo.Log(&git.LogOptions{From: from})
	if err != nil {
		return 0, err
	}
	defer fromIter.Close()

	// Get all commits reachable from 'to'
	toCommits := make(map[plumbing.Hash]bool)
	toIter, err := repo.Log(&git.LogOptions{From: to})
	if err != nil {
		return 0, err
	}
	defer toIter.Close()

	err = toIter.ForEach(func(c *object.Commit) error {
		toCommits[c.Hash] = true
		return nil
	})
	if err != nil {
		return 0, err
	}

	// Count commits in 'from' that are not in 'to'
	count := 0
	err = fromIter.ForEach(func(c *object.Commit) error {
		if !toCommits[c.Hash] {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return count, nil
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
