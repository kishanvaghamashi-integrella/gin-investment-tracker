# ── dev stage ──────────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS dev

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8000

CMD ["air", "-c", ".air.toml"]

# ── builder stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /investment-tracker ./cmd/gin-investment-tracker

# ── prod stage ──────────────────────────────────────────────────────────────────
FROM alpine:3.21 AS prod

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /investment-tracker .

EXPOSE 8000

CMD ["./investment-tracker"]