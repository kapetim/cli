# cli toolchain image — one lean, pinned image with the `cli` binary and the
# native linters it shells out to (shellcheck, hadolint, actionlint).
#
# Other repositories use this image as the base for their workflows:
#   docker run --rm -v "$PWD:/repo" -w /repo kapetim/cli:<v> cli lint all
#
# No Node, no Python, no Go toolchain — markdown is linted by cli's built-in
# Go rules. Versions are pinned in docker/runtime/versions.env.

FROM alpine:3.21 AS build

COPY docker/runtime /install/
RUN . /install/versions.env \
    && apk add --no-cache "bash=${APK_BASH}" \
    && bash /install/install-core.sh \
    && bash /install/install-go.sh

WORKDIR /src
COPY go.mod ./
COPY src ./src
RUN go build -trimpath -ldflags "-s -w" -o /out/cli ./src

FROM alpine:3.21 AS final

COPY --from=build /out/cli /usr/local/bin/cli
COPY docker/runtime /install/
RUN . /install/versions.env \
    && apk add --no-cache "bash=${APK_BASH}" \
    && bash /install/install-core.sh \
    && bash /install/install-tools.sh \
    && rm -rf /install

WORKDIR /repo
CMD ["cli", "help"]
