ARG GO_VERSION=1.17.7

# Build
FROM golang:${GO_VERSION}-alpine AS build
 
RUN apk add --no-cache git

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/*.go

# Final
FROM gcr.io/distroless/static AS final

USER nonroot:nonroot

COPY --chown=nonroot:nonroot ./config /config/
 
COPY --from=build --chown=nonroot:nonroot /app /app


ENTRYPOINT ["/app"]