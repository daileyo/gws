package xdg

import (
	"errors"
	"path/filepath"
	"testing"
)

// withHome points resolution at a temp directory and clears the XDG variables,
// so no test reads or writes the real home directory.
func withHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	orig := homeDirFunc
	homeDirFunc = func() (string, error) { return home, nil }
	t.Cleanup(func() { homeDirFunc = orig })
	t.Setenv(EnvConfigHome, "")
	t.Setenv(EnvDataHome, "")
	return home
}

func TestConfigDir(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want func(home string) string
	}{
		{
			name: "unset falls back to home",
			env:  "",
			want: func(home string) string { return filepath.Join(home, ".config", "gws") },
		},
		{
			name: "absolute value is honored",
			env:  filepath.FromSlash("/custom/config"),
			want: func(string) string { return filepath.FromSlash("/custom/config/gws") },
		},
		{
			name: "relative value is treated as unset",
			env:  filepath.Join("relative", "config"),
			want: func(home string) string { return filepath.Join(home, ".config", "gws") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := withHome(t)
			t.Setenv(EnvConfigHome, tt.env)

			got, err := ConfigDir()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if want := tt.want(home); got != want {
				t.Errorf("ConfigDir() = %q, want %q", got, want)
			}
		})
	}
}

func TestProjectsDir(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want func(home string) string
	}{
		{
			name: "unset falls back to home",
			env:  "",
			want: func(home string) string {
				return filepath.Join(home, ".local", "share", "gws", "projects")
			},
		},
		{
			name: "absolute value is honored",
			env:  filepath.FromSlash("/custom/data"),
			want: func(string) string { return filepath.FromSlash("/custom/data/gws/projects") },
		},
		{
			name: "relative value is treated as unset",
			env:  filepath.Join("relative", "data"),
			want: func(home string) string {
				return filepath.Join(home, ".local", "share", "gws", "projects")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := withHome(t)
			t.Setenv(EnvDataHome, tt.env)

			got, err := ProjectsDir()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if want := tt.want(home); got != want {
				t.Errorf("ProjectsDir() = %q, want %q", got, want)
			}
		})
	}
}

// TestLayoutIsPlatformIndependent is the guard on the spec's cross-platform
// parity requirement: the layout relative to home must not vary by platform.
// A regression here would most likely arrive as a runtime.GOOS branch or a
// switch to os.UserConfigDir, which returns %AppData% on Windows.
func TestLayoutIsPlatformIndependent(t *testing.T) {
	home := withHome(t)

	cases := []struct {
		name string
		fn   func() (string, error)
		want []string
	}{
		{"config", ConfigDir, []string{".config", "gws"}},
		{"data", DataDir, []string{".local", "share", "gws"}},
		{"projects", ProjectsDir, []string{".local", "share", "gws", "projects"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.fn()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			rel, err := filepath.Rel(home, got)
			if err != nil {
				t.Fatalf("path is not under home: %v", err)
			}
			if want := filepath.Join(tc.want...); rel != want {
				t.Errorf("layout relative to home = %q, want %q (must be identical on every platform)", rel, want)
			}
		})
	}
}

func TestConfigFile(t *testing.T) {
	home := withHome(t)

	got, err := ConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := filepath.Join(home, ".config", "gws", "config.json"); got != want {
		t.Errorf("ConfigFile() = %q, want %q", got, want)
	}
}

func TestRepoProjectsDir(t *testing.T) {
	home := withHome(t)

	got, err := RepoProjectsDir("my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(home, ".local", "share", "gws", "projects", "my-repo")
	if got != want {
		t.Errorf("RepoProjectsDir() = %q, want %q", got, want)
	}
}

func TestLegacyConfig(t *testing.T) {
	home := withHome(t)

	// The legacy location must ignore XDG_CONFIG_HOME — it is a fixed historical
	// path, not a configurable one.
	t.Setenv(EnvConfigHome, filepath.FromSlash("/somewhere/else"))

	dir, err := LegacyConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := filepath.Join(home, ".gws"); dir != want {
		t.Errorf("LegacyConfigDir() = %q, want %q", dir, want)
	}

	file, err := LegacyConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := filepath.Join(home, ".gws", "config.json"); file != want {
		t.Errorf("LegacyConfigFile() = %q, want %q", file, want)
	}
}

func TestHomeDirFailurePropagates(t *testing.T) {
	orig := homeDirFunc
	sentinel := errors.New("no home")
	homeDirFunc = func() (string, error) { return "", sentinel }
	t.Cleanup(func() { homeDirFunc = orig })
	t.Setenv(EnvConfigHome, "")

	if _, err := ConfigDir(); !errors.Is(err, sentinel) {
		t.Errorf("ConfigDir() error = %v, want it to wrap %v", err, sentinel)
	}
}
