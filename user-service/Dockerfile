# STAGE 1: Build the application
FROM golang:1.26-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# STAGE 2: Run the application
# Use 'scratch' for the smallest possible image, or 'alpine' for a full shell
FROM alpine:latest  

WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/app .

# Expose the application port (if applicable)
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
