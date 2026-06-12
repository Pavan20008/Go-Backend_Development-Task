# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

# Build the static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# ---- Run stage ----
FROM alpine:3.20

RUN adduser -D -u 10001 appuser
WORKDIR /app

COPY --from=builder /app/server /app/server

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/server"]
