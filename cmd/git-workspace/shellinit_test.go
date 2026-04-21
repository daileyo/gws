package main

import (
	"strings"
	"testing"
)

func TestShellTemplatesRouteSubcommands(t *testing.T) {
	// Verify shell templates include all subcommand names in the routing pattern
	subcommands := []string{"list", "init", "add", "refresh", "print-workspace", "tag", "user"}

	templates := []struct {
		name     string
		template string
	}{
		{"zsh", zshInitTemplate},
		{"bash", bashInitTemplate},
		{"powershell", powershellInitTemplate},
	}

	for _, tmpl := range templates {
		t.Run(tmpl.name, func(t *testing.T) {
			for _, cmd := range subcommands {
				// Check that the subcommand name appears in the case pattern
				if !strings.Contains(tmpl.template, cmd+"|") && !strings.Contains(tmpl.template, "|"+cmd) {
					t.Errorf("%s template does not route subcommand %q to binary", tmpl.name, cmd)
				}
			}
		})
	}
}

func TestShellTemplatesContainBinPlaceholder(t *testing.T) {
	if !strings.Contains(zshInitTemplate, "{BIN}") {
		t.Error("zsh template missing {BIN} placeholder")
	}
	if !strings.Contains(bashInitTemplate, "{BIN}") {
		t.Error("bash template missing {BIN} placeholder")
	}
	if !strings.Contains(powershellInitTemplate, "{BIN}") {
		t.Error("powershell template missing {BIN} placeholder")
	}
}

func TestShellTemplatesContainNavigationFallthrough(t *testing.T) {
	// The default/fallthrough case should handle navigation via cd/Set-Location
	if !strings.Contains(zshInitTemplate, `_dest="$({BIN}`) {
		t.Error("zsh template missing navigation fallthrough")
	}
	if !strings.Contains(bashInitTemplate, `dest="$({BIN}`) {
		t.Error("bash template missing navigation fallthrough")
	}
	// PowerShell uses Set-Location and & {BIN} for navigation
	if !strings.Contains(powershellInitTemplate, `& {BIN}`) {
		t.Error("powershell template missing binary invocation")
	}
	if !strings.Contains(powershellInitTemplate, `Set-Location`) {
		t.Error("powershell template missing Set-Location navigation")
	}
}

func TestShellTemplatesContainWorktreeNavigation(t *testing.T) {
	templates := []struct {
		name     string
		template string
	}{
		{"zsh", zshInitTemplate},
		{"bash", bashInitTemplate},
	}
	for _, tmpl := range templates {
		t.Run(tmpl.name, func(t *testing.T) {
			// -wt flag detection for repo-scoped worktree navigation
			if !strings.Contains(tmpl.template, `"-wt"`) {
				t.Errorf("%s template missing -wt flag check", tmpl.name)
			}
			// --worktree flag in repo-scoped navigation command
			if !strings.Contains(tmpl.template, `--worktree "$3"`) {
				t.Errorf("%s template missing --worktree with branch argument", tmpl.name)
			}
			// Bare --worktree (no value) for selection mode
			if !strings.Contains(tmpl.template, `--worktree -q`) {
				t.Errorf("%s template missing bare --worktree invocation", tmpl.name)
			}
			// worktree has its own case block (not in passthrough list)
			if !strings.Contains(tmpl.template, "worktree)") {
				t.Errorf("%s template missing worktree case block", tmpl.name)
			}
			// worktree subcommands (list, align, add) are passed through
			if !strings.Contains(tmpl.template, "list|align|add") {
				t.Errorf("%s template missing worktree subcommand passthrough", tmpl.name)
			}
		})
	}

	// PowerShell-specific worktree checks
	t.Run("powershell", func(t *testing.T) {
		// -wt flag detection
		if !strings.Contains(powershellInitTemplate, `'-wt'`) {
			t.Error("powershell template missing -wt flag check")
		}
		// --worktree flag in navigation command
		if !strings.Contains(powershellInitTemplate, `--worktree $branch`) {
			t.Error("powershell template missing --worktree with branch argument")
		}
		// Bare --worktree (no value) for selection mode
		if !strings.Contains(powershellInitTemplate, `--worktree -q`) {
			t.Error("powershell template missing bare --worktree invocation")
		}
		// worktree case block
		if !strings.Contains(powershellInitTemplate, `'^worktree$'`) {
			t.Error("powershell template missing worktree case block")
		}
		// worktree subcommands (list, align, add) are passed through
		if !strings.Contains(powershellInitTemplate, "'list', 'align', 'add'") {
			t.Error("powershell template missing worktree subcommand passthrough")
		}
	})
}

func TestShellTemplatesContainParentNavigation(t *testing.T) {
	templates := []struct {
		name     string
		template string
	}{
		{"zsh", zshInitTemplate},
		{"bash", bashInitTemplate},
	}
	for _, tmpl := range templates {
		t.Run(tmpl.name, func(t *testing.T) {
			// All three invocation styles must route through the parent subcommand
			if !strings.Contains(tmpl.template, `"-p"`) {
				t.Errorf("%s template missing -p flag check", tmpl.name)
			}
			if !strings.Contains(tmpl.template, `"--parent"`) {
				t.Errorf("%s template missing --parent flag check", tmpl.name)
			}
			if !strings.Contains(tmpl.template, "parent)") {
				t.Errorf("%s template missing 'parent' keyword case", tmpl.name)
			}
			if !strings.Contains(tmpl.template, `{BIN} parent`) {
				t.Errorf("%s template missing '{BIN} parent' invocation", tmpl.name)
			}
		})
	}

	// PowerShell-specific parent navigation checks
	t.Run("powershell", func(t *testing.T) {
		if !strings.Contains(powershellInitTemplate, `'-p'`) {
			t.Error("powershell template missing -p flag check")
		}
		if !strings.Contains(powershellInitTemplate, `'--parent'`) {
			t.Error("powershell template missing --parent flag check")
		}
		if !strings.Contains(powershellInitTemplate, `parent`) {
			t.Error("powershell template missing 'parent' keyword case")
		}
		if !strings.Contains(powershellInitTemplate, `{BIN} parent`) {
			t.Error("powershell template missing '{BIN} parent' invocation")
		}
	})
}
