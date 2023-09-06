# Use the official Go image as the base image
FROM golang:1.17 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download and cache the Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a minimal Alpine Linux image as the base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the build stage to the working directory
COPY --from=build /app/app .
COPY .env .

# Expose the port that your Go Fiber application will listen on
EXPOSE 6000

# Run the Go Fiber application when the container starts
CMD ["./app"]