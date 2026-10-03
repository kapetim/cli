#!/usr/bin/env bash
# Install the minimal workflow tooling shipped with the cli image.
# Kept deliberately lean; see docs/docker.md for the future-additions list.
set -euo pipefail
# shellcheck source=/dev/null
source "$(dirname "${BASH_SOURCE[0]}")/versions.env"

apk add --no-cache \
  "git=${APK_GIT}" \
  "git-lfs=${APK_GIT_LFS}" \
  "jq=${APK_JQ}" \
  "zip=${APK_ZIP}" \
  "unzip=${APK_UNZIP}" \
  "shellcheck=${APK_SHELLCHECK}"
