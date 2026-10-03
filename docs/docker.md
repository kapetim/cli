# Docker

One lean, pinned image, `kapetim/cli:<v>`, built on Alpine. It ships the `cli`
binary and the native linters it calls. **No Node, no Python, no Go toolchain.**

## Contents

| Tool | Purpose |
| --- | --- |
| `cli` | the repository toolchain binary |
| `shellcheck` | shell lint |
| `actionlint` | GitHub Actions lint |
| `hadolint` | Dockerfile lint |
| `git`, `git-lfs` | clone/attr/ls-files, LFS repositories |
| `jq` | JSON processing |
| `zip`, `unzip` | archives |
| `bash`, `coreutils`, `curl`, `ca-certificates` | base runtime |

## Build and run

```bash
docker build -t kapetim/cli:0.1.0 .
docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:0.1.0 cli lint all
```

## Rules

- **Pinned.** Alpine base and apk packages are pinned in
  [`docker/runtime/versions.env`](../docker/runtime/versions.env).
- **Multi-stage.** The `build` stage compiles the binary; `final` ships it plus
  the linters (no Go).
- **One image.** Tags are versioned only (`<v>`, no suffix, no `latest`).

## Future additions (deliberately deferred)

| Addition | Purpose | Est. compressed | Note |
| --- | --- | --- | --- |
| `yq-go` | YAML ops | ~4 MB | apk |
| `make` | build/release glue | ~1 MB | apk |
| `docker-cli` (+ buildx/compose) | image ops | ~20–30 MB | needs socket mount |
| `gh` | workflow/asset retrieval | ~8 MB | |
| `openssh-client`, `gnupg` | signed/ssh git | ~6 MB | |
| `ripgrep`, `fd`, `tree`, `file` | search/util | ~3 MB | |
| `markdownlint` | markdown lint | +Node ~20–30 MB | needs Node — avoided; cli has built-in rules |
| `yamllint` | YAML lint | +Python ~20–30 MB | needs Python |
