# cli

The kapetim repository toolchain — a Go **module** and a **binary**, shipped as
one lean Docker image. Repositories run the binary from the image, or import the
module in a small Go script.

- **Binary:** `cli` — linting and table validation.
- **Module:** `github.com/kapetim/cli/src/pkg/...` — reusable functions.
- **Image:** `kapetim/cli:<v>` — the binary plus native linters (`shellcheck`,
  `hadolint`, `actionlint`). No Node, no Python, no Go toolchain.

## Commands

| Command | Purpose |
| --- | --- |
| `cli lint [kinds...]` | `markdown` (built-in Go rules) + `shell` / `docker` / `ci` (native tools) |
| `cli validate [kinds...]` | checks — `tables` (manifest) |
| `cli scan tables` | write the table manifest (`tables.json`) |
| `cli version` | print the version |
| `cli help` | show usage |

Common flags: `--dir <path>` (repo root), `--manifest <path>` (default `tables.json`).

## Use

Shell / workflow (no code):

```bash
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.2.0 cli lint all
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.2.0 cli scan tables
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.2.0 cli validate tables
```

Go repo — import the module:

```go
import "github.com/kapetim/cli/src/pkg/run"

func main() { os.Exit(run.Main(run.Options{RepoDir: ".", Lint: []string{"all"}})) }
```

## Build

```bash
go build -o cli ./src
./cli version
```

## Layout

```text
src/main.go     binary entrypoint
src/pkg/        public packages (importable)
src/test/       unit tests
tests/          integration tests (shell)
docker/runtime/ pinned installers + versions.env
scripts/        operational shell
docs/           documentation
```

## Versioning

`VERSION` is bare semver (`x.y.z`). Git tags follow the Go rule (`vX.Y.Z`);
Docker tags follow the Docker rule (`kapetim/cli:X.Y.Z`) — same number.
Tag-driven: `main` pushes auto-tag, and the release publishes the image and the
`cli` binaries.

## Docs

- [`docs/architecture.md`](docs/architecture.md)
- [`docs/docker.md`](docs/docker.md)
- [`docs/api.md`](docs/api.md)
