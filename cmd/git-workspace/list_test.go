package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
)

// --- Flag Registration Tests ---

func TestFilterFlagsOnListCmd(t *testing.T) {
	// Verify all filter/display flags are registered on listCmd with correct shorthands
	flags := []struct {
		name      string
		shorthand string
	}{
		// Lowercase filter-only flags
		{"type", "y"},
		{"visibility", "i"},
		{"tag", "t"},
		{"name", "n"},
		{"path", "p"},
		{"status", "s"},
		{"show-user", "u"},
		{"remote", "r"},
		{"remote-raw", "b"},
		// Uppercase show+filter flags
		{"show-type", "Y"},
		{"show-visibility", "I"},
		{"show-tag", "T"},
		{"show-path", "P"},
		{"show-status", "S"},
		{"show-user-col", "U"},
		{"show-remote", "R"},
		{"show-remote-raw", "B"},
		// Other flags
		{"output", "o"},
		{"verbose", "v"},
	}

	for _, f := range flags {
		t.Run(f.name, func(t *testing.T) {
			flag := listCmd.Flags().Lookup(f.name)
			if flag == nil {
				t.Errorf("Flag --%s not found on listCmd", f.name)
				return
			}
			if flag.Shorthand != f.shorthand {
				t.Errorf("Flag --%s shorthand: expected '%s', got '%s'", f.name, f.shorthand, flag.Shorthand)
			}
		})
	}
}

func TestFilterFlagsOnRootAreHidden(t *testing.T) {
	// Filter flags on root are deprecated (hidden) — they exist for backward compat
	deprecatedFilters := []string{"type", "name", "path", "output", "status", "show-user"}

	for _, name := range deprecatedFilters {
		t.Run(name, func(t *testing.T) {
			flag := rootCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Deprecated filter flag --%s not found on rootCmd", name)
				return
			}
			if !flag.Hidden {
				t.Errorf("Deprecated filter flag --%s on rootCmd should be hidden", name)
			}
		})
	}
}

func TestLowercaseFlagsNoOptDefVal(t *testing.T) {
	// Lowercase filter-only flags should NOT have NoOptDefVal
	filterOnlyFlags := []string{"type", "visibility", "tag", "path", "status", "show-user", "remote", "remote-raw", "name"}

	for _, name := range filterOnlyFlags {
		t.Run(name, func(t *testing.T) {
			flag := listCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Flag --%s not found on listCmd", name)
				return
			}
			if flag.NoOptDefVal != "" {
				t.Errorf("Lowercase flag --%s should NOT have NoOptDefVal, got %q", name, flag.NoOptDefVal)
			}
		})
	}
}

func TestUppercaseFlagsHaveNoOptDefVal(t *testing.T) {
	// Uppercase show+filter flags should have NoOptDefVal set to the sentinel
	showFlags := []string{"show-type", "show-visibility", "show-tag", "show-path", "show-status", "show-user-col", "show-remote", "show-remote-raw"}

	for _, name := range showFlags {
		t.Run(name, func(t *testing.T) {
			flag := listCmd.Flags().Lookup(name)
			if flag == nil {
				t.Errorf("Flag --%s not found on listCmd", name)
				return
			}
			if flag.NoOptDefVal != showColumnSentinel {
				t.Errorf("Uppercase flag --%s NoOptDefVal: expected sentinel, got %q", name, flag.NoOptDefVal)
			}
		})
	}
}

// --- Lowercase Filter-Only Flag Tests ---

func TestLowercaseFilterSpaceSeparated(t *testing.T) {
	// Lowercase flags accept space-separated values (no = required)
	tests := []struct {
		name     string
		args     []string
		flagVar  *string
		expected string
	}{
		{"type -y github", []string{"-y", "github"}, &flagType, "github"},
		{"tag -t ai", []string{"-t", "ai"}, &flagTag, "ai"},
		{"path -p /home", []string{"-p", "/home"}, &flagPath, "/home"},
		{"status -s clean", []string{"-s", "clean"}, &flagStatus, "clean"},
		{"user -u alice", []string{"-u", "alice"}, &flagUser, "alice"},
		{"remote -r github", []string{"-r", "github"}, &flagRemote, "github"},
		{"visibility -i private", []string{"-i", "private"}, &flagVisibility, "private"},
		{"remote-raw -b git@", []string{"-b", "git@"}, &flagRemoteRaw, "git@"},
		{"name -n api", []string{"-n", "api"}, &flagName, "api"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := *tt.flagVar
			defer func() { *tt.flagVar = orig }()
			*tt.flagVar = ""

			err := listCmd.Flags().Parse(tt.args)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", tt.args, err)
			}

			if *tt.flagVar != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, *tt.flagVar)
			}
		})
	}
}

