#!/usr/bin/env bash
# Open an interactive shell inside the sandbox

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
  printf 'Run: make setup && make run\n' >&2
  exit 1
fi

# Ensure sandbox is running
if ! sandbox_running; then
  printf 'Starting sandbox: %s\n' "${SANDBOX_NAME}"
  docker sandbox start "${SANDBOX_NAME}"
fi

printf 'Opening shell in sandbox: %s\n' "${SANDBOX_NAME}"
printf 'Workspace: %s\n\n' "${WORKSPACE}"

docker sandbox exec \
  -u agent \
  -w "${WORKSPACE}" \
  -it \
  "${SANDBOX_NAME}" \
  bash