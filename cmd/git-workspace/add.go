package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/classifier"
	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/discovery"
)

// addRecursive is the --recursive flag for the add subcommand.
var addRecursive bool

// addCmd is the Cobra subcommand for adding repositories.
var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a git repository to the workspace",
	Long: `Add a git repository to the workspace. Defaults to the current directory if no path is given.
Use --recursive to scan and add all git repositories found under the given path.

Examples:
  gws add                          # Add current directory
  gws add ~/projects/my-repo       # Add a specific repo
  gws add -v                       # Recursively add all repos in current directory
  gws add ~/projects -v            # Recursively add all repos under a path`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		return runAdd(path, addRecursive)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&addRecursive, "recursive", "v", false, "Recursively add all git repositories found under the path")
}

// runAdd handles the add logic for adding a single repository.
// path is the directory to add (use "." for current directory).
// recursive controls whether to scan recursively.
func runAdd(path string, recursive bool) error {
	if recursive {
		return runAddRecursive()
	}

	// Guard: workspace must be initialized before adding repos
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check workspace status: %w", err)
	}
	if !exists {
		fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Run gws init first to create a workspace.")
		return fmt.Errorf("workspace not initialized")
	}

	// Resolve target path: "." means current working directory
	targetPath := path
	if targetPath == "." {
		targetPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to resolve current directory: %w", err)
		}
	}

	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Validate target is a git repository
	gitDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s is not a git repository\n", absPath)
		return fmt.Errorf("%s is not a git repository", absPath)
	} else if err != nil {
		return fmt.Errorf("failed to check git directory: %w", err)
	}

	// Load config and check for duplicate
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load workspace configuration: %w", err)
	}

	repoName := filepath.Base(absPath)

	for _, repo := range cfg.Repositories {
		if repo.Path == absPath {
			fmt.Printf("%s is already tracked, skipping.\n", repoName)
			return nil
		}
	}

	// Extract repository metadata using go-git
	newRepo, err := buildRepository(absPath)
	if err != nil {
		return fmt.Errorf("failed to read repository metadata: %w", err)
	}

	// Detect git user configuration and auto-link to profiles
	singleRepo := []config.Repository{*newRepo}
	detectUserForRepos(singleRepo, cfg.Profiles)
	*newRepo = singleRepo[0]

	// Sync profiles from detected repo users
	syncProfilesFromRepos(cfg)

	// Create symlink in workspace if the repo is outside the workspace directory
	symlinkCreated, err := createSymlinkIfExternal(cfg, absPath, repoName)
	if err != nil {
		// Non-fatal: warn and continue
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	// Save updated config
	cfg.Repositories = append(cfg.Repositories, *newRepo)
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Added %s to workspace.\n", repoName)
	if symlinkCreated {
		symlinkPath := filepath.Join(cfg.Workspace, repoName)
		fmt.Printf("Created symlink: %s → %s\n", symlinkPath, absPath)
	}

	return nil
}

// runAddRecursive handles the --add --recursive path: scans CWD for all git
// repositories and adds any that are not already tracked.
func runAddRecursive() error {
	// Guard: workspace must be initialized before adding repos
	exists, err := config.Exists()
	if err != nil {
		return fmt.Errorf("failed to check workspace status: %w", err)
	}
	if !exists {
		fmt.Fprintln(os.Stderr, "Error: workspace not initialized")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Run gws --init first to create a workspace.")
		return fmt.Errorf("workspace not initialized")
	}

	// Recursive always operates on the current directory
	scanRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to resolve current directory: %w", err)
	}

	// Scan for all git repositories under CWD
	result, err := discovery.Scan(scanRoot)
	if err != nil {
		return fmt.Errorf("failed to scan for repositories: %w", err)
	}

	// Load existing config and build a fast-lookup set of tracked paths
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load workspace configuration: %w", err)
	}

	trackedPaths := make(map[string]bool, len(cfg.Repositories))
	for _, repo := range cfg.Repositories {
		trackedPaths[repo.Path] = true
	}

	var newRepos []config.Repository
	var skippedNames []string

	for _, repo := range result.Repositories {
		if trackedPaths[repo.Path] {
			skippedNames = append(skippedNames, repo.Name)
			continue
		}

		// Create symlink if repo is outside the workspace directory
		if _, symErr := createSymlinkIfExternal(cfg, repo.Path, repo.Name); symErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", symErr)
		}

		newRepos = append(newRepos, repo)
	}

	// Warn about already-tracked repos
	for _, name := range skippedNames {
		fmt.Printf("%s is already tracked, skipping.\n", name)
	}

	if len(newRepos) == 0 {
		fmt.Println("No new repositories found.")
		return nil
	}

	// Detect git user configuration and auto-link to profiles
	detectUserForRepos(newRepos, cfg.Profiles)

	// Append new repos and sync profiles from detected users
	cfg.Repositories = append(cfg.Repositories, newRepos...)
	syncProfilesFromRepos(cfg)

	// Save updated config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Print summary: "Added N repositories: repo-a, repo-b, repo-c."
	names := make([]string, len(newRepos))
	for i, repo := range newRepos {
		names[i] = repo.Name
	}
	fmt.Printf("Added %d %s: %s.\n",
		len(newRepos),
		pluralize(len(newRepos), "repository", "repositories"),
		strings.Join(names, ", "),
	)

	return nil
}

// buildRepository opens a git repo at absPath and returns a populated config.Repository.
func buildRepository(absPath string) (*config.Repository, error) {
	repo, err := gogit.PlainOpen(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	remoteURL := ""
	remotes, err := repo.Remotes()
	if err == nil {
		for _, remote := range remotes {
			if remote.Config().Name == "origin" {
				if urls := remote.Config().URLs; len(urls) > 0 {
					remoteURL = urls[0]
				}
				break
			}
		}
		if remoteURL == "" && len(remotes) > 0 {
			if urls := remotes[0].Config().URLs; len(urls) > 0 {
				remoteURL = urls[0]
			}
		}
	}

	gwsRepo := &config.Repository{
		Name:      filepath.Base(absPath),
		Path:      absPath,
		RemoteURL: remoteURL,
		Tags:      []string{},
	}
	classifier.Classify(gwsRepo)

	return gwsRepo, nil
}

// createSymlinkIfExternal creates a symlink inside the workspace directory when
// the target repository lives outside the workspace root. Returns true if a
// symlink was created. A non-nil error indicates a skip condition (already
// exists) rather than a fatal failure.
func createSymlinkIfExternal(cfg *config.Config, repoAbsPath, repoName string) (bool, error) {
	// Repos inside the workspace don't need a symlink
	workspacePath := filepath.Clean(cfg.Workspace)
	repoPath := filepath.Clean(repoAbsPath)
	if strings.HasPrefix(repoPath, workspacePath+string(filepath.Separator)) || repoPath == workspacePath {
		return false, nil
	}

	symlinkPath := filepath.Join(cfg.Workspace, repoName)

	// Check if a file or symlink already exists at the target
	if _, err := os.Lstat(symlinkPath); err == nil {
		return false, fmt.Errorf("symlink path %s already exists, skipping symlink creation", symlinkPath)
	}

	if err := os.Symlink(repoAbsPath, symlinkPath); err != nil {
		return false, fmt.Errorf("failed to create symlink at %s: %w", symlinkPath, err)
	}

	return true, nil
}
