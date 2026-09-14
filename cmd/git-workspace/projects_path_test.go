package main

import (
	"path/filepath"
	"testing"

	"github.com/daileyo/gws/internal/xdg"
)

// projectsPath returns the expected worktree location for a repo, resolving
// symlinks so comparisons hold on systems where the temp home is a symlink.
// Tests call this instead of building "<repo>.wt/<branch>" by hand.
func projectsPath(t *testing.T, repoName string, parts ...string) string {
	t.Helper()

	dir, err := xdg.RepoProjectsDir(repoName)
	if err != nil {
		t.Fatalf("failed to resolve projects dir: %v", err)
	}
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}

	return filepath.Join(append([]string{dir}, parts...)...)
}
