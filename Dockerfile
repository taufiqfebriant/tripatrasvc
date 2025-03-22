FROM golang:1.24-alpine

WORKDIR /app

# Install Air
RUN go install github.com/air-verse/air@latest

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Expose port 1323 (matches your Echo server port)
EXPOSE 1323

# Command to run Air
CMD ["air"]