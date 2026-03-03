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
}

func TestShellTemplatesContainNavigationFallthrough(t *testing.T) {
	// The * case should handle navigation via cd
	if !strings.Contains(zshInitTemplate, `_dest="$({BIN}`) {
		t.Error("zsh template missing navigation fallthrough")
	}
	if !strings.Contains(bashInitTemplate, `dest="$({BIN}`) {
		t.Error("bash template missing navigation fallthrough")
	}
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
}
