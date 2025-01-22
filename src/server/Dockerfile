# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.23.5 AS builder

# Create and change to the work directory.
WORKDIR /app

# Copy the go mod and sum files.
COPY go.mod .
COPY go.sum .

# Download dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the application. This creates the executable.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server server.go

# Use a Docker multi-stage build to create a lean production image.
# Start from the alpine image to keep the final image small.
FROM alpine:latest

# Install ca-certificates in case you need HTTPS/TLS.
RUN apk --no-cache add ca-certificates

# Set the working directory in the container.
WORKDIR /root/

# Copy the binary from the builder stage to the production image.
COPY --from=builder /app/server .

# Run the server executable when the container starts.
CMD ["./server"]

# Document that the service listens on port 50051.
EXPOSE 50051
