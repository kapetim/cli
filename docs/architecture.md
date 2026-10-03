# Architecture

`cli` is the kapetim repository toolchain: **one Go binary** and **one Docker
image**. Repositories use it as the base for their workflows instead of each
installing and wiring their own validation tooling.

## Model

- **Go** — the only language here. Local file validation, markdown rendering,
  and generation live in the binary, so there is no Node/Python runtime to
  carry around.
- **Shell** — operational glue for this repo only (`scripts/`).
- **No Node, no Python.** Markdown/YAML/JSON rules are implemented in Go;
  anything else is a pinned native binary in the image.

Complex, dynamic web scraping is **not** here — it lives in the data-science
repo, where those dependencies belong.

## Boundaries

| Concern | Where |
| --- | --- |
| Validate markdown (tables, cells, structure) | `cli` binary |
| Render markdown to HTML (planned) | `cli` binary |
| Generate repo manifests (planned) | `cli` binary |
| Shared tooling (git, jq/yq, lint trio) | image |
| Build/test Go | image + Go toolchain |

## Image

One image, `kapetim/cli:<v>`. It carries the `cli` binary, the Go toolchain,
and the shared tooling. Other repos do:

```bash
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:<v> cli validate .
```

## Naming

`cli` everywhere: the repo, the binary, and the image (`kapetim/cli`). The Go
module is `github.com/kapetim/cli`.
