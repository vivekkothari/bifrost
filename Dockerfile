# Use the official Golang 1.23.1 image to build the application
FROM golang:1.23.1 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Ensure that the binary is built for Linux (matching the target OS) using CGO_ENABLED=0 and building statically
RUN CGO_ENABLED=0 GOOS=linux go build -o bifrost .

# Start a new stage from scratch
FROM alpine:latest

# Set working directory inside the new stage
WORKDIR /root/

# Install certificates in case HTTPS connections are needed by the app
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/bifrost .

# Ensure the binary has execution permissions
RUN chmod +x ./bifrost

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./bifrost"]
