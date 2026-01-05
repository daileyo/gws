package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/daileyo/gws/internal/classifier"
	"github.com/daileyo/gws/internal/config"
	"github.com/go-git/go-git/v5"
)

// ScanResult contains the results of a repository scan
type ScanResult struct {
	Repositories []config.Repository
	Count        int
	Errors       []error
}

// Scan recursively scans the given directory for git repositories
// A directory is considered a git repository if it contains a .git directory
func Scan(rootPath string) (*ScanResult, error) {
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

	result := &ScanResult{
		Repositories: []config.Repository{},
		Errors:       []error{},
	}

	// Walk the directory tree
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Record error but continue scanning
			result.Errors = append(result.Errors, fmt.Errorf("error accessing %s: %w", path, err))
			return filepath.SkipDir
		}

		// Check if this is a .git directory
		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			repo, err := parseRepository(repoPath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to parse repository at %s: %w", repoPath, err))
			} else {
				result.Repositories = append(result.Repositories, *repo)
			}
			// Don't descend into .git directories
			return filepath.SkipDir
		}

		// Skip common directories that shouldn't be scanned
		if info.IsDir() && shouldSkipDir(info.Name()) {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory tree: %w", err)
	}

	result.Count = len(result.Repositories)
	return result, nil
}

// parseRepository extracts repository information from a git repository directory
func parseRepository(repoPath string) (*config.Repository, error) {
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

// shouldSkipDir determines if a directory should be skipped during scanning
func shouldSkipDir(dirName string) bool {
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
