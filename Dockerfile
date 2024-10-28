# Use a lightweight base image
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go source code
COPY . .

# Build the Go application
RUN go build -o app .

# Create a runtime image
FROM alpine:latest

# Copy the built binary
COPY --from=builder /app/app .

# Expose the port
EXPOSE 8009
ENV POSTGRES_HOST postgres
ENV POSTGRES_PORT 5432

# Start the server
CMD ["./app"]