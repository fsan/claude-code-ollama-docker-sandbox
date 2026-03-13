#!/usr/bin/env bash
# Remove the warm template image

set -euo pipefail

# Source common functions
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

# Check prerequisites
require_cmd docker

# Check if template exists
if ! template_exists; then
  printf 'Template does not exist: %s\n' "${TEMPLATE_TAG}"
  exit 0
fi

printf 'Removing template: %s\n' "${TEMPLATE_TAG}"

docker image rm "${TEMPLATE_TAG}" 2>/dev/null || {
  printf 'Warning: Could not remove template image.\n' >&2
}

printf 'Template removed.\n'