# Use an official Go runtime as a parent image
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app from the correct main file location
RUN go build -o main ./cmd/server/main.go

# Start a new stage from scratch
FROM alpine:latest  

# Install ca-certificates for SSL (if required)
RUN apk --no-cache add ca-certificates

# Add a new non-root user
RUN adduser -D appuser

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Change ownership of the application files to the non-root user
RUN chown -R appuser:appuser /root

# Switch to the non-root user
USER appuser

# Expose the port the app runs on
EXPOSE 80

# Run the Go binary
CMD ["./main"]
