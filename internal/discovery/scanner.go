package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gogit "github.com/go-git/go-git/v5"

	"github.com/daileyo/gws/internal/classifier"
	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
)

// Options controls how a workspace scan traverses the filesystem.
type Options struct {
	// MaxDepth is the number of directory levels below the scan root that
	// traversal descends. A repository found at exactly MaxDepth is
	// registered; nothing below it is traversed. Zero or negative means
	// config.DefaultScanMaxDepth.
	MaxDepth int
}

// ScanResult contains the results of a repository scan
type ScanResult struct {
	Repositories []config.Repository
	Count        int
	Errors       []error

	// ReachablePaths is the set of resolved repository paths that were
	// reachable from the scan root, including those reached by following
	// symlinks. Callers use it to tell "tracked but no longer reachable"
	// apart from "tracked and still present".
	ReachablePaths map[string]bool

	// WorktreeMainRepos is the set of main repository paths inferred from
	// linked worktrees seen during the scan. A main repository may live
	// outside the scan root.
	WorktreeMainRepos map[string]bool
}

// Scan recursively scans the given directory for git repositories.
//
// Traversal is repository-boundary aware: once a directory is identified as a
// repository root it is registered and traversal stops there, so submodules
// and nested clones are not registered separately. Sibling and ancestor
// container directories continue to be scanned.
func Scan(rootPath string, opts Options) (*ScanResult, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify the root path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}

	// Resolve the root itself so symlink guards compare real paths.
	rootReal, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve directory path: %w", err)
	}

	s := &scanner{
		result: &ScanResult{
			Repositories:      []config.Repository{},
			Errors:            []error{},
			ReachablePaths:    make(map[string]bool),
			WorktreeMainRepos: make(map[string]bool),
		},
		maxDepth: effectiveMaxDepth(opts.MaxDepth),
		root:     rootReal,
		visited:  make(map[string]bool),
	}

	s.walk(rootReal, 0)

	s.result.Count = len(s.result.Repositories)
	return s.result, nil
}

// effectiveMaxDepth resolves a caller-supplied depth, falling back to the
// configured default so callers cannot accidentally disable the cap.
func effectiveMaxDepth(depth int) int {
	if depth > 0 {
		return depth
	}
	return config.DefaultScanMaxDepth
}

// scanner carries the mutable state of a single scan.
type scanner struct {
	result   *ScanResult
	maxDepth int
	// root is the resolved real path of the scan root.
	root string
	// visited holds resolved real paths already processed, so a directory
	// reachable by several paths is handled exactly once and symlink loops
	// terminate. This is what makes deduplication order-independent.
	visited map[string]bool
}

// walk visits dir at the given depth below the scan root. dir must already be
// a resolved real path. It registers dir as a repository and prunes traversal
// when dir is a repository root, otherwise it recurses into eligible
// subdirectories.
func (s *scanner) walk(dir string, depth int) {
	if s.visited[dir] {
		return
	}
	s.visited[dir] = true

	// A repository root ends traversal down this branch regardless of depth.
	switch classifyDir(dir) {
	case dirRepository:
		s.register(dir)
		return
	case dirGitLinked:
		// A .git file marks a linked worktree, a submodule, or a clone whose
		// git directory lives elsewhere. Only linked worktrees are excluded;
		// the rest are ordinary repositories. Either way, traversal stops.
		isWorktree, mainRepo, err := git.IsLinkedWorktree(dir)
		if err == nil && isWorktree {
			if mainRepo != "" {
				s.result.WorktreeMainRepos[mainRepo] = true
			}
			return
		}
		s.register(dir)
		return
	case dirPlain:
	}

	if depth >= s.maxDepth {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		s.result.Errors = append(s.result.Errors, fmt.Errorf("error accessing %s: %w", dir, err))
		return
	}

	for _, entry := range entries {
		if shouldSkipDir(entry.Name()) {
			continue
		}

		child := filepath.Join(dir, entry.Name())

		if entry.Type()&os.ModeSymlink != 0 {
			if target, ok := s.resolveSymlinkedDir(child, dir); ok {
				s.walk(target, depth+1)
			}
			continue
		}

		if !entry.IsDir() {
			continue
		}
		s.walk(child, depth+1)
	}
}

