## Build stage 
FROM golang:1.22.2 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy everything else into the current image.
COPY ./ ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# Start fresh from a smaller image
FROM alpine:3.14

# Copy the binary from the build stage
COPY --from=build /main .

# Expose port 8000
EXPOSE 8000

# Run the binary when the container starts
CMD ["/main"]