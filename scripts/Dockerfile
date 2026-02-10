# Build stage
FROM golang:1.25.5-alpine AS builder

# Install necessary packages for building
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files and env template
COPY go.mod go.sum .env.example ./

# Download dependencies
RUN go mod download

# Install swag for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate Swagger documentation
RUN /go/bin/swag init -g cmd/main.go -o docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o parser ./cmd

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary, env template, and Swagger docs from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env
COPY --from=builder /app/docs ./docs/

# Expose port
EXPOSE 3030

# Command to run
CMD ["./parser"]