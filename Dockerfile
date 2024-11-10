# Use an official Golang image as the build stage
FROM golang:1.23.1 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install CA certificates (important for HTTPS requests)
RUN apt-get update && apt-get install -y ca-certificates

# Copy the source code into the container
COPY . .

# Build the application
RUN go build

# Use Ubuntu as the base image for the final container
FROM ubuntu:22.04

# Set the working directory in the final stage
WORKDIR /app

# Install required libraries and CA certificates
RUN apt-get update && apt-get install -y libc6 ca-certificates

# Copy the built executable from the builder stage
COPY --from=builder /app/fampayAssignment .

# Expose any required ports (e.g., 8080)
EXPOSE 8080

# Run the application
CMD ["./fampayAssignment"]