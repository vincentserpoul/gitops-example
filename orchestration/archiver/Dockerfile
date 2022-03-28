ARG GO_VERSION=1.18.0

# Build
FROM golang:${GO_VERSION} AS build

RUN apt-get update && apt-get install -y xz-utils && rm -rf /var/lib/apt/lists/*

ENV UPX_VERSION=3.96

ADD https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-${UPX_VERSION}-amd64_linux.tar.xz | tar -xOf - upx-${UPX_VERSION}-amd64_linux/upx > /bin/upx && chmod a+x /bin/upx

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/serve/*.go

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /migrations ./cmd/migrations/*.go

RUN strip --strip-unneeded /app
RUN upx /app

RUN strip --strip-unneeded /migrations
RUN upx /migrations

# Final
FROM gcr.io/distroless/static AS final

USER nonroot:nonroot

COPY --chown=nonroot:nonroot ./sql/migrations /sql/migrations

COPY --from=build --chown=nonroot:nonroot /migrations /migrations

COPY --chown=nonroot:nonroot ./config /config

COPY --from=build --chown=nonroot:nonroot /app /app

ENTRYPOINT ["/app"]