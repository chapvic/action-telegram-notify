FROM golang:1.26 AS builder
WORKDIR /app

# Cache dependencies separately from sources
COPY go.mod go.sum ./
RUN go mod download

# Sources
COPY . .

# Static optimized binary
RUN CGO_ENABLED=0 go build \
      -trimpath \
      -ldflags="-s -w" \
      -o telegram-notify .

FROM alpine:3.20

LABEL org.opencontainers.image.source="https://github.com/chapvic/action-telegram-notify" \
      org.opencontainers.image.description="Send notification to Telegram" \
      org.opencontainers.image.licenses="MIT"

# CA certificates are REQUIRED for HTTPS to api.telegram.org
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/telegram-notify /telegram-notify

# NOTE: root is required so the action can write to GITHUB_OUTPUT.
ENTRYPOINT ["/telegram-notify"]
