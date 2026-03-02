package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
	"github.com/daileyo/gws/internal/user"
)

// runUserDelete handles the --user --delete/-D flag logic
func runUserDelete(_ *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Find target repositories
	var repos []*config.Repository
	if len(filterTags) > 0 {
		// Tag-based batch mode: filter repos by tag(s)
		matched := filter.Apply(cfg.Repositories, filter.Criteria{Tags: filterTags})
		if len(matched) == 0 {
			return fmt.Errorf("no repositories found with tag(s): %s", strings.Join(filterTags, ", "))
		}
		// Convert to pointers into cfg.Repositories so in-place updates persist
		for i := range cfg.Repositories {
			for _, m := range matched {
				if cfg.Repositories[i].Name == m.Name && cfg.Repositories[i].Path == m.Path {
					repos = append(repos, &cfg.Repositories[i])
					break
				}
			}
		}
	} else {
		// Single/multi repo by identifier
		if len(args) == 0 {
			return fmt.Errorf("usage: gws --user -D <repo>\n       gws --user -D --tag <tag>")
		}
		repoIdentifier := args[0]
		repos = findRepositories(cfg, repoIdentifier)
		if len(repos) == 0 {
			return fmt.Errorf("no repositories found matching: %s", repoIdentifier)
		}
	}

	// Process each matching repo
	deletedCount := 0
	for _, repo := range repos {
		// Capture current config before deletion
		var oldName, oldEmail, oldSource string
		oldCfg, _ := git.GetUserConfig(repo.Path)
		if oldCfg != nil {
			oldName = oldCfg.Name
			oldEmail = oldCfg.Email
			oldSource = string(oldCfg.Source)
		}

		// Skip repos that don't have local config
		if oldCfg == nil || oldCfg.Source != config.UserSourceLocal {
			if !flagQuiet {
				fmt.Printf("%s: no local user config to remove\n", repo.Name)
			}
			continue
		}

		// Delete local config
		if err := user.DeleteLocal(repo.Path, depAll); err != nil {
			if !flagQuiet {
				fmt.Fprintf(os.Stderr, "  %s: failed to delete: %s\n", repo.Name, err)
			}
			continue
		}

		// Re-detect effective config after deletion (will be global or includeIf now)
		newCfg, _ := git.GetUserConfig(repo.Path)
		if newCfg != nil {
			repo.User = newCfg.Name
			repo.Email = newCfg.Email
			repo.SigningEnabled = newCfg.SignCommits
			repo.UserSource = newCfg.Source
		} else {
			repo.User = ""
			repo.Email = ""
			repo.SigningEnabled = false
			repo.UserSource = config.UserSourceUnknown
		}

		deletedCount++

		// Output
		if flagQuiet {
			continue
		}

		if depVerbose {
			fmt.Printf("%s:\n", repo.Name)
			fmt.Printf("  removed: user.name %q, user.email %q (was %s)\n", oldName, oldEmail, oldSource)
			if depAll {
				fmt.Printf("  removed: signing config\n")
			}
			if newCfg != nil && (newCfg.Name != "" || newCfg.Email != "") {
				fmt.Printf("  now using %s: %q <%s>\n", newCfg.Source, newCfg.Name, newCfg.Email)
			}
			fmt.Printf("  path: %s\n", repo.Path)
		} else {
			// Moderate output
			var nowUsing string
			if newCfg != nil && (newCfg.Name != "" || newCfg.Email != "") {
				nowUsing = fmt.Sprintf(" (now using %s: %q <%s>)", newCfg.Source, newCfg.Name, newCfg.Email)
			}
			suffix := ""
			if depAll {
				suffix = " + signing config"
			}
			fmt.Printf("%s: removed local user config%s%s\n", repo.Name, suffix, nowUsing)
		}
	}

	// Save config.json
	if deletedCount > 0 {
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
	}

	// Summary for multiple repos
	if !flagQuiet && len(repos) > 1 {
		fmt.Printf("\nDeleted local user config from %d %s\n",
			deletedCount, pluralize(deletedCount, "repository", "repositories"))
	}

	return nil
}
