#!/bin/bash

set -euo pipefail

MOCK_TARGET_DIR=""
MOCK_DST_ROOT=""
MODULE_PATH=""

print_help() {
  echo "Usage: $0 [-t TARGET_DIR] [-d DST_DIR] [-m MODULE_PATH]"
  echo ""
  echo "Options:"
  echo "  -t, --target     Target directory to scan for interfaces (e.g., internal/domain/service)"
  echo "  -d, --dst        Destination directory to output mocks (default: test/mock)"
  echo "  -m, --module     Go module import path (e.g., mandacode.com/accounts/auth)"
  echo "  -h, --help       Show this help message"
  exit 0
}

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
  -t | --target)
    MOCK_TARGET_DIR="$2"
    shift 2
    ;;
  -d | --dst)
    MOCK_DST_ROOT="$2"
    shift 2
    ;;
  -m | --module)
    MODULE_PATH="$2"
    shift 2
    ;;
  -h | --help)
    print_help
    ;;
  *)
    echo "❌ Unknown argument: $1"
    print_help
    ;;
  esac
done

# Prompt interactively if values missing
if [ -z "$MOCK_TARGET_DIR" ]; then
  read -rp "📂 Enter target directory (e.g., internal/domain/service): " MOCK_TARGET_DIR
fi

if [ -z "$MODULE_PATH" ]; then
  MODULE_PATH=$(go list -f '{{.ImportPath}}' "$MOCK_TARGET_DIR" 2>/dev/null || true)
  if [ -z "$MODULE_PATH" ]; then
    read -rp "🧩 Enter Go module import path (e.g., mandacode.com/accounts/auth): " MODULE_PATH
  fi
fi
if [ -z "$MOCK_DST_ROOT" ]; then
  read -rp "📦 Enter destination directory for mocks (default: test/mock): " MOCK_DST_ROOT
  MOCK_DST_ROOT=${MOCK_DST_ROOT:-test/mock}
fi

echo "🔍 Generating mocks from '$MOCK_TARGET_DIR'"
echo "📦 Module import path: $MODULE_PATH"
echo "📁 Output to: $MOCK_DST_ROOT"

find "$MOCK_TARGET_DIR" -type f -name '*.go' -not -path '*/mock_*' | while IFS= read -r src; do
  pkg_dir=$(dirname "$src")
  rel_dir=${pkg_dir#./}
  import_path="${MODULE_PATH}/${rel_dir}"

  interfaces=$(grep -E 'type [A-Z][a-zA-Z0-9_]+ interface' "$src" | awk '{print $2}')
  for iface in $interfaces; do
    out_dir="${MOCK_DST_ROOT}/${rel_dir}"
    mkdir -p "$out_dir"
    out_file="${out_dir}/mock_${iface}.go"

    echo "⚙️  Generating $iface → $out_file"
    if ! mockgen "$import_path" "$iface" >"$out_file"; then
      echo "❌ Failed to generate mock for $iface in $src"
      exit 1
    fi
  done
done

echo "✅ Mock generation complete."
