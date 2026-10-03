# Docker

One lean, pinned image, `kapetim/cli:<v>`, built on Alpine. It ships the `cli`
binary and the minimal workflow tooling.

## Contents

| Tool | Purpose |
| --- | --- |
| `cli` | the repository toolchain binary |
| `git`, `git-lfs` | clone/attr/ls-files, LFS repositories, branch archives |
| `zip`, `unzip` | repo/branch archives |
| `jq` | JSON processing in scripts |
| `shellcheck` | shell lint |
| `bash`, `coreutils`, `curl`, `ca-certificates` | base runtime |

The image intentionally has **no Go toolchain** — Go repos import the module
with their own Go setup.

## Build and run

```bash
docker build -t kapetim/cli:0.1.0 .
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.1.0 cli validate .
```

## Rules

- **Pinned.** Alpine base and apk packages are pinned in
  [`docker/runtime/versions.env`](../docker/runtime/versions.env).
- **Multi-stage.** The `build` stage compiles the binary; `final` ships it plus
  the tooling (no Go).
- **One image.** Tags are versioned only (`<v>`, no suffix, no `latest`).

## Future additions (deliberately deferred)

| Addition | Purpose | Est. compressed | Note |
| --- | --- | --- | --- |
| `actionlint` | CI lint | ~3 MB | apk |
| `hadolint` | Dockerfile lint | ~8 MB | release binary |
| `yq-go` | YAML ops | ~4 MB | apk |
| `make` | build/release glue | ~1 MB | apk |
| `docker-cli` (+ buildx/compose) | image ops | ~20–30 MB | needs socket mount |
| `gh` | workflow/asset retrieval | ~8 MB | |
| `openssh-client`, `gnupg` | signed/ssh git | ~6 MB | |
| `ripgrep`, `fd`, `tree`, `file` | search/util | ~3 MB | |
| `tar`, `gzip`, `xz` | archives | ~1 MB | |
| `lychee` | link checking | ~5 MB | if not covered by `cli` |
| `gitleaks` | secret scanning | ~4 MB | if not covered by `cli` |
| `markdownlint` | markdown lint | +Node ~20–30 MB | needs Node runtime |
| `yamllint` | YAML lint | +Python ~20–30 MB | needs Python runtime |

Prefer `cli`'s Go rules for markdown/YAML over the runtime-backed linters.
