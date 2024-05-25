FROM golang:1.22.2 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy everything
COPY ./ ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# Start fresh from a smaller image
FROM alpine:3.14

# Copy the binary from the build stage
COPY --from=build /main .
COPY --from=build /app/web /web

# Expose port
EXPOSE 8000

# Run the binary
CMD ["/main"]