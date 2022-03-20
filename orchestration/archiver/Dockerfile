ARG GO_VERSION=1.18.0

# Build
FROM golang:${GO_VERSION}-alpine AS build
 
RUN apk add --no-cache git

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

# Final
FROM gcr.io/distroless/static:nonroot AS final

USER nonroot:nonroot

COPY ./sql/migrations /sql/migrations

COPY --from=build /migrations /migrations

COPY ./config /config/

COPY --from=build /app /app

ENTRYPOINT ["/app"]