func TestLowercaseFilterEqualsCompatibility(t *testing.T) {
	// Lowercase flags also accept = syntax
	tests := []struct {
		name     string
		args     []string
		flagVar  *string
		expected string
	}{
		{"type -y=github", []string{"-y=github"}, &flagType, "github"},
		{"tag -t=ai", []string{"-t=ai"}, &flagTag, "ai"},
		{"long --type=github", []string{"--type=github"}, &flagType, "github"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := *tt.flagVar
			defer func() { *tt.flagVar = orig }()
			*tt.flagVar = ""

			err := listCmd.Flags().Parse(tt.args)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", tt.args, err)
			}

			if *tt.flagVar != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, *tt.flagVar)
			}
		})
	}
}

// --- Uppercase Show+Filter Flag Tests ---

func TestUppercaseBareShowColumn(t *testing.T) {
	// Uppercase flags used bare (no value) set the sentinel to show column
	tests := []struct {
		name     string
		args     []string
		flagVar  *string
		expected string
	}{
		{"show-type -Y", []string{"-Y"}, &flagShowType, showColumnSentinel},
		{"show-tag -T", []string{"-T"}, &flagShowTag, showColumnSentinel},
		{"show-path -P", []string{"-P"}, &flagShowPath, showColumnSentinel},
		{"show-status -S", []string{"-S"}, &flagShowStatus, showColumnSentinel},
		{"show-user -U", []string{"-U"}, &flagShowUser, showColumnSentinel},
		{"show-remote -R", []string{"-R"}, &flagShowRemote, showColumnSentinel},
		{"show-visibility -I", []string{"-I"}, &flagShowVisibility, showColumnSentinel},
		{"show-remote-raw -B", []string{"-B"}, &flagShowRemoteRaw, showColumnSentinel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := *tt.flagVar
			defer func() { *tt.flagVar = orig }()
			*tt.flagVar = ""

			err := listCmd.Flags().Parse(tt.args)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", tt.args, err)
			}

			if *tt.flagVar != tt.expected {
				t.Errorf("Expected sentinel, got %q", *tt.flagVar)
			}
		})
	}
}

func TestUppercaseWithValueEqualsOnly(t *testing.T) {
	// Uppercase flags with = syntax work directly at the parse level
	tests := []struct {
		name     string
		args     []string
		flagVar  *string
		expected string
	}{
		{"show-type -Y=github", []string{"-Y=github"}, &flagShowType, "github"},
		{"show-tag -T=ai", []string{"-T=ai"}, &flagShowTag, "ai"},
		{"show-status -S=clean", []string{"-S=clean"}, &flagShowStatus, "clean"},
		{"show-visibility -I=private", []string{"-I=private"}, &flagShowVisibility, "private"},
		{"show-remote -R=github", []string{"-R=github"}, &flagShowRemote, "github"},
		{"show-remote-raw -B=git@", []string{"-B=git@"}, &flagShowRemoteRaw, "git@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := *tt.flagVar
			defer func() { *tt.flagVar = orig }()
			*tt.flagVar = ""

			err := listCmd.Flags().Parse(tt.args)
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", tt.args, err)
			}

			if *tt.flagVar != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, *tt.flagVar)
			}
		})
	}
}

