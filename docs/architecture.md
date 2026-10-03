# Architecture

`cli` is both a **Go module** and a **binary**, shipped as one lean Docker
image. Repositories import the module when they have code; when they have none,
they run the binary from the image.

## Two usage modes

| Mode | Consumer | How |
| --- | --- | --- |
| **CLI** | repos with shell scripts only | `kapetim/cli:<v> cli <command>` |
| **Import** | repos with Go code | `import "github.com/kapetim/cli/src/pkg/..."` |

The module is `github.com/kapetim/cli`; public packages live under `src/pkg/`,
private wiring under `src/internal/`, and the entrypoint under `src/cmd/cli`.

## Image

One lean image, `kapetim/cli:<v>`: the `cli` binary plus the minimal workflow
tooling (`git`/`git-lfs`, `jq`, `shellcheck`, `zip`/`unzip`). **No Go toolchain**
— Go repos bring their own. See [docker.md](docker.md) for the future-additions
list.

## Boundaries

- **Local validation/rendering** lives in the binary/module (Go), so there is no
  Node/Python runtime to carry.
- **Products keep their language**: browser-extensions is TS, data-science is
  Python — only their markdown/doc validation moves to `cli`.
- **Shell** in this repo (`scripts/`) is operational only.

## Versioning

`VERSION` is bare semver (`x.y.z`). Git tags use the Go rule (`vX.Y.Z`); Docker
tags use the Docker rule (`kapetim/cli:X.Y.Z`) — the same number.
