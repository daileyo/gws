package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// workspaceFixture describes a workspace built for reconciliation tests, so
// both init and refresh tests can assert against the same shape.
type workspaceFixture struct {
	// Root is the workspace directory itself.
	Root string
	// External is a directory outside the workspace holding repositories that
	// are reached through symlinks.
	External string

	// TopLevel repositories sit directly under the workspace root.
	TopLevel []string
	// Nested is a repository two container directories deep.
	Nested string
	// Symlinked is an external repository reached through a workspace symlink.
	Symlinked string
	// WithWorktrees is a repository owning one aligned and one unaligned worktree.
	WithWorktrees string
	// AlignedWorktree is inside <repo>.wt/; UnalignedWorktree is not.
	AlignedWorktree   string
	UnalignedWorktree string
}

// ExpectedRepoCount is the number of repositories a full scan should find.
func (f workspaceFixture) ExpectedRepoCount() int {
	return len(f.TopLevel) + 3 // nested + symlinked + withWorktrees
}

// newWorkspaceFixture builds a workspace exercising every discovery path this
// spec cares about: top-level repositories, a repository nested under
// container directories, an external repository reached by symlink, and a
// repository with both an aligned and an unaligned worktree.
//
// HOME is redirected to a temp directory so the developer's real
// ~/.gws/config.json is never read or written.
func newWorkspaceFixture(t *testing.T) workspaceFixture {
	t.Helper()

	t.Setenv("HOME", t.TempDir())

	root := resolvedTempDir(t)
	external := resolvedTempDir(t)

	f := workspaceFixture{Root: root, External: external}

	for _, name := range []string{"alpha", "bravo"} {
		f.TopLevel = append(f.TopLevel, fixtureRepo(t, filepath.Join(root, name)))
	}

	f.Nested = fixtureRepo(t, filepath.Join(root, "org-a", "team-one", "service-api"))

	extRepo := fixtureRepo(t, filepath.Join(external, "external-repo"))
	f.Symlinked = extRepo
	if err := os.Symlink(extRepo, filepath.Join(root, "external-repo")); err != nil {
		t.Fatalf("Failed to create workspace symlink: %v", err)
	}

	f.WithWorktrees = fixtureRepo(t, filepath.Join(root, "with-worktrees"))
	f.AlignedWorktree = filepath.Join(root, "with-worktrees.wt", "feature")
	f.UnalignedWorktree = filepath.Join(root, "loose-worktrees", "hotfix")
	fixtureWorktree(t, f.WithWorktrees, "feature", f.AlignedWorktree)
	fixtureWorktree(t, f.WithWorktrees, "hotfix", f.UnalignedWorktree)

	return f
}

// resolvedTempDir returns a temp directory with symlinks resolved, so paths
// compare equal to what the scanner stores.
func resolvedTempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return dir
	}
	return resolved
}

// fixtureRepo creates a git repository with one commit and returns its
// resolved path.
func fixtureRepo(t *testing.T, path string) string {
	t.Helper()

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create %s: %v", path, err)
	}

	cmds := [][]string{
		{"init"},
		{"config", "user.email", "fixture@example.com"},
		{"config", "user.name", "Fixture User"},
		{"remote", "add", "origin", "https://github.com/test/" + filepath.Base(path) + ".git"},
		{"commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		runFixtureGit(t, path, args...)
	}

	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolved
}

// fixtureWorktree adds a linked worktree on a new branch.
func fixtureWorktree(t *testing.T, repoDir, branch, destPath string) {
	t.Helper()
	runFixtureGit(t, repoDir, "worktree", "add", "-b", branch, destPath)
}

// runFixtureGit runs a git command in dir, failing the test on error.
func runFixtureGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...) // #nosec G204 -- fixed test arguments
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s failed: %s\n%s", args, dir, err, out)
	}
}

// newSymlinkOnlyWorkspace builds a curated workspace containing nothing but
// symlinks to external repositories — the layout the previous implementation
// could not discover at all.
func newSymlinkOnlyWorkspace(t *testing.T, names ...string) (root string, repoPaths []string) {
	t.Helper()

	t.Setenv("HOME", t.TempDir())

	root = resolvedTempDir(t)
	external := resolvedTempDir(t)

	for _, name := range names {
		repoPath := fixtureRepo(t, filepath.Join(external, name))
		repoPaths = append(repoPaths, repoPath)
		if err := os.Symlink(repoPath, filepath.Join(root, name)); err != nil {
			t.Fatalf("Failed to create symlink for %s: %v", name, err)
		}
	}

	return root, repoPaths
}

// workspaceListing returns a sorted recursive listing of a directory,
// recording symlinks as links rather than following them.
func workspaceListing(t *testing.T, root string) []string {
	t.Helper()

	var entries []string
	var walk func(dir, prefix string)
	walk = func(dir, prefix string) {
		items, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, item := range items {
			name := filepath.Join(prefix, item.Name())
			if item.Type()&os.ModeSymlink != 0 {
				target, _ := os.Readlink(filepath.Join(dir, item.Name()))
				entries = append(entries, "link:"+name+" -> "+target)
				continue
			}
			entries = append(entries, name)
			if item.IsDir() {
				walk(filepath.Join(dir, item.Name()), name)
			}
		}
	}
	walk(root, "")

	return entries
}

// symlinkEntries filters a listing down to symlinks.
func symlinkEntries(listing []string) []string {
	var links []string
	for _, entry := range listing {
		if len(entry) > 5 && entry[:5] == "link:" {
			links = append(links, entry)
		}
	}
	return links
}
