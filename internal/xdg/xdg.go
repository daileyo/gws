// Package xdg resolves where gws keeps its files, following the XDG Base
// Directory Specification.
//
// The layout is deliberately identical on every platform: <home>/.config/gws
// and <home>/.local/share/gws/projects, with <home> being %USERPROFILE% on
// Windows. Windows has no XDG specification, but git itself already reads
// $HOME/.config/git/config there (Git for Windows sets $HOME to %USERPROFILE%),
// so a git-adjacent tool storing its config beside git's own is following the
// convention of its ecosystem rather than inventing one.
//
// This is why os.UserConfigDir is not used: it returns %AppData% on Windows,
// which would give Windows users a different layout from everyone else.
package xdg

import (
	"fmt"
	"os"
	"path/filepath"
)

// appName is the directory gws owns inside each base directory.
const appName = "gws"

// Environment variables from the XDG Base Directory Specification.
const (
	EnvConfigHome = "XDG_CONFIG_HOME"
	EnvDataHome   = "XDG_DATA_HOME"
)

// homeDirFunc resolves the user's home directory. It is a variable so tests can
// drive resolution without reading the real home directory.
var homeDirFunc = os.UserHomeDir

// baseDir returns the base directory for the given XDG environment variable,
// falling back to <home>/<relative> when the variable is unset.
//
// Per the specification, a value that is not an absolute path is treated as
// though it were unset. Honoring a relative value would resolve it against the
// current working directory, letting a stray cd change where gws writes.
func baseDir(env string, relative ...string) (string, error) {
	if v := os.Getenv(env); v != "" && filepath.IsAbs(v) {
		return v, nil
	}

	home, err := homeDirFunc()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}

	return filepath.Join(append([]string{home}, relative...)...), nil
}

// ConfigDir returns the directory holding gws configuration,
// $XDG_CONFIG_HOME/gws or <home>/.config/gws.
func ConfigDir() (string, error) {
	base, err := baseDir(EnvConfigHome, ".config")
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

// ConfigFile returns the path to the gws config file.
func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// DataDir returns the directory holding gws data,
// $XDG_DATA_HOME/gws or <home>/.local/share/gws.
func DataDir() (string, error) {
	base, err := baseDir(EnvDataHome, ".local", "share")
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

// ProjectsDir returns the root under which worktrees are organized, one
// subdirectory per repository. Worktrees are real user data — losing one loses
// work — so they belong under the data directory rather than cache or state.
func ProjectsDir() (string, error) {
	dir, err := DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "projects"), nil
}

// RepoProjectsDir returns the directory holding the worktrees of a single
// repository, keyed by repository name.
func RepoProjectsDir(repoName string) (string, error) {
	dir, err := ProjectsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, repoName), nil
}

// LegacyConfigDir returns the pre-XDG configuration directory, ~/.gws.
// It exists solely so the migration path can find a config left by an older
// version; nothing else should write there.
func LegacyConfigDir() (string, error) {
	home, err := homeDirFunc()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	return filepath.Join(home, ".gws"), nil
}

// LegacyConfigFile returns the path to the pre-XDG config file.
func LegacyConfigFile() (string, error) {
	dir, err := LegacyConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}
