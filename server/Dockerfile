FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api

FROM alpine:3.22

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /build/api .

USER appuser

# The server listens on the port configured via LISTEN_ADDRESS (default :8080).
EXPOSE 8080

CMD ["./api"]