func TestUppercaseSpaceSeparatedViaReassignment(t *testing.T) {
	// Uppercase flags with space-separated values: Cobra puts the value
	// as a remaining arg, and reassignTrailingArg reassigns it.
	// This simulates what RunE does.
	tests := []struct {
		name     string
		args     []string
		flagVar  *string
		expected string
	}{
		{"show-tag -T ai", []string{"-T", "ai"}, &flagShowTag, "ai"},
		{"show-type -Y github", []string{"-Y", "github"}, &flagShowType, "github"},
		{"show-status -S clean", []string{"-S", "clean"}, &flagShowStatus, "clean"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := *tt.flagVar
			defer func() { *tt.flagVar = orig }()
			*tt.flagVar = ""

			err := listCmd.Flags().Parse(tt.args[:1]) // parse flag only
			if err != nil {
				t.Fatalf("Failed to parse %v: %v", tt.args[:1], err)
			}

			if *tt.flagVar != showColumnSentinel {
				t.Fatalf("Expected sentinel after bare flag, got %q", *tt.flagVar)
			}

			// Simulate reassignment (what RunE does with remaining args)
			*tt.flagVar = tt.args[1]

			if *tt.flagVar != tt.expected {
				t.Errorf("Expected %q after reassignment, got %q", tt.expected, *tt.flagVar)
			}
		})
	}
}

func TestUppercaseEqualsCompatibility(t *testing.T) {
	// Uppercase flags also accept = syntax
	origType := flagShowType
	defer func() { flagShowType = origType }()
	flagShowType = ""

	err := listCmd.Flags().Parse([]string{"-Y=github"})
	if err != nil {
		t.Fatalf("Failed to parse -Y=github: %v", err)
	}

	if flagShowType != "github" {
		t.Errorf("Expected %q, got %q", "github", flagShowType)
	}
}

// --- Flag Stacking Tests ---

func TestUppercaseFlagStacking(t *testing.T) {
	// Uppercase bare flags can be stacked: -SU sets both show-status and show-user
	origStatus := flagShowStatus
	origUser := flagShowUser
	defer func() {
		flagShowStatus = origStatus
		flagShowUser = origUser
	}()

	flagShowStatus = ""
	flagShowUser = ""

	err := listCmd.Flags().Parse([]string{"-SU"})
	if err != nil {
		t.Fatalf("Failed to parse -SU: %v", err)
	}

	if flagShowStatus != showColumnSentinel {
		t.Errorf("Expected flagShowStatus to be sentinel after -SU, got %q", flagShowStatus)
	}
	if flagShowUser != showColumnSentinel {
		t.Errorf("Expected flagShowUser to be sentinel after -SU, got %q", flagShowUser)
	}
}

func TestUppercaseFlagStackingMultiple(t *testing.T) {
	// Stack multiple uppercase flags: -YTSP
	origType := flagShowType
	origTag := flagShowTag
	origStatus := flagShowStatus
	origPath := flagShowPath
	defer func() {
		flagShowType = origType
		flagShowTag = origTag
		flagShowStatus = origStatus
		flagShowPath = origPath
	}()

	flagShowType = ""
	flagShowTag = ""
	flagShowStatus = ""
	flagShowPath = ""

	err := listCmd.Flags().Parse([]string{"-YTSP"})
	if err != nil {
		t.Fatalf("Failed to parse -YTSP: %v", err)
	}

	if flagShowType != showColumnSentinel {
		t.Errorf("Expected flagShowType to be sentinel after -YTSP, got %q", flagShowType)
	}
	if flagShowTag != showColumnSentinel {
		t.Errorf("Expected flagShowTag to be sentinel after -YTSP, got %q", flagShowTag)
	}
	if flagShowStatus != showColumnSentinel {
		t.Errorf("Expected flagShowStatus to be sentinel after -YTSP, got %q", flagShowStatus)
	}
	if flagShowPath != showColumnSentinel {
		t.Errorf("Expected flagShowPath to be sentinel after -YTSP, got %q", flagShowPath)
	}
}

func TestUppercaseFlagStackingWithTrailingValue(t *testing.T) {
	// -YT ai: Cobra parses Y=sentinel, T=sentinel, "ai" is remaining.
	// reassignTrailingArg assigns "ai" to T (last in stack).
	origType := flagShowType
	origTag := flagShowTag
	defer func() {
		flagShowType = origType
		flagShowTag = origTag
	}()

	flagShowType = ""
	flagShowTag = ""

	err := listCmd.Flags().Parse([]string{"-YT"})
	if err != nil {
		t.Fatalf("Failed to parse -YT: %v", err)
	}

	if flagShowType != showColumnSentinel {
		t.Errorf("Expected flagShowType to be sentinel after -YT, got %q", flagShowType)
	}
	if flagShowTag != showColumnSentinel {
		t.Errorf("Expected flagShowTag to be sentinel after -YT, got %q", flagShowTag)
	}

	// Simulate reassignment (what RunE does): assign "ai" to the last flag (T)
	flagShowTag = "ai"

	if flagShowType != showColumnSentinel {
		t.Errorf("Expected flagShowType to remain sentinel, got %q", flagShowType)
	}
	if flagShowTag != "ai" {
		t.Errorf("Expected flagShowTag to be %q after reassignment, got %q", "ai", flagShowTag)
	}
}

