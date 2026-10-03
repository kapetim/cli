# cli

The kapetim repository toolchain — one Go binary and one pinned Docker image.
Repositories use it as the base for their workflows (validation, markdown
rendering, and shared tooling) instead of each wiring their own.

- **Binary:** `cli` — Go, static; local file validation and rendering.
- **Image:** `kapetim/cli:<v>` — the binary, the Go toolchain, and the shared
  tooling (`git`, `jq`/`yq`, `shellcheck`, `hadolint`, `actionlint`).

## Commands

| Command | Purpose |
| --- | --- |
| `cli validate [paths...]` | validate markdown table structure and cell widths |
| `cli render [paths...]` | render markdown to HTML (planned) |
| `cli version` | print the version |
| `cli help` | show usage |

## Build

```bash
go build -o cli ./src/cmd/cli
./cli version
```

## Image

```bash
docker build -t kapetim/cli:0.1.0 .
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.1.0 cli validate .
```

## Layout

```text
src/            Go source (cmd/cli + internal packages)
Dockerfile      the single image
docker/runtime/ pinned installers + versions.env
scripts/        operational shell (bootstrap, release)
docs/           documentation
```

## Docs

- [`docs/architecture.md`](docs/architecture.md) — what this repo is and why
- [`docs/docker.md`](docs/docker.md) — the image, contents, and rules
- [`docs/README.md`](docs/README.md) — docs index

## Release

Tag-driven: pushing a bare `X.Y.Z` tag creates a GitHub release and publishes
the image. `VERSION` is bumped in the merged PR; `main` pushes auto-tag it.
