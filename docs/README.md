# cli docs

1. [Architecture](architecture.md) — the single-tool model and its boundaries
2. [Docker](docker.md) — the image, contents, and rules

**CI:** pull requests run Go checks plus an image smoke test — see
[`.github/workflows/test.yml`](../.github/workflows/test.yml). Releases are
tag-driven — see [`.github/workflows/release.yml`](../.github/workflows/release.yml).
