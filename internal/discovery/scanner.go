package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"

	"github.com/daileyo/gws/internal/classifier"
	"github.com/daileyo/gws/internal/config"
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

	s := &scanner{
		result: &ScanResult{
			Repositories: []config.Repository{},
			Errors:       []error{},
		},
		maxDepth: effectiveMaxDepth(opts.MaxDepth),
	}

	s.walk(absPath, 0)

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
}

// walk visits dir at the given depth below the scan root. It registers dir as
// a repository and prunes traversal when dir is a repository root, otherwise
// it recurses into eligible subdirectories.
func (s *scanner) walk(dir string, depth int) {
	// A repository root ends traversal down this branch regardless of depth.
	switch classifyDir(dir) {
	case dirRepository:
		s.register(dir)
		return
	case dirGitLinked:
		// A .git file marks a linked worktree or a repository whose git
		// directory lives elsewhere. Neither is traversed into.
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
		if !entry.IsDir() {
			continue
		}
		if shouldSkipDir(entry.Name()) {
			continue
		}
		s.walk(filepath.Join(dir, entry.Name()), depth+1)
	}
}

// register builds repository metadata for dir and adds it to the result.
func (s *scanner) register(dir string) {
	repo, err := BuildRepository(dir)
	if err != nil {
		s.result.Errors = append(s.result.Errors, fmt.Errorf("failed to parse repository at %s: %w", dir, err))
		return
	}
	s.result.Repositories = append(s.result.Repositories, *repo)
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
	repo, err := git.PlainOpen(repoPath)
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
