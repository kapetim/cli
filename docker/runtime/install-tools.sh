#!/usr/bin/env bash
# Install the native linters the cli shells out to. Kept tiny — no Node/Python.
set -euo pipefail
# shellcheck source=/dev/null
source "$(dirname "${BASH_SOURCE[0]}")/versions.env"

apk add --no-cache \
  "git=${APK_GIT}" \
  "git-lfs=${APK_GIT_LFS}" \
  "jq=${APK_JQ}" \
  "zip=${APK_ZIP}" \
  "unzip=${APK_UNZIP}" \
  "shellcheck=${APK_SHELLCHECK}" \
  "actionlint=${APK_ACTIONLINT}"

# hadolint — Dockerfile linter (pinned release binary; not packaged for Alpine).
# Best-effort so the image still builds where GitHub is unreachable.
if curl -fsSL "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-Linux-x86_64" \
  -o /usr/local/bin/hadolint; then
  chmod +x /usr/local/bin/hadolint
else
  echo "[warn] hadolint download failed (no GitHub access?) — skipping" >&2
  rm -f /usr/local/bin/hadolint
fi
