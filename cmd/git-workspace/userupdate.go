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

// runListUsers handles --user (alone) and --list-users: displays available profiles
func runListUsers(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	stored, detected := user.ListProfiles(cfg)

	if len(stored) == 0 && len(detected) == 0 {
		fmt.Println("No profiles found.")
		fmt.Println("\nAdd a profile with: gws user add <name> --email <email>")
		return nil
	}

	if len(stored) > 0 {
		fmt.Printf("Stored Profiles (%d):\n\n", len(stored))
		displayProfileTable(stored)
	}

	if len(detected) > 0 {
		if len(stored) > 0 {
			fmt.Println()
		}
		fmt.Printf("Auto-Detected Profiles (%d):\n\n", len(detected))
		displayProfileTable(detected)
		fmt.Println("\n(Auto-detected from ~/.gitconfig includeIf directives)")
	}

	return nil
}

// resolveProfile resolves the target profile from args and inline flags.
// If a profile name is provided, it looks up stored then auto-detected profiles.
// If inline --git-name/--git-email are provided, they override the profile values.
// Returns an error if neither a profile name nor inline email is provided.
func resolveProfile(cfg *config.Config, args []string) (config.Profile, error) {
	var profile config.Profile
	hasProfileName := false

	// Check for profile name in positional args
	// Args layout: [repo-identifier] [profile-name] or just [profile-name] when using --tag
	profileName := ""
	if len(filterTags) > 0 {
		// With --tag, first arg is profile name (no repo identifier needed)
		if len(args) > 0 {
			profileName = args[0]
		}
	} else {
		// Without --tag, second arg is profile name (first is repo identifier)
		if len(args) > 1 {
			profileName = args[1]
		}
	}

	if profileName != "" {
		// Look up stored profile first
		p, err := user.GetProfile(cfg, profileName)
		if err != nil {
			// Try auto-detected profiles
			detected, detectErr := user.DetectProfiles()
			if detectErr == nil {
				for i := range detected {
					if strings.EqualFold(detected[i].Name, profileName) {
						p = &detected[i]
						break
					}
				}
			}
			if p == nil {
				return profile, fmt.Errorf("profile '%s' not found", profileName)
			}
		}
		profile = *p
		hasProfileName = true
	}

	// Apply inline overrides
	if depInlineName != "" {
		profile.GitName = depInlineName
	}
	if depInlineEmail != "" {
		profile.Email = depInlineEmail
	}

	// Validate we have enough info
	if !hasProfileName && depInlineEmail == "" {
		hint := "\n\nUse 'gws --user' to see available profiles."
		if len(filterTags) > 0 {
			return profile, fmt.Errorf("usage: gws --user -u --tag <tag> <profile>\n       gws --user -u --tag <tag> --git-email <email> [--git-name <name>]" + hint)
		}
		return profile, fmt.Errorf("usage: gws --user -u <repo> <profile>\n       gws --user -u <repo> --git-email <email> [--git-name <name>]" + hint)
	}

	// If only inline values with no git name, use email prefix as name
	if profile.GitName == "" && profile.Email != "" {
		parts := strings.SplitN(profile.Email, "@", 2)
		profile.GitName = parts[0]
	}

	return profile, nil
}

// runUserUpdate handles the --user --update/-u flag logic
func runUserUpdate(_ *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Resolve target profile
	profile, err := resolveProfile(cfg, args)
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
			return fmt.Errorf("usage: gws --user -u <repo> <profile>\n       gws --user -u <repo> --git-email <email> [--git-name <name>]")
		}
		repoIdentifier := args[0]
		repos = findRepositories(cfg, repoIdentifier)
		if len(repos) == 0 {
			return fmt.Errorf("no repositories found matching: %s", repoIdentifier)
		}
	}

	// Process each matching repo
	updatedCount := 0
	for _, repo := range repos {
		// Capture current config for comparison
		var oldName, oldEmail string
		oldCfg, _ := git.GetUserConfig(repo.Path)
		if oldCfg != nil {
			oldName = oldCfg.Name
			oldEmail = oldCfg.Email
		}

		// Apply the profile to repo's local .git/config
		if err := user.AssignLocal(repo.Path, profile); err != nil {
			if !flagQuiet {
				fmt.Fprintf(os.Stderr, "  %s: failed to update: %s\n", repo.Name, err)
			}
			continue
		}

		// Update config.json fields
		repo.User = profile.GitName
		repo.Email = profile.Email
		repo.SigningEnabled = profile.SignCommits
		repo.UserSource = config.UserSourceLocal

		updatedCount++

		// Output
		if flagQuiet {
			continue
		}

		if depVerbose {
			fmt.Printf("%s:\n", repo.Name)
			fmt.Printf("  user.name:  %q → %q\n", oldName, profile.GitName)
			fmt.Printf("  user.email: %q → %q\n", oldEmail, profile.Email)
			if profile.SigningKey != "" {
				fmt.Printf("  signingkey: %s\n", profile.SigningKey)
			}
			if profile.SignCommits {
				fmt.Printf("  gpgsign:    true\n")
			}
			fmt.Printf("  path:       %s\n", repo.Path)
		} else {
			// Moderate output
			var changes []string
			if oldName != profile.GitName {
				changes = append(changes, fmt.Sprintf("user.name %q → %q", oldName, profile.GitName))
			}
			if oldEmail != profile.Email {
				changes = append(changes, fmt.Sprintf("user.email %q → %q", oldEmail, profile.Email))
			}
			if len(changes) > 0 {
				fmt.Printf("%s: %s\n", repo.Name, strings.Join(changes, ", "))
			} else {
				fmt.Printf("%s: no changes (already configured)\n", repo.Name)
			}
		}
	}

	// Save config.json
	if updatedCount > 0 {
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
	}

	// Summary for multiple repos
	if !flagQuiet && len(repos) > 1 {
		fmt.Printf("\nUpdated %d %s\n", updatedCount, pluralize(updatedCount, "repository", "repositories"))
	}

	return nil
}
