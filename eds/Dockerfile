# build stage
FROM golang:1.23-alpine AS builder

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
COPY eds/environments.yaml .

EXPOSE 8090

ENV PORT=8090
ENV HOST=0.0.0.0
ENV ENVIRONMENTS_CONFIG=/app/environments.yaml

CMD ["./disclosure-server"]
