# Use official Golang image for building
FROM golang:1.24-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside container
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Use golang:alpine as the final runtime image (Fix for missing Go)
FROM golang:1.24-alpine

# Set working directory inside container
WORKDIR /app

# Install required dependencies
RUN apk add --no-cache git

# Copy application source code for live reloading
COPY --from=builder /go/bin/air /bin/air
COPY . .

# Expose the application port
EXPOSE 3000

# Command to run with live reloading
CMD ["air", "-c", ".air.toml"]
