package main

import (
	"fmt"
	"strings"

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

// syncProfilesFromRepos scans repos for unique user identities and adds any
// that don't already exist as profiles to cfg.Profiles. This ensures the
// profile list reflects all users that are in use across repos.
func syncProfilesFromRepos(cfg *config.Config) {
	// Build set of existing profile emails for dedup
	existingEmails := make(map[string]bool)
	for _, p := range cfg.Profiles {
		existingEmails[strings.ToLower(p.Email)] = true
	}

	// Also collect from auto-detected profiles (includeIf)
	detected, _ := user.DetectProfiles()
	for _, p := range detected {
		if existingEmails[strings.ToLower(p.Email)] {
			continue
		}
		existingEmails[strings.ToLower(p.Email)] = true
		cfg.Profiles = append(cfg.Profiles, p)
	}

	// Collect unique users from repos (stored config + live local config)
	type userIdentity struct {
		name  string
		email string
	}
	var identities []userIdentity

	for _, repo := range cfg.Repositories {
		// Add the stored (global/includeIf) identity
		if repo.Email != "" {
			identities = append(identities, userIdentity{repo.User, repo.Email})
		}

		// Also check live local config (detectUserForRepos skips persisting these)
		liveCfg, err := git.GetUserConfig(repo.Path)
		if err == nil && liveCfg != nil && liveCfg.Source == config.UserSourceLocal && liveCfg.Email != "" {
			identities = append(identities, userIdentity{liveCfg.Name, liveCfg.Email})
		}
	}

	for _, id := range identities {
		emailLower := strings.ToLower(id.email)
		if existingEmails[emailLower] {
			continue
		}
		existingEmails[emailLower] = true

		// Derive profile name from git name (lowercase, spaces to hyphens)
		name := strings.ToLower(id.name)
		name = strings.ReplaceAll(name, " ", "-")
		if name == "" {
			// Fall back to email prefix
			name = strings.SplitN(id.email, "@", 2)[0]
		}

		// Ensure unique profile name
		baseName := name
		suffix := 1
		for profileNameExists(cfg.Profiles, name) {
			suffix++
			name = fmt.Sprintf("%s-%d", baseName, suffix)
		}

		cfg.Profiles = append(cfg.Profiles, config.Profile{
			Name:    name,
			GitName: id.name,
			Email:   id.email,
		})
	}
}

// profileNameExists checks if a profile with the given name already exists
func profileNameExists(profiles []config.Profile, name string) bool {
	for _, p := range profiles {
		if strings.EqualFold(p.Name, name) {
			return true
		}
	}
	return false
}
