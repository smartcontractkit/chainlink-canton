# build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY . .

WORKDIR /build/eds
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /disclosure-server ./cmd/server

# runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /disclosure-server .

ENV CONFIG_FILE=/app/config.toml

CMD ["sh", "-c", "./disclosure-server", "-config", "$CONFIG_FILE"]
