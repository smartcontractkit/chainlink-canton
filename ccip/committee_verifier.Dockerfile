# syntax=docker/dockerfile:1.7
FROM golang:1.25.6-alpine3.23 AS builder
WORKDIR /build

# Copy module files first to maximize cache reuse (deps layer invalidated only when go.mod/go.sum change)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=canton-go-mod \
    go mod download

# Copy the rest of the source; build layer invalidated when source changes
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build,id=canton-go-build \
    --mount=type=cache,target=/go/pkg/mod,id=canton-go-mod \
    CGO_ENABLED=0 go build -ldflags='-s -w' -o /bin/committeeverifier ./cmd/committeeverifier

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/committeeverifier /bin/
CMD ["/bin/committeeverifier"]
