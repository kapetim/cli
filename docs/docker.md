# Docker

One pinned image, `kapetim/cli:<v>`, built on Alpine. It ships the `cli` binary,
the Go toolchain, and the shared repository tooling.

## Contents

| Tool | Purpose |
| --- | --- |
| `cli` | the repository toolchain binary |
| `go` | Go toolchain (build/test Go in the image) |
| `git`, `git-lfs` | clone/attr/ls-files, LFS repositories |
| `jq`, `yq` | JSON / YAML processing |
| `shellcheck` | shell lint |
| `hadolint` | Dockerfile lint |
| `actionlint` | GitHub Actions lint |
| `bash`, `coreutils`, `curl`, `ca-certificates` | base runtime |
| `zip`, `unzip`, `make` | archives / build glue |

## Build and run

```bash
docker build -t kapetim/cli:0.1.0 .
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.1.0 cli validate .
```

## Rules

- **Pinned.** Alpine base, apk packages, and release binaries are pinned in
  [`docker/runtime/versions.env`](../docker/runtime/versions.env).
- **Multi-stage.** The `build` stage compiles the binary; `final` ships it plus
  the toolchain.
- **One image.** No variants, no suffixes; tags are versioned only (`<v>`).

## Optional additions (not in the image)

`gh`, `rclone`, `docker` CLI, and PDF tooling (`qpdf`/`poppler`) are deliberately
excluded to keep the image lean. Add them via a consuming repo when needed.
