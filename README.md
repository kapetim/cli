# cli

The kapetim repository toolchain — a Go **module** and a **binary**, shipped as
one lean Docker image. Repositories either import the module in a small Go
script, or run the binary from the image when they have no code.

- **Binary:** `cli` — validation, rendering, generation.
- **Module:** `github.com/kapetim/cli/src/pkg/...` — import reusable functions.
- **Image:** `kapetim/cli:<v>` — the binary plus the minimal workflow tooling
  (`git`/`git-lfs`, `jq`, `shellcheck`, `zip`/`unzip`).

## Commands

| Command | Purpose |
| --- | --- |
| `cli validate [paths...]` | validate markdown table structure and cell widths |
| `cli render [paths...]` | render markdown to HTML (planned) |
| `cli version` | print the version |
| `cli help` | show usage |

## Use

Shell-only (no code) — run the image:

```bash
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.1.0 cli validate .
```

Go repo — import the module in a small script and use your own toolchain:

```go
import "github.com/kapetim/cli/src/pkg/validate"

func main() { _ = validate.Run(os.Args[1:]) }
```

## Build

```bash
go build -o cli ./src/cmd/cli
./cli version
```

## Layout

```text
src/cmd/cli     binary entrypoint
src/pkg/        public, importable packages
src/internal/   private wiring
Dockerfile      the single lean image
docker/runtime/ pinned installers + versions.env
scripts/        operational shell (bootstrap, release)
docs/           documentation
```

## Versioning

`VERSION` is bare semver (`x.y.z`). The git tag follows the Go rule (`vX.Y.Z`);
the Docker tag follows the Docker rule (`kapetim/cli:X.Y.Z`) — same number.
Tag-driven: `main` pushes auto-tag, and the release publishes the image and the
`cli` binaries.

## Docs

- [`docs/architecture.md`](docs/architecture.md) — what this repo is and why
- [`docs/docker.md`](docs/docker.md) — the image, contents, and future additions
- [`docs/README.md`](docs/README.md) — docs index
