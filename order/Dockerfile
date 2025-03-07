# Use an official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go modules manifests first (for efficient caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to include the internal package
COPY . .

# Build the order service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /order ./order

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /order .

# Expose the service port
EXPOSE 8080

# Run the order service
CMD ["./order"]
