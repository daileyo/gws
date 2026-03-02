package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/user"
)

// User subcommand flags for short-flag invocation (user -a / -d / -l / -s).
var (
	userFlagAdd    bool
	userFlagDelete bool
	userFlagList   bool
	userFlagShow   bool
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage git user profiles",
	Long: `Manage git user profiles for different contexts (work, personal, etc).

Profiles define git user.name, user.email, and optional signing configuration.
They can be manually created or auto-detected from your ~/.gitconfig includeIf directives.

Commands:
  gws user list                                    # List all profiles
  gws user add work --email work@company.com       # Add a new profile
  gws user show work                               # Show profile details
  gws user remove work                             # Remove a profile
  gws user assign my-repo work                     # Assign profile to repo
  gws user sync                                    # Sync user info

Short flags:
  gws user -l                                      # List all profiles
  gws user -a work --email work@company.com        # Add a new profile
  gws user -s work                                 # Show profile details
  gws user -d work                                 # Remove a profile`,
	Args:              cobra.ArbitraryArgs,
	ValidArgsFunction: completeProfileThenNone,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Count active short flags for mutual exclusivity
		activeCount := 0
		if userFlagAdd {
			activeCount++
		}
		if userFlagDelete {
			activeCount++
		}
		if userFlagList {
			activeCount++
		}
		if userFlagShow {
			activeCount++
		}

		if activeCount > 1 {
			return fmt.Errorf("only one of -a, -d, -l, -s can be used at a time")
		}

		// Dispatch based on short flag
		if userFlagList {
			return userListCmd.RunE(cmd, args)
		}

		if userFlagAdd {
			if len(args) < 1 {
				return fmt.Errorf("user -a requires a profile name: gws user -a <name> --email <email>")
			}
			return userAddCmd.RunE(cmd, args)
		}

		if userFlagShow {
			if len(args) != 1 {
				return fmt.Errorf("user -s requires exactly 1 argument: gws user -s <name>")
			}
			return userShowCmd.RunE(cmd, args)
		}

		if userFlagDelete {
			if len(args) != 1 {
				return fmt.Errorf("user -d requires exactly 1 argument: gws user -d <name>")
			}
			return userRemoveCmd.RunE(cmd, args)
		}

		// No sub-operation specified — show help
		return cmd.Help()
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all user profiles",
	Long: `List all stored and auto-detected user profiles.

Stored profiles are ones you've added via 'gws user add'.
Auto-detected profiles are discovered from ~/.gitconfig includeIf directives.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Display stored profiles
		if len(stored) > 0 {
			fmt.Printf("Stored Profiles (%d):\n\n", len(stored))
			displayProfileTable(stored)
		}

		// Display detected profiles
		if len(detected) > 0 {
			if len(stored) > 0 {
				fmt.Println()
			}
			fmt.Printf("Auto-Detected Profiles (%d):\n\n", len(detected))
			displayProfileTable(detected)
			fmt.Println("\n(Auto-detected from ~/.gitconfig includeIf directives)")
		}

		return nil
	},
}

// Flags for user add command
var (
	addEmail      string
	addGitName    string
	addSigningKey string
	addSignCommit bool
)

var userAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new user profile",
	Long: `Add a new user profile with the specified name.

The profile name is used to identify this profile when assigning it to repositories.
Common names: work, personal, opensource, client-name, etc.

Examples:
  gws user add work --email work@company.com --name "John Doe"
  gws user add personal --email john@personal.com --name "John"
  gws user add secure --email john@company.com --signing-key ABCD1234 --sign-commits`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		// Validate required email flag
		if addEmail == "" {
			return fmt.Errorf("--email is required")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		profile := config.Profile{
			Name:        profileName,
			GitName:     addGitName,
			Email:       addEmail,
			SigningKey:  addSigningKey,
			SignCommits: addSignCommit,
		}

		// If git name not provided, use profile name as default
		if profile.GitName == "" {
			profile.GitName = profileName
		}

		if err := user.AddProfile(cfg, profile); err != nil {
			return err
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Added profile '%s'\n", profileName)
		fmt.Printf("  Name:  %s\n", profile.GitName)
		fmt.Printf("  Email: %s\n", profile.Email)
		if profile.SigningKey != "" {
			fmt.Printf("  Signing Key: %s\n", profile.SigningKey)
		}
		if profile.SignCommits {
			fmt.Printf("  Sign Commits: yes\n")
		}

		return nil
	},
}

var userShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show profile details",
	Long: `Show detailed information about a user profile.

Examples:
  gws user show work
  gws user show personal`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// First check stored profiles
		profile, err := user.GetProfile(cfg, profileName)
		if err != nil {
			// Try to find in auto-detected profiles
			detected, detectErr := user.DetectProfiles()
			if detectErr != nil {
				return err // Return original error
			}

			for i := range detected {
				if strings.EqualFold(detected[i].Name, profileName) {
					profile = &detected[i]
					break
				}
			}

			if profile == nil {
				return err // Return original error
			}
		}

		fmt.Printf("Profile: %s\n\n", profile.Name)
		fmt.Printf("  Git Name:     %s\n", profile.GitName)
		fmt.Printf("  Email:        %s\n", profile.Email)

		if profile.SigningKey != "" {
			fmt.Printf("  Signing Key:  %s\n", profile.SigningKey)
		} else {
			fmt.Printf("  Signing Key:  (not configured)\n")
		}

		if profile.SignCommits {
			fmt.Printf("  Sign Commits: yes\n")
		} else {
			fmt.Printf("  Sign Commits: no\n")
		}

		// Show repositories using this profile
		repoCount := 0
		for _, repo := range cfg.Repositories {
			if repo.User == profile.GitName && repo.Email == profile.Email {
				repoCount++
			}
		}

		fmt.Printf("\n  Repositories using this profile: %d\n", repoCount)

		return nil
	},
}

var userRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a user profile",
	Long: `Remove a user profile from the configuration.

If repositories are using this profile, you will be prompted for confirmation.
Note: This does not change the git configuration in repositories.

Examples:
  gws user remove work
  gws user remove old-project`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Check if profile exists first
		_, err = user.GetProfile(cfg, profileName)
		if err != nil {
			return err
		}

		// Get repos using this profile (dry run)
		reposUsing, err := user.RemoveProfile(cfg, profileName)
		if err != nil {
			return err
		}

		// If repos are using this profile, ask for confirmation
		if len(reposUsing) > 0 {
			fmt.Printf("Warning: %d %s using this profile:\n",
				len(reposUsing), pluralize(len(reposUsing), "repository is", "repositories are"))
			for _, name := range reposUsing {
				fmt.Printf("  - %s\n", name)
			}
			fmt.Print("\nProceed with removal? [y/N]: ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("Canceled.")
				return nil
			}
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Removed profile '%s'\n", profileName)
		return nil
	},
}

// Flags for user assign command
var (
	assignUseSubdirs bool
	assignDryRun     bool
)

var userAssignCmd = &cobra.Command{
	Use:   "assign <repository> <profile>",
	Short: "Assign a user profile to a repository",
	Long: `Assign a user profile to a repository, setting user.name and user.email
in the repository's local .git/config file.

The repository can be specified by name (partial match) or path.
The profile must exist (either stored or auto-detected).

Options:
  --use-subdirs  Move the repository to a profile subdirectory and set up
                 includeIf in ~/.gitconfig (advanced workflow)
  --dry-run      Preview changes without applying them

Examples:
  gws user assign my-project work         # Set work profile on my-project
  gws user assign my-project personal     # Set personal profile
  gws user assign my-project work --dry-run  # Preview changes`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoIdentifier := args[0]
		profileName := args[1]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Find the repository
		repos := findRepositories(cfg, repoIdentifier)
		if len(repos) == 0 {
			return fmt.Errorf("no repository found matching: %s", repoIdentifier)
		}
		if len(repos) > 1 {
			fmt.Printf("Multiple repositories match '%s':\n", repoIdentifier)
			for _, r := range repos {
				fmt.Printf("  - %s (%s)\n", r.Name, r.Path)
			}
			return fmt.Errorf("please specify a more specific repository identifier")
		}

		repo := repos[0]

		// Find the profile (check stored first, then auto-detected)
		profile, err := user.GetProfile(cfg, profileName)
		if err != nil {
			// Try auto-detected profiles
			detected, detectErr := user.DetectProfiles()
			if detectErr != nil {
				return err
			}
			for i := range detected {
				if strings.EqualFold(detected[i].Name, profileName) {
					profile = &detected[i]
					break
				}
			}
			if profile == nil {
				return err
			}
		}

		// Dry run mode
		if assignDryRun {
			fmt.Printf("Dry run - would assign profile '%s' to repository '%s':\n\n", profileName, repo.Name)

			changes, err := user.PreviewAssignLocal(repo.Path, *profile)
			if err != nil {
				return err
			}

			if len(changes) == 0 {
				fmt.Println("  No changes needed (already configured)")
			} else {
				for _, change := range changes {
					fmt.Printf("  %s\n", change)
				}
			}

			if assignUseSubdirs {
				newPath := filepath.Join(cfg.Workspace, profile.Name, repo.Name)
				fmt.Printf("\n  Would move repository to: %s\n", newPath)
				fmt.Printf("  Would add includeIf directive to ~/.gitconfig\n")
			}

			return nil
		}

		// Subdirectory mode
		if assignUseSubdirs {
			fmt.Printf("Warning: This will move the repository to %s/%s/%s\n",
				cfg.Workspace, profile.Name, repo.Name)
			fmt.Print("Continue? [y/N]: ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("Canceled.")
				return nil
			}

			newPath, err := user.AssignWithSubdirs(cfg, repo.Path, *profile, cfg.Workspace)
			if err != nil {
				return fmt.Errorf("failed to assign with subdirs: %w", err)
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			fmt.Printf("Moved repository to: %s\n", newPath)
			fmt.Printf("Added includeIf directive to ~/.gitconfig\n")
			fmt.Printf("Profile '%s' is now active for this repository\n", profileName)
			return nil
		}

		// Standard local assignment
		if err := user.AssignLocal(repo.Path, *profile); err != nil {
			return fmt.Errorf("failed to assign profile: %w", err)
		}

		// Update stored repo info
		repo.User = profile.GitName
		repo.Email = profile.Email
		repo.SigningEnabled = profile.SignCommits
		repo.UserSource = config.UserSourceLocal

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Assigned profile '%s' to repository '%s'\n", profileName, repo.Name)
		fmt.Printf("  user.name:  %s\n", profile.GitName)
		fmt.Printf("  user.email: %s\n", profile.Email)
		if profile.SignCommits {
			fmt.Printf("  signing:    enabled\n")
		}

		return nil
	},
}

var userSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync stored user info with effective git config",
	Long: `Update stored user information for all repositories to match their
current effective git configuration.

This is useful after manually changing git config files or when repositories
may have been modified outside of gws.

Examples:
  gws user sync                    # Sync all repositories`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Println("Syncing user information for all repositories...")

		updated, err := user.SyncUserInfo(cfg)
		if err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Sync complete. Updated %d %s.\n",
			updated, pluralize(updated, "repository", "repositories"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Add subcommands
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userShowCmd)
	userCmd.AddCommand(userRemoveCmd)
	userCmd.AddCommand(userAssignCmd)
	userCmd.AddCommand(userSyncCmd)

	// Short flags on the user command for quick invocation
	userCmd.Flags().BoolVarP(&userFlagAdd, "add", "a", false, "Add a profile (equivalent to 'user add')")
	userCmd.Flags().BoolVarP(&userFlagDelete, "delete", "d", false, "Remove a profile (equivalent to 'user remove')")
	userCmd.Flags().BoolVarP(&userFlagList, "list", "l", false, "List all profiles (equivalent to 'user list')")
	userCmd.Flags().BoolVarP(&userFlagShow, "show", "s", false, "Show profile details (equivalent to 'user show')")

	// Add-related flags on userCmd so they work with -a (same backing vars as userAddCmd)
	userCmd.Flags().StringVar(&addEmail, "email", "", "Email address for git commits (required with -a)")
	userCmd.Flags().StringVar(&addGitName, "name", "", "Git user name (defaults to profile name)")
	userCmd.Flags().StringVar(&addSigningKey, "signing-key", "", "GPG signing key ID")
	userCmd.Flags().BoolVar(&addSignCommit, "sign-commits", false, "Enable commit signing")

	// Add flags for user add subcommand
	userAddCmd.Flags().StringVar(&addEmail, "email", "", "Email address for git commits (required)")
	userAddCmd.Flags().StringVar(&addGitName, "name", "", "Git user name (defaults to profile name)")
	userAddCmd.Flags().StringVar(&addSigningKey, "signing-key", "", "GPG signing key ID")
	userAddCmd.Flags().BoolVar(&addSignCommit, "sign-commits", false, "Enable commit signing")

	// Add flags for user assign
	userAssignCmd.Flags().BoolVar(&assignUseSubdirs, "use-subdirs", false, "Move repository to profile subdirectory")
	userAssignCmd.Flags().BoolVar(&assignDryRun, "dry-run", false, "Preview changes without applying")

	_ = userAddCmd.MarkFlagRequired("email")

	// Reset to Cobra's default template so user help doesn't inherit root's custom template
	userCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

