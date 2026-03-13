#!/usr/bin/env bash
# Setup script for Claude Code Docker sandbox
# Performs full setup: checks prerequisites, creates template, runs doctor

set -euo pipefail

# Source common functions
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

init_context "${1:-$PWD}"

printf '=== Claude Code Docker Sandbox Setup ===\n\n'

# Check prerequisites
printf 'Checking prerequisites...\n'
require_cmd docker
require_cmd curl

printf '\n'

# Ensure Docker Desktop sandbox plugin
printf 'Checking Docker Desktop sandbox plugin...\n'
ensure_sandbox_plugin

printf '\n'

# Wait for Ollama
printf 'Checking Ollama connectivity...\n'
wait_for_ollama

printf '\n'

# Ensure model exists
ensure_model

printf '\n'

# Create warm template if needed
if template_exists; then
  printf 'Warm template already exists: %s\n' "${TEMPLATE_TAG}"
  printf 'To recreate, run: make template-clean template\n'
else
  printf 'Creating warm template...\n'
  "${ROOT_DIR}/scripts/bake-template.sh"
fi

printf '\n'

# Run doctor to verify setup
printf 'Running doctor to verify setup...\n'
"${ROOT_DIR}/scripts/doctor.sh" "${WORKSPACE}"

printf '\n'
printf '=== Setup Complete ===\n\n'
printf 'Ready to run Claude Code in sandbox!\n'
printf 'Next: make run\n'