# Use the official Golang image
FROM golang:1.23 AS builder

# Set environment variables
ENV GO111MODULE=on \
    GOPATH=/go \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside container
WORKDIR /app

# Install Air (hot reloading tool)
# Install Air (hot reloading tool) with the correct repository
RUN go install github.com/air-verse/air@latest

# Copy only go.mod and go.sum first (to optimize caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Expose application port (change it based on your app)
EXPOSE 8080

# Start the application with Air for live reloading
CMD ["air"]

