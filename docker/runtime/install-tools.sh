#!/usr/bin/env bash
# Install the shared repository tooling: git, JSON/YAML, linters, archives.
set -euo pipefail
# shellcheck source=/dev/null
source "$(dirname "${BASH_SOURCE[0]}")/versions.env"

apk add --no-cache \
  git \
  git-lfs \
  jq \
  "yq-go=${APK_YQ_GO}" \
  "shellcheck=${APK_SHELLCHECK}" \
  "actionlint=${APK_ACTIONLINT}" \
  zip \
  unzip \
  make

# hadolint — Dockerfile linter. Not packaged for Alpine; fetched from the pinned
# release. Best-effort so the image still builds where GitHub is unreachable
# (the core tooling above is unaffected).
if curl -fsSL \
  "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-Linux-x86_64" \
  -o /usr/local/bin/hadolint; then
  chmod +x /usr/local/bin/hadolint
else
  echo "[warn] hadolint download failed (no GitHub access?) — skipping" >&2
  rm -f /usr/local/bin/hadolint
fi
