#!/usr/bin/env bash
# Doctor script to validate Claude Code Docker sandbox setup

set -euo pipefail

# Source common functions
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

init_context "${1:-$PWD}"

printf '=== Claude Code Docker Doctor ===\n\n'

errors=0
warnings=0

# Check 1: Docker is installed
printf 'Checking Docker installation... '
if command -v docker >/dev/null 2>&1; then
  printf '\033[32mOK\033[0m\n'
  docker --version
else
  printf '\033[31mFAIL\033[0m\n'
  printf '  Docker is not installed or not in PATH.\n'
  ((errors++))
fi

# Check 2: Docker Desktop sandbox plugin
printf 'Checking Docker Desktop sandbox plugin... '
if docker sandbox version >/dev/null 2>&1; then
  printf '\033[32mOK\033[0m\n'
else
  printf '\033[31mFAIL\033[0m\n'
  printf '  Docker Desktop sandbox plugin is required.\n'
  printf '  Requires Docker Desktop 4.58 or later.\n'
  printf '  Enable sandbox plugin in Docker Desktop settings.\n'
  ((errors++))
fi

# Check 3: Ollama is running
printf 'Checking Ollama connectivity... '
if curl -fsS "${OLLAMA_URL}/api/tags" >/dev/null 2>&1; then
  printf '\033[32mOK\033[0m\n'
else
  printf '\033[31mFAIL\033[0m\n'
  printf '  Cannot reach Ollama at %s\n' "${OLLAMA_URL}"
  printf '  Ensure Ollama is running: ollama serve\n'
  ((errors++))
fi

# Check 4: Model exists in Ollama
printf 'Checking model %s... ' "${MODEL}"
if curl -fsS -o /dev/null "${OLLAMA_URL}/api/show" -d "{\"model\":\"${MODEL}\"}" 2>/dev/null; then
  printf '\033[32mOK\033[0m\n'
else
  printf '\033[31mFAIL\033[0m\n'
  printf '  Model %s not found in Ollama.\n' "${MODEL}"
  printf '  Pull it first: ollama pull %s\n' "${MODEL}"
  ((errors++))
fi

# Check 5: Workspace exists
printf 'Checking workspace directory... '
if [ -d "${WORKSPACE}" ]; then
  printf '\033[32mOK\033[0m\n'
  printf '  %s\n' "${WORKSPACE}"
else
  printf '\033[31mFAIL\033[0m\n'
  printf '  Workspace directory does not exist: %s\n' "${WORKSPACE}"
  ((errors++))
fi

# Check 6: Template exists (warning only)
printf 'Checking warm template... '
if template_exists; then
  printf '\033[32mOK\033[0m\n'
  printf '  %s\n' "${TEMPLATE_TAG}"
else
  printf '\033[33mWARN\033[0m\n'
  printf '  Warm template not found: %s\n' "${TEMPLATE_TAG}"
  printf '  First run will be slower. Run: make template\n'
  ((warnings++))
fi

# Check 7: Sandbox exists (info only)
printf 'Checking sandbox... '
if sandbox_exists; then
  printf '\033[32mEXISTS\033[0m\n'
  printf '  %s\n' "${SANDBOX_NAME}"
  if sandbox_running; then
    printf '  Status: \033[32mrunning\033[0m\n'
  else
    printf '  Status: \033[33mstopped\033[0m\n'
  fi
else
  printf '\033[33mNOT FOUND\033[0m\n'
  printf '  Will be created on first run.\n'
fi

# Summary
printf '\n'
printf '=== Summary ===\n'

if [ $errors -eq 0 ] && [ $warnings -eq 0 ]; then
  printf '\033[32mAll checks passed!\033[0m\n'
  printf 'Ready to run: make run\n'
  exit 0
elif [ $errors -eq 0 ]; then
  printf '\033[33m%d warning(s), 0 error(s)\033[0m\n' $warnings
  printf 'Setup is functional but could be improved.\n'
  exit 0
else
  printf '\033[31m%d error(s), %d warning(s)\033[0m\n' $errors $warnings
  printf 'Please fix the errors above before running.\n'
  exit 1
fi