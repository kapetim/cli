# Architecture

`cli` is both a **Go module** and a **binary**, shipped as one lean Docker
image. Repositories import the module when they have code; when they have none,
they run the binary from the image.

## Two usage modes

| Mode | Consumer | How |
| --- | --- | --- |
| **CLI** | repos with shell workflows only | `kapetim/cli:<v> cli <command>` |
| **Import** | repos with Go code | `import "github.com/kapetim/cli/src/pkg/..."` |

The module is `github.com/kapetim/cli`; the entrypoint is `src/main.go`, public
packages live under `src/pkg/`, and unit tests under `src/test/`. Integration
tests (`tests/`) and shell (`scripts/`) sit outside the Go tree.

## Boundaries

- **Local validation** lives in the binary/module (Go) — markdown is linted by
  built-in Go rules, so there is no Node/Python runtime to carry.
- **Native linters** are shelled out to (`shellcheck`, `hadolint`, `actionlint`)
  and skipped when absent.
- **Products keep their language**: browser-extensions is TS, data-science is
  Python — only their markdown/doc validation moves to `cli`.

## Versioning

`VERSION` is bare semver (`x.y.z`). Git tags use the Go rule (`vX.Y.Z`); Docker
tags use the Docker rule (`kapetim/cli:X.Y.Z`) — the same number.
