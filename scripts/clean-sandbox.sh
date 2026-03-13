#!/usr/bin/env bash
# Remove sandbox completely

set -euo pipefail

# Source common functions
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

init_context "${1:-$PWD}"

# Check prerequisites
require_cmd docker
ensure_sandbox_plugin

# Check if exists
if ! sandbox_exists; then
  printf 'Sandbox does not exist: %s\n' "${SANDBOX_NAME}"
  exit 0
fi

printf 'Removing sandbox: %s\n' "${SANDBOX_NAME}"

# Stop first if running
if sandbox_running; then
  printf 'Stopping sandbox first...\n'
  docker sandbox stop "${SANDBOX_NAME}" 2>/dev/null || true
fi

# Remove
docker sandbox rm "${SANDBOX_NAME}" 2>/dev/null || {
  printf 'Warning: Could not remove sandbox completely.\n' >&2
}

printf 'Sandbox removed.\n'