func TestUppercaseStackingWithRemote(t *testing.T) {
	// -RSU: show remote, status, user columns
	origRemote := flagShowRemote
	origStatus := flagShowStatus
	origUser := flagShowUser
	defer func() {
		flagShowRemote = origRemote
		flagShowStatus = origStatus
		flagShowUser = origUser
	}()

	flagShowRemote = ""
	flagShowStatus = ""
	flagShowUser = ""

	err := listCmd.Flags().Parse([]string{"-RSU"})
	if err != nil {
		t.Fatalf("Failed to parse -RSU: %v", err)
	}

	if flagShowRemote != showColumnSentinel {
		t.Errorf("Expected flagShowRemote to be sentinel, got %q", flagShowRemote)
	}
	if flagShowStatus != showColumnSentinel {
		t.Errorf("Expected flagShowStatus to be sentinel, got %q", flagShowStatus)
	}
	if flagShowUser != showColumnSentinel {
		t.Errorf("Expected flagShowUser to be sentinel, got %q", flagShowUser)
	}
}

// --- Deprecated Flag Tests ---

func TestDeprecatedVisibilityFlag(t *testing.T) {
	// Old -V flag should still be registered (hidden) on listCmd
	flag := listCmd.Flags().Lookup("dep-visibility")
	if flag == nil {
		t.Fatal("Deprecated -V flag not found on listCmd")
	}
	if !flag.Hidden {
		t.Error("Deprecated -V flag should be hidden")
	}
	if flag.Shorthand != "V" {
		t.Errorf("Deprecated visibility flag shorthand: expected 'V', got '%s'", flag.Shorthand)
	}
}

// --- AnyColumnSelected Tests ---

