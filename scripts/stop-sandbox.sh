#!/usr/bin/env bash
# Stop running sandbox

set -euo pipefail

# Source common functions
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

init_context "${1:-$PWD}"

# Check prerequisites
require_cmd docker
ensure_sandbox_plugin

# Ensure sandbox exists
if ! sandbox_exists; then
  printf 'Sandbox does not exist: %s\n' "${SANDBOX_NAME}" >&2
  exit 0
fi

# Check if running
if ! sandbox_running; then
  printf 'Sandbox is not running: %s\n' "${SANDBOX_NAME}"
  exit 0
fi

printf 'Stopping sandbox: %s\n' "${SANDBOX_NAME}"

docker sandbox stop "${SANDBOX_NAME}"

printf 'Sandbox stopped.\n'