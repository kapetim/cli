# cli toolchain image — one pinned image with the `cli` binary, the Go toolchain,
# and the shared repository tooling (git, jq/yq, shellcheck/hadolint/actionlint).
#
# Other repositories use this image as the base for their workflows:
#   docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:<v> cli validate .
#
# Multi-stage: the `build` stage compiles the binary; the `final` stage ships it
# plus the toolchain. All versions are pinned in docker/runtime/versions.env.

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
    && bash /install/install-go.sh \
    && rm -rf /install

WORKDIR /repo
CMD ["cli", "help"]
