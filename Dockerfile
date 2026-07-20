# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create a minimal production image
FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .
# Copy .env file if present (useful for local container tests, though Azure uses environment variables)
COPY --from=builder /app/.env* ./

# Expose port (Azure App Service will map this automatically)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
