# Start with a minimal base image
FROM golang:1.23.5 as builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies and build the API
RUN go mod tidy
RUN go build -o api .

# Use a minimal image for deployment
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install CA certificates (needed for Keycloak HTTPS)
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app .

# Expose API port
EXPOSE 8080

# Run the API
CMD ["./api"]
