package classifier

import (
	"strings"

	"github.com/daileyo/gws/internal/config"
)

// DetectType determines the repository type based on the remote URL
func DetectType(remoteURL string) config.RepositoryType {
	if remoteURL == "" {
		return config.TypeUnknown
	}

	urlLower := strings.ToLower(remoteURL)

	// GitHub
	if strings.Contains(urlLower, "github.com") {
		return config.TypeGitHub
	}

	// GitLab (both gitlab.com and self-hosted)
	if strings.Contains(urlLower, "gitlab.com") || strings.Contains(urlLower, "gitlab") {
		return config.TypeGitLab
	}

	// Azure DevOps
	if strings.Contains(urlLower, "dev.azure.com") || strings.Contains(urlLower, "visualstudio.com") {
		return config.TypeADO
	}

	// Bitbucket
	if strings.Contains(urlLower, "bitbucket.org") {
		return config.TypeBitbucket
	}

	return config.TypeUnknown
}

// DetectVisibility determines repository visibility based on the remote URL protocol
// SSH URLs (git@...) are typically used for private repos
// HTTPS URLs can be either public or private, so we mark them as unknown
func DetectVisibility(remoteURL string) config.RepositoryVisibility {
	if remoteURL == "" {
		return config.VisibilityUnknown
	}

	// SSH URLs typically indicate private repositories (requires SSH key authentication)
	if strings.HasPrefix(remoteURL, "git@") || strings.HasPrefix(remoteURL, "ssh://") {
		return config.VisibilityPrivate
	}

	// HTTPS URLs could be public or private, we can't determine from URL alone
	if strings.HasPrefix(remoteURL, "https://") || strings.HasPrefix(remoteURL, "http://") {
		return config.VisibilityUnknown
	}

	return config.VisibilityUnknown
}

// Classify performs both type and visibility detection on a repository
func Classify(repo *config.Repository) {
	repo.Type = DetectType(repo.RemoteURL)
	repo.Visibility = DetectVisibility(repo.RemoteURL)
}
