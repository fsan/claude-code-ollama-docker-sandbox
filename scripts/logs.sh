#!/usr/bin/env bash
# View sandbox logs

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
  exit 1
fi

printf 'Logs for sandbox: %s\n\n' "${SANDBOX_NAME}"

docker sandbox logs "${SANDBOX_NAME}" 2>&1 || {
  printf 'No logs available or sandbox not running.\n'
}