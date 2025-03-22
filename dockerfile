# Use official Golang image as the base image
FROM golang:1.21-alpine as builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (for better caching)
COPY go.mod ./
# If you had a go.sum file, you'd copy it too: COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
# -o specifies output name, CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a smaller base image for the final container
FROM alpine:latest

# Add some basic metadata
LABEL maintainer="Your Name <your.email@example.com>"
LABEL version="1.0"
LABEL description="Simple CRUD API in Golang"

# Install ca-certificates for HTTPS calls if needed
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]