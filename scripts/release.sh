#!/usr/bin/env bash
# release.sh — build the cli image locally, tagged with VERSION.
#
# Usage: scripts/release.sh [--push]
#   --push   also push kapetim/cli:<v> to Docker Hub (requires `docker login`).
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="$(tr -d '[:space:]' < "$REPO_DIR/src/VERSION")"
IMAGE="kapetim/cli:${VERSION}"

PUSH=0
[[ "${1:-}" == "--push" ]] && PUSH=1

echo "[release] building ${IMAGE}"
docker build -t "$IMAGE" "$REPO_DIR"

if (( PUSH )); then
  echo "[release] pushing ${IMAGE}"
  docker push "$IMAGE"
fi

echo "[release] done — ${IMAGE}"
