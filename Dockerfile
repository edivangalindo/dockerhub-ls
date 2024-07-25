# Use the official Golang image as a build stage
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o dockerhub-ls

# Use a minimal image as the final stage
FROM alpine:latest

# Set environment variable
ENV GO_ENV=production

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/dockerhub-ls /usr/local/bin/dockerhub-ls

# Command to run the executable
ENTRYPOINT ["/usr/local/bin/dockerhub-ls"]
