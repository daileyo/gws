package git

import (
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
)

// RemoteInfo contains information about a repository's remotes
type RemoteInfo struct {
	OriginURL   string // URL of the origin remote (empty if no origin)
	HasMultiple bool   // True if non-origin remotes exist
}

// FormatRemoteURL converts a raw git remote URL into a clean HTTPS URL.
// SSH URLs (git@host:path) are converted to https://host/path.
// HTTPS URLs have user info stripped. Azure DevOps SSH URLs are converted
// to their HTTPS equivalent. Unrecognized formats are returned unchanged.
func FormatRemoteURL(rawURL string) string {
	if rawURL == "" {
		return rawURL
	}

	// Handle SSH format: git@host:path
	if strings.HasPrefix(rawURL, "git@") {
		return formatSSHURL(rawURL)
	}

	// Handle scheme-based URLs (https://, http://, ssh://, file://, etc.)
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" {
		// No scheme and not git@ — could be a custom SSH alias (e.g., myhost:repo.git)
		return rawURL
	}

	switch parsed.Scheme {
	case "https", "http":
		parsed.User = nil
		return parsed.String()
	default:
		// file://, ssh://, or other schemes — return unchanged
		return rawURL
	}
}

// formatSSHURL converts git@host:path to https://host/path.
// Azure DevOps SSH (git@ssh.dev.azure.com:v3/org/project/repo) is converted
// to https://dev.azure.com/org/project/_git/repo.
func formatSSHURL(rawURL string) string {
	// Strip "git@" prefix
	rest := rawURL[len("git@"):]

	// Split on first ":"
	colonIdx := strings.Index(rest, ":")
	if colonIdx < 0 {
		return rawURL
	}

	host := rest[:colonIdx]
	path := rest[colonIdx+1:]

	// Azure DevOps SSH: git@ssh.dev.azure.com:v3/org/project/repo
	if host == "ssh.dev.azure.com" && strings.HasPrefix(path, "v3/") {
		parts := strings.SplitN(path, "/", 4) // ["v3", "org", "project", "repo"]
		if len(parts) == 4 {
			return "https://dev.azure.com/" + parts[1] + "/" + parts[2] + "/_git/" + parts[3]
		}
	}

	return "https://" + host + "/" + path
}

// GetRemoteInfo inspects a repository's remotes and returns origin URL
// and whether non-origin remotes exist.
func GetRemoteInfo(repoPath string) (*RemoteInfo, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}

	info := &RemoteInfo{}
	nonOriginCount := 0

	for _, remote := range remotes {
		name := remote.Config().Name
		if name == "origin" {
			urls := remote.Config().URLs
			if len(urls) > 0 {
				info.OriginURL = urls[0]
			}
		} else {
			nonOriginCount++
		}
	}

	info.HasMultiple = nonOriginCount > 0
	return info, nil
}