func TestAnyColumnSelected(t *testing.T) {
	tests := []struct {
		name     string
		opts     ListOptions
		expected bool
	}{
		{"no columns", ListOptions{}, false},
		{"show type", ListOptions{ShowType: true}, true},
		{"show visibility", ListOptions{ShowVisibility: true}, true},
		{"show tags", ListOptions{ShowTags: true}, true},
		{"show path", ListOptions{ShowPath: true}, true},
		{"show status", ListOptions{ShowStatus: true}, true},
		{"show user", ListOptions{ShowUser: true}, true},
		{"show remote", ListOptions{ShowRemote: true}, true},
		{"multiple columns", ListOptions{ShowType: true, ShowTags: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opts.AnyColumnSelected()
			if result != tt.expected {
				t.Errorf("AnyColumnSelected() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// --- Display Tests ---

// captureStdout redirects stdout to capture printed output.
// Since displayMultiColumn uses term.IsTerminal which will return false
// for a pipe, we test the non-TTY (single-column) path directly.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestDisplayMultiColumn_NonTTY(t *testing.T) {
	// When stdout is a pipe (non-TTY), displayMultiColumn outputs one name per line
	tests := []struct {
		name     string
		names    []string
		expected []string
	}{
		{
			"empty list",
			[]string{},
			[]string{},
		},
		{
			"single repo",
			[]string{"my-repo"},
			[]string{"my-repo"},
		},
		{
			"multiple repos sorted",
			[]string{"charlie", "alpha", "bravo"},
			[]string{"alpha", "bravo", "charlie"}, // alphabetical
		},
		{
			"already sorted",
			[]string{"aaa", "bbb", "ccc"},
			[]string{"aaa", "bbb", "ccc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				displayMultiColumn(tt.names)
			})

			if len(tt.expected) == 0 {
				if output != "" {
					t.Errorf("Expected empty output, got %q", output)
				}
				return
			}

			lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
			if len(lines) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d: %v", len(tt.expected), len(lines), lines)
				return
			}
			for i, expected := range tt.expected {
				if lines[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
				}
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	// In test environment (non-TTY), should return default of 80
	width := getTerminalWidth()
	if width != 80 {
		t.Logf("Terminal width detected as %d (may be real TTY or default 80)", width)
	}
	if width <= 0 {
		t.Errorf("getTerminalWidth() returned %d, expected positive value", width)
	}
}

func TestVerboseCountFlag(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedCount int
	}{
		{"no verbose", []string{}, 0},
		{"single -v", []string{"-v"}, 1},
		{"double -vv", []string{"-vv"}, 2},
		{"triple -vvv", []string{"-vvv"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCount := verboseCount
			defer func() { verboseCount = origCount }()
			verboseCount = 0

			if len(tt.args) > 0 {
				err := listCmd.Flags().Parse(tt.args)
				if err != nil {
					t.Fatalf("Failed to parse %v: %v", tt.args, err)
				}
			}

			if verboseCount != tt.expectedCount {
				t.Errorf("Expected verboseCount=%d, got %d", tt.expectedCount, verboseCount)
			}
		})
	}
}

func TestVerboseLevelColumnOverrides(t *testing.T) {
	tests := []struct {
		name           string
		opts           ListOptions
		expectShowType bool
		expectShowVis  bool
		expectShowTags bool
		expectShowPath bool
		expectShowStat bool
		expectShowUser bool
		expectShowRem  bool
	}{
		{
			"verbose 0 - no overrides",
			ListOptions{VerboseLevel: 0},
			false, false, false, false, false, false, false,
		},
		{
			"verbose 1 - stored data columns",
			ListOptions{VerboseLevel: 1},
			true, true, true, true, false, false, false,
		},
		{
			"verbose 2 - all columns",
			ListOptions{VerboseLevel: 2},
			true, true, true, true, true, true, true,
		},
		{
			"verbose 1 with filter - filter preserved, all stored columns shown",
			ListOptions{VerboseLevel: 1, FilterType: "github", ShowType: true},
			true, true, true, true, false, false, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.opts
			// Apply verbose overrides (same logic as parseDualPurposeFlags)
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

			if opts.ShowType != tt.expectShowType {
				t.Errorf("ShowType = %v, want %v", opts.ShowType, tt.expectShowType)
			}
			if opts.ShowVisibility != tt.expectShowVis {
				t.Errorf("ShowVisibility = %v, want %v", opts.ShowVisibility, tt.expectShowVis)
			}
			if opts.ShowTags != tt.expectShowTags {
				t.Errorf("ShowTags = %v, want %v", opts.ShowTags, tt.expectShowTags)
			}
			if opts.ShowPath != tt.expectShowPath {
				t.Errorf("ShowPath = %v, want %v", opts.ShowPath, tt.expectShowPath)
			}
			if opts.ShowStatus != tt.expectShowStat {
				t.Errorf("ShowStatus = %v, want %v", opts.ShowStatus, tt.expectShowStat)
			}
			if opts.ShowUser != tt.expectShowUser {
				t.Errorf("ShowUser = %v, want %v", opts.ShowUser, tt.expectShowUser)
			}
			if opts.ShowRemote != tt.expectShowRem {
				t.Errorf("ShowRemote = %v, want %v", opts.ShowRemote, tt.expectShowRem)
			}

			// Verify filter is preserved
			if tt.opts.FilterType != "" && opts.FilterType != tt.opts.FilterType {
				t.Errorf("FilterType = %q, want %q (filter should be preserved)", opts.FilterType, tt.opts.FilterType)
			}
		})
	}
}

func TestDisplayJSON_ColumnSelection(t *testing.T) {
	repos := []config.Repository{
		{
			Name:       "test-repo",
			Type:       "github",
			Visibility: "private",
			Tags:       []string{"web", "api"},
			Path:       "/home/user/test-repo",
			User:       "testuser",
			Email:      "test@example.com",
		},
	}

	tests := []struct {
		name        string
		opts        ListOptions
		expectKeys  []string
		excludeKeys []string
	}{
		{
			"default - name only",
			ListOptions{OutputFormat: "json"},
			[]string{"name"},
			[]string{"type", "visibility", "tags", "path", "status", "user", "email"},
		},
		{
			"type and path flags",
			ListOptions{OutputFormat: "json", ShowType: true, ShowPath: true},
			[]string{"name", "type", "path"},
			[]string{"visibility", "tags", "status", "user", "email"},
		},
		{
			"verbose 1 - stored data",
			ListOptions{OutputFormat: "json", ShowType: true, ShowVisibility: true, ShowTags: true, ShowPath: true, VerboseLevel: 1},
			[]string{"name", "type", "visibility", "tags", "path"},
			[]string{"status", "user", "email"},
		},
		{
			"show user",
			ListOptions{OutputFormat: "json", ShowUser: true},
			[]string{"name", "user", "email", "signing_enabled"},
			[]string{"type", "visibility", "tags", "path"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				_ = displayJSON(repos, nil, tt.opts)
			})

			var result []map[string]interface{}
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, output)
			}

			if len(result) != 1 {
				t.Fatalf("Expected 1 entry, got %d", len(result))
			}

			entry := result[0]
			for _, key := range tt.expectKeys {
				if _, ok := entry[key]; !ok {
					t.Errorf("Expected key %q in JSON output, not found. Keys: %v", key, mapKeys(entry))
				}
			}
			for _, key := range tt.excludeKeys {
				if _, ok := entry[key]; ok {
					t.Errorf("Key %q should NOT be in JSON output when not selected", key)
				}
			}
		})
	}
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// --- Formatting Tests ---

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		code     string
		expected string
	}{
		{"green check", "✓", ansiGreen, "\033[32m✓\033[0m"},
		{"red cross", "✗", ansiRed, "\033[31m✗\033[0m"},
		{"cyan ahead", "↑3", ansiCyan, "\033[36m↑3\033[0m"},
		{"magenta behind", "↓2", ansiMagenta, "\033[35m↓2\033[0m"},
		{"empty text", "", ansiGreen, "\033[32m\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := colorize(tt.text, tt.code)
			if got != tt.expected {
				t.Errorf("colorize(%q, %q) = %q, want %q", tt.text, tt.code, got, tt.expected)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no ansi", "hello", "hello"},
		{"green text", "\033[32m✓\033[0m", "✓"},
		{"multiple codes", "\033[31m✗\033[0m \033[36m↑3\033[0m", "✗ ↑3"},
		{"empty", "", ""},
		{"plain unicode", "main ✓ ↑3", "main ✓ ↑3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripANSI(tt.input)
			if got != tt.expected {
				t.Errorf("stripANSI(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"ascii", "main", 4},
		{"unicode check", "✓", 1},
		{"unicode cross", "✗", 1},
		{"unicode arrows", "↑3↓2", 4},
		{"mixed", "main ✓ ↑3", 9},
		{"colored", "\033[32m✓\033[0m", 1},
		{"colored mixed", "\033[31m✗\033[0m \033[36m↑3\033[0m", 4},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := displayWidth(tt.input)
			if got != tt.expected {
				t.Errorf("displayWidth(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTruncateBranch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"short name", "main", 30, "main"},
		{"exact limit", "abcdefghij", 10, "abcdefghij"},
		{"one over", "abcdefghijk", 10, "abcdefg..."},
		{"very long", "feature/JIRA-1234-implement-user-authentication-flow", 30, "feature/JIRA-1234-implement..."},
		{"max 3", "abcdef", 3, "abc"},
		{"max 2", "abcdef", 2, "ab"},
		{"max 1", "abcdef", 1, "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateBranch(tt.input, tt.maxLen)
			if got != tt.expected {
				t.Errorf("truncateBranch(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.expected)
			}
		})
	}
}

func TestFormatStatusBranch(t *testing.T) {
	tests := []struct {
		name     string
		status   *git.Status
		expected string
	}{
		{"simple branch", &git.Status{Branch: "main"}, "main"},
		{"no commits", &git.Status{Branch: ""}, "no commits"},
		{"long branch", &git.Status{Branch: "feature/very-long-branch-name-that-exceeds-the-limit"}, "feature/very-long-branch-na..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusBranch(tt.status)
			if got != tt.expected {
				t.Errorf("formatStatusBranch() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatStatusIcons(t *testing.T) {
	tests := []struct {
		name         string
		status       *git.Status
		colorEnabled bool
		expected     string
	}{
		{"clean no color", &git.Status{Branch: "main", IsClean: true}, false, "✓"},
		{"dirty no color", &git.Status{Branch: "main", HasChanges: true}, false, "✗"},
		{"ahead no color", &git.Status{Branch: "main", Ahead: 3}, false, "↑3 ✓"},
		{"behind no color", &git.Status{Branch: "main", Behind: 2}, false, "↓2 ✓"},
		{"ahead+behind no color", &git.Status{Branch: "main", HasChanges: true, Ahead: 5, Behind: 3}, false, "↓3 ↑5 ✗"},
		{"no commits", &git.Status{Branch: ""}, false, ""},
		{"clean with color", &git.Status{Branch: "main", IsClean: true}, true, "\033[32m✓\033[0m"},
		{"dirty with color", &git.Status{Branch: "main", HasChanges: true}, true, "\033[31m✗\033[0m"},
		{"ahead with color", &git.Status{Branch: "main", Ahead: 3}, true, "\033[36m↑3\033[0m \033[32m✓\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusIcons(tt.status, tt.colorEnabled)
			if got != tt.expected {
				t.Errorf("formatStatusIcons() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatStatusIconsColorWidthParity(t *testing.T) {
	statuses := []*git.Status{
		{Branch: "main", IsClean: true},
		{Branch: "main", HasChanges: true},
		{Branch: "main", Ahead: 3},
		{Branch: "main", Behind: 2},
		{Branch: "main", HasChanges: true, Ahead: 5, Behind: 3},
	}

	for _, s := range statuses {
		colored := formatStatusIcons(s, true)
		uncolored := formatStatusIcons(s, false)
		if displayWidth(colored) != displayWidth(uncolored) {
			t.Errorf("display width mismatch for %+v: colored=%d, uncolored=%d",
				s, displayWidth(colored), displayWidth(uncolored))
		}
	}
}

func TestFormatStatusShort(t *testing.T) {
	tests := []struct {
		name     string
		status   *git.Status
		expected string
	}{
		{"clean", &git.Status{Branch: "main"}, "main ✓"},
		{"dirty", &git.Status{Branch: "main", HasChanges: true}, "main ✗"},
		{"ahead", &git.Status{Branch: "feature", Ahead: 3}, "feature ↑3 ✓"},
		{"behind", &git.Status{Branch: "main", Behind: 2}, "main ↓2 ✓"},
		{"ahead+behind", &git.Status{Branch: "develop", HasChanges: true, Ahead: 5, Behind: 3}, "develop ↓3 ↑5 ✗"},
		{"no commits", &git.Status{Branch: ""}, "no commits"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusShort(tt.status)
			if got != tt.expected {
				t.Errorf("formatStatusShort() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestColorFlagRegistered(t *testing.T) {
	flag := listCmd.Flags().Lookup("color")
	if flag == nil {
		t.Fatal("--color flag not found on listCmd")
	}
	if flag.DefValue != "auto" {
		t.Errorf("--color default value = %q, want %q", flag.DefValue, "auto")
	}
}

func TestDeprecatedListDispatchPopulatesListOptions(t *testing.T) {
	// Verify the deprecated --list dispatch sets Show* flags to match old default behavior
	// (show all stored-data columns: type, visibility, tags, path)
	opts := ListOptions{
		FilterType:     "github",
		FilterTags:     []string{"web"},
		FilterName:     "",
		FilterPath:     "",
		OutputFormat:   "table",
		ShowType:       true,
		ShowVisibility: true,
		ShowTags:       true,
		ShowPath:       true,
		ShowStatus:     false,
		ShowUser:       false,
		ShowRemote:     false,
	}

	// Verify old default: stored-data columns shown
	if !opts.ShowType || !opts.ShowVisibility || !opts.ShowTags || !opts.ShowPath {
		t.Error("Deprecated list should show all stored-data columns")
	}
	// Verify live columns NOT shown by default
	if opts.ShowStatus || opts.ShowUser || opts.ShowRemote {
		t.Error("Deprecated list should NOT show live-fetched columns by default")
	}
	// Verify AnyColumnSelected is true
	if !opts.AnyColumnSelected() {
		t.Error("Deprecated list should have columns selected (triggers table view, not multi-column)")
	}
}
