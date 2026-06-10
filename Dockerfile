# Stage 1: Development (hot-reload via Air)
FROM golang:1.26-alpine AS development

RUN apk add --no-cache git curl
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 9900
CMD ["air", "-c", ".air.toml"]


# Stage 2: Builder
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-w -s" -o /app/bin/server ./cmd/server


# Stage 3: Production
FROM alpine:3.19 AS production

RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/server .
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 9900

HEALTHCHECK --interval=30s --timeout=10s --start-period=15s --retries=3 \
  CMD wget -qO- http://localhost:9900/health || exit 1

CMD ["./server"]
