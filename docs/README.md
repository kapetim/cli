# cli docs

1. [Architecture](architecture.md) — the module + binary model and its boundaries
2. [Docker](docker.md) — the lean image, contents, future additions
3. [API](api.md) — importable packages and the table manifest

**CI:** pull requests run Go vet/build/tests, the shell integration tests, and
an image smoke test — see [`.github/workflows/test.yml`](../.github/workflows/test.yml).
Releases are tag-driven — see [`.github/workflows/release.yml`](../.github/workflows/release.yml).
