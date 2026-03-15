# Cloma CLI Design Spec

**Date:** 2026-03-15
**Status:** Implemented

## Overview

`cloma` is a Go CLI application that manages Docker Desktop sandboxes for running code agents in isolation. It provides workspace management, sandbox lifecycle operations, and connectivity to Ollama running on the host machine.

## Architecture

```
cloma/
├── cmd/cloma/main.go          # Entry point
├── internal/
│   ├── cmd/                    # Cobra commands
│   │   ├── root.go            # Root command + global flags
│   │   ├── run.go             # `cloma run` (default command)
│   │   ├── list.go            # `cloma list`
│   │   ├── shell.go           # `cloma shell`
│   │   ├── stop.go            # `cloma stop`
│   │   ├── clean.go           # `cloma clean`
│   │   ├── doctor.go          # `cloma doctor`
│   │   └── version.go         # `cloma version`
│   ├── sandbox/               # Docker sandbox operations
│   │   ├── client.go          # Docker sandbox API client
│   │   ├── create.go          # Create sandbox
│   │   ├── exec.go            # Exec into sandbox
│   │   ├── list.go            # List sandboxes
│   │   ├── network.go         # Network proxy config
│   │   ├── stop.go            # Stop sandbox
│   │   └── clean.go           # Remove sandbox
│   ├── workspace/             # Workspace management
│   │   ├── resolve.go         # Path resolution
│   │   ├── random.go          # Random workspace creation
│   │   └── naming.go          # Sandbox name generation
│   ├── ollama/                # Ollama connectivity
│   │   └── client.go          # Check model, connectivity
│   └── config/                # Configuration
│       └── config.go          # State dir, defaults
└── go.mod
```

## Command Structure

### Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Config file (default: `$HOME/.cloma/config.yaml`) |
| `-v, --verbose` | Verbose output (stackable: `-v`, `-vv`) |
| `--json` | Output in JSON format |

### Subcommands

#### `cloma run` (default)

Launch an agent in an isolated Docker sandbox.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--workspace` | `-w` | `.` (current dir) | Workspace directory |
| `--model` | `-m` | `glm-5:cloud` | AI model to use |
| `--port` | `-p` | `11434` | Ollama port |
| `--flags` | `-f` | (empty) | Additional agent flags |

Behavior:
1. Resolve workspace path (create random if not specified)
2. Check prerequisites (Docker, sandbox plugin, Ollama, model)
3. Create sandbox if needed
4. Configure network proxy for host access
5. Launch agent interactively

#### `cloma list`

List all cloma-managed sandboxes.

Output shows:
- Name (format: `cloma-{slug}-{hash}`)
- Status (running/stopped)
- Workspace (decoded from name)

Supports `--json` for machine-readable output.

#### `cloma shell`

Open an interactive shell in the sandbox.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--workspace` | `-w` | `.` (current dir) | Workspace directory |

#### `cloma stop`

Stop a running sandbox.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--workspace` | `-w` | `.` (current dir) | Workspace directory |

#### `cloma clean`

Remove a sandbox completely.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--workspace` | `-w` | `.` (current dir) | Workspace directory |
| `--force` | `-f` | `false` | Skip confirmation |

#### `cloma doctor`

Run health checks on the system.

Checks:
1. Docker installation
2. Docker Desktop sandbox plugin
3. Ollama connectivity
4. Model availability
5. Workspace directory
6. Warm template availability
7. Sandbox status

Exit codes:
- `0`: All checks passed or only warnings
- `1`: One or more checks failed

Supports `--json` for machine-readable output.

#### `cloma version`

Print version information.

Supports `--json` for machine-readable output.

## Workspace Resolution

1. **No workspace specified**: Create random workspace in `~/.cloma/workspaces/cloma-XXXXXX`
2. **`.` (dot)**: Resolve to current working directory
3. **`~` or `~/path`**: Resolve to home directory
4. **Relative path**: Resolve to absolute path
5. **Absolute path**: Use as-is

## Sandbox Naming

Sandbox names follow the pattern: `cloma-{slug}-{hash}`

- **slug**: Lowercase basename of workspace, special chars replaced with hyphens
- **hash**: First 8 characters of SHA256 hash of workspace absolute path

Example:
- Workspace: `/Users/fox/myproject`
- Sandbox: `cloma-myproject-a1b2c3d4`

## State Directory

All state is stored in `~/.cloma/`:

```
~/.cloma/
├── config.yaml           # Configuration (optional)
└── workspaces/          # Random workspaces created by `cloma`
    ├── cloma-a1b2c3d4/
    └── cloma-e5f6g7h8/
```

## Configuration

Environment variables override defaults:

| Variable | Description |
|----------|-------------|
| `CLOMA_MODEL` | AI model to use |
| `OLLAMA_PORT` | Host Ollama port |
| `OLLAMA_URL` | Ollama base URL |
| `CLOMA_TEMPLATE_TAG` | Template image tag |
| `CLOMA_STATE_DIR` | State directory |
| `CLOMA_WORKSPACES_DIR` | Workspaces directory |

## Build and Install

```bash
# Build
make build
# Output: bin/cloma

# Install to /usr/local/bin
make install

# Run directly
make cloma ARGS="--help"
```

## Dependencies

- Go 1.22+
- github.com/spf13/cobra v1.8.1
- github.com/spf13/viper v1.19.0

## Migration from Bash Scripts

The Go CLI is a full port of the existing bash scripts:

| Bash Script | Go Command |
|-------------|------------|
| `scripts/run-claude-code.sh` | `cloma run` |
| `scripts/shell.sh` | `cloma shell` |
| `scripts/stop-sandbox.sh` | `cloma stop` |
| `scripts/clean-sandbox.sh` | `cloma clean` |
| `scripts/doctor.sh` | `cloma doctor` |

The Makefile preserves the old bash targets for backwards compatibility.