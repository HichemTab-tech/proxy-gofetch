# Build stage
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Set default value for ALLOWED_DOMAINS
ENV ALLOWED_DOMAINS=""

# Expose port 3001
EXPOSE 3001

# Run the binary
CMD ["./main"]
