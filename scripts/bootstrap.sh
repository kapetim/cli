#!/usr/bin/env bash
# bootstrap.sh — install the `cli` binary for local use.
#
# Downloads the pinned release binary for this OS/arch. The Docker image
# (kapetim/cli:<v>) is the primary distribution; this is for running `cli`
# without Docker.
#
# Usage: scripts/bootstrap.sh [--version <v>] [--bin-dir <dir>]
set -euo pipefail

OWNER="kapetim"
REPO="cli"
VERSION="${CLI_VERSION:-$(tr -d '[:space:]' < "$(dirname "${BASH_SOURCE[0]}")/../VERSION")}"
BIN_DIR="${CLI_BIN_DIR:-$PWD/.bin}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version) VERSION="${2:?--version requires a value}"; shift 2 ;;
    --bin-dir) BIN_DIR="${2:?--bin-dir requires a value}"; shift 2 ;;
    *) echo "unknown argument: $1" >&2; exit 1 ;;
  esac
done

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac

asset="cli-${os}-${arch}"
url="https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}/${asset}"

mkdir -p "$BIN_DIR"
echo "[bootstrap] fetching ${asset} (${VERSION})"
curl -fsSL "$url" -o "$BIN_DIR/cli"
chmod +x "$BIN_DIR/cli"

echo "[bootstrap] installed $("$BIN_DIR/cli" version) -> $BIN_DIR/cli"
echo "[bootstrap] add to PATH: export PATH=\"$BIN_DIR:\$PATH\""
