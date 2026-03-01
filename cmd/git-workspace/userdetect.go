package main

import (
	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
)

// detectUserForRepos iterates over repos, calls git.GetUserConfig() for each,
// and populates User/Email/SigningEnabled/UserSource fields. Returns the count
// of repos with detected user info.
func detectUserForRepos(repos []config.Repository) int {
	userDetectedCount := 0
	for i := range repos {
		repo := &repos[i]
		userCfg, err := git.GetUserConfig(repo.Path)
		if err == nil && userCfg != nil {
			repo.User = userCfg.Name
			repo.Email = userCfg.Email
			repo.SigningEnabled = userCfg.SignCommits
			repo.UserSource = userCfg.Source
			if userCfg.Name != "" || userCfg.Email != "" {
				userDetectedCount++
			}
		}
	}
	return userDetectedCount
}
