# ==========================================
# STAGE 1: Build the Svelte Frontend SPA
# ==========================================
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend

# Install dependencies
COPY ../frontend/package*.json ./
RUN npm ci

# Copy code and build Svelte application
COPY ../frontend/ ./
RUN npm run build

# ==========================================
# STAGE 2: Build the Go Backend
# ==========================================
FROM golang:1.25-alpine AS backend-builder
WORKDIR /backend
RUN apk add --no-cache git gcc musl-dev

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code and compile binary (strip debug symbols and symbols table with -s -w)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go

# ==========================================
# STAGE 3: Final Minimal Production Image
# ==========================================
FROM alpine:3.19
WORKDIR /app

# Install certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Copy compiled Go binary
COPY --from=backend-builder /backend/main ./main

# Copy compiled Svelte assets
COPY --from=frontend-builder /frontend/dist ./dist

# Create folder for persistent keys
RUN mkdir -p /app/keys

# Set environment variables for production
ENV ENV=prod
ENV PORT=8080
ENV RSA_PRIVATE_KEY_PATH=/app/keys/private_key.pem
ENV RSA_PUBLIC_KEY_PATH=/app/keys/public_key.pem

EXPOSE 8080

CMD ["./main"]
