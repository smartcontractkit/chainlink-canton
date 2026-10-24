# syntax=docker/dockerfile:1.7
FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder
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

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/committeeverifier /bin/
CMD ["/bin/committeeverifier"]
