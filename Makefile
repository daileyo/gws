# Makefile for gws - Git Workspace CLI

# Binary name
BINARY_NAME=git-workspace

# Build directory
BUILD_DIR=./build

# Coverage directory
COVERAGE_DIR=./coverage

# Install locations (XDG-compliant, overridable)
INSTALL_BIN  ?= $(HOME)/.local/bin
INSTALL_ZSH  ?= $(HOME)/.local/share/zsh/site-functions
INSTALL_BASH ?= $(HOME)/.local/share/bash-completion/completions

# Version information
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags to embed version information
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build test clean install uninstall help lint coverage snapshot ci vet fmt setup-hooks docs use-dev use-release

all: build

## setup-hooks: Install git hooks (pre-push linting, commit-msg formatting)
setup-hooks:
	@echo "Setting up git hooks..."
	git config core.hooksPath .githooks
	@echo "Git hooks installed:"
	@echo "  pre-push   — runs go vet, golangci-lint, and tests before each push"
	@echo "  commit-msg — pads conventional commit types for aligned git log output"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/git-workspace
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete"

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -v -race ./...
	@echo "Tests complete"

## coverage: Run tests with coverage report
coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	go test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report: $(COVERAGE_DIR)/coverage.html"
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -n 1

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

## snapshot: Build snapshot release with goreleaser (local testing)
snapshot:
	@echo "Building snapshot release..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser build --snapshot --clean; \
	else \
		echo "goreleaser not installed. Install with:"; \
		echo "  go install github.com/goreleaser/goreleaser/v2@latest"; \
		exit 1; \
	fi

## docs: Serve MkDocs site locally for preview (http://127.0.0.1:8000)
docs:
	@if [ ! -d .venv ]; then \
		echo "Creating Python virtual environment..."; \
		python3 -m venv .venv; \
		.venv/bin/pip install -r docs/requirements.txt; \
	fi
	@echo "Starting MkDocs dev server..."
	.venv/bin/python -m mkdocs serve

## ci: Run all CI checks (vet, lint, test with race detector)
ci: vet lint test-race
	@echo "All CI checks passed!"

## clean: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -rf dist
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

## uninstall: Remove installed binary and shell completions
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_BIN)/$(BINARY_NAME)
	@rm -f $(INSTALL_ZSH)/_$(BINARY_NAME)
	@rm -f $(INSTALL_BASH)/$(BINARY_NAME)
	@echo "Uninstall complete"

## install: Install binary and shell completions
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_BIN)..."
	@mkdir -p $(INSTALL_BIN)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "Installing zsh completion to $(INSTALL_ZSH)..."
	@mkdir -p $(INSTALL_ZSH)
	@$(INSTALL_BIN)/$(BINARY_NAME) completion zsh > $(INSTALL_ZSH)/_$(BINARY_NAME)
	@echo "Installing bash completion to $(INSTALL_BASH)..."
	@mkdir -p $(INSTALL_BASH)
	@$(INSTALL_BIN)/$(BINARY_NAME) completion bash > $(INSTALL_BASH)/$(BINARY_NAME)
	@echo ""
	@echo "Install complete!"
	@echo ""
	@echo "Add to your shell config:"
	@echo ""
	@echo "  zsh (~/.zshrc):"
	@echo "    export PATH=\"$(INSTALL_BIN):\$$PATH\""
	@echo "    eval \"\$$(git-workspace shell-init zsh)\""
	@echo ""
	@echo "  bash (~/.bashrc):"
	@echo "    export PATH=\"$(INSTALL_BIN):\$$PATH\""
	@echo "    eval \"\$$(git-workspace shell-init bash)\""

## use-dev: Switch to dev build (build and install to ~/.local/bin)
use-dev: install
	@echo ""
	@echo "Switched to dev build:"
	@$(INSTALL_BIN)/$(BINARY_NAME) --version

## use-release: Switch to released (Homebrew) build (remove dev binary)
use-release:
	@if [ ! -f $(INSTALL_BIN)/$(BINARY_NAME) ]; then \
		echo "Already using released build:"; \
	else \
		rm -f $(INSTALL_BIN)/$(BINARY_NAME); \
		echo "Removed dev binary from $(INSTALL_BIN)"; \
		echo "Switched to released build:"; \
	fi
	@$(BINARY_NAME) --version

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
