#!/bin/sh
set -eu

lockfile_fingerprint="$(sha256sum package.json pnpm-lock.yaml | sha256sum | cut -d' ' -f1)"
fingerprint_file="node_modules/.deps-fingerprint"

if [ ! -x node_modules/.bin/vite ] || [ ! -f "$fingerprint_file" ] || [ "$(cat "$fingerprint_file" 2>/dev/null || true)" != "$lockfile_fingerprint" ]; then
  CI=true pnpm install --frozen-lockfile
  printf '%s' "$lockfile_fingerprint" > "$fingerprint_file"
fi

exec pnpm dev --host 0.0.0.0 --port "${VITE_PORT:-5173}"
