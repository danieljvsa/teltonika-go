# Use the official Golang image to build the binary
FROM golang:1.21 as builder

# Set working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o myserver main.go

# Create a minimal image to run the Go server
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/myserver .

# Expose ports for TCP and UDP
EXPOSE 8080 9090/udp

# Run the binary
CMD ["./myserver"]
