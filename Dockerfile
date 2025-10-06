## Multi-stage build for Secret Hitler (web + Go backend)

# 1) Web build stage
FROM node:20-slim AS web-builder
WORKDIR /app/web

# Install deps first for better layer caching
COPY web/package*.json ./
RUN npm ci --no-audit --no-fund

# Copy the rest of the web source and build
COPY web/ .

# Provide prod defaults for API and WS base URLs if not supplied at build
ARG VITE_BASE_URL=""
ARG VITE_BASE_URL_WS=""
ENV VITE_BASE_URL=$VITE_BASE_URL \
    VITE_BASE_URL_WS=$VITE_BASE_URL_WS

RUN npm run build


# 2) Go build stage
FROM golang:1.23-alpine AS go-builder
WORKDIR /app

# Install build tools
RUN apk add --no-cache ca-certificates git

# Go module files
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

# Copy backend source
COPY backend/ .

# Copy built web assets into backend/static where the router serves from
COPY --from=web-builder /app/web/dist ./static

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main main.go


# 3) Final runtime image
FROM alpine:3.20
WORKDIR /app

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

# Copy binary and static assets
COPY --from=go-builder /app/backend/bin/main ./main
COPY --from=go-builder /app/backend/static ./static

# Certificates for HTTPS requests if needed
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Environment for production
ENV ENV=production

# Expose backend port
EXPOSE 8080

USER app

CMD ["/app/main"]


