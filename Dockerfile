# ==================== Stage 1: Builder ====================
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
# Build the Migrate binary
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# ==================== Stage 2: Runtime ====================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/api .
COPY --from=builder /build/migrate .

# Copy migrations
COPY --from=builder /build/migrations ./migrations

# Copy entrypoint script
COPY --from=builder /build/entrypoint.sh .
RUN chmod +x entrypoint.sh

# Create uploads directory
RUN mkdir -p uploads/team/logos

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
