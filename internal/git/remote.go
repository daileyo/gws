package git

import (
	"github.com/go-git/go-git/v5"
)

// RemoteInfo contains information about a repository's remotes
type RemoteInfo struct {
	OriginURL   string // URL of the origin remote (empty if no origin)
	HasMultiple bool   // True if non-origin remotes exist
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