// completeProfileThenNone completes profile names for the first arg, nothing after.
func completeProfileThenNone(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeProfileNames(toComplete)
}

// completeProfileNames returns profile names (stored + auto-detected) matching the prefix.
func completeProfileNames(toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, err := config.Load()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	seen := make(map[string]bool)
	var names []string

	// Stored profiles
	for _, p := range cfg.Profiles {
		if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(toComplete)) && !seen[p.Name] {
			seen[p.Name] = true
			names = append(names, p.Name)
		}
	}

	// Auto-detected profiles
	detected, err := user.DetectProfiles()
	if err == nil {
		for _, p := range detected {
			if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(toComplete)) && !seen[p.Name] {
				seen[p.Name] = true
				names = append(names, p.Name)
			}
		}
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

// displayProfileTable shows profiles in a formatted table
func displayProfileTable(profiles []config.Profile) {
	// Calculate column widths
	maxNameLen := 4  // "NAME"
	maxGitLen := 8   // "GIT NAME"
	maxEmailLen := 5 // "EMAIL"

	for _, p := range profiles {
		if len(p.Name) > maxNameLen {
			maxNameLen = len(p.Name)
		}
		if len(p.GitName) > maxGitLen {
			maxGitLen = len(p.GitName)
		}
		if len(p.Email) > maxEmailLen {
			maxEmailLen = len(p.Email)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s  %s\n",
		maxNameLen, "NAME",
		maxGitLen, "GIT NAME",
		maxEmailLen, "EMAIL",
		"SIGN")
	fmt.Printf("%s  %s  %s  %s\n",
		strings.Repeat("-", maxNameLen),
		strings.Repeat("-", maxGitLen),
		strings.Repeat("-", maxEmailLen),
		strings.Repeat("-", 4))

	// Print profiles
	for _, p := range profiles {
		signStr := ""
		if p.SignCommits {
			signStr = "✓"
		}

		fmt.Printf("%-*s  %-*s  %-*s  %s\n",
			maxNameLen, p.Name,
			maxGitLen, p.GitName,
			maxEmailLen, p.Email,
			signStr)
	}
}
