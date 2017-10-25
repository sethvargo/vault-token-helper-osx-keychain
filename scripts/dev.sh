#!/usr/bin/env bash
set -e

#
# Helper script for local development. Automatically builds and registers the
# token helper. Requires `vault` is installed and available on $PATH.
#

# Get the right dir
DIR="$(cd "$(dirname "$(readlink "$0")")" && pwd)"
SCRATCH="$DIR/tmp"
mkdir -p "$SCRATCH"

function cleanup {
  echo ""
  echo "==> Cleaning up"
  kill -INT "$VAULT_PID"
  rm -rf "$SCRATCH"
}
trap cleanup EXIT

tee "$SCRATCH/vault.hcl" > /dev/null <<EOF
token_helper = "$SCRATCH/vault-token-helper"
EOF
export VAULT_CONFIG_PATH="$SCRATCH/vault.hcl"

go build -o "$SCRATCH/vault-token-helper"

vault server \
  -dev \
  -dev-root-token-id="root" \
  -log-level="debug" \
  &
sleep 2
VAULT_PID=$!

echo "==> Ready!"
wait $!
