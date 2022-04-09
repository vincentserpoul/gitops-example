ARG GO_VERSION=1.18.0

# Build
FROM golang:${GO_VERSION} AS build

ENV XZ_UTILS_VERSION=5.2.5-2
ENV UPX_VERSION=3.96

RUN apt-get update && \
    apt-get install -y xz-utils=${XZ_UTILS_VERSION} --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

ADD https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz /usr/local
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN xz -d -c /usr/local/upx-${UPX_VERSION}-amd64_linux.tar.xz | tar -xOf - upx-${UPX_VERSION}-amd64_linux/upx > /bin/upx && chmod a+x /bin/upx

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

ARG LAST_MAIN_COMMIT_HASH
ARG LAST_MAIN_COMMIT_TIME
ENV FLAG="-X main.CommitTime=${LAST_MAIN_COMMIT_TIME}"
ENV FLAG="$FLAG -X main.CommitHash=${LAST_MAIN_COMMIT_HASH}"

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -ldflags "-s -w $FLAG" \
    -buildvcs=true \
    -o /app ./cmd/serve/*.go

# Final
FROM gcr.io/distroless/static:latest-amd64 AS final

USER nonroot:nonroot

COPY --chown=nonroot:nonroot ./config /config

COPY --from=build --chown=nonroot:nonroot /app /app

ENTRYPOINT ["/app"]