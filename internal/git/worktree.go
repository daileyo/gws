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

// IsLinkedWorktree reports whether dir is a linked git worktree rather than a
// repository clone, and returns the path of the worktree's main repository.
//
// Detection is structural and never depends on directory naming: a clone's
// .git is a directory, while a linked worktree's .git is a file containing
// "gitdir: <main>/.git/worktrees/<name>". A .git file may also belong to a
// submodule or to a clone created with --separate-git-dir; neither is a
// linked worktree, and both are reported as false.
//
// The fast path reads the .git file only. When that file is missing, malformed,
// or its target does not have the expected shape, detection falls back to
// comparing "git rev-parse --git-dir" with "git rev-parse --git-common-dir",
// which differ only inside a linked worktree.
func IsLinkedWorktree(dir string) (bool, string, error) {
	gitPath := filepath.Join(filepath.Clean(dir), ".git")

	info, err := os.Lstat(gitPath)
	if err != nil {
		return false, "", err
	}
	if info.IsDir() {
		// A .git directory is always a repository clone.
		return false, "", nil
	}

	// Fast path: read the .git file and interpret its gitdir target.
	// #nosec G304 -- gitPath is a cleaned path built from a directory the
	// scanner is already traversing; only the fixed ".git" name is appended.
	if data, readErr := os.ReadFile(gitPath); readErr == nil {
		if target := parseGitdirPointer(string(data)); target != "" {
			if !filepath.IsAbs(target) {
				target = filepath.Join(filepath.Clean(dir), target)
			}
			if mainRepo, ok := mainRepoFromWorktreeGitDir(target); ok {
				return true, mainRepo, nil
			}
		}
	}

	// Ambiguous: unreadable, malformed, or an unexpected target shape.
	return linkedWorktreeViaGit(dir)
}

// parseGitdirPointer extracts the path from a "gitdir: <path>" .git file.
// Returns an empty string when the content does not have that form.
func parseGitdirPointer(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if after, found := strings.CutPrefix(line, "gitdir:"); found {
			return strings.TrimSpace(after)
		}
	}
	return ""
}

// mainRepoFromWorktreeGitDir derives the main repository path from a linked
// worktree's git directory, which git writes as <main>/.git/worktrees/<name>.
// Returns false when the path does not have that shape.
func mainRepoFromWorktreeGitDir(gitDir string) (string, bool) {
	worktreesDir := filepath.Dir(filepath.Clean(gitDir))
	if filepath.Base(worktreesDir) != "worktrees" {
		return "", false
	}
	commonDir := filepath.Dir(worktreesDir)
	if filepath.Base(commonDir) != ".git" {
		return "", false
	}
	return filepath.Dir(commonDir), true
}

// linkedWorktreeViaGit asks git directly whether dir is a linked worktree.
// Inside a linked worktree the git directory and the common git directory
// differ; everywhere else they are the same.
func linkedWorktreeViaGit(dir string) (bool, string, error) {
	gitDir, err := gitCommand(dir, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return false, "", err
	}

	commonDir, err := gitCommand(dir, "rev-parse", "--git-common-dir")
	if err != nil {
		return false, "", err
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(filepath.Clean(dir), commonDir)
	}

	if filepath.Clean(gitDir) == filepath.Clean(commonDir) {
		return false, "", nil
	}

	// commonDir points at <main>/.git for a linked worktree.
	return true, filepath.Dir(filepath.Clean(commonDir)), nil
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
		// Repair git's worktree tracking to match the restored filesystem state
		_, _ = gitCommand(repoPath, "worktree", "repair")
		return fmt.Errorf("%w (rolled back to original location)", err)
	}

	if !srcGone && dstExists {
		// Destination already exists (likely from a prior failed attempt).
		// Clean up the stale destination so a retry can succeed.
		if rmErr := os.RemoveAll(newPath); rmErr != nil {
			return fmt.Errorf("%w (destination %s already exists and cleanup failed: %v)", err, newPath, rmErr)
		}
		// Prune stale worktree entries left by prior failed attempts
		_, _ = gitCommand(repoPath, "worktree", "prune")
		// Retry the move after cleanup
		_, retryErr := gitCommand(repoPath, "worktree", "move", currentPath, newPath)
		if retryErr != nil {
			return fmt.Errorf("%w (retry after cleanup also failed: %v)", err, retryErr)
		}
		return nil
	}

	// For any other failure, attempt repair to keep git state consistent
	_, _ = gitCommand(repoPath, "worktree", "repair")
	return err
}

// RepairWorktrees runs "git worktree repair" to fix broken internal links
// without removing entries. This is safer than prune — it preserves worktrees
// whose directories still exist but have broken git pointers.
func RepairWorktrees(repoPath string) error {
	_, err := gitCommand(repoPath, "worktree", "repair")
	return err
}

// PruneWorktrees removes worktree entries whose directories no longer exist on disk.
// Call RepairWorktrees first to salvage fixable entries before pruning truly dead ones.
func PruneWorktrees(repoPath string) error {
	_, err := gitCommand(repoPath, "worktree", "prune")
	return err
}

// IsWorktreeLocked checks whether a worktree is locked by looking for a "locked"
// file in git's internal worktree directory.
func IsWorktreeLocked(repoPath, worktreePath string) (bool, string) {
	// Git stores worktree state in .git/worktrees/<name>/
	// A "locked" file in that directory means the worktree is locked.
	// The file contents (if any) are the lock reason.
	gitDir := filepath.Join(repoPath, ".git", "worktrees")
	entries, err := os.ReadDir(gitDir)
	if err != nil {
		return false, ""
	}

	// Find the worktree entry matching this path
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		gitdirFile := filepath.Join(gitDir, entry.Name(), "gitdir")
		content, err := os.ReadFile(gitdirFile)
		if err != nil {
			continue
		}
		// gitdir file contains the path to the worktree's .git file
		linkedPath := strings.TrimSpace(string(content))
		// The gitdir points to <worktreePath>/.git
		wtPath := filepath.Dir(linkedPath)
		if filepath.Clean(wtPath) == filepath.Clean(worktreePath) {
			lockFile := filepath.Join(gitDir, entry.Name(), "locked")
			if reason, err := os.ReadFile(lockFile); err == nil {
				return true, strings.TrimSpace(string(reason))
			}
			return false, ""
		}
	}
	return false, ""
}

// IsAligned checks whether a worktree path is inside the <repoPath>.wt/ directory.
func IsAligned(worktreePath, repoPath string) bool {
	wtDir := filepath.Clean(repoPath) + ".wt"
	cleanWt := filepath.Clean(worktreePath)
	return cleanWt == wtDir || strings.HasPrefix(cleanWt, wtDir+string(filepath.Separator))
}
