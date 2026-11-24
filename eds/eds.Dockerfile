# build stage
FROM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS builder

# Layer 1: Download dependencies first to leverage Docker layer caching
WORKDIR /build
COPY go.mod go.sum ./
COPY ./authentication ./authentication
COPY ./contracts ./contracts
COPY ./eds ./eds
WORKDIR /build/eds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Layer 2: Build the application
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /disclosure-server ./cmd/server

FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

RUN apk --no-cache add ca-certificates

# Create a non-root user and group for running the application
RUN addgroup -S eds && adduser -S eds -G eds

COPY --from=builder /disclosure-server /app/disclosure-server

# Ensure the non-root user can read the application and config
RUN chown -R eds:eds /app && chmod -R 555 /app

ENV CONFIG_FILE=/app/config.toml

USER eds:eds

CMD ["/app/disclosure-server"]
