#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-readiness}"
PROBE_NAME="${MODE}Probe"
SECRET_FILE="/vault/secrets/secrets"
LOGFILE="/app/probe.log"
CMD_BINARY="/app/cmd"

# Load secrets if the file is readable
if [[ -r "$SECRET_FILE" ]]; then
  # shellcheck source=/dev/null
  source "$SECRET_FILE"
fi

# Log start of the probe
{
  ts=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  echo "${ts} [${PROBE_NAME}] Starting healthcheck (mode=${MODE})..."
} >>"$LOGFILE" 2>&1

# Execute the probe
if [[ -x "$CMD_BINARY" ]]; then
  "$CMD_BINARY" healthcheck "--${MODE}" >>"$LOGFILE" 2>&1
  status=$?
else
  {
    ts=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "${ts} [${PROBE_NAME}] ERROR: '$CMD_BINARY' not found or not executable"
  } >>"$LOGFILE" 2>&1
  exit 127
fi

# Log exit status
{
  ts=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  echo "${ts} [${PROBE_NAME}] Exit status = ${status}"
} >>"$LOGFILE" 2>&1

exit "$status"
