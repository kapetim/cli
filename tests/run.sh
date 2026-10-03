#!/usr/bin/env bash
# Integration tests — build the cli and run it against a fixture repo.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="$(mktemp -d)/cli"
( cd "$ROOT" && go build -o "$BIN" ./src )

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK" "${BIN%/*}"' EXIT

cat > "$WORK/doc.md" <<'EOF'
# Doc

<!-- begin table -->
| A | B |
| --- | --- |
| 1 | 2 |
<!-- end table -->
EOF

# scan creates the manifest
"$BIN" scan tables --dir "$WORK" --manifest tables.json
test -f "$WORK/tables.json"

# validate passes against a fresh manifest
"$BIN" validate tables --dir "$WORK" --manifest tables.json

# markdown lint passes on balanced markers
"$BIN" lint markdown --dir "$WORK"

# drift: an extra table makes validate fail (exit 2)
cat >> "$WORK/doc.md" <<'EOF'

<!-- begin table -->
| C |
| --- |
| 3 |
<!-- end table -->
EOF
set +e
"$BIN" validate tables --dir "$WORK" --manifest tables.json
code=$?
set -e
if [ "$code" -ne 2 ]; then
  echo "expected exit 2 on drift, got $code" >&2
  exit 1
fi

echo "[OK] integration tests passed"
