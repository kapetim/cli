# cli toolchain image — one lean, pinned image with the `cli` binary and the
# minimal workflow tooling (git, git-lfs, zip/unzip, jq, shellcheck).
#
# Other repositories use this image as the base for their workflows:
#   docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:<v> cli validate .
#
# Go repos do NOT use this image to build; they import the module
# (github.com/kapetim/cli/src/pkg/...) with their own Go toolchain.
#
# Multi-stage: `build` compiles the binary; `final` ships it with the tooling.
# Versions are pinned in docker/runtime/versions.env.

FROM alpine:3.21 AS build

COPY docker/runtime /install/
RUN apk add --no-cache bash \
    && bash /install/install-core.sh \
    && bash /install/install-go.sh

WORKDIR /src
COPY go.mod ./
COPY src ./src
RUN go build -trimpath -ldflags "-s -w" -o /out/cli ./src/cmd/cli

FROM alpine:3.21 AS final

COPY --from=build /out/cli /usr/local/bin/cli
COPY docker/runtime /install/
RUN apk add --no-cache bash \
    && bash /install/install-core.sh \
    && bash /install/install-tools.sh \
    && rm -rf /install

WORKDIR /repo
CMD ["cli", "help"]
