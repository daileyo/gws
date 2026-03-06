package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/filter"
	"github.com/daileyo/gws/internal/git"
)

// showColumnSentinel is the NoOptDefVal used to distinguish "flag present
// without a value" (show column only) from "flag not present at all".
const showColumnSentinel = "\x00show"

// Dual-purpose flag variables scoped to the list subcommand.
// When the value equals showColumnSentinel the column is shown but no filter
// is applied. Any other non-empty value means filter + show.
var (
	flagType       string
	flagVisibility string
	flagTag        string
	flagPath       string
	flagStatus     string
	flagUser       string
	flagRemote     string
	flagName       string
	outputFormat   string
	verboseCount   int
)

// listCmd is the Cobra subcommand for listing repositories.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked repositories",
	Long: `List all repositories tracked by gws with optional filtering.

By default, only repository names are shown in a compact multi-column layout.
Use flags to add columns or filter results. A flag without a value shows the
column; a flag with a value filters by that value AND shows the column.

Examples:
  gws list                          # Multi-column repo names
  gws list -v                       # Table with type, visibility, tags, path
  gws list -vv                      # Table with all columns (status, user, remote, etc.)
  gws list -yVtp                    # Show type, visibility, tags, path columns
  gws list -y=github                # Filter by GitHub, show type column
  gws list -y=github -V=private     # Filter by GitHub + private
  gws list -y=github -tp            # Filter by GitHub, show type, tags, path
  gws list -n "api-*"               # Filter by name pattern
  gws list -su                      # Show status and user columns
  gws list -r                       # Show remote URL column
  gws list -o json                  # Output as JSON

Note: dual-purpose flags use = for values (e.g., -y=github, --type=github).
Without =, the flag shows the column without filtering.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := parseDualPurposeFlags(cmd)
		return runList(opts)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Dual-purpose flags: NoOptDefVal makes them work as optional-value flags.
	// -y (no value) = show TYPE column; -y github = filter by github + show TYPE
	listCmd.Flags().StringVarP(&flagType, "type", "y", "", "Show type column, or filter by type value")
	listCmd.Flags().Lookup("type").NoOptDefVal = showColumnSentinel

	listCmd.Flags().StringVarP(&flagVisibility, "visibility", "V", "", "Show visibility column, or filter by visibility value")
	listCmd.Flags().Lookup("visibility").NoOptDefVal = showColumnSentinel

	listCmd.Flags().StringVarP(&flagTag, "tag", "t", "", "Show tags column, or filter by tag value (repeatable for AND logic)")
	listCmd.Flags().Lookup("tag").NoOptDefVal = showColumnSentinel

	listCmd.Flags().StringVarP(&flagPath, "path", "p", "", "Show path column, or filter by path pattern")
	listCmd.Flags().Lookup("path").NoOptDefVal = showColumnSentinel

	listCmd.Flags().StringVarP(&flagStatus, "status", "s", "", "Show git status column, or filter by status pattern")
	listCmd.Flags().Lookup("status").NoOptDefVal = showColumnSentinel

	listCmd.Flags().StringVarP(&flagUser, "show-user", "u", "", "Show user/email/sign columns, or filter by user name")
	listCmd.Flags().Lookup("show-user").NoOptDefVal = showColumnSentinel

	listCmd.Flags().StringVarP(&flagRemote, "remote", "r", "", "Show remote URL column, or filter by remote URL pattern")
	listCmd.Flags().Lookup("remote").NoOptDefVal = showColumnSentinel

	// Filter-only flags
	listCmd.Flags().StringVarP(&flagName, "name", "n", "", "Filter by repository name (partial match, supports wildcards)")
	listCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")

	// Verbose count flag: -v = 1, -vv = 2
	listCmd.Flags().CountVarP(&verboseCount, "verbose", "v", "Verbose output (-v stored data, -vv all columns)")

	// Custom help function to clean up NoOptDefVal display artifacts.
	// Cobra renders string flags with NoOptDefVal as: --flag string[="sentinel"]
	// We strip the string[="..."] part for cleaner help output.
	defaultUsages := listCmd.Flags().FlagUsages
	listCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		raw := defaultUsages()
		// Remove all string[="..."] artifacts from NoOptDefVal flags
		cleaned := raw
		for strings.Contains(cleaned, "string[=\"") {
			start := strings.Index(cleaned, "string[=\"")
			end := strings.Index(cleaned[start:], "\"]")
			if end < 0 {
				break
			}
			cleaned = cleaned[:start] + cleaned[start+end+2:]
		}
		// Also clean up any remaining " string" type hints on the cleaned flags
		// (non-NoOptDefVal string flags like --output and --name keep their type hint)

		fmt.Fprint(cmd.OutOrStdout(), cmd.Long+"\n\n")
		fmt.Fprintf(cmd.OutOrStdout(), "Usage:\n  %s\n\nFlags:\n%s", cmd.UseLine(), cleaned)
	})
}

// ListOptions holds all parameters for the list operation.
type ListOptions struct {
	// Filter values (empty = no filter)
	FilterType       string
	FilterVisibility string
	FilterTags       []string
	FilterName       string
	FilterPath       string
	FilterStatus     string
	FilterUser       string
	FilterRemote     string

	// Column display toggles
	ShowType       bool
	ShowVisibility bool
	ShowTags       bool
	ShowPath       bool
	ShowStatus     bool
	ShowUser       bool
	ShowRemote     bool

	// Output
	OutputFormat string
	VerboseLevel int
}

// AnyColumnSelected returns true if any column display flag is set.
func (o ListOptions) AnyColumnSelected() bool {
	return o.ShowType || o.ShowVisibility || o.ShowTags || o.ShowPath ||
		o.ShowStatus || o.ShowUser || o.ShowRemote
}

// parseDualPurposeFlags reads cobra flag state and builds ListOptions.
func parseDualPurposeFlags(cmd *cobra.Command) ListOptions {
	opts := ListOptions{
		OutputFormat: outputFormat,
		VerboseLevel: verboseCount,
	}

	// Helper: for a dual-purpose flag, if Changed and value is sentinel → show only.
	// If Changed and value is not sentinel → filter + show.
	parseDual := func(flagName string, value string) (filterVal string, show bool) {
		if !cmd.Flags().Changed(flagName) {
			return "", false
		}
		if value == showColumnSentinel {
			return "", true
		}
		return value, true
	}

	opts.FilterType, opts.ShowType = parseDual("type", flagType)
	opts.FilterVisibility, opts.ShowVisibility = parseDual("visibility", flagVisibility)
	opts.FilterPath, opts.ShowPath = parseDual("path", flagPath)
	opts.FilterStatus, opts.ShowStatus = parseDual("status", flagStatus)
	opts.FilterUser, opts.ShowUser = parseDual("show-user", flagUser)
	opts.FilterRemote, opts.ShowRemote = parseDual("remote", flagRemote)

	// Tag: dual-purpose. Multiple -t values collected via repeated flag.
	if cmd.Flags().Changed("tag") {
		opts.ShowTags = true
		if flagTag != showColumnSentinel && flagTag != "" {
			opts.FilterTags = []string{flagTag}
		}
	}

	// Name is filter-only
	opts.FilterName = flagName

	// Apply verbose overrides
	if opts.VerboseLevel >= 1 {
		opts.ShowType = true
		opts.ShowVisibility = true
		opts.ShowTags = true
		opts.ShowPath = true
	}
	if opts.VerboseLevel >= 2 {
		opts.ShowStatus = true
		opts.ShowUser = true
		opts.ShowRemote = true
	}

	return opts
}

// runList handles the list logic with explicit options.
func runList(opts ListOptions) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(cfg.Repositories) == 0 {
		if opts.OutputFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Println("No repositories found. Run 'gws init' to discover repositories.")
		}
		return nil
	}

	// Build filter criteria
	criteria := filter.Criteria{
		Type:       opts.FilterType,
		Visibility: opts.FilterVisibility,
		Tags:       opts.FilterTags,
		Name:       opts.FilterName,
		Path:       opts.FilterPath,
	}

	// Apply filters
	filtered := filter.Apply(cfg.Repositories, criteria)

	if len(filtered) == 0 {
		if opts.OutputFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Println("No repositories match the specified filters.")
		}
		return nil
	}

	// Load git status cache if status display is enabled
	var statusCache *git.Cache
	if opts.ShowStatus {
		statusCache = git.NewCache(git.DefaultTTL)
		cachePath, err := git.GetCachePath()
		if err == nil {
			_ = statusCache.Load(cachePath) // Ignore errors, cache might not exist yet
		}
	}

	// Display results based on format
	switch opts.OutputFormat {
	case "json":
		return displayJSON(filtered, opts)
	default:
		displayTable(filtered, statusCache, opts)
	}

	return nil
}

// displayTable shows repositories in a formatted table with selective columns.
func displayTable(repos []config.Repository, statusCache *git.Cache, opts ListOptions) {
	fmt.Printf("Found %d %s:\n\n", len(repos), pluralize(len(repos), "repository", "repositories"))

	// If no columns selected and no verbose, show compact multi-column names
	if !opts.AnyColumnSelected() && opts.VerboseLevel == 0 {
		names := make([]string, len(repos))
		for i, repo := range repos {
			names[i] = repo.Name
		}
		displayMultiColumn(names)
		return
	}

	// Calculate column widths
	maxNameLen := 4   // "NAME"
	maxTypeLen := 4   // "TYPE"
	maxVisLen := 10   // "VISIBILITY"
	maxTagsLen := 4   // "TAGS"
	maxStatusLen := 6 // "STATUS"
	maxUserLen := 4   // "USER"
	maxEmailLen := 5  // "EMAIL"
	maxRemoteLen := 6 // "REMOTE"
	maxPathLen := 4   // "PATH"

	// Pre-compute user info and drift status for width calculation and display
	type repoUserInfo struct {
		userDisplay string
		email       string
		signStr     string
		hasDrift    bool
	}
	userInfoMap := make(map[string]repoUserInfo)

	// Pre-compute remote display strings
	remoteDisplayMap := make(map[string]string)
	if opts.ShowRemote {
		for _, repo := range repos {
			remoteInfo, err := git.GetRemoteInfo(repo.Path)
			if err != nil {
				// Graceful fallback: use stored URL without asterisk
				remoteDisplayMap[repo.Path] = repo.RemoteURL
				continue
			}
			switch {
			case remoteInfo.OriginURL != "" && remoteInfo.HasMultiple:
				remoteDisplayMap[repo.Path] = "* " + remoteInfo.OriginURL
			case remoteInfo.OriginURL != "":
				remoteDisplayMap[repo.Path] = remoteInfo.OriginURL
			case remoteInfo.HasMultiple:
				remoteDisplayMap[repo.Path] = "*"
			default:
				remoteDisplayMap[repo.Path] = repo.RemoteURL
			}
		}
	}

	// Get global default user for marking repos using the default identity
	var globalDefaultUser *git.GlobalUserConfig
	if opts.ShowUser {
		globalDefaultUser, _ = git.GetGlobalDefaultUser()
	}

	for _, repo := range repos {
		if len(repo.Name) > maxNameLen {
			maxNameLen = len(repo.Name)
		}
		if opts.ShowType {
			if len(repo.Type) > maxTypeLen {
				maxTypeLen = len(repo.Type)
			}
		}
		if opts.ShowVisibility {
			if len(repo.Visibility) > maxVisLen {
				maxVisLen = len(repo.Visibility)
			}
		}
		if opts.ShowTags {
			tags := strings.Join(repo.Tags, ", ")
			if len(tags) > maxTagsLen {
				maxTagsLen = len(tags)
			}
		}

		// Calculate path and remote column widths
		if opts.ShowPath || opts.ShowRemote {
			if len(repo.Path) > maxPathLen {
				maxPathLen = len(repo.Path)
			}
		}
		if opts.ShowRemote {
			if rd := remoteDisplayMap[repo.Path]; len(rd) > maxRemoteLen {
				maxRemoteLen = len(rd)
			}
		}

		// Calculate status column width
		if statusCache != nil {
			status, _ := statusCache.GetOrFetch(repo.Path)
			if status != nil {
				statusStr := formatStatusShort(status)
				if len(statusStr) > maxStatusLen {
					maxStatusLen = len(statusStr)
				}
			}
		}

		// Calculate user column widths and pre-compute display values
		if opts.ShowUser {
			info := repoUserInfo{}

			// Start with stored config values
			displayUser := repo.User
			displayEmail := repo.Email
			displaySign := repo.SigningEnabled
			displaySource := repo.UserSource

			// Read effective config at display time to detect local overrides
			currentCfg, err := git.GetUserConfig(repo.Path)
			if err == nil && currentCfg != nil {
				if currentCfg.Source == config.UserSourceLocal {
					// Local config detected - show effective local values
					displayUser = currentCfg.Name
					displayEmail = currentCfg.Email
					displaySign = currentCfg.SignCommits
					displaySource = config.UserSourceLocal
				} else if (repo.User != "" && repo.User != currentCfg.Name) ||
					(repo.Email != "" && repo.Email != currentCfg.Email) {
					info.hasDrift = true
				}
			}

			// Format user display with source marker
			info.userDisplay = displayUser
			if info.userDisplay == "" {
				info.userDisplay = "-"
			} else if displaySource == config.UserSourceLocal {
				info.userDisplay += " (local)"
			} else if displaySource == config.UserSourceGlobal && globalDefaultUser != nil &&
				globalDefaultUser.Name != "" && displayUser == globalDefaultUser.Name {
				info.userDisplay += " *"
			}

			info.email = displayEmail
			if info.email == "" {
				info.email = "-"
			}

			if displaySign {
				info.signStr = "✓"
			} else {
				info.signStr = ""
			}

			userInfoMap[repo.Path] = info

			if len(info.userDisplay) > maxUserLen {
				maxUserLen = len(info.userDisplay)
			}
			if len(info.email) > maxEmailLen {
				maxEmailLen = len(info.email)
			}
		}
	}

	// Build and print header
	var headerParts []string
	var separatorParts []string

	headerParts = append(headerParts, fmt.Sprintf("%-*s", maxNameLen, "NAME"))
	separatorParts = append(separatorParts, strings.Repeat("-", maxNameLen))

	if statusCache != nil {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxStatusLen, "STATUS"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxStatusLen))
	}

	if opts.ShowUser {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxUserLen, "USER"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxUserLen))
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxEmailLen, "EMAIL"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxEmailLen))
		headerParts = append(headerParts, fmt.Sprintf("%-4s", "SIGN"))
		separatorParts = append(separatorParts, strings.Repeat("-", 4))
	}

	if opts.ShowType {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxTypeLen, "TYPE"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxTypeLen))
	}

	if opts.ShowVisibility {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxVisLen, "VISIBILITY"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxVisLen))
	}

	if opts.ShowTags {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxTagsLen, "TAGS"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxTagsLen))
	}

	if opts.ShowPath && opts.ShowRemote {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxPathLen, "PATH"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxPathLen))
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxRemoteLen, "REMOTE"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxRemoteLen))
	} else if opts.ShowPath {
		headerParts = append(headerParts, "PATH")
		separatorParts = append(separatorParts, strings.Repeat("-", 4))
	} else if opts.ShowRemote {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", maxRemoteLen, "REMOTE"))
		separatorParts = append(separatorParts, strings.Repeat("-", maxRemoteLen))
	}

	fmt.Println(strings.Join(headerParts, "  "))
	fmt.Println(strings.Join(separatorParts, "  "))

	// Print repositories
	for _, repo := range repos {
		var rowParts []string

		// Name column - add drift indicator if applicable
		nameDisplay := repo.Name
		if opts.ShowUser {
			if info, ok := userInfoMap[repo.Path]; ok && info.hasDrift {
				nameDisplay += " ⚠"
			}
		}
		rowParts = append(rowParts, fmt.Sprintf("%-*s", maxNameLen, nameDisplay))

		// Status column (optional)
		if statusCache != nil {
			status, _ := statusCache.GetOrFetch(repo.Path)
			statusStr := "-"
			if status != nil {
				statusStr = formatStatusShort(status)
			}
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxStatusLen, statusStr))
		}

		// User columns (optional)
		if opts.ShowUser {
			info := userInfoMap[repo.Path]
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxUserLen, info.userDisplay))
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxEmailLen, info.email))
			rowParts = append(rowParts, fmt.Sprintf("%-4s", info.signStr))
		}

		// Type column (optional)
		if opts.ShowType {
			typeStr := string(repo.Type)
			if typeStr == "" {
				typeStr = "unknown"
			}
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxTypeLen, typeStr))
		}

		// Visibility column (optional)
		if opts.ShowVisibility {
			visStr := string(repo.Visibility)
			if visStr == "" {
				visStr = "unknown"
			}
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxVisLen, visStr))
		}

		// Tags column (optional)
		if opts.ShowTags {
			tags := strings.Join(repo.Tags, ", ")
			if tags == "" {
				tags = "-"
			}
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxTagsLen, tags))
		}

		// Path column (optional, padded when remote column follows)
		if opts.ShowPath && opts.ShowRemote {
			rowParts = append(rowParts, fmt.Sprintf("%-*s", maxPathLen, repo.Path))
			rowParts = append(rowParts, remoteDisplayMap[repo.Path])
		} else if opts.ShowPath {
			rowParts = append(rowParts, repo.Path)
		} else if opts.ShowRemote {
			rowParts = append(rowParts, remoteDisplayMap[repo.Path])
		}

		fmt.Println(strings.Join(rowParts, "  "))
	}

	// Save cache if we fetched any statuses
	if statusCache != nil {
		cachePath, err := git.GetCachePath()
		if err == nil {
			_ = statusCache.Save(cachePath) // Ignore save errors
		}
	}
}

// getTerminalWidth returns the terminal width, or 80 if detection fails.
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 80
	}
	return width
}

// displayMultiColumn prints names in a compact multi-column layout like ls.
// Names are sorted alphabetically and fill top-to-bottom, then left-to-right
// (column-first order). Falls back to single-column when stdout is not a TTY.
func displayMultiColumn(names []string) {
	if len(names) == 0 {
		return
	}

	sort.Strings(names)

	// Non-TTY: one name per line
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		for _, name := range names {
			fmt.Println(name)
		}
		return
	}

	termWidth := getTerminalWidth()
	colGap := 2

	// Find the widest name
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	// Calculate number of columns that fit
	colWidth := maxLen + colGap
	numCols := termWidth / colWidth
	if numCols < 1 {
		numCols = 1
	}
	// Cap columns so we get at least 3 rows for readability,
	// unless there are very few items (≤3 always single column).
	if len(names) > 3 {
		maxCols := len(names) / 3
		if maxCols < 1 {
			maxCols = 1
		}
		if numCols > maxCols {
			numCols = maxCols
		}
	}

	// Calculate rows needed (column-first layout)
	numRows := (len(names) + numCols - 1) / numCols

	// Print in column-first order: row by row, picking items that belong
	// to each column. Item at column c, row r = names[c*numRows + r].
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			idx := c*numRows + r
			if idx >= len(names) {
				break
			}
			// Last column in this row or last item: no trailing padding
			isLastCol := c == numCols-1 || (c+1)*numRows+r >= len(names)
			if isLastCol {
				fmt.Print(names[idx])
			} else {
				fmt.Printf("%-*s", colWidth, names[idx])
			}
		}
		fmt.Println()
	}
}

// formatStatusShort returns a short status string with visual indicators
func formatStatusShort(status *git.Status) string {
	if status.Branch == "" {
		return "no commits"
	}

	result := status.Branch

	// Add clean/dirty indicator
	if status.HasChanges {
		result += " ✗"
	} else {
		result += " ✓"
	}

	// Add ahead/behind indicators
	if status.Ahead > 0 && status.Behind > 0 {
		result += fmt.Sprintf(" ↑%d↓%d", status.Ahead, status.Behind)
	} else if status.Ahead > 0 {
		result += fmt.Sprintf(" ↑%d", status.Ahead)
	} else if status.Behind > 0 {
		result += fmt.Sprintf(" ↓%d", status.Behind)
	}

	return result
}

// displayJSON outputs repositories in JSON format respecting column selection.
// Only fields for displayed columns are included; "name" is always present.
func displayJSON(repos []config.Repository, opts ListOptions) error {
	entries := make([]map[string]interface{}, len(repos))

	// Load status cache if needed
	var statusCache *git.Cache
	if opts.ShowStatus {
		statusCache = git.NewCache(git.DefaultTTL)
		cachePath, err := git.GetCachePath()
		if err == nil {
			_ = statusCache.Load(cachePath)
		}
	}

	for i, repo := range repos {
		entry := map[string]interface{}{
			"name": repo.Name,
		}

		if opts.ShowType {
			t := string(repo.Type)
			if t == "" {
				t = "unknown"
			}
			entry["type"] = t
		}
		if opts.ShowVisibility {
			v := string(repo.Visibility)
			if v == "" {
				v = "unknown"
			}
			entry["visibility"] = v
		}
		if opts.ShowTags {
			if repo.Tags != nil {
				entry["tags"] = repo.Tags
			} else {
				entry["tags"] = []string{}
			}
		}
		if opts.ShowPath {
			entry["path"] = repo.Path
		}
		if opts.ShowStatus && statusCache != nil {
			status, _ := statusCache.GetOrFetch(repo.Path)
			if status != nil {
				entry["status"] = formatStatusShort(status)
			} else {
				entry["status"] = "-"
			}
		}
		if opts.ShowUser {
			entry["user"] = repo.User
			entry["email"] = repo.Email
			entry["signing_enabled"] = repo.SigningEnabled
		}
		if opts.ShowRemote {
			remoteInfo, err := git.GetRemoteInfo(repo.Path)
			if err == nil {
				if remoteInfo.OriginURL != "" {
					entry["remote_url"] = remoteInfo.OriginURL
				} else {
					entry["remote_url"] = repo.RemoteURL
				}
				entry["has_multiple_remotes"] = remoteInfo.HasMultiple
			} else {
				entry["remote_url"] = repo.RemoteURL
				entry["has_multiple_remotes"] = false
			}
		}

		entries[i] = entry
	}

	// Save status cache if used
	if statusCache != nil {
		cachePath, err := git.GetCachePath()
		if err == nil {
			_ = statusCache.Save(cachePath)
		}
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
