package config

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/daileyo/gws/internal/xdg"
)

// withMigrationEnv isolates a test from the real home directory and resets the
// once-per-process migration guard so each case runs the migration itself.
// It returns the captured notice buffer.
func withMigrationEnv(t *testing.T) *bytes.Buffer {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv(xdg.EnvConfigHome, "")
	t.Setenv(xdg.EnvDataHome, "")

	migrateOnce = sync.Once{}
	t.Cleanup(func() { migrateOnce = sync.Once{} })

	var notice bytes.Buffer
	orig := migrationNotice
	migrationNotice = &notice
	t.Cleanup(func() { migrationNotice = orig })

	return &notice
}

// seedLegacyConfig writes a config at the pre-XDG location.
func seedLegacyConfig(t *testing.T, workspace string) string {
	t.Helper()

	legacyPath, err := xdg.LegacyConfigFile()
	if err != nil {
		t.Fatalf("failed to resolve legacy path: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0755); err != nil {
		t.Fatalf("failed to create legacy dir: %v", err)
	}

	body := `{"version":"1.1.0","workspace":"` + workspace + `","repositories":[]}`
	if err := os.WriteFile(legacyPath, []byte(body), 0600); err != nil {
		t.Fatalf("failed to seed legacy config: %v", err)
	}

	return legacyPath
}

func TestMigrate_MovesLegacyConfig(t *testing.T) {
	notice := withMigrationEnv(t)
	legacyPath := seedLegacyConfig(t, "/workspace/root")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Workspace != "/workspace/root" {
		t.Errorf("workspace = %q, want %q", cfg.Workspace, "/workspace/root")
	}

	newPath, _ := xdg.ConfigFile()
	if _, err := os.Stat(newPath); err != nil {
		t.Errorf("config was not moved to %s: %v", newPath, err)
	}
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Errorf("legacy config still present at %s", legacyPath)
	}
	if !strings.Contains(notice.String(), newPath) {
		t.Errorf("notice should name the destination, got %q", notice.String())
	}
}

func TestMigrate_PreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permission bits are not meaningful on Windows")
	}

	withMigrationEnv(t)
	seedLegacyConfig(t, "/workspace/root")

	if _, err := Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newPath, _ := xdg.ConfigFile()
	info, err := os.Stat(newPath)
	if err != nil {
		t.Fatalf("migrated config missing: %v", err)
	}
	// The config records repository paths and git identity; widening it would
	// expose that to other local users.
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("migrated config mode = %04o, want 0600", perm)
	}
}

func TestMigrate_RemovesEmptyLegacyDir(t *testing.T) {
	notice := withMigrationEnv(t)
	legacyPath := seedLegacyConfig(t, "/workspace/root")

	if _, err := Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	legacyDir := filepath.Dir(legacyPath)
	if _, err := os.Stat(legacyDir); !os.IsNotExist(err) {
		t.Errorf("empty legacy dir %s should have been removed", legacyDir)
	}
	if !strings.Contains(notice.String(), "removed empty") {
		t.Errorf("notice should mention the cleanup, got %q", notice.String())
	}
}

func TestMigrate_KeepsNonEmptyLegacyDir(t *testing.T) {
	notice := withMigrationEnv(t)
	legacyPath := seedLegacyConfig(t, "/workspace/root")

	// Something the user put there themselves must survive.
	legacyDir := filepath.Dir(legacyPath)
	keep := filepath.Join(legacyDir, "notes.txt")
	if err := os.WriteFile(keep, []byte("mine"), 0600); err != nil {
		t.Fatalf("failed to write extra file: %v", err)
	}

	if _, err := Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(keep); err != nil {
		t.Errorf("user file in the legacy dir was removed: %v", err)
	}
	if !strings.Contains(notice.String(), "notes.txt") {
		t.Errorf("notice should name what was left behind, got %q", notice.String())
	}
}

func TestMigrate_NoOpWhenNewConfigExists(t *testing.T) {
	notice := withMigrationEnv(t)
	seedLegacyConfig(t, "/legacy/workspace")

	// A config already at the new location must win, untouched.
	newPath, _ := xdg.ConfigFile()
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	body := `{"version":"1.2.0","workspace":"/current/workspace","repositories":[]}`
	if err := os.WriteFile(newPath, []byte(body), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Workspace != "/current/workspace" {
		t.Errorf("workspace = %q, want the existing config to win", cfg.Workspace)
	}
	if notice.Len() != 0 {
		t.Errorf("no migration should be reported, got %q", notice.String())
	}
}

func TestMigrate_RefusesSymlinkedLegacyConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation needs elevation on Windows")
	}

	notice := withMigrationEnv(t)

	// A symlink at the legacy path could redirect the move elsewhere.
	target := filepath.Join(t.TempDir(), "elsewhere.json")
	if err := os.WriteFile(target, []byte(`{"version":"1.1.0","workspace":"/x","repositories":[]}`), 0600); err != nil {
		t.Fatalf("failed to write target: %v", err)
	}
	legacyPath, _ := xdg.LegacyConfigFile()
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0755); err != nil {
		t.Fatalf("failed to create legacy dir: %v", err)
	}
	if err := os.Symlink(target, legacyPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// The workspace still loads — via the legacy fallback — but is not moved.
	if _, err := Load(); err != nil {
		t.Fatalf("a symlinked legacy config should still load: %v", err)
	}

	newPath, _ := xdg.ConfigFile()
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Errorf("symlinked config should not have been migrated to %s", newPath)
	}
	if !strings.Contains(notice.String(), "symlink") {
		t.Errorf("notice should explain the refusal, got %q", notice.String())
	}
}

func TestMigrate_FallsBackToLegacyOnFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions do not block writes the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}

	withMigrationEnv(t)
	seedLegacyConfig(t, "/workspace/root")

	// Make the destination unwritable so the move cannot succeed.
	configDir, _ := xdg.ConfigDir()
	parent := filepath.Dir(configDir)
	if err := os.MkdirAll(parent, 0755); err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}
	if err := os.Chmod(parent, 0500); err != nil {
		t.Fatalf("failed to chmod parent: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0755) })

	// A failed migration must not present a working workspace as uninitialized.
	cfg, err := Load()
	if err != nil {
		t.Fatalf("load should fall back to the legacy path, got: %v", err)
	}
	if cfg.Workspace != "/workspace/root" {
		t.Errorf("workspace = %q, want the legacy config to be read", cfg.Workspace)
	}
}

func TestMigrate_UninitializedWhenNeitherExists(t *testing.T) {
	withMigrationEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("expected an error when no config exists anywhere")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("error = %q, want the not-initialized guidance", err.Error())
	}
}

func TestExists_FindsLegacyConfig(t *testing.T) {
	withMigrationEnv(t)
	seedLegacyConfig(t, "/workspace/root")

	// Exists runs before any migration, so it has to look in both places or a
	// pending-migration workspace reports as uninitialized.
	exists, err := Exists()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("Exists() should find a config still at the legacy path")
	}
}

func TestExists_FalseWhenNeitherExists(t *testing.T) {
	withMigrationEnv(t)

	exists, err := Exists()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("Exists() should be false when no config exists anywhere")
	}
}
