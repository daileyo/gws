package filter

import (
	"path/filepath"
	"strings"

	"github.com/daileyo/gws/internal/config"
)

// Criteria represents filtering criteria for repositories
type Criteria struct {
	Type string   // Repository type filter (github, gitlab, ado, bitbucket)
	Tags []string // Custom tags filter (AND logic - repository must have all tags)
	Name string   // Repository name filter (partial, case-insensitive)
	Path string   // Repository path filter (partial match)
}

// Apply filters a list of repositories based on the provided criteria
func Apply(repos []config.Repository, criteria Criteria) []config.Repository {
	filtered := []config.Repository{}

	for _, repo := range repos {
		if matchesCriteria(repo, criteria) {
			filtered = append(filtered, repo)
		}
	}

	return filtered
}

// MatchesPattern checks if a value matches a pattern. If the pattern contains
// wildcard characters (* or ?), it uses glob-style matching. Otherwise, it
// falls back to partial, case-insensitive substring matching.
func MatchesPattern(value, pattern string) bool {
	valueLower := strings.ToLower(value)
	patternLower := strings.ToLower(pattern)

	if strings.ContainsAny(pattern, "*?") {
		matched, err := filepath.Match(patternLower, valueLower)
		if err != nil {
			return false
		}
		return matched
	}

	// Non-wildcard: partial, case-insensitive match
	return strings.Contains(valueLower, patternLower)
}

// matchesCriteria checks if a repository matches all filter criteria
func matchesCriteria(repo config.Repository, criteria Criteria) bool {
	// Apply type filter
	if criteria.Type != "" && !MatchesPattern(string(repo.Type), criteria.Type) {
		return false
	}

	// Apply tags filter (AND logic - must have all specified tags)
	if len(criteria.Tags) > 0 {
		for _, filterTag := range criteria.Tags {
			hasTag := false
			for _, repoTag := range repo.Tags {
				if MatchesPattern(repoTag, filterTag) {
					hasTag = true
					break
				}
			}
			if !hasTag {
				return false
			}
		}
	}

	// Apply name filter
	if criteria.Name != "" {
		if !MatchesPattern(repo.Name, criteria.Name) {
			return false
		}
	}

	// Apply path filter
	if criteria.Path != "" {
		if !MatchesPattern(repo.Path, criteria.Path) {
			return false
		}
	}

	return true
}

// ByType filters repositories by type
func ByType(repos []config.Repository, repoType string) []config.Repository {
	return Apply(repos, Criteria{Type: repoType})
}

// ByTags filters repositories that have all specified tags (AND logic)
func ByTags(repos []config.Repository, tags []string) []config.Repository {
	return Apply(repos, Criteria{Tags: tags})
}

// ByName filters repositories by name (partial match)
func ByName(repos []config.Repository, name string) []config.Repository {
	return Apply(repos, Criteria{Name: name})
}

// ByPath filters repositories by path (partial match)
func ByPath(repos []config.Repository, path string) []config.Repository {
	return Apply(repos, Criteria{Path: path})
}
