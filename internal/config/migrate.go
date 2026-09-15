package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/daileyo/omgitworks/internal/xdg"
)

// migrateOnce ensures the legacy config is migrated at most once per process.
var migrateOnce sync.Once

// migrationNotice is where the migration reports what it did. It is a variable
// so tests can capture the notice instead of writing to the real stderr.
var migrationNotice io.Writer = os.Stderr

// ensureMigrated runs the legacy migration a single time, reporting what it did.
//
// Failures are deliberately non-fatal: if the move cannot be completed, callers
// fall back to reading the legacy path. A user with a working workspace must
// never be told it is uninitialized because a file could not be moved.
func ensureMigrated() {
	migrateOnce.Do(func() {
		// A successful migration reports itself; only the failure needs handling here.
		if _, err := migrateLegacyConfig(); err != nil {
			fmt.Fprintf(migrationNotice, "warning: could not migrate config to its new location: %v\n", err)
		}
	})
}

// migrateLegacyConfig moves ~/.gws/config.json to the XDG config location.
//
// It is a no-op when the new config already exists or the legacy one does not,
// so it is safe to call on every load.
func migrateLegacyConfig() (bool, error) {
	newPath, err := xdg.ConfigFile()
	if err != nil {
		return false, err
	}

	// Already migrated, or a fresh install — nothing to do.
	if _, err := os.Stat(newPath); err == nil {
		return false, nil
	}

	legacyPath, err := xdg.LegacyConfigFile()
	if err != nil {
		return false, err
	}

	// Lstat, not Stat: a symlink here could redirect the move outside the
	// intended directory, so refuse rather than follow it.
	info, err := os.Lstat(legacyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("legacy config %s is a symlink; move it manually", legacyPath)
	}
	if !info.Mode().IsRegular() {
		return false, fmt.Errorf("legacy config %s is not a regular file", legacyPath)
	}

	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return false, fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := moveFile(legacyPath, newPath, info.Mode().Perm()); err != nil {
		return false, err
	}

	fmt.Fprintf(migrationNotice, "note: moved config %s -> %s\n", legacyPath, newPath)

	reportLegacyDirCleanup(filepath.Dir(legacyPath))

	return true, nil
}

// moveFile relocates src to dst, preserving perm. It prefers a rename and falls
// back to copy-then-remove when the two paths live on different filesystems.
func moveFile(src, dst string, perm os.FileMode) error {
	if err := os.Rename(src, dst); err == nil {
		// Rename preserves the mode, but be explicit: the config records
		// repository paths and git identity, so it must not widen to 0644.
		return os.Chmod(dst, perm)
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read legacy config: %w", err)
	}
	// #nosec G703 -- dst is derived from the user's own XDG config location,
	// not from external input; there is no privilege boundary to cross here.
	if err := os.WriteFile(dst, data, perm); err != nil {
		return fmt.Errorf("failed to write config to new location: %w", err)
	}
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("config copied to %s but the original could not be removed: %w", dst, err)
	}
	return nil
}

// reportLegacyDirCleanup removes the legacy directory when the migrated config
// was all it held, and otherwise says what was left behind. Anything the user
// put there themselves is preserved.
func reportLegacyDirCleanup(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	if len(entries) == 0 {
		if err := os.Remove(dir); err == nil {
			fmt.Fprintf(migrationNotice, "note: removed empty %s\n", dir)
		}
		return
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	fmt.Fprintf(migrationNotice, "note: left %s in place; it still contains: %v\n", dir, names)
}
