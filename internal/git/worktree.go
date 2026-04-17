package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorktreeEntry represents a single git worktree discovered via git worktree list.
type WorktreeEntry struct {
	Path   string // Absolute filesystem path to the worktree
	Branch string // Branch checked out (without refs/heads/ prefix)
}

// ListWorktrees runs "git worktree list --porcelain" for the given repo and
// returns all worktrees except the main one (whose path matches repoPath).
func ListWorktrees(repoPath string) ([]WorktreeEntry, error) {
	output, err := gitCommand(repoPath, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}

	// Resolve symlinks so we correctly match the main worktree path
	// (e.g., macOS /var → /private/var)
	resolved, err := filepath.EvalSymlinks(repoPath)
	if err != nil {
		resolved = repoPath
	}

	return parseWorktreeListPorcelain(output, resolved), nil
}

// parseWorktreeListPorcelain parses the porcelain output of git worktree list.
// Each entry is separated by a blank line. Lines of interest:
//
//	worktree <path>
//	branch refs/heads/<name>
func parseWorktreeListPorcelain(output, repoPath string) []WorktreeEntry {
	if output == "" {
		return nil
	}

	repoClean := filepath.Clean(repoPath)
	var entries []WorktreeEntry

	// Split into blocks separated by blank lines
	blocks := splitWorktreeBlocks(output)
	for _, block := range blocks {
		var path, branch string
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "worktree ") {
				path = strings.TrimPrefix(line, "worktree ")
			} else if strings.HasPrefix(line, "branch ") {
				ref := strings.TrimPrefix(line, "branch ")
				branch = strings.TrimPrefix(ref, "refs/heads/")
			}
		}

		// Skip the main worktree (same path as the repo itself)
		if path == "" || filepath.Clean(path) == repoClean {
			continue
		}

		entries = append(entries, WorktreeEntry{
			Path:   path,
			Branch: branch,
		})
	}

	return entries
}

// splitWorktreeBlocks splits porcelain output into blocks separated by blank lines.
func splitWorktreeBlocks(output string) []string {
	var blocks []string
	var current strings.Builder

	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if current.Len() > 0 {
				blocks = append(blocks, current.String())
				current.Reset()
			}
			continue
		}
		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(line)
	}
	if current.Len() > 0 {
		blocks = append(blocks, current.String())
	}

	return blocks
}

// AddWorktree creates a new git worktree at destPath for the given branch.
// If the branch already exists, it is checked out. If not, a new branch is created.
func AddWorktree(repoPath, branch, destPath string) error {
	// Try checking out existing branch first
	_, err := gitCommand(repoPath, "rev-parse", "--verify", branch)
	if err == nil {
		// Branch exists — check it out into the worktree
		_, err = gitCommand(repoPath, "worktree", "add", destPath, branch)
		return err
	}
	// Branch doesn't exist — create it
	_, err = gitCommand(repoPath, "worktree", "add", "-b", branch, destPath)
	return err
}

// MoveWorktree moves a worktree from currentPath to newPath using git worktree move.
// If the move fails, it attempts to detect and recover from partial moves.
func MoveWorktree(repoPath, currentPath, newPath string) error {
	_, err := gitCommand(repoPath, "worktree", "move", currentPath, newPath)
	if err == nil {
		return nil
	}

	// Check for partial move: destination exists but source is gone
	_, srcErr := os.Stat(currentPath)
	_, dstErr := os.Stat(newPath)
	srcGone := os.IsNotExist(srcErr)
	dstExists := dstErr == nil

	if srcGone && dstExists {
		// Directory was moved but git internals weren't updated.
		// Move the directory back so the worktree isn't left in a broken state.
		if mvErr := os.Rename(newPath, currentPath); mvErr != nil {
			return fmt.Errorf("%w (rollback also failed: %v)", err, mvErr)
		}
		return fmt.Errorf("%w (rolled back to original location)", err)
	}

	if !srcGone && dstExists {
		// Destination already exists (likely from a prior failed attempt).
		// Clean up the stale destination so a retry can succeed.
		if rmErr := os.RemoveAll(newPath); rmErr != nil {
			return fmt.Errorf("%w (destination %s already exists and cleanup failed: %v)", err, newPath, rmErr)
		}
		// Retry the move after cleanup
		_, retryErr := gitCommand(repoPath, "worktree", "move", currentPath, newPath)
		if retryErr != nil {
			return fmt.Errorf("%w (retry after cleanup also failed: %v)", err, retryErr)
		}
		return nil
	}

	return err
}

// IsAligned checks whether a worktree path is inside the <repoPath>.wt/ directory.
func IsAligned(worktreePath, repoPath string) bool {
	wtDir := filepath.Clean(repoPath) + ".wt"
	cleanWt := filepath.Clean(worktreePath)
	return cleanWt == wtDir || strings.HasPrefix(cleanWt, wtDir+string(filepath.Separator))
}
