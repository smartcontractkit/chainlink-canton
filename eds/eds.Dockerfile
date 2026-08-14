# build stage
FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

# Layer 1: Download dependencies first to leverage Docker layer caching
WORKDIR /build
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Layer 2: Build the application
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /disclosure-server ./eds/cmd/server

FROM alpine:3.22.4@sha256:310c62b5e7ca5b08167e4384c68db0fd2905dd9c7493756d356e893909057601

RUN apk --no-cache add ca-certificates

# Create a non-root user and group for running the application
RUN addgroup -S eds && adduser -S eds -G eds

COPY --from=builder /disclosure-server /app/disclosure-server

# Ensure the non-root user can read the application and config
RUN chown -R eds:eds /app && chmod -R 555 /app

ENV CONFIG_FILE=/app/config.toml

USER eds:eds

CMD ["/app/disclosure-server"]
