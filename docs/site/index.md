<p align="center"><img src="assets/images/gws-logo-hero.png" alt="gws logo" width="250"></p>

# git-workspace

[![CI](https://github.com/daileyo/gws/actions/workflows/ci.yml/badge.svg)](https://github.com/daileyo/gws/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/daileyo/gws)](https://goreportcard.com/report/github.com/daileyo/gws)
[![Snyk Security](https://snyk.io/test/github/daileyo/gws/badge.svg)](https://snyk.io/test/github/daileyo/gws)
[![Release](https://img.shields.io/github/v/release/daileyo/gws)](https://github.com/daileyo/gws/releases/latest)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A lightweight, cross-platform CLI tool for discovering, organizing, and navigating git repositories on your local system.

## Features

- **Repository Discovery**: Automatically find all git repositories in a directory tree
- **Automatic Classification**: Detect repository type (GitHub, GitLab, Azure DevOps, Bitbucket) from remote URLs
- **Subcommand-Based CLI**: Clean command structure with `list`, `init`, `add`, `refresh`, `tag`, and `user` subcommands
- **Git Status Integration**: View branch, clean/dirty state, and ahead/behind indicators at a glance
- **Smart Caching**: Fast status display with configurable cache (5-minute TTL)
- **Custom Tagging**: Organize repositories with `gws tag add` / `gws tag remove`
- **Advanced Filtering**: Search and filter repositories by type, tags, name, or path
- **User Profile Management**: Manage git user profiles across repositories with `gws user`
- **Repository Navigation**: Jump to any repository instantly with `gws <repo-name>`
- **Workspace Management**: Track and organize repositories in a centralized configuration
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Lightweight**: Single binary with no external dependencies

## Documentation

- [Getting Started](getting-started.md) — Install `git-workspace` and set up your workspace for the first time
- **Commands:**
    - [Core Commands](commands-core.md) — Listing, initializing, adding, refreshing, and navigating repositories
    - [User Management](commands-user.md) — Managing git user profiles across repositories
    - [Tagging](commands-tagging.md) — Organizing repositories with custom tags
    - [Legacy Flags](commands-legacy.md) — Deprecated flag-to-subcommand migration reference
- [Shell Integration](shell-integration.md) — Set up the `gws` shell function and workspace navigation
- [Configuration](configuration.md) — Config file location, structure, and field reference
