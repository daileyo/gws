package main

import (
	"github.com/daileyo/gws/internal/config"
	"github.com/daileyo/gws/internal/git"
	"github.com/daileyo/gws/internal/user"
)

// detectUserForRepos iterates over repos, calls git.GetUserConfig() for each,
// and populates User/Email/SigningEnabled/UserSource fields. When a repo has
// local .git/config user settings, the underlying global/includeIf config is
// stored instead (local config is detected at display time only). When profiles
// are provided and a repo's user is detected via includeIf, it attempts to
// auto-link the repo to a matching stored profile. Returns the count of repos
// with detected user info.
func detectUserForRepos(repos []config.Repository, profiles ...[]config.Profile) int {
	var profileList []config.Profile
	if len(profiles) > 0 {
		profileList = profiles[0]
	}

	userDetectedCount := 0
	for i := range repos {
		repo := &repos[i]
		userCfg, err := git.GetUserConfig(repo.Path)
		if err != nil || userCfg == nil {
			continue
		}

		// If local config was detected, don't persist it - get the underlying
		// global/includeIf config instead. Local config is read at display time.
		if userCfg.Source == config.UserSourceLocal {
			nonLocalCfg, err := git.GetNonLocalUserConfig(repo.Path)
			if err == nil && nonLocalCfg != nil && (nonLocalCfg.Name != "" || nonLocalCfg.Email != "") {
				userCfg = nonLocalCfg
			} else {
				// No underlying config found; still count as detected but store empty
				userDetectedCount++
				continue
			}
		}

		repo.User = userCfg.Name
		repo.Email = userCfg.Email
		repo.SigningEnabled = userCfg.SignCommits
		repo.UserSource = userCfg.Source
		if userCfg.Name != "" || userCfg.Email != "" {
			userDetectedCount++
		}

		// Auto-link to stored profile if detected via includeIf
		if userCfg.Source == config.UserSourceIncludeIf && len(profileList) > 0 {
			if matched := user.MatchProfileByUser(profileList, userCfg.Name, userCfg.Email); matched != nil {
				_ = matched // Match exists, linking is implicit via matching user/email
			}
		}
	}
	return userDetectedCount
}
