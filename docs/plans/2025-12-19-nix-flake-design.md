# Nix Flake Setup Design

**Date:** 2025-12-19
**Status:** Approved

## Overview

Add Nix flake support to procmap to provide reproducible builds and a standardized development environment. The flake will support both package building (`nix build`) and development shells (`nix develop`).

## Requirements

- Build procmap binary via Nix
- Provide development environment with Go tooling
- Support both bash and nushell users
- Use latest stable Go version (1.23.x)
- Multi-platform support (Linux x86_64/aarch64, macOS x86_64/aarch64)

## Flake Outputs

### 1. Package Output (`packages.default`)

Uses `buildGoModule` to create a reproducible build of procmap:

- **Builder:** `pkgs.buildGoModule`
- **Go Version:** Latest stable from nixpkgs (1.23.x)
- **Metadata:** Derived from git (version) and go.mod (module name)
- **vendorHash:** Fixed-output hash of Go dependencies
  - Initial build will fail with hash mismatch
  - Copy correct hash from error message into flake.nix
  - Standard Nix workflow for Go projects

**Build outputs:**
- Binary: `procmap`
- Accessible via: `nix build`, `nix run`, or as flake dependency

### 2. Development Shell (`devShells.default`)

Provides complete Go development environment:

**Tools included:**
- Go 1.23.x toolchain (gofmt, goimports, etc.)
- gopls - Official Go Language Server Protocol implementation
- delve - Go debugger
- staticcheck - Comprehensive static analyzer

**Behavior:**
- Does not vendor dependencies - uses standard `go mod download`
- Respects existing go.mod/go.sum files
- All tools available on PATH
- Normal Go workflow: `go run`, `go build`, `go test`

## Shell Integration

### Bash/Zsh Support

Standard `shellHook` that displays welcome message with:
- Available tools and versions
- Quick command reference
- Nix build instructions

### Nushell Support

- Generate `.nix-shell.nu` environment script
- Include nushell-formatted welcome message
- shellHook detects nushell and suggests sourcing the script
- User can optionally add to nushell config for auto-sourcing

**Welcome message format:**
```
🚀 procmap development environment

Tools available:
  go 1.23.x - Go toolchain
  gopls - Language server
  delve - Debugger
  staticcheck - Linter

Quick commands:
  go run main.go - Run the app
  go test ./... - Run tests
  nix build - Build with Nix
```

## File Changes

### New Files

1. **flake.nix** - Main flake definition
   - Package build configuration
   - Development shell definition
   - Multi-system support via flake-utils or manual

2. **flake.lock** - Auto-generated dependency lock file
   - Tracks nixpkgs version
   - Ensures reproducible builds

3. **.nix-shell.nu** - Generated nushell environment (gitignored)

### Modified Files

1. **.gitignore**
   - Add `.nix-shell.nu`
   - Add `result` (symlink from `nix build`)

2. **go.mod**
   - Update `go 1.21` → `go 1.23`
   - Dependencies remain unchanged

## Flake Metadata

- **Description:** "A terminal-based process visualization tool"
- **Inputs:** nixpkgs (unstable channel for latest Go)
- **Systems:** x86_64-linux, aarch64-linux, x86_64-darwin, aarch64-darwin
- **License:** TBD (check existing or suggest MIT/Apache-2.0)

## Initial Setup Workflow

1. Create flake.nix with placeholder vendorHash
2. Run `nix build` - will fail with hash mismatch
3. Copy actual hash from error message
4. Update vendorHash in flake.nix
5. Rebuild - should succeed
6. Test `nix develop` shell functionality

## Benefits

- **Reproducibility:** Same build everywhere via vendorHash
- **Isolation:** Development environment doesn't pollute system
- **Discoverability:** New contributors run `nix develop` and get everything
- **CI/CD Ready:** Can use in GitHub Actions, GitLab CI, etc.
- **Multi-platform:** Same flake works on Linux and macOS, x86 and ARM

## Non-Goals

- NixOS module for system-wide installation (future enhancement)
- Container image generation (future enhancement)
- Automatic go.mod updates (user maintains this manually)
- Pre-vendoring in devShell (uses standard Go module workflow)