// resolveSymlinkedDir resolves a symlink entry to a real directory path,
// applying the loop guards. It reports false when the link is broken, does not
// point at a directory, or would take traversal somewhere it must not go.
func (s *scanner) resolveSymlinkedDir(linkPath, currentDir string) (string, bool) {
	target, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		s.result.Errors = append(s.result.Errors,
			fmt.Errorf("broken symlink or inaccessible entry %s: %w", linkPath, err))
		return "", false
	}

	info, err := os.Stat(target)
	if err != nil {
		s.result.Errors = append(s.result.Errors,
			fmt.Errorf("broken symlink or inaccessible entry %s: %w", linkPath, err))
		return "", false
	}
	if !info.IsDir() {
		return "", false
	}

	// Following a link to the workspace root or above it would re-scan the
	// whole tree from the inside.
	if isAncestorOrEqual(target, s.root) {
		return "", false
	}

	// Following a link back up the current branch is a cycle.
	if isAncestorOrEqual(target, currentDir) {
		return "", false
	}

	return target, true
}

// isAncestorOrEqual reports whether candidate is path itself or one of its
// ancestor directories.
func isAncestorOrEqual(candidate, path string) bool {
	c := filepath.Clean(candidate)
	p := filepath.Clean(path)
	if c == p {
		return true
	}
	return strings.HasPrefix(p, c+string(filepath.Separator))
}

// register builds repository metadata for dir and adds it to the result.
// dir must be a resolved real path, which is what makes a repository reachable
// through several paths register exactly once.
func (s *scanner) register(dir string) {
	repo, err := BuildRepository(dir)
	if err != nil {
		s.result.Errors = append(s.result.Errors, fmt.Errorf("failed to parse repository at %s: %w", dir, err))
		return
	}
	s.result.Repositories = append(s.result.Repositories, *repo)
	s.result.ReachablePaths[dir] = true
}

// dirKind describes what a directory is from the scanner's point of view.
type dirKind int

const (
	// dirPlain is an ordinary directory that traversal may descend into.
	dirPlain dirKind = iota
	// dirRepository contains a .git directory and is a repository root.
	dirRepository
	// dirGitLinked contains a .git file: a linked worktree, a submodule, or
	// a clone whose git directory lives elsewhere.
	dirGitLinked
)

// classifyDir reports whether dir is a repository root, a git-linked
// directory, or an ordinary directory.
func classifyDir(dir string) dirKind {
	info, err := os.Lstat(filepath.Join(dir, ".git"))
	if err != nil {
		return dirPlain
	}
	if info.IsDir() {
		return dirRepository
	}
	return dirGitLinked
}

// BuildRepository extracts repository information from a git repository
// directory, classifying its type and visibility from the remote URL.
func BuildRepository(repoPath string) (*config.Repository, error) {
	// Open the repository
	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get repository name from directory
	repoName := filepath.Base(repoPath)

	// Get remote URL (origin)
	remoteURL := ""
	remotes, err := repo.Remotes()
	if err == nil && len(remotes) > 0 {
		// Try to find origin first
		for _, remote := range remotes {
			if remote.Config().Name == "origin" {
				urls := remote.Config().URLs
				if len(urls) > 0 {
					remoteURL = urls[0]
				}
				break
			}
		}
		// If no origin found, use first remote
		if remoteURL == "" && len(remotes) > 0 {
			urls := remotes[0].Config().URLs
			if len(urls) > 0 {
				remoteURL = urls[0]
			}
		}
	}

	gwsRepo := &config.Repository{
		Name:      repoName,
		Path:      repoPath,
		RemoteURL: remoteURL,
		Tags:      []string{},
	}

	// Automatically classify the repository
	classifier.Classify(gwsRepo)

	return gwsRepo, nil
}

// shouldSkipDir determines if a directory should be skipped during scanning.
// Hidden directories are skipped at every level, along with a list of
// directories that commonly contain vendored or generated content.
func shouldSkipDir(dirName string) bool {
	if len(dirName) > 1 && dirName[0] == '.' {
		return true
	}

	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".venv":        true,
		"venv":         true,
		"__pycache__":  true,
		".tox":         true,
		"target":       true, // Rust/Java
		"build":        true,
		"dist":         true,
		".cache":       true,
		".terraform":   true,
	}
	return skipDirs[dirName]
}
