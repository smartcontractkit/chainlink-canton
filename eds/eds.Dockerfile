# build stage
FROM golang:1.26-alpine AS builder

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

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /disclosure-server /app/disclosure-server

ENV CONFIG_FILE=/app/config.toml
CMD ["sh", "-c", "./disclosure-server", "-config", "$CONFIG_FILE"]
