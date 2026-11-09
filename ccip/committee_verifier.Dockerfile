# syntax=docker/dockerfile:1.7
FROM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS builder
WORKDIR /build

# Copy module files first to maximize cache reuse (deps layer invalidated only when go.mod/go.sum change)
COPY go.mod go.sum ./
COPY ./authentication ./authentication
COPY ./contracts ./contracts
COPY ./eds ./eds
RUN --mount=type=cache,target=/go/pkg/mod,id=canton-go-mod \
    go mod download

# Copy the rest of the source; build layer invalidated when source changes
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build,id=canton-go-build \
    --mount=type=cache,target=/go/pkg/mod,id=canton-go-mod \
    CGO_ENABLED=0 go build -ldflags='-s -w' -o /bin/committeeverifier ./cmd/committeeverifier

FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/committeeverifier /bin/

CMD ["/bin/committeeverifier"]
